# OpenAPI 无秘钥改造 — SDK & 长连接 实现任务清单

> 代码仓库：`larksuite/oapi-sdk-go`，分支：`v3_main`
>
> 所有设计决策已锁定，可直接进入编码实现。

---

## 设计决策总览

| 决策项 | 结论 |
|---|---|
| 接口签名 | `RetrieveToken(ctx context.Context, aud string) (*Token, error)` |
| aud 格式 | 纯域名，如 `open.feishu.cn`、`open.larksuite.com`、`fsopen.bytedance.net` |
| aud 来源（SDK） | 从 `Config.BaseUrl` 提取域名部分 |
| aud 来源（长连接） | 从长连接 Domain 提取域名部分 |
| Token Endpoint | `/open-apis/authen/v2/oauth/token`，JSON Body |
| Token Endpoint Body 字段 | `grant_type`、`client_assertion_type`、`client_assertion`、`app_id`（全 snake_case，不含 scope） |
| grant_type | `urn:ietf:params:oauth:grant-type:jwt-bearer` |
| AppAccessToken | ClientAssertion 模式下不需要，不调用 `/auth/v3/app_access_token/internal` |
| UserAccessToken | 本次不涉及 SDK 自动获取，仍通过 `WithUserAccessToken` 手动传入 |
| 商店应用 | 不支持 ClientAssertion，仅限自建应用 |
| 优先级 | ClientAssertionProvider 优先，忽略 AppSecret |
| Option 命名 | SDK 和 ws 包统一用 `WithClientAssertionProvider` |
| TargetInfo.TargetService | **纯域名**，不携带 `https://`，SDK 拼接 URL 时需主动补 `https://` |
| TargetInfo 代理路径 | `https:// + TargetService + TargetPrefix + APIPath` |
| TargetPrefix 格式 | 始终以 `/` 开头，不以 `/` 结尾 |
| X-Target-Service Header | 值为 aud（纯域名），即 `RetrieveToken` 传入的值 |
| 代理适用范围 | 仅 Token 获取和 WS Bootstrap，普通 OpenAPI 请求不走代理 |
| 长连接 Bootstrap Body | `{"AppID": "xxx", "ClientAssertion": "xxx"}`（PascalCase），走代理和不走代理请求体一致 |
| 长连接重连 | 每次重连重新调用 `RetrieveToken`（JWT 有 TTL） |
| Token 缓存 | 沿用现有 TenantAccessToken 缓存逻辑 |
| 重试策略 | 沿用现有 SDK 重试策略 |
| ClientAssertion 过期/失效 | IAM 返回的错误直接透传给调用方，不自动重试 |

---

## 一、核心类型与接口定义

### 任务 1：定义 ClientAssertionProvider 接口及相关类型

**文件**：`core/config.go`

新增以下类型定义：

```go
type TargetInfo struct {
    TargetService string // 代理服务纯域名（不含 https://），如 proxy.example.com
    TargetPrefix  string // 路径前缀，始终以 / 开头、不以 / 结尾，如 /v1/proxy
}

type Token struct {
    Value      string
    TargetInfo *TargetInfo // nil 时直接请求 IAM / 长连接服务
}

type ClientAssertionProvider interface {
    RetrieveToken(ctx context.Context, aud string) (*Token, error)
}
```

在 `Config` 结构体中新增字段：

```go
ClientAssertionProvider ClientAssertionProvider
```

**要点**：
- `TargetInfo`、`Token`、`ClientAssertionProvider` 均为导出类型，定义在 `core` 包中
- `Token.TargetInfo` 为指针类型，nil 表示不走代理
- `TargetService` 为纯域名，SDK 在拼接 URL 时需主动补 `https://`

---

### 任务 2：新增常量与错误码

**文件**：`core/constants.go`

新增常量：

| 常量名 | 值 |
|---|---|
| `OAuthTokenUrlPath` | `/open-apis/authen/v2/oauth/token` |
| `GrantTypeJWTBearer` | `urn:ietf:params:oauth:grant-type:jwt-bearer` |
| `ClientAssertionTypeJWTBearer` | `urn:ietf:params:oauth:client-assertion-type:jwt-bearer` |
| `HeaderXTargetService` | `X-Target-Service` |

新增错误码：

