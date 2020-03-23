package daos

import (
	"chorelist/chore-service/database"
	"chorelist/chore-service/models"
	"context"
	"encoding/json"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

// ChoreDAO is the chore data access object.
type ChoreDAO struct{}

// GetFamilyChores returns all database entries for family chores.
// TODO: pagination
func (c *ChoreDAO) GetFamilyChores(familyID string) ([]models.Chore, error) {

	collection := database.DB.GetChoreCollection()

	ctx, cancel := context.WithTimeout(context.Background(), database.DB.Timeout)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{"familyID": familyID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var chores []models.Chore
	for cursor.Next(ctx) {
		var chore models.Chore
		if err := cursor.Decode(&chore); err != nil {
			return nil, err
		}

		chores = append(chores, chore)
	}

	return chores, nil
}

// InsertChore into the database.
func (c *ChoreDAO) InsertChore(chore models.Chore) (string, error) {
	collection := database.DB.GetChoreCollection()

	ctx, cancel := context.WithTimeout(context.Background(), database.DB.Timeout)
	defer cancel()

	resource, err := collection.InsertOne(ctx, chore)
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
