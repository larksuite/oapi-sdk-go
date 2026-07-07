package registration

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
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

func TestRegisterAppOmitsAppPresetWhenNil(t *testing.T) {
	qrURL := captureQRURL(t, Options{})
	query := qrURL.Query()

	if _, ok := query["avatar"]; ok {
		t.Fatalf("unexpected avatar param: %s", qrURL.String())
	}
	if _, ok := query["name"]; ok {
		t.Fatalf("unexpected name param: %s", qrURL.String())
	}
	if _, ok := query["desc"]; ok {
		t.Fatalf("unexpected desc param: %s", qrURL.String())
	}
	if query.Get("from") != "sdk" || query.Get("source") != "go-sdk" || query.Get("tp") != "sdk" {
		t.Fatalf("unexpected qr query: %s", query.Encode())
	}
}

func TestRegisterAppAppPresetSingleAvatar(t *testing.T) {
	qrURL := captureQRURL(t, Options{
		AppPreset: &AppPreset{
			Avatar: []string{"https://example.com/a.png"},
		},
	})

	if !reflect.DeepEqual(qrURL.Query()["avatar"], []string{"https://example.com/a.png"}) {
		t.Fatalf("unexpected avatar params: %#v", qrURL.Query()["avatar"])
	}
}

func TestRegisterAppAppPresetMultipleAvatarsPreserveOrder(t *testing.T) {
	avatars := []string{
		"https://example.com/a.png",
		"https://example.com/b.webp",
		"https://example.com/c.gif",
	}
	qrURL := captureQRURL(t, Options{
		AppPreset: &AppPreset{
			Avatar: avatars,
		},
	})

	if !reflect.DeepEqual(qrURL.Query()["avatar"], avatars) {
		t.Fatalf("unexpected avatar params: %#v", qrURL.Query()["avatar"])
	}
}

func TestRegisterAppAppPresetAcceptsExactlySixAvatars(t *testing.T) {
	avatars := []string{
		"https://example.com/0.png",
		"https://example.com/1.png",
		"https://example.com/2.png",
		"https://example.com/3.png",
		"https://example.com/4.png",
		"https://example.com/5.png",
	}
	qrURL := captureQRURL(t, Options{
		AppPreset: &AppPreset{
			Avatar: avatars,
		},
	})

	if !reflect.DeepEqual(qrURL.Query()["avatar"], avatars) {
		t.Fatalf("unexpected avatar params: %#v", qrURL.Query()["avatar"])
	}
}

func TestRegisterAppAppPresetRejectsMoreThanSixAvatars(t *testing.T) {
	expectAppPresetReject(t, &AppPreset{
		Avatar: []string{
			"https://example.com/0.png",
			"https://example.com/1.png",
			"https://example.com/2.png",
			"https://example.com/3.png",
			"https://example.com/4.png",
			"https://example.com/5.png",
			"https://example.com/6.png",
		},
	}, "AppPreset.Avatar supports at most 6 URLs, got 7")
}

func TestRegisterAppAppPresetRejectsEmptyAvatarEntry(t *testing.T) {
	expectAppPresetReject(t, &AppPreset{
		Avatar: []string{"https://example.com/a.png", ""},
	}, "AppPreset.Avatar[1] must be a non-empty string")
}

func TestRegisterAppAppPresetEncodesNameWithUserPlaceholder(t *testing.T) {
	qrURL := captureQRURL(t, Options{
		AppPreset: &AppPreset{
			Name: "{user}的应用",
		},
	})

	if got := qrURL.Query().Get("name"); got != "{user}的应用" {
		t.Fatalf("unexpected name: %s", got)
	}
	if !strings.Contains(qrURL.String(), "name=%7Buser%7D%E7%9A%84%E5%BA%94%E7%94%A8") {
		t.Fatalf("expected encoded name in raw qr url: %s", qrURL.String())
	}
}

func TestRegisterAppAppPresetEncodesDesc(t *testing.T) {
	qrURL := captureQRURL(t, Options{
		AppPreset: &AppPreset{
			Desc: "由业务平台自动生成",
		},
	})

	if got := qrURL.Query().Get("desc"); got != "由业务平台自动生成" {
		t.Fatalf("unexpected desc: %s", got)
	}
}

