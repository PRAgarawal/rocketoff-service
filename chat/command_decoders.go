package chat

import (
	"context"
	"net/http"
)

// CommandDecoder
type CommandDecoder interface {
	// DecodeCommand takes an http.Request and parses a chat Command request, after performing whatever verification or pre-decode steps may be necessary.
	DecodeCommand(ctx context.Context, r *http.Request) (*Command, error)
}

type Command struct {
	WebhookURL string
	RequestingUserName string
}
