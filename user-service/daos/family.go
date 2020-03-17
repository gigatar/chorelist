package daos

import (
	"chorelist/user-service/database"
	"chorelist/user-service/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// FamilyDAO is the Family data access object.
type FamilyDAO struct{}

// CreateFamily in database.
func (f *FamilyDAO) CreateFamily(family models.Family) error {

	collection := database.DB.GetFamilyCollection()

	ctx, cancel := context.WithTimeout(context.Background(), database.DB.Timeout)
	defer cancel()

	_, err := collection.InsertOne(ctx, family)
	if err != nil {
		return err
	}

	return nil
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