func TestRegisterAppAppPresetEmitsAllFieldsTogether(t *testing.T) {
	qrURL := captureQRURL(t, Options{
		AppPreset: &AppPreset{
			Avatar: []string{"https://example.com/a.png", "https://example.com/b.png"},
			Name:   "MyApp",
			Desc:   "demo",
		},
	})
	query := qrURL.Query()

	if !reflect.DeepEqual(query["avatar"], []string{"https://example.com/a.png", "https://example.com/b.png"}) {
		t.Fatalf("unexpected avatar params: %#v", query["avatar"])
	}
	if query.Get("name") != "MyApp" {
		t.Fatalf("unexpected name: %s", query.Get("name"))
	}
	if query.Get("desc") != "demo" {
		t.Fatalf("unexpected desc: %s", query.Get("desc"))
	}
}

func TestRegisterAppAppPresetDoesNotInterfereWithSource(t *testing.T) {
	qrURL := captureQRURL(t, Options{
		Source: "lark-cli",
		AppPreset: &AppPreset{
			Name: "X",
		},
	})
	query := qrURL.Query()

	if query.Get("source") != "go-sdk/lark-cli" {
		t.Fatalf("unexpected source: %s", query.Get("source"))
	}
	if query.Get("name") != "X" {
		t.Fatalf("unexpected name: %s", query.Get("name"))
	}
}

func TestRegisterAppOmitsAddonsCreateOnlyAndAppIDWhenUnset(t *testing.T) {
	qrURL := captureQRURL(t, Options{})
	query := qrURL.Query()

	for _, key := range []string{"addons", "createOnly", "clientID"} {
		if _, ok := query[key]; ok {
			t.Fatalf("unexpected %s param: %s", key, qrURL.String())
		}
	}
}

func TestRegisterAppAddonsEncodesURLSafeParam(t *testing.T) {
	addons := &AppAddons{
		Scopes: AppAddonsScopes{
			Tenant: []string{"im:message:send_as_bot", "drive:drive.metadata:readonly"},
			User:   []string{"calendar:calendar:read"},
		},
		Events: AppAddonsEvents{
			Items: AppAddonsEventItems{
				Tenant: []string{"im.message.receive_v1"},
				User:   []string{"calendar.calendar.event.changed_v4"},
			},
		},
		Callbacks: AppAddonsCallbacks{
			Items: []string{"card.action.trigger"},
		},
	}
	qrURL := captureQRURL(t, Options{
		Addons: addons,
	})
	encoded := qrURL.Query().Get("addons")

	if encoded == "" {
		t.Fatalf("expected addons param: %s", qrURL.String())
	}
	if strings.ContainsAny(encoded, "+/=") {
		t.Fatalf("expected URL-safe base64 without padding, got %q", encoded)
	}
	decoded := decodeAddonsParam(t, qrURL)
	if !reflect.DeepEqual(decoded, *addons) {
		t.Fatalf("unexpected decoded addons:\nwant: %#v\n got: %#v", *addons, decoded)
	}
}

func TestRegisterAppAddonsRejectsEmpty(t *testing.T) {
	expectOptionsReject(t, Options{
		Addons: &AppAddons{},
	}, "Addons must contain at least one scope, event or callback")
}

func TestRegisterAppAddonsRejectsEmptyItem(t *testing.T) {
	expectOptionsReject(t, Options{
		Addons: &AppAddons{
			Callbacks: AppAddonsCallbacks{
				Items: []string{"card.action.trigger", ""},
			},
		},
	}, "Addons.Callbacks.Items[1] must be a non-empty string")
}

