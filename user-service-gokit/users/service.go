package users

import (
	"chorelist/user-service-gokit/gigatarerrors"
	"context"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"

	"go.mongodb.org/mongo-driver/bson"
)

// Service defines the beahvior of our users service.
type Service interface {
	GetUsers(ctx context.Context, familyID string) ([]User, error)
	GetUserByID(ctx context.Context, inputUser User) (User, error)
	Login(ctx context.Context, inputUser User) (User, error)
	ChangeName(ctx context.Context, inputUser User) error
	ChangePassword(ctx context.Context, inputUser User) error
	CreateUser(ctx context.Context, inputUser User) (interface{}, error)
	DeleteUser(ctx context.Context, inputUser User) error
}

const bcryptPasswordCost = 10

// var (
// 	errNotFound   = errors.New("not found")
// 	errDuplicate  = errors.New("duplicate")
// 	errBadRequest = errors.New("bad request")
// )

// User Datatype
type User struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty" example:""`
	Name        string             `json:"name" bson:"name" example:"John Doe"`
	Email       string             `json:"email,omitempty" bson:"email,omitempty" example:"johndoe@gmail.com"`
	Type        string             `json:"type" bson:"type" example:"Parent"`
	Password    string             `json:"password,omitempty" bson:"password" example:"ABC123"`
	OldPassword string             `json:"oldPassword,omitempty" bson:"oldPassword,omitempty" example:"ABC123"`
	FamilyID    string             `json:"familyID,omitempty" bson:"familyID,omitempty" example:""`
	LastLogin   int64              `json:"lastLogin,omitempty" bson:"lastLogin,omitempty" example:"1584588677"`
}

// NewService returns our user service.
func NewService() Service {
	return User{}
}

// DeleteUser removes a user from the system
func (User) DeleteUser(ctx context.Context, inputUser User) error {
	var db Database
	collection, err := db.GetPersonCollection(ctx)
	if err != nil {
		return err
	}

	result, err := collection.DeleteOne(ctx, bson.M{"_id": inputUser.ID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return gigatarerrors.ErrNotFound
	}

	return nil
}

// CreateUser creates a new user and returns the location
func (User) CreateUser(ctx context.Context, inputUser User) (interface{}, error) {

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(inputUser.Password), bcryptPasswordCost)
	if err != nil {
		return nil, err
	}
	inputUser.Password = string(hash)

	var db Database
	collection, err := db.GetPersonCollection(ctx)
	if err != nil {
		return nil, err
	}

	result, err := collection.InsertOne(ctx, inputUser)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return nil, gigatarerrors.ErrDuplicate
		}
		return nil, err
	}

	return result.InsertedID, nil
}

// ChangePassword replaces the password of a user.
func (User) ChangePassword(ctx context.Context, inputUser User) error {
	var db Database
	collection, err := db.GetPersonCollection(ctx)
	if err != nil {
		return err
	}

	// Get user information from database.
	var current User
	err = collection.FindOne(ctx, bson.M{"_id": inputUser.ID}).Decode(&current)
	if err != nil {
		if strings.Contains(err.Error(), "no documents") {
			return gigatarerrors.ErrBadRequest
		}
		return err
	}

	// Check Hash
	err = bcrypt.CompareHashAndPassword([]byte(current.Password), []byte(inputUser.OldPassword))
	if err != nil {
		return gigatarerrors.ErrBadRequest
	}

	password, err := bcrypt.GenerateFromPassword([]byte(inputUser.Password), bcryptPasswordCost)
	if err != nil {
		return err
	}

	// We don't care about the modified count because it will always mutate (hopefully if bcrypt works...).
	_, err = collection.UpdateOne(ctx, bson.M{"_id": inputUser.ID}, bson.M{"$set": bson.M{"password": string(password)}})
	if err != nil {
		return err
	}

	return nil
}

// ChangeName modifies the name of a user.
func (User) ChangeName(ctx context.Context, inputUser User) error {
	var db Database
	collection, err := db.GetPersonCollection(ctx)
	if err != nil {
		return err
	}

	// We don't care about the modified count because we want to be idemopotent.
	_, err = collection.UpdateOne(ctx, bson.M{"_id": inputUser.ID}, bson.M{"$set": bson.M{"name": inputUser.Name}})
	if err != nil {
		return err
	}

	return nil
}

// GetUsers returns all family users from our system.
func (User) GetUsers(ctx context.Context, familyID string) ([]User, error) {
	var db Database

	collection, err := db.GetPersonCollection(ctx)
	if err != nil {
		return nil, err
	}

	cursor, err := collection.Find(ctx, bson.M{"familyID": familyID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []User
	for cursor.Next(ctx) {
		var user User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		user.StripSensitive()
		users = append(users, user)
	}

	if len(users) == 0 {
		return nil, gigatarerrors.ErrNotFound
	}
	return users, nil
}

// GetUserByID returns a specific user from our database.
func (User) GetUserByID(ctx context.Context, inputUser User) (User, error) {
	var db Database

	collection, err := db.GetPersonCollection(ctx)
	if err != nil {
		return User{}, err
	}

	var user User
	err = collection.FindOne(ctx, bson.M{"_id": inputUser.ID, "familyID": inputUser.FamilyID}).Decode(&user)
	if err != nil {
		if strings.Contains(err.Error(), "no documents") {
			return User{}, gigatarerrors.ErrNotFound
		}
		return User{}, err
	}

	user.StripSensitive()
	return user, nil
}

// Login checks a users authentication and returns user if good.
func (User) Login(ctx context.Context, inputUser User) (User, error) {
	var db Database

	// Get person from database
	collection, err := db.GetPersonCollection(ctx)
	if err != nil {
		return User{}, err
	}
	var user User
	err = collection.FindOne(ctx, bson.M{"email": inputUser.Email}).Decode(&user)
	if err != nil {
		return User{}, gigatarerrors.ErrBadRequest
	}

	// Validate password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(inputUser.Password)); err != nil {
		return User{}, gigatarerrors.ErrBadRequest
	}

	user.StripSensitive()
	return user, nil
}
