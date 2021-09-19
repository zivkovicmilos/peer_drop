package workspaces

import (
	"encoding/json"
	"net/http"

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
