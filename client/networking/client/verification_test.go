package client

import (
	"crypto/rsa"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zivkovicmilos/peer_drop/crypto"
	"github.com/zivkovicmilos/peer_drop/rest/types"
)

func TestConstructPasswordChallenge(t *testing.T) {
	password := "myPassword"
	unencrypted := "dummy"

	passwordChallenge, err := ConstructPasswordChallenge([]byte(unencrypted), password)
	assert.NoError(t, err)

	result, err := SolvePasswordChallenge(passwordChallenge, password)
	assert.NoError(t, err)

	assert.Equal(t, unencrypted, string(result.DecryptedValue))
}

func TestConstructPublicKeyChallenge(t *testing.T) {
	unencrypted := "dummy"

	privateKeyPEM, genErr := crypto.GenerateKeyPair(types.GenerateKeyPairRequest{
		KeySize: 2048,
		Name:    "Milos Zivkovic",
		Email:   "milos@zmilos.com",
	})

	privateKey, parseErr := crypto.ParsePrivateKeyFromPemStr(privateKeyPEM)
	assert.NoError(t, parseErr)

	publicKey := privateKey.PublicKey.PublicKey.(*rsa.PublicKey)
	publicKeyPEM, encodeErr := crypto.EncodePublicKeyStr(publicKey)
	assert.NoError(t, encodeErr)

	assert.NoError(t, genErr)

	publicKeyChallenge, err := ConstructPublicKeyChallenge([]byte(unencrypted), publicKeyPEM)
	assert.NoError(t, err)

	solution, solErr := SolvePublicKeyChallenge(publicKeyChallenge, privateKeyPEM)
	assert.NoError(t, solErr)

	assert.Equal(t, unencrypted, string(solution.DecryptedValue))
}
