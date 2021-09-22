package crypto

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	libp2pCrypto "github.com/libp2p/go-libp2p-core/crypto"
	"github.com/zivkovicmilos/peer_drop/rest/types"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/packet"
	"golang.org/x/crypto/pbkdf2"
)

// ParseRsaPublicKeyFromPemStr parses a PEM formatted public key
func ParseRsaPublicKeyFromPemStr(pubPEM string) (*packet.PublicKey, error) {
	block, err := armor.Decode(strings.NewReader(pubPEM))
	if err != nil {
		return nil, fmt.Errorf("unable to decode armor, %v", err)
	}

	if block.Type != openpgp.PublicKeyType {
		return nil, errors.New("invalid public key")
	}

	reader := packet.NewReader(block.Body)
	pkt, err := reader.Next()
	if err != nil {
		return nil, fmt.Errorf("unable to read public key, %v", err)
	}

	key, ok := pkt.(*packet.PublicKey)
	if !ok {
		return nil, fmt.Errorf("unable to read public key, %v", err)
	}

	return key, nil
}

// ParsePrivateKeyFromPemStr parses a PEM formatted private key
func ParsePrivateKeyFromPemStr(privPEM string) (*packet.PrivateKey, error) {
	return parsePrivateKey(strings.NewReader(privPEM))
}

// parsePrivateKey parses the private key from the specified input
func parsePrivateKey(inputReader io.Reader) (*packet.PrivateKey, error) {
	block, err := armor.Decode(inputReader)
	if err != nil {
		return nil, fmt.Errorf("unable to decode armor, %v", err)
	}

	if block.Type != openpgp.PrivateKeyType {
		return nil, errors.New("invalid private key")
	}

	reader := packet.NewReader(block.Body)
	pkt, err := reader.Next()
	if err != nil {
		return nil, fmt.Errorf("unable to read private key, %v", err)
	}

	key, ok := pkt.(*packet.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("unable to read private key, %v", err)
	}

	return key, nil
}

// ReadPrivateKey reads the private key from the specified file path
func ReadPrivateKey(dataDir string, keyFileName string) (*packet.PrivateKey, error) {
	in, openErr := os.Open(filepath.Join(dataDir, keyFileName))
	defer in.Close()
	if openErr != nil {
		return nil, fmt.Errorf("Unable to open key file, %v", openErr)
	}

	return parsePrivateKey(in)
}

// EncodePrivateKey returns the PEM private key block
func EncodePrivateKey(output io.Writer, key *rsa.PrivateKey) error {
	w, err := armor.Encode(output, openpgp.PrivateKeyType, make(map[string]string))
	if err != nil {
		return fmt.Errorf("unable to encode armor, %v", err)
	}

	pgpKey := packet.NewRSAPrivateKey(time.Now(), key)
	err = pgpKey.Serialize(w)
	if err != nil {
		return fmt.Errorf("unable to serialize private key, %v", err)
	}
	err = w.Close()
	if err != nil {
		return fmt.Errorf("unable to serialize private key, %v", err)
	}

	return nil
}

// EncodePrivateKeyStr is a wrapper function for converting a private key object to a PEM block string
func EncodePrivateKeyStr(key *rsa.PrivateKey) (string, error) {
	buf := new(bytes.Buffer)

	encodeErr := EncodePrivateKey(buf, key)
	if encodeErr != nil {
		return "", encodeErr
	}

	return buf.String(), nil
}

// EncodePublicKeyStr is a wrapper function for converting a public key object to a PEM block string
func EncodePublicKeyStr(key *rsa.PublicKey) (string, error) {
	buf := new(bytes.Buffer)

	encodeErr := EncodePublicKey(buf, key)
	if encodeErr != nil {
		return "", encodeErr
	}

	return buf.String(), nil
}

// EncodePublicKey returns the PEM public key block
func EncodePublicKey(output io.Writer, key *rsa.PublicKey) error {
	w, err := armor.Encode(output, openpgp.PublicKeyType, make(map[string]string))
	if err != nil {
		return fmt.Errorf("unable to encode armor, %v", err)
	}

	pgpKey := packet.NewRSAPublicKey(time.Now(), key)
	err = pgpKey.Serialize(w)
	if err != nil {
		return fmt.Errorf("unable to serialize public key, %v", err)
	}
	err = w.Close()
	if err != nil {
		return fmt.Errorf("unable to serialize public key, %v", err)
	}

	return nil
}

