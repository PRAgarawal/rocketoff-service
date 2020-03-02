package rocketoff

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/PRAgarawal/rocketoff/chat"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

// MakeHTTPHandler initializes all the available http routes
func MakeHTTPHandler(e Endpoints, commandDecoder chat.CommandDecoder) http.Handler {
	router := mux.NewRouter()
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(encodeError),
	}
	decodeSlashCommandRequest := makeSlashCommandRequestDecoder(commandDecoder)
	router.Methods(http.MethodPost).Path("/show_em_the_beard/").
		Handler(kithttp.NewServer(
			e.ShowEmTheBeard,
			decodeSlashCommandRequest,
			encodeResponse,
			opts...,
		))
	router.Methods(http.MethodPost).Path("/show_em_the_point_god/").
		Handler(kithttp.NewServer(
			e.ShowEmThePointGod,
			decodeSlashCommandRequest,
			encodeResponse,
			opts...,
		))

	router.Methods(http.MethodGet).Path("/oauth_complete/").
		Queries(
			"state", "{state:[0-9a-zA-Z\\W]+|}",
			"code", "{code:[0-9a-zA-Z\\W]+|}").
		Handler(kithttp.NewServer(
			e.OAuthComplete,
			decodeOAuthCompleteRequest,
			encodeRedirectResponse,
			opts...,
		))
	router.Methods(http.MethodGet).Path("/oauth_complete/").
		Queries(
			"code", "{code:[0-9a-zA-Z\\W]+|}").
		Handler(kithttp.NewServer(
			e.OAuthComplete,
			decodeOAuthCompleteRequest,
			encodeRedirectResponse,
			opts...,
		))

	router.Methods(http.MethodGet).Path("/redirect/").
		Handler(kithttp.NewServer(
			e.OAuthRedirect,
			decodeRedirectRequest,
			encodeRedirectResponse,
			opts...,
		))

	return router
}

func makeSlashCommandRequestDecoder(commandDecoder chat.CommandDecoder) kithttp.DecodeRequestFunc {
	return func(ctx context.Context, request *http.Request) (interface{}, error) {
		command, err := commandDecoder.DecodeCommand(ctx, request)
		if err != nil {
			return nil, err
		}

		return &commandRequest{
			webhookURL:       command.WebhookURL,
			requestingUserID: command.RequestingUserID,
		}, nil
	}
}

func encodeResponse(_ context.Context, _ http.ResponseWriter, _ interface{}) error {
	return nil
}

func decodeOAuthCompleteRequest(_ context.Context, request *http.Request) (interface{}, error) {
	vars := mux.Vars(request)
	if vars == nil {
		return nil, fmt.Errorf("invalid or missing query params on oauth redirect")
	}
	return &oauthCompleteRequest{vars["code"], vars["state"]}, nil
}

func encodeRedirectResponse(_ context.Context, writer http.ResponseWriter, resp interface{}) error {
	if redirectURI, ok := resp.(string); ok {
		writer.Header().Set("Location", redirectURI)
		writer.WriteHeader(http.StatusFound)
		return nil
	}

	return fmt.Errorf("no redirect URI found")
}

func decodeRedirectRequest(_ context.Context, _ *http.Request) (interface{}, error) {
	return nil, nil
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err.(type) {
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
