下面我把我们前面已经确认的内容整理成一版 **Web2API MVP 产品需求文档（PRD v1.0）**。未明确的地方，我按“**第一版怎么简单怎么来**”处理，避免继续把范围做大。

# Web2API 项目 PRD v1.0

## 1. 产品概述

### 1.1 产品名称

Web2API

### 1.2 产品定位

Web2API 是一个面向个人开发者的 Web 大模型 API 网关。

用户可以在 Web 管理后台中维护第三方大模型网页账号，创建 API Key。外部应用通过标准的 OpenAI-compatible API 调用 Web2API，系统内部负责：

* API Key 鉴权
* API Key 限流
* Qwen 账号池调度
* Web 模型请求执行
* Web 模型响应转换
* 流式响应输出
* 调用日志记录

第一阶段以 `chat.qwen.ai` 作为第一个**模型接入实例**。

### 1.3 第一阶段核心目标

完成下面这条完整链路：

```text
第三方应用
    ↓
OpenAI-compatible API
    ↓
API Key 鉴权
    ↓
10 RPM 限流
    ↓
Qwen 账号池调度
    ↓
Playwright / Qwen Provider
    ↓
chat.qwen.ai
    ↓
流式响应转换
    ↓
OpenAI SSE
    ↓
第三方应用
```

第一阶段以“**先把完整框架跑通**”为目标，不追求复杂的账号调度、数据库和部署能力。

---

# 2. 技术栈

## 2.1 后端

* Go
* Gin
* Playwright-Go

## 2.2 前端

* React
* Vite

## 2.3 数据存储

第一阶段不使用 MySQL / PostgreSQL 等数据库。

使用 JSON 文件模拟数据库表：

```text
data/
├── admin.json
├── accounts.json
├── api_keys.json
└── api_logs.json
```

## 2.4 运行方式

开发阶段：

```text
React/Vite
localhost:5173

Go/Gin
localhost:8080
```

最终运行：

```text
Go
localhost:8080

├── React 静态资源
├── /admin/*
└── /v1/*
```

第一阶段暂不考虑 Docker、Kubernetes、云部署。

---

# 3. 产品角色

第一阶段只有一个角色：

## 管理员

管理员可以：

* 登录后台
* 管理 Qwen 账号
* 管理 API Key
* 查看 Dashboard
* 使用 API Key 进行测试对话

暂不支持：

* 多管理员
* 用户注册
* RBAC
* 多角色权限
* 多租户

---

# 4. 管理员登录

## 4.1 登录方式

采用：

```text
用户名 + 密码
```

登录成功后使用 JWT 作为管理后台登录凭证。

请求：

```http
POST /admin/login
```

成功后返回 JWT。

后续管理接口：

```http
Authorization: Bearer <JWT>
```

## 4.2 JWT

第一阶段：

* 有效期：24 小时
* 不使用 Refresh Token
* 不实现 JWT 黑名单
* Token 过期后重新登录

## 4.3 默认管理员

系统首次启动时：

```text
admin.json 不存在
        ↓
自动创建默认管理员
```

默认用户名：

```text
admin
```

默认密码由第一版实现统一约定，并写入 README。

不强制首次修改密码。

> 注意：第一版主要用于本地自用，因此管理员密码明文存储，仅适合作为 MVP 方案，不建议直接用于公网生产环境。

---

# 5. Qwen 账号管理

## 5.1 账号概念

一个 Qwen 账号对应一个独立的 Playwright Browser Context。

整体关系：

```text
Browser
├── Context A → Qwen Account A
├── Context B → Qwen Account B
└── Context C → Qwen Account C
```

整个 Go 服务只启动一个 Browser。

每个账号使用独立 Context，确保登录态相互隔离。

## 5.2 添加账号

管理后台提供：

```text
添加 Qwen 账号
```

字段：

```text
Username     必填
Password     可选
Access Token 可选
```

