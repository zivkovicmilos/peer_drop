package types

type Contact struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	DateAdded   string `json:"dateAdded"`
	PublicKeyID string `json:"publicKeyID"`
	PublicKey   string `json:"publicKey"`
}

type ContactsResponse struct {
	Data  []*Contact `json:"data"`
	Count int        `json:"count"`
}
