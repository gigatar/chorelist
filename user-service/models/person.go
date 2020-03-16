package models

// Person data structure.
// Note: Type is essentially a string enum
type Person struct {
	ID       string `json:"id" bson:"_id" example:""`
	Name     string `json:"name" bson:"name" example:"John Doe"`
	Email    string `json:"email,omitempty" bson:"email,omitempty" example:"johndoe@gmail.com"`
	Type     string `json:"type" bson:"type" example:"Parent"`
	Password string `json:"password" bson:"password" example:"ABC123"`
	FamilyID string `json:"familyID" bson:"familyID" example:""`
}