要求：

```text
Password / Access Token
至少填写一个
```

合法情况：

```text
Username + Password
Username + Token
Username + Password + Token
```

非法情况：

```text
只有 Username
```

## 5.3 凭证使用优先级

优先使用：

```text
Access Token
```

如果 Token 失效：

```text
Token 失效
    ↓
使用 Username + Password
    ↓
尝试重新建立登录态
```

密码恢复最多自动尝试 1 次。

如果仍然失败：

```text
账号 → 异常
```

## 5.4 添加账号时验证

账号创建后必须立即验证。

流程：

```text
添加账号
    ↓
创建 Browser Context
    ↓
建立登录态
    ↓
访问 / 验证 chat.qwen.ai
    ↓
验证成功 → 正常
验证失败 → 异常
```

只有正常账号才可以进入账号池。

## 5.5 账号状态

持久化状态：

```text
normal
error
disabled
```

运行时状态：

```text
idle
busy
```

运行时状态不写入 JSON。

例如：

```text
Account A

persistent:
status = normal

runtime:
busy = true
currentRequestID = xxx
```

## 5.6 账号编辑

支持修改：

```text
Username
Password
Access Token
```

保存后重新验证：

```text
修改凭证
    ↓
重新建立 Playwright Context
    ↓
验证
    ├── 成功 → normal
    └── 失败 → error
```

## 5.7 账号禁用

管理员禁用账号后：

```text
disabled
    ↓
不再接受新请求
```

如果账号空闲：

```text
立即关闭 Browser Context
```

如果账号正在执行请求：

```text
等待当前请求结束
    ↓
关闭 Context
```

重新启用：

```text
创建新的 Context
    ↓
重新验证
    ↓
验证成功 → normal
```

## 5.8 账号删除

如果账号空闲：

```text
关闭 Context
    ↓
从 accounts.json 删除
```

如果账号正在执行：

```text
标记待删除
    ↓
不再接受新请求
    ↓
当前请求继续执行
    ↓
请求结束
    ↓
关闭 Context
    ↓
从 accounts.json 删除
```

---

# 6. Qwen 账号池

所有 API Key 共用一个 Qwen 账号池：

```text
API Key A ──┐
API Key B ──┼──→ Qwen Account Pool
API Key C ──┘
```

## 6.1 单账号并发

第一阶段：

> 一个 Qwen 账号同时只允许处理 1 个请求。

因此：

```text
Account A
├── idle → 可分配
└── busy → 不可分配
```

以后可以扩展为单账号多并发，第一阶段不实现。

## 6.2 调度策略

采用 Round Robin：

```text
Request 1 → A
Request 2 → B
Request 3 → C
Request 4 → A
```

但实际调度时只考虑：

```text
normal
+
idle
```

自动跳过：

```text
error
disabled
busy
```

## 6.3 账号不可用时

如果当前没有空闲账号：

```text
请求
 ↓
FIFO 等待队列
```

最大等待时间：

```text
30 秒
```

30 秒内：

```text
有账号释放 → 获取账号并执行
```

超过 30 秒：

```text
返回错误
```

第一阶段不设置队列最大长度。

---

# 7. API Key 管理

## 7.1 API Key 用途

API Key 用于第三方应用访问 Web2API。

请求方式：

```http
Authorization: Bearer <API-Key>
```

暂不支持：

```text
x-api-key
api-key
query 参数
```

## 7.2 API Key 功能

第一阶段只提供：

```text
创建
查看
启用 / 禁用
删除
```

## 7.3 Key 生成

系统自动生成 Key。

格式采用：

```text
sk- + 32 位随机字符串
```

例如：

```text
sk-a83f91c2d7e14b6a9f0c2e81d35a7b42
```

## 7.4 存储方式

第一阶段：

> API Key 明文存储在 `api_keys.json`。

主要目的是快速完成 MVP。

---

# 8. API Key 限流

