package crypto

import (
	"bytes"
	"crypto/rsa"
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/zivkovicmilos/peer_drop/rest/types"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/packet"
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
	block, err := armor.Decode(strings.NewReader(privPEM))
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

// EncodePrivateKey returns the PEM private key block
func EncodePrivateKey(key *rsa.PrivateKey) (string, error) {
	buf := new(bytes.Buffer)
	w, err := armor.Encode(buf, openpgp.PrivateKeyType, make(map[string]string))
	if err != nil {
		return "", fmt.Errorf("unable to encode armor, %v", err)
	}

	pgpKey := packet.NewRSAPrivateKey(time.Now(), key)
	err = pgpKey.Serialize(w)
	if err != nil {
		return "", fmt.Errorf("unable to serialize private key, %v", err)
	}
	err = w.Close()
	if err != nil {
		return "", fmt.Errorf("unable to serialize private key, %v", err)
	}

	return buf.String(), nil
}

// EncodePublicKey returns the PEM public key block
func EncodePublicKey(key *rsa.PublicKey) (string, error) {
	buf := new(bytes.Buffer)
	w, err := armor.Encode(buf, openpgp.PublicKeyType, make(map[string]string))
	if err != nil {
		return "", fmt.Errorf("unable to encode armor, %v", err)
	}

	pgpKey := packet.NewRSAPublicKey(time.Now(), key)
	err = pgpKey.Serialize(w)
	if err != nil {
		return "", fmt.Errorf("unable to serialize public key, %v", err)
	}
	err = w.Close()
	if err != nil {
		return "", fmt.Errorf("unable to serialize public key, %v", err)
	}

	return buf.String(), nil
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
