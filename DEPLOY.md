# PanFlow 部署指南

## 环境要求

| 组件 | 最低版本 |
|------|---------|
| Go | 1.21+ |
| MySQL | 5.7+ / 8.0 |
| Redis | 6.0+ |
| OS | Linux / Windows / macOS |

---

## 方式一：Docker 部署（推荐）

### 1. 克隆项目

```bash
git clone https://github.com/yourname/panflow.git
cd panflow
```

### 2. 创建配置文件

```bash
cp config.example.yaml config.yaml
```

编辑 `config.yaml`，填写数据库和 Redis 连接信息：

```yaml
database:
  host: "mysql"      # docker-compose 服务名
  port: "3306"
  user: "panflow"
  password: "panflow"
  name: "panflow"

redis:
  host: "redis"      # docker-compose 服务名
  port: "6379"
  password: ""
  db: 0

hklist:
  admin_password: "your-admin-password"
```

### 3. 启动服务

```bash
docker compose up -d
```

### 4. 验证

```bash
curl http://localhost:8080/ping
# {"message":"pong"}
```

---

## 方式二：裸机部署

### 1. 编译

```bash
go build -ldflags="-s -w" -o panflow ./cmd/server
```

### 2. 配置

```bash
cp config.example.yaml config.yaml
# 编辑 config.yaml 填写实际连接信息
```

### 3. 运行

```bash
./panflow
```

### 4. 使用 systemd 管理（Linux）

创建 `/etc/systemd/system/panflow.service`：

```ini
[Unit]
Description=PanFlow Service
After=network.target mysql.service redis.service

[Service]
Type=simple
WorkingDirectory=/opt/panflow
ExecStart=/opt/panflow/panflow
Restart=on-failure
RestartSec=5s
User=www-data

[Install]
WantedBy=multi-user.target
```

```bash
systemctl daemon-reload
systemctl enable panflow
systemctl start panflow
systemctl status panflow
```

---

## 配置说明

### server

| 字段 | 默认值 | 说明 |
|------|--------|------|
| host | 0.0.0.0 | 监听地址 |
| port | 8080 | 监听端口 |
| mode | release | gin 模式（debug/release） |

### database

| 字段 | 默认值 | 说明 |
|------|--------|------|
| host | 127.0.0.1 | MySQL 地址 |
| port | 3306 | MySQL 端口 |
| user | root | 用户名 |
| password | - | 密码 |
| name | panflow | 数据库名 |

### redis

| 字段 | 默认值 | 说明 |
|------|--------|------|
| host | 127.0.0.1 | Redis 地址 |
| port | 6379 | Redis 端口 |
| password | - | 密码（无则留空） |
| db | 0 | 数据库编号 |

### hklist（核心配置）

| 字段 | 说明 |
|------|------|
| admin_password | 管理端密码 |
| parse_password | 解析端密码（留空则不校验） |
| debug | 调试模式（跳过黑名单检查） |
| max_once | 单次最多解析文件数 |
| max_single_filesize | 单文件最大字节数 |
| max_all_filesize | 单次总文件最大字节数 |
| proxy_enable | 是否启用全局代理 |
| proxy_http | HTTP 代理地址 |
| mail_switch | 是否启用邮件通知 |

---

## 反向代理（Nginx）

```nginx
server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}
```

---

## 故障排查

**服务无法启动**
- 检查 `config.yaml` 是否存在且格式正确
- 检查 MySQL / Redis 是否可连接
- 查看日志输出

**数据库连接失败**
- 确认 MySQL 用户有 `panflow` 数据库的权限
- 确认防火墙未拦截 3306 端口

**Redis 连接失败**
- Redis 不可用时服务仍可启动，L2 缓存降级为不可用
- 检查 Redis 密码配置是否正确
