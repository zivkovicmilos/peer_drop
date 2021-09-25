package dispatcher

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"github.com/rs/cors"
	"github.com/zivkovicmilos/peer_drop/config"
	"github.com/zivkovicmilos/peer_drop/rest/contacts"
	"github.com/zivkovicmilos/peer_drop/rest/crypto"
	"github.com/zivkovicmilos/peer_drop/rest/identities"
	"github.com/zivkovicmilos/peer_drop/rest/rendezvous"
	"github.com/zivkovicmilos/peer_drop/rest/search"
	"github.com/zivkovicmilos/peer_drop/rest/workspaces"
	"github.com/zivkovicmilos/peer_drop/storage"
)

type Dispatcher struct {
	logger     hclog.Logger
	nodeConfig *config.NodeConfig
	router     *mux.Router
	server     *http.Server
	storage    *storage.StorageHandler
}

// NewDispatcher returns a new instance of the dispatcher
func NewDispatcher(
	logger hclog.Logger,
	nodeConfig *config.NodeConfig,
) *Dispatcher {
	return &Dispatcher{
		logger:     logger.Named("dispatcher"),
		nodeConfig: nodeConfig,
	}
}

// commonMiddleware defines a middleware that all requests go through
func (d *Dispatcher) commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//d.setupHeader(w, r)
		w.Header().Set("Content-Type", "application/json")

		d.logger.Info(fmt.Sprintf("%s %s", r.Method, r.URL))
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Start starts the http dispatcher service
func (d *Dispatcher) Start(closeChannel chan struct{}) {
	// Set up the router and server
	//d.setupRouter()
	d.router = mux.NewRouter()

	// Set up the middleware
	d.router.Use(d.commonMiddleware)

	// Register the available endpoints
	d.registerEndpoints()

	corsConfig := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000", "http://localhost:4000"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	})

	d.setupServer(corsConfig.Handler(d.router))

	// Start the handle loop
	go func() {
		if err := d.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	d.logger.Info(
		fmt.Sprintf("HTTP server started on: %s:%d", d.nodeConfig.HostAddress, d.nodeConfig.HttpPort),
	)

	<-closeChannel
	d.logger.Info("Caught stop signal...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err := d.server.Shutdown(ctx); err != nil {
		d.logger.Error(fmt.Sprintf("Unable to gracefully shut down server, %v", err))
	}

	d.logger.Info("HTTP server stopped gracefully")
}

// setupServer is a helper method for creating the http server
func (d *Dispatcher) setupServer(handler http.Handler) {

	d.server = &http.Server{
		Addr:         fmt.Sprintf("%s:%d", d.nodeConfig.HostAddress, d.nodeConfig.HttpPort),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		Handler:      handler,
	}
}

// NotFoundHandler returns a route not found message
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Route not found", http.StatusNotFound)
}

// ShutdownHandler closes the peer_drop server service
func ShutdownHandler(w http.ResponseWriter, r *http.Request) {
	p, err := os.FindProcess(os.Getpid())
	if err != nil {
		http.Error(w, "Unable to find process", http.StatusInternalServerError)
		return
	}

	encodeErr := json.NewEncoder(w).Encode("Shutting down...")
	if encodeErr != nil {
		http.Error(w, "Unable to encode response", http.StatusInternalServerError)
	}

	_ = p.Signal(syscall.SIGQUIT)
}

// registerEndpoints registers all available REST endpoints
func (d *Dispatcher) registerEndpoints() {
	// Contacts
	d.router.HandleFunc("/api/contacts", contacts.GetContacts).Methods("GET")
	d.router.HandleFunc("/api/contacts", contacts.CreateContact).Methods("POST")
	d.router.HandleFunc("/api/contacts/{contactId}", contacts.GetContact).Methods("GET")
	d.router.HandleFunc("/api/contacts/{contactId}", contacts.UpdateContact).Methods("PUT")
	d.router.HandleFunc("/api/contacts/{contactId}", contacts.DeleteContact).Methods("DELETE")

	// Identities
	d.router.HandleFunc("/api/identities", identities.GetIdentities).Methods("GET")
	d.router.HandleFunc("/api/identities", identities.CreateIdentity).Methods("POST")
	d.router.HandleFunc("/api/identities/{identityId}", identities.GetIdentity).Methods("GET")
	d.router.HandleFunc("/api/identities/{identityId}/public-key", identities.GetPublicKey).Methods("GET")
	d.router.HandleFunc("/api/identities/{identityId}/private-key", identities.GetPrivateKey).Methods("GET")
	d.router.HandleFunc("/api/identities/{identityId}", identities.UpdateIdentity).Methods("PUT")
	d.router.HandleFunc("/api/identities/{identityId}/set-primary", identities.SetAsPrimary).Methods("PUT")
	d.router.HandleFunc("/api/identities/{identityId}", identities.DeleteIdentity).Methods("DELETE")
	d.router.HandleFunc("/api/me", identities.GetPrimaryIdentity).Methods("GET")

	// Crypto
	d.router.HandleFunc("/api/crypto/validate-public-key", crypto.ValidatePublicKey).Methods("POST")
	d.router.HandleFunc("/api/crypto/validate-private-key", crypto.ValidatePrivateKey).Methods("POST")
	d.router.HandleFunc("/api/crypto/generate-key-pair", crypto.GenerateKeyPair).Methods("POST")

	// Workspaces
	d.router.HandleFunc("/api/workspaces", workspaces.CreateWorkspace).Methods("POST")
	d.router.HandleFunc("/api/workspaces", workspaces.GetWorkspaces).Methods("GET")
	d.router.HandleFunc("/api/workspaces/{mnemonic}", workspaces.LeaveWorkspace).Methods("DELETE")
	d.router.HandleFunc("/api/workspaces/{mnemonic}/files", workspaces.GetWorkspaceFiles).Methods("GET")
	d.router.HandleFunc("/api/workspaces/{mnemonic}/peers", workspaces.GetWorkspaceNumPeers).Methods("GET")
	d.router.HandleFunc("/api/join-workspace", workspaces.JoinWorkspace).Methods("POST")
	d.router.HandleFunc("/api/workspaces/{mnemonic}", workspaces.GetWorkspaceInfo).Methods("GET")
	d.router.HandleFunc("/api/workspaces/upload", workspaces.AddFileToWorkspace).Methods("POST")
	d.router.HandleFunc("/api/workspaces/download", workspaces.DownloadWorkspaceFile).Methods("POST")

	// Rendezvous
	d.router.HandleFunc("/api/rendezvous", rendezvous.GetRendezvousNodes).Methods("GET")
	d.router.HandleFunc("/api/rendezvous", rendezvous.AddRendezvousNode).Methods("POST")
	d.router.HandleFunc("/api/rendezvous", rendezvous.RemoveRendezvousNode).Methods("DELETE")

	// Search
	d.router.HandleFunc("/api/search", search.GetSearchResults).Methods("GET")

	// Shutdown handler
	d.router.HandleFunc("/api/shutdown", ShutdownHandler).Methods("POST")

	d.router.NotFoundHandler = http.HandlerFunc(NotFoundHandler)

	http.Handle("/", d.router)
}

// setupHeaders sets up the base CORS config
func (d *Dispatcher) setupHeader(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding")
}
