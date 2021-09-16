package storage

import (
	"bytes"
	"math/big"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"github.com/zivkovicmilos/peer_drop/rest/types"
	"github.com/zivkovicmilos/peer_drop/rest/utils"
)

// Data organization
// tableName:id:attributeName => value

var (
	delimiter = []byte(":")
)

// Entity prefixes
var (
	CONTACTS   = []byte("contacts")
	IDENTITIES = []byte("identities")
	WORKSPACES = []byte("workspaces")

	IDENTITY_PRIMARY = []byte("identityPrimary")
)

// Sub-prefixes
var (
	// CONTACTS //

	CONTACT_NAME          = []byte("name")
	CONTACT_EMAIL         = []byte("email")
	CONTACT_DATE_ADDED    = []byte("dateAdded")
	CONTACT_PUBLIC_KEY    = []byte("publicKey")
	CONTACT_PUBLIC_KEY_ID = []byte("publicKeyID")

	// IDENTITIES //
	IDENTITY_NAME           = []byte("name")
	IDENTITY_PICTURE        = []byte("picture")
	IDENTITY_DATE_CREATED   = []byte("dateCreated")
	IDENTITY_PRIVATE_KEY    = []byte("privateKey")
	IDENTITY_PUBLIC_KEY     = []byte("publicKey")
	IDENTITY_PUBLIC_KEY_ID  = []byte("publicKeyID")
	IDENTITY_IS_PRIMARY     = []byte("isPrimary")
	IDENTITY_NUM_WORKSPACES = []byte("numWorkspaces")
)

type StorageHandler struct {
	logger       hclog.Logger
	db           *leveldb.DB
	closeChannel chan struct{}
}

// SetLogger sets the storage handler's logger
func (sh *StorageHandler) SetLogger(logger hclog.Logger) {
	GetStorageHandler().logger = logger.Named("storage-handler")
}

// SetCloseChannel sets the Storage Handler's close channel
func (sh *StorageHandler) SetCloseChannel(closeChannel chan struct{}) {
	sh.closeChannel = closeChannel

	// Set the listener
	go func() {
		<-closeChannel

		sh.logger.Info("Caught stop signal...")
		closeErr := sh.Close()
		if closeErr != nil {
			sh.logger.Info("Storage interface unable to stop gracefully")
		} else {
			sh.logger.Info("Storage interface stopped gracefully")
		}
	}()
}

// OpenDB sets the storage handler's logger
func (sh *StorageHandler) OpenDB(directory string) error {
	db, err := leveldb.OpenFile(directory, nil) // TODO check of if exists
	if err != nil {
		return err
	}
	sh.db = db

	return nil
}

var storageHandlerInstance StorageHandler
var once sync.Once

// GetStorageHandler initializes the storage handler singleton
func GetStorageHandler() *StorageHandler {
	once.Do(func() {
		storageHandlerInstance = StorageHandler{}
	})

	return &storageHandlerInstance
}

// Close closes the DB connection
func (sh *StorageHandler) Close() error {
	return sh.db.Close()
}

// CONTACTS //

// GetContact fetches a single contact based on its id
func (sh *StorageHandler) GetContact(id string) (*types.Contact, error) {
	var foundContact *types.Contact
	foundContact = nil

	keyBase := append(CONTACTS, delimiter...)
	iter := sh.db.NewIterator(util.BytesPrefix(append(keyBase, []byte(id)...)), nil)
	for iter.Next() {
		if foundContact == nil {
			foundContact = &types.Contact{ID: id}
		}
		// tableName:id:attributeName => value
		keyParts := strings.Split(string(iter.Key()), ":")
		attributeName := keyParts[len(keyParts)-1]

		value := string(iter.Value())
		switch attributeName {
		case "name":
			foundContact.Name = value
		case "email":
			foundContact.Email = value
		case "dateAdded":
			foundContact.DateAdded = value
		case "publicKey":
			foundContact.PublicKey = value
		case "publicKeyID":
			foundContact.PublicKeyID = value
		}
	}

	iter.Release()
	err := iter.Error()

	return foundContact, err
}

