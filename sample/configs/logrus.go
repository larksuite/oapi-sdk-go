package configs

import (
	"context"
	"github.com/sirupsen/logrus"
	"os"
)

// use logrus implement log.Logger
type Logrus struct {
}

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)
}

func (Logrus) Debug(ctx context.Context, args ...interface{}) {
	logrus.Debug(args...)
}
func (Logrus) Info(ctx context.Context, args ...interface{}) {
	logrus.Info(args...)
}
func (Logrus) Warn(ctx context.Context, args ...interface{}) {
	logrus.Debug(args...)
}
func (Logrus) Error(ctx context.Context, args ...interface{}) {
	logrus.Debug(args...)
}
