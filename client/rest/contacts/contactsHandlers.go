package contacts

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/zivkovicmilos/peer_drop/crypto"
	"github.com/zivkovicmilos/peer_drop/rest/types"
	"github.com/zivkovicmilos/peer_drop/storage"
)

type ContactsHandler struct {
}

// GetContacts fetches all the contacts
func GetContacts(w http.ResponseWriter, r *http.Request) {

}

// GetContact fetches a single contact
func GetContact(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	contact, contactError := storage.GetStorageHandler().GetContact(params["contactId"])
	if contactError != nil {
		http.Error(w, "Unable to fetch contact", http.StatusInternalServerError)
		return
	}

	if contact == nil {
		http.Error(w, "Contact not found", http.StatusNotFound)
		return
	}

	encodeErr := json.NewEncoder(w).Encode(contact)
	if encodeErr != nil {
		http.Error(w, "Unable to encode response", http.StatusInternalServerError)
	}
}

// UpdateContact updates a single contact
func UpdateContact(w http.ResponseWriter, r *http.Request) {

}

// CreateContact creates a single contact
func CreateContact(w http.ResponseWriter, r *http.Request) {
	var contact types.Contact

	decodeErr := json.NewDecoder(r.Body).Decode(&contact)
	if decodeErr != nil {
		http.Error(w, "Unable to parse input", http.StatusBadRequest)
		return
	}

	// Parse the public key
	publicKey, parseError := crypto.ParseRsaPublicKeyFromPemStr(contact.PublicKey)
	if parseError != nil {
		http.Error(w, "Invalid public key", http.StatusBadRequest)
		return
	}

	contact.ID = uuid.New().String()
	contact.PublicKeyID = crypto.GetKeyID(publicKey.N.Bytes())

	storageHandler := storage.GetStorageHandler()
	createError := storageHandler.CreateContact(contact)
	if createError != nil {
		http.Error(w, "Unable to create contact", http.StatusInternalServerError)
		return
	}

	if encodeErr := json.NewEncoder(w).Encode(contact); encodeErr != nil {
		http.Error(w, "Unable to encode response", http.StatusInternalServerError)
		return
	}
}

// DeleteContact deletes a single contact
func DeleteContact(w http.ResponseWriter, r *http.Request) {

}
