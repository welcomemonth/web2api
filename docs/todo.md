# Web2API MVP 开发 Todo

> 基于 Web2API PRD v1.0 制定。
>
> **技术栈**：Go + Gin + Playwright-Go + React + Vite + JSON 文件存储  
> **开发模式**：React `localhost:5173`，Go `localhost:8080`  
> **最终模式**：Go 单体进程托管 React 静态资源，并提供 `/admin/*` 与 `/v1/*`  
> **第一版原则**：先完成核心链路，不引入数据库、Docker、Kubernetes、复杂监控、多租户等额外复杂度。
>
> 核心链路：
>
> `管理后台 → Qwen 账号 → 账号池 → API Key → 限流 → 调度 → Qwen Provider → OpenAI-compatible SSE`

---

## 0. 任务约定

- [x] **TASK-000** 建立项目任务执行规则；**依赖：无**
  - [x] **TASK-000-S01** 所有任务完成后再勾选对应任务；**依赖：TASK-000**
  - [x] **TASK-000-S02** 每完成一个阶段至少执行一次完整测试；**依赖：TASK-000**
  - [x] **TASK-000-S03** 不在 MVP 中临时加入 PRD 未定义的大功能；**依赖：TASK-000**
  - [x] **TASK-000-S04** 如果实现过程中发现需求冲突，优先保持现有 API 链路和数据模型稳定；**依赖：TASK-000**

---

# 1. 项目初始化与开发环境

## 1.1 项目骨架

- [x] **TASK-001** 创建 Web2API 项目根目录和基础工程结构；**依赖：TASK-000**
  - [x] **TASK-001-S01** 创建 `cmd/server` 目录并准备 `main.go`；**依赖：TASK-001**
  - [x] **TASK-001-S02** 创建 `internal/handler`、`service`、`repository`、`middleware`、`model`、`scheduler`、`provider`、`storage` 目录；**依赖：TASK-001**
  - [x] **TASK-001-S03** 创建 `web` 前端工程目录；**依赖：TASK-001**
  - [x] **TASK-001-S04** 创建 `data` 数据目录；**依赖：TASK-001**
  - [x] **TASK-001-S05** 创建 `.gitignore`，忽略运行数据、构建产物、浏览器缓存和本地敏感文件；**依赖：TASK-001**

- [x] **TASK-002** 初始化 Go Module 和后端依赖；**依赖：TASK-001**
  - [x] **TASK-002-S01** 初始化 `go.mod`；**依赖：TASK-002**
  - [x] **TASK-002-S02** 添加 Gin 依赖；**依赖：TASK-002**
  - [x] **TASK-002-S03** 添加 Playwright-Go 依赖；**依赖：TASK-002**
  - [x] **TASK-002-S04** 添加 JWT 实现所需依赖；**依赖：TASK-002**
  - [x] **TASK-002-S05** 执行 `go mod tidy` 并确认依赖可正常解析；**依赖：TASK-002**

- [x] **TASK-003** 初始化 React + Vite 前端工程；**依赖：TASK-001**
  - [x] **TASK-003-S01** 初始化 React + Vite 项目；**依赖：TASK-003**
  - [x] **TASK-003-S02** 配置开发服务器端口为 `5173`；**依赖：TASK-003**
  - [x] **TASK-003-S03** 配置前端请求开发阶段指向 Go `localhost:8080`；**依赖：TASK-003**
  - [x] **TASK-003-S04** 删除无关示例页面和默认模板代码；**依赖：TASK-003**

## 1.2 最小运行闭环

- [x] **TASK-004** 实现最小 Go HTTP 服务；**依赖：TASK-002**
  - [x] **TASK-004-S01** 创建 Gin Engine；**依赖：TASK-004**
  - [x] **TASK-004-S02** 配置 `localhost:8080` 监听；**依赖：TASK-004**
  - [x] **TASK-004-S03** 添加 `/health` 健康检查接口；**依赖：TASK-004**
  - [x] **TASK-004-S04** 验证服务启动和健康检查正常；**依赖：TASK-004**

- [ ] **TASK-005** 完善基础 README 开发说明；**依赖：TASK-004、TASK-003**
  - [ ] **TASK-005-S01** 记录 Go 和 Node 环境要求；**依赖：TASK-005**
  - [ ] **TASK-005-S02** 记录后端启动命令；**依赖：TASK-005**
  - [ ] **TASK-005-S03** 记录前端启动命令；**依赖：TASK-005**
  - [ ] **TASK-005-S04** 记录默认开发端口；**依赖：TASK-005**
  - [ ] **TASK-005-S05** 记录 MVP 默认管理员账号规则；**依赖：TASK-005**

---

# 2. 配置、数据模型与 JSON Repository

## 2.1 配置系统

- [ ] **TASK-006** 建立后端基础配置模型；**依赖：TASK-004**
  - [ ] **TASK-006-S01** 定义 HTTP 服务端口配置；**依赖：TASK-006**
  - [ ] **TASK-006-S02** 定义数据目录配置；**依赖：TASK-006**
  - [ ] **TASK-006-S03** 定义 JWT 有效期配置，默认 24 小时；**依赖：TASK-006**
  - [ ] **TASK-006-S04** 定义 API Key 固定 10 RPM 配置；**依赖：TASK-006**
  - [ ] **TASK-006-S05** 定义账号池最大等待时间 30 秒；**依赖：TASK-006**

## 2.2 数据模型

- [ ] **TASK-007** 定义管理员数据模型；**依赖：TASK-006**
  - [ ] **TASK-007-S01** 定义管理员用户名字段；**依赖：TASK-007**
  - [ ] **TASK-007-S02** 定义管理员密码字段；**依赖：TASK-007**
  - [ ] **TASK-007-S03** 保持第一版明文存储约定；**依赖：TASK-007**

- [ ] **TASK-008** 定义 Qwen Account 数据模型；**依赖：TASK-006**
  - [ ] **TASK-008-S01** 定义账号 ID、Username、Password、Access Token、创建时间等持久化字段；**依赖：TASK-008**
  - [ ] **TASK-008-S02** 定义持久化状态 `normal/error/disabled`；**依赖：TASK-008**
  - [ ] **TASK-008-S03** 明确 `busy/idle` 为运行时状态，不写入 JSON；**依赖：TASK-008**
  - [ ] **TASK-008-S04** 为运行时对象预留 Browser Context、busy 状态和当前请求 ID；**依赖：TASK-008**

- [ ] **TASK-009** 定义 API Key 数据模型；**依赖：TASK-006**
  - [ ] **TASK-009-S01** 定义 Key ID、名称、Key、状态、创建时间字段；**依赖：TASK-009**
  - [ ] **TASK-009-S02** 定义启用和禁用状态；**依赖：TASK-009**
  - [ ] **TASK-009-S03** 明确第一版 API Key 明文存储；**依赖：TASK-009**

- [ ] **TASK-010** 定义 API Log 数据模型；**依赖：TASK-006**
  - [ ] **TASK-010-S01** 定义时间字段；**依赖：TASK-010**
  - [ ] **TASK-010-S02** 定义 API Key、Model、Qwen Account、耗时、结果、错误信息字段；**依赖：TASK-010**
  - [ ] **TASK-010-S03** 明确不记录 messages 和用户实际对话内容；**依赖：TASK-010**
  - [ ] **TASK-010-S04** 明确最多保留最近 1000 条日志；**依赖：TASK-010**

## 2.3 JSON 存储基础设施

- [ ] **TASK-011** 实现通用 JSON 文件存储组件；**依赖：TASK-007、TASK-008、TASK-009、TASK-010**
  - [ ] **TASK-011-S01** 实现 JSON 文件读取；**依赖：TASK-011**
  - [ ] **TASK-011-S02** 实现 JSON 文件写入；**依赖：TASK-011**
  - [ ] **TASK-011-S03** 为写操作增加 Mutex 保护；**依赖：TASK-011**
  - [ ] **TASK-011-S04** 采用“写临时文件 → Rename”方式完成原子替换；**依赖：TASK-011**
  - [ ] **TASK-011-S05** 处理文件不存在时的初始化场景；**依赖：TASK-011**
  - [ ] **TASK-011-S06** 处理 JSON 损坏和读写错误；**依赖：TASK-011**

- [ ] **TASK-012** 实现 Admin Repository；**依赖：TASK-011、TASK-007**
  - [ ] **TASK-012-S01** 支持读取管理员；**依赖：TASK-012**
  - [ ] **TASK-012-S02** 支持保存管理员；**依赖：TASK-012**
  - [ ] **TASK-012-S03** 支持首次启动自动创建管理员；**依赖：TASK-012**

- [ ] **TASK-013** 实现 Account Repository；**依赖：TASK-011、TASK-008**
  - [ ] **TASK-013-S01** 支持加载全部账号到内存；**依赖：TASK-013**
  - [ ] **TASK-013-S02** 支持新增账号；**依赖：TASK-013**
  - [ ] **TASK-013-S03** 支持更新账号；**依赖：TASK-013**
  - [ ] **TASK-013-S04** 支持删除账号；**依赖：TASK-013**

- [ ] **TASK-014** 实现 API Key Repository；**依赖：TASK-011、TASK-009**
  - [ ] **TASK-014-S01** 支持加载全部 API Key 到内存；**依赖：TASK-014**
  - [ ] **TASK-014-S02** 支持新增 API Key；**依赖：TASK-014**
  - [ ] **TASK-014-S03** 支持更新 API Key 状态；**依赖：TASK-014**
  - [ ] **TASK-014-S04** 支持删除 API Key；**依赖：TASK-014**

- [ ] **TASK-015** 实现 API Log Repository；**依赖：TASK-011、TASK-010**
  - [ ] **TASK-015-S01** 支持追加 API 调用日志；**依赖：TASK-015**
  - [ ] **TASK-015-S02** 写入后最多保留 1000 条；**依赖：TASK-015**
  - [ ] **TASK-015-S03** 超过上限时删除最旧记录；**依赖：TASK-015**