## 8.1 限流维度

以 API Key 为主要限流维度。

## 8.2 限流规则

固定：

```text
10 RPM
```

即：

```text
每一个 API Key
每分钟最多 10 个请求
```

第一阶段不支持单独配置 RPM。

## 8.3 限流算法

使用：

```text
Fixed Window
```

例如：

```text
10:00:00 ~ 10:00:59
最多 10 次
```

下一分钟重新计算。

暂不实现：

* Token Bucket
* Leaky Bucket
* Sliding Window
* 全局限流

---

# 9. 对外 API

## 9.1 OpenAI-compatible API

第一阶段正式支持：

```http
POST /v1/chat/completions
```

请求：

```http
Authorization: Bearer sk-xxx
Content-Type: application/json
```

请求示例：

```json
{
  "model": "any-model",
  "messages": [
    {
      "role": "system",
      "content": "你是一个 Go 语言专家"
    },
    {
      "role": "user",
      "content": "什么是 goroutine？"
    }
  ],
  "stream": true
}
```

## 9.2 model

第一阶段：

> 接收任意 `model` 字段，不校验模型名称。

例如：

```text
qwen
qwen-plus
gpt-xxx
claude-xxx
abc
```

均可以进入系统。

内部统一路由到当前 Qwen Web 模型实例。

为了兼容客户端，响应中的 `model` 第一版可直接回显请求中的 model。

---

# 10. messages 规则

第一阶段仅支持：

```text
system
user
assistant
```

不支持：

```text
tool
function
developer
其他 role
```

不支持时：

```text
400 Bad Request
```

## 10.1 content

只支持：

```text
string
```

例如：

```json
{
  "role": "user",
  "content": "你好"
}
```

暂不支持：

```json
{
  "content": [
    {
      "type": "text",
      "text": "你好"
    }
  ]
}
```

---

# 11. Stream

第一阶段：

> 只支持流式请求。

因此：

```json
{
  "stream": true
}
```

是必须条件。

如果：

```json
{
  "stream": false
}
```

直接：

```text
400 Bad Request
```

---

# 12. OpenAI SSE 响应

Qwen 内部返回的数据格式对 API 客户端透明。

系统内部统一转换后，对外返回 OpenAI-compatible SSE。

示例：

```text
data: {"id":"chatcmpl-xxx","object":"chat.completion.chunk","choices":[{"index":0,"delta":{"role":"assistant","content":"你"}}]}

data: {"id":"chatcmpl-xxx","object":"chat.completion.chunk","choices":[{"index":0,"delta":{"content":"好"}}]}

data: {"id":"chatcmpl-xxx","object":"chat.completion.chunk","choices":[{"index":0,"delta":{"content":"！"}}]}

data: [DONE]
```

最终必须：

```text
data: [DONE]
```

结束流。

---

# 13. API 请求整体流程

```text
Client
  ↓
POST /v1/chat/completions
  ↓
Authorization 校验
  ↓
API Key 是否有效？
  ├── 否 → 401
  └── 是
       ↓
10 RPM 限流
       ├── 超限 → 429
       └── 正常
            ↓
        参数校验
            ↓
        stream=true？
            ├── 否 → 400
            └── 是
                 ↓
        获取可用 Qwen 账号
                 ├── 当前有空闲账号
                 │       ↓
                 │   Round Robin
                 │
                 └── 没有
                         ↓
                    FIFO 等待
                         ↓
                       最多30s
                         ├── 成功 → 获取账号
                         └── 超时 → 503
                                 ↓
                          Qwen Provider
                                 ↓
                         流式获取结果
                                 ↓
                         转换 OpenAI SSE
                                 ↓
                              Client
```

---

# 14. 多轮对话策略

API 层采用：

> **无状态模式**

后端不维护 API Session / Conversation。

客户端每次请求自己携带完整：

```text
messages
```

例如：

