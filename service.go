package rocketoff

import (
	"context"

	kitlog "github.com/go-kit/kit/log"

	"github.com/nlopes/slack"
)

type Service interface {
	// ShowEmTheBeard will post to the slack channel specified by its input a gif of James Harden, reminding those in the channel that they're lame. Not all of them. They will know who they are.
	ShowEmTheBeard(ctx context.Context) error

	// ShowEmThePointGod will post to the slack channel specified by its input a gif of Chris Paul, reminding those in the channel that they are not funny. Not all of them. They will know who they are.
	ShowEmThePointGod(ctx context.Context) error
}

type Svc struct {
	logger kitlog.Logger
	rtm    *slack.RTM
}

func New(logger kitlog.Logger, rtm *slack.RTM) Service {
	return &Svc{
		logger,
		rtm,
	}
}

func (s *Svc) ShowEmTheBeard(ctx context.Context) error {
	println("BEARD INCOMING")
	return nil
}

func (s *Svc) ShowEmThePointGod(ctx context.Context) error {
	println("POINT GOD INCOMING")
	return nil
}
