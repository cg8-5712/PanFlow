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

### 阶段 2：数据库设计与模型定义 🔄

**目标：** 完成数据库表结构设计和 GORM 模型定义

#### 2.1 数据表设计

- [ ] `accounts` 表 - SVIP 账号管理
  - [ ] 定义 GORM 模型
  - [ ] 实现 JSON 字段序列化（account_data）
  - [ ] 添加软删除支持
  - [ ] 添加 `provider_user_id` 字段（记录提供账号的用户）
  - [ ] 编写数据库迁移

- [ ] `tokens` 表 - Token/卡密管理
  - [ ] 定义 GORM 模型
  - [ ] 添加 `user_type` 字段（guest/vip/svip/admin）
  - [ ] 添加 `provider_user_id` 字段（SVIP 用户关联）
  - [ ] 实现 guest token 自动初始化
  - [ ] 添加唯一索引

- [ ] `users` 表 - 用户管理（新增）
  - [ ] 定义 GORM 模型
  - [ ] 用户类型：guest/vip/svip/admin
  - [ ] VIP 用户充值次数管理
  - [ ] SVIP 用户百度账号绑定

- [ ] `configs` 表 - 数据库配置（新增）
  - [ ] 定义 GORM 模型
  - [ ] key-value 结构
  - [ ] 支持 string/int/bool/json 类型
  - [ ] 配置热更新机制

- [ ] `file_lists` 表 - 文件索引缓存
  - [ ] 定义 GORM 模型
  - [ ] 添加 fs_id 唯一索引

- [ ] `records` 表 - 解析记录
  - [ ] 定义 GORM 模型
  - [ ] 添加外键关联
  - [ ] 添加 `user_id` 字段

- [ ] `black_lists` 表 - 黑名单管理
  - [ ] 定义 GORM 模型
  - [ ] 实现过期时间自动清理

- [ ] `proxies` 表 - 代理服务器配置
  - [ ] 定义 GORM 模型
  - [ ] 添加账号关联

#### 2.2 Repository 层实现

- [ ] `repository/db.go` - 数据库连接和初始化
- [ ] `repository/account.go` - 账号数据访问
- [ ] `repository/token.go` - Token 数据访问
- [ ] `repository/user.go` - 用户数据访问（新增）
- [ ] `repository/config.go` - 配置数据访问（新增）
- [ ] `repository/record.go` - 记录数据访问
- [ ] `repository/file_list.go` - 文件列表数据访问
- [ ] `repository/black_list.go` - 黑名单数据访问
- [ ] `repository/proxy.go` - 代理数据访问

#### 2.3 缓存层设计（新增）

- [ ] 集成 Moka 内存缓存
  - [ ] L1 缓存：热点配置、Token 信息
  - [ ] TTL: 1-5 分钟
  - [ ] LRU 淘汰策略
  - [ ] 容量限制配置

- [ ] 集成 Redis 分布式缓存
  - [ ] L2 缓存：Token、账号信息、配置
  - [ ] TTL: 10-30 分钟
  - [ ] 支持集群模式
  - [ ] 缓存预热机制

- [ ] 多级缓存策略
  - [ ] L1 Miss → L2 查询
  - [ ] L2 Miss → 数据库查询
  - [ ] 缓存更新策略（Write-Through）
  - [ ] 缓存失效策略

**预计完成：** 第 1-2 周

---

### 阶段 3：中间件开发 🔄

**目标：** 实现请求过滤和权限控制中间件

#### 3.1 IdentifierFilter 中间件

- [ ] 实现 IP 黑名单检查
- [ ] 实现浏览器指纹黑名单检查
- [ ] 集成 ip2region IP 归属地查询
- [ ] Debug 模式放行逻辑
- [ ] 编写单元测试

#### 3.2 PassFilter 中间件

- [ ] 实现 ADMIN 模式密码校验
  - [ ] Header 参数获取
  - [ ] Query 参数获取
  - [ ] Body 参数获取
- [ ] 实现 USER 模式密码校验
- [ ] 编写单元测试

#### 3.3 CORS 中间件

- [ ] 配置跨域策略
- [ ] 支持预检请求

**预计完成：** 第 2 周

---

### 阶段 4：百度网盘 API 封装 🔄

**目标：** 封装百度网盘官方 API 调用

#### 4.1 账号类型支持

- [ ] Cookie 类型账号
  - [ ] 获取账户信息
  - [ ] 获取 VIP 类型和到期时间
  - [ ] 检查账号封禁状态

