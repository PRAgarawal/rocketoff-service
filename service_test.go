package rocketoff

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PRAgarawal/rocketoff/chat"
	kitlog "github.com/go-kit/kit/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	testWebhook = "https://clutchCITY.net"
	testUser    = "Rudy T"
	redirectURI = "https://redirect.uri"
)

func TestShowEmTheBeard(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		msgr := &mockMessenger{}
		svc := New(kitlog.NewNopLogger(), msgr, nil)
		shouldAccept := &chat.CommandReply{
			WebhookURL:       testWebhook,
			RequestingUserID: testUser,
			ImageURL:         theBeardGif,
		}
		msgr.On("SendImageReply", shouldAccept).Return(nil)
		command := &ImageCommand{
			WebhookURL:       testWebhook,
			RequestingUserID: testUser,
		}

		assert.NoError(t, svc.ShowEmTheBeard(context.Background(), command))
	})

	t.Run("messenger error", func(t *testing.T) {
		msgr := &mockMessenger{}
		svc := New(kitlog.NewNopLogger(), msgr, nil)
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
		svc := New(kitlog.NewNopLogger(), msgr, nil)
		shouldAccept := &chat.CommandReply{
			WebhookURL:       testWebhook,
			RequestingUserID: testUser,
			ImageURL:         thePointGodGif,
		}
		msgr.On("SendImageReply", shouldAccept).Return(nil)
		command := &ImageCommand{
			WebhookURL:       testWebhook,
			RequestingUserID: testUser,
		}

		assert.NoError(t, svc.ShowEmThePointGod(context.Background(), command))
	})

	t.Run("messenger error", func(t *testing.T) {
		msgr := &mockMessenger{}
		svc := New(kitlog.NewNopLogger(), msgr, nil)
		msgr.On("SendImageReply", mock.Anything).Return(fmt.Errorf("OH GOD WESTBROOK IS OPEN FROM 18 FEET"))
		err := svc.ShowEmThePointGod(context.Background(), &ImageCommand{})

		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "WESTBROOK IS OPEN")
		}
	})
}

func TestCompleteChatOAuth(t *testing.T) {
	// Set up a service to respond to authorization and token exchange requests in these tests
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.String() != "/token" {
			t.Errorf("Unexpected exchange request URL %q", r.URL)
		}
		headerAuth := r.Header.Get("Authorization")
		assert.Equal(t, "Basic MzozMzM=", headerAuth, fmt.Sprintf("Unexpected authorization header %q, want %q", headerAuth, "Basic MzozMzM="))
		headerContentType := r.Header.Get("Content-Type")
		assert.Equal(t, "application/x-www-form-urlencoded", headerContentType, fmt.Sprintf("Unexpected Content-Type header %q", headerContentType))
		body, err := ioutil.ReadAll(r.Body)
		assert.Nil(t, err, "Failed reading request body")
		assert.Equal(t, "code=code&grant_type=authorization_code", string(body), fmt.Sprintf("Unexpected exchange payload; got %q", body))
		w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
		w.Write([]byte("access_token=90d64460d14870c08c81352a05dedd3465940a7c&scope=user&token_type=bearer&refresh_token=80d64460d14870c07c81353a05dedd3465940a7d&expires_in=3600"))
	}))
	defer ts.Close()

	t.Run("happy path - completing chat OAuth", func(t *testing.T) {
		config := &ChatConfig{
			ClientID:                 "3",
			ClientSecret:             "333",
			TokenEndpoint:            ts.URL + "/token",
			OAuthCompleteRedirectURL: redirectURI,
		}
		svc := New(kitlog.NewNopLogger(), nil, config)
		oauthOptions := &OAuthCompleteOptions{Code: "code"}
		redirect, err := svc.CompleteChatOAuth(context.Background(), oauthOptions)

		assert.Nil(t, err)
		assert.Equal(t, redirectURI, redirect)
	})
}

func TestRedirectForOAuth(t *testing.T) {
	config := &ChatConfig{
		ClientID:                 "123",
		AuthorizationEndpoint:    "https://authorize.me/",
		OAuthRedirectURL:         "localhost/oauth_complete",
		Scopes:                   "scope1,scope2",
	}
	svc := New(kitlog.NewNopLogger(), nil, config)
	expectedURI := "https://authorize.me/?client_id=123&redirect_uri=localhost%2Foauth_complete&response_type=code&scope=scope1+scope2"
	redirect, err := svc.RedirectForOAuth(context.Background())

	assert.Equal(t, expectedURI, redirect)
	assert.Nil(t, err)
}

func TestBuildRedirectURI(t *testing.T) {
	baseURI := "http://redirect.uri?key1=value1&key2=value2"
	resultURI, err := buildRedirectURI(baseURI, ErrInvalidValue{"WE'RE ALL GONNA DIE"})

	assert.Nil(t, err)
	assert.Equal(t, "http://redirect.uri?error=invalid+value%3A+%27WE%27RE+ALL+GONNA+DIE%27&key1=value1&key2=value2", resultURI)
}

type mockMessenger struct {
	mock.Mock
}

func (m *mockMessenger) SendImageReply(reply *chat.CommandReply) error {
	args := m.Called(reply)
	return args.Error(0)
}
