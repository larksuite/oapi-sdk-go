package lark

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

type LoggerProxy struct {
	LogLevel LogLevel
	Logger   Logger
}

func NewLoggerProxy(logLevel LogLevel, logger Logger) *LoggerProxy {
	return &LoggerProxy{
		LogLevel: logLevel,
		Logger:   logger,
	}
}

func (p *LoggerProxy) Debug(ctx context.Context, args ...interface{}) {
	if p.LogLevel <= LogLevelDebug && p.Logger != nil {
		p.Logger.Debug(ctx, args...)
	}
}

func (p *LoggerProxy) Info(ctx context.Context, args ...interface{}) {
	if p.LogLevel <= LogLevelInfo && p.Logger != nil {
		p.Logger.Info(ctx, args...)
	}
}

func (p *LoggerProxy) Warn(ctx context.Context, args ...interface{}) {
	if p.LogLevel <= LogLevelWarn && p.Logger != nil {
		p.Logger.Warn(ctx, args...)
	}
}

func (p *LoggerProxy) Error(ctx context.Context, args ...interface{}) {
	if p.LogLevel <= LogLevelError && p.Logger != nil {
		p.Logger.Error(ctx, args...)
	}
}

type Logger interface {
	Debug(context.Context, ...interface{})
	Info(context.Context, ...interface{})
	Warn(context.Context, ...interface{})
	Error(context.Context, ...interface{})
}

func NewDefaultLogger() Logger {
	return &defaultLogger{
		logger: log.New(os.Stdout, "", log.LstdFlags),
	}
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
