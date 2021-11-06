package crypto

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
	libp2pCrypto "github.com/libp2p/go-libp2p-core/crypto"
	"github.com/zivkovicmilos/peer_drop/rest/types"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/packet"
)

// ParseRSAKey parses a private or public RSA key from PEM encoded string
func ParseRSAKey(keyPEM string) (*crypto.Key, error) {
	key, err := crypto.NewKeyFromArmored(keyPEM)
	if err != nil {
		return nil, fmt.Errorf("unable to read key armor, %v", err)
	}

	return key, nil
}

// ReadLibp2pPrivateKey reads the private key from the specified file path
func ReadLibp2pPrivateKey(dataDir string, keyFileName string) ([]byte, error) {
	b, openErr := os.ReadFile(filepath.Join(dataDir, keyFileName))
	if openErr != nil {
		return nil, fmt.Errorf("unable to open key file, %v", openErr)
	}

	return b, nil
}

// EncodeKeyToFile encodes an RSA key to a provided file
func EncodeKeyToFile(output *os.File, key *crypto.Key) error {
	keyArmor, armorErr := key.Armor()
	if armorErr != nil {
		return fmt.Errorf("unable to encode armor, %v", armorErr)
	}

	_, err := output.WriteString(keyArmor)
	if err != nil {
		return fmt.Errorf("unable to write key, %v", err)
	}

	return nil
}

// GetEmailFromIdentityName extracts the email address from the public key identity
func GetEmailFromIdentityName(name string) (string, error) {
	re := regexp.MustCompile("[a-z0-9._%+\\-]+@[a-z0-9.\\-]+\\.[a-z]{2,4}")
	match := re.FindStringSubmatch(name)

	if len(match) > 0 {
		return match[0], nil
	} else {
		return "unknown email", nil
	}
}

// GetIdentityFromPublicKey gets the identity information from the public key
func GetIdentityFromPublicKey(pubPEM string) (*openpgp.Identity, error) {
	block, err := armor.Decode(strings.NewReader(pubPEM))
	if err != nil {
		return nil, fmt.Errorf("unable to decode armor, %v", err)
	}

	if block.Type != openpgp.PublicKeyType {
		return nil, errors.New("invalid public key")
	}

	reader := packet.NewReader(block.Body)
	entity, entityError := openpgp.ReadEntity(reader)
	if entityError != nil {
		return nil, fmt.Errorf("unable to read entity, %v", entityError)
	}

	var identity *openpgp.Identity
	for _, presentIdentity := range entity.Identities {
		identity = presentIdentity
		break
	}

	return identity, nil
}

// GenerateKeyPair generates a new RSA key pair using the bit size.
// Returns a private key
func GenerateKeyPair(request types.GenerateKeyPairRequest) (*crypto.Key, error) {
	return crypto.GenerateKey(request.Name, request.Email, "rsa", request.KeySize)
}

// GetKeyIDFromPEM returns the public key's ID from the given PEM string
func GetKeyIDFromPEM(publicKeyPEM string) (string, error) {
	publicKey, err := ParseRSAKey(publicKeyPEM)
	if err != nil {
		return "", err
	}

	return publicKey.GetHexKeyID(), nil
}

// GetKeyID returns the public key's ID
func GetKeyID(modulus []byte, long bool) string {
	var size int
	if long {
		size = 8
	} else {
		size = 4
	}

	last4Bytes := modulus[len(modulus)-size:]

	return strings.ToUpper(hex.EncodeToString(last4Bytes))
}

// ReadLibp2pKey reads the libp2p private key from the passed in directory,
// or creates it if it doesn't exist
func ReadLibp2pKey(dataDir string, keyFileName string) (libp2pCrypto.PrivKey, error) {
	path := filepath.Join(dataDir, keyFileName)
	_, err := os.Stat(path)
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("unable to stat file %s, %v", path, err)
	}

	if os.IsNotExist(err) {
		// Key doesn't exist yet, generate it
		privateKey, _, generateErr := libp2pCrypto.GenerateKeyPair(libp2pCrypto.RSA, 2048)
		if generateErr != nil {
			return nil, fmt.Errorf("unable to generate private key, %v", generateErr)
		}

		privKeyFile, fileErr := os.Create(filepath.Join(dataDir, keyFileName))
		defer privKeyFile.Close()
		if fileErr != nil {
			return nil, fmt.Errorf("unable to generate private key file, %v", fileErr)
		}

		keyRaw, rawErr := privateKey.Raw()
		if rawErr != nil {
			return nil, fmt.Errorf("unable to get raw private key, %v", rawErr)
		}

		_, encodeErr := privKeyFile.Write(keyRaw)
		if encodeErr != nil {
			return nil, fmt.Errorf("unable to encode private key, %v", encodeErr)
		}

		return rsaPrivToLibp2pPriv(keyRaw)
	}

	// Key exists, read it from the disk
	privKey, readErr := ReadLibp2pPrivateKey(dataDir, keyFileName)
	if readErr != nil {
		return nil, fmt.Errorf("unable to read private key file, %v", readErr)
	}

	return rsaPrivToLibp2pPriv(privKey)
}

// rsaPrivToLibp2pPriv does a conversion to the libp2p crypto standard
func rsaPrivToLibp2pPriv(key []byte) (libp2pCrypto.PrivKey, error) {
	result, unmarshalError := libp2pCrypto.UnmarshalRsaPrivateKey(key)
	if unmarshalError != nil {
		return nil, fmt.Errorf("unable to unmarshal rsa private key, %v", unmarshalError)
	}

	return result, nil
}

// NewSHA256 hashes the input data
func NewSHA256(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}
