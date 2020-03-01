package rocketoff

import (
	"context"

	"github.com/PRAgarawal/rocketoff/chat"
	kitlog "github.com/go-kit/kit/log"
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
}

type ImageCommand struct {
	// WebhookURL specifies the (likely ephemeral) URL provided by the chat application to use when interacting with the commanding user. POSTing messages to this URL should send messages for some period of time.
	WebhookURL string

	// RequestingUserID is the ID of the chat app user who sent the command
	RequestingUserID string
}

type Svc struct {
	logger kitlog.Logger
	msgr   chat.Messenger
}

func New(logger kitlog.Logger, msgr chat.Messenger) Service {
	return &Svc{
		logger: logger,
		msgr:   msgr,
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
