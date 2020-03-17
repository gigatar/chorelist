package daos

import (
	"chorelist/user-service/database"
	"chorelist/user-service/models"
	"context"
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
