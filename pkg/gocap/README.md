# gocap

一个可嵌入业务服务的 Go CAPTCHA 库（单机内存版），实现 `challenge -> redeem -> siteverify` 核心闭环。

## 为什么用 gocap

- 无需额外部署独立 CAPTCHA 服务
- 直接挂载到现有 `net/http` 或 Gin 路由
- 内置防重放与一次性 token 消费语义
- 单机内存存储，接入简单

## 当前边界

- 支持：单机、内存态、核心闭环
- 不支持（当前版本）：分布式一致性、持久化存储、管理后台、instrumentation

## 与官方实现对比（`@cap.js/server`）

> 目标不是逐行复刻官方实现，而是在 Go 生态里保持核心协议兼容，同时强化可嵌入性与工程可维护性。

| 维度 | 官方 `@cap.js/server`（文档语义） | 本项目 `gocap`（当前实现） | 取舍与考量 |
| --- | --- | --- | --- |
| 形态定位 | 偏“服务端库”：提供 `createChallenge` / `redeemChallenge` / `validateToken` 方法，由业务框架自行挂路由 | 同时提供库能力 + 开箱即用 HTTP Handler（`challenge/redeem/siteverify`） | Go 业务常希望直接挂 `net/http`/Gin，降低接入成本 |
| 对外流程 | `challenge -> redeem -> validateToken/siteverify` | `challenge -> redeem -> siteverify` | 保持核心闭环一致，便于对接 widget 与后端校验 |
| `siteverify` 请求语义 | `secret + response` | `secret + response` | 与官方/recaptcha 风格保持一致，迁移成本低 |
| 存储抽象 | 文档强调 `challenges`/`tokens` 存储接口，可接数据库 | 当前默认内存存储（`memstore`），并抽象 `store.Store` 接口 | 先做单机简化，保留后续扩展 Redis/DB 的演进点 |
| challenge token 机制 | 官方内部实现细节可替换，强调 API 行为 | 使用 JWT-like 签名 challenge token（HS256）+ 服务端验签 | 在 Go 中实现简单、可审计，减少额外依赖 |
| 防重放策略 | 验证 token 的一次性语义（默认） | 显式记录 challenge 签名已使用 + redeem token 消费 | 更直观可控，便于排查重复提交/重放问题 |
| 默认参数与 TTL | 官方文档给出默认 challenge 参数（如 `50/32/4`、`expiresMs`） | 当前默认值与 TTL 策略按本项目配置为主 | 偏向服务端可配置与业务稳态，可按需再向官方默认对齐 |
| 限流/CORS/Body 限制 | 官方示例通常由框架或 standalone 层处理 | 内建可选限流、CORS、请求体大小限制 | 将通用防护前置到 transport，减少业务侧重复代码 |
| 多站点能力 | Standalone 场景支持多 site key | 引擎支持 `RegisterSite` 动态注册站点配置 | 贴近多租户/多业务线场景，便于统一接入 |
| instrumentation 挑战 | 官方生态支持（尤其 standalone） | 当前未实现（README 已声明） | 控制复杂度，先聚焦 PoW 核心闭环 |
| 分布式一致性/持久化 | 官方可结合数据库部署 | 当前为单机内存态 | 先保证小体量稳定，再逐步扩展分布式能力 |
| 错误模型 | 官方返回以成功语义为主（框架层可自定义） | 统一错误码（`bad_request`/`forbidden`/`rate_limit` 等） | 便于业务监控、告警归类与灰度排障 |

### 结论

- 协议层面：本项目已与官方核心流程保持同向兼容。
- 工程层面：本项目更偏 Go 服务内嵌与防护增强，不追求与官方内部实现逐行一致。
- 演进策略：在保持当前架构的前提下，可逐步补齐持久化存储、分布式一致性与 instrumentation 能力。

## 快速开始（net/http）

```go
package main

import (
	"log"
	"net/http"
	"time"

	"github.com/lin-snow/ech0/pkg/gocap/cap"
)

func main() {
	engine, err := cap.New(
		cap.WithInMemoryStore(),
		cap.WithEnableCORS(true),
		cap.WithRateLimit(30, 5*time.Second),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer engine.Close()

	if err := engine.RegisterSite(cap.SiteRegistration{
		SiteKey:        "example-site",
		Secret:         "example-secret",
		Difficulty:     4,
		ChallengeCount: 80,
		SaltSize:       32,
	}); err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.Handle("/cap/", http.StripPrefix("/cap", engine.Handler()))
	log.Fatal(http.ListenAndServe(":8080", mux))
}
```

## Gin 接入

```go
r.Any("/cap/*any", gin.WrapH(http.StripPrefix("/cap", engine.Handler())))
```

## HTTP API

支持两种路径格式：

- `/{siteKey}/challenge`
- `/{siteKey}/redeem`
- `/{siteKey}/siteverify`

或带前缀：

- `/cap/{siteKey}/challenge`
- `/cap/{siteKey}/redeem`
- `/cap/{siteKey}/siteverify`

### `POST /{siteKey}/challenge`

返回：

```json
{
  "challenge": { "c": 80, "s": 32, "d": 4 },
  "token": "xxx.yyy.zzz",
  "expires": 1760000000000
}
```

### `POST /{siteKey}/redeem`

请求：

```json
{
  "token": "xxx.yyy.zzz",
  "solutions": [123, 456]
}
```

返回：

```json
{
  "success": true,
  "token": "redeem_token",
  "expires": 1760000000000
}
```

说明：`redeem` 允许扩展字段（例如 `instr`、`instr_timeout`、`instr_blocked`），不会因为未知字段直接失败。

### `POST /{siteKey}/siteverify`

请求：

```json
{
  "secret": "example-secret",
  "response": "redeem_token"
}
```

返回：

```json
{ "success": true }
```

## 错误响应

统一格式：

```json
{
  "success": false,
  "code": "bad_request",
  "error": "Malformed JSON body"
}
```

常见错误码：

- `bad_request`
- `forbidden`
- `not_found`
- `rate_limit`
- `internal`
- `method_not_allowed`

## 默认行为

- challenge TTL：15 分钟
- redeem token TTL：2 小时
- 默认限流：
  - `challenge` 开启
  - `redeem` 关闭（可配置开启）
  - `siteverify` 关闭（可配置开启）
- 请求体大小限制：1 MiB（可配置）

## 配置项（Option）

- `WithChallengeTTL(time.Duration)`
- `WithRedeemTTL(time.Duration)`
- `WithGCInterval(time.Duration)`
- `WithSecretPepper([]byte)`
- `WithStore(store.Store)`
- `WithInMemoryStore()`
- `WithRateLimit(max int, window time.Duration)`
- `WithRateLimitScope(scope string)`
- `WithRateLimitOnRedeem(bool)`
- `WithRateLimitOnSiteVerify(bool)`
- `WithEnableCORS(bool)`
- `WithIPHeader(string)`
- `WithMaxBodyBytes(int64)`

