package users

import (
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

const (
	databaseName         = "ChoreList_UserService"
	personCollectionName = "persons"
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

	// Create indexes
	if err := db.createUniquePersonIndex(ctx); err != nil {
		return errors.New("createUserIndex:" + err.Error())
	}
	return nil
}

// createUniquePersonIndex ensures that all emails are unique.
func (db *Database) createUniquePersonIndex(ctx context.Context) error {
	// Create the required index
	collection := db.Client.Database(databaseName).Collection(personCollectionName)

	options := options.Index()
	options.SetUnique(true)
	options.SetName("uniqueEmail")
	index := mongo.IndexModel{
		Keys:    bsonx.Doc{{Key: "email", Value: bsonx.Int32(1)}},
		Options: options,
	}

	_, err := collection.Indexes().CreateOne(ctx, index)
	if err != nil {
		return err
	}

	return nil
}

// GetPersonCollection is a helper function to dial and return our person collection.
func (db *Database) GetPersonCollection(ctx context.Context) (*mongo.Collection, error) {
	clientOptions := options.Client()
	clientOptions.ApplyURI("mongodb://localhost:27017")

	var err error

	db.Client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return db.Client.Database(databaseName).Collection(personCollectionName), nil
}
