package main

import (
	"chorelist/chore-service/controllers"
	"chorelist/chore-service/database"
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

	var chore controllers.ChoreController
	var jwt token.JWTToken

	router := mux.NewRouter()
	rest := router.PathPrefix("/rest/v1").Subrouter()
	choreEndpoint := rest.PathPrefix("/chores").Subrouter()

	// Unauthenticated endpoints

	// Authenticated endpoints
	choreEndpoint.HandleFunc("", jwt.ValidateMiddleware(chore.ListFamilyChores)).Methods("GET")
	choreEndpoint.HandleFunc("", jwt.ValidateMiddleware(chore.AddChore)).Methods("POST")

	// Configure CORS
	allowedMethods := handlers.AllowedMethods([]string{
		"GET", "PATCH", "POST", "OPTIONS",
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
