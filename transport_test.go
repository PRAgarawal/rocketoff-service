package rocketoff

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeHTTPHandler(t *testing.T) {
	e := Endpoints{}

	result := MakeHTTPHandler(e, "")
	assert.Implements(t, (*http.Handler)(nil), result)
}

func TestSlashCommandRequestDecoder(t *testing.T) {
	//TODO
	t.Error("unimplemented")
}
