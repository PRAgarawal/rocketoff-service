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
		ShowEmTheBeard:    makeShowEmTheBeardEndpoint(s),
		ShowEmThePointGod: makeShowEmThePointGodEndpoint(s),
	}
}

func makeShowEmTheBeardEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		imgReply, err := svc.ShowEmTheBeard(ctx)
		if err != nil {
			return nil, err
		}
		return &commandResponse{imgReply.ImageURL}, nil
	}
}

func makeShowEmThePointGodEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		imgReply, err := svc.ShowEmThePointGod(ctx)
		if err != nil {
			return nil, err
		}
		return &commandResponse{imgReply.ImageURL}, nil
	}
}

type commandResponse struct {
	imageUrl string
}
