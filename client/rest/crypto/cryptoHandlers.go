package crypto

import (
	"crypto/rsa"
	"encoding/json"
	"net/http"

	"github.com/zivkovicmilos/peer_drop/crypto"
	"github.com/zivkovicmilos/peer_drop/rest/types"
)

// ValidatePublicKey validates that the requested public key is valid
func ValidatePublicKey(w http.ResponseWriter, r *http.Request) {
	var validatePublicKeyRequest types.ValidatePublicKeyRequest

	decodeErr := json.NewDecoder(r.Body).Decode(&validatePublicKeyRequest)
	if decodeErr != nil {
		http.Error(w, "Unable to parse input", http.StatusBadRequest)
		return
	}

	// Parse the public key
	publicKey, parseError := crypto.ParseRsaPublicKeyFromPemStr(validatePublicKeyRequest.PublicKey)
	if parseError != nil {
		http.Error(w, "Invalid public key", http.StatusBadRequest)
		return
	}

	response := &types.ValidatePublicKeyResponse{
		PublicKeyID: publicKey.KeyIdString(),
	}

	if encodeErr := json.NewEncoder(w).Encode(response); encodeErr != nil {
		http.Error(w, "Unable to encode response", http.StatusInternalServerError)
		return
	}
}

// ValidatePrivateKey validates that the requested private key is valid
func ValidatePrivateKey(w http.ResponseWriter, r *http.Request) {
	var validatePrivateKeyRequest types.ValidatePrivateKeyRequest

	decodeErr := json.NewDecoder(r.Body).Decode(&validatePrivateKeyRequest)
	if decodeErr != nil {
		http.Error(w, "Unable to parse input", http.StatusBadRequest)
		return
	}

	// Parse the public / private key
	privateKey, parseError := crypto.ParsePrivateKeyFromPemStr(validatePrivateKeyRequest.PrivateKey)
	if parseError != nil {
		http.Error(w, "Invalid private key", http.StatusBadRequest)
		return
	}

	publicKeyEncoded, encodeErr := crypto.EncodePublicKey(privateKey.PublicKey.PublicKey.(*rsa.PublicKey))
	if encodeErr != nil {
		http.Error(w, "Unable to encode public key", http.StatusInternalServerError)
		return
	}

	response := &types.ValidatePrivateKeyResponse{
		PublicKey:   publicKeyEncoded,
		PublicKeyID: privateKey.PublicKey.KeyIdString(),
	}

	if encodeErr := json.NewEncoder(w).Encode(response); encodeErr != nil {
		http.Error(w, "Unable to encode response", http.StatusInternalServerError)
		return
	}
}

// GenerateKeyPair generates a new key pair based on the passed in params
func GenerateKeyPair(w http.ResponseWriter, r *http.Request) {
	var gKeyPairRequest types.GenerateKeyPairRequest

	decodeErr := json.NewDecoder(r.Body).Decode(&gKeyPairRequest)
	if decodeErr != nil {
		http.Error(w, "Unable to parse input", http.StatusBadRequest)
		return
	}

	if gKeyPairRequest.KeySize != 2048 && gKeyPairRequest.KeySize != 4096 {
		http.Error(w, "Invalid key size", http.StatusBadRequest)
		return
	}

	privateKeyStr, generateError := crypto.GenerateKeyPair(gKeyPairRequest)
	if generateError != nil {
		http.Error(w, "Unable to generate key pair", http.StatusInternalServerError)
		return
	}

	privateKey, parseError := crypto.ParsePrivateKeyFromPemStr(privateKeyStr)
	if parseError != nil {
		http.Error(w, "Invalid private key", http.StatusInternalServerError)
		return
	}

	response := &types.GenerateKeyPairResponse{
		PrivateKey:  privateKeyStr,
		PublicKeyID: privateKey.PublicKey.KeyIdString(),
	}

	if encodeErr := json.NewEncoder(w).Encode(response); encodeErr != nil {
		http.Error(w, "Unable to encode response", http.StatusInternalServerError)
		return
	}
}