- [ ] Open Platform 类型账号
  - [ ] AccessToken 换取
  - [ ] RefreshToken 刷新
  - [ ] 获取 VIP 信息

- [ ] Enterprise Cookie 类型账号
  - [ ] 企业账号信息获取
  - [ ] BDSToken 获取
  - [ ] Dlink Cookie 处理

- [ ] Download Ticket 类型账号
  - [ ] 下载卷模式支持

#### 4.2 核心 API 封装

- [ ] 获取分享链接信息
- [ ] 获取文件列表
- [ ] 获取验证码
- [ ] 转存文件到网盘
- [ ] 获取下载直链（locatedownload）
- [ ] 删除网盘文件

#### 4.3 Service 层实现

- [ ] `service/bdwp.go` - 百度网盘 API 封装
- [ ] `service/account.go` - 账号管理服务
- [ ] 实现账号随机选择算法
- [ ] 实现账号使用统计更新
- [ ] 编写单元测试

**预计完成：** 第 3-4 周

---

### 阶段 5：解析服务开发 🔄

**目标：** 实现核心解析业务逻辑

#### 5.1 解析流程实现

- [ ] 用户类型识别与权限校验（新增）
  - [ ] Guest 用户：5次/天（可配置）
  - [ ] VIP 用户：按充值次数扣减
  - [ ] SVIP 用户：100次/天（可配置），仅使用自己的账号
  - [ ] Admin 用户：无限制

- [ ] Token 配额校验
  - [ ] 次数限制检查
  - [ ] 大小限制检查
  - [ ] 每日限额检查
  - [ ] IP 限制检查
  - [ ] 用户类型限制检查（新增）

- [ ] 文件大小限制校验
  - [ ] 单文件大小限制
  - [ ] 总大小限制

- [ ] 账号选择逻辑
  - [ ] 可用账号筛选
  - [ ] SVIP 用户仅选择自己提供的账号（新增）
  - [ ] 随机选择算法
  - [ ] 负载均衡

- [ ] 文件转存流程
  - [ ] 分享链接解析
  - [ ] 文件转存到账号
  - [ ] 转存失败重试

- [ ] 下载链接生成
  - [ ] 调用 locatedownload API
  - [ ] 限速检测
  - [ ] 代理服务器中转

- [ ] 解析记录保存
  - [ ] 记录解析历史
  - [ ] 更新账号统计
  - [ ] 更新 Token 统计
  - [ ] 更新用户统计（新增）

#### 5.2 Service 层实现

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

### 阶段 6：用户端 API 开发 📝

**目标：** 实现用户端接口

#### 6.1 解析相关接口

- [ ] GET `/api/v1/user/parse/config` - 获取解析配置
- [ ] GET `/api/v1/user/parse/limit` - 获取限制信息
- [ ] POST `/api/v1/user/parse/get_file_list` - 获取文件列表
- [ ] POST `/api/v1/user/parse/get_vcode` - 获取验证码
- [ ] POST `/api/v1/user/parse/get_download_links` - 获取下载链接

#### 6.2 Token 和历史接口

- [ ] GET `/api/v1/user/token` - 查询 Token 信息
- [ ] GET `/api/v1/user/history` - 查询解析历史

#### 6.3 Handler 层实现

- [ ] `handler/parse.go` - 解析处理器
- [ ] `handler/response.go` - 统一响应格式
- [ ] 参数校验
- [ ] 错误处理
- [ ] 编写 API 测试

**预计完成：** 第 7 周

---

### 阶段 7：管理端 API 开发 📝

**目标：** 实现管理端接口

#### 7.1 账号管理

- [ ] GET `/api/v1/admin/account` - 获取账号列表
- [ ] POST `/api/v1/admin/account` - 添加账号
- [ ] PATCH `/api/v1/admin/account` - 更新账号
- [ ] DELETE `/api/v1/admin/account` - 删除账号

#### 7.2 Token 管理

- [ ] GET `/api/v1/admin/token` - 获取 Token 列表
- [ ] POST `/api/v1/admin/token` - 创建 Token
- [ ] PATCH `/api/v1/admin/token` - 更新 Token
- [ ] DELETE `/api/v1/admin/token` - 删除 Token

#### 7.3 用户管理（新增）

- [ ] GET `/api/v1/admin/user` - 获取用户列表
- [ ] POST `/api/v1/admin/user` - 创建用户
- [ ] PATCH `/api/v1/admin/user` - 更新用户
- [ ] DELETE `/api/v1/admin/user` - 删除用户
- [ ] POST `/api/v1/admin/user/recharge` - VIP 用户充值

