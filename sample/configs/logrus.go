package configs

import (
	"context"
	"github.com/sirupsen/logrus"
)

// use logrus implement log.Logger
type Logrus struct {
}

func (Logrus) Debug(ctx context.Context, args ...interface{}) {
	logrus.Debug(args...)
}
func (Logrus) Info(ctx context.Context, args ...interface{}) {
	logrus.Debug(args...)
}
func (Logrus) Warn(ctx context.Context, args ...interface{}) {
	logrus.Debug(args...)
}
func (Logrus) Error(ctx context.Context, args ...interface{}) {
	logrus.Debug(args...)
}
