package controllers

import (
	"chorelist/chore-service/daos"
	"chorelist/chore-service/models"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gigatar/chorelist/token"
)

// ChoreController data type.
type ChoreController struct {
	dao daos.ChoreDAO
}

// ListFamilyChores returns all chores from the system.
func (c *ChoreController) ListFamilyChores(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get familyID to ensure we can only family chores.
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

	chores, err := c.dao.GetFamilyChores(familyID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(chores) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(chores)
}

// CreateChore adds a new chore in the system.
func (c *ChoreController) CreateChore(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Deserialize request
	var inputChore models.Chore
	if err := json.NewDecoder(r.Body).Decode(&inputChore); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Validate chore
	if !inputChore.ValidateCreate() {
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
	// Get familyID from JWT
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

	// TODO: Lookup and validate Assigned
	if strings.Compare(strings.ToLower(inputChore.Status), "unassigned") == 0 {
		inputChore.AssigneeID = nil
	}

	// Add Family ID and Creator ID
	inputChore.FamilyID = familyID
	inputChore.CreatorID = userID
	inputChore.DateCreated = time.Now().Unix()

	id, err := c.dao.InsertChore(inputChore)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Location", id)
	w.WriteHeader(http.StatusCreated)
}
