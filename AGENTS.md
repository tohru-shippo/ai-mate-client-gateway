# AGENTS.md

本文件补充 `ai-mate-client-gateway/` 的项目级规范；根目录 `../AGENTS.md` 中的全局规则继续生效。

## 1. 项目定位

- 技术栈：Go 1.26.3 + Gin + Viper + slog + OpenTelemetry。
- 职责：AI-Mate 用户端 HTTP 入口，面向 `chat-kiro`，目标路径统一为 `/api/*`。
- 迁移期：按路径逐步承接 Java Gateway 现有 `/api/client/*` 和用户态 `/api/common/*` 能力；新路由优先使用去掉入口前缀后的 `/api/*`。
- 当前阶段：先建立干净的 Client Gateway 架构，不做业务表直连，不做 Brain 直连；用户端能力按模块逐步通过 `ai-mate-server` gRPC 接入。

## 2. Git 工作流

- 默认工作分支为 `dev`，与其他 AI-Mate 项目保持一致。
- 禁止在 `main` / `master` 上直接做业务提交或推送。
- 提交格式遵守根目录规范：`功能`、`修复`、`优化`、`文档` 四类中文前缀。

## 3. 分层约定

- `cmd/server/`：服务启动、优雅关闭、资源初始化。
- `configs/`：本地和生产环境配置。
- `internal/app/`：配置、日志、链路追踪、HTTP 服务装配。
- `internal/middleware/`：请求 ID、访问日志、用户鉴权和入口态中间件。
- `internal/handler/`：HTTP 接口入参解析、响应投影、错误响应。
- `internal/client/core/`：调用 `ai-mate-server` 的 gRPC client。
- `internal/service/`：仅放用户鉴权、OAuth 回调、验证码、token 等入口态协议适配，不承载核心业务写表逻辑。
- `internal/response/`：统一用户端响应结构。

## 4. API 约定

- 对外 HTTP 前缀统一为 `/api`，目标路径不再保留 `/client` 前缀。
- 迁移期如需兼容旧路径，优先通过 Nginx rewrite 承接；只有确认前端切换成本更高时才在 Gateway 中增加兼容路由。
- 对外 JSON 字段统一使用 `camelCase`。
- 成功响应统一为 `{ "code": 0, "message": "success", "data": any }`。
- 错误响应统一为 `{ "code": int, "message": string }`，不返回 `data`。
- 用户鉴权、语言选择、当前用户上下文注入和 HTTP 响应投影在 Client Gateway 完成。
- 核心业务读写通过 `ai-mate-server` 完成，不直接访问其他服务内部表。

## 5. 错误与语言

- Core gRPC 错误必须通过 `internal/errors/core.go` 转换成用户端错误响应。
- Client Gateway 默认语言为 `en-US`。
- 语言选择优先级：已登录用户上下文中的 `language` -> `X-Language` -> `en-US`。
- 当前未接入真实用户鉴权前，先按 `X-Language` -> `en-US` 处理。
- 支持语言白名单与 `chat-kiro` 对齐：`en-US`、`zh-CN`、`zh-TW`、`ja-JP`、`ko-KR`、`ru-RU`。
- `X-Language` 必须传完整语言码；不支持或为空时回退 `en-US`。

## 6. 代码注释规范

- 所有手写的 Go `struct`、`interface`、`type`、导出函数和方法必须有注释。
- 结构体字段必须补字段注释，说明字段含义、来源、取值约束或对外语义。
- Handler 注释必须说明对应 HTTP 路由、鉴权要求和主要响应语义。
- 依赖装配问题必须在启动阶段失败，不要在 Handler、middleware 等运行期路径里写无意义兜底。

## 7. 验证命令

```powershell
go test ./...
```
