package dispatcher

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkevent "github.com/larksuite/oapi-sdk-go/v3/event"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

type constructorLog struct{ level, message string }
type constructorLogger struct{ events []constructorLog }

func (l *constructorLogger) Debug(_ context.Context, args ...interface{}) {
	l.events = append(l.events, constructorLog{"debug", fmt.Sprint(args...)})
}
func (l *constructorLogger) Info(_ context.Context, args ...interface{}) {
	l.events = append(l.events, constructorLog{"info", fmt.Sprint(args...)})
}
func (l *constructorLogger) Warn(_ context.Context, args ...interface{}) {
	l.events = append(l.events, constructorLog{"warn", fmt.Sprint(args...)})
}
func (l *constructorLogger) Error(_ context.Context, args ...interface{}) {
	l.events = append(l.events, constructorLog{"error", fmt.Sprint(args...)})
}

func TestDispatcherConstructorLogger(t *testing.T) {
	for _, test := range []struct {
		name  string
		level larkcore.LogLevel
		want  []constructorLog
	}{
		{"default", 0, []constructorLog{{"info", "event-dispatch is ready"}, {"info", "info"}, {"warn", "warn"}, {"error", "error"}}},
		{"debug", larkcore.LogLevelDebug, []constructorLog{{"info", "event-dispatch is ready"}, {"debug", "debug"}, {"info", "info"}, {"warn", "warn"}, {"error", "error"}}},
		{"warn", larkcore.LogLevelWarn, []constructorLog{{"info", "event-dispatch is ready"}, {"warn", "warn"}, {"error", "error"}}},
	} {
		t.Run(test.name, func(t *testing.T) {
			logger := &constructorLogger{}
			dispatcher := NewEventDispatcherWithLogger("", "", logger)
			// Constructor injection must not pre-wrap the logger and prevent the
			// existing InitConfig entry point from selecting the requested level.
			dispatcher.InitConfig(larkevent.WithLogLevel(test.level))
			dispatcher.Logger().Debug(context.Background(), "debug")
			dispatcher.Logger().Info(context.Background(), "info")
			dispatcher.Logger().Warn(context.Background(), "warn")
			dispatcher.Logger().Error(context.Background(), "error")
			if !reflect.DeepEqual(logger.events, test.want) {
				t.Fatalf("dispatcher logs = %#v, want %#v", logger.events, test.want)
			}
			var messageID string
			dispatcher.OnP2MessageReceiveV1(func(_ context.Context, event *larkim.P2MessageReceiveV1) error {
				messageID = *event.Event.Message.MessageId
				return nil
			})
			_, err := dispatcher.Do(context.Background(), []byte(`{"schema":"2.0","header":{"event_id":"event-1","event_type":"im.message.receive_v1"},"event":{"message":{"message_id":"message-1"}}}`))
			if err != nil || messageID != "message-1" {
				t.Fatalf("message dispatch = %q, %v", messageID, err)
			}
		})
	}
}