// GetContacts fetches all contacts
func (sh *StorageHandler) GetContacts(paginationLimits utils.PaginationLimits) ([]*types.Contact, int, error) {
	var foundContacts []*types.Contact

	keyBase := append(CONTACTS, delimiter...)
	iter := sh.db.NewIterator(util.BytesPrefix(keyBase), nil)
	var currentContact *types.Contact
	for iter.Next() {
		// tableName:id:attributeName => value
		keyParts := strings.Split(string(iter.Key()), ":")

		if currentContact == nil || (currentContact != nil && currentContact.ID != keyParts[1]) {
			foundContact, findErr := sh.GetContact(keyParts[1])
			if findErr != nil {
				return nil, 0, findErr
			}

			foundContacts = append(foundContacts, foundContact)
			currentContact = foundContact
		}
	}

	iter.Release()
	err := iter.Error()

	totalContacts := len(foundContacts)

	if paginationLimits != utils.NoPagination && totalContacts != 0 {
		offset := (paginationLimits.Page - 1) * paginationLimits.Limit

		upperBound := offset + paginationLimits.Limit
		if upperBound > totalContacts {
			upperBound = totalContacts
		}

		foundContacts = foundContacts[offset:upperBound]
	} else {
		foundContacts = []*types.Contact{}
	}

	return foundContacts, totalContacts, err
}

// CreateContact stores the contact into the DB
func (sh *StorageHandler) CreateContact(contact types.Contact) error {
	fieldPairs := []struct {
		key   []byte
		value []byte
	}{
		{
			CONTACT_NAME,
			[]byte(contact.Name),
		},
		{
			CONTACT_EMAIL,
			[]byte(contact.Email),
		},
		{
			CONTACT_DATE_ADDED,
			[]byte(contact.DateAdded),
		},
		{
			CONTACT_PUBLIC_KEY,
			[]byte(contact.PublicKey),
		},
		{
			CONTACT_PUBLIC_KEY_ID,
			[]byte(contact.PublicKeyID),
		},
	}

	entityKeyBase := append(append(CONTACTS, delimiter...), append([]byte(contact.ID), delimiter...)...)
	for _, field := range fieldPairs {
		colonSeparated := append(entityKeyBase, delimiter...)
		putError := sh.db.Put(append(colonSeparated, field.key...), field.value, nil)
		if putError != nil {
			return putError
		}
	}

	return nil
}

// DeleteContact deletes a contact from the DB
func (sh *StorageHandler) DeleteContact(contact types.Contact) error {
	fieldPairs := []struct {
		key []byte
	}{
		{
			CONTACT_NAME,
		},
		{
			CONTACT_EMAIL,
		},
		{
			CONTACT_DATE_ADDED,
		},
		{
			CONTACT_PUBLIC_KEY,
		},
		{
			CONTACT_PUBLIC_KEY_ID,
		},
	}

	entityKeyBase := append(append(CONTACTS, delimiter...), append([]byte(contact.ID), delimiter...)...)
	for _, field := range fieldPairs {
		colonSeparated := append(entityKeyBase, delimiter...)
		deleteError := sh.db.Delete(append(colonSeparated, field.key...), nil)
		if deleteError != nil {
			return deleteError
		}
	}

	return nil
}

// IDENTITIES //

// GetPrimaryIdentity returns the ID of the identity that's primary
func (sh *StorageHandler) GetPrimaryIdentity() string {
	entityKeyBase := append(IDENTITY_PRIMARY, delimiter...)

	value, err := sh.db.Get(entityKeyBase, nil)
	if err != nil {
		return ""
	}

	return string(value)
}

// SetPrimaryIdentity sets the primary identity.
// Only updates the IDENTITY_PRIMARY keyset
func (sh *StorageHandler) SetPrimaryIdentity(id string) error {
	entityKeyBase := append(append(IDENTITY_PRIMARY, delimiter...), []byte(id)...)

	return sh.db.Put(entityKeyBase, []byte(id), nil)
}

