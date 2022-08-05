package logger

import (
	"sync"

	"go.uber.org/zap"
)

// Logger represents a standard logging interface.
type Logger interface {
	// Log a notice statement
	Noticef(format string, v ...interface{})

	// Log a warning statement
	Warnf(format string, v ...interface{})

	// Log a fatal error
	Fatalf(format string, v ...interface{})

	// Log an error
	Errorf(format string, v ...interface{})

	// Log a debug statement
	Debugf(format string, v ...interface{})

	// Log a trace statement
	Tracef(format string, v ...interface{})
}

var _ Logger = (*logger)(nil)

type logger struct {
	opts *Opts
	sync.RWMutex
}

// Opt ...
type Opt func(*Opts)

// Opts ...
type Opts struct {
	Logger *zap.Logger
}

// Configure ...
func (o *Opts) configure(opts ...Opt) {
	for _, opt := range opts {
		opt(o)
	}
}

// NewLogger ...
func NewLogger(o ...Opt) Logger {
	options := new(Opts)
	options.configure(o...)

	l := new(logger)
	l.opts = options

	return l
}

// Errorf ...
func (l *logger) Errorf(format string, v ...interface{}) {
	l.logFunc(func(log *zap.Logger, format string, v ...interface{}) {
		log.Sugar().Errorf(format, v...)
	}, format, v...)
}

// Debugf ...
func (l *logger) Debugf(format string, v ...interface{}) {
	l.logFunc(func(log *zap.Logger, format string, v ...interface{}) {
		log.Sugar().Debugf(format, v...)
	}, format, v...)
}

// Fatalf ...
func (l *logger) Fatalf(format string, v ...interface{}) {
	l.logFunc(func(log *zap.Logger, format string, v ...interface{}) {
		log.Sugar().Fatalf(format, v...)
	}, format, v...)
}

// Noticef ...
func (l *logger) Noticef(format string, v ...interface{}) {
	l.logFunc(func(log *zap.Logger, format string, v ...interface{}) {
		log.Sugar().Infof(format, v...)
	}, format, v...)
}

// Warnf ...
func (l *logger) Warnf(format string, v ...interface{}) {
	l.logFunc(func(log *zap.Logger, format string, v ...interface{}) {
		log.Sugar().Warnf(format, v...)
	}, format, v...)
}

// Tracef ...
func (l *logger) Tracef(format string, v ...interface{}) {
	l.logFunc(func(log *zap.Logger, format string, v ...interface{}) {
		log.Sugar().Debugf(format, v...)
	}, format, v...)
}

func (l *logger) logFunc(f func(log *zap.Logger, format string, v ...interface{}), format string, args ...interface{}) {
	l.Lock()
	defer l.Unlock()

	if l.opts.Logger == nil {
		return
	}

	f(l.opts.Logger, format, args...)
}
