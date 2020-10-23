package v1

type AccessTokenReqBody struct {
	GrantType string `json:"grant_type"`
	Code      string `json:"code"`
}

type RefreshAccessTokenReqBody struct {
	GrantType    string `json:"grant_type"`
	RefreshToken string `json:"refresh_token"`
}

type AccessTokenResult struct {
	AccessToken      string `json:"access_token"`
	AvatarUrl        string `json:"avatar_url"`
	AvatarThumb      string `json:"avatar_thumb"`
	AvatarMiddle     string `json:"avatar_middle"`
	AvatarBig        string `json:"avatar_big"`
	ExpiresIn        int    `json:"expires_in"`
	Name             string `json:"name"`
	EnName           string `json:"en_name"`
	OpenID           string `json:"open_id"`
	UnionID          string `json:"union_id"`
	UserID           string `json:"user_id"`
	TenantKey        string `json:"tenant_key"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
}

type UserInfoResult struct {
	Name         string `json:"name"`
	AvatarUrl    string `json:"avatar_url"`
	AvatarThumb  string `json:"avatar_thumb"`
	AvatarMiddle string `json:"avatar_middle"`
	AvatarBig    string `json:"avatar_big"`
	Email        string `json:"email"`
	OpenID       string `json:"open_id"`
	UnionID      string `json:"union_id"`
	UserID       string `json:"user_id"`
	Mobile       string `json:"mobile"`
}
