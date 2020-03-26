package families

import (
	"chorelist/user-service-gokit/gigatarerrors"
	"context"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Service defines the beahvior of our famlies service.
type Service interface {
	GetFamilyByID(ctx context.Context, inputFamily Family) (Family, error)
	ChangeName(ctx context.Context, inputFamily Family) error
	CreateFamily(ctx context.Context, inputFamily Family) (interface{}, error)
	// Add Family Member
	// Remove Family Member
	// Delete Family
}

// Family data type
type Family struct {
	ID     primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty" example:""`
	Name   string             `json:"name" bson:"name" example:"Doe"`
	Person []string           `json:"person" bson:"person" example:""`
}

// NewService returns our user service.
func NewService() Service {
	return Family{}
}

// CreateFamily Adds a new family and returns the resourceID.
func (Family) CreateFamily(ctx context.Context, inputFamily Family) (interface{}, error) {
	var db Database
	collection, err := db.GetFamilyCollection(ctx)
	if err != nil {
		return nil, err
	}

	result, err := collection.InsertOne(ctx, inputFamily)
	if err != nil {
		return nil, err
	}

	return result.InsertedID, nil
}

// GetFamilyByID returns a specific family from the database.
func (Family) GetFamilyByID(ctx context.Context, inputFamily Family) (Family, error) {
	var db Database
	collection, err := db.GetFamilyCollection(ctx)
	if err != nil {
		return Family{}, err
	}

	var family Family
	err = collection.FindOne(ctx, bson.M{"_id": inputFamily.ID}).Decode(&family)
	if err != nil {
		if strings.Contains(err.Error(), "no documents") {
			return Family{}, gigatarerrors.ErrNotFound

		}
		return Family{}, err
	}
	return family, nil
}

// ChangeName modifies the name of a family.
func (Family) ChangeName(ctx context.Context, inputFamily Family) error {
	var db Database
	collection, err := db.GetFamilyCollection(ctx)
	if err != nil {
		return err
	}

	// We don't care about the modified count because we want to be idemopotent.
	_, err = collection.UpdateOne(ctx, bson.M{"_id": inputFamily.ID}, bson.M{"$set": bson.M{"name": inputFamily.Name}})
	if err != nil {
		return err
	}

	return nil
}
