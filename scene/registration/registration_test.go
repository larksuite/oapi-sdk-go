package registration

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
	"time"
)

func TestRegisterAppRequiresOnQRCode(t *testing.T) {
	_, err := RegisterApp(context.Background(), &Options{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "registration: OnQRCode is required" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRegisterAppSuccess(t *testing.T) {
	restoreWait := stubWaitForInterval()
	defer restoreWait()

	pollCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Fatalf("parse form failed: %v", err)
		}
		switch r.Form.Get("action") {
		case "begin":
			writeJSON(w, `{"device_code":"device-1","verification_uri_complete":"https://qr.example.com/scan?foo=bar","interval":0,"expire_in":60}`)
		case "poll":
			pollCount++
			if pollCount == 1 {
				writeJSON(w, `{"error":"authorization_pending","error_description":"pending"}`)
				return
			}
			writeJSON(w, `{"client_id":"cli_a","client_secret":"sec_a","user_info":{"open_id":"ou_x","tenant_brand":"feishu"}}`)
		default:
			t.Fatalf("unexpected action: %s", r.Form.Get("action"))
		}
	}))
	defer server.Close()

	var qrInfo *QRCodeInfo
	statuses := make([]string, 0, 1)
	result, err := RegisterApp(context.Background(), &Options{
		Domain: server.URL,
		Source: "my-cli",
		OnQRCode: func(info *QRCodeInfo) {
			qrInfo = info
		},
		OnStatusChange: func(info *StatusChangeInfo) {
			statuses = append(statuses, info.Status)
		},
	})
	if err != nil {
		t.Fatalf("register app failed: %v", err)
	}
	if result.ClientID != "cli_a" || result.ClientSecret != "sec_a" {
		t.Fatalf("unexpected result: %#v", result)
	}
	if result.UserInfo == nil || result.UserInfo.OpenID != "ou_x" || result.UserInfo.TenantBrand != "feishu" {
		t.Fatalf("unexpected user info: %#v", result.UserInfo)
	}
	if qrInfo == nil {
		t.Fatal("expected qr info")
	}
	if qrInfo.ExpireIn != 60 {
		t.Fatalf("unexpected expire_in: %d", qrInfo.ExpireIn)
	}
	parsedQR, err := url.Parse(qrInfo.URL)
	if err != nil {
		t.Fatalf("parse qr url failed: %v", err)
	}
	query := parsedQR.Query()
	if query.Get("foo") != "bar" {
		t.Fatalf("unexpected original query: %s", query.Encode())
	}
	if query.Get("from") != "sdk" || query.Get("tp") != "sdk" || query.Get("source") != "go-sdk/my-cli" {
		t.Fatalf("unexpected qr query: %s", query.Encode())
	}
	if len(statuses) != 1 || statuses[0] != StatusPolling {
		t.Fatalf("unexpected statuses: %#v", statuses)
	}
}

