package models

// Family data structure.
// Note: Person will be a HATEOAS location.
type Family struct {
	ID     string   `json:"id" bson:"_id" example:""`
	Name   string   `json:"name" bson:"name" example:"Doe"`
	Person []Person `json:"person" bson:"person" example:""`
}
