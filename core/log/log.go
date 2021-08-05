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
	logLevel Level
	logger   Logger
}

func (p *LoggerProxy) SetLogLevel(l Level) {
	p.logLevel = l
}

func (p *LoggerProxy) SetLogger(logger Logger) {
	p.logger = logger
}

func NewLoggerProxy(logLevel Level, logger Logger) *LoggerProxy {
	return &LoggerProxy{
		logLevel: logLevel,
		logger:   logger,
	}
}

func (p *LoggerProxy) Debug(ctx context.Context, args ...interface{}) {
	if p.logLevel <= LevelDebug {
		p.logger.Debug(ctx, args...)
	}
}

func (p *LoggerProxy) Info(ctx context.Context, args ...interface{}) {
	if p.logLevel <= LevelInfo {
		p.logger.Info(ctx, args...)
	}
}

func (p *LoggerProxy) Warn(ctx context.Context, args ...interface{}) {
	if p.logLevel <= LevelWarn {
		p.logger.Warn(ctx, args...)
	}
}

func (p *LoggerProxy) Error(ctx context.Context, args ...interface{}) {
	if p.logLevel <= LevelError {
		p.logger.Error(ctx, args...)
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
