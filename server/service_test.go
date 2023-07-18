package server

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLookupName(t *testing.T) {
	assert.Equal(t, "server.test", Service.Name())

	os.Setenv("NAME", "test")

	env := ServiceEnv{"NAME"}

	Service.Lookup(env...)
	assert.Equal(t, "test", Service.Name())

	os.Setenv("NAME", "foo")
	Service.Lookup(env...)
	assert.NotEqual(t, "foo", Service.Name())
}

func TestDefaultEnv(t *testing.T) {
	assert.Equal(t, "test", Service.Name())
}
