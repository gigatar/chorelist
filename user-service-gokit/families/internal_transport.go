package families

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

type createFamilyRequest struct {
	Family Family `json:"family"`
}
type createFamilyResponse struct {
	Location primitive.ObjectID `json:"location,omitempty"`
}

func decodeCreateFamilyRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var family Family
	if err := json.NewDecoder(r.Body).Decode(&family); err != nil {
		return nil, err
	}

	if !family.Validate() {
		log.Println("Fail name validation")
		return nil, gigatarerrors.ErrBadRequest
	}

	return family, nil
}

func encodeCreateFamilyResponse(ctx context.Context, w http.ResponseWriter, r interface{}) error {
	response := r.(createFamilyResponse)
	if len(response.Location) != len(primitive.NilObjectID) {
		return errors.New("Invalid Location")
	}

	// Convert Location to string (I know this is dirty and should be broken apart...)
	location := strings.TrimSuffix(strings.TrimPrefix(response.Location.String(), "ObjectID(\""), "\")")

	w.Header().Add("Location", location)
	w.WriteHeader(http.StatusCreated)

	return nil
}

type deleteFamilyRequest struct {
	Family Family `json:"family"`
}
type deleteFamilyResponse struct{}

func decodeDeleteFamilyRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var family Family

	vars := mux.Vars(r)
	if _, ok := vars["id"]; !ok {
		return Family{}, gigatarerrors.ErrBadRequest
	}

	inputID, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		return nil, gigatarerrors.ErrBadRequest
	}

	family.ID = inputID

	return family, nil
}

func encodeDeleteFamilyResponse(ctx context.Context, w http.ResponseWriter, r interface{}) error {
	w.WriteHeader(http.StatusNoContent)
	return nil
}
