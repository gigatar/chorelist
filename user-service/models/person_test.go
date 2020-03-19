package models

import "testing"

func TestPersonValidate(t *testing.T) {
	testCases := []struct {
		name           string
		in             Person
		expectedStatus bool
	}{
		{
			name: "Good test",
			in: Person{
				Name:     "Test User",
				Email:    "user@test.com",
				Password: "TestP@ssw0rd123",
				Type:     "Parent",
			},
			expectedStatus: true,
		},
		{
			name: "Good test child type",
			in: Person{
				Name:     "Test User",
				Email:    "user@test.com",
				Password: "TestP@ssw0rd123",
				Type:     "Child",
			},
			expectedStatus: true,
		},
		{
			name: "Bad name",
			in: Person{
				Name:     "",
				Email:    "user@test.com",
				Password: "TestP@ssw0rd123",
				Type:     "Parent",
			},
			expectedStatus: false,
		},
		{
			name: "Bad Email",
			in: Person{
				Name:     "Test User",
				Email:    "bad",
				Password: "TestP@ssw0rd123",
				Type:     "Parent",
			},
			expectedStatus: false,
		},
		{
			name: "Bad Password - too short",
			in: Person{
				Name:     "Test User",
				Email:    "user@test.com",
				Password: "test",
				Type:     "Parent",
			},
			expectedStatus: false,
		},
		{
			name: "Bad Password - too long",
			in: Person{
				Name:     "Test User",
				Email:    "user@test.com",
				Password: "TestP@ssw0rd123TestP@ssw0rd123TestP@ssw0rd123TestP@ssw0rd123TestP@ssw0rd123TestP@ssw0rd123TestP@ssw0rd123TestP@ssw0rd123TestP@sSs",
				Type:     "Parent",
			},
			expectedStatus: false,
		},
		{
			name: "Bad Type",
			in: Person{
				Name:     "Test User",
				Email:    "user@test.com",
				Password: "TestP@ssw0rd123",
				Type:     "bad",
			},
			expectedStatus: false,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			if !test.in.Validate() == test.expectedStatus {
				t.Error("Failed to properly validate Person")
				t.Fail()
			}
		})
	}
}

func TestStripSensitive(t *testing.T) {
	testCases := []struct {
		name           string
		in             Person
		expectedStatus bool
	}{
		{
			name: "Test with all Fields",
			in: Person{
				Name:        "Test User",
				Email:       "user@test.com",
				Password:    "TestP@ssw0rd123",
				OldPassword: "Test",
				Type:        "Parent",
			},
			expectedStatus: true,
		},
		{
			name: "Test with just password",
			in: Person{
				Name:     "Test User",
				Email:    "user@test.com",
				Password: "TestP@ssw0rd123",
				Type:     "Parent",
			},
			expectedStatus: true,
		},
		{
			name: "Test no sensitive",
			in: Person{
				Name:  "",
				Email: "user@test.com",
				Type:  "Parent",
			},
			expectedStatus: true,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			test.in.StripSensitive()
			if len(test.in.Password) > 0 || len(test.in.OldPassword) > 0 {
				t.Error("Still have sensitive data")
				t.Fail()
			}
		})
	}
}
