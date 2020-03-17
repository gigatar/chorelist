package controllers

import (
	"chorelist/user-service/database"
	"chorelist/user-service/models"
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
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

func randSeq(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyz")
	rand.Seed(time.Now().UnixNano())

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
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

func TestLogin(t *testing.T) {
	// Create user for setup purposes
	person := models.Person{
		Email:    string(randSeq(5) + "@test.com"),
		Password: string(randSeq(15)),
		Name:     string(randSeq(10)),
		Type:     "Parent",
	}

	testCases := []struct {
		name           string
		in             *http.Request
		out            *httptest.ResponseRecorder
		expectedStatus int
	}{
		{
			name:           "Create User (Setup)",
			in:             httptest.NewRequest("GET", "/rest/v1/users", createReader(person)),
			out:            httptest.NewRecorder(),
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "Login Success",
			in:             httptest.NewRequest("GET", "/rest/v1/users/login", createReader(person)),
			out:            httptest.NewRecorder(),
			expectedStatus: http.StatusOK,
		},
		{
			name: "Login BadPassword",
			in: httptest.NewRequest("GET", "/rest/v1/users/login", createReader(models.Person{
				Email:    "user@test.com",
				Password: "bad",
			})),
			out:            httptest.NewRecorder(),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Login Bad User",
			in: httptest.NewRequest("GET", "/rest/v1/users/login", createReader(models.Person{
				Email:    "bad",
				Password: "TestP@ssw0rd123",
			})),
			out:            httptest.NewRecorder(),
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, test := range testCases {
		p := PersonController{}
		t.Run(test.name, func(t *testing.T) {
			if strings.Compare(test.name, "Create User (Setup)") == 0 {
				p.CreatePerson(test.out, test.in)
			} else {
				p.Login(test.out, test.in)
			}
			if test.out.Code != test.expectedStatus {
				t.Error("Invalid response code:", test.out.Code)
			}

		})
	}
}

func TestChangeName(t *testing.T) {
	// Create user for setup purposes
	person := models.Person{
		Email:    string(randSeq(5) + "@test.com"),
		Password: string(randSeq(15)),
		Name:     string(randSeq(10)),
		Type:     "Parent",
	}

	// Authorization header
	var auth string

	testCases := []struct {
		name           string
		in             *http.Request
		out            *httptest.ResponseRecorder
		expectedStatus int
	}{
		{
			name:           "Create User (Setup)",
			in:             httptest.NewRequest("GET", "/rest/v1/users", createReader(person)),
			out:            httptest.NewRecorder(),
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "Login User (Setup)",
			in:             httptest.NewRequest("GET", "/rest/v1/users/login", createReader(person)),
			out:            httptest.NewRecorder(),
			expectedStatus: http.StatusOK,
		},
		{
			name: "ChangeName success",
			in: httptest.NewRequest("PATCH", "/rest/v1/users/", createReader(models.Person{
				Name: "Updated",
			})),
			out:            httptest.NewRecorder(),
			expectedStatus: http.StatusNoContent,
		},
	}

	for _, test := range testCases {
		p := PersonController{}
		t.Run(test.name, func(t *testing.T) {
			if strings.Compare(test.name, "Create User (Setup)") == 0 {
				p.CreatePerson(test.out, test.in)
			} else if strings.Compare(test.name, "Login User (Setup)") == 0 {
				p.Login(test.out, test.in)
			} else {
				test.in.Header.Add("authorization", string("Bearer "+auth))
				p.ChangeName(test.out, test.in)
			}
			if test.out.Code != test.expectedStatus {
				t.Error("Invalid response code:", test.out.Code)
			}

			if strings.Compare(test.name, "Login User (Setup)") == 0 {
				auth = test.out.Header().Get("Authorization")
			}

		})
	}
}
