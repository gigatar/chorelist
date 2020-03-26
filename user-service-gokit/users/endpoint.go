package users

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/go-kit/kit/endpoint"
)

// Endpoints that are exposed.
type Endpoints struct {
	GetUsers       endpoint.Endpoint
	GetUserByID    endpoint.Endpoint
	Login          endpoint.Endpoint
	ChangeName     endpoint.Endpoint
	ChangePassword endpoint.Endpoint
	CreateUser     endpoint.Endpoint
	DeleteUser     endpoint.Endpoint
}

// MakeServerEndpoints returns the struct with the endpoint mapping.
func MakeServerEndpoints(srv Service) Endpoints {
	return Endpoints{
		GetUsers:       MakeGetUsersEndpoint(srv),
		GetUserByID:    MakeGetUserByIDEndpoint(srv),
		Login:          MakeLoginEndpoint(srv),
		ChangeName:     MakeChangeNameEndpoint(srv),
		ChangePassword: MakeChangePasswordEndpoint(srv),
		CreateUser:     MakeCreateUserEndpoint(srv),
		DeleteUser:     MakeDeleteUserEndpoint(srv),
	}
}

// MakeCreateUserEndpoint returns the response from our service "CreateUser".
func MakeCreateUserEndpoint(srv Service) endpoint.Endpoint {
	return func(ctx context.Context, inputRequest interface{}) (interface{}, error) {
		request := inputRequest.(User)
		response, err := srv.CreateUser(ctx, request)
		if err != nil {
			return nil, err
		}

		var out createUserResponse
		out.Location = response.(primitive.ObjectID)
		return out, nil
	}
}

// MakeGetUsersEndpoint returns the response from our service "GetUsers".
func MakeGetUsersEndpoint(srv Service) endpoint.Endpoint {
	return func(ctx context.Context, inputRequest interface{}) (interface{}, error) {
		request := inputRequest.(getUsersRequest)
		response, err := srv.GetUsers(ctx, request.FamilyID)
		if err != nil {
			return nil, err
		}

		var out getUsersResponse
		out.Users = response
		return out, nil
	}
}

// MakeGetUserByIDEndpoint returns the response from our service "GetUserByID"
func MakeGetUserByIDEndpoint(srv Service) endpoint.Endpoint {
	return func(ctx context.Context, inputRequest interface{}) (interface{}, error) {
		user := inputRequest.(User)
		response, err := srv.GetUserByID(ctx, user)
		if err != nil {
			return nil, err
		}

		var out getUserByIDResponse
		out.User = response
		return out, nil
	}
}

// MakeLoginEndpoint returns the response from our service "Login".
func MakeLoginEndpoint(srv Service) endpoint.Endpoint {
	return func(ctx context.Context, inputRequest interface{}) (interface{}, error) {
		user := inputRequest.(User)
		response, err := srv.Login(ctx, user)
		if err != nil {
			return nil, err
		}

		var out loginResponse
		out.Login = response
		return out, nil
	}
}

// MakeChangeNameEndpoint returns the response from our service "ChangeName".
func MakeChangeNameEndpoint(srv Service) endpoint.Endpoint {
	return func(ctx context.Context, inputRequest interface{}) (interface{}, error) {
		user := inputRequest.(User)
		err := srv.ChangeName(ctx, user)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}
}

// MakeChangePasswordEndpoint returns the response from our service "ChangePassword".
func MakeChangePasswordEndpoint(srv Service) endpoint.Endpoint {
	return func(ctx context.Context, inputRequest interface{}) (interface{}, error) {
		user := inputRequest.(User)
		err := srv.ChangePassword(ctx, user)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}
}

// MakeDeleteUserEndpoint returns the response from our service "DeleteUser".
func MakeDeleteUserEndpoint(srv Service) endpoint.Endpoint {
	return func(ctx context.Context, inputRequest interface{}) (interface{}, error) {
		user := inputRequest.(User)
		err := srv.DeleteUser(ctx, user)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}
}
