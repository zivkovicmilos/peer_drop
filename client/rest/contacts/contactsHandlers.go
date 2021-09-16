package contacts

import (
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

var (
	dateFormat = "02.01.2006."
)

// GetContacts fetches all the contacts
func GetContacts(w http.ResponseWriter, r *http.Request) {
	limit := r.URL.Query().Get("limit")
	page := r.URL.Query().Get("page")
	paginationLimits := utils.ParsePagination(limit, page)

	contacts, totalContacts, contactsError := storage.GetStorageHandler().GetContacts(paginationLimits)
	if contactsError != nil {
		http.Error(w, "Unable to fetch contacts", http.StatusInternalServerError)
		return
	}

	encodeErr := json.NewEncoder(w).Encode(types.ContactsResponse{
		Data:  contacts,
		Count: totalContacts,
	})
	if encodeErr != nil {
		http.Error(w, "Unable to encode response", http.StatusInternalServerError)
	}
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

	email, emailError := crypto.GetEmailFromPublicKey(contact.PublicKey)
	if emailError != nil {
		http.Error(w, "Invalid public key", http.StatusBadRequest)
		return
	}

	contact.ID = uuid.New().String()
	contact.PublicKeyID = publicKey.KeyIdString()
	contact.Email = email
	contact.DateAdded = time.Now().Format(dateFormat)

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
