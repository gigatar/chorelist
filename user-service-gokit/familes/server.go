package familes

import (
	"crypto/tls"
	"log"
	"net/http"
	"time"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// NewHTTPServer returns a new http server.
func NewHTTPServer(s Service) *http.Server {
	endpoints := MakeServerEndpoints(s)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeError),
	}

	router := mux.NewRouter()
	familiesAuth := router.PathPrefix("/rest/v1/families").Subrouter()

	// Add middleware
	router.Use(commonMiddleware)
	familiesAuth.Use(validateJWT)

	// Unauthenticated Endpoints

	// Authenticated Endpoints
	familiesAuth.Methods("GET").Path("/{id}").Handler(httptransport.NewServer(
		endpoints.GetFamilyByID,
		decodeGetFamilyByIDRequest,
		encodeGetFamilyByIDResponse,
		options...,
	))

	familiesAuth.Methods("PATCH").Path("/name").Handler(httptransport.NewServer(
		endpoints.ChangeName,
		decodeChangeNameRequest,
		encodeChangeNameResponse,
		options...,
	))

	// Configure CORS
	allowedMethods := handlers.AllowedMethods([]string{
		"GET", "POST", "PATCH", "DELETE", "OPTIONS",
	})

	// Accept everyone - probably should be the LB or something.
	allowedOrigin := handlers.AllowedOrigins([]string{"*"})
	allowedHeaders := handlers.AllowedHeaders([]string{"Content-Type", "Authorization"})
	exposedHeaders := handlers.ExposedHeaders([]string{"Authorization"})

	// Setup HTTPS HERE.
	server := &http.Server{
		Handler:      handlers.CORS(allowedMethods, allowedHeaders, allowedOrigin, exposedHeaders)(router),
		Addr:         ":8081",
		WriteTimeout: 60 * time.Second,
		ReadTimeout:  60 * time.Second,
		TLSConfig: &tls.Config{
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			},

			PreferServerCipherSuites: true,
			InsecureSkipVerify:       false,
			MinVersion:               tls.VersionTLS12,
			MaxVersion:               tls.VersionTLS12,
		},
	}

	log.Println("Starting Families Server")
	return server
}
