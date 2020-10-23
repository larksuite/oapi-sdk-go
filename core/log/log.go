package log

import (
	"context"
	"log"
	"os"
)

type Level int

const (
	LevelDebug Level = 1
	LevelInfo  Level = 2
	LevelWarn  Level = 3
	LevelError Level = 4
)

type LoggerProxy struct {
	LogLevel Level
	Logger   Logger
}

func NewLoggerProxy(logLevel Level, logger Logger) *LoggerProxy {
	return &LoggerProxy{
		LogLevel: logLevel,
		Logger:   logger,
	}
}

func (p *LoggerProxy) Debug(ctx context.Context, args ...interface{}) {
	if p.LogLevel <= LevelDebug {
		p.Logger.Debug(ctx, args...)
	}
}

func (p *LoggerProxy) Info(ctx context.Context, args ...interface{}) {
	if p.LogLevel <= LevelInfo {
		p.Logger.Info(ctx, args...)
	}
}

func (p *LoggerProxy) Warn(ctx context.Context, args ...interface{}) {
	if p.LogLevel <= LevelWarn {
		p.Logger.Warn(ctx, args...)
	}
}

func (p *LoggerProxy) Error(ctx context.Context, args ...interface{}) {
	if p.LogLevel <= LevelError {
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
