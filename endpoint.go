package rocketoff

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	ShowEmThePointGod endpoint.Endpoint
	ShowEmTheBeard    endpoint.Endpoint
}

// MakeServerEndpoints initializes the endpoints for the service
func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		ShowEmThePointGod: makeShowEmThePointGodEndpoint(s),
		ShowEmTheBeard:    makeShowEmTheBeardEndpoint(s),
	}
}

func makeShowEmThePointGodEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		err := svc.ShowEmThePointGod(ctx)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}
}

func makeShowEmTheBeardEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		err := svc.ShowEmTheBeard(ctx)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}
}
