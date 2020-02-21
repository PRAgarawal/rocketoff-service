package rocketoff

import (
	"context"
	"testing"

	kitlog "github.com/go-kit/kit/log"
	"github.com/stretchr/testify/assert"
)

func TestShowEmTheBeard(t *testing.T) {
	svc := New(kitlog.NewNopLogger(), nil)

	assert.NoError(t, svc.ShowEmTheBeard(context.Background()))
}

func TestShowEmThePointGod(t *testing.T) {
	svc := New(kitlog.NewNopLogger(), nil)

	assert.NoError(t, svc.ShowEmThePointGod(context.Background()))
}