- [ ] **TASK-016** 实现应用启动时的数据加载；**依赖：TASK-012、TASK-013、TASK-014、TASK-015**
  - [ ] **TASK-016-S01** 启动时加载 `admin.json`；**依赖：TASK-016**
  - [ ] **TASK-016-S02** 启动时加载 `accounts.json`；**依赖：TASK-016**
  - [ ] **TASK-016-S03** 启动时加载 `api_keys.json`；**依赖：TASK-016**
  - [ ] **TASK-016-S04** 启动时加载 `api_logs.json`；**依赖：TASK-016**
  - [ ] **TASK-016-S05** 确认业务层运行期间优先访问内存数据；**依赖：TASK-016**

---

# 3. 管理员认证与 JWT

## 3.1 管理员初始化

- [ ] **TASK-017** 实现默认管理员初始化服务；**依赖：TASK-016**
  - [ ] **TASK-017-S01** 检测 `admin.json` 是否存在；**依赖：TASK-017**
  - [ ] **TASK-017-S02** 不存在时创建用户名 `admin`；**依赖：TASK-017**
  - [ ] **TASK-017-S03** 设置第一版默认密码；**依赖：TASK-017**
  - [ ] **TASK-017-S04** 将默认密码写入 README；**依赖：TASK-017**

## 3.2 JWT

- [ ] **TASK-018** 实现 JWT Token 生成与校验；**依赖：TASK-017**
  - [ ] **TASK-018-S01** 定义 JWT Claims；**依赖：TASK-018**
  - [ ] **TASK-018-S02** 设置 24 小时有效期；**依赖：TASK-018**
  - [ ] **TASK-018-S03** 实现 Token 签发；**依赖：TASK-018**
  - [ ] **TASK-018-S04** 实现 Token 校验和过期判断；**依赖：TASK-018**
  - [ ] **TASK-018-S05** 第一版不实现 Refresh Token 和黑名单；**依赖：TASK-018**

- [ ] **TASK-019** 实现管理员认证 Service；**依赖：TASK-012、TASK-018**
  - [ ] **TASK-019-S01** 实现用户名密码校验；**依赖：TASK-019**
  - [ ] **TASK-019-S02** 登录成功后生成 JWT；**依赖：TASK-019**
  - [ ] **TASK-019-S03** 登录失败返回统一认证错误；**依赖：TASK-019**

- [ ] **TASK-020** 实现管理员 JWT Middleware；**依赖：TASK-018**
  - [ ] **TASK-020-S01** 从 `Authorization: Bearer <JWT>` 获取 Token；**依赖：TASK-020**
  - [ ] **TASK-020-S02** 校验 JWT；**依赖：TASK-020**
  - [ ] **TASK-020-S03** 无 Token 返回 401；**依赖：TASK-020**
  - [ ] **TASK-020-S04** Token 无效或过期返回 401；**依赖：TASK-020**

- [ ] **TASK-021** 实现管理员登录 API；**依赖：TASK-019**
  - [ ] **TASK-021-S01** 注册 `POST /admin/login`；**依赖：TASK-021**
  - [ ] **TASK-021-S02** 校验用户名和密码参数；**依赖：TASK-021**
  - [ ] **TASK-021-S03** 成功返回 JWT；**依赖：TASK-021**
  - [ ] **TASK-021-S04** 测试正确密码、错误密码和缺失参数；**依赖：TASK-021**

---

# 4. Qwen Browser / Playwright 基础设施

## 4.1 Browser 生命周期

- [ ] **TASK-022** 实现全局 Playwright Browser 生命周期管理；**依赖：TASK-002**
  - [ ] **TASK-022-S01** 创建 Playwright 实例；**依赖：TASK-022**
  - [ ] **TASK-022-S02** 启动单个 Browser；**依赖：TASK-022**
  - [ ] **TASK-022-S03** 应用退出时关闭 Browser；**依赖：TASK-022**
  - [ ] **TASK-022-S04** 处理 Browser 启动失败；**依赖：TASK-022**

- [ ] **TASK-023** 实现 Qwen Account Browser Context 管理器；**依赖：TASK-022、TASK-008**
  - [ ] **TASK-023-S01** 为每个账号创建独立 Browser Context；**依赖：TASK-023**
  - [ ] **TASK-023-S02** 保存账号对应 Context 的运行时引用；**依赖：TASK-023**
  - [ ] **TASK-023-S03** 实现 Context 关闭；**依赖：TASK-023**
  - [ ] **TASK-023-S04** 确保不同账号 Context 相互隔离；**依赖：TASK-023**

## 4.2 Qwen 登录与验证

- [ ] **TASK-024** 定义 Qwen 账号登录/验证接口；**依赖：TASK-023**
  - [ ] **TASK-024-S01** 定义 Token 登录路径；**依赖：TASK-024**
  - [ ] **TASK-024-S02** 定义 Username + Password 登录路径；**依赖：TASK-024**
  - [ ] **TASK-024-S03** 定义账号验证结果；**依赖：TASK-024**
  - [ ] **TASK-024-S04** 明确验证成功对应 `normal`，失败对应 `error`；**依赖：TASK-024**

- [ ] **TASK-025** 实现 Access Token 优先的 Qwen 登录流程；**依赖：TASK-024**
  - [ ] **TASK-025-S01** 优先使用 Access Token 建立登录态；**依赖：TASK-025**
  - [ ] **TASK-025-S02** 验证 `chat.qwen.ai` 登录状态；**依赖：TASK-025**
  - [ ] **TASK-025-S03** Token 失效时进入密码恢复流程；**依赖：TASK-025**

- [ ] **TASK-026** 实现 Username + Password 登录流程；**依赖：TASK-024**
  - [ ] **TASK-026-S01** 建立登录页面访问流程；**依赖：TASK-026**
  - [ ] **TASK-026-S02** 填写用户名和密码；**依赖：TASK-026**
  - [ ] **TASK-026-S03** 判断登录是否成功；**依赖：TASK-026**
  - [ ] **TASK-026-S04** 登录失败返回明确错误；**依赖：TASK-026**

- [ ] **TASK-027** 实现账号凭证恢复策略；**依赖：TASK-025、TASK-026**
  - [ ] **TASK-027-S01** Token 失效时最多自动使用密码恢复 1 次；**依赖：TASK-027**
  - [ ] **TASK-027-S02** 恢复失败后将账号状态设置为 `error`；**依赖：TASK-027**
  - [ ] **TASK-027-S03** 恢复成功后恢复为 `normal`；**依赖：TASK-027**

---

# 5. Qwen Account Service 与后台账号管理

## 5.1 Account Service

- [ ] **TASK-028** 实现 Account Service 基础 CRUD；**依赖：TASK-013、TASK-023**
  - [ ] **TASK-028-S01** 实现创建账号；**依赖：TASK-028**
  - [ ] **TASK-028-S02** 实现查询账号列表；**依赖：TASK-028**
  - [ ] **TASK-028-S03** 实现编辑账号；**依赖：TASK-028**
  - [ ] **TASK-028-S04** 实现删除账号；**依赖：TASK-028**

- [ ] **TASK-029** 实现账号创建后的立即验证；**依赖：TASK-028、TASK-027**
  - [ ] **TASK-029-S01** 创建账号后建立 Browser Context；**依赖：TASK-029**
  - [ ] **TASK-029-S02** 执行 Qwen 登录/验证；**依赖：TASK-029**
  - [ ] **TASK-029-S03** 验证成功后保存 `normal`；**依赖：TASK-029**
  - [ ] **TASK-029-S04** 验证失败后保存 `error`；**依赖：TASK-029**

- [ ] **TASK-030** 实现账号编辑后的重新验证；**依赖：TASK-028、TASK-027**
  - [ ] **TASK-030-S01** 修改凭证前处理旧 Context；**依赖：TASK-030**
  - [ ] **TASK-030-S02** 根据新凭证创建新的 Context；**依赖：TASK-030**
  - [ ] **TASK-030-S03** 执行重新验证；**依赖：TASK-030**
  - [ ] **TASK-030-S04** 验证失败时设置 `error`；**依赖：TASK-030**

- [ ] **TASK-031** 实现账号禁用/启用；**依赖：TASK-028、TASK-023**
  - [ ] **TASK-031-S01** 禁用账号后禁止新请求分配；**依赖：TASK-031**
  - [ ] **TASK-031-S02** 空闲账号禁用时立即关闭 Context；**依赖：TASK-031**
  - [ ] **TASK-031-S03** 忙碌账号禁用时等待当前请求结束；**依赖：TASK-031**
  - [ ] **TASK-031-S04** 启用账号时重新创建 Context；**依赖：TASK-031**
  - [ ] **TASK-031-S05** 启用后重新验证并恢复 `normal` 或 `error`；**依赖：TASK-031**

- [ ] **TASK-032** 实现账号删除生命周期；**依赖：TASK-028、TASK-023**
  - [ ] **TASK-032-S01** 空闲账号直接关闭 Context 后删除；**依赖：TASK-032**
  - [ ] **TASK-032-S02** 忙碌账号标记待删除；**依赖：TASK-032**
  - [ ] **TASK-032-S03** 待删除账号禁止新请求进入；**依赖：TASK-032**
  - [ ] **TASK-032-S04** 当前请求结束后关闭 Context 并删除持久化数据；**依赖：TASK-032**

## 5.2 Account Handler

- [ ] **TASK-033** 实现管理员账号管理 API；**依赖：TASK-020、TASK-028**
  - [ ] **TASK-033-S01** 实现账号列表接口；**依赖：TASK-033**
  - [ ] **TASK-033-S02** 实现账号创建接口；**依赖：TASK-033**
  - [ ] **TASK-033-S03** 实现账号编辑接口；**依赖：TASK-033**
  - [ ] **TASK-033-S04** 实现启用接口；**依赖：TASK-033**
  - [ ] **TASK-033-S05** 实现禁用接口；**依赖：TASK-033**
  - [ ] **TASK-033-S06** 实现重新验证接口；**依赖：TASK-033**
  - [ ] **TASK-033-S07** 实现删除接口；**依赖：TASK-033**
  - [ ] **TASK-033-S08** 所有接口统一使用管理员 JWT Middleware；**依赖：TASK-033**

---

# 6. Qwen Account Pool 与调度器

## 6.1 调度器核心状态

- [ ] **TASK-034** 建立 Account Scheduler 核心数据结构；**依赖：TASK-028、TASK-031、TASK-032**
  - [ ] **TASK-034-S01** 管理可调度账号集合；**依赖：TASK-034**
  - [ ] **TASK-034-S02** 为账号维护 `busy/idle` 运行时状态；**依赖：TASK-034**
  - [ ] **TASK-034-S03** 增加 Round Robin 游标；**依赖：TASK-034**
  - [ ] **TASK-034-S04** 为并发访问增加 Mutex 保护；**依赖：TASK-034**

