# PanFlow 开发计划

## 项目概述

PanFlow 是 Go/Gin 的百度网盘高速下载链接解析服务。本文档记录完整的开发计划、进度追踪和技术决策。

---

## 开发阶段

### 阶段 1：基础架构搭建 ✅

**目标：** 建立项目基础框架和核心依赖

- [x] 初始化 Go 项目结构
- [x] 配置 Gin 框架
- [x] 集成 GORM ORM
- [x] 配置 Viper 配置管理
- [x] 集成 Zap 日志系统
- [x] 设计分层架构（Handler → Service → Repository）
- [x] 编写项目文档（README、CLAUDE、ERRORS）

**完成时间：** 已完成

---

### 阶段 2：数据库设计与模型定义 ✅

**目标：** 完成数据库表结构设计和 GORM 模型定义

- [x] 所有 8 张表 GORM 模型定义（accounts/tokens/users/configs/file_lists/records/black_lists/proxies）
- [x] JSONMap / JSONSlice 自定义序列化
- [x] 软删除支持（accounts/tokens/users）
- [x] provider_user_id、user_type 字段
- [x] AutoMigrate + seedGuestToken + seedDefaultConfigs
- [x] 全部 Repository 层（db/account/token/user/config/record/file_list/black_list/proxy）
- [x] L1 缓存（Otter，变长 TTL，10000 容量）
- [x] L2 缓存（Redis，go-redis/v9）
- [x] 多级缓存策略（CacheGet/CacheSet/CacheDelete）

**完成时间：** 已完成

---

### 阶段 3：中间件开发 ✅

**目标：** 实现请求过滤和权限控制中间件

- [x] IdentifierFilter：IP + 指纹黑名单检查，L2 缓存加速，debug 放行
- [x] PassFilterAdmin：Header/Query/Body 三路密码校验
- [x] PassFilterUser：解析密码校验
- [x] CORS 中间件：预检请求支持

**完成时间：** 已完成

---

### 阶段 4：百度网盘 API 封装 ✅

**目标：** 封装百度网盘官方 API 调用

- [x] `service/bdwp.go` - HTTP 客户端封装（支持代理）
- [x] GetShareInfo - 获取分享链接信息（shareid/uk/bdstoken）
- [x] GetFileList - 获取分享文件列表
- [x] TransferFiles - 转存文件到「我的资源」
- [x] LocateDownload - 获取高速下载直链
- [x] DeleteFile - 删除网盘文件

**完成时间：** 已完成

---

### 阶段 5：解析服务开发 ✅

**目标：** 实现核心解析业务逻辑

- [x] `service/token.go` - Token 校验（次数/大小/每日/IP/过期）+ 缓存
- [x] `service/user.go` - 用户配额校验（guest/vip/svip/admin）+ VIP 扣减
- [x] `service/account.go` - 账号选择（SVIP 专属账号 / 公共随机池）
- [x] `service/record.go` - 解析记录保存
- [x] `service/config.go` - 配置读取（L1 缓存 + 热更新）
- [x] `service/cache.go` - 多级缓存协调（L1 Otter + L2 Redis）
- [x] `service/parse.go` - 完整解析流程编排
- [x] `handler/parse.go` - 解析 API 处理器
- [x] `internal/router/router.go` - 全路由注册

**完成时间：** 已完成
- [ ] `service/parse.go` - 解析服务
- [ ] `service/token.go` - Token 管理服务
- [ ] `service/user.go` - 用户管理服务（新增）
- [ ] `service/record.go` - 记录管理服务
- [ ] `service/cache.go` - 缓存管理服务（新增）
  - [ ] Moka 内存缓存封装
  - [ ] Redis 缓存封装
  - [ ] 多级缓存协调
- [ ] 编写单元测试
- [ ] 编写集成测试

**预计完成：** 第 5-6 周

---

### 阶段 6：用户端 API 开发 ✅

**目标：** 实现用户端接口

- [x] GET `/api/v1/user/parse/config`
- [x] GET `/api/v1/user/parse/limit`
- [x] POST `/api/v1/user/parse/get_file_list`
- [x] POST `/api/v1/user/parse/get_vcode`
- [x] POST `/api/v1/user/parse/get_download_links`
- [x] GET `/api/v1/user/token`
- [x] GET `/api/v1/user/history`
- [x] `handler/response.go` 统一响应格式