```text
Request 1
messages = [user]

Request 2
messages = [user, assistant, user]

Request 3
messages = [user, assistant, user, assistant, user]
```

## 14.1 Qwen Web 层

Qwen Web 层复用账号当前的 Browser Context。

但是：

> 每一次 API 请求，都根据当前请求中的 `messages` 重新构建 Qwen 当前对话。

因此：

```text
API 层：无状态

Qwen Browser Context：有状态

每次 API 请求：
messages → 重新构建 Qwen 对话
```

目标是让一次 API 请求尽量独立，不受该账号之前 API 请求内容影响。

---

# 15. 客户端断开连接

客户端主动断开连接时：

```text
Client Disconnect
       ↓
Go Request Context Cancel
       ↓
取消当前 Qwen 任务
       ↓
停止对应 Playwright 操作
       ↓
释放 Qwen Account
```

第一阶段不做后台继续执行机制。

---

# 16. Qwen Provider 抽象

Qwen 的具体 Playwright 实现不进入本次 PRD 的核心业务设计。

第一阶段只定义 Provider 抽象。

例如：

```go
type ModelProvider interface {
    Chat(
        ctx context.Context,
        req *ChatRequest,
        onChunk func(string) error,
    ) error
}
```

业务层只依赖：

```text
ChatRequest
        ↓
ModelProvider
        ↓
Stream Chunk
```

目录建议：

```text
provider/
├── provider.go
└── qwen/
    └── qwen.go
```

后续可以扩展：

```text
provider/
├── qwen/
├── doubao/
├── kimi/
└── ...
```

API 层不需要知道具体 Web 模型的实现细节。

---

# 17. Anthropic API 预留

第一阶段预留：

```http
POST /v1/messages
```

但不实现业务逻辑。

统一返回：

```text
HTTP 501 Not Implemented
```

目的：

> 在项目第一版中提前确定 Anthropic-compatible API 的接口位置，后续可以在不破坏现有架构的情况下实现协议转换。

---

# 18. API 错误格式

第一阶段统一使用 OpenAI-compatible Error 格式。

示例：

```json
{
  "error": {
    "message": "Invalid API key",
    "type": "authentication_error",
    "code": "invalid_api_key"
  }
}
```

主要错误：

| 场景                | HTTP |
| ----------------- | ---: |
| API Key 无效        |  401 |
| 请求参数错误            |  400 |
| stream != true    |  400 |
| role 不支持          |  400 |
| content 类型错误      |  400 |
| API Key 超过 10 RPM |  429 |
| 30 秒没有可用账号        |  503 |
| Qwen 执行失败         |  500 |
| Anthropic 接口未实现   |  501 |

---

# 19. API 调用日志

每次 API 调用记录基础信息：

```text
时间
API Key
Model
Qwen Account
请求耗时
请求结果
错误信息
```

不记录：

```text
messages
用户实际对话内容
```

日志文件：

```text
data/api_logs.json
```

最多保留：

```text
最近 1000 条
```

超过 1000 条后删除最旧记录。

第一阶段不提供日志管理页面。

---

# 20. JSON 数据层

JSON 文件对应“数据库表”：

```text
data/
├── admin.json
├── accounts.json
├── api_keys.json
└── api_logs.json
```

业务层不能直接到处操作 JSON。

采用：

```text
Handler
   ↓
Service
   ↓
Repository
   ↓
JSON
```

例如：

```text
AccountHandler
    ↓
AccountService
    ↓
AccountRepository
    ↓
accounts.json
```

## 20.1 内存缓存

程序启动时：

```text
JSON
 ↓
加载到内存
```

运行期间：

```text
业务代码 → 内存
```

不需要每次查询都读取 JSON。

## 20.2 写入机制

写数据时：

```text
Mutex
 ↓
修改内存
 ↓
写临时文件
 ↓
Rename
 ↓
完成
```

第一阶段不做：

