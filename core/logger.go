package core

import (
	"context"
	"log"
	"os"
)

type LogLevel int

const (
	LogLevelDebug LogLevel = 1
	LogLevelInfo  LogLevel = 2
	LogLevelWarn  LogLevel = 3
	LogLevelError LogLevel = 4
)

type loggerProxy struct {
	LogLevel LogLevel
	Logger   Logger
}

func newLoggerProxy(logLevel LogLevel, logger Logger) *loggerProxy {
	return &loggerProxy{
		LogLevel: logLevel,
		Logger:   logger,
	}
}

func (p *loggerProxy) Debug(ctx context.Context, args ...interface{}) {
	if p.Logger != nil && p.LogLevel <= LogLevelDebug {
		p.Logger.Debug(ctx, args...)
	}
}

func (p *loggerProxy) Info(ctx context.Context, args ...interface{}) {
	if p.Logger != nil && p.LogLevel <= LogLevelInfo {
		p.Logger.Info(ctx, args...)
	}
}

func (p *loggerProxy) Warn(ctx context.Context, args ...interface{}) {
	if p.Logger != nil && p.LogLevel <= LogLevelWarn {
		p.Logger.Warn(ctx, args...)
	}
}

func (p *loggerProxy) Error(ctx context.Context, args ...interface{}) {
	if p.Logger != nil && p.LogLevel <= LogLevelError {
		p.Logger.Error(ctx, args...)
	}
}

type Logger interface {
	Debug(context.Context, ...interface{})
	Info(context.Context, ...interface{})
	Warn(context.Context, ...interface{})
	Error(context.Context, ...interface{})
}

func NewLogger(config *Config) {
	if config.Logger != nil {
		logLevel := LogLevelInfo
		if config.LogLevel != 0 {
			logLevel = config.LogLevel
		}
		logger := newLoggerProxy(logLevel, config.Logger)
		config.Logger = logger
	} else {
		logLevel := LogLevelInfo
		if config.LogLevel != 0 {
			logLevel = config.LogLevel
		}
		logger := newLoggerProxy(logLevel, defaultLogger{
			logger: log.New(os.Stdout, "", log.LstdFlags),
		})
		config.Logger = logger
	}

}

func NewEventLogger() Logger {
	logger := newLoggerProxy(LogLevelInfo, defaultLogger{
		logger: log.New(os.Stdout, "", log.LstdFlags),
	})
	return logger
}

type defaultLogger struct {
	logger *log.Logger
}

func (l defaultLogger) Debug(ctx context.Context, args ...interface{}) {
	l.logger.Printf("[Debug] %v", args)
}

func (l defaultLogger) Info(ctx context.Context, args ...interface{}) {
	l.logger.Printf("[Info] %v", args)
}

func (l defaultLogger) Warn(ctx context.Context, args ...interface{}) {
	l.logger.Printf("[Warn] %v", args)
}

func (l defaultLogger) Error(ctx context.Context, args ...interface{}) {
	l.logger.Printf("[Error] %v", args)
}
