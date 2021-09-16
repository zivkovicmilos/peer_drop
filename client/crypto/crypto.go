package crypto

import (
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/packet"
)

// ParseRsaPublicKeyFromPemStr parses a PEM formatted string
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
