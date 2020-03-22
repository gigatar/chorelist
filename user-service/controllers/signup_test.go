package controllers

// func TestCreateSignup(t *testing.T) {
// 	// Create signup request
// 	signup := models.Signup{
// 		Person: models.Person{
// 			Email:    string(randSeq(5) + "@test"),
// 			Password: string(randSeq(15)),
// 			Name:     string(randSeq(10)),
// 			Type:     "Parent",
// 		},
// 		Family: models.Family{
// 			Name: string(randSeq(5)),
// 		},
// 	}
// 	testCases := []struct {
// 		name           string
// 		in             *http.Request
// 		out            *httptest.ResponseRecorder
// 		expectedStatus int
// 	}{
// 		{
// 			name:           "Create Signup Success",
// 			in:             httptest.NewRequest("GET", "/rest/v1/signup", createReader(signup)),
// 			out:            httptest.NewRecorder(),
// 			expectedStatus: http.StatusAccepted,
// 		},
// 		{
// 			name:           "Create Signup Duplicate",
// 			in:             httptest.NewRequest("GET", "/rest/v1/signup", createReader(signup)),
// 			out:            httptest.NewRecorder(),
// 			expectedStatus: http.StatusConflict,
// 		},
// 	}

// 	for _, test := range testCases {
// 		var s SignupController
// 		t.Run(test.name, func(t *testing.T) {
// 			s.CreateSignup(test.out, test.in)
// 			if test.out.Code != test.expectedStatus {
// 				t.Error("Invalid response code:", test.out.Code)
// 			}
// 		})
// 	}
// }
