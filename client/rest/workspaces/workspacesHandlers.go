package workspaces

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/golang/protobuf/jsonpb"
	"github.com/gorilla/mux"
	"github.com/zivkovicmilos/peer_drop/rest/types"
	servicehandler "github.com/zivkovicmilos/peer_drop/service-handler"
	"github.com/zivkovicmilos/peer_drop/storage"
)

// CreateWorkspace creates a new workspace
func CreateWorkspace(w http.ResponseWriter, r *http.Request) {
	var workspaceRequest types.NewWorkspaceRequest

	decodeErr := json.NewDecoder(r.Body).Decode(&workspaceRequest)
	if decodeErr != nil {
		http.Error(w, "Unable to parse input", http.StatusBadRequest)
		return
	}

	// Fill in the initial workspace owner
	identity, identityErr := storage.GetStorageHandler().GetIdentityByPublicKeyID(
		workspaceRequest.BaseWorkspaceOwnerKeyID,
	)
	if identityErr != nil {
		http.Error(w, "Unknown identity", http.StatusNotFound)
		return
	}

	workspaceRequest.WorkspaceOwners = append(workspaceRequest.WorkspaceOwners, identity.PublicKey)

	// Contact the Rendezvous servers with the creation request
	clientServer := servicehandler.GetServiceHandler().GetClientServer()
	workspaceInfo, createErr := clientServer.CreateWorkspace(workspaceRequest)
	if createErr != nil {
		http.Error(w, "Unable to create workspace", http.StatusInternalServerError)
		return
	}

	createError := storage.GetStorageHandler().CreateWorkspaceInfo(workspaceInfo)
	if createError != nil {
		http.Error(w, "Unable to create workspace", http.StatusInternalServerError)
		return
	}

	if workspaceRequest.WorkspaceAccessControlType == "password" {
		password := workspaceRequest.WorkspaceAccessControl.Password
		// Save the workspace credentials
		if createErr = storage.GetStorageHandler().CreateWorkspaceCredentials(
			workspaceInfo.Mnemonic,
			nil,
			nil,
			&password,
		); createErr != nil {
			http.Error(w, "Unable to save workspace credentials", http.StatusInternalServerError)
			return
		}
	} else {
		// Save the workspace credentials
		if createErr = storage.GetStorageHandler().CreateWorkspaceCredentials(
			workspaceInfo.Mnemonic,
			&identity.PublicKey,
			&identity.PublicKey,
			nil,
		); createErr != nil {
			http.Error(w, "Unable to save workspace credentials", http.StatusInternalServerError)
			return
		}
	}

	// Initialize the workspace locally
	clientServer.TriggerWorkspaceInit(workspaceInfo)

	if encodeErr := json.NewEncoder(w).Encode(&types.NewWorkspaceResponse{Mnemonic: workspaceInfo.Mnemonic}); encodeErr != nil {
		http.Error(w, "Unable to encode response", http.StatusInternalServerError)
		return
	}
}

// GetWorkspaceInfo fetches workspace info based on the mnemonic
func GetWorkspaceInfo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	outputArr := strings.Split(params["mnemonic"], "-")
	mnemonic := strings.Join(outputArr[:], " ")

	workspaceInfo, workspaceError := storage.GetStorageHandler().GetWorkspaceInfo(mnemonic)
	if workspaceError != nil {
		http.Error(w, "Unable to fetch workspace info", http.StatusInternalServerError)
		return
	}

	if workspaceInfo == nil {
		// Not found locally, check the rendezvous
		clientServer := servicehandler.GetServiceHandler().GetClientServer()
		workspaceInfo, workspaceError = clientServer.GetWorkspaceInfo(mnemonic)
		if workspaceError != nil {
			http.Error(w, "Unable to find workspace", http.StatusInternalServerError)
			return
		}

		if workspaceInfo == nil {
			http.Error(w, "Workspace info not found", http.StatusNotFound)
			return
		}
	}

	marshaler := jsonpb.Marshaler{}
	buf := new(bytes.Buffer)
	err := marshaler.Marshal(buf, workspaceInfo)
	if err != nil {
		http.Error(w, "Unable to encode response", http.StatusInternalServerError)
	}

	_, writeErr := w.Write(buf.Bytes())
	if writeErr != nil {
		http.Error(w, "Unable to encode response", http.StatusInternalServerError)
	}
}

func JoinWorkspace(w http.ResponseWriter, r *http.Request) {
	var joinWorkspaceRequest types.JoinWorkspaceRequest

	if joinWorkspaceRequest.Password == "" && joinWorkspaceRequest.PublicKeyID == "" {
		http.Error(w, "Invalid params", http.StatusBadRequest)
		return
	}

	// Grab the workspace info
	clientServer := servicehandler.GetServiceHandler().GetClientServer()
	workspaceInfo, workspaceError := clientServer.GetWorkspaceInfo(joinWorkspaceRequest.Mnemonic)
	if workspaceError != nil {
		http.Error(w, "Unable to find workspace", http.StatusInternalServerError)
		return
	}

	if workspaceInfo == nil {
		http.Error(w, "Workspace info not found", http.StatusNotFound)
		return
	}

	confirmed := false
	if workspaceInfo.WorkspaceType == "password" {
		// Password authentication
		confirmed = clientServer.JoinWorkspacePassword(workspaceInfo, joinWorkspaceRequest.Password)
		if createErr := storage.GetStorageHandler().CreateWorkspaceInfo(workspaceInfo); createErr != nil {
			http.Error(w, "Unable to save workspace", http.StatusInternalServerError)
			return
		}

		// Save the workspace credentials
		if createErr := storage.GetStorageHandler().CreateWorkspaceCredentials(
			workspaceInfo.Mnemonic,
			nil,
			nil,
			&joinWorkspaceRequest.Password,
		); createErr != nil {
			http.Error(w, "Unable to save workspace credentials", http.StatusInternalServerError)
			return
		}
	} else {
		// Public key authentication
		identity, identityErr := storage.GetStorageHandler().GetIdentityByPublicKeyID(joinWorkspaceRequest.PublicKeyID)
		if identityErr != nil {
			http.Error(w, "Unable to find identity", http.StatusInternalServerError)
			return
		}

		if identity == nil {
			http.Error(w, "Unable to find identity", http.StatusNotFound)
			return
		}

		confirmed = clientServer.JoinWorkspacePublicKey(workspaceInfo, identity.PublicKey)

		if createErr := storage.GetStorageHandler().CreateWorkspaceInfo(workspaceInfo); createErr != nil {
			http.Error(w, "Unable to save workspace", http.StatusInternalServerError)
			return
		}

		// Save the workspace credentials
		if createErr := storage.GetStorageHandler().CreateWorkspaceCredentials(
			workspaceInfo.Mnemonic,
			&identity.PublicKey,
			&identity.PrivateKey,
			nil,
		); createErr != nil {
			http.Error(w, "Unable to save workspace credentials", http.StatusInternalServerError)
			return
		}
	}

	if confirmed {
		// Alert the client server listeners of joining
		clientServer.TriggerWorkspaceInit(workspaceInfo)

		if encodeErr := json.NewEncoder(w).Encode("Workspace joined!"); encodeErr != nil {
			http.Error(w, "Unable to encode response", http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "Unable to join workspace", http.StatusBadRequest)
	}

	return
}
