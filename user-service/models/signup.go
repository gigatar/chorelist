package models

// Signup data model.
type Signup struct {
	Person     Person `json:"person" bson:"person"`
	Family     Family `json:"family" bson:"family"`
	SignupTime int64  `json:"signupTime,omitempty" bson:"signupTime,omitempty"`
	Code       string `json:"code,omitempty" bson:"code,omitempty"`
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