// UpdateIdentityPrimary sets a new primary identity
func (sh *StorageHandler) UpdateIdentityPrimary(id string) error {
	// Unmark the previous primary identity
	previousPrimary := sh.GetPrimaryIdentity()
	if previousPrimary != "" {
		// Set its isPrimary field to false
		// TODO fix this
		entityKeyBase := append(append(IDENTITIES, delimiter...), append([]byte(previousPrimary), delimiter...)...)
		entityKeyBase = append(entityKeyBase, IDENTITY_IS_PRIMARY...)

		putError := sh.db.Put(entityKeyBase, []byte{0}, nil)
		if putError != nil {
			return putError
		}
	}

	// Update the IDENTITY_PRIMARY index
	if identityPrimaryError := sh.SetPrimaryIdentity(id); identityPrimaryError != nil {
		return identityPrimaryError
	}

	// Update the main DB instance field
	entityKeyBase := append(append(IDENTITIES, delimiter...), append([]byte(id), delimiter...)...)
	entityKeyBase = append(entityKeyBase, IDENTITY_IS_PRIMARY...)

	return sh.db.Put(entityKeyBase, []byte{1}, nil)
}

// CreateIdentity stores the identity into the DB
func (sh *StorageHandler) CreateIdentity(identity types.Identity) error {
	var isPrimaryValue []byte
	if identity.IsPrimary {
		isPrimaryValue = []byte{1}
	} else {
		isPrimaryValue = []byte{0}
	}

	primaryIdentityValue := sh.GetPrimaryIdentity()
	if primaryIdentityValue == "" {
		// No primary identity set, this is the first identity, so it should be the primary one
		isPrimaryValue = []byte{1}
		if setPrimaryErr := sh.SetPrimaryIdentity(identity.ID); setPrimaryErr != nil {
			return setPrimaryErr
		}
	}

	fieldPairs := []struct {
		key   []byte
		value []byte
	}{
		{
			IDENTITY_NAME,
			[]byte(identity.Name),
		},
		{
			IDENTITY_PICTURE,
			[]byte(identity.Picture),
		},
		{
			IDENTITY_DATE_CREATED,
			[]byte(identity.DateCreated),
		},
		{
			IDENTITY_PRIVATE_KEY,
			[]byte(identity.PrivateKey),
		},
		{
			IDENTITY_PUBLIC_KEY,
			[]byte(identity.PublicKey),
		},
		{
			IDENTITY_PUBLIC_KEY_ID,
			[]byte(identity.PublicKeyID),
		},
		{
			IDENTITY_IS_PRIMARY,
			isPrimaryValue,
		},
		{
			IDENTITY_NUM_WORKSPACES,
			big.NewInt(int64(identity.NumWorkspaces)).Bytes(),
		},
	}

	entityKeyBase := append(append(IDENTITIES, delimiter...), append([]byte(identity.ID), delimiter...)...)
	for _, field := range fieldPairs {
		colonSeparated := append(entityKeyBase, delimiter...)
		putError := sh.db.Put(append(colonSeparated, field.key...), field.value, nil)
		if putError != nil {
			return putError
		}
	}

	return nil
}

// DeleteIdentity deletes the identity from the DB
func (sh *StorageHandler) DeleteIdentity(identity types.Identity) error {

	fieldPairs := []struct {
		key []byte
	}{
		{
			IDENTITY_NAME,
		},
		{
			IDENTITY_PICTURE,
		},
		{
			IDENTITY_DATE_CREATED,
		},
		{
			IDENTITY_PRIVATE_KEY,
		},
		{
			IDENTITY_PUBLIC_KEY,
		},
		{
			IDENTITY_PUBLIC_KEY_ID,
		},
		{
			IDENTITY_IS_PRIMARY,
		},
		{
			IDENTITY_NUM_WORKSPACES,
		},
	}

	entityKeyBase := append(append(IDENTITIES, delimiter...), append([]byte(identity.ID), delimiter...)...)
	for _, field := range fieldPairs {
		colonSeparated := append(entityKeyBase, delimiter...)
		deleteError := sh.db.Delete(append(colonSeparated, field.key...), nil)
		if deleteError != nil {
			return deleteError
		}
	}

	// TODO Update the primary identity if primary is deleted

	return nil
}

