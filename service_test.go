package rocketoff

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	kitlog "github.com/go-kit/kit/log"
	"github.com/nlopes/slack"
	"github.com/stretchr/testify/assert"
)

const (
	testUser = "Rudy T"
)

func TestShowEmTheBeard(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		ts := httptest.NewServer(makeHandlerFunc(theBeardGif, t))
		svc := New(kitlog.NewNopLogger())
		command := &ImageCommand{
			ResponseURL:        ts.URL,
			RequestingUsername: testUser,
		}
		err := svc.ShowEmTheBeard(context.Background(), command)

		assert.NoError(t, err)
	})
}

func TestShowEmThePointGod(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		ts := httptest.NewServer(makeHandlerFunc(thePointGodGif, t))
		svc := New(kitlog.NewNopLogger())
		command := &ImageCommand{
			ResponseURL:        ts.URL,
			RequestingUsername: testUser,
		}
		err := svc.ShowEmThePointGod(context.Background(), command)

		assert.NoError(t, err)
	})
}

func makeHandlerFunc(imageURL string, t *testing.T) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		headerContentType := r.Header.Get("Content-Type")
		assert.Equal(t, "application/json", headerContentType)
		body, err := ioutil.ReadAll(r.Body)
		assert.Nil(t, err, "Failed reading request body")
		expected := &slack.Msg{
			Text:         fmt.Sprintf("with warm regards from %s", testUser),
			ResponseType: slack.ResponseTypeInChannel,
			Attachments: []slack.Attachment{
				{
					ImageURL: imageURL,
				},
			},
		}
		expectedRaw, _ := json.Marshal(expected)
		assert.Equal(t, expectedRaw, body)
	}
}
