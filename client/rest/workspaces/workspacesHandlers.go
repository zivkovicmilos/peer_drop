package workspaces

import (
	"encoding/json"
	"net/http"
	"strings"

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

	if encodeErr := json.NewEncoder(w).Encode(&types.NewWorkspaceResponse{Mnemonic: workspaceInfo.Mnemonic}); encodeErr != nil {
		http.Error(w, "Unable to encode response", http.StatusInternalServerError)
		return
	}
}

// GetWorkspaceInfo fetches workspace info based on the mnemonic
func GetWorkspaceInfo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	outputArr := strings.Split(params["mnemonic"], "-")

	workspaceInfo, workspaceError := storage.GetStorageHandler().GetWorkspaceInfo(strings.Join(outputArr[:], " "))
	if workspaceError != nil {
		http.Error(w, "Unable to fetch workspace info", http.StatusInternalServerError)
		return
	}

	if workspaceInfo == nil {
		http.Error(w, "Workspace info not found", http.StatusNotFound)
		return
	}

	encodeErr := json.NewEncoder(w).Encode(workspaceInfo)
	if encodeErr != nil {
		http.Error(w, "Unable to encode response", http.StatusInternalServerError)
	}
}