func TestRegisterAppAddonsOmitsPresetKeyWhenUnset(t *testing.T) {
	addons := &AppAddons{
		Scopes: AppAddonsScopes{
			Tenant: []string{"im:message:send_as_bot"},
			User:   []string{"calendar:calendar:read"},
		},
		Events: AppAddonsEvents{
			Items: AppAddonsEventItems{
				Tenant: []string{"im.message.receive_v1"},
			},
		},
		Callbacks: AppAddonsCallbacks{
			Items: []string{"card.action.trigger"},
		},
	}
	qrURL := captureQRURL(t, Options{
		Addons: addons,
	})

	payload := decodeAddonsParamRaw(t, qrURL)
	if _, ok := payload["preset"]; ok {
		t.Fatalf("unexpected preset key in payload: %#v", payload)
	}
	decoded := decodeAddonsParam(t, qrURL)
	if !reflect.DeepEqual(decoded, *addons) {
		t.Fatalf("unexpected decoded addons:\nwant: %#v\n got: %#v", *addons, decoded)
	}
}

func TestRegisterAppAddonsPresetFalseAllowsEmptyIncrements(t *testing.T) {
	qrURL := captureQRURL(t, Options{
		Addons: &AppAddons{
			Preset: boolPtr(false),
		},
	})

	payload := decodeAddonsParamRaw(t, qrURL)
	// {"preset":false} alone is a valid payload: the confirmation page renders
	// a near-empty app on the minimal template base.
	want := map[string]interface{}{
		"preset": false,
	}
	if !reflect.DeepEqual(payload, want) {
		t.Fatalf("unexpected payload:\nwant: %#v\n got: %#v", want, payload)
	}
}

func TestRegisterAppAddonsPresetFalseCoexistsWithIncrements(t *testing.T) {
	qrURL := captureQRURL(t, Options{
		Addons: &AppAddons{
			Preset: boolPtr(false),
			Scopes: AppAddonsScopes{
				Tenant: []string{"im:message:send_as_bot"},
			},
		},
	})

	payload := decodeAddonsParamRaw(t, qrURL)
	want := map[string]interface{}{
		"preset": false,
		"scopes": map[string]interface{}{
			"tenant": []interface{}{"im:message:send_as_bot"},
		},
	}
	if !reflect.DeepEqual(payload, want) {
		t.Fatalf("unexpected payload:\nwant: %#v\n got: %#v", want, payload)
	}
}

func TestRegisterAppAddonsPresetTrueEncodesPresetKey(t *testing.T) {
	qrURL := captureQRURL(t, Options{
		Addons: &AppAddons{
			Preset: boolPtr(true),
			Callbacks: AppAddonsCallbacks{
				Items: []string{"card.action.trigger"},
			},
		},
	})

	payload := decodeAddonsParamRaw(t, qrURL)
	want := map[string]interface{}{
		"preset": true,
		"callbacks": map[string]interface{}{
			"items": []interface{}{"card.action.trigger"},
		},
	}
	if !reflect.DeepEqual(payload, want) {
		t.Fatalf("unexpected payload:\nwant: %#v\n got: %#v", want, payload)
	}
}

func TestRegisterAppAddonsPresetTrueRejectsEmptyIncrements(t *testing.T) {
	// Explicit true selects the same default template base as leaving Preset
	// unset, so an otherwise empty addons payload stays invalid. The spec only
	// fixes the error prefix; the exact wording is up to the implementation.
	expectOptionsReject(t, Options{
		Addons: &AppAddons{
			Preset: boolPtr(true),
		},
	}, "registration:")
}

func TestRegisterAppAddonsPresetFalseKeepsNonEmptyStringValidation(t *testing.T) {
	expectOptionsReject(t, Options{
		Addons: &AppAddons{
			Preset: boolPtr(false),
			Callbacks: AppAddonsCallbacks{
				Items: []string{"card.action.trigger", ""},
			},
		},
	}, "Addons.Callbacks.Items[1] must be a non-empty string")
}

