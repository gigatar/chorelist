package users

import (
	"regexp"
	"strings"
)

// Validate input data for Person.
func (u User) Validate() bool {
	if !u.ValidateName() {
		return false
	} else if !u.ValidateEmail() {
		return false
	} else if !u.ValidatePassword() {
		return false
	} else if !u.ValidateType() {
		return false
	}

	return true
}

// ValidateEmail validates the Person.Email field.
func (u User) ValidateEmail() bool {
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if !emailRegex.MatchString(u.Email) {
		return false
	}
	return true
}

// ValidatePassword validates the Person.Password field.
func (u User) ValidatePassword() bool {
	if len(u.Password) < 8 || len(u.Password) > 128 {
		return false
	}

	return true
}

// ValidateType validates the Person.Type field.
func (u User) ValidateType() bool {
	if strings.Compare(strings.ToLower(u.Type), "parent") != 0 && strings.Compare(strings.ToLower(u.Type), "child") != 0 {
		return false
	}

	return true
}

// ValidateName validates the Person.Name field.
func (u User) ValidateName() bool {
	if len(u.Name) < 1 || len(u.Name) > 128 {
		return false
	}
	return true
}

// StripSensitive removes all sensitive variables from Person.
func (u *User) StripSensitive() {
	u.Password = ""
	u.OldPassword = ""
}