| 错误码 | 常量名 | 含义 |
|---|---|---|
| 7100 | `ErrCodeClientAssertionProviderNotConfigured` | ClientAssertionProvider 未配置 / ClientAssertion 模式下不支持 AppAccessToken |
| 7101 | `ErrCodeClientAssertionTokenEmpty` | RetrieveToken 返回的 Token.Value 为空 |
| 7102 | `ErrCodeClientAssertionRetrieveFailed` | RetrieveToken 调用失败 |
| 7103 | `ErrCodeClientAssertionModeNotSupported` | 该 API 不支持 ClientAssertion 模式（仅支持 AppAccessToken 的 API） |

---

### 任务 3：新增 OAuthTokenResp 响应结构体

**文件**：`core/tokenmanager.go`（或独立文件）

```go
type OAuthTokenResp struct {
    Code                  int    `json:"code"`
    Error                 string `json:"error"`
    ErrorDescription      string `json:"error_description"`
    AccessToken           string `json:"access_token"`
    ExpiresIn             int    `json:"expires_in"`
    RefreshToken          string `json:"refresh_token"`
    RefreshTokenExpiresIn int    `json:"refresh_token_expires_in"`
    Scope                 string `json:"scope"`
    TokenType             string `json:"token_type"`
}
```

**要点**：与现有 `TenantAccessTokenResp` 不同——新结构使用 `access_token` 字段（OAuth2 标准），旧结构使用 `tenant_access_token` 字段。必须新建结构体。

---

## 二、SDK 入口改造

### 任务 4：新增 WithClientAssertionProvider Option

**文件**：`client.go`

```go
func WithClientAssertionProvider(provider core.ClientAssertionProvider) ClientOptionFunc {
    return func(config *core.Config) {
        config.ClientAssertionProvider = provider
    }
}
```

**使用方式**：

```go
client := lark.NewClient(appID, "", // appSecret 可为空
    lark.WithClientAssertionProvider(myProvider),
)
```

**要点**：
- 不新增专用构造函数（如 `NewClientWithZTI`），仅通过 Option 模式接入
- `appSecret` 参数传空字符串即可
- ClientAssertionProvider 优先级高于 AppSecret

---

## 三、工具函数

### 任务 5：新增 aud 提取函数

**文件**：`core/` 包内（如 `core/util.go`）

```go
func extractAudFromURL(rawURL string) (string, error)
```

- 输入：`https://open.feishu.cn/open-apis` → 输出：`open.feishu.cn`
- 使用 `net/url.Parse` 取 `.Host`
- 处理无 scheme 的情况（自动补全后解析）

---

### 任务 6：新增代理 URL 拼接函数

**文件**：`core/` 包内（如 `core/util.go`）

```go
func buildProxyURL(targetService, targetPrefix, apiPath string) string
```

- **`targetService` 为纯域名**（如 `proxy.example.com`），函数内部需主动补 `https://`
- 拼接规则：`"https://" + targetService + targetPrefix + apiPath`
- TargetPrefix 格式约定：始终以 `/` 开头，不以 `/` 结尾
- 示例：`proxy.example.com` + `/v1` + `/open-apis/authen/v2/oauth/token` → `https://proxy.example.com/v1/open-apis/authen/v2/oauth/token`

---

## 四、Token 获取流程改造

### 任务 7：修改 getTenantAccessToken 入口

**文件**：`core/tokenmanager.go`

在现有 `getTenantAccessToken` 入口处增加分支判断：

```go
if config.ClientAssertionProvider != nil {
    // 仅自建应用支持
    return getTenantTokenByClientAssertion(ctx)
}
// 现有逻辑不变
```

**要点**：
- ClientAssertionProvider 优先级高于 AppSecret
- 仅自建应用支持（`AppType != AppTypeISV`），商店应用不进入此分支

---

### 任务 8：新增 getTenantTokenByClientAssertion 方法

**文件**：`core/tokenmanager.go`

完整流程：

1. **提取 aud**：`extractAudFromURL(config.BaseUrl)` → 纯域名
2. **调用 Provider**：`provider.RetrieveToken(ctx, aud)` → `*Token`
3. **校验**：Token 为 nil 或 Value 为空 → 错误码 7101；调用失败 → 错误码 7102
4. **确定请求 URL**：
    - `Token.TargetInfo == nil`：`config.BaseUrl + OAuthTokenUrlPath`
    - `Token.TargetInfo != nil`：`buildProxyURL(TargetService, TargetPrefix, OAuthTokenUrlPath)`
        - 即 `"https://" + TargetService + TargetPrefix + "/open-apis/authen/v2/oauth/token"`
        - Header 添加 `X-Target-Service: {aud}`（纯域名）
