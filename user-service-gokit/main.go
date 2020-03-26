package main

import (
	"chorelist/user-service-gokit/families"
	"chorelist/user-service-gokit/users"
	"context"
	"fmt"
	"log"
)

func main() {

	fmt.Println("Starting User-Service")

	// Start Family http server
	go func() {
		var db families.Database
		db.Init(context.TODO())
		log.Fatal(families.NewHTTPServer(families.NewService()).ListenAndServe())
	}()

	// Main HTTP server
	var db users.Database
	db.Init(context.TODO())
	log.Fatal(users.NewHTTPServer(users.NewService()).ListenAndServe())
}