- [ ] **TASK-035** 实现 Round Robin 空闲账号选择；**依赖：TASK-034**
  - [ ] **TASK-035-S01** 只选择 `normal + idle` 账号；**依赖：TASK-035**
  - [ ] **TASK-035-S02** 跳过 `error` 账号；**依赖：TASK-035**
  - [ ] **TASK-035-S03** 跳过 `disabled` 账号；**依赖：TASK-035**
  - [ ] **TASK-035-S04** 跳过 `busy` 账号；**依赖：TASK-035**
  - [ ] **TASK-035-S05** 选择成功后更新 Round Robin 游标；**依赖：TASK-035**

## 6.2 FIFO 等待队列

- [ ] **TASK-036** 实现无空闲账号时的 FIFO 等待机制；**依赖：TASK-035**
  - [ ] **TASK-036-S01** 定义等待请求结构；**依赖：TASK-036**
  - [ ] **TASK-036-S02** 按进入顺序保存等待请求；**依赖：TASK-036**
  - [ ] **TASK-036-S03** 账号释放时优先唤醒队首请求；**依赖：TASK-036**
  - [ ] **TASK-036-S04** 为每个等待请求绑定 30 秒超时 Context；**依赖：TASK-036**
  - [ ] **TASK-036-S05** 超时后从等待队列安全移除；**依赖：TASK-036**
  - [ ] **TASK-036-S06** 超时返回“无可用账号”的业务错误；**依赖：TASK-036**

- [ ] **TASK-037** 实现账号租借与释放机制；**依赖：TASK-035、TASK-036**
  - [ ] **TASK-037-S01** 成功分配时将账号标记为 `busy`；**依赖：TASK-037**
  - [ ] **TASK-037-S02** 请求结束后释放账号；**依赖：TASK-037**
  - [ ] **TASK-037-S03** 释放账号时唤醒 FIFO 队列；**依赖：TASK-037**
  - [ ] **TASK-037-S04** 使用 `defer` 保证异常路径也能释放账号；**依赖：TASK-037**

- [ ] **TASK-038** 实现调度器与账号生命周期联动；**依赖：TASK-031、TASK-032、TASK-037**
  - [ ] **TASK-038-S01** 禁用账号后从可调度集合移除；**依赖：TASK-038**
  - [ ] **TASK-038-S02** 删除账号后从调度器移除；**依赖：TASK-038**
  - [ ] **TASK-038-S03** 重新启用且验证成功后重新加入调度器；**依赖：TASK-038**
  - [ ] **TASK-038-S04** 账号进入 `error` 后停止新请求分配；**依赖：TASK-038**

- [ ] **TASK-039** 为调度器增加并发安全测试；**依赖：TASK-037、TASK-038**
  - [ ] **TASK-039-S01** 测试多个请求同时获取不同账号；**依赖：TASK-039**
  - [ ] **TASK-039-S02** 测试同一账号不会被同时租借；**依赖：TASK-039**
  - [ ] **TASK-039-S03** 测试 FIFO 顺序；**依赖：TASK-039**
  - [ ] **TASK-039-S04** 测试 30 秒超时；**依赖：TASK-039**
  - [ ] **TASK-039-S05** 使用 Go race detector 检查调度器；**依赖：TASK-039**

---

# 7. API Key 管理与限流

## 7.1 API Key Service

- [ ] **TASK-040** 实现 API Key Service；**依赖：TASK-014、TASK-009**
  - [ ] **TASK-040-S01** 实现 `sk-` 前缀 Key 生成；**依赖：TASK-040**
  - [ ] **TASK-040-S02** 生成 32 位随机字符串；**依赖：TASK-040**
  - [ ] **TASK-040-S03** 检查 Key 唯一性；**依赖：TASK-040**
  - [ ] **TASK-040-S04** 实现创建 Key；**依赖：TASK-040**
  - [ ] **TASK-040-S05** 实现查询 Key；**依赖：TASK-040**
  - [ ] **TASK-040-S06** 实现启用/禁用 Key；**依赖：TASK-040**
  - [ ] **TASK-040-S07** 实现删除 Key；**依赖：TASK-040**

- [ ] **TASK-041** 实现 API Key 管理 API；**依赖：TASK-020、TASK-040**
  - [ ] **TASK-041-S01** 实现 Key 列表接口；**依赖：TASK-041**
  - [ ] **TASK-041-S02** 实现 Key 创建接口；**依赖：TASK-041**
  - [ ] **TASK-041-S03** 实现 Key 启用接口；**依赖：TASK-041**
  - [ ] **TASK-041-S04** 实现 Key 禁用接口；**依赖：TASK-041**
  - [ ] **TASK-041-S05** 实现 Key 删除接口；**依赖：TASK-041**
  - [ ] **TASK-041-S06** 所有接口使用管理员 JWT Middleware；**依赖：TASK-041**

## 7.2 API Key 鉴权 Middleware

- [ ] **TASK-042** 实现 API Key 鉴权 Middleware；**依赖：TASK-040**
  - [ ] **TASK-042-S01** 读取 `Authorization` Header；**依赖：TASK-042**
  - [ ] **TASK-042-S02** 校验 `Bearer` 格式；**依赖：TASK-042**
  - [ ] **TASK-042-S03** 根据 Key 查找 API Key；**依赖：TASK-042**
  - [ ] **TASK-042-S04** 检查 Key 是否启用；**依赖：TASK-042**
  - [ ] **TASK-042-S05** 无效 Key 返回 401；**依赖：TASK-042**
  - [ ] **TASK-042-S06** 禁用 Key 返回 401；**依赖：TASK-042**
  - [ ] **TASK-042-S07** 将 API Key 信息写入 Gin Context 供后续 Service 使用；**依赖：TASK-042**

## 7.3 固定窗口限流

- [ ] **TASK-043** 实现 API Key 级 Fixed Window 10 RPM 限流器；**依赖：TASK-042、TASK-006**
  - [ ] **TASK-043-S01** 为每个 API Key 维护当前窗口起始时间；**依赖：TASK-043**
  - [ ] **TASK-043-S02** 为每个 API Key 维护当前窗口请求数；**依赖：TASK-043**
  - [ ] **TASK-043-S03** 窗口切换时重置计数；**依赖：TASK-043**
  - [ ] **TASK-043-S04** 第 1～10 次请求允许通过；**依赖：TASK-043**
  - [ ] **TASK-043-S05** 第 11 次请求返回 429；**依赖：TASK-043**
  - [ ] **TASK-043-S06** 为并发访问增加 Mutex 或等价保护；**依赖：TASK-043**

---

# 8. OpenAI Chat API 数据结构与错误处理

## 8.1 Chat Model

- [ ] **TASK-044** 定义 OpenAI-compatible Chat 请求模型；**依赖：TASK-010**
  - [ ] **TASK-044-S01** 定义 `model`；**依赖：TASK-044**
  - [ ] **TASK-044-S02** 定义 `messages`；**依赖：TASK-044**
  - [ ] **TASK-044-S03** 定义 `stream`；**依赖：TASK-044**
  - [ ] **TASK-044-S04** 定义 message 的 `role/content`；**依赖：TASK-044**

- [ ] **TASK-045** 定义 OpenAI-compatible SSE Chunk 模型；**依赖：TASK-044**
  - [ ] **TASK-045-S01** 定义 `id`；**依赖：TASK-045**
  - [ ] **TASK-045-S02** 定义 `object`；**依赖：TASK-045**
  - [ ] **TASK-045-S03** 定义 `model`；**依赖：TASK-045**
  - [ ] **TASK-045-S04** 定义 choices、delta、content 字段；**依赖：TASK-045**
  - [ ] **TASK-045-S05** 定义 `[DONE]` 结束格式；**依赖：TASK-045**

## 8.2 参数校验

- [ ] **TASK-046** 实现 Chat 请求参数校验器；**依赖：TASK-044**
  - [ ] **TASK-046-S01** 校验请求 JSON 格式；**依赖：TASK-046**
  - [ ] **TASK-046-S02** 校验 `model` 存在；**依赖：TASK-046**
  - [ ] **TASK-046-S03** 校验 `messages` 非空且结构正确；**依赖：TASK-046**
  - [ ] **TASK-046-S04** 只允许 `system/user/assistant` role；**依赖：TASK-046**
  - [ ] **TASK-046-S05** 只允许 string 类型 content；**依赖：TASK-046**
  - [ ] **TASK-046-S06** 强制 `stream=true`；**依赖：TASK-046**
  - [ ] **TASK-046-S07** 不支持的 role/content 返回 400；**依赖：TASK-046**

## 8.3 统一错误格式

- [ ] **TASK-047** 实现统一 OpenAI-compatible Error Response；**依赖：TASK-046**
  - [ ] **TASK-047-S01** 定义 error message/type/code；**依赖：TASK-047**
  - [ ] **TASK-047-S02** 实现 400 错误；**依赖：TASK-047**
  - [ ] **TASK-047-S03** 实现 401 错误；**依赖：TASK-047**
  - [ ] **TASK-047-S04** 实现 429 错误；**依赖：TASK-047**
  - [ ] **TASK-047-S05** 实现 503 错误；**依赖：TASK-047**
  - [ ] **TASK-047-S06** 实现 500 错误；**依赖：TASK-047**
  - [ ] **TASK-047-S07** 实现 501 错误；**依赖：TASK-047**

---

# 9. Model Provider 抽象与 Qwen Provider

## 9.1 Provider 抽象

- [ ] **TASK-048** 定义 `ModelProvider` 接口；**依赖：TASK-044、TASK-045**
  - [ ] **TASK-048-S01** 定义 `Chat(ctx, req, onChunk)` 接口；**依赖：TASK-048**
  - [ ] **TASK-048-S02** 定义 Provider 统一错误返回方式；**依赖：TASK-048**
  - [ ] **TASK-048-S03** 确保 Service 层不依赖 Qwen 具体实现；**依赖：TASK-048**

## 9.2 Qwen Provider

