package users

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gorilla/mux"

	"github.com/gigatar/chorelist/token"
)

type getUsersRequest struct{}
type getUsersResponse struct {
	Users []User `json:"users,omitempty"`
}

type getUserByIDRequest struct{}
type getUserByIDResponse struct {
	User User `json:"user,omitempty"`
}

type loginRequest struct {
	Login User `json:"user"`
}
type loginResponse struct {
	Login User `json:"user,omitempty"`
}

type changeNameRequest struct {
	User User `json:"user"`
}
type changeNameResponse struct{}
type changePasswordRequest struct {
	User User `json:"user"`
}
type changePasswordResponse struct{}

// Decoders
func decodeGetUsersRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	// Nothing to really do here...
	var request getUsersRequest
	return request, nil
}

func decodeChangeNameRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		return User{}, err
	}
	if !user.ValidateName() {
		return User{}, errBadRequest
	}

	// Get userID from JWT
	var jwt token.JWTToken
	userID, err := jwt.GetUser(r.Header.Get("authorization"))
	if err != nil {
		return User{}, err
	}

	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return User{}, errBadRequest
	}
	user.ID = id

	return user, nil
}
func decodeChangePasswordRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		return User{}, err
	}
	if !user.ValidatePassword() {
		return User{}, errBadRequest
	}

	// Get userID from JWT
	var jwt token.JWTToken
	userID, err := jwt.GetUser(r.Header.Get("authorization"))
	if err != nil {
		return User{}, err
	}

	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return User{}, errBadRequest
	}
	user.ID = id

	return user, nil
}

func decodeGetUserByIDRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var user User

	vars := mux.Vars(r)
	if _, ok := vars["id"]; !ok {
		return User{}, errBadRequest
	}

	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		return User{}, errBadRequest
	}
	user.ID = id

	// Get Family ID
	var jwt token.JWTToken
	familyID, err := jwt.GetFamily(r.Header.Get("authorization"))
	if err != nil {
		return User{}, err
	}
	user.FamilyID = familyID

	return user, nil
}

func decodeLoginRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		return nil, err
	}

	return user, nil
}

// Encoders
func encodeGetUsersResponse(ctx context.Context, w http.ResponseWriter, r interface{}) error {
	response := r.(getUsersResponse)
	if len(response.Users) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return nil
	}
	return json.NewEncoder(w).Encode(response)
}

func encodeGetUserByIDResponse(ctx context.Context, w http.ResponseWriter, r interface{}) error {
	response := r.(getUserByIDResponse)
	if len(response.User.Name) < 1 {
		return errNotFound
	}

	return json.NewEncoder(w).Encode(response)
}

func encodeLoginResponse(ctx context.Context, w http.ResponseWriter, r interface{}) error {
	response := r.(loginResponse)

	id, err := response.Login.ID.MarshalJSON()
	idString := strings.Trim(string(id), "\"")
	var jwt token.JWTToken
	token, err := jwt.CreateJWT("127.0.0.1:80", idString, response.Login.FamilyID)
	if err != nil {
		return err
	}

	w.Header().Add("Authorization", token)

	return json.NewEncoder(w).Encode(response.Login)
}

func encodeChangeNameResponse(ctx context.Context, w http.ResponseWriter, r interface{}) error {
	w.WriteHeader(http.StatusNoContent)
	return nil
}

func encodeChangePasswordResponse(ctx context.Context, w http.ResponseWriter, r interface{}) error {
	w.WriteHeader(http.StatusNoContent)
	return nil
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		log.Fatal("encodeError with nil error!")
	}
	errCode := errorCodes(err)
	w.WriteHeader(errCode)
	// if errCode != 500 {
	json.NewEncoder(w).Encode(map[string]interface{}{"error": err.Error()})
	// }
}

func errorCodes(err error) int {
	switch err {

	case errBadRequest:
		return http.StatusBadRequest
	case errDuplicate:
		return http.StatusConflict
	case errNotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}

}
