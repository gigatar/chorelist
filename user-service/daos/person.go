package daos

import (
	"chorelist/user-service/database"
	"chorelist/user-service/models"
	"context"
	"encoding/json"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
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

// Login gets the needed details for a login.
func (p *PersonDAO) Login(email string) (models.Person, error) {

	var person models.Person

	collection := database.DB.GetPersonCollection()

	ctx, cancel := context.WithTimeout(context.Background(), database.DB.Timeout)
	defer cancel()

	// We'll keep the id and pass it back to the user.
	findOptions := options.FindOne()
	findOptions.SetProjection(bson.M{"email": 0})

	err := collection.FindOne(ctx, bson.M{"email": email}, findOptions).Decode(&person)
	if err != nil {
		return person, err
	}

	return person, nil
}
