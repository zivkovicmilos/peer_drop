package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