* 数据库事务
* WAL
* 分布式锁
* 文件锁
* 数据库连接池

---

# 21. React 管理后台

后台页面：

```text
/login

/dashboard

/accounts
/api-keys
/chat-test
```

## 21.1 Login

功能：

* 用户名
* 密码
* 登录
* JWT 获取

## 21.2 Dashboard

只做基础状态概览：

### Qwen 账号

```text
总数
正常
忙碌
异常
禁用
```

### 当前请求

```text
执行中
等待中
```

### API 调用

```text
总请求
成功
失败
```

不做：

* 图表
* 数据分析
* 高级监控
* 实时大盘

---

# 22. Qwen 账号管理页面

列表展示：

```text
账号名称 / Username
状态
创建时间
```

支持：

```text
添加
编辑
禁用
启用
重新验证
删除
```

状态例如：

```text
正常
异常
禁用
```

添加/编辑页面：

```text
Username
Password
Access Token
```

---

# 23. API Key 管理页面

列表：

```text
名称
Key
状态
创建时间
```

功能：

```text
创建
启用
禁用
删除
```

创建时：

```text
填写名称
 ↓
系统生成 Key
```

---

# 24. 对话测试页面

这是后台中的一个非常重要的辅助功能。

它不是独立的“测试链路”。

而是：

> **直接使用已经创建的 API Key，调用真正的 `/v1/chat/completions` 接口。**

页面：

```text
API Key
[ My Test Key ▼ ]

模型
[ qwen-test ]

┌────────────────────────────┐
│                            │
│      Chat Messages         │
│                            │
└────────────────────────────┘

[ 输入问题........................ ]

[ Send ]
```

发送之后：

```text
React
 ↓
选择 API Key
 ↓
POST /v1/chat/completions
Authorization: Bearer <API-Key>
 ↓
正常经过 API 鉴权
 ↓
10 RPM
 ↓
账号池
 ↓
Qwen
 ↓
OpenAI SSE
 ↓
React 流式显示
```

因此：

> **后台对话测试和真实第三方客户端调用使用完全相同的 API 业务链路。**

这样测试页面本身也可以用于验证 Web2API 是否正常工作。

---

# 25. 账号调度器设计

第一版建议抽象：

```go
type AccountScheduler interface {
    Acquire(ctx context.Context) (*AccountRuntime, error)
    Release(accountID string)
}
```

核心逻辑：

```text
Acquire
  ↓
寻找 normal + idle
  ↓
Round Robin
  ↓
找到
  └── 标记 busy

找不到
  ↓
进入 FIFO Queue
  ↓
等待最多 30 秒
  ↓
Context 取消 / 超时
```

执行结束：

```text
Release
  ↓
busy = false
  ↓
唤醒等待队列
```

---

# 26. 账号生命周期

```text
                 添加
                  ↓
                验证
            ┌─────┴─────┐
            ↓           ↓
          正常          异常
            │            │
            │       管理员重新验证
            │            │
            │       ┌────┴────┐
            │       ↓         ↓
            │     成功       失败
            │       ↓         ↓
            └────→ 正常      异常
            │
         管理员禁用
            ↓
           禁用
            │
         重新启用
            ↓
          验证
```

请求执行失败：

```text
正常
 ↓
Qwen 请求失败
 ↓
异常
```

第一阶段不自动恢复。

管理员点击：

```text
重新验证
```

才尝试恢复。

---

# 27. 第一阶段明确不做

为了控制 MVP 范围，以下功能全部暂不实现：

```text
Claude 真正调用
多模型路由
多个模型 Provider 同时运行
图片输入
文件上传
多模态
非流式接口
Tools / Function Calling
Conversation Session 持久化
数据库
Docker
Kubernetes
云部署
多用户
RBAC
全局限流
复杂限流策略
自动恢复异常账号
请求失败自动重试
等待队列最大长度
请求执行超时
日志管理页面
数据统计图表
配额管理
计费
```

