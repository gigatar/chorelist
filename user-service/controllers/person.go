package controllers

import (
	"chorelist/user-service/daos"
	"chorelist/user-service/models"
	"encoding/json"
	"errors"
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
			w.WriteHeader(http.StatusUnauthorized)
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

// ChangePassword of a person.
func (p *PersonController) ChangePassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Deserialize request
	var inputPerson models.Person
	err := json.NewDecoder(r.Body).Decode(&inputPerson)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Validate input
	if !inputPerson.ValidatePassword() {
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
			w.WriteHeader(http.StatusUnauthorized)
		}
		return
	}

	if userID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get old password from database
	oldPassword, err := p.dao.GetEncryptedPassword(userID)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Compare new and old with bcrypt
	if err := bcrypt.CompareHashAndPassword([]byte(oldPassword), []byte(inputPerson.OldPassword)); err != nil {
		if strings.Contains(err.Error(), "hashedPassword is not the hash of the given password") {
			w.WriteHeader(http.StatusForbidden)
		} else {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	// Hash new
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(inputPerson.Password), bcryptPasswordCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Insert into database
	err = p.dao.UpdatePassword(userID, string(hashedPassword))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeletePerson removes a person from the system
func (p *PersonController) DeletePerson(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get userID to ensure we can only modify ourselves
	var jwt token.JWTToken
	userID, err := jwt.GetUser(r.Header.Get("authorization"))
	if err != nil {
		if strings.Contains(err.Error(), "Invalid token") {
			w.WriteHeader(http.StatusUnauthorized)
		} else {
			log.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
		}
		return
	}

	if userID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = p.dao.DeletePerson(userID)
	if err != nil {
		if strings.Contains(err.Error(), "no documents") {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// getPersonType returns the type of person based on UID
func (p *PersonController) getPersonType(personID string) (string, error) {

	personType, err := p.dao.GetPersonType(personID)
	if err != nil {
		return "", err
	}

	return personType, nil
}

// createPerson based on input.
func (p *PersonController) createPerson(person models.Person) error {
	if !person.Validate() {
		return errors.New("Invalid input")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(person.Password), bcryptPasswordCost)
	if err != nil {
		return err
	}

	person.Password = string(hashedPassword)
	_, err = p.dao.CreatePerson(person)
	if err != nil {
		return err
	}

	return nil
}
