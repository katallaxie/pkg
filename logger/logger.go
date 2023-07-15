package logger

import (
	"sync"

	"go.uber.org/zap"
)

// LogSink is the logger sink.
var LogSink *zap.Logger

func init() {
	l, err := NewLogSink()
	if err != nil {
		panic(err)
	}

	LogSink = l
}

// NewLogSink returns a new logger sink.
func NewLogSink() (*zap.Logger, error) {
	return zap.NewProduction()
}

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

// LogFunc is a bridge between Logger and any third party logger.
type LogFunc func(string, ...interface{})

// Printf is a bridge between Logger and any third party logger.
func (f LogFunc) Printf(msg string, args ...interface{}) { f(msg, args...) }

type logger struct {
	opts *Opts
	sync.RWMutex
}

// Opt is a logger option.
type Opt func(*Opts)

// Opts are the options for the logger.
type Opts struct {
	Logger *zap.Logger
}

// Configure is configuring the logger.
func (o *Opts) Configure(opts ...Opt) {
	for _, opt := range opts {
		opt(o)
	}
}

// NewLogger is creating a new logger.
func NewLogger(o ...Opt) Logger {
	options := new(Opts)
	options.Logger = LogSink

	options.Configure(o...)

	l := new(logger)
	l.opts = options

	return l
}

// Errorf is logging an error.
func (l *logger) Errorf(format string, v ...interface{}) {
	l.logFunc(func(log *zap.Logger, format string, v ...interface{}) {
		log.Sugar().Errorf(format, v...)
	}, format, v...)
}

// Debugf is logging a debug statement.
func (l *logger) Debugf(format string, v ...interface{}) {
	l.logFunc(func(log *zap.Logger, format string, v ...interface{}) {
		log.Sugar().Debugf(format, v...)
	}, format, v...)
}

// Fatalf is logging a fatal error.
func (l *logger) Fatalf(format string, v ...interface{}) {
	l.logFunc(func(log *zap.Logger, format string, v ...interface{}) {
		log.Sugar().Fatalf(format, v...)
	}, format, v...)
}

// Noticef is logging a notice statement.
func (l *logger) Noticef(format string, v ...interface{}) {
	l.logFunc(func(log *zap.Logger, format string, v ...interface{}) {
		log.Sugar().Infof(format, v...)
	}, format, v...)
}

// Warnf is logging a warning statement.
func (l *logger) Warnf(format string, v ...interface{}) {
	l.logFunc(func(log *zap.Logger, format string, v ...interface{}) {
		log.Sugar().Warnf(format, v...)
	}, format, v...)
}

// Tracef is logging a trace statement.
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
