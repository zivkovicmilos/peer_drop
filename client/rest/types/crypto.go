package types

type ValidatePublicKeyRequest struct {
	PublicKey string `json:"publicKey"`
}

type ValidatePublicKeyResponse struct {
	PublicKeyID string `json:"publicKeyID"`
}

type ValidatePrivateKeyRequest struct {
	PrivateKey string `json:"privateKey"`
}

type ValidatePrivateKeyResponse struct {
	PublicKey   string `json:"publicKey"`
	PublicKeyID string `json:"publicKeyID"`
}

type GenerateKeyPairRequest struct {
	KeySize int    `json:"keySize"`
	Name    string `json:"name"`
	Email   string `json:"email"`
}

type GenerateKeyPairResponse struct {
	PrivateKey  string `json:"privateKey"`
	PublicKeyID string `json:"publicKeyID"`
}
