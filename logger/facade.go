package logger

// Printf ..
func Printf(format string, args ...interface{}) {
	LogSink.Sugar().Infof(format, args...)
}

// Debugf ...
func Debugf(format string, args ...interface{}) {
	LogSink.Sugar().Debugf(format, args...)
}

// Infof ...
func Infof(format string, args ...interface{}) {
	LogSink.Sugar().Infof(format, args...)
}

// Errorf ...
func Errorf(format string, args ...interface{}) {
	LogSink.Sugar().Errorf(format, args...)
}

// Warnf ...
func Warnf(format string, args ...interface{}) {
	LogSink.Sugar().Warnf(format, args...)
}

// Panicf ...
func Panicf(format string, args ...interface{}) {
	LogSink.Sugar().Panicf(format, args...)
}

// Fatalf ...
func Fatalf(format string, args ...interface{}) {
	LogSink.Sugar().Fatalf(format, args...)
}

// Debugw ...
func Debugw(msg string, keysAndValues ...interface{}) {
	LogSink.Sugar().Debugw(msg, keysAndValues...)
}

// Infow ...
func Infow(msg string, keysAndValues ...interface{}) {
	LogSink.Sugar().Infow(msg, keysAndValues...)
}

// Warnw ...
func Warnw(msg string, keysAndValues ...interface{}) {
	LogSink.Sugar().Warnw(msg, keysAndValues...)
}

// Errorw ...
func Errorw(msg string, keysAndValues ...interface{}) {
	LogSink.Sugar().Errorw(msg, keysAndValues...)
}

// DPanicw ...
func DPanicw(msg string, keysAndValues ...interface{}) {
	LogSink.Sugar().DPanicw(msg, keysAndValues...)
}

// Panicw ...
func Panicw(msg string, keysAndValues ...interface{}) {
	LogSink.Sugar().Panicw(msg, keysAndValues...)
}

// Fatalw ...
func Fatalw(msg string, keysAndValues ...interface{}) {
	LogSink.Sugar().Fatalw(msg, keysAndValues...)
}
