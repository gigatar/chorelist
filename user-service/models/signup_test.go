package models

import (
	"math/rand"
	"testing"
	"time"
)

func randSeq(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyz")
	rand.Seed(time.Now().UnixNano())

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
func TestSignupValidate(t *testing.T) {
	person := Person{
		Email:    string(randSeq(5) + "@test.com"),
		Password: string(randSeq(15)),
		Name:     string(randSeq(10)),
		Type:     "Parent",
	}

	family := Family{
		Name: string(randSeq(5)),
	}

	testCases := []struct {
		name           string
		in             Signup
		expectedStatus bool
	}{
		{
			name: "Good test",
			in: Signup{
				Person: person,
				Family: family,
			},
			expectedStatus: true,
		},
		{
			name: "Missing Family",
			in: Signup{
				Person: person,
			},
			expectedStatus: false,
		},
		{
			name: "Missing Person",
			in: Signup{
				Family: family,
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
