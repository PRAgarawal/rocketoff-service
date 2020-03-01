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

// Command is all the relevant data we need from a command a chat user has issued
type Command struct {
	// WebhookURL is the ephemeral URL created for this command interaction on the given chat application
	WebhookURL       string

	// RequestingUserID is the chat app's unique ID for the user we will reply to
	RequestingUserID string
}
