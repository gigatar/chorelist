package daos

import (
	"chorelist/user-service/database"
	"chorelist/user-service/models"
	"context"
	"encoding/json"
	"strings"
)

// PersonDAO is the person data access object.
type PersonDAO struct{}

// CreatePerson creates a new person in the database.
func (p *PersonDAO) CreatePerson(person models.Person) (string, error) {
	collection := database.DB.GetPersonCollection()

	ctx, cancel := context.WithTimeout(context.Background(), database.DB.Timeout)
	defer cancel()

	resource, err := collection.InsertOne(ctx, person)
	if err != nil {
		return "", err
	}

	js, err := json.Marshal(resource.InsertedID)
	if err != nil {
		return "", err
	}
	// Strip quotes from string
	ret := strings.Replace(string(js), "\"", "", -1)
	return ret, nil
}
