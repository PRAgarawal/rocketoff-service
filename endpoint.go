package rocketoff

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	ShowEmThePointGod endpoint.Endpoint
	ShowEmTheBeard    endpoint.Endpoint
	OAuthComplete     endpoint.Endpoint
}

// MakeServerEndpoints initializes the endpoints for the service
func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		ShowEmTheBeard:    makeShowEmTheBeardEndpoint(s),
		ShowEmThePointGod: makeShowEmThePointGodEndpoint(s),
		OAuthComplete:     makeOAuthCompleteEndpoint(s),
	}
}

func makeShowEmTheBeardEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		command, ok := request.(*commandRequest)
		if !ok {
			return nil, ErrInvalidType{"commandRequest"}
		}
		return nil, svc.ShowEmTheBeard(ctx, &ImageCommand{
			WebhookURL:       command.webhookURL,
			RequestingUserID: command.requestingUserID,
		})
	}
}

func makeShowEmThePointGodEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		command, ok := request.(*commandRequest)
		if !ok {
			return nil, ErrInvalidType{"commandRequest"}
		}
		return nil, svc.ShowEmThePointGod(ctx, &ImageCommand{
			WebhookURL:       command.webhookURL,
			RequestingUserID: command.requestingUserID,
		})
	}
}

func makeOAuthCompleteEndpoint(svc Service) endpoint.Endpoint {
	return func(_ context.Context, _ interface{}) (interface{}, error) {
		return nil, nil
	}
}

type commandRequest struct {
	webhookURL       string
	requestingUserID string
}
