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