func TestRegisterAppMockFlowForLocalRun(t *testing.T) {
	restoreWait := stubWaitForInterval()
	defer restoreWait()

	beginCalled := 0
	pollCalled := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Fatalf("parse form failed: %v", err)
		}

		action := r.Form.Get("action")
		switch action {
		case "begin":
			beginCalled++
			writeJSON(w, `{
				"device_code":"mock-device-code",
				"verification_uri_complete":"https://mock-qr.example.com/activate?ticket=abc123",
				"interval":0,
				"expire_in":120
			}`)
		case "poll":
			pollCalled++
			if r.Form.Get("device_code") != "mock-device-code" {
				t.Fatalf("unexpected device_code: %s", r.Form.Get("device_code"))
			}
			if pollCalled == 1 {
				writeJSON(w, `{
					"error":"authorization_pending",
					"error_description":"waiting for user scan"
				}`)
				return
			}
			writeJSON(w, `{
				"client_id":"cli_mock_local",
				"client_secret":"sec_mock_local",
				"user_info":{
					"open_id":"ou_mock_user",
					"tenant_brand":"feishu"
				}
			}`)
		default:
			t.Fatalf("unexpected action: %s", action)
		}
	}))
	defer server.Close()

	var qrInfo *QRCodeInfo
	statuses := make([]string, 0, 2)
	result, err := RegisterApp(context.Background(), &Options{
		Domain: server.URL,
		Source: "local-test",
		OnQRCode: func(info *QRCodeInfo) {
			qrInfo = info
			t.Logf("mock qr generated: url=%s expire_in=%d", info.URL, info.ExpireIn)
		},
		OnStatusChange: func(info *StatusChangeInfo) {
			statuses = append(statuses, info.Status)
			t.Logf("mock status changed: status=%s interval=%d", info.Status, info.Interval)
		},
	})
	if err != nil {
		t.Fatalf("register app failed: %v", err)
	}

	if beginCalled != 1 {
		t.Fatalf("unexpected begin count: %d", beginCalled)
	}
	if pollCalled != 2 {
		t.Fatalf("unexpected poll count: %d", pollCalled)
	}
	if qrInfo == nil {
		t.Fatal("expected qr info")
	}

	parsedQR, err := url.Parse(qrInfo.URL)
	if err != nil {
		t.Fatalf("parse qr url failed: %v", err)
	}
	query := parsedQR.Query()
	if query.Get("ticket") != "abc123" {
		t.Fatalf("unexpected qr query: %s", query.Encode())
	}
	if query.Get("from") != "sdk" || query.Get("source") != "go-sdk/local-test" || query.Get("tp") != "sdk" {
		t.Fatalf("unexpected qr query params: %s", query.Encode())
	}

	if result.ClientID != "cli_mock_local" || result.ClientSecret != "sec_mock_local" {
		t.Fatalf("unexpected result: %#v", result)
	}
	if result.UserInfo == nil || result.UserInfo.OpenID != "ou_mock_user" {
		t.Fatalf("unexpected user info: %#v", result.UserInfo)
	}
	if len(statuses) != 1 || statuses[0] != StatusPolling {
		t.Fatalf("unexpected statuses: %#v", statuses)
	}

	t.Logf("mock registration finished: client_id=%s poll_count=%d", result.ClientID, pollCalled)
}

func TestRegisterAppPollsImmediatelyThenUsesReturnedInterval(t *testing.T) {
	getIntervals, restoreWait := stubWaitRecorder()
	defer restoreWait()

	pollCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Fatalf("parse form failed: %v", err)
		}
		switch r.Form.Get("action") {
		case "begin":
			writeJSON(w, `{"device_code":"device-immediate","verification_uri_complete":"https://qr.example.com/scan","interval":7,"expire_in":60}`)
		case "poll":
			pollCount++
			if pollCount == 1 {
				writeJSON(w, `{"error":"authorization_pending","error_description":"pending"}`)
				return
			}
			writeJSON(w, `{"client_id":"cli_immediate","client_secret":"sec_immediate"}`)
		default:
			t.Fatalf("unexpected action: %s", r.Form.Get("action"))
		}
	}))
	defer server.Close()

	result, err := RegisterApp(context.Background(), &Options{
		Domain:   server.URL,
		OnQRCode: func(info *QRCodeInfo) {},
	})
	if err != nil {
		t.Fatalf("register app failed: %v", err)
	}
	if result.ClientID != "cli_immediate" {
		t.Fatalf("unexpected result: %#v", result)
	}
	if pollCount != 2 {
		t.Fatalf("unexpected poll count: %d", pollCount)
	}
	intervals := getIntervals()
	if !reflect.DeepEqual(intervals, []time.Duration{7 * time.Second}) {
		t.Fatalf("unexpected wait intervals: %#v", intervals)
	}
}

func TestRegisterAppUsesDefaultIntervalAndExpireIn(t *testing.T) {
	getIntervals, restoreWait := stubWaitRecorder()
	defer restoreWait()

	pollCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Fatalf("parse form failed: %v", err)
		}
		switch r.Form.Get("action") {
		case "begin":
			writeJSON(w, `{"device_code":"device-default","verification_uri_complete":"https://qr.example.com/scan","interval":0,"expire_in":0}`)
		case "poll":
			pollCount++
			writeJSON(w, `{"client_id":"cli_default","client_secret":"sec_default"}`)
		default:
			t.Fatalf("unexpected action: %s", r.Form.Get("action"))
		}
	}))
	defer server.Close()

	var qrInfo *QRCodeInfo
	result, err := RegisterApp(context.Background(), &Options{
		Domain: server.URL,
		OnQRCode: func(info *QRCodeInfo) {
			qrInfo = info
		},
	})
	if err != nil {
		t.Fatalf("register app failed: %v", err)
	}
	if result.ClientID != "cli_default" {
		t.Fatalf("unexpected result: %#v", result)
	}
	if qrInfo == nil || qrInfo.ExpireIn != defaultExpireInSeconds {
		t.Fatalf("unexpected qr info: %#v", qrInfo)
	}
	intervals := getIntervals()
	if len(intervals) != 0 {
		t.Fatalf("unexpected wait intervals: %#v", intervals)
	}
}

