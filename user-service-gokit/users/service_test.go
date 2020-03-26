package users

import (
	"chorelist/user-service-gokit/gigatarerrors"
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetAllUsers(t *testing.T) {
	testCases := []struct {
		name          string
		in            string
		expectedError error
	}{
		{
			name:          "Test get all users invalid familyID",
			in:            "",
			expectedError: gigatarerrors.ErrNotFound,
		},
	}
	for _, test := range testCases {
		var userController User
		users, err := userController.GetUsers(context.TODO(), test.in)
		if err != test.expectedError {
			t.Error("Invalid Error:", err, "Expect:", test.expectedError)
			t.Fail()
		}
		if len(users) > 1 {
			for _, user := range users {
				if len(user.Password) > 1 || len(user.OldPassword) > 1 {
					t.Error("Sensitive information leak")
					t.FailNow()
				}
			}

		}
	}
}

func TestGetUserByID(t *testing.T) {
	id, _ := primitive.ObjectIDFromHex("5e7bdf9b6e0a707af2aeae75")
	testCases := []struct {
		name          string
		in            User
		expectedError error
	}{
		{
			name: "Test get user",
			in: User{
				ID:       id,
				FamilyID: "5e7bdf9b6e0a707af2aeae76",
			},
			expectedError: nil,
		},
		{
			name:          "Test Bad Family",
			in:            User{ID: id, FamilyID: ""},
			expectedError: gigatarerrors.ErrNotFound,
		},
		{
			name:          "Test get non-exist",
			in:            User{ID: primitive.NilObjectID, FamilyID: ""},
			expectedError: gigatarerrors.ErrNotFound,
		},
	}
	for _, test := range testCases {
		var userController User
		user, err := userController.GetUserByID(context.TODO(), test.in)
		if err != test.expectedError {
			t.Error("Invalid Error:", err, "Expect:", test.expectedError)
			t.Fail()
		}

		if len(user.Password) > 1 || len(user.OldPassword) > 1 {
			t.Error("Sensitive information leak")
			t.FailNow()
		}
	}

}

func TestLogin(t *testing.T) {
	testCases := []struct {
		name          string
		in            User
		expectedError error
	}{
		{
			name: "Test good login",
			in: User{
				Email:    "user@test.com",
				Password: "Testtest",
			},
			expectedError: nil,
		},
		{
			name: "Test bad email",
			in: User{
				Email:    "bad@bad.com",
				Password: "Testtest",
			},
			expectedError: gigatarerrors.ErrBadRequest,
		},
		{
			name: "Test bad password",
			in: User{
				Email:    "user@test.com",
				Password: "badpassword",
			},
			expectedError: gigatarerrors.ErrBadRequest,
		},
	}

	for _, test := range testCases {
		var userController User
		user, err := userController.Login(context.TODO(), test.in)
		if err != test.expectedError {
			t.Error("Invalid Error:", err, "Expect:", test.expectedError)
			t.Fail()
		}

		if len(user.Password) > 1 || len(user.OldPassword) > 1 {
			t.Error("Sensitive information leak")
			t.FailNow()
		}
	}
}

// Login(ctx context.Context, inputUser User) (User, error)
// ChangeName(ctx context.Context, inputUser User) error
// ChangePassword(ctx context.Context, inputUser User) error