#### 7.4 配置管理（新增）

- [ ] GET `/api/v1/admin/config` - 获取配置列表
- [ ] PATCH `/api/v1/admin/config` - 更新配置
- [ ] POST `/api/v1/admin/config/reload` - 重载配置缓存

#### 7.5 黑名单管理

- [ ] GET `/api/v1/admin/black_list` - 获取黑名单列表
- [ ] POST `/api/v1/admin/black_list` - 添加黑名单
- [ ] PATCH `/api/v1/admin/black_list` - 更新黑名单
- [ ] DELETE `/api/v1/admin/black_list` - 删除黑名单

#### 7.4 记录管理

- [ ] GET `/api/v1/admin/record` - 获取解析记录
- [ ] GET `/api/v1/admin/record/history` - 获取历史记录

#### 7.5 代理管理

- [ ] GET `/api/v1/admin/proxy` - 获取代理列表
- [ ] POST `/api/v1/admin/proxy` - 添加代理
- [ ] PATCH `/api/v1/admin/proxy` - 更新代理
- [ ] DELETE `/api/v1/admin/proxy` - 删除代理

#### 7.6 其他管理接口

- [ ] POST `/api/v1/admin/check_password` - 校验管理密码
- [ ] POST `/api/v1/install` - 系统安装

#### 7.7 Handler 层实现

- [ ] `handler/account.go` - 账号管理处理器
- [ ] `handler/token.go` - Token 管理处理器
- [ ] `handler/black_list.go` - 黑名单管理处理器
- [ ] `handler/record.go` - 记录管理处理器
- [ ] `handler/proxy.go` - 代理管理处理器
- [ ] `handler/config.go` - 配置管理处理器
- [ ] 编写 API 测试

**预计完成：** 第 8-9 周

---

### 阶段 8：辅助功能开发 📝

**目标：** 实现邮件、代理等辅助功能

#### 8.1 邮件服务

- [ ] `service/mail.go` - 邮件发送服务
- [ ] 配置 SMTP 服务器
- [ ] 实现错误通知邮件
- [ ] 实现账号异常通知

#### 8.2 代理服务

- [ ] `service/proxy.go` - 代理服务
- [ ] 支持 HTTP/HTTPS 代理
- [ ] 支持 SOCKS5 代理
- [ ] 代理健康检查
- [ ] 代理自动切换

#### 8.3 配置管理

- [ ] `service/config.go` - 配置管理服务
- [ ] 动态配置更新
- [ ] 配置持久化

**预计完成：** 第 10 周

---

### 阶段 9：测试与优化 🧪

**目标：** 完善测试覆盖率和性能优化

#### 9.1 单元测试

- [ ] Repository 层测试覆盖率 > 80%
- [ ] Service 层测试覆盖率 > 80%
- [ ] Handler 层测试覆盖率 > 70%
- [ ] 中间件测试覆盖率 > 90%

#### 9.2 集成测试

- [ ] 完整解析流程测试
- [ ] 账号管理流程测试
- [ ] Token 管理流程测试
- [ ] 黑名单功能测试

#### 9.3 性能测试

- [ ] 压力测试（并发 1000+）
- [ ] 数据库查询优化
- [ ] 缓存策略优化
- [ ] 连接池配置优化

#### 9.4 代码质量

- [ ] 使用 golangci-lint 静态检查
- [ ] 代码格式化（gofmt）
- [ ] 代码注释完善
- [ ] 错误处理规范化

**预计完成：** 第 11-12 周

---

### 阶段 10：部署与文档 🚀

**目标：** 完成部署配置和文档编写

#### 10.1 Docker 支持

- [ ] 编写 Dockerfile
- [ ] 编写 docker-compose.yml
- [ ] 多阶段构建优化
- [ ] 镜像大小优化

#### 10.2 部署文档

- [ ] 编写部署指南
- [ ] 编写配置说明
- [ ] 编写运维手册
- [ ] 编写故障排查指南

#### 10.3 API 文档

- [ ] 使用 Swagger 生成 API 文档
- [ ] 编写接口调用示例
- [ ] 编写错误码说明

#### 10.4 迁移指南

- [ ] 编写从 HkList 迁移指南
- [ ] 数据迁移脚本
- [ ] 配置迁移说明

**预计完成：** 第 13 周

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