**完成时间：** 已完成

---

### 阶段 7：管理端 API 开发 ✅

**目标：** 实现管理端接口

- [x] 账号 CRUD（GET/POST/PATCH/DELETE `/api/v1/admin/account`）
- [x] Token CRUD（GET/POST/PATCH/DELETE `/api/v1/admin/token`）
- [x] 用户 CRUD + 充值（GET/POST/PATCH/DELETE `/api/v1/admin/user`，POST `/api/v1/admin/user/recharge`）
- [x] 配置管理（GET/PATCH `/api/v1/admin/config`，POST `/api/v1/admin/config/reload`）
- [x] 黑名单 CRUD（GET/POST/PATCH/DELETE `/api/v1/admin/black_list`）
- [x] 记录查询（GET `/api/v1/admin/record`，GET `/api/v1/admin/record/history`）
- [x] 代理 CRUD（GET/POST/PATCH/DELETE `/api/v1/admin/proxy`）
- [x] POST `/api/v1/admin/check_password`

**完成时间：** 已完成

---

### 阶段 8：辅助功能开发 ✅

**目标：** 实现邮件、代理等辅助功能

- [x] `service/mail.go` - SMTP 邮件服务（账号异常通知、解析失败通知）
- [x] `service/proxy.go` - 代理服务（账号级代理选择、HTTP 客户端构建、URL 包装）

**完成时间：** 已完成

---

### 阶段 9：测试与优化 🔄

**目标：** 完善测试覆盖率和性能优化

- [x] `test/service/service_test.go` - Service 层单元测试（24 个用例，全部通过）
- [x] `test/handler/handler_test.go` - Handler 层单元测试（9 个用例，全部通过）
- [x] `test/middleware/middleware_test.go` - 中间件单元测试（8 个用例，全部通过）
- [x] `go vet ./...` - 零警告
- [x] `gofmt -w` - 全部文件格式化
- **测试总计：41 个用例，全部通过**
- [ ] Repository 层测试（需要 testcontainers 或 mock DB）
- [ ] 数据库连接池配置优化

**预计完成：** 进行中

---

### 阶段 10：部署与文档 ✅

**目标：** 完成部署配置和文档编写

- [x] `Dockerfile` - 多阶段构建（golang:1.24-alpine → alpine:3.19，含 ca-certificates + tzdata）
- [x] `docker-compose.yml` - app + MySQL 8.0 + Redis 7（含 healthcheck）
- [x] `API.md` - 完整 RESTful API 文档（用户端 + 管理端，含请求/响应示例）
- [x] `DEPLOY.md` - 部署指南（Docker + 裸机 + systemd + Nginx 反向代理）
- [x] `MIGRATION.md` - 从 HkList 迁移指南（配置映射 + SQL + API 路径对比）

**完成时间：** 已完成

---

## 技术选型与决策

### 核心依赖

| 依赖 | 版本 | 用途 |
|------|------|------|
| Go | 1.21+ | 编程语言 |
| Gin | v1.9.1 | HTTP 框架 |
| GORM | v1.25+ | ORM 框架 |
| Viper | v1.18+ | 配置管理 |
| Zap | v1.26+ | 日志系统 |
| ip2region | v3.0+ | IP 归属地 |
| gomail | v2 | 邮件发送 |
| Moka | latest | 内存缓存（L1） |
| go-redis | v9+ | Redis 客户端（L2） |

### 架构设计

```
┌─────────────────────────────────────────┐
│           HTTP Request (Gin)            │
└─────────────────┬───────────────────────┘
                  │
         ┌────────▼────────┐
         │   Middleware    │
         │  - CORS         │
         │  - PassFilter   │
         │  - Identifier   │
         └────────┬────────┘
                  │
         ┌────────▼────────┐
         │     Handler     │
         │  - 参数校验      │
         │  - 错误处理      │
         └────────┬────────┘
                  │
         ┌────────▼────────┐
         │     Service     │
         │  - 业务逻辑      │
         │  - 事务管理      │
         │  - 缓存协调      │
         └────────┬────────┘
                  │
         ┌────────▼────────┐
         │   Cache Layer   │
         │  L1: Moka       │
         │  L2: Redis      │
         └────────┬────────┘
                  │
         ┌────────▼────────┐
         │   Repository    │
         │  - 数据访问      │
         │  - SQL 封装      │
         └────────┬────────┘
                  │
         ┌────────▼────────┐
         │   Database      │
         │    (MySQL)      │
         └─────────────────┘
```

