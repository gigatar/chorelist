package database

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

// DB is our global database variable used across the service.
// Note that this is initialized in the main package.
var DB Database

const (
	databaseName         = "ChoreList"
	personCollectionName = "persons"
)

// Database is our instance.
type Database struct {
	Client  *mongo.Client
	Timeout time.Duration
}

// Init is how we initialize our database instance.
func (db *Database) Init() error {
	clientOptions := options.Client().ApplyURI("mongodb://mongodb:27017")
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	db.Client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}

	db.Timeout = time.Second * 5 // Default of 5 seconds

	// Create indexes
	if err := db.createUniquePersonIndex(); err != nil {
		return errors.New("createUserIndex:" + err.Error())
	}

	return nil
}

// createUniquePersonIndex ensures that all emails are unique.
func (db *Database) createUniquePersonIndex() error {
	// Create the required index
	collection := DB.GetPersonCollection()

	options := options.Index()
	options.SetUnique(true)
	options.SetName("uniqueEmail")
	index := mongo.IndexModel{
		Keys:    bsonx.Doc{{Key: "email", Value: bsonx.Int32(1)}},
		Options: options,
	}

	ctx, cancel := context.WithTimeout(context.Background(), DB.Timeout)
	defer cancel()

	_, err := collection.Indexes().CreateOne(ctx, index)
	if err != nil {
		return err
	}

	return nil
}

// GetPersonCollection is a helper function to return our person collection.
func (db *Database) GetPersonCollection() *mongo.Collection {
	return DB.Client.Database(databaseName).Collection(personCollectionName)
}
