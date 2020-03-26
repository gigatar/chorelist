package users

import (
	"chorelist/user-service-gokit/gigatarerrors"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type createUserRequest struct {
	User User `json:"user"`
}
type createUserResponse struct {
	Location primitive.ObjectID `json:"location,omitempty"`
}

func decodeCreateUserRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		return nil, err
	}

	if !user.Validate() {
		return nil, gigatarerrors.ErrBadRequest
	}

	return user, nil
}
func encodeCreateUserResponse(ctx context.Context, w http.ResponseWriter, r interface{}) error {
	response := r.(createUserResponse)
	if len(response.Location) != len(primitive.NilObjectID) {
		return errors.New("Invalid Location")
	}

	// Convert Location to string (I know this is dirty and should be broken apart...)
	location := strings.TrimSuffix(strings.TrimPrefix(response.Location.String(), "ObjectID(\""), "\")")

	w.Header().Add("Location", location)
	w.WriteHeader(http.StatusCreated)

	return nil
}

type deleteUserRequest struct{}
type deleteUserResponse struct{}

func decodeDeleteUserRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var user User

	vars := mux.Vars(r)
	if _, ok := vars["id"]; !ok {
		return User{}, gigatarerrors.ErrBadRequest
	}

	inputID, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		return nil, gigatarerrors.ErrBadRequest
	}

	log.Println(inputID)
	user.ID = inputID

	return user, nil
}

func encodeDeleteUserResponse(ctx context.Context, w http.ResponseWriter, r interface{}) error {
	w.WriteHeader(http.StatusNoContent)
	return nil
}
