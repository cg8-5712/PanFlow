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
| 配置 | github.com/spf13/viper（YAML + 环境变量） |
| 日志 | go.uber.org/zap（封装在 pkg/logger） |
| IP 归属地 | github.com/lionsoul2014/ip2region/binding/golang |
| 邮件 | gopkg.in/gomail.v2 |

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
│   │   ├── record.go
│   │   ├── file_list.go
│   │   ├── black_list.go
│   │   └── proxy.go
│   ├── service/
│   │   ├── bdwp.go
│   │   ├── parse.go
│   │   ├── account.go
│   │   ├── token.go
│   │   ├── record.go
│   │   ├── black_list.go
│   │   ├── proxy.go
│   │   ├── config.go
│   │   └── mail.go
│   ├── handler/
│   │   ├── parse.go
│   │   ├── account.go
│   │   ├── token.go
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
3. repository.NewDB(cfg.Database)  // AutoMigrate + seed guest token
4. 初始化各 Repository（注入 *gorm.DB）
5. 初始化各 Service（注入 repo）
6. gin.SetMode
7. r := gin.Default()
8. r.Use(middleware.Cors())
9. 注册路由
10. r.Run(addr)

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
  used_count BIGINT DEFAULT 0,
  used_size BIGINT DEFAULT 0,
  total_size BIGINT DEFAULT 0,
  total_size_updated_at DATETIME,
  last_use_at DATETIME,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME
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
  id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  token VARCHAR(255) UNIQUE,
  token_type VARCHAR(20),
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
  deleted_at DATETIME
);
guest token 默认值
token=guest
token_type=daily
count=10
size=10GB
day=1
can_use_ip_count=99999
file_lists
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
CREATE TABLE records (
  id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  ip VARCHAR(100),
  fingerprint VARCHAR(255),
  fs_id BIGINT UNSIGNED,
  urls TEXT,
  ua VARCHAR(500),
  token_id BIGINT UNSIGNED,
  account_id BIGINT UNSIGNED,
  created_at DATETIME,
  updated_at DATETIME
);
black_lists
CREATE TABLE black_lists (
  id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  type VARCHAR(20),
  identifier VARCHAR(255),
  reason VARCHAR(500),
  expires_at DATETIME,
  created_at DATETIME,
  updated_at DATETIME
);
proxies
CREATE TABLE proxies (
  id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  type VARCHAR(20),
  proxy VARCHAR(500),
  enable TINYINT(1) DEFAULT 1,
  reason VARCHAR(500),
  account_id BIGINT UNSIGNED,
  created_at DATETIME,
  updated_at DATETIME
);
路由与中间件
路由表（前缀 /api/v1）
公开路由
POST /install
用户端路由

Middleware：IdentifierFilter

GET  /user/parse/config
GET  /user/parse/limit
POST /user/parse/get_file_list
POST /user/parse/get_vcode
POST /user/parse/get_download_links
GET  /user/token
GET  /user/history
管理端路由

Middleware：PassFilter:ADMIN

POST   /admin/check_password

GET    /admin/account
POST   /admin/account
PATCH  /admin/account
DELETE /admin/account

GET    /admin/token
POST   /admin/token
PATCH  /admin/token
DELETE /admin/token

GET    /admin/black_list
POST   /admin/black_list
PATCH  /admin/black_list
DELETE /admin/black_list

GET    /admin/record
GET    /admin/record/history

GET    /admin/proxy
POST   /admin/proxy
PATCH  /admin/proxy
DELETE /admin/proxy
中间件
IdentifierFilter

文件：

internal/middleware/identifier_filter.go

逻辑：

1. debug 模式直接放行
2. 获取客户端 IP
3. 查询 black_lists 表
4. 若命中 ip 黑名单
   返回 code=20014
5. 获取浏览器指纹 rand2
6. 若命中 fingerprint 黑名单
   返回 code=20014
7. 否则放行
PassFilter

文件：

internal/middleware/pass_filter.go
ADMIN
Header[admin_password]
Query[admin_password]
Body[admin_password]

校验：

cfg.Hklist.AdminPassword
USER
Query[parse_password]
Body[parse_password]

校验：

cfg.Hklist.ParsePassword
统一响应格式

文件：

internal/handler/response.go
type Response struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data"`
}
成功
func Success(c *gin.Context, data interface{})
失败
func Fail(c *gin.Context, httpStatus, code int, msg string)