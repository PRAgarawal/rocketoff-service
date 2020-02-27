package rocketoff

import (
	"context"
	"testing"

	kitlog "github.com/go-kit/kit/log"
	"github.com/stretchr/testify/assert"
)

func TestShowEmTheBeard(t *testing.T) {
	svc := New(kitlog.NewNopLogger())
	imgReply, err := svc.ShowEmTheBeard(context.Background())
	expected := &ImageReply{ImageURL:theBeardGif}

	assert.NoError(t, err)
	assert.Equal(t, expected, imgReply)
}

func TestShowEmThePointGod(t *testing.T) {
	svc := New(kitlog.NewNopLogger())
	imgReply, err := svc.ShowEmThePointGod(context.Background())
	expected := &ImageReply{ImageURL:thePointGodGif}

	assert.NoError(t, err)
	assert.Equal(t, expected, imgReply)
}
