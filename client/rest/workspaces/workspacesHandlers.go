package workspaces

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/golang/protobuf/jsonpb"
	"github.com/gorilla/mux"
	"github.com/zivkovicmilos/peer_drop/proto"
	"github.com/zivkovicmilos/peer_drop/rest/types"
	"github.com/zivkovicmilos/peer_drop/rest/utils"
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

	if workspaceInfo.SecurityType == "password" {
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

func formatWorkspaceListResponse(workspaceInfos []*proto.WorkspaceInfo, total int) types.WorkspaceListResponse {
	workspaceWrappers := make([]types.WorkspaceInfoWrapper, 0)
	for _, workspaceInfo := range workspaceInfos {
		workspaceWrappers = append(workspaceWrappers, types.WorkspaceInfoWrapper{
			WorkspaceMnemonic: workspaceInfo.Mnemonic,
			WorkspaceName:     workspaceInfo.Name,
		})
	}

	return types.WorkspaceListResponse{
		WorkspaceWrappers: workspaceWrappers,
		Count:             total,
	}
}

// GetWorkspaces fetches all the workspaces
func GetWorkspaces(w http.ResponseWriter, r *http.Request) {
	limit := r.URL.Query().Get("limit")
	page := r.URL.Query().Get("page")
	paginationLimits := utils.ParsePagination(limit, page)

	workspaces, totalWorkspaces, workspacesError := storage.GetStorageHandler().GetWorkspaces(paginationLimits)
	if workspacesError != nil {
		http.Error(w, "Unable to fetch workspaces", http.StatusInternalServerError)
		return
	}

	if encodeErr := json.NewEncoder(w).Encode(formatWorkspaceListResponse(workspaces, totalWorkspaces)); encodeErr != nil {
		http.Error(w, "Unable to encode response", http.StatusInternalServerError)
		return
	}
}

// formatFileList is a helper function for formatting the retrieved file list
func formatFileList(fileList []*proto.File) []types.FileInfo {
	responseList := make([]types.FileInfo, 0)

	for _, file := range fileList {
		responseList = append(responseList, types.FileInfo{
			Name:         file.Name,
			Extension:    file.Extension,
			Size:         file.Size,
			DateModified: file.DateModified,
			Checksum:     file.FileChecksum,
		})
	}

	return responseList
}

// GetWorkspaceFiles fetches all available files for download
func GetWorkspaceFiles(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	outputArr := strings.Split(params["mnemonic"], "-")
	mnemonic := strings.Join(outputArr[:], " ")

	// Check if we know this workspace
	workspaceInfo, workspaceError := storage.GetStorageHandler().GetWorkspaceInfo(mnemonic)
	if workspaceError != nil {
		http.Error(w, "Unable to fetch workspace info", http.StatusInternalServerError)
		return
	}

	if workspaceInfo == nil {
		http.Error(w, "Workspace not found", http.StatusNotFound)
		return
	}

	// Grab the file list from the clientServer
	clientServer := servicehandler.GetServiceHandler().GetClientServer()
	fileList := clientServer.GetFileList(mnemonic)

	formattedList := formatFileList(fileList)

	detailedResponse := &types.WorkspaceDetailedResponse{
		WorkspaceMnemonic: workspaceInfo.Mnemonic,
		WorkspaceName:     workspaceInfo.Name,
		WorkspaceType:     workspaceInfo.WorkspaceType,
		WorkspaceFiles:    formattedList,
	}

	if encodeErr := json.NewEncoder(w).Encode(detailedResponse); encodeErr != nil {
		http.Error(w, "Unable to encode response", http.StatusInternalServerError)
		return
	}
}

// GetWorkspaceNumPeers returns the number of connected workspace peers
func GetWorkspaceNumPeers(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	outputArr := strings.Split(params["mnemonic"], "-")
	mnemonic := strings.Join(outputArr[:], " ")

	// Check if we know this workspace
	workspaceInfo, workspaceError := storage.GetStorageHandler().GetWorkspaceInfo(mnemonic)
	if workspaceError != nil {
		http.Error(w, "Unable to fetch workspace info", http.StatusInternalServerError)
		return
	}

	if workspaceInfo == nil {
		http.Error(w, "Workspace not found", http.StatusNotFound)
		return
	}

	// Grab the number of peers from the clientServer
	clientServer := servicehandler.GetServiceHandler().GetClientServer()
	numPeers := clientServer.GetNumberOfPeers(mnemonic)

	if encodeErr := json.NewEncoder(w).Encode(&types.WorkspacePeersResponse{NumPeers: numPeers}); encodeErr != nil {
		http.Error(w, "Unable to encode response", http.StatusInternalServerError)
		return
	}
}

func JoinWorkspace(w http.ResponseWriter, r *http.Request) {
	var joinWorkspaceRequest types.JoinWorkspaceRequest

	decodeErr := json.NewDecoder(r.Body).Decode(&joinWorkspaceRequest)
	if decodeErr != nil {
		http.Error(w, "Unable to parse input", http.StatusBadRequest)
		return
	}

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
	if workspaceInfo.SecurityType == "password" {
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

// AddFileToWorkspace uploads a new file to the workspace
func AddFileToWorkspace(w http.ResponseWriter, r *http.Request) {
	// Max size for now is 1 GB
	if parseErr := r.ParseMultipartForm(1 * 1024 * 1024 * 1024); parseErr != nil {
		http.Error(w, "Unable to parse file", http.StatusBadRequest)
		return
	}

	formMnemonic := r.FormValue("mnemonic")
	outputArr := strings.Split(formMnemonic, "-")
	mnemonic := strings.Join(outputArr[:], " ")

	formFile, handler, err := r.FormFile("workspaceFile")
	if err != nil {
		http.Error(w, "Unable to find form file", http.StatusBadRequest)
		return
	}
	defer formFile.Close()

	// Grab the save directory from the client server
	clientServer := servicehandler.GetServiceHandler().GetClientServer()
	saveDirectory, findErr := clientServer.GetWorkspaceSaveDir(mnemonic)
	if findErr != nil {
		http.Error(w, "Unknown workspace", http.StatusInternalServerError)
		return
	}

	// Create the file in the save directory
	saveFile, createErr := os.Create(
		fmt.Sprintf("%s/%s", saveDirectory, handler.Filename),
	)
	if createErr != nil {
		http.Error(w, "Unable to create file", http.StatusInternalServerError)
		return
	}

	defer saveFile.Close()

	// Read from the uploaded file
	fileBytes, err := ioutil.ReadAll(formFile)
	if err != nil {
		http.Error(w, "Unable to read file", http.StatusInternalServerError)
		return
	}

	// write this byte array to the created file
	_, writeErr := saveFile.Write(fileBytes)
	if writeErr != nil {
		http.Error(w, "Unable to save file", http.StatusInternalServerError)
		return
	}

	if encodeErr := json.NewEncoder(w).Encode("Workspace file uploaded!"); encodeErr != nil {
		http.Error(w, "Unable to encode response", http.StatusInternalServerError)
		return
	}
}

func DownloadWorkspaceFile(w http.ResponseWriter, r *http.Request) {
	var downloadFileRequest types.FileDownloadRequest

	decodeErr := json.NewDecoder(r.Body).Decode(&downloadFileRequest)
	if decodeErr != nil {
		http.Error(w, "Unable to parse input", http.StatusBadRequest)
		return
	}

	clientServer := servicehandler.GetServiceHandler().GetClientServer()
	downloadInfo, downloadErr := clientServer.HandleFileDownload(
		downloadFileRequest.WorkspaceMnemonic,
		downloadFileRequest.FileChecksum,
	)

	if downloadErr != nil {
		http.Error(w, "Unable to download file", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+strconv.Quote(downloadInfo.FileName))
	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeFile(w, r, downloadInfo.FilePath)
}
