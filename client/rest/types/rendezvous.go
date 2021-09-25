package types

type RendezvousNodesResponse struct {
	Data  []*RendezvousNodeResponse `json:"data"`
	Count int                       `json:"count"`
}

type RendezvousNodeResponse struct {
	Address string `json:"address"`
	Status  string `json:"status"`
}

type RendezvousRequestWrapper struct {
	Address string `json:"address"`
}
