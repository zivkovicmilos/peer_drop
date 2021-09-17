package identities

import (
	"crypto/rsa"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/zivkovicmilos/peer_drop/crypto"
	"github.com/zivkovicmilos/peer_drop/rest/types"
	"github.com/zivkovicmilos/peer_drop/rest/utils"
	"github.com/zivkovicmilos/peer_drop/storage"
)

// CreateIdentity creates a single identity
func CreateIdentity(w http.ResponseWriter, r *http.Request) {
	var identity types.NewIdentityRequest
	newIdentity := &types.Identity{}

	decodeErr := json.NewDecoder(r.Body).Decode(&identity)
	if decodeErr != nil {
		http.Error(w, "Unable to parse input", http.StatusBadRequest)
		return
	}

	// Parse the private key
	privateKey, parseError := crypto.ParsePrivateKeyFromPemStr(identity.PrivateKey)
	if parseError != nil {
		http.Error(w, "Invalid private key", http.StatusBadRequest)
		return
	}

	publicKeyEncoded, encodeErr := crypto.EncodePublicKey(privateKey.PublicKey.PublicKey.(*rsa.PublicKey))
	if encodeErr != nil {
		http.Error(w, "Unable to encode public key", http.StatusInternalServerError)
		return
	}

	newIdentity.ID = uuid.New().String()
	newIdentity.Name = identity.Name
	newIdentity.Picture = identity.Picture
	newIdentity.DateCreated = time.Now().Format(utils.DateFormat)
	newIdentity.NumWorkspaces = 0
	newIdentity.PrivateKey = identity.PrivateKey
	newIdentity.PublicKey = publicKeyEncoded
	newIdentity.PublicKeyID = privateKey.PublicKey.KeyIdString()
	newIdentity.IsPrimary = false

	createError := storage.GetStorageHandler().CreateIdentity(*newIdentity)
	if createError != nil {
		http.Error(w, "Unable to create identity", http.StatusInternalServerError)
		return
	}

	if encodeErr := json.NewEncoder(w).Encode(newIdentity); encodeErr != nil {
		http.Error(w, "Unable to encode response", http.StatusInternalServerError)
		return
	}
}

// GetIdentity fetches a single identity
func GetIdentity(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	identity, identityError := storage.GetStorageHandler().GetIdentity(params["identityId"])
	if identityError != nil {
		http.Error(w, "Unable to fetch identity", http.StatusInternalServerError)
		return
	}

	if identity == nil {
		http.Error(w, "Identity not found", http.StatusNotFound)
		return
	}

	encodeErr := json.NewEncoder(w).Encode(identity)
	if encodeErr != nil {
		http.Error(w, "Unable to encode response", http.StatusInternalServerError)
	}
}

// GetIdentities fetches all the identities
func GetIdentities(w http.ResponseWriter, r *http.Request) {
	limit := r.URL.Query().Get("limit")
	page := r.URL.Query().Get("page")
	paginationLimits := utils.ParsePagination(limit, page)

	sortParam := r.URL.Query().Get("sortParam")
	sortDirection := r.URL.Query().Get("sortDirection")
	sortParams := utils.ParseSortParams(sortParam, sortDirection)

	identities, totalIdentities, identitiesError := storage.GetStorageHandler().GetIdentities(paginationLimits, sortParams)
	if identitiesError != nil {
		http.Error(w, "Unable to fetch identities", http.StatusInternalServerError)
		return
	}

	encodeErr := json.NewEncoder(w).Encode(types.IdentitiesResponse{
		Data:  identities,
		Count: totalIdentities,
	})
	if encodeErr != nil {
		http.Error(w, "Unable to encode response", http.StatusInternalServerError)
	}
}

// UpdateIdentity updates a single identity
func UpdateIdentity(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	// Check to see if the identity exists
	identity, identityError := storage.GetStorageHandler().GetIdentity(params["identityId"])
	if identityError != nil {
		http.Error(w, "Unable to fetch identity", http.StatusInternalServerError)
		return
	}

	if identity == nil {
		http.Error(w, "Identity not found", http.StatusNotFound)
		return
	}

	// Grab the updated identity

	// Exists because it needs to have the private key encoded
	var updateRequest types.IdentityUpdateRequest

	var updatedIdentity types.Identity

	decodeErr := json.NewDecoder(r.Body).Decode(&updateRequest)
	if decodeErr != nil {
		http.Error(w, "Unable to parse input", http.StatusBadRequest)
		return
	}

	if updateRequest.PrivateKey != "" {
		// The private key has been updated.

		// Parse the private key
		privateKey, parseError := crypto.ParsePrivateKeyFromPemStr(updateRequest.PrivateKey)
		if parseError != nil {
			http.Error(w, "Invalid private key", http.StatusBadRequest)
			return
		}

		publicKeyEncoded, encodeErr := crypto.EncodePublicKey(privateKey.PublicKey.PublicKey.(*rsa.PublicKey))
		if encodeErr != nil {
			http.Error(w, "Unable to encode public key", http.StatusInternalServerError)
			return
		}

		updatedIdentity.PublicKey = publicKeyEncoded
		updatedIdentity.PublicKeyID = privateKey.PublicKey.KeyIdString()
	} else {
		// No private key update, use the old one
		updatedIdentity.PublicKey = identity.PublicKey
		updatedIdentity.PrivateKey = identity.PrivateKey
		updatedIdentity.PublicKeyID = identity.PublicKeyID
	}

	// Copy the existing fields and populate new ones
	updatedIdentity.ID = identity.ID
	updatedIdentity.Picture = updateRequest.Picture
	updatedIdentity.Name = updateRequest.Name
	updatedIdentity.DateCreated = identity.DateCreated
	updatedIdentity.IsPrimary = identity.IsPrimary
	updatedIdentity.NumWorkspaces = identity.NumWorkspaces

	updateError := storage.GetStorageHandler().CreateIdentity(updatedIdentity)
	if updateError != nil {
		http.Error(w, "Unable to update identity", http.StatusInternalServerError)
		return
	}

	encodeErr := json.NewEncoder(w).Encode("Identity updated")
	if encodeErr != nil {
		http.Error(w, "Unable to encode response", http.StatusInternalServerError)
	}
}

