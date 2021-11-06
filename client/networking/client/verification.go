package client

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"time"

	libCrypto "github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/google/uuid"
	"github.com/zivkovicmilos/peer_drop/crypto"
	"github.com/zivkovicmilos/peer_drop/proto"
)

func ConstructPasswordChallenge(unencryptedData []byte, password string) (*proto.Challenge, error) {
	passwordChallenge := &proto.Challenge{
		ChallengeId: uuid.New().String(),
	}

	// Set the timestamp to prevent replay attacks
	passwordChallenge.Timestamp = time.Now().Unix()

	block, _ := aes.NewCipher(crypto.NewSHA256([]byte(password)))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	passwordChallenge.EncryptedValue = gcm.Seal(nonce, nonce, unencryptedData, nil)

	return passwordChallenge, nil
}

func ConstructPublicKeyChallenge(unencryptedData []byte, publicKeyPEM string) (*proto.Challenge, error) {
	publicKeyChallenge := &proto.Challenge{
		ChallengeId: uuid.New().String(),
	}

	// Set the timestamp to prevent replay attacks
	publicKeyChallenge.Timestamp = time.Now().Unix()

	// Convert the public key PEM block to a key object
	publicKey, parseErr := crypto.ParseRSAKey(publicKeyPEM)
	if parseErr != nil {
		return nil, parseErr
	}

	publicKeyRing, err := libCrypto.NewKeyRing(publicKey)
	encryptedData, err := publicKeyRing.Encrypt(libCrypto.NewPlainMessage(unencryptedData), nil)
	if err != nil {
		return nil, err
	}
	if err != nil {
		panic(err)
	}

	publicKeyChallenge.EncryptedValue = encryptedData.GetBinary()

	return publicKeyChallenge, nil
}

// ConstructVerificationResponse constructs a verification response
func ConstructVerificationResponse(message string, confirmed bool) *proto.VerificationResponse {
	return &proto.VerificationResponse{
		Message:   message,
		Confirmed: confirmed,
	}
}

// SolvePasswordChallenge attempts to solve the password handshake challenge
func SolvePasswordChallenge(
	challenge *proto.Challenge,
	password string,
) (*proto.ChallengeSolution, error) {
	solution := &proto.ChallengeSolution{}
	solution.ChallengeId = challenge.ChallengeId

	key := crypto.NewSHA256([]byte(password))

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()

	nonce, subEncryptedDataCypher := challenge.EncryptedValue[:nonceSize], challenge.EncryptedValue[nonceSize:]
	decryptedData, err := gcm.Open(nil, nonce, subEncryptedDataCypher, nil)
	if err != nil {
		return nil, err
	}

	solution.DecryptedValue = decryptedData

	return solution, nil
}

// SolvePublicKeyChallenge attempts to solve the public key handshake challenge
func SolvePublicKeyChallenge(
	challenge *proto.Challenge,
	privateKeyPEM string,
) (*proto.ChallengeSolution, error) {
	solution := &proto.ChallengeSolution{}
	solution.ChallengeId = challenge.ChallengeId

	privateKey, parseErr := crypto.ParseRSAKey(privateKeyPEM)
	if parseErr != nil {
		return nil, parseErr
	}

	privateKeyRing, err := libCrypto.NewKeyRing(privateKey)
	if err != nil {
		return nil, fmt.Errorf("unable to parse RSA key ring, %v", err)
	}

	encryptedData := libCrypto.NewPGPMessage(challenge.EncryptedValue)
	decryptedData, err := privateKeyRing.Decrypt(
		encryptedData,
		nil,
		0,
	)
	if err != nil {
		return nil, err
	}

	solution.DecryptedValue = decryptedData.GetBinary()

	return solution, nil
}
