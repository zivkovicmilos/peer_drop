package search

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/zivkovicmilos/peer_drop/rest/types"
	"github.com/zivkovicmilos/peer_drop/rest/utils"
	"github.com/zivkovicmilos/peer_drop/storage"
)

func GetSearchResults(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("input")

	// Get all identities
	foundIdentities := make([]*types.Identity, 0)
	identities, _, err := storage.GetStorageHandler().GetIdentities(utils.NoPagination, utils.DefaultSort)
	if err != nil {
		http.Error(w, "Unable to perform identity search", http.StatusInternalServerError)
		return
	}

	for _, identity := range identities {
		if strings.Contains(strings.ToLower(identity.Name), strings.ToLower(search)) {
			foundIdentities = append(foundIdentities, identity)
		}
	}

	// Get all workspaces
	foundWorkspaces := make([]*types.WorkspaceInfoWrapper, 0)
	workspaces, _, err := storage.GetStorageHandler().GetWorkspaces(utils.NoPagination)
	if err != nil {
		http.Error(w, "Unable to perform workspace search", http.StatusInternalServerError)
		return
	}

	for _, workspace := range workspaces {
		if strings.Contains(strings.ToLower(workspace.Name), strings.ToLower(search)) ||
			strings.Contains(strings.ToLower(workspace.Mnemonic), strings.ToLower(search)) {
			foundWorkspaces = append(foundWorkspaces, &types.WorkspaceInfoWrapper{
				WorkspaceMnemonic: workspace.Mnemonic,
				WorkspaceName:     workspace.Name,
			})
		}
	}

	// Get all contacts
	foundContacts := make([]*types.Contact, 0)
	contacts, _, err := storage.GetStorageHandler().GetContacts(utils.NoPagination)
	if err != nil {
		http.Error(w, "Unable to perform contacts search", http.StatusInternalServerError)
		return
	}

	for _, contact := range contacts {
		if strings.Contains(strings.ToLower(contact.Name), strings.ToLower(search)) {
			foundContacts = append(foundContacts, contact)
		}
	}

	if encodeErr := json.NewEncoder(w).Encode(&types.SearchResults{
		Identities: foundIdentities,
		Workspaces: foundWorkspaces,
		Contacts:   foundContacts,
	}); encodeErr != nil {
		http.Error(w, "Unable to encode response", http.StatusInternalServerError)
		return
	}
}
