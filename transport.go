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
// TODO: This layer is very slack-specific right now. Can a different transport implementation handle this functionality for a different chat apps (keybase?) given the same Service interface and Endpoint layer?
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
		_, err = slack.SlashCommandParse(r)
		if err != nil {
			return nil, err
		}

		if err = verifier.Ensure(); err != nil {
			return nil, err
		}

		return nil, nil
	}
}

func encodeResponse(_ context.Context, w http.ResponseWriter, resp interface{}) error {
	msg, err := buildSlackMsg(resp)
	if err != nil {
		return err
	}

	bytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err = w.Write(bytes); err != nil {
		return err
	}

	return nil
}

// buildSlackMsg takes a commandResponse interface from the endpoint layer, and builds a slack.Msg struct that can be written to an http.ResponseWriter. It returns an error if the resp is not of the proper type.
func buildSlackMsg(resp interface{}) (*slack.Msg, error) {
	cmdResponse, ok := resp.(*commandResponse)
	if !ok {
		return nil, ErrInvalidType{"commandResponse"}
	}
	return &slack.Msg{
		ResponseType: slack.ResponseTypeInChannel,
		Attachments: []slack.Attachment{
			{
				ImageURL: cmdResponse.imageUrl,
			},
		},
	}, nil
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
