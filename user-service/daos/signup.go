package daos

import (
	"chorelist/user-service/database"
	"chorelist/user-service/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

// SignupDAO is the signup data access object
type SignupDAO struct{}

// CreateSignup adds a new signup request to the database.
func (s *SignupDAO) CreateSignup(signup models.Signup) error {

	collection := database.DB.GetSignupCollection()

	ctx, cancel := context.WithTimeout(context.Background(), database.DB.Timeout)
	defer cancel()

	_, err := collection.InsertOne(ctx, signup)
	if err != nil {
		return err
	}
	return nil
}

// DeleteStale removes stale signups from the database.
func (s *SignupDAO) DeleteStale(expire int64) (int64, error) {
	collection := database.DB.GetSignupCollection()

	ctx, cancel := context.WithTimeout(context.Background(), database.DB.Timeout)
	defer cancel()

	deleteResult, err := collection.DeleteMany(ctx, bson.M{"signupTime": bson.M{"$lt": expire}})
	if err != nil {
		return 0, err
	}

	return deleteResult.DeletedCount, nil
}
