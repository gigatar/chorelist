package models

// Family data structure.
// Note: Person will be a HATEOAS location.
type Family struct {
	ID     string   `json:"id,omitempty" bson:"_id,omitempty" example:""`
	Name   string   `json:"name" bson:"name" example:"Doe"`
	Person []string `json:"person" bson:"person" example:""`
}

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
