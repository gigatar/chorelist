package main

import (
	"chorelist/user-service/controllers"
	"chorelist/user-service/database"
	"crypto/tls"
	"log"
	"net/http"
	"time"

	"github.com/gigatar/chorelist/token"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {

	// Initialize Database
	if err := database.DB.Init(); err != nil {
		log.Fatal(err)
	}

	var person controllers.PersonController
	var jwt token.JWTToken

	router := mux.NewRouter()
	rest := router.PathPrefix("/rest/v1").Subrouter()
	personEndpoint := rest.PathPrefix("/users").Subrouter()

	// Unauthenticated endpoints
	personEndpoint.HandleFunc("/login", person.Login).Methods("POST")

	// Authenticated endpoints
	personEndpoint.HandleFunc("", jwt.ValidateMiddleware(person.CreatePerson)).Methods("POST")

	// Configure CORS
	allowedMethods := handlers.AllowedMethods([]string{
		"GET", "POST", "PUT", "DELETE", "OPTIONS",
	})

	// Accept everyone - probably should be the LB or something.
	allowedOrigin := handlers.AllowedOrigins([]string{"*"})
	allowedHeaders := handlers.AllowedHeaders([]string{"Content-Type", "Authorization"})
	exposedHeaders := handlers.ExposedHeaders([]string{"Authorization"})

	// Setup HTTPS HERE.
	server := &http.Server{
		Handler:      handlers.CORS(allowedMethods, allowedHeaders, allowedOrigin, exposedHeaders)(router),
		Addr:         ":8080",
		WriteTimeout: 120 * time.Second,
		ReadTimeout:  120 * time.Second,
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

	// Listen HTTP
	log.Fatal(server.ListenAndServe())

	// Listen TLS
	// log.Fatal(server.ListenAndServeTLS(CertificatePath, KeyPath))
}
