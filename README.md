[English](README.md) | [简体中文](README_CN.md)

<div align="center">
  <a href="https://github.com/welcomemonth/web2api">
    <img src="https://img.shields.io/badge/Web-2API-1677ff?style=for-the-badge&logo=alibabacloud&logoColor=white" alt="Web2API" height="80">
  </a>

  <h1>Web2API</h1>

  <p>
    A browser-first Qwen Web gateway that drives the Qwen Web chat directly via Playwright, exposing OpenAI, Anthropic, and Gemini compatible APIs.
  </p>

  <p>
    <a href="https://github.com/welcomemonth/web2api">GitHub</a> ·
    <a href="./README_CN.md">中文说明</a>
  </p>

  <p>
    <a href="https://github.com/welcomemonth/web2api/releases">
      <img src="https://img.shields.io/github/v/release/welcomemonth/web2api?logo=github&label=Version&style=flat-square" alt="Release">
    </a>
    <a href="https://github.com/welcomemonth/web2api/stargazers">
      <img src="https://img.shields.io/github/stars/welcomemonth/web2api?logo=github&style=flat-square&label=Stars" alt="Stars">
    </a>
    <img src="https://img.shields.io/badge/Backend-Go%201.26-00ADD8?logo=go&style=flat-square" alt="Go">
    <img src="https://img.shields.io/badge/WebUI-React%2019-61DAFB?logo=react&style=flat-square" alt="React">
    <img src="https://img.shields.io/badge/Driver-Playwright-2EAD33?logo=playwright&style=flat-square" alt="Playwright">
    <img src="https://img.shields.io/badge/License-GPL--3.0-blue?style=flat-square" alt="License">
  </p>
</div>

## 一、项目简介 / Project Overview

web2api is a browser-first Qwen Web gateway: it drives the real Qwen Web chat interface with Playwright and exposes it as OpenAI, Anthropic, and Gemini compatible APIs. A local WebUI covers account management, downstream API keys, runtime settings, and model/image/video tests.

