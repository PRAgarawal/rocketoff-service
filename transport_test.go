package rocketoff

import (
	"context"
	"fmt"
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

func TestCommandDecoder(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		request := &http.Request{}
		mockCD := &mockCommandDecoder{}
		mockCD.On("DecodeCommand", request).Return(&chat.Command{}, nil)
		decoder := makeSlashCommandRequestDecoder(mockCD)
		req, err := decoder(context.TODO(), request)

		assert.Nil(t, err)
		assert.Equal(t, &commandRequest{}, req)
	})

	t.Run("decoder error", func(t *testing.T) {
		request := &http.Request{}
		mockCD := &mockCommandDecoder{}
		mockCD.On("DecodeCommand", mock.Anything).Return(nil, fmt.Errorf("OMG WE'RE $5 INTO LUXURY TAX TERRITORY"))
		decoder := makeSlashCommandRequestDecoder(mockCD)
		req, err := decoder(context.TODO(), request)

		assert.Nil(t, req)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "LUXURY TAX")
		}
	})
}

type mockCommandDecoder struct {
	mock.Mock
}

func (m *mockCommandDecoder) DecodeCommand(_ context.Context, request *http.Request) (*chat.Command, error) {
	args := m.Called(request)
	command, _ := args.Get(0).(*chat.Command)
	return command, args.Error(1)
}
