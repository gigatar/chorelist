package familes

import (
	"chorelist/user-service-gokit/gigatarerrors"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gigatar/chorelist/token"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type getFamilyByIDRequest struct{}
type getFamilyByIDResponse struct {
	Family Family `json:"family,omitempty"`
}

type changeNameRequest struct {
	Family Family `json:"family,omitempty"`
}
type changeNameResponse struct{}

// decoders
func decodeGetFamilyByIDRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var family Family

	vars := mux.Vars(r)
	if _, ok := vars["id"]; !ok {
		return Family{}, gigatarerrors.ErrBadRequest
	}

	inputID := vars["id"]

	// Get Family ID
	var jwt token.JWTToken
	familyID, err := jwt.GetFamily(r.Header.Get("authorization"))
	if err != nil {
		return Family{}, err
	}

	if strings.Compare(inputID, familyID) != 0 {
		return Family{}, gigatarerrors.ErrNotFound
	}

	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		return Family{}, gigatarerrors.ErrBadRequest
	}

	family.ID = id

	fmt.Println(id)
	return family, nil
}
func decodeChangeNameRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var family Family
	if err := json.NewDecoder(r.Body).Decode(&family); err != nil {
		return Family{}, err
	}

	if !family.ValidateName() {
		return Family{}, gigatarerrors.ErrBadRequest
	}

	// Get familyID from JWT
	var jwt token.JWTToken
	familyID, err := jwt.GetFamily(r.Header.Get("authorization"))
	if err != nil {
		return Family{}, gigatarerrors.ErrBadRequest
	}

	family.ID, err = primitive.ObjectIDFromHex(familyID)
	if err != nil {
		return Family{}, gigatarerrors.ErrBadRequest
	}

	return family, nil
}

// encoders
func encodeGetFamilyByIDResponse(ctx context.Context, w http.ResponseWriter, r interface{}) error {
	response := r.(getFamilyByIDResponse)
	if len(response.Family.Name) < 1 {
		return gigatarerrors.ErrNotFound
	}

	return json.NewEncoder(w).Encode(response)
}

func encodeChangeNameResponse(ctx context.Context, w http.ResponseWriter, r interface{}) error {
	w.WriteHeader(http.StatusNoContent)
	return nil
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		log.Fatal("encodeError with nil error!")
	}
	errCode := gigatarerrors.WebErrorCodes(err)
	w.WriteHeader(errCode)
	// if errCode != 500 {
	json.NewEncoder(w).Encode(map[string]interface{}{"error": err.Error()})
	// }
}
