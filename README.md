# ai-mate-client-gateway

Go Client Gateway for AI-Mate.

This service is the user-facing API entry for `chat-kiro`. New routes use `/api/*`; migration from the Java Gateway can temporarily keep `/api/client/*` and user-facing `/api/common/*` through reverse-proxy rewrite when needed.

## Run

```powershell
go run ./cmd/server
```

## Verify

```powershell
go test ./...
```
