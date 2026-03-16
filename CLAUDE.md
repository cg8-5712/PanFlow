# PanFlow 项目概述

**PHP/Laravel 版本（HkList） → Go 版本（PanFlow）重构**
**架构对齐：** F:/bian/goproject/Easy-Stream

---

## 1. 项目功能

PanFlow 通过正版百度网盘 SVIP 账号代理提取高速下载链接。

### 核心流程

1. 用户提交分享链接（surl）、提取码（pwd）、文件 ID 列表（fs_id[]）
2. 系统校验 token 配额、文件大小限制
3. 随机选取一个可用 SVIP 账号
4. 调用百度官方接口将分享文件转存到该账号的「我的资源」目录
5. 调用 locatedownload 接口生成高速链接
6. 写入解析记录，更新账号/token 用量统计
7. 返回下载链接（可选经过代理服务器中转）

---

## 2. 技术栈

| 组件 | 选型 |
|------|------|
| HTTP 框架 | github.com/gin-gonic/gin v1.9.1 |
| ORM | gorm.io/gorm + gorm.io/driver/mysql |
| 配置 | github.com/spf13/viper（YAML + 数据库） |
| 日志 | go.uber.org/zap（封装在 pkg/logger） |
| IP 归属地 | github.com/lionsoul2014/ip2region/binding/golang |
| 邮件 | gopkg.in/gomail.v2 |
| 缓存 L1 | Moka（内存缓存） |
| 缓存 L2 | go-redis/redis（分布式缓存） |

---

## 3. 目录结构

```
PanFlow/
├── cmd/server/main.go
├── internal/
│   ├── config/config.go
│   ├── model/model.go
│   ├── repository/
│   │   ├── db.go
│   │   ├── account.go
│   │   ├── token.go
│   │   ├── user.go          # 新增：用户管理
│   │   ├── config.go        # 新增：配置管理
│   │   ├── record.go
│   │   ├── file_list.go
│   │   ├── black_list.go
│   │   └── proxy.go
│   ├── service/
│   │   ├── bdwp.go
│   │   ├── parse.go
│   │   ├── account.go
│   │   ├── token.go
│   │   ├── user.go          # 新增：用户服务
│   │   ├── cache.go         # 新增：缓存服务
│   │   ├── record.go
│   │   ├── black_list.go
│   │   ├── proxy.go
│   │   ├── config.go
│   │   └── mail.go
│   ├── handler/
│   │   ├── parse.go
│   │   ├── account.go
│   │   ├── token.go
│   │   ├── user.go          # 新增：用户处理器
│   │   ├── record.go
│   │   ├── black_list.go
│   │   ├── proxy.go
│   │   ├── config.go
│   │   └── response.go
│   └── middleware/
│       ├── pass_filter.go
│       ├── identifier_filter.go
│       └── cors.go
├── pkg/
│   ├── logger/logger.go
│   ├── cache/               # 新增：缓存封装
│   │   ├── moka.go
│   │   └── redis.go
│   └── utils/utils.go
├── ip2region.xdb
├── config.yaml
├── config.example.yaml
└── go.mod
```

---

## 4. main.go 初始化顺序

1. viper config.Load()
2. logger.Init(cfg.Log.Level)
3. cache.InitMoka()  // 初始化 L1 缓存
4. cache.InitRedis(cfg.Redis)  // 初始化 L2 缓存
5. repository.NewDB(cfg.Database)  // AutoMigrate + seed guest token
6. 初始化各 Repository（注入 *gorm.DB）
7. 初始化各 Service（注入 repo + cache）
8. 从数据库加载配置到缓存
9. gin.SetMode
10. r := gin.Default()
11. r.Use(middleware.Cors())
12. 注册路由
13. r.Run(addr)

---

## 5. 数据库表结构

### accounts

```sql
CREATE TABLE accounts (
  id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  baidu_name VARCHAR(255),
  uk VARCHAR(255),
  account_type VARCHAR(50),
  account_data TEXT,
  `switch` TINYINT(1) DEFAULT 1,
  reason VARCHAR(500),
  prov VARCHAR(100),
  provider_user_id BIGINT,  -- 新增：提供账号的用户ID
  used_count BIGINT DEFAULT 0,
  used_size BIGINT DEFAULT 0,
  total_size BIGINT DEFAULT 0,
  total_size_updated_at DATETIME,
  last_use_at DATETIME,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  INDEX idx_provider_user (provider_user_id)
);
```

#### account_data JSON 结构

**cookie 类型**

```json
{
  "cookie": "BDUSS=xxx; STOKEN=xxx;",
  "vip_type": "超级会员",
  "expires_at": "2025-01-01 00:00:00"
}
```

vip_type 可能值：
- 超级会员
- 普通会员
- 普通用户

