package controllers

import (
	"chorelist/user-service/models"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateFamily(t *testing.T) {
	testCases := []struct {
		name           string
		in             *http.Request
		out            *httptest.ResponseRecorder
		expectedStatus int
	}{
		{
			name: "Create Family Success",
			in: httptest.NewRequest("GET", "/rest/v1/family", createReader(models.Family{
				Name: "Test Family",
			})),
			out:            httptest.NewRecorder(),
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "Create Family Validation fail",
			in:             httptest.NewRequest("GET", "/rest/v1/family", createReader(models.Family{})),
			out:            httptest.NewRecorder(),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Create Family Duplicate",
			in: httptest.NewRequest("GET", "/rest/v1/family", createReader(models.Family{
				Name: "Test Family",
			})),
			out:            httptest.NewRecorder(),
			expectedStatus: http.StatusCreated,
		},
	}

	for _, test := range testCases {
		f := FamilyController{}
		t.Run(test.name, func(t *testing.T) {
			f.CreateFamily(test.out, test.in)
			if test.out.Code != test.expectedStatus {
				t.Error("Invalid response code:", test.out.Code)
			}

		})
	}
}
