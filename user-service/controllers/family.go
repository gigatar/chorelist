package controllers

import (
	"chorelist/user-service/daos"
	"chorelist/user-service/models"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"github.com/gigatar/chorelist/token"
)

// FamilyController data type.
type FamilyController struct {
	dao daos.FamilyDAO
}

// DeleteFamily deletes a family and all Persons assigned.
func (f *FamilyController) DeleteFamily(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get FamilyID to ensure we can only modify our family
	var jwt token.JWTToken
	familyID, err := jwt.GetFamily(r.Header.Get("authorization"))
	if err != nil {
		if strings.Contains(err.Error(), "Invalid token") {
			w.WriteHeader(http.StatusUnauthorized)
		} else {
			log.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
		}
		return
	}

	if familyID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get userID to ensure we are a parent.
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

	// Check person is Parent.
	var p PersonController
	userType, err := p.getPersonType(userID)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if strings.Compare(strings.ToLower(userType), "parent") != 0 {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Get family details
	family, err := f.dao.GetFamily(familyID)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Delete all family persons.
	for _, person := range family.Person {
		err = p.dao.DeletePerson(person)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	// TODO: Delete chores

	// Delete Family
	err = f.dao.DeleteFamily(familyID)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (f *FamilyController) createFamily(family models.Family) (string, error) {

	if !family.Validate() {
		return "", errors.New("Invalid input")
	}

	id, err := f.dao.CreateFamily(family)
	if err != nil {
		return "", err
	}
	return id, nil
}

// AddFamilyMember adds a new family member to the system.
func (f *FamilyController) AddFamilyMember(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Decode person input
	var inputPerson models.Person
	err := json.NewDecoder(r.Body).Decode(&inputPerson)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Validate input
	if !inputPerson.Validate() {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get userID from JWT
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

	// Get FamilyID from JWT
	familyID, err := jwt.GetFamily(r.Header.Get("authorization"))
	if err != nil {
		if strings.Contains(err.Error(), "Invalid token") {
			w.WriteHeader(http.StatusUnauthorized)
		} else {
			log.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
		}
		return
	}

	if familyID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Make sure Parent
	var p PersonController
	personType, err := p.getPersonType(userID)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if strings.Compare(strings.ToLower(personType), "parent") != 0 {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Add familyID to person
	inputPerson.FamilyID = familyID

	// Create person
	newPersonID, err := p.createPerson(inputPerson)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			w.WriteHeader(http.StatusConflict)
		} else {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	// Get Family
	family, err := f.dao.GetFamily(familyID)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Add person to family
	family.Person = append(family.Person, newPersonID)
	err = f.dao.UpdateFamilyMember(family)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Return success
	w.Header().Add("Location", newPersonID)
	w.WriteHeader(http.StatusNoContent)
}

// RemoveFamilyMember removeos a family member and deletes the person.
func (f *FamilyController) RemoveFamilyMember(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	personID := mux.Vars(r)["personID"]
	if len(personID) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get userID from JWT
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

	// Get FamilyID from JWT
	familyID, err := jwt.GetFamily(r.Header.Get("authorization"))
	if err != nil {
		if strings.Contains(err.Error(), "Invalid token") {
			w.WriteHeader(http.StatusUnauthorized)
		} else {
			log.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
		}
		return
	}

	if familyID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Make sure Parent
	var p PersonController
	personType, err := p.getPersonType(userID)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if strings.Compare(strings.ToLower(personType), "parent") != 0 {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Make sure not removing self.
	if strings.Compare(personID, userID) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return

	}

	// Get Family
	family, err := f.dao.GetFamily(familyID)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Remove the person from the family
	for i := 0; i < len(family.Person); i++ {
		if family.Person[i] == personID {
			family.Person = append(family.Person[:i], family.Person[i+1:]...)
		}
	}

	err = f.dao.UpdateFamilyMember(family)
	if err != nil {
		if strings.Contains(err.Error(), "Family members not updated") {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	// Delete person from system.
	err = p.dao.DeletePerson(personID)
	if err != nil {
		if strings.Contains(err.Error(), "no documents") {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ViewFamily returns the family object.
func (f *FamilyController) ViewFamily(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get familyID to ensure we can only view our family.
	var jwt token.JWTToken
	familyID, err := jwt.GetFamily(r.Header.Get("authorization"))
	if err != nil {
		if strings.Contains(err.Error(), "Invalid token") {
			w.WriteHeader(http.StatusUnauthorized)
		} else {
			log.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
		}
		return
	}

	if familyID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	family, err := f.dao.GetFamily(familyID)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(family)
}

// ChangeName changes the name of a family.
func (f *FamilyController) ChangeName(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Deserialize request
	var inputFamily models.Family
	err := json.NewDecoder(r.Body).Decode(&inputFamily)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Validate input
	if !inputFamily.ValidateName() {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get userID to ensure we are a parent
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

	// Get familyID to ensure we can only modify our family.
	familyID, err := jwt.GetFamily(r.Header.Get("authorization"))
	if err != nil {
		if strings.Contains(err.Error(), "Invalid token") {
			w.WriteHeader(http.StatusUnauthorized)
		} else {
			log.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
		}
		return
	}

	if familyID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var p PersonController
	personType, err := p.getPersonType(userID)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if strings.Compare(strings.ToLower(personType), "parent") != 0 {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Change name
	err = f.dao.ChangeName(familyID, inputFamily.Name)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
