package crypto

import (
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
