package rest

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"github.com/zivkovicmilos/peer_drop/rest/contacts"
	"github.com/zivkovicmilos/peer_drop/storage"
)

var (
	serverHost = "localhost"
	serverPort = 5000
)

type Dispatcher struct {
	logger  hclog.Logger
	router  *mux.Router
	server  *http.Server
	storage *storage.StorageHandler
}

// NewDispatcher returns a new instance of the dispatcher
func NewDispatcher(logger hclog.Logger, storage *storage.StorageHandler) *Dispatcher {
	return &Dispatcher{
		logger:  logger.Named("dispatcher"),
		storage: storage,
	}
}

// GetStorageHandler returns the reference to the storage handler
func (d *Dispatcher) GetStorageHandler() *storage.StorageHandler {
	return d.storage
}

// commonMiddleware defines a middleware that all requests go through
func (d *Dispatcher) commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		d.logger.Info(fmt.Sprintf("%s %s", r.Method, r.URL))

		next.ServeHTTP(w, r)
	})
}

// Start starts the http dispatcher service
func (d *Dispatcher) Start(closeChannel <-chan os.Signal) {
	// Set up the router and server
	d.setupRouter()

	// Set up the middleware
	d.router.Use(d.commonMiddleware)

	// Register the available endpoints
	d.registerEndpoints()

	// Start the handle loop
	go func() {
		if err := d.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	d.logger.Info(
		fmt.Sprintf("HTTP server started on: %s:%d", serverHost, serverPort),
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

// setupRouter is a helper method for creating the mux router and http server
func (d *Dispatcher) setupRouter() {
	d.router = mux.NewRouter()
	// Add your routes as needed

	d.server = &http.Server{
		Addr:         fmt.Sprintf("%s:%d", serverHost, serverPort),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		Handler:      d.router,
	}
}

// NotFoundHandler returns a route not found message
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Route not found", http.StatusNotFound)
}

// registerEndpoints registers all available REST endpoints
func (d *Dispatcher) registerEndpoints() {
	// Contacts
	d.router.HandleFunc("/contacts", contacts.GetContacts).Methods("GET")
	d.router.HandleFunc("/contacts", contacts.CreateContact).Methods("POST")
	d.router.HandleFunc("/contacts/{contactId}", contacts.GetContact).Methods("GET")
	d.router.HandleFunc("/contacts/{contactId}", contacts.UpdateContact).Methods("PUT")
	d.router.HandleFunc("/contacts/{contactId}", contacts.DeleteContact).Methods("DELETE")

	d.router.NotFoundHandler = http.HandlerFunc(NotFoundHandler)
}
