package rocketoff

import (
	"github.com/nlopes/slack"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeHTTPHandler(t *testing.T) {
	e := Endpoints{}

	result := MakeHTTPHandler(e, "")
	assert.Implements(t, (*http.Handler)(nil), result)
}

func TestSlashCommandRequestDecoder(t *testing.T) {
	//TODO
	t.Error("unimplemented")
}

func TestEncodeResponse(t *testing.T) {
	//TODO
	t.Error("unimplemented")
}

func TestBuildSlackMsg(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		resp := &commandResponse{imageUrl: theBeardGif}
		msg, err := buildSlackMsg(resp)
		expected := &slack.Msg{
			ResponseType: slack.ResponseTypeInChannel,
			Attachments: []slack.Attachment{
				{
					ImageURL: theBeardGif,
				},
			},
		}

		assert.Nil(t, err)
		assert.Equal(t, expected, msg)
	})

	t.Run("happy path", func(t *testing.T) {
		_, err := buildSlackMsg(nil)

		assert.Equal(t, ErrInvalidType{"commandResponse"}, err)
	})
}
