package controllers

import (
	"chorelist/user-service/models"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

func TestAddFamilyMember(t *testing.T) {
	child := models.Person{
		Email:    string(randSeq(5) + "@test.com"),
		Password: string(randSeq(15)),
		Name:     string(randSeq(10)),
		Type:     "Child",
	}
	parent := models.Person{
		Email:    string(randSeq(5) + "@test.com"),
		Password: string(randSeq(15)),
		Name:     string(randSeq(10)),
		Type:     "Parent",
	}
	badType := models.Person{
		Email:    string(randSeq(5) + "@test.com"),
		Password: string(randSeq(15)),
		Name:     string(randSeq(10)),
		Type:     "Invalid",
	}

	// Setup
	auth, err := createUserAndLogin()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	var childAuth string

	testCases := []struct {
		name           string
		in             *http.Request
		out            *httptest.ResponseRecorder
		expectedStatus int
	}{
		{
			name:           "Add Child to Family",
			in:             httptest.NewRequest("GET", "/rest/v1/families/add", createReader(child)),
			out:            httptest.NewRecorder(),
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "Duplicate Add User to Family",
			in:             httptest.NewRequest("GET", "/rest/v1/families/add", createReader(child)),
			out:            httptest.NewRecorder(),
			expectedStatus: http.StatusConflict,
		},
		{
			name:           "Add Parent to Family",
			in:             httptest.NewRequest("GET", "/rest/v1/families/add", createReader(parent)),
			out:            httptest.NewRecorder(),
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "Add Bad Type to Family",
			in:             httptest.NewRequest("GET", "/rest/v1/families/add", createReader(badType)),
			out:            httptest.NewRecorder(),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Login as child",
			in:             httptest.NewRequest("GET", "/rest/v1/families/add", createReader(child)),
			out:            httptest.NewRecorder(),
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Add to Family as child",
			in:             httptest.NewRequest("GET", "/rest/v1/families/add", createReader(parent)),
			out:            httptest.NewRecorder(),
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, test := range testCases {
		var f FamilyController
		t.Run(test.name, func(t *testing.T) {
			if strings.Compare(test.name, "Login as child") == 0 {
				var p PersonController
				p.Login(test.out, test.in)
				if test.out.Code != test.expectedStatus {
					t.Error("Invalid response code:", test.out.Code)
					t.Fail()
				}

				childAuth = test.out.Header().Get("Authorization")

			} else {
				if strings.Compare(test.name, "Add to Family as child") == 0 {
					test.in.Header.Add("Authorization", "Bearer "+childAuth)
				} else {
					test.in.Header.Add("Authorization", "Bearer "+auth)
				}
				f.AddFamilyMember(test.out, test.in)
				if test.out.Code != test.expectedStatus {
					t.Error("Invalid response code:", test.out.Code)
				}
			}
		})
	}
}

func TestRemoveFamilyMember(t *testing.T) {
	child := models.Person{
		Email:    string(randSeq(5) + "@test.com"),
		Password: string(randSeq(15)),
		Name:     string(randSeq(10)),
		Type:     "Child",
	}

	// Setup
	auth, err := createUserAndLogin()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	var childAuth string
	var childURI string

	testCases := []struct {
		name           string
		in             *http.Request
		out            *httptest.ResponseRecorder
		expectedStatus int
	}{
		{
			name:           "Add Child to Family",
			in:             httptest.NewRequest("GET", "/rest/v1/families/add", createReader(child)),
			out:            httptest.NewRecorder(),
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "Login as child",
			in:             httptest.NewRequest("POST", "/rest/v1/users/login", createReader(child)),
			out:            httptest.NewRecorder(),
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Remove from Family as child",
			in:             httptest.NewRequest("DELETE", "/rest/v1/families/persons/", nil),
			out:            httptest.NewRecorder(),
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "Remove Child Success",
			in:             httptest.NewRequest("DELETE", "/rest/v1/families/persons/", nil),
			out:            httptest.NewRecorder(),
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "Remove No-exist",
			in:             httptest.NewRequest("DELETE", "/rest/v1/families/persons/", nil),
			out:            httptest.NewRecorder(),
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, test := range testCases {
		var f FamilyController
		t.Run(test.name, func(t *testing.T) {
			if strings.Compare(test.name, "Login as child") == 0 {
				var p PersonController
				p.Login(test.out, test.in)
				if test.out.Code != test.expectedStatus {
					t.Error("Invalid response code:", test.out.Code)
					t.Fail()
				}

				childAuth = test.out.Header().Get("Authorization")

			} else {
				if strings.Compare(test.name, "Remove from Family as child") == 0 {
					test.in.Header.Add("Authorization", "Bearer "+childAuth)
				} else {
					test.in.Header.Add("Authorization", "Bearer "+auth)
				}
				if strings.Compare(test.name, "Add Child to Family") == 0 {
					f.AddFamilyMember(test.out, test.in)
					childURI = strings.TrimPrefix(test.out.Header().Get("Location"), "/rest/v1/users/")
				} else {
					test.in.RequestURI += childURI
					test.in.URL.Path += childURI

					r := mux.NewRouter()
					r.HandleFunc("/rest/v1/families/persons/{personID}", f.RemoveFamilyMember).Methods("DELETE")
					r.ServeHTTP(test.out, test.in)
				}
				if test.out.Code != test.expectedStatus {
					t.Error("Invalid response code:", test.out.Code)
				}
			}
		})
	}
}
