package rocketoff

import (
	"context"
	"encoding/json"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

// MakeHTTPHandler initializes all the available http routes
func MakeHTTPHandler(e Endpoints) http.Handler {
	r := mux.NewRouter()
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(encodeError),
	}
	r.Methods(http.MethodPost).Path("/show_em_the_beard/").
		Handler(kithttp.NewServer(
			e.ShowEmTheBeard,
			decodeShowEmTheBeardRequest,
			encodeShowEmTheBeardResponse,
			opts...,
		))
	r.Methods(http.MethodPost).Path("/show_em_the_point_god/").
		Handler(kithttp.NewServer(
			e.ShowEmThePointGod,
			decodeShowEmThePointGodRequest,
			encodeShowEmThePointGodResponse,
			opts...,
		))

	return r
}

func decodeShowEmTheBeardRequest(_ context.Context, _ *http.Request) (interface{}, error) {
	return nil, nil
}

func encodeShowEmTheBeardResponse(_ context.Context, _ http.ResponseWriter, _ interface{}) error {
	return nil
}

func decodeShowEmThePointGodRequest(_ context.Context, _ *http.Request) (interface{}, error) {
	return nil, nil
}

func encodeShowEmThePointGodResponse(_ context.Context, _ http.ResponseWriter, _ interface{}) error {
	return nil
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	// validate the type of error if it came from struct that implements the errors interface
	switch err.(type) {
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
