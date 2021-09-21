package types

type NewWorkspaceRequest struct {
	WorkspaceName              string `json:"workspaceName"`
	WorkspaceType              string `json:"workspaceType"`
	WorkspaceAccessControlType string `json:"workspaceAccessControlType"`
	BaseWorkspaceOwnerKeyID    string `json:"baseWorkspaceOwnerKeyID"`

	WorkspaceAccessControl NewWorkspaceACType `json:"workspaceAccessControl"`
	WorkspaceOwners        []string           `json:"workspaceAdditionalOwnerPublicKeys"`
}

type NewWorkspaceACType struct {
	Password string     `json:"password"`
	Contacts []*Contact `json:"contacts"`
}

type NewWorkspaceResponse struct {
	Mnemonic string `json:"mnemonic"`
}

type JoinWorkspaceRequest struct {
	Mnemonic    string `json:"mnemonic"`
	Password    string `json:"password"`
	PublicKeyID string `json:"publicKeyID"`
}

type WorkspaceCredentials struct {
	Mnemonic   string
	PublicKey  *string
	PrivateKey *string
	Password   *string
}

type WorkspaceListResponse struct {
	WorkspaceWrappers []WorkspaceInfoWrapper `json:"workspaceWrappers"`
	Count             int                    `json:"count"`
}

type WorkspaceInfoWrapper struct {
	WorkspaceMnemonic string `json:"workspaceMnemonic"`
	WorkspaceName     string `json:"workspaceName"`
}

type WorkspaceDetailedResponse struct {
	WorkspaceMnemonic string `json:"workspaceMnemonic"`
	WorkspaceName     string `json:"workspaceName"`
	WorkspaceType     string `json:"workspaceType"`

	WorkspaceFiles []FileInfo `json:"workspaceFiles"`
}

type FileInfo struct {
	Name         string `json:"name"`
	Extension    string `json:"extension"`
	Size         int64  `json:"size"`
	DateModified int64  `json:"dateModified"`
	Checksum     string `json:"checksum"`
}

type WorkspacePeersResponse struct {
	NumPeers int `json:"numPeers"`
}
