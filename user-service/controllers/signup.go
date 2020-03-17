package controllers

import (
	"chorelist/user-service/daos"
	"chorelist/user-service/models"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// SignupController defines all endpoints for interacting with signup.
type SignupController struct {
	dao daos.SignupDAO
}

// CreateSignup creates a new signup in the system.
func (s *SignupController) CreateSignup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Deserialize request
	var inputSignup models.Signup
	err := json.NewDecoder(r.Body).Decode(&inputSignup)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Validate input
	if !inputSignup.Validate() {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Verify email unqiue
	var p PersonController
	exists, err := p.dao.EmailExists(inputSignup.Person.Email)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if exists {
		w.WriteHeader(http.StatusConflict)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(inputSignup.Person.Password), bcryptPasswordCost)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	inputSignup.Person.Password = string(hashedPassword)

	// Create Unique Code
	code := s.generateCode(15)
	if len(code) != 15 {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	inputSignup.Code = code

	// Set Expiration time
	inputSignup.SignupTime = time.Now().Unix()

	// Insert
	err = s.dao.CreateSignup(inputSignup)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			w.WriteHeader(http.StatusConflict)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
		return
	}

	//TODO: Email code

	w.WriteHeader(http.StatusAccepted)
}

// generateCode returns a unique code of length.
func (s *SignupController) generateCode(length int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyz0123456789"
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
