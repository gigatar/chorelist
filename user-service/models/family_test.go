package models

import "testing"

func TestFamilyValidate(t *testing.T) {
	testCases := []struct {
		name           string
		in             Family
		expectedStatus bool
	}{
		{
			name: "Good test",
			in: Family{
				Name: "Test Family",
			},
			expectedStatus: true,
		},
		{
			name: "Bad name too short",
			in: Family{
				Name: "",
			},
			expectedStatus: false,
		},
		{
			name: "Bad name too long",
			in: Family{
				Name: "abc123abc123abc123abc123abc123abc123abc123abc123abc123abc123abc123abc123abc123abc123abc123abc123abc123abc123abc123abc123abc123456",
			},
			expectedStatus: false,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			if !test.in.Validate() == test.expectedStatus {
				t.Error("Failed to properly validate Family")
				t.Fail()
			}
		})
	}
}
