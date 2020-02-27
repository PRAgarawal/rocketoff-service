package rocketoff

import (
	"context"
	"net/http"
	"testing"

	"github.com/PRAgarawal/rocketoff/chat"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMakeHTTPHandler(t *testing.T) {
	e := Endpoints{}

	result := MakeHTTPHandler(e, &mockCommandDecoder{})
	assert.Implements(t, (*http.Handler)(nil), result)
}

type mockCommandDecoder struct {
	mock.Mock
}

func (m *mockCommandDecoder) DecodeCommand(_ context.Context, request *http.Request) (*chat.Command, error) {
	args := m.Called(request)
	return args.Get(0).(*chat.Command), args.Error(1)
}
