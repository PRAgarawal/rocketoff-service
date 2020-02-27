package slack

import (
	"context"
	"github.com/PRAgarawal/rocketoff/chat"
	"github.com/nlopes/slack"
	"io"
	"io/ioutil"
	"net/http"
)

type CommandDecoder struct {
	signingSecret string
}

func NewCommandDecoder(signingSecret string) chat.CommandDecoder {
	return &CommandDecoder{
		signingSecret: signingSecret,
	}
}

func (cd *CommandDecoder) DecodeCommand(_ context.Context, request *http.Request) (*chat.Command, error) {
	verifier, err := slack.NewSecretsVerifier(request.Header, cd.signingSecret)
	if err != nil {
		return nil, err
	}

	request.Body = ioutil.NopCloser(io.TeeReader(request.Body, &verifier))
	command, err := slack.SlashCommandParse(request)
	if err != nil {
		return nil, err
	}

	if err = verifier.Ensure(); err != nil {
		return nil, err
	}

	return &chat.Command{
		WebhookURL:         command.ResponseURL,
		RequestingUserName: command.UserName,
	}, nil
}
