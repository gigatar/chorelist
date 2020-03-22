package controllers

import (
	"bytes"
	"chorelist/user-service/daos"
	"chorelist/user-service/models"
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"

	"golang.org/x/crypto/bcrypt"
)

const codeLength = 15

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
	code := s.generateCode(codeLength)
	if len(code) != codeLength {
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

	// Email code
	if err := s.SendEmail(inputSignup.Person.Email, inputSignup.Person.Name, inputSignup.Code); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

// SignupVerify approves the signup request and creates the Person and Family
func (s *SignupController) SignupVerify(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	code := mux.Vars(r)["code"]
	if len(code) != codeLength {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get Request for Code
	signup, err := s.dao.GetSignup(code)
	if err != nil {
		if strings.Contains(err.Error(), "no documents") {
			w.WriteHeader(http.StatusNotFound)
		} else {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	// Create Person
	// Note: we use the dao method because the controller method will re-hash the password.
	var p PersonController
	personID, err := p.dao.CreatePerson(signup.Person)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Create Family
	signup.Family.Person = []string{personID}
	var f FamilyController

	familyID, err := f.createFamily(signup.Family)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Update Person with Family ID
	err = p.dao.ChangeFamilyID(personID, familyID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Delete signup request
	if err := s.dao.DeleteSignup(code); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// RemoveStaleSignups will remove any signups older than
// the specified time.
func (s *SignupController) RemoveStaleSignups() {
	deleteTime := time.Now().Unix() - 172800 // 172800s‬ == 48hrs

	count, err := s.dao.DeleteStale(deleteTime)
	if err != nil {
		log.Println(err)
	}
	if count > 0 {
		log.Println("[stale signups] Deleted:", count)
	}
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

// SendEmail sends our signup email to the email service.
func (s *SignupController) SendEmail(email, name, code string) error {

	signupRequest := struct {
		Email string `json:"email"`
		Name  string `json:"name"`
		Code  string `json:"code"`
	}{
		Email: email,
		Name:  name,
		Code:  code,
	}

	// Call email-service
	emailRequest, err := json.Marshal(signupRequest)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", "http://email-service:8080/rest/v1/emails/signup", bytes.NewBuffer(emailRequest))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {

		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return errors.New(resp.Status)
	}
	return nil
}
