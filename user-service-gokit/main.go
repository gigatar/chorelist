package main

import (
	"chorelist/user-service-gokit/users"
	"log"
)

func main() {

	log.Fatal(users.NewHTTPServer(users.NewService()).ListenAndServe())
}
