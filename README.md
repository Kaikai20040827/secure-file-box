## 中文教程
# 1. 项目概述

本项目是一个用于加密文件和上传的 Web 应用，采用 Go 作为后端，前端使用静态资源（HTML/CSS/JS）。功能包括用户注册/登录、文件上传/下载、时间表查看与管理、基本权限校验和健康检查接口。项目目录结构组织为后端服务、路由、处理器、服务层与静态前端资源。

# 2. 目录结构（重要文件）

- `config.yaml`：应用配置文件（端口、数据库等）。
- `server/main.go`：程序入口。
- `internal/`：内部包（配置、处理器、服务、模型、middleware）。
- `api/swagger/`：API 文档与 Swagger 定义（若有）(正在编写)。
- `web/static/`：前端静态文件（JS、CSS、images）。
- `web/templates/`：页面模板文件（HTML）。
- `tests/`：包含单元与集成测试目录。

# 3. 前置要求

- Go 1.18+（或本项目指定的最低 Go 版本）。 => https://go.dev/doc/install
- 可选：Node.js（若需对前端资源进行打包/构建）, 本项目最低运行没要求使用nodejs。
- 数据库：请参考 `config.yaml` 中的配置（可能使用本地或远程数据库）。

# 4. 安装与构建

1) 克隆仓库：
- Https:
```
  git clone https://github.com/Kaikai20040827/graduation.git
```

- SSH:
```
  git clone git@github.com:Kaikai20040827/graduation.git
```

- Github CLI:
```
  gh repo clone Kaikai20040827/graduation
```

2) 进入项目根目录并设置配置文件 `config.yaml`（根据环境调整数据库、端口、密码等）。
jwt.secret 为 base64 编码加密

3) 下载数据库: MySQL Workbench => https://dev.mysql.com/downloads/workbench/
创建一个名字为 "secure_file_box" 的数据库

4) 构建后端：
先在config.yaml的password那一行内容换成自己的数据库密码
``` 
  go build -o ./bin/app ./server
``` 

  或直接运行：
``` 
  go run ./server/main.go
```

# 5. 运行

- 本地运行（开发）：

  1. 确保 `config.yaml` 配置正确。
  2. 执行 `go run ./cmd/server/main.go` 或运行编译后的可执行文件 `./bin`。
  3. 访问浏览器：`http://localhost:<port>`（端口在 `config.yaml` 中配置），默认8080端口。

# 6. 配置说明 🔧

- 请编辑根目录的 `config.yaml`，主要配置项通常包括：服务端口、数据库连接字符串、JWT 密钥、上传目录等。
- JWT 密钥：请在 `config.yaml` 的 `jwt.secret` 字段中设置一个强随机密钥，并避免将其提交到版本控制。

# 7. 使用说明（主要功能）

- 用户：注册、登录、查看与修改个人信息。
- 时间表：上传/导入时间表文件、查看日/周/月视图。
- 文件：支持文件上传、存储管理及下载接口。
- 管理接口：健康检查、系统信息与日志查看（按项目实现）。

# 8. API 文档 (正在开发)

- 项目包含 `api/swagger/` 目录，如存在请使用 Swagger UI 打开对应的 JSON/YAML 文档，或联系开发者生成最新 API 文档。

# 9. 测试

- 运行所有测试 (尚未开发)：
```
  go test .
```
- 单元测试目录：`tests/unit/`，集成测试目录：`tests/integration/`。请根据测试说明准备依赖（如测试数据库）。

# 10. 部署建议

- 在生产环境中：
  - 使用配置管理（例如环境变量或 secret 管理）替换明文配置。
  - 使用进程管理工具或容器化（Docker）运行服务，并配合反向代理（如 Nginx）和 TLS。
  - 配置日志与监控、定期备份数据库。

# 11. 贡献指南

- 欢迎提 PR：请先创建 issue 描述问题或功能，再提交分支与 PR。遵循代码风格、补充测试，并保持提交信息清晰。

# 12. 联系方式

- 如需帮助或有疑问，请联系项目维护者或在仓库中创建 issue。

---

