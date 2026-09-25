# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this project is

`web2api` is a self-hosted gateway that drives the **Qwen Web chat** (chat.qwen.ai) and exposes it as OpenAI-, Anthropic-, and Gemini-compatible APIs, with a React WebUI for account/key management. It is a fork of [qwen2API](https://github.com/YuJunZhiXue/qwen2API) whose main change is "browser-first": it uses Playwright to operate the real Qwen Web interface instead of reverse-engineering the internal HTTP API.

## Commands

Run from the repo root unless noted.

- One-command local start (builds backend, installs frontend deps, starts both): `go run start-all.go`
- Install Playwright browsers (needed once, or before first run): `go run start-all.go --install-browsers`
- Backend (Go, in `backend/`):
  - Run: `go run .`
  - Build: `go build -trimpath -ldflags="-s -w" -o ../bin/web2api-backend.exe .`
  - Test: `go test ./...` — note there are currently **no test files**; this only verifies compilation.
- Frontend (React, in `frontend/`):
  - Install: `npm ci`
  - Dev server: `npm run dev`
  - Build: `npm run build` (runs `tsc --noEmit` typecheck, then `vite build`)
  - Lint: `npm run lint` / typecheck only: `npm run typecheck`
- Docker (build local image): `docker compose -f docker-compose.yml -f docker-compose.build.yml build` then `... up -d`

## Architecture

### Repository layout

- `backend/` — Go server (module name **`qwen2api-go`**, still the original name despite the rename).
- `frontend/` — React 19 + Vite 6 + Tailwind 4 + react-router 7 SPA.
- `start-all.go` — launcher that builds the backend, `npm install`s the frontend, and runs both with Vite proxying to the backend.
- `Dockerfile` — multi-stage build that also runs `--install-browsers` and ships the SPA under `frontend/dist`.
- `auth.json` — Playwright storage state (real Qwen login cookies/localStorage). Sensitive; do not commit.

### `backend/main.go` is the runtime source of truth

`backend/main.go` is a ~9,000-line monolith that contains the real implementation: the `App` struct, `Account`/`AccountPool`, `MailSession` (email verification), `ChatIDPool`, `Settings`/`LoadSettings`, `QwenClient`, `JSONStore`, the full route table (`routes()`), and every HTTP handler plus the large tool-call recovery/filtering logic. When fixing runtime behavior, this file is where it lives.

### Migrated-package scaffolding (duplication warning)

The sub-packages under `backend/` are a partial migration; some are live and some are dead scaffolding, and several symbols are **duplicated**. Only `main.go` imports them.

- **Live:** `adapter/` (`BuildChatStandardRequest`, `MessagesToPrompt`, `ModelMode`), `browser/` (`BrowserManager`), `toolcall/` (parse/format tool calls), `upstream/` (`ParseQwenEvent`, `BuildChatPayload`, `NormalizeChatType`, `ExtractUpstreamError`), `utils/` (`WriteJSON`, `WriteError`, `DecodeJSON`, `IsTransientUpstreamErrorMessage`), and part of `services/` (`truncation_recovery`, `model_catalog`, `tool_parser`).
- **Dead / self-check-only:** `api/` (`apidesc` route descriptors, only counted in `runMigratedPackageSelfCheck`), `core/` (only `core.ChooseEngine`/`core.LoadSettings` used in that self-check), `runtime/` (only `rt.VisibleText`).
- **Duplicated symbols** (edit `main.go`, not these, to change runtime behavior): `Settings`/`LoadSettings` exist in both `main.go` and `core/config.go`; `QwenHeaders` exists in `main.go`, `core/httpx_engine.go`, and `services/qwen_client.go`; `FormatUpstreamError` exists in `services/qwen_client.go` and `upstream/qwen_executor.go`. `services/qwen_client.go` as a whole is not the active client.

### Request flow

HTTP handler (`main.go` or `api/*.go`) → auth (`resolveAuth`: downstream API key from `api_keys.json`/env) → `adapter.BuildChatStandardRequest` (normalize the OpenAI/Anthropic/Gemini body into an internal `StandardRequest`) → `runCompletion` → `AccountPool.AcquireFor` (pick an account) → `QwenClient` (create/stream a `Chat_ID`) → tool-call parse + response formatting per surface.

### Upstream: two Qwen paths

Both ultimately target `https://chat.qwen.ai/api/v2/chat/completions`.

- `QwenClient.StreamChat` — the **browser path** and the primary one. Opens a new Playwright page, navigates the Qwen Web chat UI, and reads the stream from the browser **console**. `browser/brower.go` injects an init script that wraps `window.fetch`, clones the `/api/v2/chat/completions` response, and `console.log`s each SSE chunk prefixed `[QWEN-SSE-CHUNK]`; `StreamChat` listens via `page.OnConsole` and parses those lines.
- `QwenClient.StreamWebChat` — a direct bearer-token HTTP SSE call. Currently **unused** (no callers).
- `PostChatCompletionOnce` / `GetChatDetail` / `GetVisionTaskStatus` are used for image/video generation polling.

`browser.NewBrowserManager` launches a **headed** Firefox via Playwright (the `ws://localhost:9222/camoufox` endpoint argument is ignored/commented out) and loads `auth.json` as the browser storage state so it is already logged in.

### Chat_ID session model

Each Qwen account is bound to a real Web chat session (`chat_id`), so only the current message is sent upstream instead of full prompt history. `ChatIDPool` prewarms a pool of `chat_id`s keyed by `(email, model, chatType)` (`CHAT_ID_PREWARM_*` settings) and cleans them up on a TTL.

### Account pool & rate limiting

`AccountPool` rotates accounts, enforces per-account concurrency (`MAX_INFLIGHT_PER_ACCOUNT`), and keeps separate cooldowns per usage class (chat / image / video). Rate limits are marked via `MarkRateLimited`/`MarkRateLimitedFor` with base/max cooldown windows.

### Tool calling

Qwen does not emit structured tool calls, so tool calls are recovered from model text. `toolcall/` parses JSON, XML, and QNML (`<|QNML|...>`) encodings; `adapter/tool_profile_*.go` reformats them per client (Claude Code, Hermes, OpenClaw, OpenCode). `main.go` holds a large amount of filtering/recovery logic: repeated/out-of-order/invalid tool calls, truncated-call continuation, and prompt-based tool disallowal.

### Compatibility surfaces & admin API

- OpenAI: `/v1/chat/completions`, `/v1/responses`, `/v1/models`, `/v1/embeddings`, `/v1/images/generations`, `/v1/videos/generations`, `/v1/files`.
- Anthropic: `/v1/messages` (+ `/anthropic/v1/messages`), `/v1/messages/count_tokens`.
- Gemini: `/v1beta/models/{model}:generateContent`, `:streamGenerateContent`.
- Admin (WebUI backend), guarded by `Authorization: Bearer <ADMIN_KEY>`: `/api/admin/*` (accounts, keys, settings, verify), plus `/healthz`, `/readyz`, `/keepalive`, and the SPA at `/`.

OpenAI model names are mapped to internal Qwen models in `modelMap` (e.g. `gpt-4o` → `qwen3.6-plus`).

### Configuration & data files

- Config is env-var driven; `.env.example` is the authoritative list. `main.go`'s `LoadSettings()` is what actually loads it (see duplication note).
- Persistent state lives as JSON files under `data/` (`accounts.json`, `api_keys.json`, `users.json`, `captures.json`, `config.json`, `context_cache.json`, `uploaded_files.json`, `session_affinity.json`), managed by the `JSONStore` atomic file store in `main.go`.

## Gotchas

- The Go module is still named `qwen2api-go`; imports use `qwen2api-go/...`. Don't "fix" this casually — renaming touches every file.
- There are no Go tests. `go test ./...` is a compile check only.
- Real secrets must never be committed: `auth.json` (browser cookies), `data/`, `logs/`, `.env`, Qwen tokens/passwords, and downstream API keys.
- `start-all.go` hardcodes a Go fallback path (`D:\go\bin\go.exe`); the backend exe name is still `qwen2api-backend` / `qwen2api`.
