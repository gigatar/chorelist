package controllers

import (
	"chorelist/user-service/database"
	"chorelist/user-service/models"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func init() {
	// Initialize Database
	if err := database.DB.Init(); err != nil {
		log.Fatal(err)
	}
}

func createReader(input interface{}) io.Reader {
	json, _ := json.Marshal(input)

	return strings.NewReader(string(json))

}
func TestCreatePerson(t *testing.T) {
	testCases := []struct {
		name           string
		in             *http.Request
		out            *httptest.ResponseRecorder
		expectedStatus int
	}{
		{
			name: "Create User Success",
			in: httptest.NewRequest("GET", "/rest/v1/users", createReader(models.Person{
				Name:     "Test User",
				Email:    "user@test.com",
				Password: "TestP@ssw0rd123",
				Type:     "Parent",
			})),
			out:            httptest.NewRecorder(),
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Create User Validation fail",
			in: httptest.NewRequest("GET", "/rest/v1/users", createReader(models.Person{
				Name:     "Test User",
				Email:    "user@test.com",
				Password: "TestP@ssw0rd123",
			})),
			out:            httptest.NewRecorder(),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Create User Duplicate",
			in: httptest.NewRequest("GET", "/rest/v1/users", createReader(models.Person{
				Name:     "Test User",
				Email:    "user@test.com",
				Password: "TestP@ssw0rd123",
				Type:     "Parent",
			})),
			out:            httptest.NewRecorder(),
			expectedStatus: http.StatusConflict,
		},
	}

	for _, test := range testCases {
		p := PersonController{}
		t.Run(test.name, func(t *testing.T) {
			p.CreatePerson(test.out, test.in)
			if test.out.Code != test.expectedStatus {
				t.Error("Invalid response code:", test.out.Code)
			}

		})
	}

}
