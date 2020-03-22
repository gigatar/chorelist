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

func TestValidateSignupEmail(t *testing.T) {
	testCases := []struct {
		name           string
		in             SignupEmail
		expectedStatus bool
	}{
		{
			name: "Good Test",
			in: SignupEmail{
				Name:  randSeq(10),
				Email: randSeq(4) + "@" + randSeq(3) + ".com",
				Code:  randSeq(15),
			},
			expectedStatus: true,
		},
		{
			name: "Bad Test Name Too Short",
			in: SignupEmail{
				Name:  "",
				Email: randSeq(4) + "@" + randSeq(3) + ".com",
				Code:  randSeq(15),
			},
			expectedStatus: false,
		},
		{
			name: "Bad Test Name Too Long",
			in: SignupEmail{
				Name:  randSeq(129),
				Email: randSeq(4) + "@" + randSeq(3) + ".com",
				Code:  randSeq(15),
			},
			expectedStatus: false,
		},
		{
			name: "Bad Test Invalid Email",
			in: SignupEmail{
				Name:  randSeq(5),
				Email: randSeq(4),
				Code:  randSeq(15),
			},
			expectedStatus: false,
		},
		{
			name: "Bad Test Code Too Long",
			in: SignupEmail{
				Name:  randSeq(5),
				Email: randSeq(4) + "@" + randSeq(3) + ".com",
				Code:  randSeq(20),
			},
			expectedStatus: false,
		},
		{
			name: "Bad Test Code Too Short",
			in: SignupEmail{
				Name:  randSeq(5),
				Email: randSeq(4) + "@" + randSeq(3) + ".com",
				Code:  randSeq(5),
			},
			expectedStatus: false,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			if !test.in.Validate() == test.expectedStatus {
				t.Error("Failed to properly validate signup input.")
				t.Fail()
			}
		})
	}
}