func TestRegisterAppAddonsPresetFalseFullFlowQRCodeURL(t *testing.T) {
	restoreWait := stubWaitForInterval()
	defer restoreWait()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Fatalf("parse form failed: %v", err)
		}
		switch r.Form.Get("action") {
		case "begin":
			writeJSON(w, `{"device_code":"device-preset","verification_uri_complete":"https://qr.example.com/scan?ticket=preset-1","interval":0,"expire_in":60}`)
		case "poll":
			writeJSON(w, `{"client_id":"cli_preset","client_secret":"sec_preset"}`)
		default:
			t.Fatalf("unexpected action: %s", r.Form.Get("action"))
		}
	}))
	defer server.Close()

	var qrInfo *QRCodeInfo
	result, err := RegisterApp(context.Background(), &Options{
		Domain: server.URL,
		Source: "preset-cli",
		Addons: &AppAddons{
			Preset: boolPtr(false),
			Scopes: AppAddonsScopes{
				Tenant: []string{"im:message:send_as_bot"},
			},
		},
		OnQRCode: func(info *QRCodeInfo) {
			qrInfo = info
		},
	})
	if err != nil {
		t.Fatalf("register app failed: %v", err)
	}
	if result.ClientID != "cli_preset" || result.ClientSecret != "sec_preset" {
		t.Fatalf("unexpected result: %#v", result)
	}
	if qrInfo == nil {
		t.Fatal("expected qr info")
	}

	qrURL, err := url.Parse(qrInfo.URL)
	if err != nil {
		t.Fatalf("parse qr url failed: %v", err)
	}
	query := qrURL.Query()
	if query.Get("ticket") != "preset-1" {
		t.Fatalf("unexpected original query: %s", query.Encode())
	}
	if query.Get("from") != "sdk" || query.Get("tp") != "sdk" || query.Get("source") != "go-sdk/preset-cli" {
		t.Fatalf("unexpected qr query: %s", query.Encode())
	}
	if encoded := query.Get("addons"); strings.ContainsAny(encoded, "+/=") {
		t.Fatalf("expected URL-safe base64 without padding, got %q", encoded)
	}
	payload := decodeAddonsParamRaw(t, qrURL)
	want := map[string]interface{}{
		"preset": false,
		"scopes": map[string]interface{}{
			"tenant": []interface{}{"im:message:send_as_bot"},
		},
	}
	if !reflect.DeepEqual(payload, want) {
		t.Fatalf("unexpected payload:\nwant: %#v\n got: %#v", want, payload)
	}
}

func TestRegisterAppCreateOnlySetsParamOnlyWhenTrue(t *testing.T) {
	enabledURL := captureQRURL(t, Options{
		CreateOnly: true,
	})
	if got := enabledURL.Query().Get("createOnly"); got != "true" {
		t.Fatalf("unexpected createOnly: %s", got)
	}

	disabledURL := captureQRURL(t, Options{
		CreateOnly: false,
	})
	if _, ok := disabledURL.Query()["createOnly"]; ok {
		t.Fatalf("unexpected createOnly param: %s", disabledURL.String())
	}
}

func TestRegisterAppAppIDSetsClientID(t *testing.T) {
	qrURL := captureQRURL(t, Options{
		AppID: "cli_a1b2c3",
	})

	if got := qrURL.Query().Get("clientID"); got != "cli_a1b2c3" {
		t.Fatalf("unexpected clientID: %s", got)
	}
}

func TestRegisterAppNewParamsCoexistWithAppPreset(t *testing.T) {
	qrURL := captureQRURL(t, Options{
		AppPreset: &AppPreset{
			Name: "MyApp",
		},
		Addons: &AppAddons{
			Scopes: AppAddonsScopes{
				Tenant: []string{"im:message:send_as_bot"},
			},
		},
		CreateOnly: true,
		AppID:      "cli_a1b2c3",
	})
	query := qrURL.Query()

	if got := query.Get("name"); got != "MyApp" {
		t.Fatalf("unexpected name: %s", got)
	}
	if got := query.Get("createOnly"); got != "true" {
		t.Fatalf("unexpected createOnly: %s", got)
	}
	if got := query.Get("clientID"); got != "cli_a1b2c3" {
		t.Fatalf("unexpected clientID: %s", got)
	}
	decoded := decodeAddonsParam(t, qrURL)
	want := AppAddons{
		Scopes: AppAddonsScopes{
			Tenant: []string{"im:message:send_as_bot"},
		},
	}
	if !reflect.DeepEqual(decoded, want) {
		t.Fatalf("unexpected decoded addons:\nwant: %#v\n got: %#v", want, decoded)
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

func captureQRURL(t *testing.T, opts Options) *url.URL {
	t.Helper()

	restoreWait := stubWaitForInterval()
	defer restoreWait()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Fatalf("parse form failed: %v", err)
		}
		switch r.Form.Get("action") {
		case "begin":
			writeJSON(w, `{"device_code":"dev-1","verification_uri_complete":"https://accounts.feishu.cn/page/launcher","verification_uri":"https://accounts.feishu.cn/page/launcher","user_code":"uc-1","interval":0,"expire_in":600}`)
		case "poll":
			writeJSON(w, `{"client_id":"cli_test","client_secret":"sec_test"}`)
		default:
			t.Fatalf("unexpected action: %s", r.Form.Get("action"))
		}
	}))
	defer server.Close()

	var captured string
	opts.Domain = server.URL
	if opts.OnQRCode == nil {
		opts.OnQRCode = func(info *QRCodeInfo) {
			captured = info.URL
		}
	}

	if _, err := RegisterApp(context.Background(), &opts); err != nil {
		t.Fatalf("register app failed: %v", err)
	}
	if captured == "" {
		t.Fatal("expected qr url")
	}
	parsed, err := url.Parse(captured)
	if err != nil {
		t.Fatalf("parse qr url failed: %v", err)
	}
	return parsed
}

