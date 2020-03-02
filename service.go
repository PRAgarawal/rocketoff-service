package rocketoff

import (
	"context"
	"net/url"
	"strings"

	"github.com/PRAgarawal/rocketoff/chat"
	kitlog "github.com/go-kit/kit/log"
	"golang.org/x/oauth2"
)

const (
	theBeardGif    = "https://i.imgur.com/t7ddUae.gif"
	thePointGodGif = "https://i.imgur.com/c2qPNN2.gif"
)

// Service contains simple methods that can be used to generate responses to chat app slash commands.
type Service interface {
	// ShowEmTheBeard will reply with the specified URL a gif of James Harden, reminding those in the chat that they're lame. Not all of them. They will know who they are.
	ShowEmTheBeard(ctx context.Context, command *ImageCommand) error

	// ShowEmThePointGod will reply with a gif of Chris Paul, reminding those in the chat that they are not funny. Not all of them. They will know who they are.
	ShowEmThePointGod(ctx context.Context, command *ImageCommand) error

	// CompleteChatOAuth takes the requesting user's authorization code an retrieves an access token. This process should complete the signup for the application, but we don't actually need the authorization token.
	CompleteChatOAuth(ctx context.Context, oauthOptions *OAuthCompleteOptions) (string, error)

	// RedirectForOAuth builds the oauth authentication URI that a chat app store may use to initiate the authentication process
	RedirectForOAuth(_ context.Context) (string, error)
}

type ImageCommand struct {
	// WebhookURL specifies the (likely ephemeral) URL provided by the chat application to use when interacting with the commanding user. POSTing messages to this URL should send messages for some period of time.
	WebhookURL string

	// RequestingUserID is the ID of the chat app user who sent the command
	RequestingUserID string
}

// OAuthCompleteOptions contains the values provided by the external OAuth2 authorization server that we will use to complete authentication of the corresponding integration, and procure the user's OAuth2 tokens.
type OAuthCompleteOptions struct {
	// Code is the OAuth2 authorization code from the authorization server. Rocketoff will use this value to collect the access and refresh tokens for the given integration
	Code string

	// State corresponds to the CSRFState value on the integration, and it's not very important for the purposes of this application
	State string
}

type ChatConfig struct {
	ClientID                 string
	ClientSecret             string
	AuthorizationEndpoint    string
	TokenEndpoint            string
	OAuthRedirectURL         string
	OAuthCompleteRedirectURL string
	Scopes                   string
}

type Svc struct {
	logger kitlog.Logger
	msgr   chat.Messenger
	config *ChatConfig
}

func New(logger kitlog.Logger, msgr chat.Messenger, config *ChatConfig) Service {
	return &Svc{
		logger: logger,
		msgr:   msgr,
		config: config,
	}
}

func (s *Svc) ShowEmTheBeard(_ context.Context, command *ImageCommand) error {
	reply := &chat.CommandReply{
		RequestingUserID: command.RequestingUserID,
		WebhookURL:       command.WebhookURL,
		ImageURL:         theBeardGif,
	}
	if err := s.msgr.SendImageReply(reply); err != nil {
		return err
	}

	s.logger.Log(
		"message", "successfully sent The Beard",
		"severity", "INFO",
	)

	return nil
}

func (s *Svc) ShowEmThePointGod(_ context.Context, command *ImageCommand) error {
	reply := &chat.CommandReply{
		RequestingUserID: command.RequestingUserID,
		WebhookURL:       command.WebhookURL,
		ImageURL:         thePointGodGif,
	}
	if err := s.msgr.SendImageReply(reply); err != nil {
		return err
	}

	s.logger.Log(
		"message", "successfully sent The Beard",
		"severity", "INFO",
	)

	return nil
}

// CompleteChatOauth performs the actual token exchange to get the [useless for now] tokens
func (s *Svc) CompleteChatOAuth(ctx context.Context, oauthOptions *OAuthCompleteOptions) (string, error) {
	conf := s.OAuth2Config()
	// We don't actually need the token to use webhook URLs to post messages
	_, err := conf.Exchange(ctx, oauthOptions.Code)
	return buildRedirectURI(s.config.OAuthCompleteRedirectURL, err)
}

// RedirectForOAuth uses an oauth2 config to build and return the oauth authentication URI
func (s *Svc) RedirectForOAuth(_ context.Context) (string, error) {
	conf := s.OAuth2Config()
	return conf.AuthCodeURL(""), nil
}

// buildRedirectURI takes a base redirect URI and will add an `error` query parameter in the return URI string if an err is provided
func buildRedirectURI(redirectBaseURI string, err error) (string, error) {
	if err == nil {
		return redirectBaseURI, nil
	}

	// There was an error, so add it to the redirect URI
	redirectURL, parseErr := url.Parse(redirectBaseURI)
	if parseErr != nil {
		return "", parseErr
	}
	query, parseErr := url.ParseQuery(redirectURL.RawQuery)
	if parseErr != nil {
		return "", parseErr
	}
	query.Add("error", err.Error())
	redirectURL.RawQuery = query.Encode()

	return redirectURL.String(), nil
}

func (s *Svc) OAuth2Config() *oauth2.Config {
	conf := &oauth2.Config{}
	conf.Endpoint = oauth2.Endpoint{
		AuthURL:  s.config.AuthorizationEndpoint,
		TokenURL: s.config.TokenEndpoint,
	}
	conf.ClientID = s.config.ClientID
	conf.ClientSecret = s.config.ClientSecret
	conf.RedirectURL = s.config.OAuthRedirectURL
	conf.Scopes = strings.Split(s.config.Scopes, ",")

	return conf
}
