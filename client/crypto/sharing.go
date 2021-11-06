package crypto

import (
	"crypto/rand"
	"crypto/sha512"
	"fmt"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/google/uuid"
	"golang.org/x/crypto/pbkdf2"
)

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

const ivSize int = 16   // B
const saltSize int = 32 // B

// generateIV generates a randomly filled IV
func generateIV() ([]byte, error) {
	iv := make([]byte, ivSize)
	_, err := rand.Read(iv)
	if err != nil {
		return nil, err
	}

	return iv, nil
}

// generateSalt generates a random salt
func generateSalt() ([]byte, error) {
	salt := make([]byte, saltSize)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}

	return salt, nil
}

// fileSharingSolution is a wrapper for calculated key values
type fileSharingSolution struct {
	AESKey  []byte
	HMACKey []byte
}

// GeneratePasswordFileSharingSolution generates the required
// aes / hmac values based on the password, salt and IV
func GeneratePasswordFileSharingSolution(password string, salt []byte) *fileSharingSolution {
	// We need to generate a 512bit key for AES / HMAC
	generatedKey := pbkdf2.Key([]byte(password), salt, 4096, 64, sha512.New)

	// First half is for AES, second half is for HMAC
	aesKey := generatedKey[:32]
	hmacKey := generatedKey[32:]

	return &fileSharingSolution{
		AESKey:  aesKey,
		HMACKey: hmacKey,
	}
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

// GenerateKeyFileSharingSolution decrypts the
// aes / hmac keys using the private key
func GenerateKeyFileSharingSolution(
	privateKeyPEM string,
	aesEncrypted []byte,
	hmacEncrypted []byte,
) (*fileSharingSolution, error) {
	// RSA key setup
	privateKey, err := ParseRSAKey(privateKeyPEM)
	if err != nil {
		return nil, fmt.Errorf("unable to parse RSA key, %v", err)
	}
	privateKeyRing, err := crypto.NewKeyRing(privateKey)
	if err != nil {
		return nil, fmt.Errorf("unable to parse RSA key ring, %v", err)
	}

	// AES
	aesCipher := crypto.NewPGPMessage(aesEncrypted)
	decryptedAES, err := privateKeyRing.Decrypt(
		aesCipher,
		nil,
		0,
	)
	if err != nil {
		return nil, err
	}

	// HMAC
	hmacCipher := crypto.NewPGPMessage(hmacEncrypted)
	decryptedHMAC, err := privateKeyRing.Decrypt(
		hmacCipher,
		nil,
		0,
	)
	if err != nil {
		return nil, err
	}

	return &fileSharingSolution{
		AESKey:  decryptedAES.GetBinary(),
		HMACKey: decryptedHMAC.GetBinary(),
	}, nil
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
	// Convert the public key PEM block to a key object
	publicKey, parseErr := ParseRSAKey(publicKeyPEM)
	if parseErr != nil {
		return nil, parseErr
	}
	publicKeyRing, err := crypto.NewKeyRing(publicKey)

	aesMessage, err := publicKeyRing.Encrypt(crypto.NewPlainMessage(aesKey), nil)
	if err != nil {
		return nil, err
	}

	hmacMessage, err := publicKeyRing.Encrypt(crypto.NewPlainMessage(hmacKey), nil)
	if err != nil {
		return nil, err
	}

	metadata.EncryptedAESKey = aesMessage.GetBinary()
	metadata.EncryptedHMACKey = hmacMessage.GetBinary()

	return metadata, nil
}
