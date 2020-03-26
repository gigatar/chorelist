package users

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
	users := router.PathPrefix("/rest/v1/users").Subrouter()
	usersInternal := router.PathPrefix("/internal/v1/users").Subrouter()
	usersAuth := router.PathPrefix("/rest/v1/users").Subrouter()

	// Add middleware
	router.Use(commonMiddleware)
	router.Use(listBackendMiddleware)
	usersAuth.Use(validateJWT)

	// Unauthenticated Endpoints
	users.Methods("POST").Path("/login").Handler(httptransport.NewServer(
		endpoints.Login,
		decodeLoginRequest,
		encodeLoginResponse,
		options...,
	))

	// Authenticated Endpoints
	usersAuth.Methods("GET").Path("").Handler(httptransport.NewServer(
		endpoints.GetUsers,
		decodeGetUsersRequest,
		encodeGetUsersResponse,
		options...,
	))
	usersAuth.Methods("GET").Path("/{id}").Handler(httptransport.NewServer(
		endpoints.GetUserByID,
		decodeGetUserByIDRequest,
		encodeGetUserByIDResponse,
		options...,
	))
	usersAuth.Methods("PATCH").Path("/name").Handler(httptransport.NewServer(
		endpoints.ChangeName,
		decodeChangeNameRequest,
		encodeChangeNameResponse,
		options...,
	))
	usersAuth.Methods("PATCH").Path("/password").Handler(httptransport.NewServer(
		endpoints.ChangePassword,
		decodeChangePasswordRequest,
		encodeChangePasswordResponse,
		options...,
	))

	// Internal Endpoints
	usersInternal.Methods("POST").Path("/create").Handler(httptransport.NewServer(
		endpoints.CreateUser,
		decodeCreateUserRequest,
		encodeCreateUserResponse,
		options...,
	))
	usersInternal.Methods("DELETE").Path("/{id}").Handler(httptransport.NewServer(
		endpoints.DeleteUser,
		decodeDeleteUserRequest,
		encodeDeleteUserResponse,
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
		Addr:         ":9000",
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

	log.Println("Starting User Server")
	return server
}
