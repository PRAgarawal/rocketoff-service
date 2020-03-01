package slack

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PRAgarawal/rocketoff/chat"
	"github.com/nlopes/slack"
	"github.com/stretchr/testify/assert"
)

const (
	testUser = "JVG"
	testImage = "T-Mac.jpg"
)

func TestSendImageReply(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			headerContentType := r.Header.Get("Content-Type")
			assert.Equal(t, "application/json", headerContentType)
			body, err := ioutil.ReadAll(r.Body)
			assert.Nil(t, err, "Failed reading request body")
			expected := &slack.Msg{
				Text:         fmt.Sprintf("with warm regards from <@%s>", testUser),
				ResponseType: slack.ResponseTypeInChannel,
				Attachments: []slack.Attachment{
					{
						ImageURL: testImage,
					},
				},
			}
			expectedRaw, _ := json.Marshal(expected)
			assert.Equal(t, expectedRaw, body)
		}))

		reply := &chat.CommandReply{
			RequestingUserID: testUser,
			WebhookURL:       ts.URL,
			ImageURL:         testImage,
		}
		msgr := NewMessenger()

		assert.NoError(t, msgr.SendImageReply(reply))
	})
}
