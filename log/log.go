// Package log 是 SDK 的 logger 接口定义与内置的 logger。
package log

// DefaultLogger 默认logger
var DefaultLogger = Logger(&consoleLogger{Level: InfoLevel})

func Trace(v ...interface{}) {
	DefaultLogger.Trace(v...)
}

func Debug(v ...interface{}) {
	DefaultLogger.Debug(v...)
}

// Info log.Info
func Info(v ...interface{}) {
	DefaultLogger.Info(v...)
}

// Warn log.Warn
func Warn(v ...interface{}) {
	DefaultLogger.Warn(v...)
}

// Error log.Error
func Error(v ...interface{}) {
	DefaultLogger.Error(v...)
}

func Tracef(format string, v ...interface{}) {
	DefaultLogger.Tracef(format, v...)
}

// Debugf log.Debugf
func Debugf(format string, v ...interface{}) {
	DefaultLogger.Debugf(format, v...)
}

// Infof log.Infof
func Infof(format string, v ...interface{}) {
	DefaultLogger.Infof(format, v...)
}

// Warnf log.Warnf
func Warnf(format string, v ...interface{}) {
	DefaultLogger.Warnf(format, v...)
}

// Errorf log.Errorf
func Errorf(format string, v ...interface{}) {
	DefaultLogger.Errorf(format, v...)
}

// Sync logger Sync calls to flush buffer
func Sync() {
	_ = DefaultLogger.Sync()
}
