package types

type ValidatePublicKeyRequest struct {
	PublicKey string `json:"publicKey"`
}

type ValidatePublicKeyResponse struct {
	PublicKeyID string `json:"publicKeyID"`
}
