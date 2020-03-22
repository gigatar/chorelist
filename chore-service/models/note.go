package models

// Note Data Structure
type Note struct {
	ID          string `json:"id,omitempty" bson:"_id,omitempty" example:""`
	FamilyID    string `json:"familyID,omitempty" bson:"familyID,omitempty" example:""`
	Text        string `json:"text,omitempty" bson:"text,omitempty" example:"I am note"`
	CreatorID   string `json:"creatorID,omitempty" bson:"creatorID,omitempty" example:""`
	DateCreated int64  `json:"dateCreated,omitempty" bson:"dateCreated,omitempty" example:"1584588677"`
}

// Validate note input fields.
func (n Note) Validate() bool {
	if !n.ValidateText() {
		return false
	}

	return true
}

// ValidateText is within length bounds.
func (n Note) ValidateText() bool {
	if len(n.Text) < 1 || len(n.Text) > 4096 {
		return false
	}
	return true
}
