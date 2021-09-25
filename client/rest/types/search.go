package types

type SearchResults struct {
	Identities []*Identity             `json:"identities"`
	Workspaces []*WorkspaceInfoWrapper `json:"workspaces"`
	Contacts   []*Contact              `json:"contacts"`
}