- [ ] **TASK-049** 建立 Qwen Provider 基础结构；**依赖：TASK-048、TASK-025**
  - [ ] **TASK-049-S01** 创建 `provider/qwen` 实现目录；**依赖：TASK-049**
  - [ ] **TASK-049-S02** 注入账号 Browser Context；**依赖：TASK-049**
  - [ ] **TASK-049-S03** 实现 Qwen Chat 调用入口；**依赖：TASK-049**

- [ ] **TASK-050** 实现 Qwen 多轮对话重建流程；**依赖：TASK-049、TASK-044**
  - [ ] **TASK-050-S01** 根据请求 messages 创建当前请求所需的 Qwen 对话；**依赖：TASK-050**
  - [ ] **TASK-050-S02** 按消息顺序写入 system/user/assistant 内容；**依赖：TASK-050**
  - [ ] **TASK-050-S03** 确保一次 API 请求不依赖账号之前的 API 对话；**依赖：TASK-050**

- [ ] **TASK-051** 实现 Qwen 页面交互与请求执行；**依赖：TASK-050**
  - [ ] **TASK-051-S01** 定位 Qwen 输入区域；**依赖：TASK-051**
  - [ ] **TASK-051-S02** 提交当前请求；**依赖：TASK-051**
  - [ ] **TASK-051-S03** 监听 Qwen 流式响应数据；**依赖：TASK-051**
  - [ ] **TASK-051-S04** 将可识别的文本片段通过 `onChunk` 回调输出；**依赖：TASK-051**
  - [ ] **TASK-051-S05** 处理 Qwen 页面异常和请求失败；**依赖：TASK-051**

- [ ] **TASK-052** 实现 Qwen Provider 的 Context 取消；**依赖：TASK-051**
  - [ ] **TASK-052-S01** 监听 Go `ctx.Done()`；**依赖：TASK-052**
  - [ ] **TASK-052-S02** 请求取消时停止 Playwright 操作；**依赖：TASK-052**
  - [ ] **TASK-052-S03** 确保取消后账号最终能够释放；**依赖：TASK-052**

- [ ] **TASK-053** 实现 Qwen Provider 错误状态联动；**依赖：TASK-051、TASK-038**
  - [ ] **TASK-053-S01** Qwen 执行失败时返回 Provider 错误；**依赖：TASK-053**
  - [ ] **TASK-053-S02** 按 PRD 将异常账号标记为 `error`；**依赖：TASK-053**
  - [ ] **TASK-053-S03** 确保异常账号不再被新请求调度；**依赖：TASK-053**

---

# 10. OpenAI-compatible Chat API 核心链路

## 10.1 Chat Service

- [ ] **TASK-054** 实现 Chat Service；**依赖：TASK-043、TASK-046、TASK-047、TASK-037、TASK-048**
  - [ ] **TASK-054-S01** 接收已通过 API Key 鉴权的请求；**依赖：TASK-054**
  - [ ] **TASK-054-S02** 执行 Chat 参数校验；**依赖：TASK-054**
  - [ ] **TASK-054-S03** 从账号池获取可用 Qwen 账号；**依赖：TASK-054**
  - [ ] **TASK-054-S04** 调用 ModelProvider；**依赖：TASK-054**
  - [ ] **TASK-054-S05** 通过回调处理 Provider 流式 chunk；**依赖：TASK-054**
  - [ ] **TASK-054-S06** 使用 `defer` 确保账号释放；**依赖：TASK-054**
  - [ ] **TASK-054-S07** 将客户端取消传递到 Provider；**依赖：TASK-054**

## 10.2 SSE 输出

- [ ] **TASK-055** 实现 OpenAI-compatible SSE Writer；**依赖：TASK-045、TASK-054**
  - [ ] **TASK-055-S01** 设置 `Content-Type: text/event-stream`；**依赖：TASK-055**
  - [ ] **TASK-055-S02** 设置必要的缓存/连接 Header；**依赖：TASK-055**
  - [ ] **TASK-055-S03** 将 Provider chunk 转换成 OpenAI SSE chunk；**依赖：TASK-055**
  - [ ] **TASK-055-S04** 每个 chunk 写入 `data: ...`；**依赖：TASK-055**
  - [ ] **TASK-055-S05** 每次写入后 Flush；**依赖：TASK-055**
  - [ ] **TASK-055-S06** 正常完成后输出 `data: [DONE]`；**依赖：TASK-055**

## 10.3 Chat Handler

- [ ] **TASK-056** 实现 `POST /v1/chat/completions`；**依赖：TASK-042、TASK-043、TASK-054、TASK-055**
  - [ ] **TASK-056-S01** 注册 `/v1/chat/completions` 路由；**依赖：TASK-056**
  - [ ] **TASK-056-S02** 绑定 API Key Middleware；**依赖：TASK-056**
  - [ ] **TASK-056-S03** 绑定 10 RPM Middleware；**依赖：TASK-056**
  - [ ] **TASK-056-S04** 解析请求体；**依赖：TASK-056**
  - [ ] **TASK-056-S05** 调用 Chat Service；**依赖：TASK-056**
  - [ ] **TASK-056-S06** 将 Service 错误转换为统一 Error Response；**依赖：TASK-056**
  - [ ] **TASK-056-S07** 将正常结果以 SSE 输出；**依赖：TASK-056**

## 10.4 Anthropic 占位接口

- [ ] **TASK-057** 实现 `POST /v1/messages` 占位接口；**依赖：TASK-047**
  - [ ] **TASK-057-S01** 注册 `/v1/messages`；**依赖：TASK-057**
  - [ ] **TASK-057-S02** 返回 HTTP 501；**依赖：TASK-057**
  - [ ] **TASK-057-S03** 返回统一错误格式；**依赖：TASK-057**

---

# 11. API 调用日志

- [ ] **TASK-058** 将 API 日志接入 Chat 核心链路；**依赖：TASK-015、TASK-054**
  - [ ] **TASK-058-S01** 请求开始时记录开始时间；**依赖：TASK-058**
  - [ ] **TASK-058-S02** 记录 API Key；**依赖：TASK-058**
  - [ ] **TASK-058-S03** 记录请求 Model；**依赖：TASK-058**
  - [ ] **TASK-058-S04** 请求获得账号后记录 Qwen Account；**依赖：TASK-058**
  - [ ] **TASK-058-S05** 请求结束时计算耗时；**依赖：TASK-058**
  - [ ] **TASK-058-S06** 记录 success/failure；**依赖：TASK-058**
  - [ ] **TASK-058-S07** 失败时记录错误信息；**依赖：TASK-058**
  - [ ] **TASK-058-S08** 确认 messages 和实际对话内容不会进入日志；**依赖：TASK-058**

- [ ] **TASK-059** 验证 API Log 1000 条保留策略；**依赖：TASK-058**
  - [ ] **TASK-059-S01** 构造超过 1000 条日志；**依赖：TASK-059**
  - [ ] **TASK-059-S02** 确认只保留最近 1000 条；**依赖：TASK-059**
  - [ ] **TASK-059-S03** 确认最旧记录被删除；**依赖：TASK-059**

---

# 12. React 管理后台基础设施

## 12.1 前端基础结构

- [ ] **TASK-060** 建立 React 路由和页面基础结构；**依赖：TASK-003、TASK-021**
  - [ ] **TASK-060-S01** 建立 `/login`；**依赖：TASK-060**
  - [ ] **TASK-060-S02** 建立 `/dashboard`；**依赖：TASK-060**
  - [ ] **TASK-060-S03** 建立 `/accounts`；**依赖：TASK-060**
  - [ ] **TASK-060-S04** 建立 `/api-keys`；**依赖：TASK-060**
  - [ ] **TASK-060-S05** 建立 `/chat-test`；**依赖：TASK-060**
  - [ ] **TASK-060-S06** 增加登录态路由保护；**依赖：TASK-060**

- [ ] **TASK-061** 建立前端 HTTP API 封装；**依赖：TASK-060**
  - [ ] **TASK-061-S01** 统一封装 API Base URL；**依赖：TASK-061**
  - [ ] **TASK-061-S02** 自动添加 JWT Authorization Header；**依赖：TASK-061**
  - [ ] **TASK-061-S03** 统一处理 401 并回到登录页；**依赖：TASK-061**
  - [ ] **TASK-061-S04** 封装统一错误处理；**依赖：TASK-061**

## 12.2 Login 页面

- [ ] **TASK-062** 实现管理员 Login 页面；**依赖：TASK-060、TASK-061**
  - [ ] **TASK-062-S01** 创建用户名输入框；**依赖：TASK-062**
  - [ ] **TASK-062-S02** 创建密码输入框；**依赖：TASK-062**
  - [ ] **TASK-062-S03** 创建登录按钮；**依赖：TASK-062**
  - [ ] **TASK-062-S04** 调用 `/admin/login`；**依赖：TASK-062**
  - [ ] **TASK-062-S05** 保存 JWT；**依赖：TASK-062**
  - [ ] **TASK-062-S06** 登录成功跳转 Dashboard；**依赖：TASK-062**
  - [ ] **TASK-062-S07** 登录失败显示错误信息；**依赖：TASK-062**

---

# 13. Dashboard

- [ ] **TASK-063** 实现 Dashboard 后端统计接口；**依赖：TASK-033、TASK-041、TASK-058**
  - [ ] **TASK-063-S01** 返回 Qwen 账号总数；**依赖：TASK-063**
  - [ ] **TASK-063-S02** 返回正常/忙碌/异常/禁用数量；**依赖：TASK-063**
  - [ ] **TASK-063-S03** 返回执行中请求数量；**依赖：TASK-063**
  - [ ] **TASK-063-S04** 返回等待中请求数量；**依赖：TASK-063**
  - [ ] **TASK-063-S05** 返回 API 总请求/成功/失败数量；**依赖：TASK-063**
  - [ ] **TASK-063-S06** 不实现图表和高级监控；**依赖：TASK-063**

- [ ] **TASK-064** 实现 Dashboard 页面；**依赖：TASK-063、TASK-061**
  - [ ] **TASK-064-S01** 展示 Qwen 账号总数；**依赖：TASK-064**
  - [ ] **TASK-064-S02** 展示正常/忙碌/异常/禁用状态；**依赖：TASK-064**
  - [ ] **TASK-064-S03** 展示当前执行/等待请求；**依赖：TASK-064**
  - [ ] **TASK-064-S04** 展示 API 总请求/成功/失败；**依赖：TASK-064**
  - [ ] **TASK-064-S05** 验证空数据状态下页面仍可正常显示；**依赖：TASK-064**

