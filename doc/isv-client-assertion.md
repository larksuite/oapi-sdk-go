# 商店应用 ClientAssertion 认证

本功能适用于平台已为其开放 ClientAssertion 认证的商店应用，目前服务端支持范围为字节内部场景。应用需要启用 `private_key_jwt`，目标租户需要安装应用并授予相应权限。

## 初始化

继续使用现有 `ClientAssertionProvider` 接口提供已签发的 Assertion。SDK 不生成或保管私钥，不需要 AppSecret：

```go
client := lark.NewClient(appID, "",
    lark.WithMarketplaceApp(),
    lark.WithClientAssertionProvider(provider),
    lark.WithReqTimeout(requestTimeout),
)
```

`provider` 实现 `larkcore.ClientAssertionProvider`，示例中的 `requestTimeout` 从调用方配置获取。超时配置沿用已有 SDK 选项，不作为 ISV 新增的必填配置；也可以为每次调用的 `ctx` 设置截止时间。Provider 应遵守传入的 `ctx`，避免取消后继续执行。

自定义环境通过 `WithOpenBaseUrl` 和 `WithOAuthBaseUrl` 配置业务域名和 OAuth 域名。Provider 接收真实 OAuth host 作为 `aud`；已有 TargetInfo 代理方式继续生效。

## TenantAccessToken

业务接口使用 Tenant 身份时，按请求传入目标租户：

```go
resp, err := client.Get(ctx, apiPath, nil, larkcore.AccessTokenTypeTenant,
    larkcore.WithTenantKey(tenantKey),
)
if err != nil {
    return fmt.Errorf("call tenant API: %w", err)
}
```

示例中的 `apiPath` 由具体业务接口确定，`resp` 交由业务代码处理；生成的服务方法也使用相同的 `WithTenantKey` 选项。

SDK 使用 `client_credentials + client_id + client_assertion_type + client_assertion + tenant_key` 的 JSON 请求兑换 TAT，按应用、租户和 OAuth 端点缓存。缓存复用现有 TokenManager 和 Cache，TTL 采用响应中的 `expires_in` 减去原有的提前过期余量，扣减后为负数时按原有行为将 TTL 置为 0。缓存写入失败时记录警告，继续返回已获取的 Token。本次不调整内置缓存实现、容量或淘汰策略。

商店应用无密钥初始化和 TAT 兑换不需要 AppTicket 或 AppAccessToken。设置 `WithEnableTokenCache(false)` 后，沿用现有语义：业务 API 要求调用方显式传入 Token，不再自动兑换。

## UserAccessToken

用户授权后，使用原有的授权码兑换方法：

```go
resp, err := client.AccessToken.RetrieveByAuthorizationCode(ctx,
    authorizationcode.NewTokenRequestBuilder().
        Code(code).
        RedirectUri(redirectURI).
        CodeVerifier(codeVerifier).
        Build(),
)
if err != nil {
    return fmt.Errorf("exchange user authorization code: %w", err)
}
```

`RedirectUri` 在授权流程中传过时保持一致；`CodeVerifier` 用于 PKCE，传入生成对应 challenge 的原始值。按实际授权流程省略未使用的可选字段。UAT 兑换不传 `tenant_key`，返回的 `resp.Data.AccessToken` 为 UAT。

刷新继续使用 `client.AccessToken.Refresh` 和 `refreshtoken.NewTokenRequestBuilder().RefreshToken(refreshToken).Build()`。用户授予 `offline_access` 时才会返回 RefreshToken；刷新后调用方负责保存新返回的凭证。

使用 UAT 调用业务方法时传入 `WithUserAccessToken`，SDK 不会为该请求获取 TAT。Token 交换沿用 SDK 现有的请求发送与重试行为。

## 兼容性

已有 AppSecret 接入、自建应用无密钥流程和旧版显式 Token 接口保持原有行为。同时配置 Provider 和 AppSecret 时使用 Provider，失败后不回退 AppSecret。

请勿记录 Assertion、Token、授权码或包含这些内容的请求/响应。持久化 UAT 或 RefreshToken 时，由调用方提供加密、访问控制及轮换管理。
