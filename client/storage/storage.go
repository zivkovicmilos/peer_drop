package storage

import (
	"strings"
	"sync"

	"github.com/hashicorp/go-hclog"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"github.com/zivkovicmilos/peer_drop/rest/types"
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
)

// Sub-prefixes
var (
	// CONTACTS //

	CONTACT_NAME          = []byte("name")
	CONTACT_EMAIL         = []byte("email")
	CONTACT_DATE_ADDED    = []byte("dateAdded")
	CONTACT_PUBLIC_KEY    = []byte("publicKey")
	CONTACT_PUBLIC_KEY_ID = []byte("publicKeyID")
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