---

# 14. Qwen Account 管理页面

- [ ] **TASK-065** 实现 Accounts 页面列表；**依赖：TASK-033、TASK-061**
  - [ ] **TASK-065-S01** 展示 Username；**依赖：TASK-065**
  - [ ] **TASK-065-S02** 展示账号状态；**依赖：TASK-065**
  - [ ] **TASK-065-S03** 展示创建时间；**依赖：TASK-065**
  - [ ] **TASK-065-S04** 提供添加、编辑、启用、禁用、验证、删除操作；**依赖：TASK-065**

- [ ] **TASK-066** 实现添加 Qwen Account 表单；**依赖：TASK-065**
  - [ ] **TASK-066-S01** 添加 Username 字段；**依赖：TASK-066**
  - [ ] **TASK-066-S02** 添加 Password 字段；**依赖：TASK-066**
  - [ ] **TASK-066-S03** 添加 Access Token 字段；**依赖：TASK-066**
  - [ ] **TASK-066-S04** 前端校验 Username 必填；**依赖：TASK-066**
  - [ ] **TASK-066-S05** 前端校验 Password/Token 至少填写一个；**依赖：TASK-066**
  - [ ] **TASK-066-S06** 提交后显示账号验证状态；**依赖：TASK-066**

- [ ] **TASK-067** 实现 Account 编辑和状态操作；**依赖：TASK-065**
  - [ ] **TASK-067-S01** 实现编辑表单；**依赖：TASK-067**
  - [ ] **TASK-067-S02** 实现启用操作；**依赖：TASK-067**
  - [ ] **TASK-067-S03** 实现禁用操作；**依赖：TASK-067**
  - [ ] **TASK-067-S04** 实现重新验证操作；**依赖：TASK-067**
  - [ ] **TASK-067-S05** 实现删除操作并增加确认；**依赖：TASK-067**

---

# 15. API Key 管理页面

- [ ] **TASK-068** 实现 API Keys 页面；**依赖：TASK-041、TASK-061**
  - [ ] **TASK-068-S01** 展示名称；**依赖：TASK-068**
  - [ ] **TASK-068-S02** 展示 Key；**依赖：TASK-068**
  - [ ] **TASK-068-S03** 展示状态；**依赖：TASK-068**
  - [ ] **TASK-068-S04** 展示创建时间；**依赖：TASK-068**
  - [ ] **TASK-068-S05** 提供创建、启用、禁用、删除操作；**依赖：TASK-068**

- [ ] **TASK-069** 实现 API Key 创建流程；**依赖：TASK-068**
  - [ ] **TASK-069-S01** 创建 Key 名称输入框；**依赖：TASK-069**
  - [ ] **TASK-069-S02** 调用创建 API；**依赖：TASK-069**
  - [ ] **TASK-069-S03** 展示系统生成的 `sk-...` Key；**依赖：TASK-069**
  - [ ] **TASK-069-S04** 验证创建后 Key 可以进入启用状态；**依赖：TASK-069**

- [ ] **TASK-070** 实现 API Key 状态和删除操作；**依赖：TASK-068**
  - [ ] **TASK-070-S01** 实现启用；**依赖：TASK-070**
  - [ ] **TASK-070-S02** 实现禁用；**依赖：TASK-070**
  - [ ] **TASK-070-S03** 实现删除；**依赖：TASK-070**
  - [ ] **TASK-070-S04** 删除前增加确认；**依赖：TASK-070**

---

# 16. Chat Test 页面

## 16.1 页面

- [ ] **TASK-071** 实现 Chat Test 页面基础 UI；**依赖：TASK-062、TASK-068**
  - [ ] **TASK-071-S01** 增加 API Key 下拉选择；**依赖：TASK-071**
  - [ ] **TASK-071-S02** 增加 Model 输入框；**依赖：TASK-071**
  - [ ] **TASK-071-S03** 增加消息展示区域；**依赖：TASK-071**
  - [ ] **TASK-071-S04** 增加问题输入框；**依赖：TASK-071**
  - [ ] **TASK-071-S05** 增加 Send 按钮；**依赖：TASK-071**

## 16.2 SSE 客户端

- [ ] **TASK-072** 实现 Chat Test 的 SSE 流式调用；**依赖：TASK-056、TASK-071**
  - [ ] **TASK-072-S01** 使用实际 `/v1/chat/completions` API；**依赖：TASK-072**
  - [ ] **TASK-072-S02** 使用选中的 API Key；**依赖：TASK-072**
  - [ ] **TASK-072-S03** 设置 `stream=true`；**依赖：TASK-072**
  - [ ] **TASK-072-S04** 解析 `data: ...` SSE 数据；**依赖：TASK-072**
  - [ ] **TASK-072-S05** 识别 `[DONE]`；**依赖：TASK-072**
  - [ ] **TASK-072-S06** 将文本 chunk 实时追加到回答区域；**依赖：TASK-072**
  - [ ] **TASK-072-S07** 请求取消/页面离开时正确关闭流；**依赖：TASK-072**

- [ ] **TASK-073** 确保 Chat Test 与真实第三方 API 使用完全相同业务链路；**依赖：TASK-072、TASK-054**
  - [ ] **TASK-073-S01** 不在 Chat Test 中绕过 API Key 鉴权；**依赖：TASK-073**
  - [ ] **TASK-073-S02** 不在 Chat Test 中绕过 10 RPM 限流；**依赖：TASK-073**
  - [ ] **TASK-073-S03** 不在 Chat Test 中绕过账号池；**依赖：TASK-073**
  - [ ] **TASK-073-S04** 不创建独立 Qwen 测试链路；**依赖：TASK-073**

---

# 17. 开发阶段联调

## 17.1 后端单模块联调

- [ ] **TASK-074** 联调管理员登录；**依赖：TASK-021、TASK-062**
  - [ ] **TASK-074-S01** 启动 Go 服务；**依赖：TASK-074**
  - [ ] **TASK-074-S02** 启动 React 服务；**依赖：TASK-074**
  - [ ] **TASK-074-S03** 使用默认管理员登录；**依赖：TASK-074**
  - [ ] **TASK-074-S04** 验证 JWT 后台接口访问成功；**依赖：TASK-074**
  - [ ] **TASK-074-S05** 验证错误密码无法登录；**依赖：TASK-074**

- [ ] **TASK-075** 联调 Qwen Account 管理；**依赖：TASK-033、TASK-066**
  - [ ] **TASK-075-S01** 创建一个有效 Qwen Account；**依赖：TASK-075**
  - [ ] **TASK-075-S02** 验证账号创建后立即验证；**依赖：TASK-075**
  - [ ] **TASK-075-S03** 验证正常账号进入 `normal`；**依赖：TASK-075**
  - [ ] **TASK-075-S04** 验证无效账号进入 `error`；**依赖：TASK-075**
  - [ ] **TASK-075-S05** 验证重新验证能够恢复账号；**依赖：TASK-075**

- [ ] **TASK-076** 联调 API Key 管理；**依赖：TASK-041、TASK-069**
  - [ ] **TASK-076-S01** 创建 API Key；**依赖：TASK-076**
  - [ ] **TASK-076-S02** 验证 Key 格式正确；**依赖：TASK-076**
  - [ ] **TASK-076-S03** 验证启用 Key 可用于 API；**依赖：TASK-076**
  - [ ] **TASK-076-S04** 验证禁用 Key 返回 401；**依赖：TASK-076**
  - [ ] **TASK-076-S05** 验证删除 Key 后不可用；**依赖：TASK-076**

## 17.2 核心 API 链路联调

- [ ] **TASK-077** 完成 `/v1/chat/completions` 单账号端到端联调；**依赖：TASK-056、TASK-075、TASK-076**
  - [ ] **TASK-077-S01** 使用 curl 发起 OpenAI-compatible 请求；**依赖：TASK-077**
  - [ ] **TASK-077-S02** 验证 Authorization Bearer API Key；**依赖：TASK-077**
  - [ ] **TASK-077-S03** 验证请求经过 10 RPM 限流；**依赖：TASK-077**
  - [ ] **TASK-077-S04** 验证账号池分配 Qwen Account；**依赖：TASK-077**
  - [ ] **TASK-077-S05** 验证 Qwen Provider 正常执行；**依赖：TASK-077**
  - [ ] **TASK-077-S06** 验证 OpenAI SSE chunk 持续输出；**依赖：TASK-077**
  - [ ] **TASK-077-S07** 验证最终输出 `data: [DONE]`；**依赖：TASK-077**
  - [ ] **TASK-077-S08** 验证 API Log 正确记录；**依赖：TASK-077**

- [ ] **TASK-078** 联调多轮 messages；**依赖：TASK-077**
  - [ ] **TASK-078-S01** 使用 system + user 请求；**依赖：TASK-078**
  - [ ] **TASK-078-S02** 使用 user + assistant + user 请求；**依赖：TASK-078**
  - [ ] **TASK-078-S03** 验证 Qwen 当前请求根据完整 messages 重建上下文；**依赖：TASK-078**

- [ ] **TASK-079** 联调客户端断开连接；**依赖：TASK-077、TASK-052、TASK-037**
  - [ ] **TASK-079-S01** 发起流式请求；**依赖：TASK-079**
  - [ ] **TASK-079-S02** 在流式过程中主动断开客户端；**依赖：TASK-079**
  - [ ] **TASK-079-S03** 验证 Go Context 被取消；**依赖：TASK-079**
  - [ ] **TASK-079-S04** 验证 Playwright 操作停止；**依赖：TASK-079**
  - [ ] **TASK-079-S05** 验证账号最终恢复 `idle`；**依赖：TASK-079**

---

# 18. 并发、限流与异常场景测试

- [ ] **TASK-080** 测试 API Key 10 RPM 限流；**依赖：TASK-043、TASK-077**
  - [ ] **TASK-080-S01** 连续发送 10 个请求并确认允许；**依赖：TASK-080**
  - [ ] **TASK-080-S02** 第 11 个请求确认返回 429；**依赖：TASK-080**
  - [ ] **TASK-080-S03** 跨分钟后确认计数重置；**依赖：TASK-080**
  - [ ] **TASK-080-S04** 确认不同 API Key 之间限流独立；**依赖：TASK-080**

