[English](README.md) | [简体中文](README_CN.md)

<div align="center">
  <a href="https://github.com/welcomemonth/web2api">
    <img src="https://img.shields.io/badge/Web-2API-1677ff?style=for-the-badge&logo=alibabacloud&logoColor=white" alt="Web2API" height="80">
  </a>

  <h1>Web2API</h1>

  <p>
    浏览器优先的千问 Web 网关，通过 Playwright 直接驱动千问 Web 对话，对外提供 OpenAI、Anthropic、Gemini 兼容接口。
  </p>

  <p>
    <a href="https://github.com/welcomemonth/web2api">GitHub</a> ·
    <a href="./README.md">English README</a>
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

## 一、项目简介

web2api 是一个浏览器优先的千问 Web 网关：通过 Playwright 直接驱动千问 Web 对话界面，对外提供 OpenAI、Anthropic、Gemini 兼容接口。本地 WebUI 用于管理上游账号、下游 API Key、运行配置，以及模型/图片/视频测试。

> [!NOTE]
> 本项目是基于 [qwen2API](https://github.com/YuJunZhiXue/qwen2API) 的二次开发。主要改动是改为直接使用千问 **Web** 界面（浏览器自动化），而不是逆向内部 HTTP 接口。每个千问账号绑定一个真实的 Web 会话（`Chat_ID`），因此上游只需发送当前消息，无需传递完整的历史 prompt。

### 1. 功能概览

| 模块 | 能力 |
| --- | --- |
| OpenAI 兼容接口 | `/v1/chat/completions`、`/v1/responses`、`/v1/models`、`/v1/files`、`/v1/images/generations`、`/v1/videos/generations` |
| Anthropic 兼容接口 | `/v1/messages`、`/anthropic/v1/messages`、`/v1/messages/count_tokens` |
| Gemini 兼容接口 | `/v1beta/models/{model}:generateContent`、`/v1beta/models/{model}:streamGenerateContent` |
| 浏览器上游 | Playwright 驱动的千问 Web 对话、无头浏览器池、`Chat_ID` 会话绑定与预热 |
| WebUI | 账号管理、API Key 管理、运行配置、对话测试、图片测试、视频测试 |
| 账号池 | 多账号轮询、单账号并发控制、对话/图片/视频分用途冷却记录 |
| 运维接口 | `/healthz`、`/readyz`、`/keepalive`、Docker Healthcheck、多架构镜像发布 |

### 2. 版本说明

| 版本 | 技术栈 | 状态 |
| --- | --- | --- |
| `qwen2API v1.0` | Python + FastAPI/Uvicorn | 上游旧版实现，仅保留历史说明 |
| `qwen2API v2.0` | Go 后端 + React WebUI | 本项目基于的上游主线 |
| `web2api` | Go 后端 + React WebUI + Playwright | 本项目，浏览器优先 |

## 二、快速部署

### 1. Docker 部署

大多数情况下，直接用 Docker 镜像就够了。把 `data` 和 `logs` 放在 compose 文件旁边，Docker 会把它们挂载进容器；以后更新镜像时，账号、Key 和日志都还在。

> [!NOTE]
> 如果 `welcomemonth/web2api` 镜像尚未发布，请改用 [本地 Docker 编译](#2-本地-docker-编译)。

```bash
mkdir web2api
cd web2api
mkdir -p data logs
```

创建一个简单的 `.env`：

```env
HOST_PORT=7860
HOST_DATA_DIR=./data
HOST_LOGS_DIR=./logs
ADMIN_KEY=replace-with-your-own-strong-random-key
```

创建 `docker-compose.yml`：

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

不用手动配置 `accounts.json`、`api_keys.json` 这类内部文件路径。镜像里默认就是 `/app/data` 和 `/app/logs`，上面的 volume 映射决定它们最终落在宿主机哪个目录。

拉取并启动：

```bash
docker compose pull
docker compose up -d
docker compose logs -f web2api
```

访问地址：

- WebUI：`http://127.0.0.1:7860`
- 健康检查：`http://127.0.0.1:7860/healthz`
- 保活探针：`http://127.0.0.1:7860/keepalive`

### 2. 本地 Docker 编译

当你修改了源码并需要构建自己的镜像时，使用本地 Docker 编译方式。

```bash
git clone https://github.com/welcomemonth/web2api.git
cd web2api
cp .env.example .env
docker compose -f docker-compose.yml -f docker-compose.build.yml build
docker compose -f docker-compose.yml -f docker-compose.build.yml up -d
```

## 三、架构与配置

### 1. 运行架构

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

### 2. 环境变量

不要提交真实密钥。`.env.example` 只保留空值和注释示例。

| 变量 | 说明 |
| --- | --- |
| `ADMIN_KEY` | WebUI 和 `/api/admin/*` 管理接口 Key，请自行设置强随机值 |
| `QWEN_API_KEY`、`QWEN_API_KEYS`、`QWEN_API_KEY_N` | 环境变量注入的下游 API Key，仅运行时存在，不写入 `data/api_keys.json`，不能从 WebUI 删除 |
| `QWEN_ACCOUNT_N` | 环境变量注入的上游账号，格式为 `token;optional-email;optional-password`，不写入 `data/accounts.json` |
| `BROWSER_POOL_SIZE` | 浏览器池中的无头浏览器实例数量，默认 `1` |
| `BROWSER_STREAM_TIMEOUT_SECONDS` | 浏览器驱动流式请求的超时时间，默认 `1800` |
| `MAX_INFLIGHT_PER_ACCOUNT` | 单个千问账号的最大并发请求数，默认 `2` |
| `CHAT_ID_PREWARM_TARGET_PER_ACCOUNT` | 每个账号预创建的千问 Web 会话数量，默认 `5` |
| `CHAT_ID_PREWARM_TTL_SECONDS` / `CHAT_ID_PREWARM_MAX_CONCURRENCY` | `Chat_ID` 预热池的 TTL 与最大并发，默认 `120` / `16` |
| `KEEPALIVE_URL`、`KEEPALIVE_INTERVAL` | 可选后台保活任务；环境变量存在时会锁定 WebUI 中对应配置 |
| `TOOL_RECOVERY_MAX_ATTEMPTS` | 工具结果之后上游没有产出下一次客户端工具调用时的自动恢复尝试次数；默认 `4`，限制在 `1`-`8` |
| `HOST_DATA_DIR`、`HOST_LOGS_DIR` | Docker 宿主机挂载路径，默认 `./data` 和 `./logs` |
| `DATA_DIR`、`LOGS_DIR` | 本地非 Docker 路径覆盖；留空时使用当前项目目录 |

## 四、开发指南

### 1. 环境要求

- **Go**: `1.26`
- **Node.js**: `20+`
- **npm**: 随 Node.js 安装
- **Docker**: 仅在需要容器构建或容器部署时使用
- **Playwright 浏览器**: 使用 `--install-browsers` 安装，或 `go run start-all.go --install-browsers`

### 2. 一键本地启动

```bash
# 1. 克隆项目
git clone https://github.com/welcomemonth/web2api.git
cd web2api

# 2. 安装前端依赖
cd frontend
npm ci
cd ..

# 3. 一键启动 Go 后端和 React WebUI
go run start-all.go
```

### 3. 启动与验证

1. 访问 WebUI：`http://127.0.0.1:7860`
2. 检查健康状态：`http://127.0.0.1:7860/healthz`
3. 检查保活探针：`http://127.0.0.1:7860/keepalive`
4. 修改后端代码后执行：

```powershell
cd backend
go test ./...
go build -trimpath -ldflags="-s -w" -o ..\bin\web2api-backend.exe .
```

5. 修改前端代码后执行：

```powershell
cd frontend
npm run build
```

### 4. Docker 本地调试

```bash
# 使用本地源码构建镜像
docker compose -f docker-compose.yml -f docker-compose.build.yml build

# 使用本地镜像启动
docker compose -f docker-compose.yml -f docker-compose.build.yml up -d

# 查看日志
docker compose logs -f web2api
```

### 5. 开发规范

- 🧩 **模块化开发**: 协议适配、账号池、浏览器引擎、文件上下文、图片/视频能力和 WebUI 组件应保持边界清晰。
- 🔄 **提交流程**: 小步修改 → 本地验证 → 提交 PR，避免把无关清理、功能改动和格式化混在一起。
- 📚 **文档更新**: 新增用户可见配置、Docker 参数、接口能力或部署方式时，同步更新 README 和 `.env.example`。
- 🔐 **安全边界**: 不要提交 `data/`、`logs/`、`.env`、真实 token、Cookie、密码或下游 API Key。
- 🐳 **Docker 规则**: 容器内部数据路径固定为 `/app/data` 和 `/app/logs`，宿主机路径只通过 compose volume 映射控制。

## 五、参与贡献

### 1. 贡献方式

- 🐛 [报告 Bug](https://github.com/welcomemonth/web2api/issues)
- 💡 [功能建议](https://github.com/welcomemonth/web2api/issues)
- 🔧 [提交代码](https://github.com/welcomemonth/web2api/pulls)

### 2. 贡献列表

感谢所有帮助改进 web2api 的贡献者。

[![Contributors](https://contrib.rocks/image?repo=welcomemonth/web2api)](https://github.com/welcomemonth/web2api/graphs/contributors)

## 六、其他信息

### 1. Star History

[![Star History Chart](https://api.star-history.com/svg?repos=welcomemonth/web2api&type=Timeline)](https://www.star-history.com/#welcomemonth/web2api&Timeline)

### 2. 开源协议

本项目采用 [GPL-3.0 License](./LICENSE) 开源协议。

**您可以自由地：**

- ✅ 使用、复制、修改和分发本软件。
- ✅ 用于个人、自托管或商业场景。
- ✅ 在遵循 GPL-3.0 的前提下发布修改版本。

**但需要：**

- 📋 保留原始版权声明和开源协议。
- 📋 标明对原始代码的修改。
- 📋 分发衍生作品时遵循 GPL-3.0 的开源要求。

### 3. 免责声明

- 本项目为**免费开源**项目，不保证绝对完善无误。
- 本项目是自托管协议转换网关，请自行评估部署环境、账号风险和使用方式。
- 使用前请自行确认所在地法律、平台规则和上游账号政策。
- 严禁公开泄露真实账号 token、Cookie、密码或下游 API Key。
- 如发现安全问题，请优先通过私密渠道反馈，避免公开敏感信息。

### 4. 特别鸣谢

- 特别鸣谢: [qwen2API](https://github.com/YuJunZhiXue/qwen2API) — 本项目基于 qwen2API 二次开发，感谢 [YuJunZhiXue](https://github.com/YuJunZhiXue) 及所有 qwen2API 贡献者。
- 特别鸣谢: [LinuxDo](https://linux.do/)

---

<div align="center">
  <p>如果 web2api 对你有帮助，可以考虑给项目一个 Star。</p>
  <p>由 <a href="https://github.com/welcomemonth">welcomemonth</a> 维护，基于 <a href="https://github.com/YuJunZhiXue/qwen2API">qwen2API</a>（原作者 <a href="https://github.com/YuJunZhiXue">YuJunZhiXue</a>）。</p>
</div>
