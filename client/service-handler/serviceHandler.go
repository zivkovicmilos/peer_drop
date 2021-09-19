package service_handler

import (
	"os"
	"sync"

	"github.com/zivkovicmilos/peer_drop/networking/client"
)

type ServiceHandler struct {
	closeChannel chan os.Signal
	clientServer *client.ClientServer

	serviceListeners map[string]chan struct{}
}

var serviceHandlerInstance ServiceHandler
var once sync.Once

// GetServiceHandler initializes the service handler singleton
func GetServiceHandler() *ServiceHandler {
	once.Do(func() {
		serviceHandlerInstance = ServiceHandler{
			serviceListeners: make(map[string]chan struct{}),
		}
	})

	return &serviceHandlerInstance
}

// SetCloseChannel sets the close channel for the service handler
func (sh *ServiceHandler) SetCloseChannel(closeChannel chan os.Signal) {
	sh.closeChannel = closeChannel
}

// SetClientServer sets the client server
func (sh *ServiceHandler) SetClientServer(clientServer *client.ClientServer) {
	sh.clientServer = clientServer
}

// GetClientServer returns a reference to the client service
func (sh *ServiceHandler) GetClientServer() *client.ClientServer {
	return sh.clientServer
}

// BroadcastNotifier waits for a term signal and alerts all listening services
func (sh *ServiceHandler) BroadcastNotifier() {
	<-sh.closeChannel

	// Alert all registered services
	for _, serviceChannel := range sh.serviceListeners {
		serviceChannel <- struct{}{}
	}
}

// RegisterCloseListener adds a new listener for Service Handler's close events
func (sh *ServiceHandler) RegisterCloseListener(serviceName string) chan struct{} {
	channel := make(chan struct{}, 1)
	sh.serviceListeners[serviceName] = channel

	return sh.serviceListeners[serviceName]
}
