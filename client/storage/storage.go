package storage

import (
	"bytes"
	"fmt"
	"math/big"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"github.com/zivkovicmilos/peer_drop/proto"
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
	WORKSPACE_INFO   = []byte("workspaceInfo")
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

	// RENDEZVOUS NODES //
	WORKSPACE_INFO_MNEMONIC               = []byte("mnemonic")
	WORKSPACE_INFO_WORKSPACE_OWNER        = []byte("workspaceOwner")
	WORKSPACE_INFO_SECURITY_TYPE          = []byte("securityType")
	WORKSPACE_INFO_NAME                   = []byte("name")
	WORKSPACE_INFO_PASSWORD_HASH          = []byte("passwordHash")
	WORKSPACE_INFO_CONTACT                = []byte("contact")
	WORKSPACE_INFO_WORKSPACE_OWNER_LIBP2P = []byte("libp2pAddress")
	WORKSPACE_INFO_WORKSPACE_PUBLIC_KEY   = []byte("publicKey")
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
	db, err := leveldb.OpenFile(directory, nil)
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
		putError := sh.db.Put(append(entityKeyBase, field.key...), field.value, nil)
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
		deleteError := sh.db.Delete(append(entityKeyBase, field.key...), nil)
		if deleteError != nil {
			return deleteError
		}
	}

	return nil
}

// IDENTITIES //

// GetPrimaryIdentity returns the ID of the identity that's primary
func (sh *StorageHandler) GetPrimaryIdentity() string {
	value, err := sh.db.Get(IDENTITY_PRIMARY, nil)
	if err != nil {
		return ""
	}

	return string(value)
}

// SetPrimaryIdentity sets the primary identity.
// Only updates the IDENTITY_PRIMARY keyset
func (sh *StorageHandler) SetPrimaryIdentity(id string) error {
	return sh.db.Put(IDENTITY_PRIMARY, []byte(id), nil)
}

