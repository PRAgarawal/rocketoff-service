package rocketoff

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/PRAgarawal/rocketoff/chat"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

const (
	slackOAuthURL = "https://slack.com/oauth/v2/authorize?client_id=292540839159.944358440549&scope=chat:write,commands"
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

	router.Methods(http.MethodGet).Path("/slack_oauth_complete/").
		Handler(kithttp.NewServer(
			e.OAuthComplete,
			decodeSlackOAuthCompleteRequest,
			encodeSlackOAuthCompleteResponse,
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

func decodeSlackOAuthCompleteRequest(_ context.Context, _ *http.Request) (interface{}, error) {
	return nil, nil
}

func encodeSlackOAuthCompleteResponse(_ context.Context, writer http.ResponseWriter, _ interface{}) error {
	writer.Header().Set("Location", slackOAuthURL)
	writer.WriteHeader(http.StatusFound)
	return nil
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
