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

func TestFamilyValidateAddPerson(t *testing.T) {
	goodFamily := Family{
		Name: "Good Test Family",
		Person: []string{"person1", "person2", "person3", "person4",
			"person5", "person6", "person7", "person8", "person9", "person10",
			"person11", "person12", "person13", "person14"},
	}
	badFamily := Family{
		Name: "Bad Test Family",
		Person: []string{"person1", "person2", "person3", "person4",
			"person5", "person6", "person7", "person8", "person9", "person10",
			"person11", "person12", "person13", "person14", "person15",
		},
	}

	testCases := []struct {
		name           string
		in             Family
		expectedStatus bool
	}{
		{
			name:           "Good test",
			in:             goodFamily,
			expectedStatus: true,
		},
		{
			name:           "Bad test",
			in:             badFamily,
			expectedStatus: false,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			if !test.in.ValidateAddPerson() == test.expectedStatus {
				t.Error("Failed to properly validate Family")
				t.Fail()
			}
		})
	}
}
