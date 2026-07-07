package registration

const (
	StatusPolling        = "polling"
	StatusSlowDown       = "slow_down"
	StatusDomainSwitched = "domain_switched"
)

type QRCodeInfo struct {
	URL      string
	ExpireIn int
}

type StatusChangeInfo struct {
	Status   string
	Interval int
}

// AppPreset contains values used to pre-fill the app creation page after the
// user opens the QR code URL. Users can still edit these values on the page;
// the final app values are whatever the user submits.
type AppPreset struct {
	// Avatar contains app avatar URLs. Use one entry for a single avatar, or
	// 1-6 entries for candidates. The first entry is selected by default.
	// Values are URL-encoded by the SDK when building the QR URL.
	Avatar []string

	// Name is the app name shown on the app creation page. It supports the
	// {user} placeholder, which is resolved by the web page for the scanning
	// user.
	Name string

	// Desc is the app description shown on the app creation page. It supports
	// the {user} placeholder, which is resolved by the web page for the scanning
	// user.
	Desc string
}

// AppAddons contains incremental app config pre-filled into the confirmation
// page shown after the user opens the QR code URL. Only public scopes, events,
// and callbacks are supported. Sensitive config must be updated through the
// application config OpenAPI instead.
type AppAddons struct {
	// Preset selects the base template for app creation. It is unrelated to
	// AppPreset, which pre-fills the app name, description and avatar.
	//
	//   - nil (default): the platform default template; the encoded payload
	//     carries no preset key, identical to not setting this field.
	//   - false: the minimal base template; the confirmation page shows only
	//     the scopes, events and callbacks declared in this AppAddons. With
	//     Preset set to false, an AppAddons without any scope, event or
	//     callback is valid.
	//   - true: explicitly requests the platform default template, same base
	//     as nil.
	Preset *bool `json:"preset,omitempty"`

	// Scopes contains app-identity and user-identity permission scopes.
	Scopes AppAddonsScopes `json:"scopes,omitempty"`

	// Events contains app-identity and user-identity event subscriptions.
	Events AppAddonsEvents `json:"events,omitempty"`

	// Callbacks contains callback names.
	Callbacks AppAddonsCallbacks `json:"callbacks,omitempty"`
}

type AppAddonsScopes struct {
	// Tenant contains app-identity scopes, for example
	// im:message:send_as_bot.
	Tenant []string `json:"tenant,omitempty"`

	// User contains user-identity scopes, for example calendar:calendar:read.
	User []string `json:"user,omitempty"`
}

type AppAddonsEvents struct {
	Items AppAddonsEventItems `json:"items,omitempty"`
}

type AppAddonsEventItems struct {
	// Tenant contains app-identity events, for example im.message.receive_v1.
	Tenant []string `json:"tenant,omitempty"`

	// User contains user-identity events, for example
	// calendar.calendar.event.changed_v4.
	User []string `json:"user,omitempty"`
}

type AppAddonsCallbacks struct {
	// Items contains callback names, for example card.action.trigger.
	Items []string `json:"items,omitempty"`
}

type UserInfo struct {
	OpenID      string
	TenantBrand string
}

type RegisterAppResult struct {
	ClientID     string
	ClientSecret string
	UserInfo     *UserInfo
}

type Options struct {
	Source         string
	Domain         string
	LarkDomain     string
	AppPreset      *AppPreset
	Addons         *AppAddons
	CreateOnly     bool
	AppID          string
	OnQRCode       func(info *QRCodeInfo)
	OnStatusChange func(info *StatusChangeInfo)
}

type beginResponse struct {
	DeviceCode              string `json:"device_code"`
	VerificationURIComplete string `json:"verification_uri_complete"`
	VerificationURI         string `json:"verification_uri"`
	UserCode                string `json:"user_code"`
	Interval                int    `json:"interval"`
	ExpireIn                int    `json:"expire_in"`
}

type pollResponse struct {
	ClientID     string    `json:"client_id,omitempty"`
	ClientSecret string    `json:"client_secret,omitempty"`
	UserInfo     *userInfo `json:"user_info,omitempty"`
	Error        string    `json:"error,omitempty"`
	ErrorDesc    string    `json:"error_description,omitempty"`
}

type userInfo struct {
	OpenID      string `json:"open_id,omitempty"`
	TenantBrand string `json:"tenant_brand,omitempty"`
}