// GetEmailFromPublicKey extracts the email address from the public key identity
func GetEmailFromPublicKey(pubPEM string) (string, error) {
	identity, identityError := GetIdentityFromPublicKey(pubPEM)
	if identityError != nil {
		return "", fmt.Errorf("Unable to retrieve identity, %v", identityError)
	}

	re := regexp.MustCompile("[a-z0-9._%+\\-]+@[a-z0-9.\\-]+\\.[a-z]{2,4}")
	match := re.FindStringSubmatch(identity.Name)

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
// Returns a PEM encoded private key block
func GenerateKeyPair(request types.GenerateKeyPairRequest) (string, error) {
	config := &packet.Config{
		RSABits: request.KeySize,
	}

	// Create the entity
	entity, entityError := openpgp.NewEntity(request.Name, "", request.Email, config)
	if entityError != nil {
		return "", fmt.Errorf("unable to create entity, %v", entityError)
	}

	// Encode the initial armor
	buf := new(bytes.Buffer)
	w, encodeError := armor.Encode(buf, openpgp.PrivateKeyType, nil)
	if encodeError != nil {
		return "", fmt.Errorf("unable to encode armor, %v", encodeError)
	}

	// Serialize the private data + identity information
	serializeError := entity.SerializePrivate(w, config)
	_ = w.Close()
	if serializeError != nil {
		return "", fmt.Errorf("unable to serialize entity, %v", serializeError)
	}

	return buf.String(), nil
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
		privKey, generateErr := rsa.GenerateKey(rand.Reader, 2048)
		if generateErr != nil {
			return nil, fmt.Errorf("unable to generate private key, %v", generateErr)
		}

		privKeyFile, fileErr := os.Create(filepath.Join(dataDir, keyFileName))
		defer privKeyFile.Close()
		if fileErr != nil {
			return nil, fmt.Errorf("unable to generate private key file, %v", fileErr)
		}

		encodeErr := EncodePrivateKey(privKeyFile, privKey)
		if encodeErr != nil {
			return nil, fmt.Errorf("unable to encode private key, %v", encodeErr)
		}

		return rsaPrivToLibp2pPriv(privKey)
	}

	// Key exists, read it from the disk
	privKey, readErr := ReadPrivateKey(dataDir, keyFileName)
	if readErr != nil {
		return nil, fmt.Errorf("unable to read private key file, %v", readErr)
	}

	return rsaPrivToLibp2pPriv(privKey.PrivateKey.(*rsa.PrivateKey))
}

// rsaPrivToLibp2pPriv does a conversion to the libp2p crypto standard
func rsaPrivToLibp2pPriv(key *rsa.PrivateKey) (libp2pCrypto.PrivKey, error) {
	keyBytes := x509.MarshalPKCS1PrivateKey(key)
	result, unmarshalError := libp2pCrypto.UnmarshalRsaPrivateKey(keyBytes)
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

type PasswordFileSharingMetadata struct {
	IV   []byte
	Salt []byte

	AESKey  []byte
	HMACKey []byte
}

type KeyFileSharingMetadata struct {
	IV               []byte
	EncryptedAESKey  []byte
	EncryptedHMACKey []byte

	AESKey  []byte
	HMACKey []byte
}

const IV_SIZE int = 16   // B
const SALT_SIZE int = 32 // B

// generateIV generates a randomly filled IV
func generateIV() ([]byte, error) {
	iv := make([]byte, IV_SIZE)
	_, err := rand.Read(iv)
	if err != nil {
		return nil, err
	}

	return iv, nil
}

// generateSalt generates a random salt
func generateSalt() ([]byte, error) {
	salt := make([]byte, SALT_SIZE)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}

	return salt, nil
}

// GeneratePasswordFileSharingMetadata generates the required metadata
// for password based file sharing
func GeneratePasswordFileSharingMetadata(password string) (*PasswordFileSharingMetadata, error) {
	metadata := &PasswordFileSharingMetadata{}

	// Generate IV
	iv, ivErr := generateIV()
	if ivErr != nil {
		return nil, ivErr
	}

	// Generate salt
	salt, saltErr := generateSalt()
	if saltErr != nil {
		return nil, saltErr
	}

	metadata.IV = iv
	metadata.Salt = salt

	// We need to generate a 512bit key for AES / HMAC
	generatedKey := pbkdf2.Key([]byte(password), salt, 4096, 64, sha512.New)

	// First half is for AES, second half is for HMAC
	aesKey := generatedKey[:32]
	hmacKey := generatedKey[32:]

	metadata.AESKey = aesKey
	metadata.HMACKey = hmacKey

	return metadata, nil
}

// GenerateKeyFileSharingMetadata generates the required metadata
// for public key based file sharing
func GenerateKeyFileSharingMetadata(publicKeyPEM string) (*KeyFileSharingMetadata, error) {
	metadata := &KeyFileSharingMetadata{}

	// Generate IV
	iv, ivErr := generateIV()
	if ivErr != nil {
		return nil, ivErr
	}
	metadata.IV = iv

	salt, saltErr := generateSalt()
	if saltErr != nil {
		return nil, saltErr
	}

	// We need to generate a 512bit key for AES / HMAC
	generatedKey := pbkdf2.Key([]byte(uuid.New().String()), salt, 4096, 64, sha512.New)

	// First half is for AES, second half is for HMAC
	aesKey := generatedKey[:32]
	hmacKey := generatedKey[32:]

	metadata.AESKey = aesKey
	metadata.HMACKey = hmacKey

	// Encrypt the keys
	// Convert the public key PEM block to an rsa object
	publicKey, parseErr := ParseRsaPublicKeyFromPemStr(publicKeyPEM)
	if parseErr != nil {
		return nil, parseErr
	}

	rsaPublicKey := publicKey.PublicKey.(*rsa.PublicKey)

	// Encrypt the small messages using RSA
	encryptedAES, err := rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		rsaPublicKey,
		aesKey,
		nil)
	if err != nil {
		return nil, err
	}

	encryptedHMAC, err := rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		rsaPublicKey,
		hmacKey,
		nil)
	if err != nil {
		return nil, err
	}

	metadata.EncryptedAESKey = encryptedAES
	metadata.EncryptedHMACKey = encryptedHMAC

	return metadata, nil
}
