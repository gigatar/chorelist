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

func TestNoteValidate(t *testing.T) {
	testCases := []struct {
		name           string
		in             Note
		expectedStatus bool
	}{
		{
			name: "Good test",
			in: Note{
				Text: "I'm a note",
			},
			expectedStatus: true,
		},
		{
			name: "Bad text too short",
			in: Note{
				Text: "",
			},
			expectedStatus: false,
		},
		{
			name: "Bad name too long",
			in: Note{
				Text: randSeq(4097),
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
