package types

type Identity struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	DateCreated   string `json:"dateCreated"`
	PrivateKey    string `json:"-"`
	PublicKey     string `json:"-"`
	PublicKeyID   string `json:"publicKeyID"`
	IsPrimary     bool   `json:"isPrimary"`
	NumWorkspaces int    `json:"numWorkspaces"`
}

type IdentityUpdateRequest struct {
	Name       string `json:"name"`
	Picture    string `json:"picture"`
	PrivateKey string `json:"privateKey"`
}

type IdentitiesResponse struct {
	Data  []*Identity `json:"data"`
	Count int         `json:"count"`
}

type NewIdentityRequest struct {
	Name       string `json:"name"`
	Picture    string `json:"picture"`
	PrivateKey string `json:"privateKey"`
}
