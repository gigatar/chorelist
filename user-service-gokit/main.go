package main

import (
	"chorelist/user-service-gokit/familes"
	"chorelist/user-service-gokit/users"
	"fmt"
	"log"
)

func main() {

	fmt.Println("Starting User-Service")

	// Start Family http server
	go func() {
		log.Fatal(familes.NewHTTPServer(familes.NewService()).ListenAndServe())
	}()

	// Main HTTP server
	log.Fatal(users.NewHTTPServer(users.NewService()).ListenAndServe())
}
