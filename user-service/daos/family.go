package daos

import (
	"chorelist/user-service/database"
	"chorelist/user-service/models"
	"context"
	"encoding/json"
	"errors"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// FamilyDAO is the Family data access object.
type FamilyDAO struct{}

// CreateFamily in database.
func (f *FamilyDAO) CreateFamily(family models.Family) (string, error) {

	collection := database.DB.GetFamilyCollection()

	ctx, cancel := context.WithTimeout(context.Background(), database.DB.Timeout)
	defer cancel()

	resource, err := collection.InsertOne(ctx, family)
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

// DeleteFamily from database
func (f *FamilyDAO) DeleteFamily(familyID string) error {
	id, err := primitive.ObjectIDFromHex(familyID)
	if err != nil {
		return err
	}

	collection := database.DB.GetFamilyCollection()

	ctx, cancel := context.WithTimeout(context.Background(), database.DB.Timeout)
	defer cancel()

	_, err = collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	return nil
}

// GetFamily returns the family details
func (f *FamilyDAO) GetFamily(familyID string) (models.Family, error) {
	var family models.Family

	id, err := primitive.ObjectIDFromHex(familyID)
	if err != nil {
		return family, err
	}

	collection := database.DB.GetFamilyCollection()

	ctx, cancel := context.WithTimeout(context.Background(), database.DB.Timeout)
	defer cancel()

	err = collection.FindOne(ctx, bson.M{"_id": id}).Decode(&family)
	if err != nil {
		return family, err
	}

	return family, nil
}

// UpdateFamilyMember updates a family in the database.
func (f *FamilyDAO) UpdateFamilyMember(family models.Family) error {
	id, err := primitive.ObjectIDFromHex(family.ID)
	if err != nil {
		return err
	}

	collection := database.DB.GetFamilyCollection()

	ctx, cancel := context.WithTimeout(context.Background(), database.DB.Timeout)
	defer cancel()

	result, err := collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"person": family.Person}})
	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		return errors.New("Family members not updated")
	}

	return nil
}

// ChangeName updates the family name
func (f *FamilyDAO) ChangeName(familyID, newName string) error {
	id, err := primitive.ObjectIDFromHex(familyID)
	if err != nil {
		return err
	}
	collection := database.DB.GetFamilyCollection()

	ctx, cancel := context.WithTimeout(context.Background(), database.DB.Timeout)
	defer cancel()

	_, err = collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"name": newName}})
	if err != nil {
		return err
	}

	return nil
}
