package crypto

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
)

// ParseRsaPublicKeyFromPemStr parses a PEM formatted string
func ParseRsaPublicKeyFromPemStr(pubPEM string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pubPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	switch pub := pub.(type) {
	case *rsa.PublicKey:
		return pub, nil
	default:
		break
	}

	return nil, errors.New("key type is not RSA")
}

// GetKeyID returns the public key's ID
func GetKeyID(modulus []byte) string {
	last4Bytes := modulus[len(modulus)-4:]

	return hex.EncodeToString(last4Bytes)
}
