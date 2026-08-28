package ws

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type lifecycleRecordingLogger struct {
	mu      sync.Mutex
	entries []string
}

func (l *lifecycleRecordingLogger) Debug(_ context.Context, args ...interface{}) {
	l.record(args...)
}

func (l *lifecycleRecordingLogger) Info(_ context.Context, args ...interface{}) {
	l.record(args...)
}

func (l *lifecycleRecordingLogger) Warn(_ context.Context, args ...interface{}) {
	l.record(args...)
}

func (l *lifecycleRecordingLogger) Error(_ context.Context, args ...interface{}) {
	l.record(args...)
}

func (l *lifecycleRecordingLogger) record(args ...interface{}) {
	l.mu.Lock()
	l.entries = append(l.entries, fmt.Sprint(args...))
	l.mu.Unlock()
}

func (l *lifecycleRecordingLogger) String() string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return strings.Join(l.entries, "\n")
}

func assertLifecycleTextDoesNotContain(t *testing.T, text string, markers ...string) {
	t.Helper()
	for _, marker := range markers {
		if strings.Contains(text, marker) {
			t.Errorf("lifecycle output contains sensitive marker %q: %s", marker, text)
		}
	}
}

func TestConnectedEndpointKeepsHandshakeQueryButRedactsLifecycleOutput(t *testing.T) {
	const (
		appSecretMarker = "app-secret-sensitive-marker"
		pathMarker      = "endpoint-path-sensitive-marker"
		deviceIDMarker  = "device-id-sensitive-marker"
		serviceIDMarker = "424242"
	)

	gateway, cleanupGateway := newLifecycleGateway(t)
	defer cleanupGateway()
	parsedEndpoint, err := url.Parse(gateway.endpoint())
	if err != nil {
		t.Fatalf("parse gateway endpoint: %v", err)
	}
	parsedEndpoint.Path = "/" + pathMarker
	query := parsedEndpoint.Query()
	query.Set(DeviceID, deviceIDMarker)
	query.Set(ServiceID, serviceIDMarker)
	parsedEndpoint.RawQuery = query.Encode()
	endpoint := parsedEndpoint.String()

	var bootstrapRequests int32
	domain, cleanupBootstrap := installLifecycleBootstrap(t, endpoint, &bootstrapRequests)
	defer cleanupBootstrap()
	logger := &lifecycleRecordingLogger{}
	ready := make(chan struct{}, 1)
	client := NewClient("app-id", appSecretMarker,
		WithDomain(domain),
		WithAutoReconnect(false),
		WithLogger(logger),
		WithOnReady(func() { ready <- struct{}{} }),
	)
	result := startLifecycleClient(client, context.Background())
	waitLifecycleSignal(t, ready, "Ready")

	var request lifecycleGatewayRequest
	select {
	case request = <-gateway.requests:
	case <-time.After(lifecycleTestTimeout):
		t.Fatal("timed out waiting for gateway handshake request")
	}
	if request.deviceID != deviceIDMarker {
		t.Errorf("gateway received device_id %q, want %q", request.deviceID, deviceIDMarker)
	}
	if request.serviceID != serviceIDMarker {
		t.Errorf("gateway received service_id %q, want %q", request.serviceID, serviceIDMarker)
	}
	if request.path != "/"+pathMarker {
		t.Errorf("gateway received path %q, want %q", request.path, "/"+pathMarker)
	}

	client.Close()
	startErr, ok := waitLifecycleResult(result)
	if !ok {
		t.Error("Start did not return after Close")
	} else if startErr != nil {
		t.Errorf("Start returned %v after Close, want nil", startErr)
	}
	logs := logger.String()
	if !strings.Contains(logs, "[conn_id="+deviceIDMarker+"]") {
		t.Errorf("lifecycle output did not include conn_id marker: %s", logs)
	}
	assertLifecycleTextDoesNotContain(t, logs,
		appSecretMarker,
		pathMarker,
		serviceIDMarker,
	)
}

