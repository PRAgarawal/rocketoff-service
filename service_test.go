package rocketoff

import (
	"context"
	"fmt"
	"testing"

	"github.com/PRAgarawal/rocketoff/chat"
	kitlog "github.com/go-kit/kit/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	testWebhook = "https://clutchCITY.net"
	testUser    = "Rudy T"
)

func TestShowEmTheBeard(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		msgr := &mockMessenger{}
		svc := New(kitlog.NewNopLogger(), msgr)
		shouldAccept := &chat.CommandReply{
			WebhookURL:         testWebhook,
			RequestingUserName: testUser,
			ImageURL: theBeardGif,
		}
		msgr.On("SendImageReply", shouldAccept).Return(nil)
		command := &ImageCommand{
			WebhookURL:         testWebhook,
			RequestingUserName: testUser,
		}

		assert.NoError(t, svc.ShowEmTheBeard(context.Background(), command))
	})

	t.Run("messenger error", func(t *testing.T) {
		msgr := &mockMessenger{}
		svc := New(kitlog.NewNopLogger(), msgr)
		msgr.On("SendImageReply", mock.Anything).Return(fmt.Errorf("A BASIC ASS CASUAL POSTED A YOUTUBE CLIP TO PROVE HARDEN PLAYS NO DEFENSE AND FLOPS"))
		err := svc.ShowEmTheBeard(context.Background(), &ImageCommand{})

		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "BASIC ASS")
		}
	})
}

func TestShowEmThePointGod(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		msgr := &mockMessenger{}
		svc := New(kitlog.NewNopLogger(), msgr)
		shouldAccept := &chat.CommandReply{
			WebhookURL:         testWebhook,
			RequestingUserName: testUser,
			ImageURL: thePointGodGif,
		}
		msgr.On("SendImageReply", shouldAccept).Return(nil)
		command := &ImageCommand{
			WebhookURL:         testWebhook,
			RequestingUserName: testUser,
		}

		assert.NoError(t, svc.ShowEmThePointGod(context.Background(), command))
	})

	t.Run("messenger error", func(t *testing.T) {
		msgr := &mockMessenger{}
		svc := New(kitlog.NewNopLogger(), msgr)
		msgr.On("SendImageReply", mock.Anything).Return(fmt.Errorf("OH GOD WESTBROOK IS OPEN FROM 18 FEET"))
		err := svc.ShowEmThePointGod(context.Background(), &ImageCommand{})

		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "WESTBROOK IS OPEN")
		}
	})
}

type mockMessenger struct {
	mock.Mock
}

func (m *mockMessenger) SendImageReply(reply *chat.CommandReply) error {
	args := m.Called(reply)
	return args.Error(0)
}
