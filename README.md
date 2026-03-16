# PanFlow

基于 Go + Gin + GORM 重构的百度网盘高速下载链接解析服务。

## 项目说明

PanFlow 是基于 Go 语言开发的高性能百度网盘直链解析工具，通过正版百度网盘 SVIP 账号代理提取高速下载链接。

### 重构目标

- 从 PHP/Laravel 迁移到 Go/Gin，提升性能和并发能力
- 保持与 HkList 的功能一致性

## 核心流程

1. 用户提交分享链接（surl）、提取码（pwd）、文件 ID 列表（fs_id[]）
2. 系统校验 Token 配额、文件大小限制
3. 随机选取一个可用 SVIP 账号
4. 调用百度官方接口将分享文件转存到该账号的「我的资源」目录
5. 调用 `locatedownload` 接口生成高速链接
6. 写入解析记录，更新账号 / Token 用量统计
7. 返回下载链接（可选经过代理服务器中转）

## 技术栈

| 组件 | 选型 |
|------|------|
| 语言 | Go 1.21+ |
| HTTP 框架 | github.com/gin-gonic/gin v1.9.1 |
| ORM | gorm.io/gorm + gorm.io/driver/mysql |
| 配置管理 | github.com/spf13/viper |
| 日志 | go.uber.org/zap |
| IP 归属地 | github.com/lionsoul2014/ip2region |
| 邮件 | gopkg.in/gomail.v2 |

## 目录结构

```
PanFlow/
├── cmd/server/main.go          # 入口文件
├── internal/
│   ├── config/                 # 配置加载
│   ├── model/                  # 数据模型
│   ├── repository/             # 数据访问层
│   ├── service/                # 业务逻辑层
│   ├── handler/                # HTTP 处理层
│   └── middleware/             # 中间件
├── pkg/
│   ├── logger/                 # 日志封装
│   └── utils/                  # 工具函数
├── config.yaml                 # 配置文件
├── config.example.yaml         # 配置示例
├── ip2region.xdb              # IP 数据库
├── claude.md                   # 项目架构文档
├── errors.md                   # 错误码定义
└── go.mod
```

## 快速开始

### 环境要求

- Go 1.21+
- MySQL 5.7+

### 安装步骤

1. 克隆仓库

```bash
git clone https://github.com/cg8-5712/PanFlow.git
cd PanFlow
```

2. 安装依赖

```bash
go mod download
```

3. 配置数据库

复制配置文件并修改数据库连接：

```bash
cp config.example.yaml config.yaml
```

编辑 `config.yaml`：

```yaml
database:
  host: 127.0.0.1
  port: 3306
  user: root
  password: your_password
  dbname: panflow
```

4. 运行服务

```bash
go run cmd/server/main.go
```

服务默认运行在 `http://localhost:8080`

### Docker 部署

```bash
docker build -t panflow .
docker run -d -p 8080:8080 \
  -e DB_HOST=your_db_host \
  -e DB_DATABASE=panflow \
  -e DB_USERNAME=root \
  -e DB_PASSWORD=your_password \
  panflow
```

## 配置说明

主要配置项：

```yaml
server:
  port: 8080
  mode: release  # debug/release

database:
  host: 127.0.0.1
  port: 3306
  user: root
  password: ""
  dbname: panflow

hklist:
  admin_password: admin      # 后台管理密码
  parse_password: admin      # 解析密码（留空则不需要）

log:
  level: info               # debug/info/warn/error
```

完整配置请参考 `config.example.yaml`。

## API 文档

### 路由前缀

所有 API 路由前缀为 `/api/v1`

### 用户端接口

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /user/parse/config | 获取解析配置 |
| GET | /user/parse/limit | 获取限制信息 |
| POST | /user/parse/get_file_list | 获取文件列表 |
| POST | /user/parse/get_vcode | 获取验证码 |
| POST | /user/parse/get_download_links | 获取下载链接 |
| GET | /user/token | 查询 Token 信息 |
| GET | /user/history | 查询解析历史 |

### 管理端接口

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /admin/check_password | 校验管理密码 |
| GET/POST/PATCH/DELETE | /admin/account | 账号管理 |
| GET/POST/PATCH/DELETE | /admin/token | Token 管理 |
| GET/POST/PATCH/DELETE | /admin/black_list | 黑名单管理 |
| GET | /admin/record | 解析记录 |
| GET | /admin/record/history | 历史记录 |
| GET/POST/PATCH/DELETE | /admin/proxy | 代理管理 |

详细错误码定义请查看 [errors.md](./errors.md)

## 数据库表结构

主要数据表：

- `accounts` - SVIP 账号管理
- `tokens` - Token/卡密管理
- `file_lists` - 文件索引缓存
- `records` - 解析记录
- `black_lists` - IP/指纹黑名单
- `proxies` - 代理服务器配置

详细表结构请查看 [claude.md](./claude.md)

## 开发指南

### 初始化顺序

1. 加载配置 (`viper`)
2. 初始化日志 (`zap`)
3. 连接数据库并自动迁移
4. 初始化 Repository 层
5. 初始化 Service 层
6. 设置 Gin 模式
7. 注册路由和中间件
8. 启动 HTTP 服务

### 代码规范

- 遵循 Go 官方代码规范
- 使用 `gofmt` 格式化代码
- 分层架构：Handler → Service → Repository
- 统一错误处理和响应格式


## 许可证

本项目仅供学习参考使用。 cc by nc sa 4.0

## 免责声明

项目所涉及的接口均为官方开放接口，需使用正版 SVIP 会员账号进行代理提取高速链接，无破坏官方接口行为，本身不存在违法。

仅供自己参考学习使用。若违规使用，官方会限制或封禁你的账号，包括你的 IP；如无官方授权进行商业用途会对你造成更严重后果。

源码仅供学习，如无视声明使用产生正负面结果（限速、被封等）均与作者无关。
