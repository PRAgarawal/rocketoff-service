package rocketoff

import (
	"context"

	kitlog "github.com/go-kit/kit/log"
)

const (
	theBeardGif = "https://i.imgur.com/t7ddUae.gif"
	thePointGodGif = "https://i.imgur.com/c2qPNN2.gif"
)

// Service contains simple methods that can be used to generate responses to chat app slash commands.
type Service interface {
	// ShowEmTheBeard will return a gif of James Harden, reminding those in the channel that they're lame. Not all of them. They will know who they are.
	ShowEmTheBeard(ctx context.Context) (*ImageReply, error)

	// ShowEmThePointGod will return a gif of Chris Paul, reminding those in the channel that they are not funny. Not all of them. They will know who they are.
	ShowEmThePointGod(ctx context.Context) (*ImageReply, error)
}

type ImageReply struct {
	ImageURL string
}

type Svc struct {
	logger kitlog.Logger
}

func New(logger kitlog.Logger) Service {
	return &Svc{logger}
}

func (s *Svc) ShowEmTheBeard(_ context.Context) (*ImageReply, error) {
	return &ImageReply{theBeardGif}, nil
}

func (s *Svc) ShowEmThePointGod(_ context.Context) (*ImageReply, error) {
	return &ImageReply{thePointGodGif}, nil
}
