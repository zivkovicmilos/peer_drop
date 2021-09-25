package rendezvous

import (
	"encoding/json"
	"net/http"

	"github.com/zivkovicmilos/peer_drop/rest/types"
	servicehandler "github.com/zivkovicmilos/peer_drop/service-handler"
	"github.com/zivkovicmilos/peer_drop/storage"
)

// GetRendezvousNodes fetches the current list of all known rendezvous nodes
func GetRendezvousNodes(w http.ResponseWriter, r *http.Request) {
	// Get the list of rendezvous nodes
	listResponse := make([]*types.RendezvousNodeResponse, 0)

	rendezvousNodes, findErr := storage.GetStorageHandler().GetRendezvousNodes()
	if findErr != nil {
		http.Error(w, "Unable to fetch info", http.StatusInternalServerError)
		return
	}

	for _, nodeMAdddr := range rendezvousNodes {
		// Get connected status
		var status string
		isConnected := servicehandler.GetServiceHandler().GetClientServer().IsConnectedMaddr(nodeMAdddr)

		if isConnected {
			status = "on"
		} else {
			status = "off"
		}

		listResponse = append(listResponse, &types.RendezvousNodeResponse{
			Address: nodeMAdddr,
			Status:  status,
		})
	}

	if encodeErr := json.NewEncoder(w).Encode(&types.RendezvousNodesResponse{
		Count: len(listResponse),
		Data:  listResponse,
	}); encodeErr != nil {
		http.Error(w, "Unable to encode response", http.StatusInternalServerError)
		return
	}
}

// AddRendezvousNode adds a new rendezvous node
func AddRendezvousNode(w http.ResponseWriter, r *http.Request) {
	var request types.RendezvousRequestWrapper

	decodeErr := json.NewDecoder(r.Body).Decode(&request)
	if decodeErr != nil {
		http.Error(w, "Unable to parse input", http.StatusBadRequest)
		return
	}

	// Grab the old list
	oldList, oldErr := storage.GetStorageHandler().GetRendezvousNodes()
	if oldErr != nil {
		http.Error(w, "Unable to parse fetch list", http.StatusInternalServerError)
		return
	}

	// Check if the rendezvous address already exists
	for _, address := range oldList {
		if address == request.Address {
			http.Error(w, "Rendezvous node already exists", http.StatusBadRequest)
			return
		}
	}

	newList := append(oldList, request.Address)

	setErr := storage.GetStorageHandler().SetRendezvousNodes(newList)
	if setErr != nil {
		http.Error(w, "Unable to set new list", http.StatusInternalServerError)
		return
	}

	servicehandler.GetServiceHandler().GetClientServer().SetRendezvous(newList)

	if encodeErr := json.NewEncoder(w).Encode("Rendezvous node added"); encodeErr != nil {
		http.Error(w, "Unable to encode response", http.StatusInternalServerError)
		return
	}
}

// RemoveRendezvousNode removes a rendezvous node
func RemoveRendezvousNode(w http.ResponseWriter, r *http.Request) {
	var request types.RendezvousRequestWrapper

	decodeErr := json.NewDecoder(r.Body).Decode(&request)
	if decodeErr != nil {
		http.Error(w, "Unable to parse input", http.StatusBadRequest)
		return
	}

	// Grab the old list
	oldList, oldErr := storage.GetStorageHandler().GetRendezvousNodes()
	if oldErr != nil {
		http.Error(w, "Unable to parse fetch list", http.StatusInternalServerError)
		return
	}

	foundIndx := -1
	for indx, node := range oldList {
		if node == request.Address {
			foundIndx = indx
			break
		}
	}

	if foundIndx < 0 {
		http.Error(w, "Rendezvous not found", http.StatusNotFound)
		return
	}

	newList := append(oldList[:foundIndx], oldList[foundIndx+1:]...)

	setErr := storage.GetStorageHandler().SetRendezvousNodes(newList)
	if setErr != nil {
		http.Error(w, "Unable to set new list", http.StatusInternalServerError)
		return
	}

	servicehandler.GetServiceHandler().GetClientServer().SetRendezvous(newList)

	if encodeErr := json.NewEncoder(w).Encode("Rendezvous node updated"); encodeErr != nil {
		http.Error(w, "Unable to encode response", http.StatusInternalServerError)
		return
	}
}