func TestRegisterAppKeepsPollingOnEmptyPollResponse(t *testing.T) {
	getIntervals, restoreWait := stubWaitRecorder()
	defer restoreWait()

	pollCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Fatalf("parse form failed: %v", err)
		}
		switch r.Form.Get("action") {
		case "begin":
			writeJSON(w, `{"device_code":"device-empty","verification_uri_complete":"https://qr.example.com/scan","interval":3,"expire_in":60}`)
		case "poll":
			pollCount++
			if pollCount == 1 {
				writeJSON(w, `{}`)
				return
			}
			writeJSON(w, `{"client_id":"cli_after_empty","client_secret":"sec_after_empty"}`)
		default:
			t.Fatalf("unexpected action: %s", r.Form.Get("action"))
		}
	}))
	defer server.Close()

	result, err := RegisterApp(context.Background(), &Options{
		Domain:   server.URL,
		OnQRCode: func(info *QRCodeInfo) {},
	})
	if err != nil {
		t.Fatalf("register app failed: %v", err)
	}
	if result.ClientID != "cli_after_empty" {
		t.Fatalf("unexpected result: %#v", result)
	}
	if pollCount != 2 {
		t.Fatalf("unexpected poll count: %d", pollCount)
	}
	intervals := getIntervals()
	if !reflect.DeepEqual(intervals, []time.Duration{3 * time.Second}) {
		t.Fatalf("unexpected wait intervals: %#v", intervals)
	}
}

func TestRegisterAppSwitchesToLarkDomain(t *testing.T) {
	restoreWait := stubWaitForInterval()
	defer restoreWait()

	larkPollCount := 0
	larkServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Fatalf("parse form failed: %v", err)
		}
		larkPollCount++
		writeJSON(w, `{"client_id":"cli_lark","client_secret":"sec_lark","user_info":{"open_id":"ou_lark","tenant_brand":"lark"}}`)
	}))
	defer larkServer.Close()

	feishuPollCount := 0
	feishuServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Fatalf("parse form failed: %v", err)
		}
		switch r.Form.Get("action") {
		case "begin":
			writeJSON(w, `{"device_code":"device-2","verification_uri_complete":"https://qr.example.com/scan","interval":0,"expire_in":60}`)
		case "poll":
			feishuPollCount++
			writeJSON(w, `{"user_info":{"tenant_brand":"lark"}}`)
		default:
			t.Fatalf("unexpected action: %s", r.Form.Get("action"))
		}
	}))
	defer feishuServer.Close()

	statuses := make([]string, 0, 1)
	result, err := RegisterApp(context.Background(), &Options{
		Domain:     feishuServer.URL,
		LarkDomain: larkServer.URL,
		OnQRCode:   func(info *QRCodeInfo) {},
		OnStatusChange: func(info *StatusChangeInfo) {
			statuses = append(statuses, info.Status)
		},
	})
	if err != nil {
		t.Fatalf("register app failed: %v", err)
	}
	if result.ClientID != "cli_lark" || result.ClientSecret != "sec_lark" {
		t.Fatalf("unexpected result: %#v", result)
	}
	if feishuPollCount != 1 {
		t.Fatalf("unexpected feishu poll count: %d", feishuPollCount)
	}
	if larkPollCount != 1 {
		t.Fatalf("unexpected lark poll count: %d", larkPollCount)
	}
	if len(statuses) != 1 || statuses[0] != StatusDomainSwitched {
		t.Fatalf("unexpected statuses: %#v", statuses)
	}
}

