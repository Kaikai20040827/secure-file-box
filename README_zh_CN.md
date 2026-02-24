# 安全文件盒

[**English**](README.md)

一个基于 Go + Gin 的 Web 应用，用于用户身份验证和加密文件存储，前端采用静态 HTML/CSS/JS。

**主要功能**

- 使用 JWT 进行用户注册/登录

- 加密文件上传/下载（AES-CTR + HMAC 完整性）

- 文件列表和删除

- 由后端提供的静态 Web UI

---

## 1. 项目布局

- `cmd/server/main.go`：应用程序入口点

- `internal/config/`：配置加载和验证

- `internal/handler/`：Gin HTTP 处理程序

- `internal/service/`：业务逻辑（文件加密代码位于此处）

- `internal/model/`：GORM 模型

- `internal/pkg/`：数据库、日志记录器、辅助函数

- `internal/routes/`：API + 静态路由

- `web/templates/`：HTML 页面

- `web/static/`：JS/CSS/图像

- `storage/`：加密文件 blob（运行时创建）

- `config.yaml`：运行时配置配置

---

## 2. 前提条件

- Go 1.18+（建议与 `go.mod` 版本匹配）

- MySQL 8+（或兼容版本）

---

## 3. 配置 (`config.yaml`)

必需字段：

- `database.*`：数据库连接参数

- `jwt.secret`：JWT 签名密钥（至少 32 个字符）

- `file_crypto.key`：**base64 URL 安全**密钥（解码后至少 32 字节）

示例（已包含在仓库中）：

```yaml
server:

host: 127.0.0.1

port: 8080

database:

driver: mysql

host: localhost

port: 3306

user: root

password: "0827"

name: secure_file_box

jwt:

issuer: secure_file_box

audience: secure_users

secret: <您的强密钥>

file_crypto:

key: <base64-url-encoded-32-bytes>

```

备注：

- 启动时，如果 `jwt.secret` 或 `file_crypto.key` 缺失或强度不足，应用程序将**自动**生成并将其写回 `config.yaml`。

- `file_crypto.key` 必须是 base64 URL 安全的密钥（无填充）。示例生成：

```bash

python - <<'PY'

import os, base64

print(base64.urlsafe_b64encode(os.urandom(32)).rstrip(b'=').decode())

PY

```

---

## 4. 数据库设置

创建数据库（模式名称必须与 `config.yaml` 中的模式名称匹配）：

```sql

CREATE DATABASE secure_file_box;

```

设置 MySQL root 密码，使其与 `config.yaml` 中的密码匹配（示例）：

```sql

ALTER USER 'root'@'localhost' IDENTIFIED BY 'yourpassword';

```

---

## 5. 运行（开发）

从仓库根目录运行：

```bash

go run ./cmd/server/main.go

```

打开：

- `http://127.0.0.1:8080`

---

## 6. 构建（生产）

```bash

go build -o ./bin/app ./cmd/server

./bin/app

```

---

## 7. API 概览

所有 API 都挂载在 `/api/v1` 目录下。

- `POST /api/v1/auth/signup`

- `POST /api/v1/auth/login`

- `GET /api/v1/user/profile`

- `PUT /api/v1/user/profile`

文件：

- `POST /api/v1/files/upload`（需要 JWT）

- `POST /api/v1/files/public/upload`（无需 JWT）

- `GET /api/v1/files`（需要 JWT）

- `GET /api/v1/files/download/:id`（需要 JWT）

- `DELETE /api/v1/files/:id`（需要 JWT）

---

## 8. 加密详情

上传的文件数据使用以下加密方式存储在 `storage/` 目录中：

- AES-CTR 用于流式加密

- HMAC-SHA256 用于完整性验证

下载时，将验证 HMAC 值。 **解密前**。如果完整性检查失败，下载将返回错误。

**重要提示：**更改 `file_crypto.key` 将导致现有文件无法读取。

---

## 9. 测试

目前尚未包含任何自动化测试。

---

## 10. 故障排除

- **MySQL 身份验证错误**：请验证 `database.user/password` 是否正确，以及数据库是否可访问。

- **文件魔数无效/完整性检查失败**：文件使用了不同的 `file_crypto.key` 加密或已损坏。

- **启动时密钥错误**：请确保 `file_crypto.key` 是有效的 base64 URL 安全密钥，并且解码后至少为 32 字节。

---

## 11. 部署说明

- 在生产环境中使用环境变量或密钥管理器。

- 在 Go 服务器前端部署 Nginx/Traefik 以启用 TLS。

- 同时备份 `storage/` 目录和数据库。

---

## 12. 贡献

在进行重大更改之前，请先提交 issue。尽量减少更改，并尽可能包含测试。
