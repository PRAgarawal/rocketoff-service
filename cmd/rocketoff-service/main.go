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
	kitlog "github.com/go-kit/kit/log"
)

func main() {
	var (
		signingSecret = flag.String("slack-signing-secret", os.Getenv("SLACK_SIGNING_SECRET"), "slack application signing secret to verify web requests")
	)
	logger := kitlog.With(kitlog.NewJSONLogger(os.Stderr), "ts", kitlog.DefaultTimestampUTC)
	rocketoffService := rocketoff.New(logger)
	integrationEndpoints := rocketoff.MakeServerEndpoints(rocketoffService)

	mux := http.NewServeMux()
	mux.Handle("/", rocketoff.MakeHTTPHandler(integrationEndpoints, *signingSecret))

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