**open_platform 类型**

```json
{
  "access_token": "xxx",
  "refresh_token": "xxx",
  "token_expires_at": "2025-01-01 00:00:00",
  "vip_type": "超级会员",
  "expires_at": "2025-01-01 00:00:00"
}
```

**enterprise_cookie 类型**

```json
{
  "cookie": "xxx",
  "cid": 123,
  "expires_at": "2025-01-01 00:00:00",
  "bdstoken": "xxx",
  "dlink_cookie": "xxx"
}
```

**download_ticket 类型**

```json
{
  "surl": "xxx",
  "pwd": "xxx",
  "dir": "/",
  "cid": 123,
  "save_cookie": "xxx",
  "save_bdstoken": "xxx",
  "download_cookie": "xxx",
  "download_bdstoken": "xxx"
}
```

---

### tokens

```sql
CREATE TABLE tokens (
  id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  token VARCHAR(255) UNIQUE,
  token_type VARCHAR(20),
  user_type ENUM('guest', 'vip', 'svip', 'admin') DEFAULT 'guest',  -- 新增：用户类型
  provider_user_id BIGINT,  -- 新增：SVIP用户关联
  count BIGINT DEFAULT 0,
  size BIGINT DEFAULT 0,
  day BIGINT DEFAULT 0,
  used_count BIGINT DEFAULT 0,
  used_size BIGINT DEFAULT 0,
  can_use_ip_count BIGINT DEFAULT 1,
  ip TEXT,
  `switch` TINYINT(1) DEFAULT 1,
  reason VARCHAR(500),
  expires_at DATETIME,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  INDEX idx_user_type (user_type),
  INDEX idx_provider_user (provider_user_id)
);
```

#### guest token 默认值

- token=guest
- token_type=daily
- user_type=guest
- count=10
- size=10GB
- day=1
- can_use_ip_count=99999

---

### users（新增）

```sql
CREATE TABLE users (
  id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  username VARCHAR(255) UNIQUE NOT NULL,
  email VARCHAR(255),
  user_type ENUM('guest', 'vip', 'svip', 'admin') DEFAULT 'guest',
  vip_balance BIGINT DEFAULT 0,  -- VIP剩余次数
  daily_used_count BIGINT DEFAULT 0,  -- 今日已用次数
  daily_limit INT DEFAULT 5,  -- 每日限额
  baidu_account_id BIGINT,  -- SVIP用户绑定的百度账号ID
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  INDEX idx_user_type (user_type),
  INDEX idx_username (username)
);
```

---

### configs（新增）

```sql
CREATE TABLE configs (
  id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  `key` VARCHAR(255) UNIQUE NOT NULL,
  value TEXT,
  type ENUM('string', 'int', 'bool', 'json') DEFAULT 'string',
  description VARCHAR(500),
  created_at DATETIME,
  updated_at DATETIME,
  INDEX idx_key (`key`)
);
```

#### 默认配置项

```sql
INSERT INTO configs (`key`, value, type, description) VALUES
  ('guest_daily_limit', '5', 'int', '普通用户每日次数'),
  ('vip_count_based', 'true', 'bool', 'VIP按次数计费'),
  ('svip_daily_limit', '100', 'int', 'SVIP用户每日次数'),
  ('admin_unlimited', 'true', 'bool', 'Admin无限制');
```

---

### file_lists
CREATE TABLE file_lists (
  id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  surl VARCHAR(255),
  pwd VARCHAR(100),
  fs_id VARCHAR(255) UNIQUE,
  size BIGINT,
  filename VARCHAR(500),
  created_at DATETIME,
  updated_at DATETIME
);
records
### records

```sql
CREATE TABLE records (
  id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  ip VARCHAR(100),
  fingerprint VARCHAR(255),
  fs_id BIGINT UNSIGNED,
  urls TEXT,
  ua VARCHAR(500),
  token_id BIGINT UNSIGNED,
  account_id BIGINT UNSIGNED,
  user_id BIGINT UNSIGNED,  -- 新增：用户ID
  created_at DATETIME,
  updated_at DATETIME,
  INDEX idx_user (user_id),
  INDEX idx_token (token_id),
  INDEX idx_account (account_id)
);
```

---

### black_lists

```sql
CREATE TABLE black_lists (
  id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  type VARCHAR(20),
  identifier VARCHAR(255),
  reason VARCHAR(500),
  expires_at DATETIME,
  created_at DATETIME,
  updated_at DATETIME,
  INDEX idx_type_identifier (type, identifier)
);
```

---

### proxies

```sql
CREATE TABLE proxies (
  id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  type VARCHAR(20),
  proxy VARCHAR(500),
  enable TINYINT(1) DEFAULT 1,
  reason VARCHAR(500),
  account_id BIGINT UNSIGNED,
  created_at DATETIME,
  updated_at DATETIME,
  INDEX idx_account (account_id)
);
```

