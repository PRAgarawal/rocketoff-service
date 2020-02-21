package rocketoff

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestEndpointMethods(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		e := Endpoints{
			ShowEmTheBeard: func(ctx context.Context, request interface{}) (response interface{}, err error) {
				return nil, nil
			},
			ShowEmThePointGod: func(ctx context.Context, request interface{}) (response interface{}, err error) {
				return nil, nil
			},
		}
		ctx := context.Background()

		resp, err := e.ShowEmTheBeard(ctx, nil)
		assert.Nil(t, err)
		assert.Nil(t, resp)

		resp, err = e.ShowEmThePointGod(ctx, nil)
		assert.Nil(t, err)
		assert.Nil(t, resp)
	})

	t.Run("error cases", func(t *testing.T) {
		e := Endpoints{
			ShowEmTheBeard: func(ctx context.Context, request interface{}) (response interface{}, err error) {
				return nil, errors.New("OMG HARDEN WAS ROBBED O")
			},

			ShowEmThePointGod: func(ctx context.Context, request interface{}) (repsonse interface{}, err error) {
				return nil, errors.New("OMG ")
			},
		}
		ctx := context.Background()

		resp, err := e.ShowEmTheBeard(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, resp)

		resp, err = e.ShowEmThePointGod(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}

func TestMakeServerEndpoints(t *testing.T) {
	mockSvc := new(mockService)
	e := MakeServerEndpoints(mockSvc)

	assert.NotNil(t, e.ShowEmTheBeard)
	assert.NotNil(t, e.ShowEmThePointGod)
}

func TestShowEmTheBearEndpoint(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockSvc := new(mockService)

		// set up expectations
		mockSvc.On("ShowEmTheBeard", mock.Anything).Return(nil)

		endpoint := makeShowEmTheBeardEndpoint(mockSvc)

		response, err := endpoint(context.Background(), nil)
		assert.Nil(t, err)
		assert.Nil(t, response)
	})

	t.Run("error response", func(t *testing.T) {
		mockSvc := new(mockService)

		// set up expectations
		mockSvc.On("ShowEmTheBeard", mock.Anything).Return(errors.New("JAMES HARDEN ROBBED YET AGAIN OF MVP"))

		endpoint := makeShowEmTheBeardEndpoint(mockSvc)

		response, err := endpoint(context.Background(), nil)
		assert.Error(t, err)
		assert.Nil(t, response)
	})
}

func TestShowEmThePointGodEndpoint(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockSvc := new(mockService)

		// set up expectations
		mockSvc.On("ShowEmThePointGod", mock.Anything).Return(nil)

		endpoint := makeShowEmThePointGodEndpoint(mockSvc)

		response, err := endpoint(context.Background(), nil)
		assert.Nil(t, err)
		assert.Nil(t, response)
	})

	t.Run("error response", func(t *testing.T) {
		mockSvc := new(mockService)

		// set up expectations
		mockSvc.On("ShowEmThePointGod", mock.Anything).Return(errors.New("CHRIS PAUL SUFFERS A HAMSTRING INJURY AND STARTS A FIGHT WITH TEAMMATE IN GAME 6"))

		endpoint := makeShowEmThePointGodEndpoint(mockSvc)

		response, err := endpoint(context.Background(), nil)
		assert.Error(t, err)
		assert.Nil(t, response)
	})
}

type mockService struct {
	mock.Mock
}

func (m *mockService) ShowEmTheBeard(_ context.Context) error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockService) ShowEmThePointGod(_ context.Context) error {
	args := m.Called()
	return args.Error(0)
}
