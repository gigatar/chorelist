package daos

import (
	"chorelist/user-service/database"
	"chorelist/user-service/models"
	"context"
	"encoding/json"
	"errors"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"

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

// ChangeName updates a person's name in the database.
func (p *PersonDAO) ChangeName(userID, newName string) error {

	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}
	collection := database.DB.GetPersonCollection()

	ctx, cancel := context.WithTimeout(context.Background(), database.DB.Timeout)
	defer cancel()

	// We don't need to do anything with the result.
	_, err = collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"name": newName}})
	if err != nil {
		return err
	}

	return nil
}

// GetEncryptedPassword for person from database for password change.
func (p *PersonDAO) GetEncryptedPassword(userID string) (string, error) {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return "", err
	}

	collection := database.DB.GetPersonCollection()

	ctx, cancel := context.WithTimeout(context.Background(), database.DB.Timeout)
	defer cancel()

	findOptions := options.FindOne()
	findOptions.SetProjection(bson.M{"password": 1})

	var person models.Person
	err = collection.FindOne(ctx, bson.M{"_id": id}, findOptions).Decode(&person)
	if err != nil {
		return "", err
	}

	return person.Password, nil
}

// UpdatePassword updates a users password.
func (p *PersonDAO) UpdatePassword(userID, hashedPassword string) error {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	collection := database.DB.GetPersonCollection()

	ctx, cancel := context.WithTimeout(context.Background(), database.DB.Timeout)
	defer cancel()

	// Don't need result, no documents is implicit in the error.
	_, err = collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"password": hashedPassword}})
	if err != nil {
		return err
	}

	return nil
}

// DeletePerson removes them from the database
func (p *PersonDAO) DeletePerson(userID string) error {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	collection := database.DB.GetPersonCollection()

	ctx, cancel := context.WithTimeout(context.Background(), database.DB.Timeout)
	defer cancel()

	result, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	if result.DeletedCount != 1 {
		return errors.New("no documents")
	}

	return nil
}

// GetPersonType returns the type field for a user.
func (p *PersonDAO) GetPersonType(userID string) (string, error) {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return "", err
	}

	collection := database.DB.GetPersonCollection()

	ctx, cancel := context.WithTimeout(context.Background(), database.DB.Timeout)
	defer cancel()

	var person models.Person

	findOptions := options.FindOne()
	findOptions.SetProjection(bson.M{"type": 1})

	err = collection.FindOne(ctx, bson.M{"_id": id}).Decode(&person)
	if err != nil {
		return "", err
	}

	return person.Type, nil
}

// EmailExists checks if an email address exists in the system.
func (p *PersonDAO) EmailExists(email string) (bool, error) {
	collection := database.DB.GetPersonCollection()

	ctx, cancel := context.WithTimeout(context.Background(), database.DB.Timeout)
	defer cancel()

	count, err := collection.CountDocuments(ctx, bson.M{"email": email})
	if err != nil {
		return true, err
	}

	if count == 0 {
		return false, nil
	}

	return true, nil
}

// ChangeFamilyID updates a person's name in the database.
func (p *PersonDAO) ChangeFamilyID(userID, familyID string) error {

	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}
	collection := database.DB.GetPersonCollection()

	ctx, cancel := context.WithTimeout(context.Background(), database.DB.Timeout)
	defer cancel()

	// We don't need to do anything with the result.
	_, err = collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"familyID": familyID}})
	if err != nil {
		return err
	}

	return nil
}

// UpdateLastLogin updates a person's last login time.
func (p *PersonDAO) UpdateLastLogin(userID string, loginTime int64) error {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	collection := database.DB.GetPersonCollection()

	ctx, cancel := context.WithTimeout(context.Background(), database.DB.Timeout)
	defer cancel()

	_, err = collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"lastLogin": loginTime}})
	if err != nil {
		return err
	}

	return nil
}