> [!NOTE]
> This project is a fork of [qwen2API](https://github.com/YuJunZhiXue/qwen2API). The main change is that it uses the Qwen **Web** interface directly (browser automation) instead of reverse-engineering the internal HTTP API. Each Qwen account is bound to a real Web chat session (`Chat_ID`), so only the current message is sent upstream instead of the full prompt history.

### 1. Feature Map

| Area | Capability |
| --- | --- |
| OpenAI-compatible APIs | `/v1/chat/completions`, `/v1/responses`, `/v1/models`, `/v1/files`, `/v1/images/generations`, `/v1/videos/generations` |
| Anthropic-compatible APIs | `/v1/messages`, `/anthropic/v1/messages`, `/v1/messages/count_tokens` |
| Gemini-compatible APIs | `/v1beta/models/{model}:generateContent`, `/v1beta/models/{model}:streamGenerateContent` |
| Browser upstream | Playwright-driven Qwen Web chat, headless browser pool, `Chat_ID` session binding and prewarm |
| WebUI | Accounts, API keys, runtime config, chat test, image test, video test |
| Account pool | Multi-account rotation, per-account concurrency, separate chat/image/video cooldown tracking |
| Operations | `/healthz`, `/readyz`, `/keepalive`, Docker healthcheck, multi-arch image publishing |

### 2. Version Line

| Version | Stack | Status |
| --- | --- | --- |
| `qwen2API v1.0` | Python + FastAPI/Uvicorn | Upstream legacy version, kept only as historical context |
| `qwen2API v2.0` | Go backend + React WebUI | Upstream mainline that this fork is based on |
| `web2api` | Go backend + React WebUI + Playwright | This project, browser-first |

## 二、快速部署 / Quick Deployment

### 1. Docker Deployment

For most deployments, use the Docker image directly. Keep `data` and `logs` beside your compose file; Docker will mount them into the container so upgrades do not wipe your accounts, keys, or logs.

> [!NOTE]
> If the `welcomemonth/web2api` image is not published yet, use [Build Locally With Docker](#2-build-locally-with-docker) instead.

```bash
mkdir web2api
cd web2api
mkdir -p data logs
```

Create a small `.env`:

```env
HOST_PORT=7860
HOST_DATA_DIR=./data
HOST_LOGS_DIR=./logs
ADMIN_KEY=replace-with-your-own-strong-random-key
```

Create `docker-compose.yml`:

```yaml
services:
  web2api:
    image: ${WEB2API_IMAGE:-welcomemonth/web2api:latest}
    container_name: web2api
    restart: unless-stopped
    init: true
    env_file:
      - .env
    ports:
      - "${HOST_PORT:-7860}:${PORT:-7860}"
    volumes:
      - ${HOST_DATA_DIR:-./data}:/app/data
      - ${HOST_LOGS_DIR:-./logs}:/app/logs
    shm_size: "512m"
    healthcheck:
      test: ["CMD-SHELL", "curl -fsS http://127.0.0.1:${PORT:-7860}/healthz || exit 1"]
      interval: 30s
      timeout: 10s
      start_period: 120s
      retries: 3
```

You do not need to set paths for `accounts.json`, `api_keys.json`, or other internal files. The image already uses `/app/data` and `/app/logs`; the volume mapping above decides where those files live on your host.

Start it:

```bash
docker compose pull
docker compose up -d
docker compose logs -f web2api
```

Open:

- WebUI: `http://127.0.0.1:7860`
- Health check: `http://127.0.0.1:7860/healthz`
- Keepalive probe: `http://127.0.0.1:7860/keepalive`

### 2. Build Locally With Docker

Use this path when you changed the source code and need to build your own image.

```bash
git clone https://github.com/welcomemonth/web2api.git
cd web2api
cp .env.example .env
docker compose -f docker-compose.yml -f docker-compose.build.yml build
docker compose -f docker-compose.yml -f docker-compose.build.yml up -d
```

## 三、架构与配置 / Architecture and Configuration

### 1. Runtime Architecture

```mermaid
flowchart LR
  subgraph Clients["API Clients"]
    OpenAI["OpenAI SDK / Chat Completions"]
    Anthropic["Claude / Anthropic Messages"]
    Gemini["Gemini-compatible clients"]
    CLI["Claude Code / Codex / other CLI tools"]
  end

  subgraph App["web2api"]
    WebUI["React WebUI"]
    Router["Go HTTP Router"]
    Adapter["Protocol Adapters"]
    Tools["Tool-call / Context Pipeline"]
    Pool["Qwen Account Pool"]
    Browser["Playwright Browser Pool"]
    Store["JSON Stores / Data Files"]
  end

  subgraph Runtime["Runtime"]
    Docker["Docker image"]
    Data["./data volume"]
    Logs["./logs volume"]
  end

  Qwen["Qwen Web Chat"]

  OpenAI --> Router
  Anthropic --> Router
  Gemini --> Router
  CLI --> Router
  WebUI --> Router
  Router --> Adapter
  Adapter --> Tools
  Tools --> Pool
  Pool --> Browser
  Browser --> Qwen
  Pool --> Store
  Store --> Data
  Router --> Logs
  Docker --> App
```

### 2. Environment Variables

Do not commit real secrets. `.env.example` intentionally contains empty values and commented examples only.

| Variable | Description |
| --- | --- |
| `ADMIN_KEY` | WebUI and `/api/admin/*` management key. Set a strong private value. |
| `QWEN_API_KEY`, `QWEN_API_KEYS`, `QWEN_API_KEY_N` | Runtime-only downstream API keys injected from env. They are not saved to `data/api_keys.json` and cannot be deleted from WebUI. |
| `QWEN_ACCOUNT_N` | Runtime-only upstream Qwen account, format `token;optional-email;optional-password`. It is not saved to `data/accounts.json`. |
| `BROWSER_POOL_SIZE` | Number of headless browser instances in the pool. Default `1`. |
| `BROWSER_STREAM_TIMEOUT_SECONDS` | Timeout for a browser-driven stream. Default `1800`. |
| `MAX_INFLIGHT_PER_ACCOUNT` | Maximum concurrent requests per Qwen account. Default `2`. |
| `CHAT_ID_PREWARM_TARGET_PER_ACCOUNT` | Number of Qwen Web chat sessions to pre-create per account. Default `5`. |
| `CHAT_ID_PREWARM_TTL_SECONDS` / `CHAT_ID_PREWARM_MAX_CONCURRENCY` | TTL and max concurrency for the `Chat_ID` prewarm pool. Defaults `120` / `16`. |
| `KEEPALIVE_URL`, `KEEPALIVE_INTERVAL` | Optional background keepalive task. Env values lock the same WebUI settings. |
| `TOOL_RECOVERY_MAX_ATTEMPTS` | Maximum automatic recovery attempts when an upstream response after a tool result fails to produce the next client tool call. Default `4`, clamped to `1`-`8`. |
| `HOST_DATA_DIR`, `HOST_LOGS_DIR` | Host paths mounted into Docker as `/app/data` and `/app/logs`. Defaults are `./data` and `./logs`. |
| `DATA_DIR`, `LOGS_DIR` | Local non-Docker path overrides. Leave empty to use the current project directory. |

## 四、开发指南 / Development Guide

### 1. Requirements

- Go `1.26`
- Node.js `20+`
- npm
- Docker, only if you need container builds
- Playwright browsers (run with `--install-browsers`, or `go run start-all.go --install-browsers`)

### 2. One-Command Local Startup

```powershell
go run start-all.go
```

### 3. Backend Development

```powershell
cd backend
go run .
```

Verification:

```powershell
cd backend
go test ./...
go build -trimpath -ldflags="-s -w" -o ..\bin\web2api-backend.exe .
```

### 4. Frontend Development

```powershell
cd frontend
npm ci
npm run dev
```

Production build:

```powershell
cd frontend
npm run build
```

### 5. Development Rules

- Keep the Go backend as the runtime source of truth.
- Keep Docker data paths container-internal as `/app/data` and `/app/logs`.
- Control host paths through compose volume mappings instead of hard-coded workspace paths.
- Do not commit `data/`, `logs/`, `.env`, real tokens, cookies, passwords, or downstream API keys.
- Update README and `.env.example` when adding user-visible configuration.

## 五、参与贡献 / Contribution

### 1. How to Contribute

- Report bugs through [GitHub Issues](https://github.com/welcomemonth/web2api/issues).
- Submit feature requests through [GitHub Issues](https://github.com/welcomemonth/web2api/issues).
- Open focused pull requests through [GitHub Pull Requests](https://github.com/welcomemonth/web2api/pulls).
- Include practical verification steps when possible.

### 2. Pull Request Checklist

- `go test ./...` passes in `backend`.
- `npm run build` passes in `frontend`.
- Docker-related changes are reflected in `Dockerfile`, `docker-compose.yml`, and README when needed.
- No generated data, logs, local `.env`, Qwen token, cookie, password, or downstream API key is included.

### 3. Contributors

Thanks to everyone who helps improve web2api.

[![Contributors](https://contrib.rocks/image?repo=welcomemonth/web2api)](https://github.com/welcomemonth/web2api/graphs/contributors)

## 六、其他信息 / Other Information

### 1. Star History

[![Star History Chart](https://api.star-history.com/svg?repos=welcomemonth/web2api&type=Timeline)](https://www.star-history.com/#welcomemonth/web2api&Timeline)

### 2. License

This project is released under the [GPL-3.0 License](./LICENSE).

### 3. Disclaimer

- This project is provided as an open-source self-hosted gateway.
- Review your local laws, platform rules, and upstream account policies before deployment.
- Do not publish or share real account tokens, cookies, passwords, or downstream API keys.
- If you find a security issue, please avoid public secret disclosure and report it through a private channel first.

### 4. Acknowledgements

- 特别鸣谢: [qwen2API](https://github.com/YuJunZhiXue/qwen2API) — this project is a fork of qwen2API, thanks to [YuJunZhiXue](https://github.com/YuJunZhiXue) and all qwen2API contributors.
- 特别鸣谢: [LinuxDo](https://linux.do/)

---

<div align="center">
  <p>If web2api helps you, consider giving the project a Star.</p>
  <p>Made by <a href="https://github.com/welcomemonth">welcomemonth</a>, based on <a href="https://github.com/YuJunZhiXue/qwen2API">qwen2API</a> by <a href="https://github.com/YuJunZhiXue">YuJunZhiXue</a>.</p>
</div>