- [ ] **TASK-081** 测试多账号 Round Robin；**依赖：TASK-039、TASK-077**
  - [ ] **TASK-081-S01** 准备 A/B/C 三个正常账号；**依赖：TASK-081**
  - [ ] **TASK-081-S02** 发起多个可并发请求；**依赖：TASK-081**
  - [ ] **TASK-081-S03** 验证请求按可用账号 Round Robin 分配；**依赖：TASK-081**
  - [ ] **TASK-081-S04** 验证 busy 账号不会被再次分配；**依赖：TASK-081**

- [ ] **TASK-082** 测试账号全部忙碌时的 FIFO；**依赖：TASK-081、TASK-036**
  - [ ] **TASK-082-S01** 使用单账号制造一个长时间运行请求；**依赖：TASK-082**
  - [ ] **TASK-082-S02** 连续提交多个等待请求；**依赖：TASK-082**
  - [ ] **TASK-082-S03** 验证等待顺序保持 FIFO；**依赖：TASK-082**
  - [ ] **TASK-082-S04** 释放账号后验证队首请求先获得账号；**依赖：TASK-082**

- [ ] **TASK-083** 测试 30 秒账号等待超时；**依赖：TASK-082**
  - [ ] **TASK-083-S01** 保持所有账号 busy；**依赖：TASK-083**
  - [ ] **TASK-083-S02** 发起新的 API 请求；**依赖：TASK-083**
  - [ ] **TASK-083-S03** 验证最多等待 30 秒；**依赖：TASK-083**
  - [ ] **TASK-083-S04** 验证最终返回 503；**依赖：TASK-083**

- [ ] **TASK-084** 测试 Qwen 执行失败；**依赖：TASK-053、TASK-077**
  - [ ] **TASK-084-S01** 模拟或制造 Provider 执行失败；**依赖：TASK-084**
  - [ ] **TASK-084-S02** 验证 API 返回 500；**依赖：TASK-084**
  - [ ] **TASK-084-S03** 验证账号进入 `error`；**依赖：TASK-084**
  - [ ] **TASK-084-S04** 验证异常账号不会继续被调度；**依赖：TASK-084**

- [ ] **TASK-085** 测试非法请求参数；**依赖：TASK-046、TASK-077**
  - [ ] **TASK-085-S01** 测试 `stream=false` 返回 400；**依赖：TASK-085**
  - [ ] **TASK-085-S02** 测试不支持的 role 返回 400；**依赖：TASK-085**
  - [ ] **TASK-085-S03** 测试非 string content 返回 400；**依赖：TASK-085**
  - [ ] **TASK-085-S04** 测试缺少 model/messages 返回 400；**依赖：TASK-085**

---

# 19. 数据持久化与重启恢复测试

- [ ] **TASK-086** 测试 JSON 数据持久化；**依赖：TASK-016、TASK-075、TASK-076**
  - [ ] **TASK-086-S01** 创建账号并确认 `accounts.json` 更新；**依赖：TASK-086**
  - [ ] **TASK-086-S02** 创建 API Key 并确认 `api_keys.json` 更新；**依赖：TASK-086**
  - [ ] **TASK-086-S03** 确认日志写入 `api_logs.json`；**依赖：TASK-086**
  - [ ] **TASK-086-S04** 检查 JSON 文件内容可正常解析；**依赖：TASK-086**

- [ ] **TASK-087** 测试服务重启恢复；**依赖：TASK-086、TASK-023**
  - [ ] **TASK-087-S01** 停止 Go 服务；**依赖：TASK-087**
  - [ ] **TASK-087-S02** 重新启动 Go 服务；**依赖：TASK-087**
  - [ ] **TASK-087-S03** 验证管理员数据恢复；**依赖：TASK-087**
  - [ ] **TASK-087-S04** 验证 API Key 数据恢复；**依赖：TASK-087**
  - [ ] **TASK-087-S05** 验证 Qwen Account 数据恢复；**依赖：TASK-087**
  - [ ] **TASK-087-S06** 验证账号重新建立 Browser Context；**依赖：TASK-087**
  - [ ] **TASK-087-S07** 验证正常账号重新进入可调度状态；**依赖：TASK-087**

---

# 20. 前端完整业务联调

- [ ] **TASK-088** 完成 Login → Dashboard 完整流程；**依赖：TASK-074、TASK-064**
  - [ ] **TASK-088-S01** 浏览器访问 `/login`；**依赖：TASK-088**
  - [ ] **TASK-088-S02** 登录成功进入 Dashboard；**依赖：TASK-088**
  - [ ] **TASK-088-S03** 刷新页面后保持有效 JWT；**依赖：TASK-088**
  - [ ] **TASK-088-S04** JWT 失效后重新进入 Login；**依赖：TASK-088**

- [ ] **TASK-089** 完成 Account 页面完整流程；**依赖：TASK-075、TASK-067**
  - [ ] **TASK-089-S01** 前端创建账号；**依赖：TASK-089**
  - [ ] **TASK-089-S02** 查看验证结果；**依赖：TASK-089**
  - [ ] **TASK-089-S03** 编辑账号；**依赖：TASK-089**
  - [ ] **TASK-089-S04** 禁用和启用账号；**依赖：TASK-089**
  - [ ] **TASK-089-S05** 重新验证账号；**依赖：TASK-089**
  - [ ] **TASK-089-S06** 删除账号；**依赖：TASK-089**

- [ ] **TASK-090** 完成 API Key 页面完整流程；**依赖：TASK-076、TASK-070**
  - [ ] **TASK-090-S01** 创建 API Key；**依赖：TASK-090**
  - [ ] **TASK-090-S02** 验证 Key 状态；**依赖：TASK-090**
  - [ ] **TASK-090-S03** 禁用 Key；**依赖：TASK-090**
  - [ ] **TASK-090-S04** 启用 Key；**依赖：TASK-090**
  - [ ] **TASK-090-S05** 删除 Key；**依赖：TASK-090**

- [ ] **TASK-091** 完成 Chat Test 完整流程；**依赖：TASK-073、TASK-089、TASK-090**
  - [ ] **TASK-091-S01** 选择 API Key；**依赖：TASK-091**
  - [ ] **TASK-091-S02** 输入 Model；**依赖：TASK-091**
  - [ ] **TASK-091-S03** 输入问题；**依赖：TASK-091**
  - [ ] **TASK-091-S04** 发起请求；**依赖：TASK-091**
  - [ ] **TASK-091-S05** 验证流式回答逐步显示；**依赖：TASK-091**
  - [ ] **TASK-091-S06** 验证 `[DONE]` 后请求完成；**依赖：TASK-091**
  - [ ] **TASK-091-S07** 验证异常请求能显示错误；**依赖：TASK-091**

---

# 21. 最终 Go 单体运行模式

- [ ] **TASK-092** 实现 React 静态资源构建流程；**依赖：TASK-003、TASK-091**
  - [ ] **TASK-092-S01** 执行 React production build；**依赖：TASK-092**
  - [ ] **TASK-092-S02** 确认构建产物生成；**依赖：TASK-092**
  - [ ] **TASK-092-S03** 明确 Go 需要托管的静态资源目录；**依赖：TASK-092**

- [ ] **TASK-093** 实现 Go 托管 React 静态资源；**依赖：TASK-092、TASK-004**
  - [ ] **TASK-093-S01** 配置静态资源访问；**依赖：TASK-093**
  - [ ] **TASK-093-S02** 配置 SPA History Fallback，使 `/dashboard` 等页面可直接访问；**依赖：TASK-093**
  - [ ] **TASK-093-S03** 确保 `/admin/*` API 路由不被静态资源处理逻辑拦截；**依赖：TASK-093**
  - [ ] **TASK-093-S04** 确保 `/v1/*` API 路由不被静态资源处理逻辑拦截；**依赖：TASK-093**

- [ ] **TASK-094** 验证最终单进程运行模式；**依赖：TASK-093**
  - [ ] **TASK-094-S01** 仅启动 Go 服务；**依赖：TASK-094**
  - [ ] **TASK-094-S02** 浏览器访问 Go 提供的 React 页面；**依赖：TASK-094**
  - [ ] **TASK-094-S03** 完成管理员登录；**依赖：TASK-094**
  - [ ] **TASK-094-S04** 完成账号管理；**依赖：TASK-094**
  - [ ] **TASK-094-S05** 完成 API Key 管理；**依赖：TASK-094**
  - [ ] **TASK-094-S06** 完成 Chat Test；**依赖：TASK-094**
  - [ ] **TASK-094-S07** 使用 curl 调用 `/v1/chat/completions`；**依赖：TASK-094**

---

# 22. 测试、质量与稳定性

## 22.1 单元测试

- [ ] **TASK-095** 为 JSON Storage 增加单元测试；**依赖：TASK-011**
  - [ ] **TASK-095-S01** 测试读写；**依赖：TASK-095**
  - [ ] **TASK-095-S02** 测试文件不存在初始化；**依赖：TASK-095**
  - [ ] **TASK-095-S03** 测试并发写入保护；**依赖：TASK-095**
  - [ ] **TASK-095-S04** 测试临时文件 Rename 流程；**依赖：TASK-095**

- [ ] **TASK-096** 为鉴权和限流增加单元测试；**依赖：TASK-018、TASK-043**
  - [ ] **TASK-096-S01** 测试 JWT 正常校验；**依赖：TASK-096**
  - [ ] **TASK-096-S02** 测试 JWT 过期；**依赖：TASK-096**
  - [ ] **TASK-096-S03** 测试无效 API Key；**依赖：TASK-096**
  - [ ] **TASK-096-S04** 测试禁用 API Key；**依赖：TASK-096**
  - [ ] **TASK-096-S05** 测试 10 RPM；**依赖：TASK-096**

- [ ] **TASK-097** 为 Chat 参数校验和错误响应增加单元测试；**依赖：TASK-046、TASK-047**
  - [ ] **TASK-097-S01** 测试合法 messages；**依赖：TASK-097**
  - [ ] **TASK-097-S02** 测试非法 role；**依赖：TASK-097**
  - [ ] **TASK-097-S03** 测试非法 content；**依赖：TASK-097**
  - [ ] **TASK-097-S04** 测试 stream=false；**依赖：TASK-097**
  - [ ] **TASK-097-S05** 测试各 HTTP 错误码；**依赖：TASK-097**