5. **构造请求 Body**（JSON，全 snake_case）：
   ```json
   {
     "grant_type": "urn:ietf:params:oauth:grant-type:jwt-bearer",
     "client_assertion_type": "urn:ietf:params:oauth:client-assertion-type:jwt-bearer",
     "client_assertion": "<Token.Value>",
     "app_id": "<config.AppId>"
   }
   ```
6. **解析响应**：使用 `OAuthTokenResp` 结构体
7. **缓存 Token**：沿用现有 TenantAccessToken 缓存逻辑（key 含 appId，提前刷新策略不变）

---

### 任务 9：阻断 getAppAccessToken

**文件**：`core/tokenmanager.go`

在 `getAppAccessToken` 入口处增加判断：

```go
if config.ClientAssertionProvider != nil {
    return error(ErrCodeClientAssertionProviderNotConfigured,
        "AppAccessToken is not available in ClientAssertion mode")
}
```

**为什么需要处理 AppAccessToken**：现有 SDK 代码中 `tokenManager.getAppAccessToken()` 和 `determineTokenType()` 都有通向 AppAccessToken 的路径。在 ClientAssertion 模式下，这些路径无法正常工作（没有 AppSecret，无法调用 `/auth/v3/app_access_token/internal`）。必须显式拦截并返回明确错误码 7100，避免用户遇到难以排查的隐式错误。

---

## 五、HTTP 传输层改造

### 任务 10：修改 validate 方法

**文件**：`core/httptransport.go`

当前逻辑校验 `AppId` 和 `AppSecret` 非空。修改为：

- 当 `config.ClientAssertionProvider != nil` 时，**跳过 AppSecret 的非空校验**
- `AppId` 的校验始终保留

---

### 任务 11：修改 determineTokenType 方法

**文件**：`core/httptransport.go`

ClientAssertion 模式下的分派逻辑：

| API 支持的 TokenType | ClientAssertion 模式行为 |
|---|---|
| 仅 `AccessTokenTypeTenant` | 正常，使用 TenantAccessToken |
| 仅 `AccessTokenTypeApp` | **返回错误码 7103**，提示该 API 不支持 ClientAssertion 模式 |
| `AccessTokenTypeApp \| AccessTokenTypeTenant` | 选择 TenantAccessToken |
| `AccessTokenTypeUser` | 不受影响，仍通过 `WithUserAccessToken` 传入 |

---

### 任务 12：重试与错误透传

**文件**：`core/httptransport.go`

- 错误码 7102（RetrieveToken 调用失败）：沿用现有 SDK 重试策略（重试次数、间隔与其他可重试错误码一致）
- ClientAssertion 本身过期/失效的 IAM 返回错误：**直接透传给调用方**，不自动重试

---

## 六、长连接（WebSocket）改造

### 任务 13：新增字段与 Option

**文件**：`ws/client.go`

在 ws Client 结构体中新增：

```go
clientAssertionProvider core.ClientAssertionProvider
```

新增 Option：

```go
func WithClientAssertionProvider(provider core.ClientAssertionProvider) ClientOptionFunc {
    return func(client *Client) {
        client.clientAssertionProvider = provider
    }
}
```

**使用方式**：

```go
wsClient := ws.NewClient(appID, "",
    ws.WithClientAssertionProvider(myProvider),
)
```

---

### 任务 14：修改 Bootstrap 请求模型

**文件**：`ws/model.go`

在 Bootstrap 请求结构体中新增 `ClientAssertion` 字段：

```go
type BootstrapRequest struct {
    AppID           string `json:"AppID"`
    AppSecret       string `json:"AppSecret,omitempty"`
    ClientAssertion string `json:"ClientAssertion,omitempty"`
}
```

**要点**：`AppSecret` 和 `ClientAssertion` 互斥，均为 `omitempty`。PascalCase 风格与现有 ws Body 一致。

---

### 任务 15：修改 Bootstrap 请求逻辑

**文件**：`ws/client.go`

ClientAssertion 模式下的 Bootstrap 流程：

1. **提取 aud**：从长连接 Domain 提取纯域名
2. **调用 Provider**：`clientAssertionProvider.RetrieveToken(ctx, aud)` → `*Token`
3. **确定请求 URL**：
    - `Token.TargetInfo == nil`：`Domain + /callback/ws/endpoint`（现有逻辑）
    - `Token.TargetInfo != nil`：`buildProxyURL(TargetService, TargetPrefix, "/callback/ws/endpoint")`
        - 即 `"https://" + TargetService + TargetPrefix + "/callback/ws/endpoint"`
        - Header 添加 `X-Target-Service: {aud}`（纯域名）
