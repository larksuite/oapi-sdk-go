package registration

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

func encodeAddons(addons *AppAddons) (string, error) {
	payload, err := normalizeAddons(addons)
	if err != nil {
		return "", err
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("registration: marshal Addons failed: %w", err)
	}

	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)
	if _, err := writer.Write(body); err != nil {
		if closeErr := writer.Close(); closeErr != nil {
			return "", fmt.Errorf("registration: gzip Addons failed: write: %w; close: %v", err, closeErr)
		}
		return "", fmt.Errorf("registration: gzip Addons failed: %w", err)
	}
	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("registration: gzip Addons failed: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(buf.Bytes()), nil
}

func normalizeAddons(addons *AppAddons) (map[string]interface{}, error) {
	if addons == nil {
		return nil, fmt.Errorf("registration: Addons is required")
	}

	itemCount := 0
	payload := make(map[string]interface{})

	if addons.Preset != nil {
		payload["preset"] = *addons.Preset
	}

	scopes, hasScopes, err := normalizeAddonsScopes(addons.Scopes, &itemCount)
	if err != nil {
		return nil, err
	}
	if hasScopes {
		payload["scopes"] = scopes
	}

	events, hasEvents, err := normalizeAddonsEvents(addons.Events, &itemCount)
	if err != nil {
		return nil, err
	}
	if hasEvents {
		payload["events"] = events
	}

	callbacks, hasCallbacks, err := normalizeAddonsCallbacks(addons.Callbacks, &itemCount)
	if err != nil {
		return nil, err
	}
	if hasCallbacks {
		payload["callbacks"] = callbacks
	}

	if itemCount == 0 && !isMinimalBase(addons.Preset) {
		return nil, fmt.Errorf("registration: Addons must contain at least one scope, event or callback, unless Preset is false")
	}
	return payload, nil
}

// isMinimalBase reports whether the addons explicitly select the minimal base
// template (Preset false), the only case where an empty increment set is valid.
func isMinimalBase(preset *bool) bool {
	return preset != nil && !*preset
}

func normalizeAddonsScopes(scopes AppAddonsScopes, itemCount *int) (map[string][]string, bool, error) {
	payload := make(map[string][]string)
	if values, ok, err := normalizeAddonsStringList(scopes.Tenant, "Addons.Scopes.Tenant", itemCount); err != nil {
		return nil, false, err
	} else if ok {
		payload["tenant"] = values
	}
	if values, ok, err := normalizeAddonsStringList(scopes.User, "Addons.Scopes.User", itemCount); err != nil {
		return nil, false, err
	} else if ok {
		payload["user"] = values
	}
	return payload, len(payload) > 0, nil
}

func normalizeAddonsEvents(events AppAddonsEvents, itemCount *int) (map[string]interface{}, bool, error) {
	items := make(map[string][]string)
	if values, ok, err := normalizeAddonsStringList(events.Items.Tenant, "Addons.Events.Items.Tenant", itemCount); err != nil {
		return nil, false, err
	} else if ok {
		items["tenant"] = values
	}
	if values, ok, err := normalizeAddonsStringList(events.Items.User, "Addons.Events.Items.User", itemCount); err != nil {
		return nil, false, err
	} else if ok {
		items["user"] = values
	}
	if len(items) == 0 {
		return nil, false, nil
	}
	return map[string]interface{}{
		"items": items,
	}, true, nil
}

func normalizeAddonsCallbacks(callbacks AppAddonsCallbacks, itemCount *int) (map[string][]string, bool, error) {
	values, ok, err := normalizeAddonsStringList(callbacks.Items, "Addons.Callbacks.Items", itemCount)
	if err != nil {
		return nil, false, err
	}
	if !ok {
		return nil, false, nil
	}
	return map[string][]string{
		"items": values,
	}, true, nil
}

func normalizeAddonsStringList(values []string, path string, itemCount *int) ([]string, bool, error) {
	if values == nil {
		return nil, false, nil
	}
	for idx, value := range values {
		if value == "" {
			return nil, true, fmt.Errorf("registration: %s[%d] must be a non-empty string", path, idx)
		}
	}
	*itemCount += len(values)
	return append([]string(nil), values...), true, nil
}