func TestHandshakeFailureRedactsServerMessageAndEndpointQuery(t *testing.T) {
	const (
		deviceIDMarker  = "handshake-device-sensitive-marker"
		serviceIDMarker = "515151"
		serverMsgMarker = "handshake-server-sensitive-marker"
	)

	receivedDeviceID := make(chan string, 1)
	gateway := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedDeviceID <- r.URL.Query().Get(DeviceID)
		w.Header().Set(HeaderHandshakeStatus, fmt.Sprint(Forbidden))
		w.Header().Set(HeaderHandshakeMsg, serverMsgMarker)
		w.WriteHeader(http.StatusForbidden)
	}))
	defer gateway.Close()
	endpoint := "ws" + strings.TrimPrefix(gateway.URL, "http") + "/handshake-sensitive-path?" +
		DeviceID + "=" + url.QueryEscape(deviceIDMarker) + "&" + ServiceID + "=" + url.QueryEscape(serviceIDMarker)

	var bootstrapRequests int32
	domain, cleanupBootstrap := installLifecycleBootstrap(t, endpoint, &bootstrapRequests)
	defer cleanupBootstrap()
	logger := &lifecycleRecordingLogger{}
	errorCallback := make(chan error, 1)
	client := NewClient("app-id", "app-secret",
		WithDomain(domain),
		WithAutoReconnect(false),
		WithLogger(logger),
		WithOnError(func(err error) { errorCallback <- err }),
	)

	startErr := client.Start(context.Background())
	if startErr == nil {
		t.Fatal("handshake rejection returned nil")
	}
	var callbackErr error
	select {
	case callbackErr = <-errorCallback:
	case <-time.After(lifecycleTestTimeout):
		t.Fatal("handshake rejection did not invoke OnError")
	}
	var handshakeDeviceID string
	select {
	case handshakeDeviceID = <-receivedDeviceID:
	case <-time.After(lifecycleTestTimeout):
		t.Fatal("timed out waiting for rejected handshake request")
	}
	if handshakeDeviceID != deviceIDMarker {
		t.Errorf("gateway received device_id %q, want %q", handshakeDeviceID, deviceIDMarker)
	}

	var clientErr *ClientError
	if !errors.As(startErr, &clientErr) {
		t.Errorf("handshake rejection returned %T, want *ClientError", startErr)
	}
	combined := startErr.Error() + "\n" + callbackErr.Error() + "\n" + logger.String()
	assertLifecycleTextDoesNotContain(t, combined, deviceIDMarker, serviceIDMarker, serverMsgMarker)
}

func TestMalformedEndpointErrorDoesNotExposeRawCredentials(t *testing.T) {
	const (
		usernameMarker = "endpoint-user-sensitive-marker"
		passwordMarker = "endpoint-password-sensitive-marker"
		deviceIDMarker = "parse-device-sensitive-marker"
		fragmentMarker = "parse-fragment-sensitive-marker"
	)
	rawEndpoint := "ws://" + usernameMarker + ":" + passwordMarker + "@%zz/private?" +
		DeviceID + "=" + deviceIDMarker + "#" + fragmentMarker

	var bootstrapRequests int32
	domain, cleanupBootstrap := installLifecycleBootstrap(t, rawEndpoint, &bootstrapRequests)
	defer cleanupBootstrap()
	logger := &lifecycleRecordingLogger{}
	errorCallback := make(chan error, 1)
	client := NewClient("app-id", "app-secret",
		WithDomain(domain),
		WithAutoReconnect(false),
		WithLogger(logger),
		WithOnError(func(err error) { errorCallback <- err }),
	)

	startErr := client.Start(context.Background())
	if startErr == nil {
		t.Fatal("malformed endpoint returned nil")
	}
	var callbackErr error
	select {
	case callbackErr = <-errorCallback:
	case <-time.After(lifecycleTestTimeout):
		t.Fatal("malformed endpoint did not invoke OnError")
	}
	combined := startErr.Error() + "\n" + callbackErr.Error() + "\n" + logger.String()
	assertLifecycleTextDoesNotContain(t, combined,
		usernameMarker,
		passwordMarker,
		deviceIDMarker,
		fragmentMarker,
	)
	if got := atomic.LoadInt32(&bootstrapRequests); got != 1 {
		t.Errorf("malformed endpoint made %d bootstrap requests, want 1", got)
	}
}
