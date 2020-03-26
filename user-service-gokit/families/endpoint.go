package families

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/go-kit/kit/endpoint"
)

// Endpoints that are exposed.
type Endpoints struct {
	GetFamilyByID endpoint.Endpoint
	ChangeName    endpoint.Endpoint
	CreateFamily  endpoint.Endpoint
}

// MakeServerEndpoints returns the struct with the endpoint mapping.
func MakeServerEndpoints(srv Service) Endpoints {
	return Endpoints{
		GetFamilyByID: MakeGetFamilyByIDEndpoint(srv),
		ChangeName:    MakeChangeNameEndpoint(srv),
		CreateFamily:  MakeCreateFamilyEndpoint(srv),
	}
}

// MakeGetFamilyByIDEndpoint returns the response from our service "GetFamilyByID".
func MakeGetFamilyByIDEndpoint(srv Service) endpoint.Endpoint {
	return func(ctx context.Context, inputRequest interface{}) (interface{}, error) {
		family := inputRequest.(Family)
		response, err := srv.GetFamilyByID(ctx, family)
		if err != nil {
			return nil, err
		}

		var out getFamilyByIDResponse
		out.Family = response
		return out, nil
	}
}

// MakeChangeNameEndpoint returns the response from our service "ChangeName".
func MakeChangeNameEndpoint(srv Service) endpoint.Endpoint {
	return func(ctx context.Context, inputRequest interface{}) (interface{}, error) {
		family := inputRequest.(Family)
		err := srv.ChangeName(ctx, family)
		return nil, err
	}
}

// MakeCreateFamilyEndpoint returns the response from our service "CreateFamily".
func MakeCreateFamilyEndpoint(srv Service) endpoint.Endpoint {
	return func(ctx context.Context, inputRequest interface{}) (interface{}, error) {
		family := inputRequest.(Family)
		location, err := srv.CreateFamily(ctx, family)
		if err != nil {
			return nil, err
		}
		var out createFamilyResponse
		out.Location = location.(primitive.ObjectID)
		return out, err
	}
}
