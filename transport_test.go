package rocketoff

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeHTTPHandler(t *testing.T) {
	e := Endpoints{}

	result := MakeHTTPHandler(e)
	assert.Implements(t, (*http.Handler)(nil), result)
}

func TestEncodeShowEmTheBeardResponse(t *testing.T) {
	t.Run("valid request", func(t *testing.T) {
		err := encodeShowEmTheBeardResponse(context.Background(), nil, nil)
		assert.Nil(t, err)
	})
}

func TestDecodeShowEmTheBeardRequest(t *testing.T) {
	t.Run("valid request", func(t *testing.T) {
		decodedRequest, err := decodeShowEmTheBeardRequest(context.Background(), nil)
		assert.Nil(t, err)
		assert.Nil(t, decodedRequest)
	})
}