// UpdateIdentityPrimary sets a new primary identity
func (sh *StorageHandler) UpdateIdentityPrimary(id string) error {
	// Unmark the previous primary identity
	previousPrimary := sh.GetPrimaryIdentity()
	if previousPrimary != "" {
		// Set its isPrimary field to false
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
		putError := sh.db.Put(append(entityKeyBase, field.key...), field.value, nil)
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
		deleteError := sh.db.Delete(append(entityKeyBase, field.key...), nil)
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

// RENDEZVOUS //

func (sh *StorageHandler) GetWorkspaceInfo(mnemonic string) (*proto.WorkspaceInfo, error) {
	var foundWorkspaceInfo *proto.WorkspaceInfo
	foundWorkspaceInfo = nil

	workspaceOwnerMap := make(map[int]*proto.WorkspaceOwner)
	workspaceContactPKMap := make(map[int]string)

	getWorkspaceOwner := func(index int) *proto.WorkspaceOwner {
		workspaceOwner, ok := workspaceOwnerMap[index]
		if !ok {
			return &proto.WorkspaceOwner{}
		}

		return workspaceOwner
	}

	keyBase := append(WORKSPACE_INFO, delimiter...)
	iter := sh.db.NewIterator(util.BytesPrefix(append(keyBase, []byte(mnemonic)...)), nil)
	for iter.Next() {
		if foundWorkspaceInfo == nil {
			foundWorkspaceInfo = &proto.WorkspaceInfo{Mnemonic: mnemonic}
		}
		// workspaceInfo:<mnemonic>:workspaceOwner:<index>:attributeName => value
		// workspaceInfo:<mnemonic>:contact:<index>:attributeName => value
		keyParts := strings.Split(string(iter.Key()), ":")
		attributeName := keyParts[len(keyParts)-1]

		value := string(iter.Value())
		// These attributes can be extracted in the form
		// workspaceInfo:mnemonic:attributeName => value
		switch attributeName {
		case "name":
			foundWorkspaceInfo.Name = value
		case "securityType":
			foundWorkspaceInfo.SecurityType = value
		case "passwordHash":
			foundWorkspaceInfo.SecuritySettings = &proto.WorkspaceInfo_PasswordHash{PasswordHash: value}
		case "publicKey":
			bigIndex := big.NewInt(0).SetBytes([]byte(keyParts[len(keyParts)-2])) // next to last key value
			intIndex := int(bigIndex.Int64())
			if string(keyParts[len(keyParts)-3]) == "contact" {
				workspaceContactPKMap[intIndex] = value
			} else {
				ref := getWorkspaceOwner(intIndex)
				ref.PublicKey = value
				workspaceOwnerMap[intIndex] = ref
			}
		case "libp2pAddress":
			bigIndex := big.NewInt(0).SetBytes([]byte(keyParts[len(keyParts)-2])) // next to last key value
			intIndex := int(bigIndex.Int64())
			ref := getWorkspaceOwner(intIndex)
			ref.Libp2PAddress = value
			workspaceOwnerMap[intIndex] = ref
		}
	}

	iter.Release()
	err := iter.Error()

	// Add the owners to the object
	foundWorkspaceInfo.WorkspaceOwners = make([]*proto.WorkspaceOwner, 0)
	for _, workspaceOwner := range workspaceOwnerMap {
		foundWorkspaceInfo.WorkspaceOwners = append(foundWorkspaceInfo.WorkspaceOwners, workspaceOwner)
	}

	// Add the contacts if this is the correct security type
	if foundWorkspaceInfo.SecurityType == "contacts" {
		contacts := make([]string, 0)
		for _, contactPK := range workspaceContactPKMap {
			contacts = append(contacts, contactPK)
		}

		foundWorkspaceInfo.SecuritySettings = &proto.WorkspaceInfo_ContactsWrapper{
			ContactsWrapper: &proto.ContactsWrapper{
				ContactPublicKeys: contacts,
			},
		}
	}

	return foundWorkspaceInfo, err
}

// CreateWorkspaceInfo stores the workspace info into the rendezvous DB
func (sh *StorageHandler) CreateWorkspaceInfo(workspaceInfo *proto.WorkspaceInfo) error {
	fieldPairs := []struct {
		key   []byte
		value []byte
	}{
		{
			WORKSPACE_INFO_NAME,
			[]byte(workspaceInfo.Name),
		},
		{
			WORKSPACE_INFO_SECURITY_TYPE,
			[]byte(workspaceInfo.SecurityType),
		},
	}

	// Set the base fields
	entityKeyBase := append(append(WORKSPACE_INFO, delimiter...), append([]byte(workspaceInfo.Mnemonic), delimiter...)...)
	for _, field := range fieldPairs {
		putError := sh.db.Put(append(entityKeyBase, field.key...), field.value, nil)
		if putError != nil {
			return putError
		}
	}

	// Set the workspace owners
	if len(workspaceInfo.WorkspaceOwners) > 0 {
		workspaceOwnerIndex := big.NewInt(1)
		indexFound := false
		for !indexFound {
			if workspaceOwnerIndex.Int64() > 10 {
				return fmt.Errorf("invalid number of workspace owners")
			}

			// Search index is
			// workspaceInfo:<mnemonic>:workspaceOwner:<index>
			searchIndex := append(
				append(entityKeyBase, WORKSPACE_INFO_WORKSPACE_OWNER...), append(delimiter, workspaceOwnerIndex.Bytes()...)...,
			)
			foundInfo, _ := sh.db.Get(searchIndex, nil)
			if foundInfo == nil || string(foundInfo) == "" {
				indexFound = true
			} else {
				workspaceOwnerIndex.Add(workspaceOwnerIndex, big.NewInt(1))
			}
		}

		// workspaceInfo:<mnemonic>:workspaceOwner:<index>:attributeName => value
		workspaceOwnerKeybase := append(entityKeyBase, WORKSPACE_INFO_WORKSPACE_OWNER...)
		workspaceOwnerKeybase = append(workspaceOwnerKeybase, delimiter...)
		for _, workspaceOwner := range workspaceInfo.WorkspaceOwners {
			key := append(workspaceOwnerKeybase, []byte(workspaceOwnerIndex.String())...)
			key = append(key, delimiter...)

			// libp2p address
			libp2pAddressKey := append(key, WORKSPACE_INFO_WORKSPACE_OWNER_LIBP2P...)
			putError := sh.db.Put(libp2pAddressKey, []byte(workspaceOwner.Libp2PAddress), nil)
			if putError != nil {
				return putError
			}

			// public key
			publicKeyKey := append(key, WORKSPACE_INFO_WORKSPACE_PUBLIC_KEY...)
			putError = sh.db.Put(publicKeyKey, []byte(workspaceOwner.PublicKey), nil)
			if putError != nil {
				return putError
			}

			workspaceOwnerIndex.Add(workspaceOwnerIndex, big.NewInt(1))
		}
	}

	if workspaceInfo.SecurityType == "password" {
		// Set the password hash
		settings := workspaceInfo.SecuritySettings.(*proto.WorkspaceInfo_PasswordHash)
		putError := sh.db.Put(append(entityKeyBase, WORKSPACE_INFO_PASSWORD_HASH...),
			[]byte(settings.PasswordHash), nil)
		if putError != nil {
			return putError
		}

		securityTypeKey := append(entityKeyBase, WORKSPACE_INFO_SECURITY_TYPE...)
		putError = sh.db.Put(securityTypeKey, []byte("password"), nil)
		if putError != nil {
			return putError
		}
	} else {
		// Set the permitted contacts
		settings := workspaceInfo.SecuritySettings.(*proto.WorkspaceInfo_ContactsWrapper)

		securityTypeKey := append(entityKeyBase, WORKSPACE_INFO_SECURITY_TYPE...)
		putError := sh.db.Put(securityTypeKey, []byte("contacts"), nil)
		if putError != nil {
			return putError
		}

		if len(settings.ContactsWrapper.ContactPublicKeys) > 0 {
			contactsIndex := big.NewInt(1)
			indexFound := false
			for !indexFound {
				// Search index is
				// workspaceInfo:<mnemonic>:contact:<index>
				searchIndex := append(
					append(entityKeyBase, WORKSPACE_INFO_CONTACT...), append(delimiter, contactsIndex.Bytes()...)...,
				)
				foundInfo, _ := sh.db.Get(searchIndex, nil)
				if foundInfo == nil || string(foundInfo) == "" {
					indexFound = true
				} else {
					contactsIndex.Add(contactsIndex, big.NewInt(1))
				}
			}

			// workspaceInfo:<mnemonic>:contact:<index>:attributeName => value
			contactKeybase := append(entityKeyBase, WORKSPACE_INFO_CONTACT...)
			contactKeybase = append(contactKeybase, delimiter...)
			for _, contactPublicKey := range settings.ContactsWrapper.ContactPublicKeys {
				key := append(contactKeybase, []byte(contactsIndex.String())...)
				key = append(key, delimiter...)

				// public key
				publicKeyKey := append(key, WORKSPACE_INFO_WORKSPACE_PUBLIC_KEY...)
				putError = sh.db.Put(publicKeyKey, []byte(contactPublicKey), nil)
				if putError != nil {
					return putError
				}

				contactsIndex.Add(contactsIndex, big.NewInt(1))
			}
		}
	}

	return nil
}
