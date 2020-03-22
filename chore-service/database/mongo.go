package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DB is our global database variable used across the service.
// Note that this is initialized in the main package.
var DB Database

const (
	databaseName        = "ChoreList_ChoreService"
	choreCollectionName = "chores"
	noteCollectionName  = "notes"
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

	return nil
}

// GetChoreCollection is a helper function to return our chore collection.
func (db *Database) GetChoreCollection() *mongo.Collection {
	return DB.Client.Database(databaseName).Collection(choreCollectionName)
}

// GetNoteCollection is a helper function to return our Note collection.
func (db *Database) GetNoteCollection() *mongo.Collection {
	return DB.Client.Database(databaseName).Collection(noteCollectionName)
}
