package families

// Validate input for family.
func (f Family) Validate() bool {
	if !f.ValidateName() {
		return false
	}

	return true
}

// ValidateName input for family.
func (f Family) ValidateName() bool {
	if len(f.Name) < 1 || len(f.Name) > 128 {
		return false
	}

	return true
}

// ValidateAddPerson ensures that family size will not be too large.
func (f Family) ValidateAddPerson() bool {
	if len(f.Person) >= 15 {
		return false
	}

	return true
}