// SetAsPrimary updates a single identity to be the primary one
func SetAsPrimary(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	// Check to see if the identity exists
	identity, identityError := storage.GetStorageHandler().GetIdentity(params["identityId"])
	if identityError != nil {
		http.Error(w, "Unable to fetch identity", http.StatusInternalServerError)
		return
	}

	if identity == nil {
		http.Error(w, "Identity not found", http.StatusNotFound)
		return
	}

	updateErr := storage.GetStorageHandler().UpdateIdentityPrimary(identity.ID)
	if updateErr != nil {
		http.Error(w, "Unable to update identity", http.StatusInternalServerError)
		return
	}

	encodeErr := json.NewEncoder(w).Encode("Identity updated")
	if encodeErr != nil {
		http.Error(w, "Unable to encode response", http.StatusInternalServerError)
	}
}

// GetPublicKey fetches the identity's public key
func GetPublicKey(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	identity, identityError := storage.GetStorageHandler().GetIdentity(params["identityId"])
	if identityError != nil {
		http.Error(w, "Unable to fetch identity", http.StatusInternalServerError)
		return
	}

	if identity == nil {
		http.Error(w, "Identity not found", http.StatusNotFound)
		return
	}

	encodeErr := json.NewEncoder(w).Encode(types.IdentityPublicKeyResponse{PublicKey: identity.PublicKey})
	if encodeErr != nil {
		http.Error(w, "Unable to encode response", http.StatusInternalServerError)
	}
}

// GetPrivateKey fetches the identity's private key
func GetPrivateKey(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	identity, identityError := storage.GetStorageHandler().GetIdentity(params["identityId"])
	if identityError != nil {
		http.Error(w, "Unable to fetch identity", http.StatusInternalServerError)
		return
	}

	if identity == nil {
		http.Error(w, "Identity not found", http.StatusNotFound)
		return
	}

	encodeErr := json.NewEncoder(w).Encode(types.IdentityPrivateKeyResponse{PrivateKey: identity.PrivateKey})
	if encodeErr != nil {
		http.Error(w, "Unable to encode response", http.StatusInternalServerError)
	}
}

// DeleteIdentity deletes a single identity
func DeleteIdentity(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	// TODO add a constraint that at least 1 identity needs to exist

	// Check to see if the identity exists
	identity, identityError := storage.GetStorageHandler().GetIdentity(params["identityId"])
	if identityError != nil {
		http.Error(w, "Unable to fetch identity", http.StatusInternalServerError)
		return
	}

	if identity == nil {
		http.Error(w, "Identity not found", http.StatusNotFound)
		return
	}

	deleteError := storage.GetStorageHandler().DeleteIdentity(*identity)
	if deleteError != nil {
		http.Error(w, "Unable to delete identity", http.StatusInternalServerError)
		return
	}

	encodeErr := json.NewEncoder(w).Encode("Identity deleted")
	if encodeErr != nil {
		http.Error(w, "Unable to encode response", http.StatusInternalServerError)
	}
}

// GetPrimaryIdentity fetches the primary identity
func GetPrimaryIdentity(w http.ResponseWriter, r *http.Request) {
	identityID := storage.GetStorageHandler().GetPrimaryIdentity()

	if identityID == "" {
		http.Error(w, "Primary identity not found", http.StatusNotFound)
		return
	}

	identity, identityError := storage.GetStorageHandler().GetIdentity(identityID)
	if identityError != nil {
		http.Error(w, "Unable to fetch identity", http.StatusInternalServerError)
		return
	}

	if identity == nil {
		http.Error(w, "Identity not found", http.StatusNotFound)
		return
	}

	encodeErr := json.NewEncoder(w).Encode(identity)
	if encodeErr != nil {
		http.Error(w, "Unable to encode response", http.StatusInternalServerError)
	}
}
