package users

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

// Endpoints that are exposed.
type Endpoints struct {
	GetUsers    endpoint.Endpoint
	GetUserByID endpoint.Endpoint
	Login       endpoint.Endpoint
}

// MakeServerEndpoints returns the struct with the endpoint mapping.
func MakeServerEndpoints(srv Service) Endpoints {
	return Endpoints{
		GetUsers:    MakeGetUsersEndpoint(srv),
		GetUserByID: MakeGetUserByIDEndpoint(srv),
		Login:       MakeLoginEndpoint(srv),
	}
}

// MakeGetUsersEndpoint returns the response from our service "GetUsers".
func MakeGetUsersEndpoint(srv Service) endpoint.Endpoint {
	return func(ctx context.Context, inputRequest interface{}) (interface{}, error) {
		_ = inputRequest.(getUsersRequest)
		response, err := srv.GetUsers(ctx)
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
