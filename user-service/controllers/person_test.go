package controllers

import (
	"chorelist/user-service/database"
	"chorelist/user-service/models"
	"encoding/json"
	"errors"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"
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
	time.Sleep(time.Millisecond)
	letters := []rune("abcdefghijklmnopqrstuvwxyz")
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[r.Intn(len(letters))]
	}
	return string(b)
}

func createUserAndLogin() (string, error) {
	// Create user for setup purposes
	person := models.Person{
		Email:    string(randSeq(5) + "@test.com"),
		Password: string(randSeq(15)),
		Name:     string(randSeq(10)),
		Type:     "Parent",
	}
	family := models.Family{
		Name: string(randSeq(5)),
	}

	// Authorization header
	var auth string

	actions := []struct {
		name           string
		in             *http.Request
		out            *httptest.ResponseRecorder
		expectedStatus int
	}{

		{
			name:           "Login User (Setup)",
			in:             httptest.NewRequest("GET", "/rest/v1/users/login", createReader(person)),
			out:            httptest.NewRecorder(),
			expectedStatus: http.StatusOK,
		},
	}

	var p PersonController
	var f FamilyController
	// Create user
	uid, err := p.createPerson(person)
	if err != nil {
		return "", err
	}
	family.Person = []string{uid}

	// Create family
	fid, err := f.createFamily(family)
	if err != nil {
		return "", err
	}

	// Add family to user
	person.FamilyID = fid
	if err := p.dao.ChangeFamilyID(uid, fid); err != nil {
		return "", err
	}

	for _, action := range actions {

		p.Login(action.out, action.in)

		action.in.Header.Add("authorization", string("Bearer "+auth))
		if action.out.Code != action.expectedStatus {
			return "", errors.New("Unable to complete action: " + action.name + "return: " + strconv.Itoa(action.out.Code))
		}
		auth = action.out.Header().Get("Authorization")
	}

	return auth, nil
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
				var person models.Person
				if err := json.NewDecoder(test.in.Body).Decode(&person); err != nil {
					t.Error(err)
					t.Fail()
				}
				if _, err := p.createPerson(person); err != nil {
					t.Error(err)
					t.Fail()
				}
			} else {
				p.Login(test.out, test.in)

				if test.out.Code != test.expectedStatus {
					t.Error("Invalid response code:", test.out.Code)
				}
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
		{
			name: "ChangeName fail name validation",
			in: httptest.NewRequest("PATCH", "/rest/v1/users/", createReader(models.Person{
				Name: "",
			})),
			out:            httptest.NewRecorder(),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "ChangeName bad token",
			in: httptest.NewRequest("PATCH", "/rest/v1/users/", createReader(models.Person{
				Name: "Updated2",
			})),
			out:            httptest.NewRecorder(),
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, test := range testCases {
		p := PersonController{}
		t.Run(test.name, func(t *testing.T) {
			if strings.Compare(test.name, "Create User (Setup)") == 0 {
				var person models.Person
				if err := json.NewDecoder(test.in.Body).Decode(&person); err != nil {
					t.Error(err)
					t.Fail()
				}
				if _, err := p.createPerson(person); err != nil {
					t.Error(err)
					t.Fail()
				}
			} else {
				if strings.Compare(test.name, "Login User (Setup)") == 0 {
					p.Login(test.out, test.in)
				} else {
					if strings.Compare(test.name, "ChangeName bad token") == 0 {
						test.in.Header.Add("authorization", string("Bearer bad"))
					} else {
						test.in.Header.Add("authorization", string("Bearer "+auth))
					}
					p.ChangeName(test.out, test.in)
				}
				if test.out.Code != test.expectedStatus {
					t.Error("Invalid response code:", test.out.Code)
				}

				if strings.Compare(test.name, "Login User (Setup)") == 0 {
					auth = test.out.Header().Get("Authorization")
				}
			}
		})
	}
}

func TestPasswordChange(t *testing.T) {
	// Create user for setup purposes
	person := models.Person{
		Email:    string(randSeq(5) + "@test.com"),
		Password: string(randSeq(15)),
		Name:     string(randSeq(10)),
		Type:     "Parent",
	}

	// Authorization header
	var auth string

	newPassword := string(randSeq(15))

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
			name: "ChangePassword success",
			in: httptest.NewRequest("PATCH", "/rest/v1/users/", createReader(models.Person{
				Password:    newPassword,
				OldPassword: person.Password,
			})),
			out:            httptest.NewRecorder(),
			expectedStatus: http.StatusNoContent,
		},
		{
			name: "ChangePassword bad old password",
			in: httptest.NewRequest("PATCH", "/rest/v1/users/", createReader(models.Person{
				Password:    newPassword,
				OldPassword: "bad",
			})),
			out:            httptest.NewRecorder(),
			expectedStatus: http.StatusForbidden,
		},
		{
			name: "ChangePassword bad validation fail",
			in: httptest.NewRequest("PATCH", "/rest/v1/users/", createReader(models.Person{
				Password:    "bad",
				OldPassword: newPassword,
			})),
			out:            httptest.NewRecorder(),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "ChangePassword bad token",
			in: httptest.NewRequest("PATCH", "/rest/v1/users/", createReader(models.Person{
				Password:    person.Password,
				OldPassword: newPassword,
			})),
			out:            httptest.NewRecorder(),
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, test := range testCases {
		p := PersonController{}
		t.Run(test.name, func(t *testing.T) {
			if strings.Compare(test.name, "Create User (Setup)") == 0 {
				var person models.Person
				if err := json.NewDecoder(test.in.Body).Decode(&person); err != nil {
					t.Error(err)
					t.Fail()
				}
				if _, err := p.createPerson(person); err != nil {
					t.Error(err)
					t.Fail()
				}
			} else {
				if strings.Compare(test.name, "Login User (Setup)") == 0 {
					p.Login(test.out, test.in)
				} else {
					if strings.Compare(test.name, "ChangePassword bad token") == 0 {
						test.in.Header.Add("authorization", string("Bearer bad"))
					} else {
						test.in.Header.Add("authorization", string("Bearer "+auth))
					}
					p.ChangePassword(test.out, test.in)
				}
				if test.out.Code != test.expectedStatus {
					t.Error("Invalid response code:", test.out.Code)
				}

				if strings.Compare(test.name, "Login User (Setup)") == 0 {
					auth = test.out.Header().Get("Authorization")
				}
			}
		})
	}
}

func TestPersonDelete(t *testing.T) {
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
			name:           "Delete User success",
			in:             httptest.NewRequest("DELETE", "/rest/v1/users/delete", nil),
			out:            httptest.NewRecorder(),
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "Delete User no-exist",
			in:             httptest.NewRequest("DELETE", "/rest/v1/users/delete", nil),
			out:            httptest.NewRecorder(),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Delete User bad token",
			in:             httptest.NewRequest("DELETE", "/rest/v1/users/", nil),
			out:            httptest.NewRecorder(),
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, test := range testCases {
		p := PersonController{}
		t.Run(test.name, func(t *testing.T) {
			if strings.Compare(test.name, "Create User (Setup)") == 0 {
				var person models.Person
				if err := json.NewDecoder(test.in.Body).Decode(&person); err != nil {
					t.Error(err)
					t.Fail()
				}
				if _, err := p.createPerson(person); err != nil {
					t.Error(err)
					t.Fail()
				}
			} else {
				if strings.Compare(test.name, "Login User (Setup)") == 0 {
					p.Login(test.out, test.in)
				} else {
					if strings.Compare(test.name, "Delete User bad token") == 0 {
						test.in.Header.Add("authorization", string("Bearer bad"))
					} else {
						test.in.Header.Add("authorization", string("Bearer "+auth))
					}
					p.DeletePerson(test.out, test.in)
				}
				if test.out.Code != test.expectedStatus {
					t.Error("Invalid response code:", test.out.Code)
				}

				if strings.Compare(test.name, "Login User (Setup)") == 0 {
					auth = test.out.Header().Get("Authorization")
				}
			}
		})
	}
}

func TestGetPersonType(t *testing.T) {
	// Create user for setup purposes
	person := models.Person{
		Email:    string(randSeq(5) + "@test.com"),
		Password: string(randSeq(15)),
		Name:     string(randSeq(10)),
		Type:     "Parent",
	}

	var userID string

	testCases := []struct {
		name          string
		in            string
		expectedError bool
		expectType    string
	}{
		{
			name:          "Good User",
			in:            userID,
			expectedError: false,
			expectType:    person.Type,
		},
		{
			name:          "Bad User",
			in:            "abc123",
			expectedError: true,
			expectType:    "",
		},
	}

	setup := PersonController{}
	userID, err := setup.createPerson(person)
	if err != nil {
		t.Error("Unable to create setup user:", err)
		t.Fail()
	}

	for _, test := range testCases {
		p := PersonController{}
		if strings.Contains(test.name, "Good User") {
			test.in = userID
		}
		t.Run(test.name, func(t *testing.T) {
			userType, err := p.getPersonType(test.in)
			if err != nil {
				if test.expectedError == false {
					t.Error("Got error when we didn't expect one")
					t.Fail()
				}
			} else if test.expectedError == true {
				t.Error("Didn't get error when we expected one")
				t.Fail()
			}

			if strings.Compare(test.expectType, userType) != 0 {
				t.Error("Invalid user type:", userType)
				t.Fail()
			}
		})
	}
}

func TestViewPerson(t *testing.T) {
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
			in:             httptest.NewRequest("POST", "/rest/v1/families/add", createReader(child)),
			out:            httptest.NewRecorder(),
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "View Person as parent",
			in:             httptest.NewRequest("GET", "/rest/v1/users/", nil),
			out:            httptest.NewRecorder(),
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Login as child",
			in:             httptest.NewRequest("POST", "/rest/v1/users/login", createReader(child)),
			out:            httptest.NewRecorder(),
			expectedStatus: http.StatusOK,
		},
		{
			name:           "View Person as child",
			in:             httptest.NewRequest("GET", "/rest/v1/users/", nil),
			out:            httptest.NewRecorder(),
			expectedStatus: http.StatusOK,
		},
	}

	for _, test := range testCases {
		var p PersonController
		t.Run(test.name, func(t *testing.T) {
			if strings.Compare(test.name, "Login as child") == 0 {
				p.Login(test.out, test.in)
				if test.out.Code != test.expectedStatus {
					t.Error("Invalid response code:", test.out.Code)
					t.Fail()
				}

				childAuth = test.out.Header().Get("Authorization")

			} else {
				if strings.Compare(test.name, "View Person as child") == 0 {
					test.in.Header.Add("Authorization", "Bearer "+childAuth)
				} else {
					test.in.Header.Add("Authorization", "Bearer "+auth)
				}
				if strings.Compare(test.name, "Add Child to Family") == 0 {
					var f FamilyController
					f.AddFamilyMember(test.out, test.in)
					childURI = test.out.Header().Get("Location")
					if test.out.Code != test.expectedStatus {
						t.Error("Invalid response code:", test.out.Code)
					}
				} else {
					test.in.RequestURI += childURI
					test.in.URL.Path += childURI

					r := mux.NewRouter()
					r.HandleFunc("/rest/v1/users/{personID}", p.ViewPerson).Methods("GET")
					r.ServeHTTP(test.out, test.in)
					if test.out.Code != test.expectedStatus {
						t.Error("Invalid response code:", test.out.Code)
						t.FailNow()
					}
					var responsePerson models.Person
					err := json.NewDecoder(test.out.Body).Decode(&responsePerson)
					if err != nil {
						t.Error("Unable to unmarshal")
						t.Fail()
					}
					if len(responsePerson.Password) > 0 || len(responsePerson.OldPassword) > 0 {
						t.Error("Got password in response")
						t.Fail()
					}
				}
			}
		})
	}
}