// GetIdentity fetches a single identity based on its id
func (sh *StorageHandler) GetIdentity(id string) (*types.Identity, error) {
	var foundIdentity *types.Identity
	foundIdentity = nil

	keyBase := append(IDENTITIES, delimiter...)
	iter := sh.db.NewIterator(util.BytesPrefix(append(keyBase, []byte(id)...)), nil)
	for iter.Next() {
		if foundIdentity == nil {
			foundIdentity = &types.Identity{ID: id}
		}
		// tableName:id:attributeName => value
		keyParts := strings.Split(string(iter.Key()), ":")
		attributeName := keyParts[len(keyParts)-1]

		value := string(iter.Value())
		switch attributeName {
		case "name":
			foundIdentity.Name = value
		case "picture":
			foundIdentity.Picture = value
		case "dateCreated":
			foundIdentity.DateCreated = value
		case "privateKey":
			foundIdentity.PrivateKey = value
		case "publicKey":
			foundIdentity.PublicKey = value
		case "publicKeyID":
			foundIdentity.PublicKeyID = value
		case "isPrimary":
			var isPrimaryValue bool
			if bytes.Compare(iter.Value(), []byte{1}) == 0 {
				isPrimaryValue = true
			} else {
				isPrimaryValue = false
			}
			foundIdentity.IsPrimary = isPrimaryValue
		case "numWorkspaces":
			foundIdentity.NumWorkspaces = int(big.NewInt(0).SetBytes(iter.Value()).Int64())
		}
	}

	iter.Release()
	err := iter.Error()

	return foundIdentity, err
}

// GetIdentities fetches all identities
func (sh *StorageHandler) GetIdentities(
	paginationLimits utils.PaginationLimits,
	sortParams utils.SortParams,
) ([]*types.Identity, int, error) {
	var foundIdentities []*types.Identity

	keyBase := append(IDENTITIES, delimiter...)
	iter := sh.db.NewIterator(util.BytesPrefix(keyBase), nil)
	var currentIdentity *types.Identity
	for iter.Next() {
		// tableName:id:attributeName => value
		keyParts := strings.Split(string(iter.Key()), ":")

		if currentIdentity == nil || (currentIdentity != nil && currentIdentity.ID != keyParts[1]) {
			foundIdentity, findErr := sh.GetIdentity(keyParts[1])
			if findErr != nil {
				return nil, 0, findErr
			}

			foundIdentities = append(foundIdentities, foundIdentity)
			currentIdentity = foundIdentity
		}
	}

	iter.Release()
	err := iter.Error()

	totalIdentities := len(foundIdentities)

	// Enforce sort
	switch sortParams.SortParam {
	case utils.SORT_NAME:
		sort.Slice(foundIdentities, func(i, j int) bool {
			comparison := strings.Compare(foundIdentities[i].Name, foundIdentities[j].Name)

			if sortParams.SortDirection == utils.SORT_DIR_ASC {
				return comparison <= 0
			} else {
				return comparison > 0
			}
		})
	case utils.SORT_PUBLIC_KEY_ID:
		sort.Slice(foundIdentities, func(i, j int) bool {
			comparison := strings.Compare(foundIdentities[i].PublicKeyID, foundIdentities[j].PublicKeyID)

			if sortParams.SortDirection == utils.SORT_DIR_ASC {
				return comparison <= 0
			} else {
				return comparison > 0
			}
		})
	case utils.SORT_NUM_WORKSPACES:
		sort.Slice(foundIdentities, func(i, j int) bool {
			if sortParams.SortDirection == utils.SORT_DIR_ASC {
				return foundIdentities[i].NumWorkspaces <= foundIdentities[j].NumWorkspaces
			} else {
				return foundIdentities[i].NumWorkspaces > foundIdentities[j].NumWorkspaces
			}
		})
	case utils.SORT_DATE_CREATED:
		sort.Slice(foundIdentities, func(i, j int) bool {
			dateI, _ := time.Parse(utils.DateFormat, foundIdentities[i].DateCreated)
			dateJ, _ := time.Parse(utils.DateFormat, foundIdentities[j].DateCreated)
			if sortParams.SortDirection == utils.SORT_DIR_ASC {
				return dateI.Before(dateJ)
			} else {
				return dateI.After(dateJ)
			}
		})
	}

	// Enforce pagination
	if paginationLimits != utils.NoPagination && totalIdentities != 0 {
		offset := (paginationLimits.Page - 1) * paginationLimits.Limit

		upperBound := offset + paginationLimits.Limit
		if upperBound > totalIdentities {
			upperBound = totalIdentities
		}

		foundIdentities = foundIdentities[offset:upperBound]
	} else if totalIdentities == 0 {
		foundIdentities = []*types.Identity{}
	}

	return foundIdentities, totalIdentities, err
}
