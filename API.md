# PanFlow API 文档

Base URL: `http://your-server:8080/api/v1`

所有响应格式：
```json
{"code": 0, "message": "success", "data": {...}}
```

---

## 用户端

### GET /user/parse/config
获取解析配置。

**响应**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "guest_daily_limit": 5,
    "svip_daily_limit": 100,
    "vip_count_based": true,
    "admin_unlimited": true
  }
}
```

---

### GET /user/parse/limit
获取文件大小限制。

**响应**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "max_once": 5,
    "min_single_filesize": 0,
    "max_single_filesize": 53687091200,
    "max_all_filesize": 10737418240
  }
}
```

---

### POST /user/parse/get_file_list
获取分享链接文件列表。

**请求**
```json
{
  "surl": "xxxxx",
  "pwd": "abcd"
}
```

**响应**
```json
{
  "code": 0,
  "message": "success",
  "data": {"surl": "xxxxx"}
}
```

---

### POST /user/parse/get_vcode
获取验证码（需要时）。

**请求**
```json
{
  "surl": "xxxxx",
  "pwd": "abcd"
}
```

**响应**
```json
{
  "code": 0,
  "message": "success",
  "data": {"vcode": ""}
}
```

---

### POST /user/parse/get_download_links
获取高速下载链接（核心接口）。

**请求**
```json
{
  "surl": "xxxxx",
  "pwd": "abcd",
  "fs_id": [123456789, 987654321],
  "token": "your-token"
}
```

**响应**
```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "fs_id": 123456789,
      "urls": ["https://..."],
      "size": 1048576
    }
  ]
}
```

---

### GET /user/token
查询 Token 信息。

**响应**
```json
{
  "code": 0,
  "message": "success",
  "data": null
}
```

---

### GET /user/history
查询解析历史。

**Query 参数**

| 参数 | 类型 | 说明 |
|------|------|------|
| page | int | 页码（默认 1） |
| limit | int | 每页数量（默认 10，最大 50） |

---

## 管理端

所有管理端接口（除登录外）需要在请求头中携带 JWT：

```
Authorization: Bearer <token>
```

### POST /admin/login
管理员登录，获取 JWT。

**请求**
```json
{"admin_password": "your-password"}
```

**响应**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_at": "2026-03-20 12:00:00"
  }
}
```

> JWT 默认有效期 24 小时，可通过 `hklist.jwt_expire_hours` 配置。
> 密钥通过 `hklist.jwt_secret` 配置，**生产环境必须修改默认值**。

---

### GET /admin/account
获取账号列表。

**Query 参数**

| 参数 | 类型 | 说明 |
|------|------|------|
| page | int | 页码 |
| limit | int | 每页数量（最大 100） |

**响应**
```json
{
  "code": 0,
  "message": "success",
  "data": {"total": 10, "list": [...]}
}
```

---

### POST /admin/account
添加账号。

**请求**
```json
{
  "baidu_name": "用户名",
  "uk": "123456",
  "account_type": "cookie",
  "account_data": {
    "cookie": "BDUSS=xxx; STOKEN=xxx;",
    "vip_type": "超级会员",
    "expires_at": "2026-01-01 00:00:00"
  },
  "switch": true,
  "provider_user_id": 1
}
```

---

### PATCH /admin/account
更新账号（需传 `id`）。

---

### DELETE /admin/account
删除账号。

**请求**
```json
{"id": 1}
```

---

### GET /admin/token
获取 Token 列表。

---

### POST /admin/token
创建 Token。

**请求**
```json
{
  "token": "my-token-string",
  "token_type": "normal",
  "user_type": "vip",
  "count": 100,
  "size": 10737418240,
  "can_use_ip_count": 1,
  "switch": true
}
```

---

### PATCH /admin/token
更新 Token（需传 `id`）。

---

### DELETE /admin/token
删除 Token。

**请求**
```json
{"id": 1}
```

---

### GET /admin/user
获取用户列表。

---

### POST /admin/user
创建用户。

**请求**
```json
{
  "username": "alice",
  "email": "alice@example.com",
  "user_type": "vip",
  "daily_limit": 5
}
```

---

### PATCH /admin/user
更新用户（需传 `id`）。

---

### DELETE /admin/user
删除用户。

**请求**
```json
{"id": 1}
```

---

### POST /admin/user/recharge
VIP 用户充值次数。

**请求**
```json
{
  "id": 1,
  "count": 100
}
```

---

### GET /admin/config
获取所有配置项。

**响应**
```json
{
  "code": 0,
  "message": "success",
  "data": [
    {"key": "guest_daily_limit", "value": "5", "type": "int", "description": "普通用户每日次数"},
    {"key": "svip_daily_limit", "value": "100", "type": "int", "description": "SVIP用户每日次数"}
  ]
}
```

---

### PATCH /admin/config
更新配置项。

**请求**
```json
{
  "key": "guest_daily_limit",
  "value": "10",
  "type": "int",
  "description": "普通用户每日次数"
}
```

---

### POST /admin/config/reload
重载配置缓存（使所有配置缓存失效）。

---

### GET /admin/black_list
获取黑名单列表。

---

### POST /admin/black_list
添加黑名单。

**请求**
```json
{
  "type": "ip",
  "identifier": "1.2.3.4",
  "reason": "滥用",
  "expires_at": "2026-12-31 00:00:00"
}
```

`type` 可选值：`ip` | `fingerprint`

---

### PATCH /admin/black_list
更新黑名单（需传 `id`）。

---

### DELETE /admin/black_list
删除黑名单。

**请求**
```json
{"id": 1}
```

---

### GET /admin/record
获取解析记录列表。

---

### GET /admin/record/history
按条件查询历史记录。

**Query 参数**

| 参数 | 类型 | 说明 |
|------|------|------|
| page | int | 页码 |
| limit | int | 每页数量 |
| token_id | int | 按 Token ID 过滤 |
| user_id | int | 按用户 ID 过滤 |

---

### GET /admin/proxy
获取代理列表。

---

### POST /admin/proxy
添加代理。

**请求**
```json
{
  "type": "http",
  "proxy": "http://user:pass@proxy.example.com:8080",
  "enable": true,
  "account_id": 1
}
```

---

### PATCH /admin/proxy
更新代理（需传 `id`）。

---

### DELETE /admin/proxy
删除代理。

**请求**
```json
{"id": 1}
```

---

## 错误码

详见 [ERRORS.md](./ERRORS.md)
