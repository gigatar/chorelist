package families

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	databaseName         = "ChoreList_UserService"
	familyCollectionName = "families"
)

// Database is our instance.
type Database struct {
	Client  *mongo.Client
	Timeout time.Duration
}

// Init is how we initialize our database instance.
func (db *Database) Init(ctx context.Context) error {
	clientOptions := options.Client().ApplyURI("mongodb://mongodb:27017")
	var err error

	db.Client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}

	return nil
}

// GetFamilyCollection is a helper function to dial and return our family collection.
func (db *Database) GetFamilyCollection(ctx context.Context) (*mongo.Collection, error) {
	clientOptions := options.Client()
	clientOptions.ApplyURI("mongodb://mongodb:27017")

	var err error

	db.Client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return db.Client.Database(databaseName).Collection(familyCollectionName), nil
}
