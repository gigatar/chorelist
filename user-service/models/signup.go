package models

import "time"

// Signup data model.
type Signup struct {
	Person Person    `json:"person" bson:"person"`
	Family Family    `json:"family" bson:"family"`
	Expire time.Time `json:"expire,omitempty" bson:"expire,omitempty"`
}

// Validate signup information.
func (s Signup) Validate() bool {
	if !s.Person.Validate() {
		return false
	} else if !s.Family.Validate() {
		return false
	}

	return true
}