## 22.2 静态检查与并发检查

- [ ] **TASK-098** 完成 Go 代码质量检查；**依赖：TASK-094、TASK-097**
  - [ ] **TASK-098-S01** 执行 `gofmt`；**依赖：TASK-098**
  - [ ] **TASK-098-S02** 执行 `go vet`；**依赖：TASK-098**
  - [ ] **TASK-098-S03** 执行 `go test ./...`；**依赖：TASK-098**
  - [ ] **TASK-098-S04** 执行 `go test -race ./...`；**依赖：TASK-098**
  - [ ] **TASK-098-S05** 修复发现的并发数据竞争；**依赖：TASK-098**

- [ ] **TASK-099** 完成前端构建和检查；**依赖：TASK-091**
  - [ ] **TASK-099-S01** 执行 production build；**依赖：TASK-099**
  - [ ] **TASK-099-S02** 修复 TypeScript/ESLint/构建错误；**依赖：TASK-099**
  - [ ] **TASK-099-S03** 验证 production build 可被 Go 托管；**依赖：TASK-099**

---

# 23. MVP 验收测试

- [ ] **TASK-100** 验收管理员登录；**依赖：TASK-098、TASK-099**
  - [ ] **TASK-100-S01** 启动 Go；**依赖：TASK-100**
  - [ ] **TASK-100-S02** 自动创建默认管理员；**依赖：TASK-100**
  - [ ] **TASK-100-S03** 登录成功；**依赖：TASK-100**
  - [ ] **TASK-100-S04** 登录后进入 Dashboard；**依赖：TASK-100**

- [ ] **TASK-101** 验收 Qwen 账号管理；**依赖：TASK-100**
  - [ ] **TASK-101-S01** 添加有效 Qwen 账号；**依赖：TASK-101**
  - [ ] **TASK-101-S02** 验证成功进入正常状态；**依赖：TASK-101**
  - [ ] **TASK-101-S03** 添加无效账号并确认进入异常状态；**依赖：TASK-101**
  - [ ] **TASK-101-S04** 重新验证并恢复账号；**依赖：TASK-101**
  - [ ] **TASK-101-S05** 验证禁用账号不会接受新请求；**依赖：TASK-101**
  - [ ] **TASK-101-S06** 验证删除账号生命周期；**依赖：TASK-101**

- [ ] **TASK-102** 验收 API Key；**依赖：TASK-100**
  - [ ] **TASK-102-S01** 创建 API Key；**依赖：TASK-102**
  - [ ] **TASK-102-S02** 确认格式为 `sk-` + 32 位随机字符串；**依赖：TASK-102**
  - [ ] **TASK-102-S03** 使用 Key 调用 API 成功；**依赖：TASK-102**
  - [ ] **TASK-102-S04** 禁用 Key 后调用返回 401；**依赖：TASK-102**
  - [ ] **TASK-102-S05** 删除 Key 后调用失败；**依赖：TASK-102**

- [ ] **TASK-103** 验收 OpenAI-compatible Chat API；**依赖：TASK-102、TASK-101**
  - [ ] **TASK-103-S01** 使用 curl 调用 `/v1/chat/completions`；**依赖：TASK-103**
  - [ ] **TASK-103-S02** 验证 Authorization Bearer Key；**依赖：TASK-103**
  - [ ] **TASK-103-S03** 验证 10 RPM；**依赖：TASK-103**
  - [ ] **TASK-103-S04** 验证 Qwen 账号调度；**依赖：TASK-103**
  - [ ] **TASK-103-S05** 验证流式响应；**依赖：TASK-103**
  - [ ] **TASK-103-S06** 验证最终 `data: [DONE]`；**依赖：TASK-103**
  - [ ] **TASK-103-S07** 验证 API Log；**依赖：TASK-103**

- [ ] **TASK-104** 验收多账号调度；**依赖：TASK-103**
  - [ ] **TASK-104-S01** 准备多个正常账号；**依赖：TASK-104**
  - [ ] **TASK-104-S02** 连续执行多个请求；**依赖：TASK-104**
  - [ ] **TASK-104-S03** 验证 Round Robin；**依赖：TASK-104**
  - [ ] **TASK-104-S04** 验证 busy 账号跳过；**依赖：TASK-104**
  - [ ] **TASK-104-S05** 验证全部 busy 时 FIFO；**依赖：TASK-104**
  - [ ] **TASK-104-S06** 验证 30 秒超时返回 503；**依赖：TASK-104**

- [ ] **TASK-105** 验收异常处理；**依赖：TASK-103、TASK-104**
  - [ ] **TASK-105-S01** 验证无效 Key 返回 401；**依赖：TASK-105**
  - [ ] **TASK-105-S02** 验证非法请求返回 400；**依赖：TASK-105**
  - [ ] **TASK-105-S03** 验证限流返回 429；**依赖：TASK-105**
  - [ ] **TASK-105-S04** 验证无可用账号返回 503；**依赖：TASK-105**
  - [ ] **TASK-105-S05** 验证 Qwen 执行失败返回 500；**依赖：TASK-105**
  - [ ] **TASK-105-S06** 验证 `/v1/messages` 返回 501；**依赖：TASK-105**

- [ ] **TASK-106** 验收 Chat Test 页面；**依赖：TASK-103、TASK-091**
  - [ ] **TASK-106-S01** 选择 API Key；**依赖：TASK-106**
  - [ ] **TASK-106-S02** 发起聊天请求；**依赖：TASK-106**
  - [ ] **TASK-106-S03** 验证回答流式显示；**依赖：TASK-106**
  - [ ] **TASK-106-S04** 验证错误提示；**依赖：TASK-106**

- [ ] **TASK-107** 验收重启恢复；**依赖：TASK-106、TASK-087**
  - [ ] **TASK-107-S01** 重启 Go；**依赖：TASK-107**
  - [ ] **TASK-107-S02** 验证管理员数据恢复；**依赖：TASK-107**
  - [ ] **TASK-107-S03** 验证 API Key 恢复；**依赖：TASK-107**
  - [ ] **TASK-107-S04** 验证账号恢复；**依赖：TASK-107**
  - [ ] **TASK-107-S05** 再次调用 Chat API；**依赖：TASK-107**

---

# 24. 发布前整理

- [ ] **TASK-108** 完善 README；**依赖：TASK-107**
  - [ ] **TASK-108-S01** 写明项目定位；**依赖：TASK-108**
  - [ ] **TASK-108-S02** 写明开发环境安装步骤；**依赖：TASK-108**
  - [ ] **TASK-108-S03** 写明 Playwright 浏览器安装步骤；**依赖：TASK-108**
  - [ ] **TASK-108-S04** 写明开发模式启动步骤；**依赖：TASK-108**
  - [ ] **TASK-108-S05** 写明最终单体运行步骤；**依赖：TASK-108**
  - [ ] **TASK-108-S06** 写明默认管理员账号和密码；**依赖：TASK-108**
  - [ ] **TASK-108-S07** 写明 `/v1/chat/completions` 调用示例；**依赖：TASK-108**
  - [ ] **TASK-108-S08** 写明 MVP 已知限制；**依赖：TASK-108**

- [ ] **TASK-109** 检查敏感数据和运行数据处理；**依赖：TASK-108**
  - [ ] **TASK-109-S01** 确认真实 Qwen 凭证不会进入 Git；**依赖：TASK-109**
  - [ ] **TASK-109-S02** 确认真实 API Key 不会进入示例代码；**依赖：TASK-109**
  - [ ] **TASK-109-S03** 确认 `data/*.json` 不提交真实运行数据；**依赖：TASK-109**
  - [ ] **TASK-109-S04** 确认 README 明确 MVP 明文凭证仅适合本地使用；**依赖：TASK-109**

- [ ] **TASK-110** 完成最终 MVP 回归测试；**依赖：TASK-098、TASK-099、TASK-109**
  - [ ] **TASK-110-S01** 从空数据目录启动服务；**依赖：TASK-110**
  - [ ] **TASK-110-S02** 登录管理后台；**依赖：TASK-110**
  - [ ] **TASK-110-S03** 添加 Qwen Account；**依赖：TASK-110**
  - [ ] **TASK-110-S04** 创建 API Key；**依赖：TASK-110**
  - [ ] **TASK-110-S05** 使用 Chat Test 发起请求；**依赖：TASK-110**
  - [ ] **TASK-110-S06** 使用 curl 发起请求；**依赖：TASK-110**
  - [ ] **TASK-110-S07** 验证流式响应；**依赖：TASK-110**
  - [ ] **TASK-110-S08** 验证日志；**依赖：TASK-110**
  - [ ] **TASK-110-S09** 重启服务并再次验证；**依赖：TASK-110**

- [ ] **TASK-111** 完成 MVP 交付确认；**依赖：TASK-110**
  - [ ] **TASK-111-S01** 确认管理后台可用；**依赖：TASK-111**
  - [ ] **TASK-111-S02** 确认 Qwen 账号池可用；**依赖：TASK-111**
  - [ ] **TASK-111-S03** 确认 API Key 可用；**依赖：TASK-111**
  - [ ] **TASK-111-S04** 确认 `/v1/chat/completions` 可用；**依赖：TASK-111**
  - [ ] **TASK-111-S05** 确认 SSE 流式返回可用；**依赖：TASK-111**
  - [ ] **TASK-111-S06** 确认核心异常场景已验证；**依赖：TASK-111**
  - [ ] **TASK-111-S07** 确认 README 足以让新环境重新运行项目；**依赖：TASK-111**

---

# 25. 依赖关系总览

