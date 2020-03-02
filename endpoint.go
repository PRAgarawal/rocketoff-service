package rocketoff

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	ShowEmThePointGod endpoint.Endpoint
	ShowEmTheBeard    endpoint.Endpoint
	OAuthComplete     endpoint.Endpoint
	OAuthRedirect     endpoint.Endpoint
}

// MakeServerEndpoints initializes the endpoints for the service
func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		ShowEmTheBeard:    makeShowEmTheBeardEndpoint(s),
		ShowEmThePointGod: makeShowEmThePointGodEndpoint(s),
		OAuthComplete:     makeOAuthCompleteEndpoint(s),
		OAuthRedirect:     makeOAuthRedirectEndpoint(s),
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
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(*oauthCompleteRequest)
		if !ok {
			return nil, ErrInvalidType{"oauthCompleteRequest"}
		}

		oauthOptions := &OAuthCompleteOptions{
			Code:  req.code,
			State: req.state,
		}
		if oauthOptions.Code == "" {
			return nil, ErrInvalidValue{"code must be provided"}
		}

		return svc.CompleteChatOAuth(ctx, oauthOptions)
	}
}

func makeOAuthRedirectEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, _ interface{}) (interface{}, error) {
		return svc.RedirectForOAuth(ctx)
	}
}

type commandRequest struct {
	webhookURL       string
	requestingUserID string
}

type oauthCompleteRequest struct {
	code  string
	state string
}