func TestRegisterAppSlowDown(t *testing.T) {
	restoreWait := stubWaitForInterval()
	defer restoreWait()

	pollCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Fatalf("parse form failed: %v", err)
		}
		switch r.Form.Get("action") {
		case "begin":
			writeJSON(w, `{"device_code":"device-3","verification_uri_complete":"https://qr.example.com/scan","interval":0,"expire_in":60}`)
		case "poll":
			pollCount++
			if pollCount == 1 {
				writeJSON(w, `{"error":"slow_down","error_description":"slow down"}`)
				return
			}
			writeJSON(w, `{"client_id":"cli_b","client_secret":"sec_b"}`)
		default:
			t.Fatalf("unexpected action: %s", r.Form.Get("action"))
		}
	}))
	defer server.Close()

	var slowDownInfo *StatusChangeInfo
	result, err := RegisterApp(context.Background(), &Options{
		Domain:   server.URL,
		OnQRCode: func(info *QRCodeInfo) {},
		OnStatusChange: func(info *StatusChangeInfo) {
			if info.Status == StatusSlowDown {
				slowDownInfo = info
			}
		},
	})
	if err != nil {
		t.Fatalf("register app failed: %v", err)
	}
	if result.ClientID != "cli_b" || result.ClientSecret != "sec_b" {
		t.Fatalf("unexpected result: %#v", result)
	}
	if slowDownInfo == nil {
		t.Fatal("expected slow_down status")
	}
	if slowDownInfo.Interval != 10 {
		t.Fatalf("unexpected interval: %d", slowDownInfo.Interval)
	}
}

func TestRegisterAppAccessDenied(t *testing.T) {
	restoreWait := stubWaitForInterval()
	defer restoreWait()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Fatalf("parse form failed: %v", err)
		}
		switch r.Form.Get("action") {
		case "begin":
			writeJSON(w, `{"device_code":"device-4","verification_uri_complete":"https://qr.example.com/scan","interval":0,"expire_in":60}`)
		case "poll":
			writeJSON(w, `{"error":"access_denied","error_description":"user denied"}`)
		default:
			t.Fatalf("unexpected action: %s", r.Form.Get("action"))
		}
	}))
	defer server.Close()

	_, err := RegisterApp(context.Background(), &Options{
		Domain:   server.URL,
		OnQRCode: func(info *QRCodeInfo) {},
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var accessDenied *AccessDeniedError
	if !errors.As(err, &accessDenied) {
		t.Fatalf("unexpected error type: %#v", err)
	}
	if accessDenied.Code != "access_denied" || accessDenied.Description != "user denied" {
		t.Fatalf("unexpected error: %#v", accessDenied)
	}
}

func TestRegisterAppExpiredToken(t *testing.T) {
	restoreWait := stubWaitForInterval()
	defer restoreWait()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Fatalf("parse form failed: %v", err)
		}
		switch r.Form.Get("action") {
		case "begin":
			writeJSON(w, `{"device_code":"device-5","verification_uri_complete":"https://qr.example.com/scan","interval":0,"expire_in":60}`)
		case "poll":
			writeJSON(w, `{"error":"expired_token","error_description":"qr expired"}`)
		default:
			t.Fatalf("unexpected action: %s", r.Form.Get("action"))
		}
	}))
	defer server.Close()

	_, err := RegisterApp(context.Background(), &Options{
		Domain:   server.URL,
		OnQRCode: func(info *QRCodeInfo) {},
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var expired *ExpiredError
	if !errors.As(err, &expired) {
		t.Fatalf("unexpected error type: %#v", err)
	}
	if expired.Code != "expired_token" || expired.Description != "qr expired" {
		t.Fatalf("unexpected error: %#v", expired)
	}
}

func stubWaitForInterval() func() {
	original := waitForInterval
	waitForInterval = func(ctx context.Context, intervalDuration time.Duration) error {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		return nil
	}
	return func() {
		waitForInterval = original
	}
}

func stubWaitRecorder() (func() []time.Duration, func()) {
	original := waitForInterval
	intervals := make([]time.Duration, 0, 4)
	waitForInterval = func(ctx context.Context, intervalDuration time.Duration) error {
		intervals = append(intervals, intervalDuration)
		if ctx.Err() != nil {
			return ctx.Err()
		}
		return nil
	}
	return func() []time.Duration {
			copied := make([]time.Duration, len(intervals))
			copy(copied, intervals)
			return copied
		}, func() {
			waitForInterval = original
		}
}

func writeJSON(w http.ResponseWriter, body string) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(body))
}
