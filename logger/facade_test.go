package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFacade(t *testing.T) {
	Printf("test %q", "print")

	Debugf("test %q", "debug")
	Infof("test %q", "info")
	Warnf("test %q", "warn")
	assert.Panics(t, func() { Panicf("test %q", "panic") })
	Errorf("test %q", "error")

	Debugw("test", "some", "debug")
	Infow("test", "some", "info")
	Warnw("test", "some", "warn")
	assert.Panics(t, func() { Panicw("test", "some", "panic") })
	Errorw("test", "some", "error")
}