func expectAppPresetReject(t *testing.T, preset *AppPreset, wantErr string) {
	t.Helper()

	expectOptionsReject(t, Options{
		AppPreset: preset,
	}, wantErr)
}

func expectOptionsReject(t *testing.T, opts Options, wantErr string) {
	t.Helper()

	restoreWait := stubWaitForInterval()
	defer restoreWait()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Fatalf("parse form failed: %v", err)
		}
		switch r.Form.Get("action") {
		case "begin":
			writeJSON(w, `{"device_code":"dev-1","verification_uri_complete":"https://accounts.feishu.cn/page/launcher","interval":0,"expire_in":600}`)
		case "poll":
			t.Fatal("poll should not be called when options validation fails")
		default:
			t.Fatalf("unexpected action: %s", r.Form.Get("action"))
		}
	}))
	defer server.Close()

	calledQRCode := false
	opts.Domain = server.URL
	opts.OnQRCode = func(info *QRCodeInfo) {
		calledQRCode = true
	}
	_, err := RegisterApp(context.Background(), &opts)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), wantErr) {
		t.Fatalf("unexpected error: %v", err)
	}
	if calledQRCode {
		t.Fatal("OnQRCode should not be called")
	}
}

func decodeAddonsParam(t *testing.T, qrURL *url.URL) AppAddons {
	t.Helper()

	var addons AppAddons
	if err := json.Unmarshal(decodeAddonsParamBody(t, qrURL), &addons); err != nil {
		t.Fatalf("unmarshal addons failed: %v", err)
	}
	return addons
}

// decodeAddonsParamRaw decodes the addons payload into a generic map because a
// typed AppAddons cannot distinguish an absent key from a zero value, and the
// preset contract is defined by key presence.
func decodeAddonsParamRaw(t *testing.T, qrURL *url.URL) map[string]interface{} {
	t.Helper()

	var payload map[string]interface{}
	if err := json.Unmarshal(decodeAddonsParamBody(t, qrURL), &payload); err != nil {
		t.Fatalf("unmarshal addons payload failed: %v", err)
	}
	return payload
}

func decodeAddonsParamBody(t *testing.T, qrURL *url.URL) []byte {
	t.Helper()

	encoded := qrURL.Query().Get("addons")
	if encoded == "" {
		t.Fatalf("missing addons param: %s", qrURL.String())
	}
	compressed, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		t.Fatalf("decode addons base64 failed: %v", err)
	}
	reader, err := gzip.NewReader(bytes.NewReader(compressed))
	if err != nil {
		t.Fatalf("create gzip reader failed: %v", err)
	}
	defer func() {
		if err := reader.Close(); err != nil {
			t.Fatalf("close gzip reader failed: %v", err)
		}
	}()

	body, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("read addons gzip body failed: %v", err)
	}
	return body
}

func boolPtr(value bool) *bool {
	return &value
}

func writeJSON(w http.ResponseWriter, body string) {
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write([]byte(body)); err != nil {
		panic("write test JSON response failed: " + err.Error())
	}
}
