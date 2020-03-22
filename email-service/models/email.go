package models

import "regexp"

// EmailConfig is how the server connection is configured.
type EmailConfig struct {
	BaseLink   string `json:"baseLink"`
	From       string `json:"from"`
	Server     string `json:"server"`
	ServerPort int    `json:"serverPort"`
	User       string `json:"user"`
	Password   string `json:"password"`
}

// SignupEmail data type
type SignupEmail struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	Code  string `json:"code"`
}

// Validate validates the signup email data type.
func (s SignupEmail) Validate() bool {
	if len(s.Name) < 1 || len(s.Name) > 128 {
		return false
	} else if !s.ValidateEmail() {
		return false
	} else if len(s.Code) != 15 {
		return false
	}

	return true
}

// ValidateEmail validates the SignupEmail.Email field.
func (s SignupEmail) ValidateEmail() bool {
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if !emailRegex.MatchString(s.Email) {
		return false
	}
	return true
}
