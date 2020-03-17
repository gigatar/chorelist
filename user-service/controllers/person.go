package controllers

import (
	"chorelist/user-service/daos"
	"chorelist/user-service/models"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gigatar/chorelist/token"

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

// Login to the service.
func (p *PersonController) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Decode user
	var inputPerson models.Person
	err := json.NewDecoder(r.Body).Decode(&inputPerson)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Pull user from database.
	person, err := p.dao.Login(inputPerson.Email)
	if err != nil {
		if strings.Contains(err.Error(), "no documents") {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	// Validate password
	if err := bcrypt.CompareHashAndPassword([]byte(person.Password), []byte(inputPerson.Password)); err != nil {
		if strings.Contains(err.Error(), "hashedPassword is not the hash of the given password") {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	// create jwt
	var token token.JWTToken
	tokenString, err := token.CreateJWT(r.RemoteAddr, person.ID, person.FamilyID)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Authorization", tokenString)
}

// ChangeName allows a person to change their name.
func (p *PersonController) ChangeName(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Deserialize request
	var inputPerson models.Person
	err := json.NewDecoder(r.Body).Decode(&inputPerson)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Validate input
	if !inputPerson.ValidateName() {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get userID to ensure we can only modify ourselves
	var jwt token.JWTToken
	userID, err := jwt.GetUser(r.Header.Get("authorization"))
	if err != nil {
		if strings.Contains(err.Error(), "Invalid token") {
			w.WriteHeader(http.StatusUnauthorized)
		} else {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	if userID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = p.dao.ChangeName(userID, inputPerson.Name)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// success but no body
	w.WriteHeader(http.StatusNoContent)
}
