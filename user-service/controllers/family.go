package controllers

import (
	"chorelist/user-service/daos"
	"chorelist/user-service/models"
	"errors"
	"log"
	"net/http"
	"strings"

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

	// Get userID to ensure we can only modify ourselves
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
		w.WriteHeader(http.StatusBadRequest)
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
		strippedPerson := person[len(userID):] // Strip HATEOAS location off front.
		err = p.dao.DeletePerson(strippedPerson)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	// Delete Family
	err = f.dao.DeleteFamily(familyID)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (f *FamilyController) createFamily(family models.Family) error {

	if !family.Validate() {
		return errors.New("Invalid input")
	}

	if err := f.dao.CreateFamily(family); err != nil {
		return err
	}
	return nil
}