4. **构造 Body**：`{"AppID": "<appID>", "ClientAssertion": "<Token.Value>"}`（走代理和不走代理请求体一致）
5. **响应格式不变**：使用现有 `EndpointResp` 结构体

---

### 任务 16：修改重连逻辑

**文件**：`ws/client.go`

每次 WebSocket 断线重连时：

1. **重新调用** `clientAssertionProvider.RetrieveToken(ctx, aud)` 获取新 Token（JWT 有 TTL，不能复用旧值）
2. 用新 Token 重新构造 Bootstrap Body 和请求 URL（TargetInfo 可能变化）
3. RetrieveToken 失败 → 错误直接透传给调用方

---

## 七、示例与测试

### 任务 17：新增使用示例

新增 sample 文件，演示：

- 自定义 `ClientAssertionProvider` 实现示例
- 使用 `WithClientAssertionProvider` 初始化 SDK Client + API 调用
- 使用 `WithClientAssertionProvider` 初始化 ws Client + 长连接建立

---

### 任务 18：新增单元测试

| 测试场景 | 所属文件 |
|---|---|
| Provider 未配置时行为不变 | `core/tokenmanager_test.go` |
| `getTenantTokenByClientAssertion` 正常流程 | `core/tokenmanager_test.go` |
| `getTenantTokenByClientAssertion` + TargetInfo 代理流程 | `core/tokenmanager_test.go` |
| RetrieveToken 返回空 Token → 错误码 7101 | `core/tokenmanager_test.go` |
| RetrieveToken 调用失败 → 错误码 7102 | `core/tokenmanager_test.go` |
| AppAccessToken 在 ClientAssertion 模式下返回错误码 7100 | `core/tokenmanager_test.go` |
| `determineTokenType` 对 AccessTokenTypeApp 返回错误码 7103 | `core/httptransport_test.go` |
| `validate` 跳过 AppSecret 校验 | `core/httptransport_test.go` |
| Provider 优先级高于 AppSecret | `core/httptransport_test.go` |
| aud 提取逻辑（多种 BaseUrl 格式） | `core/util_test.go` |
| `buildProxyURL` 拼接逻辑（含 https:// 补全验证） | `core/util_test.go` |
| 长连接 Bootstrap Body 构造（普通模式） | `ws/client_test.go` |
| 长连接 Bootstrap Body 构造（代理模式，验证 https:// 补全） | `ws/client_test.go` |
| 长连接重连时重新获取 Token | `ws/client_test.go` |

---

## 涉及文件清单

| 文件 | 变更类型 | 关联任务 |
|---|---|---|
| `core/config.go` | 新增类型 + 新增字段 | 任务 1 |
| `core/constants.go` | 新增常量 + 错误码 | 任务 2 |
| `core/tokenmanager.go` | 新增结构体 + 新增方法 + 修改入口 | 任务 3, 7, 8, 9 |
| `core/util.go`（新建） | 新增工具函数 | 任务 5, 6 |
| `core/httptransport.go` | 修改 validate / determineTokenType / 重试逻辑 | 任务 10, 11, 12 |
| `client.go` | 新增 Option | 任务 4 |
| `ws/client.go` | 新增字段 + Option + 修改 Bootstrap + 重连 | 任务 13, 15, 16 |
| `ws/model.go` | 新增 ClientAssertion 字段 | 任务 14 |
| `sample/`（新建） | 新增示例文件 | 任务 17 |
| `core/tokenmanager_test.go` | 新增测试 | 任务 18 |
| `core/httptransport_test.go` | 新增测试 | 任务 18 |
| `core/util_test.go`（新建） | 新增测试 | 任务 18 |
| `ws/client_test.go` | 新增测试 | 任务 18 |

---

## 执行顺序

```
阶段一（基础类型）：  任务 1 → 2 → 3 → 5 → 6
阶段二（SDK 核心）：  任务 4 → 7 → 8 → 9
阶段三（HTTP 层）：   任务 10 → 11 → 12
阶段四（长连接）：    任务 13 → 14 → 15 → 16
阶段五（收尾）：      任务 17 → 18
```

---

**总计 18 个任务，无待确认项。**
