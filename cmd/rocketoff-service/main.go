package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/PRAgarawal/rocketoff"
	"github.com/PRAgarawal/rocketoff/chat/slack"
	kitlog "github.com/go-kit/kit/log"
)

const (
	// TODO build this URL in a service method... this is ugly
	slackOAuthURL = "https://slack.com/oauth/v2/authorize?client_id=%s&scope=chat:write,commands"
)

func main() {
	var (
		signingSecret = flag.String("slack-signing-secret", os.Getenv("SLACK_SIGNING_SECRET"), "slack application signing secret to verify web requests")
		clientID      = flag.String("client-id", os.Getenv("CLIENT_ID"), "application client ID for OAuth requests")
		clientSecret  = flag.String("client-secret", os.Getenv("CLIENT_SECRET"), "application client secret for OAuth requests")
		oauthTokenEndpoint = flag.String("oauth-token-endpoint", os.Getenv("OAUTH_TOKEN_ENDPOINT"), "URL to use to procure oauth2 access tokens")
		oauthCompleteURL   = flag.String("oauth-redirect-url", os.Getenv("OAUTH_REDIRECT_URL"), "URL to take the user to after they have successfully signed up for the app")
		//TODO: Implement for keybase, and select the chat application via a flag?
	)
	logger := kitlog.With(kitlog.NewJSONLogger(os.Stderr), "ts", kitlog.DefaultTimestampUTC)
	chatConfig := &rocketoff.ChatConfig{
		ClientID:                 *clientID,
		ClientSecret:             *clientSecret,
		TokenEndpoint:            *oauthTokenEndpoint,
		OAuthCompleteRedirectURL: *oauthCompleteURL,
	}
	rocketoffService := rocketoff.New(logger, slack.NewMessenger(), chatConfig)
	integrationEndpoints := rocketoff.MakeServerEndpoints(rocketoffService)

	mux := http.NewServeMux()
	commandDecoder := slack.NewCommandDecoder(*signingSecret)
	slackOAuthCompleteURL := fmt.Sprintf(slackOAuthURL, *clientID)
	mux.Handle("/", rocketoff.MakeHTTPHandler(integrationEndpoints, commandDecoder, slackOAuthCompleteURL))

	server := &http.Server{Addr: ":8080", Handler: mux}
	listenAndServeGracefully(server, logger)
}

// listenAndServeGracefully starts the provided http server and listens for SIGINT or SIGTERM
// upon receiving either signal it will stop accepting new connections and will wait for existing
// connections to close for up to 15 seconds before forcibly closing connections.
func listenAndServeGracefully(server *http.Server, logger kitlog.Logger) {
	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)
		<-sigint
		logger.Log(
			"message", "stopping service",
			"severity", "NOTICE",
		)

		sdCtx, cancel := context.WithTimeout(context.Background(), time.Second*15)
		defer cancel()
		if err := server.Shutdown(sdCtx); err != nil {
			logger.Log(
				"message", "http server stopped",
				"err", err,
				"severity", "CRITICAL",
			)
		}

		close(idleConnsClosed)
	}()

	logger.Log(
		"message", fmt.Sprintf("listening for HTTP connections on %s", server.Addr),
		"severity", "NOTICE",
	)

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		logger.Log(
			"message", "service failure",
			"error", err,
			"severity", "CRITICAL",
		)
	} else {
		logger.Log(
			"message", "service stopped",
			"severity", "NOTICE",
		)
	}

	<-idleConnsClosed
}