---

## 6. 多级缓存架构

### 缓存层级

```
┌─────────────────────────────────────┐
│         Application Layer           │
└──────────────┬──────────────────────┘
               │
        ┌──────▼──────┐
        │   L1 Cache  │  Moka (本地内存)
        │   - 配置项   │  - TTL: 1-5分钟
        │   - 热点数据 │  - LRU淘汰
        └──────┬──────┘  - 容量限制
               │ Miss
        ┌──────▼──────┐
        │   L2 Cache  │  Redis (分布式)
        │   - Token   │  - TTL: 10-30分钟
        │   - 用户信息 │  - 支持集群
        │   - 账号信息 │  - 持久化
        └──────┬──────┘
               │ Miss
        ┌──────▼──────┐
        │   Database  │  MySQL
        │   - 持久化   │
        └─────────────┘
```

### 缓存策略

#### L1 缓存（Moka）

**缓存内容：**
- 系统配置（configs 表）
- 热点 Token 信息
- 用户每日配额统计

**配置：**
- TTL: 1-5 分钟
- 最大容量: 10000 条
- 淘汰策略: LRU

#### L2 缓存（Redis）

**缓存内容：**
- Token 详细信息
- 用户信息
- 账号信息
- 黑名单列表

**配置：**
- TTL: 10-30 分钟
- 支持集群模式
- 持久化: RDB + AOF

### 缓存更新策略

1. **Write-Through（写穿）**
   - 更新数据时同时更新缓存
   - 保证数据一致性

2. **Cache-Aside（旁路）**
   - 读取时先查缓存
   - 缓存未命中时查数据库并回填

3. **失效策略**
   - 配置更新：立即失效所有节点 L1+L2
   - 用户信息更新：失效对应 L2 缓存
   - Token 更新：失效对应 L1+L2 缓存

---

## 7. 用户类型与权限

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

---

## 8. 路由与中间件
## 8. 路由与中间件

### 路由表（前缀 /api/v1）

#### 公开路由

- POST /install（已移除）

#### 用户端路由

**Middleware：IdentifierFilter**

- GET  /user/parse/config
- GET  /user/parse/limit
- POST /user/parse/get_file_list
- POST /user/parse/get_vcode
- POST /user/parse/get_download_links
- GET  /user/token
- GET  /user/history

#### 管理端路由

**Middleware：PassFilter:ADMIN**

- POST   /admin/check_password

**账号管理**
- GET    /admin/account
- POST   /admin/account
- PATCH  /admin/account
- DELETE /admin/account

**Token 管理**
- GET    /admin/token
- POST   /admin/token
- PATCH  /admin/token
- DELETE /admin/token

**用户管理（新增）**
- GET    /admin/user
- POST   /admin/user
- PATCH  /admin/user
- DELETE /admin/user
- POST   /admin/user/recharge

**配置管理（新增）**
- GET    /admin/config
- PATCH  /admin/config
- POST   /admin/config/reload

**黑名单管理**
- GET    /admin/black_list
- POST   /admin/black_list
- PATCH  /admin/black_list
- DELETE /admin/black_list

**记录管理**
- GET    /admin/record
- GET    /admin/record/history

**代理管理**
- GET    /admin/proxy
- POST   /admin/proxy
- PATCH  /admin/proxy
- DELETE /admin/proxy

---

## 9. 中间件
## 9. 中间件

### IdentifierFilter

**文件：** `internal/middleware/identifier_filter.go`

**逻辑：**

1. debug 模式直接放行
2. 获取客户端 IP
3. 查询 black_lists 表（优先查 L2 缓存）
4. 若命中 ip 黑名单 → 返回 code=20014
5. 获取浏览器指纹 rand2
6. 若命中 fingerprint 黑名单 → 返回 code=20014
7. 否则放行

---

### PassFilter

**文件：** `internal/middleware/pass_filter.go`

#### ADMIN 模式

从以下位置获取密码：
- Header[admin_password]
- Query[admin_password]
- Body[admin_password]

**校验：** cfg.Hklist.AdminPassword

#### USER 模式

从以下位置获取密码：
- Query[parse_password]
- Body[parse_password]

**校验：** cfg.Hklist.ParsePassword

---

## 10. 统一响应格式

**文件：** `internal/handler/response.go`

```go
type Response struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data"`
}
```

### 成功响应

```go
func Success(c *gin.Context, data interface{})
```

### 失败响应

```go
func Fail(c *gin.Context, httpStatus, code int, msg string)
```

详细错误码定义请查看 [ERRORS.md](./ERRORS.md)