```text
TASK-000
  ↓
TASK-001
  ├── TASK-002 → TASK-004 → TASK-006
  └── TASK-003

TASK-006
  ↓
TASK-007 / TASK-008 / TASK-009 / TASK-010
  ↓
TASK-011
  ↓
TASK-012 / TASK-013 / TASK-014 / TASK-015
  ↓
TASK-016

TASK-016
  ↓
TASK-017 → TASK-018 → TASK-019 → TASK-021
                    └→ TASK-020

TASK-022 → TASK-023 → TASK-024
                    ├→ TASK-025
                    └→ TASK-026
                         ↓
                       TASK-027

TASK-013 + TASK-023
  ↓
TASK-028
  ├→ TASK-029
  ├→ TASK-030
  ├→ TASK-031
  └→ TASK-032
       ↓
TASK-033

TASK-028 + TASK-031 + TASK-032
  ↓
TASK-034 → TASK-035 → TASK-036 → TASK-037 → TASK-038 → TASK-039

TASK-014 → TASK-040 → TASK-041 → TASK-042 → TASK-043

TASK-044 → TASK-045 → TASK-048
TASK-044 → TASK-046 → TASK-047

TASK-048 + TASK-025
  ↓
TASK-049 → TASK-050 → TASK-051 → TASK-052
                              └→ TASK-053

TASK-043 + TASK-046 + TASK-047 + TASK-037 + TASK-048
  ↓
TASK-054 → TASK-055 → TASK-056
                      └→ TASK-057

TASK-015 + TASK-054
  ↓
TASK-058 → TASK-059

TASK-003 + TASK-021
  ↓
TASK-060 → TASK-061 → TASK-062

TASK-033 + TASK-041 + TASK-058
  ↓
TASK-063 → TASK-064

TASK-033 → TASK-065 → TASK-066 / TASK-067
TASK-041 → TASK-068 → TASK-069 / TASK-070

TASK-062 + TASK-068
  ↓
TASK-071 → TASK-072 → TASK-073

TASK-021 + TASK-062
  ↓
TASK-074

TASK-033 + TASK-066
  ↓
TASK-075

TASK-041 + TASK-069
  ↓
TASK-076

TASK-056 + TASK-075 + TASK-076
  ↓
TASK-077 → TASK-078
         └→ TASK-079

TASK-043 + TASK-077 → TASK-080
TASK-039 + TASK-077 → TASK-081 → TASK-082 → TASK-083
TASK-053 + TASK-077 → TASK-084
TASK-046 + TASK-077 → TASK-085

TASK-016 + TASK-075 + TASK-076 → TASK-086 → TASK-087

TASK-074 + TASK-064 → TASK-088
TASK-075 + TASK-067 → TASK-089
TASK-076 + TASK-070 → TASK-090
TASK-073 + TASK-089 + TASK-090 → TASK-091

TASK-003 + TASK-091 → TASK-092 → TASK-093 → TASK-094

TASK-011 → TASK-095
TASK-018 + TASK-043 → TASK-096
TASK-046 + TASK-047 → TASK-097
TASK-094 + TASK-097 → TASK-098
TASK-091 → TASK-099

TASK-098 + TASK-099 → TASK-100
TASK-100 → TASK-101 / TASK-102
TASK-101 + TASK-102 → TASK-103
TASK-103 → TASK-104
TASK-103 + TASK-104 → TASK-105
TASK-103 + TASK-091 → TASK-106
TASK-106 + TASK-087 → TASK-107

TASK-107 → TASK-108 → TASK-109
TASK-098 + TASK-099 + TASK-109 → TASK-110 → TASK-111
```

---

# 26. MVP 完成定义

- [ ] **TASK-112** 确认 MVP Definition of Done；**依赖：TASK-111**
  - [ ] **TASK-112-S01** Go 服务能够独立启动；**依赖：TASK-112**
  - [ ] **TASK-112-S02** React 管理后台能够通过 Go 单体模式访问；**依赖：TASK-112**
  - [ ] **TASK-112-S03** 管理员可以登录；**依赖：TASK-112**
  - [ ] **TASK-112-S04** 管理员可以添加、编辑、启用、禁用、验证、删除 Qwen 账号；**依赖：TASK-112**
  - [ ] **TASK-112-S05** 系统可以维护多个 Qwen Account Context；**依赖：TASK-112**
  - [ ] **TASK-112-S06** 系统可以 Round Robin 调度空闲账号；**依赖：TASK-112**
  - [ ] **TASK-112-S07** 系统可以在账号全部忙碌时 FIFO 等待最多 30 秒；**依赖：TASK-112**
  - [ ] **TASK-112-S08** 管理员可以创建、启用、禁用、删除 API Key；**依赖：TASK-112**
  - [ ] **TASK-112-S09** API Key 按固定窗口执行 10 RPM；**依赖：TASK-112**
  - [ ] **TASK-112-S10** `/v1/chat/completions` 支持 OpenAI-compatible 请求；**依赖：TASK-112**
  - [ ] **TASK-112-S11** 只支持 `system/user/assistant` 和 string content；**依赖：TASK-112**
  - [ ] **TASK-112-S12** 只接受 `stream=true`；**依赖：TASK-112**
  - [ ] **TASK-112-S13** Qwen 响应能够转换为 OpenAI-compatible SSE；**依赖：TASK-112**
  - [ ] **TASK-112-S14** 流结束时输出 `data: [DONE]`；**依赖：TASK-112**
  - [ ] **TASK-112-S15** 客户端断开后能够取消 Qwen 操作并释放账号；**依赖：TASK-112**
  - [ ] **TASK-112-S16** API 调用记录基础日志且不记录 messages；**依赖：TASK-112**
  - [ ] **TASK-112-S17** `/v1/messages` 返回 501；**依赖：TASK-112**
  - [ ] **TASK-112-S18** 服务重启后 JSON 数据可以恢复；**依赖：TASK-112**
  - [ ] **TASK-112-S19** 单元测试、race detector、前端 production build 均通过；**依赖：TASK-112**
  - [ ] **TASK-112-S20** README 能够指导开发环境和最终单体模式运行；**依赖：TASK-112**

---

## 27. 明确不进入第一版的任务

以下内容**不要因为开发过程中出现需求扩张而提前加入**：

- [ ] **TASK-113** 保持 MVP 范围边界；**依赖：TASK-112**
  - [ ] **TASK-113-S01** 不加入 MySQL/PostgreSQL；**依赖：TASK-113**
  - [ ] **TASK-113-S02** 不加入 Redis；**依赖：TASK-113**
  - [ ] **TASK-113-S03** 不加入 Docker/Kubernetes；**依赖：TASK-113**
  - [ ] **TASK-113-S04** 不加入多租户；**依赖：TASK-113**
  - [ ] **TASK-113-S05** 不加入 RBAC；**依赖：TASK-113**
  - [ ] **TASK-113-S06** 不加入多管理员；**依赖：TASK-113**
  - [ ] **TASK-113-S07** 不加入复杂 Token Bucket/Sliding Window 限流；**依赖：TASK-113**
  - [ ] **TASK-113-S08** 不加入全局限流；**依赖：TASK-113**
  - [ ] **TASK-113-S09** 不加入日志管理页面；**依赖：TASK-113**
  - [ ] **TASK-113-S10** 不加入高级监控和实时大盘；**依赖：TASK-113**
  - [ ] **TASK-113-S11** 不实现非 Qwen Provider；**依赖：TASK-113**
  - [ ] **TASK-113-S12** 不实现非流式 Chat Completion；**依赖：TASK-113**
  - [ ] **TASK-113-S13** 不实现 Anthropic `/v1/messages` 业务逻辑；**依赖：TASK-113**

---

## 28. 推荐实际执行顺序

- [ ] **TASK-114** 按以下顺序推进开发；**依赖：TASK-113**
  - [ ] **TASK-114-S01** 第一轮：完成 TASK-001 ～ TASK-021，先让项目、JSON 数据层、管理员 JWT 跑起来；**依赖：TASK-114**
  - [ ] **TASK-114-S02** 第二轮：完成 TASK-022 ～ TASK-039，解决 Playwright、Qwen 账号和账号池；**依赖：TASK-114**
  - [ ] **TASK-114-S03** 第三轮：完成 TASK-040 ～ TASK-047，解决 API Key、限流、Chat 数据结构和错误；**依赖：TASK-114**
  - [ ] **TASK-114-S04** 第四轮：完成 TASK-048 ～ TASK-059，打通 Qwen Provider、SSE 和日志；**依赖：TASK-114**
  - [ ] **TASK-114-S05** 第五轮：完成 TASK-060 ～ TASK-073，完成 React 管理后台；**依赖：TASK-114**
  - [ ] **TASK-114-S06** 第六轮：完成 TASK-074 ～ TASK-091，进行端到端联调和异常测试；**依赖：TASK-114**
  - [ ] **TASK-114-S07** 第七轮：完成 TASK-092 ～ TASK-099，完成最终单体运行和代码质量检查；**依赖：TASK-114**
  - [ ] **TASK-114-S08** 第八轮：完成 TASK-100 ～ TASK-112，进行最终 MVP 验收；**依赖：TASK-114**

---

## 29. 最终核心链路检查

- [ ] **TASK-115** 验证最终核心链路完整闭环；**依赖：TASK-112**
  - [ ] **TASK-115-S01** 浏览器进入管理后台；**依赖：TASK-115**
  - [ ] **TASK-115-S02** 管理员登录并获得 JWT；**依赖：TASK-115**
  - [ ] **TASK-115-S03** 添加并验证 Qwen Account；**依赖：TASK-115**
  - [ ] **TASK-115-S04** 创建 API Key；**依赖：TASK-115**
  - [ ] **TASK-115-S05** 第三方客户端调用 `/v1/chat/completions`；**依赖：TASK-115**
  - [ ] **TASK-115-S06** API Key 鉴权通过；**依赖：TASK-115**
  - [ ] **TASK-115-S07** 10 RPM 限流通过；**依赖：TASK-115**
  - [ ] **TASK-115-S08** Account Scheduler 分配 Qwen Account；**依赖：TASK-115**
  - [ ] **TASK-115-S09** Qwen Provider 执行 Web 模型请求；**依赖：TASK-115**
  - [ ] **TASK-115-S10** Provider 持续产生文本 chunk；**依赖：TASK-115**
  - [ ] **TASK-115-S11** Web2API 转换为 OpenAI-compatible SSE；**依赖：TASK-115**
  - [ ] **TASK-115-S12** 客户端持续接收流式内容；**依赖：TASK-115**
  - [ ] **TASK-115-S13** 输出 `data: [DONE]`；**依赖：TASK-115**
  - [ ] **TASK-115-S14** 请求结束后释放 Qwen Account；**依赖：TASK-115**
  - [ ] **TASK-115-S15** 写入调用日志；**依赖：TASK-115**
  - [ ] **TASK-115-S16** 整个链路完成后，项目达到 MVP 验收条件；**依赖：TASK-115**