---

# 28. 第一版项目目录建议

```text
web2api/
├── cmd/
│   └── server/
│       └── main.go
│
├── internal/
│   ├── handler/
│   │   ├── admin.go
│   │   ├── account.go
│   │   ├── api_key.go
│   │   └── chat.go
│   │
│   ├── service/
│   │   ├── auth.go
│   │   ├── account.go
│   │   ├── api_key.go
│   │   └── chat.go
│   │
│   ├── repository/
│   │   ├── admin.go
│   │   ├── account.go
│   │   ├── api_key.go
│   │   └── api_log.go
│   │
│   ├── scheduler/
│   │   └── account_scheduler.go
│   │
│   ├── middleware/
│   │   ├── jwt.go
│   │   ├── api_key.go
│   │   └── rate_limit.go
│   │
│   ├── provider/
│   │   ├── provider.go
│   │   └── qwen/
│   │       └── qwen.go
│   │
│   ├── model/
│   │   ├── account.go
│   │   ├── api_key.go
│   │   ├── chat.go
│   │   └── api_log.go
│   │
│   └── storage/
│       └── json.go
│
├── web/
│   └── React 项目
│
├── data/
│   ├── admin.json
│   ├── accounts.json
│   ├── api_keys.json
│   └── api_logs.json
│
├── go.mod
├── README.md
└── ...
```

---

# 29. MVP 验收标准

第一版完成后，至少能够完成以下测试。

## 管理后台

```text
启动 Go
 ↓
自动创建管理员
 ↓
React 登录成功
 ↓
进入 Dashboard
```

## 账号

```text
添加 Qwen 账号
 ↓
验证
 ↓
成功进入正常状态
```

## API Key

```text
创建 API Key
 ↓
Key 可用
```

## API

使用：

```bash
curl http://localhost:8080/v1/chat/completions \
  -H "Authorization: Bearer sk-xxx" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "test-model",
    "messages": [
      {
        "role": "user",
        "content": "你好"
      }
    ],
    "stream": true
  }'
```

能够：

```text
请求进入
 ↓
Key 鉴权
 ↓
10 RPM
 ↓
账号调度
 ↓
Qwen
 ↓
流式返回
 ↓
data: [DONE]
```

## 多账号

拥有多个 Qwen 账号时：

```text
Request 1 → Account A
Request 2 → Account B
Request 3 → Account C
```

正常执行轮询。

## 账号占用

同一个账号：

```text
Request 1 → 正在执行
Request 2 → 不能再次使用该账号
```

如果没有其他账号：

```text
Request 2
 ↓
FIFO
 ↓
最多等待 30s
```

## 异常

Qwen 请求失败：

```text
请求失败
 ↓
账号异常
 ↓
后续请求不再使用
```

管理员重新验证后才能恢复。

## Anthropic

```http
POST /v1/messages
```

返回：

```text
501 Not Implemented
```

---

# 30. 第一阶段产品原则

第一版重点不是把系统做“大”，而是把核心链路做通：

```text
账号管理
    ↓
账号池
    ↓
API Key
    ↓
限流
    ↓
调度
    ↓
Qwen Provider
    ↓
流式响应
    ↓
OpenAI-compatible API
```

其他能力全部围绕这个核心链路做最小实现。

最终第一版可以形成一个简单、清晰的本地 Web2API：

```text
浏览器管理账号
       ↓
创建 API Key
       ↓
任何支持 OpenAI API 的客户端
       ↓
调用 /v1/chat/completions
       ↓
Web2API
       ↓
自动选择 Qwen 账号
       ↓
获得流式回答
```

现在这份 PRD 已经足够进入**项目设计/编码阶段**了，不需要再继续把 MVP 拆成几十个小问题。下一步比较合适的是直接根据这份 PRD 确定 **Go 后端的模块职责、接口定义和完整项目目录**，然后再开始写代码。
