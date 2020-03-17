package controllers

import (
	"chorelist/user-service/daos"
	"chorelist/user-service/models"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

const bcryptPasswordCost = 8

// PersonController defines all the methods for interacting with Person endpoints.
type PersonController struct {
	dao daos.PersonDAO
}

// CreatePerson creates a new person in the system.
func (p *PersonController) CreatePerson(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Decode input
	var inputPerson models.Person

	err := json.NewDecoder(r.Body).Decode(&inputPerson)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Validate input
	if !inputPerson.Validate() {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Make sure we don't have an ID
	inputPerson.ID = ""

	// TODO: Check to make sure person email is unique
	// This could be done with a unique key in the db too.

	// Hash password
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(inputPerson.Password), bcryptPasswordCost)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set Request password to the hash
	inputPerson.Password = string(encryptedPassword)

	// Insert into database
	resourceID, err := p.dao.CreatePerson(inputPerson)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			w.WriteHeader(http.StatusConflict)
		} else {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
		}
		return
	}

	// Write resource location header
	location := "/rest/users/" + resourceID
	w.Header().Set("Location", location)

	w.WriteHeader(http.StatusCreated)
}