### 用户类型与权限

| 用户类型 | 每日次数 | 计费方式 | 账号选择 | 说明 |
|---------|---------|---------|---------|------|
| **Guest** | 5次（可配置） | 免费 | 公共账号池 | 游客模式 |
| **VIP** | 按充值次数 | 充值购买 | 公共账号池 | 用完即止 |
| **SVIP** | 100次（可配置） | 提供自己的百度 SVIP 号 | 仅自己的账号 | 使用自己的账号 |
| **Admin** | 无限制 | - | 所有账号 | 最高权限 |

### 百度网盘账号类型

| 账号类型 | 认证方式 | 说明 |
|---------|---------|------|
| **Cookie** | BDUSS + STOKEN | 最常用，从浏览器获取 |
| **Open Platform** | OAuth Token | 官方开放平台，支持自动刷新 |
| **Enterprise Cookie** | 企业 CID + Token | 企业网盘账号 |
| **Download Ticket** | 下载凭证 | 特殊下载模式 |

### 性能目标

- 单机 QPS > 1000
- 平均响应时间 < 200ms
- P99 响应时间 < 500ms
- 数据库连接池 > 100
- 并发 Goroutine > 10000

---

## 开发规范

### 代码规范

1. **命名规范**
   - 包名：小写，单个单词
   - 文件名：小写，下划线分隔
   - 变量名：驼峰命名
   - 常量名：大写，下划线分隔

2. **注释规范**
   - 所有导出函数必须有注释
   - 复杂逻辑必须有注释说明
   - 使用 godoc 格式

3. **错误处理**
   - 不忽略任何错误
   - 使用统一的错误码
   - 记录详细的错误日志

4. **测试规范**
   - 测试文件以 `_test.go` 结尾
   - 测试函数以 `Test` 开头
   - 使用 table-driven tests

### Git 规范

**分支策略**
- `main` - 主分支，稳定版本
- `develop` - 开发分支
- `feature/*` - 功能分支
- `bugfix/*` - 修复分支

**提交规范**
```
<type>(<scope>): <subject>

<body>

<footer>
```

Type 类型：
- `feat`: 新功能
- `fix`: 修复 bug
- `docs`: 文档更新
- `style`: 代码格式
- `refactor`: 重构
- `test`: 测试
- `chore`: 构建/工具

---

## 风险与挑战

### 技术风险

1. **百度 API 变更**
   - 风险：百度可能随时调整 API
   - 应对：抽象 API 层，便于快速适配

2. **并发性能**
   - 风险：高并发下性能瓶颈
   - 应对：压力测试，连接池优化

3. **数据一致性**
   - 风险：统计数据不一致
   - 应对：使用事务，定期校验

### 业务风险

1. **账号封禁**
   - 风险：SVIP 账号被封禁
   - 应对：自动检测，及时通知

2. **配额耗尽**
   - 风险：Token 配额快速耗尽
   - 应对：限流策略，配额预警

---

## 里程碑

| 里程碑 | 目标 | 预计时间 |
|--------|------|----------|
| M1 | 基础架构完成 | 已完成 |
| M2 | 数据库和模型完成 | 第 2 周 |
| M3 | 中间件和 API 封装完成 | 第 4 周 |
| M4 | 核心解析功能完成 | 第 6 周 |
| M5 | 用户端 API 完成 | 第 7 周 |
| M6 | 管理端 API 完成 | 第 9 周 |
| M7 | 测试和优化完成 | 第 12 周 |
| M8 | 部署和文档完成 | 第 13 周 |

---

## 参考资源

- [HkList PHP 版本](https://github.com/HkList/HkList)
- [Gin 官方文档](https://gin-gonic.com/docs/)
- [GORM 官方文档](https://gorm.io/docs/)
- [Go 编码规范](https://go.dev/doc/effective_go)
- [百度网盘开放平台](https://pan.baidu.com/union/doc/)

---

## 更新日志

| 日期 | 版本 | 更新内容 |
|------|------|----------|
| 2026-03-16 | v1.0 | 初始版本，完成开发计划编写 |
