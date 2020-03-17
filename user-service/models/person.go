package models

import (
	"regexp"
	"strings"
)

// Person data structure.
// Note: Type is essentially a string enum
type Person struct {
	ID       string `json:"id,omitempty" bson:"_id,omitempty" example:""`
	Name     string `json:"name" bson:"name" example:"John Doe"`
	Email    string `json:"email,omitempty" bson:"email,omitempty" example:"johndoe@gmail.com"`
	Type     string `json:"type" bson:"type" example:"Parent"`
	Password string `json:"password" bson:"password" example:"ABC123"`
	FamilyID string `json:"familyID" bson:"familyID" example:""`
}

// Validate input data for Person
func (p Person) Validate() bool {
	if len(p.Name) < 1 || len(p.Name) > 128 {
		return false
	} else if !p.validateEmail() {
		return false
	} else if !p.validatePassword() {
		return false
	} else if !p.validateType() {
		return false
	}

	return true
}

func (p Person) validateEmail() bool {
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if !emailRegex.MatchString(p.Email) {
		return false
	}
	return true
}

func (p Person) validatePassword() bool {
	if len(p.Password) < 8 || len(p.Password) > 128 {
		return false
	}

	return true
}

func (p Person) validateType() bool {
	if strings.Compare(strings.ToLower(p.Type), "parent") != 0 && strings.Compare(strings.ToLower(p.Type), "child") != 0 {
		return false
	}

	return true
}
