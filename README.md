# Web2API

面向个人开发者的 Web 大模型 API 网关。在 Web 管理后台维护第三方大模型网页账号、创建 API Key，外部应用通过标准 OpenAI-compatible API 调用，系统内部完成鉴权、限流、账号池调度与流式响应转换。

第一阶段以 `chat.qwen.ai` 作为首个模型接入实例。

## 环境要求

- **Go**：1.23+（本项目在 1.23.6 上验证）
- **Node.js**：20+（本项目在 Node 24 上验证）
- **npm**：10+

## 开发模式启动

开发阶段前后端分离运行：

| 服务 | 地址 | 启动命令 |
| --- | --- | --- |
| 后端（Go/Gin） | http://localhost:8080 | `go run ./cmd/server` |
| 前端（React/Vite） | http://localhost:5173 | `cd web && npm install && npm run dev` |

前端开发服务器已将 `/admin/*`、`/v1/*` 代理到后端 `localhost:8080`，开发阶段无需额外配置跨域。

## 默认管理员账号

首次启动后端时，若 `data/admin.json` 不存在，系统会自动创建默认管理员：

| 项 | 值 |
| --- | --- |
| 用户名 | `admin` |
| 密码 | `admin123` |

> 第一版面向本地自用，管理员密码**明文存储**，默认密码仅适合本地开发，请勿直接用于公网生产环境。
