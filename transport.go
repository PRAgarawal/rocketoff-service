package rocketoff

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/nlopes/slack"
)

// MakeHTTPHandler initializes all the available http routes
func MakeHTTPHandler(e Endpoints, signingSecret string) http.Handler {
	router := mux.NewRouter()
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(encodeError),
	}
	decodeSlashCommandRequest := makeSlashCommandRequestDecoder(signingSecret)
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

	return router
}

func makeSlashCommandRequestDecoder(signingSecret string) kithttp.DecodeRequestFunc {
	return func(_ context.Context, r *http.Request) (interface{}, error) {
		verifier, err := slack.NewSecretsVerifier(r.Header, signingSecret)
		if err != nil {
			return nil, err
		}

		r.Body = ioutil.NopCloser(io.TeeReader(r.Body, &verifier))
		command, err := slack.SlashCommandParse(r)
		if err != nil {
			return nil, err
		}

		if err = verifier.Ensure(); err != nil {
			return nil, err
		}

		return &commandRequest{
			responseURL:        command.ResponseURL,
			requestingUsername: command.UserName,
		}, nil
	}
}

func encodeResponse(_ context.Context, _ http.ResponseWriter, _ interface{}) error {
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
