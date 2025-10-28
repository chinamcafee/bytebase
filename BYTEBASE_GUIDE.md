# Bytebase 开发指南

本文档总结了 Bytebase 项目的关键架构、构建流程和运行方式。

## 目录

- [数据库依赖](#数据库依赖)
- [前端集成方式](#前端集成方式)
- [构建流程](#构建流程)
- [运行方式](#运行方式)
- [许可证修改说明](#许可证修改说明)

---

## 数据库依赖

### Bytebase 不强依赖外部 PostgreSQL

Bytebase 支持两种数据库运行模式：

#### 1. 嵌入式数据库模式（Embedded DB）- 默认模式

- **不需要外部 PostgreSQL**
- 当 `PG_URL` 环境变量为空时自动使用
- Bytebase 会启动内置的 PostgreSQL 实例来存储元数据
- 数据存储在 `--data` 参数指定的目录中（默认为当前目录）
- 内置数据库端口为 `主端口 + 2`（如主端口 8080，数据库端口为 8082）

**启动命令示例：**
```bash
./bytebase-build/bytebase --port 8080 --data ./data
```

#### 2. 外部 PostgreSQL 模式

- 需要提供外部 PostgreSQL 连接
- 通过 `PG_URL` 环境变量指定连接字符串
- 适用于生产环境或需要高可用的场景

**启动命令示例：**
```bash
PG_URL=postgresql://user:password@localhost/dbname ./bytebase-build/bytebase --port 8080
```

### 判断逻辑

代码位置：`backend/component/config/profile.go:58-61`

```go
// UseEmbedDB returns whether to use embedDB.
func (prof *Profile) UseEmbedDB() bool {
    return len(prof.PgURL) == 0
}
```

配置读取位置：`backend/bin/server/cmd/profile.go:25`

```go
PgURL: os.Getenv("PG_URL"),
```

---

## 前端集成方式

### Bytebase 支持两种构建模式

#### 1. 嵌入式前端模式（Embedded Frontend）- 生产环境推荐

**前端会被编译并集成进二进制文件**

- 使用 Go 的 `embed` 特性将前端 `dist` 目录打包进二进制
- 构建时需要添加 `embed_frontend` build tag
- 前后端运行在同一个端口上（单体应用）
- 二进制文件包含完整的前后端功能

**关键代码：** `backend/server/server_frontend_embed.go:19-21`

```go
//go:embed dist/assets/*
//go:embed dist
var embeddedFiles embed.FS
```

#### 2. 非嵌入式模式（Development Mode）- 开发环境使用

**前端需要独立运行**

- 不使用 `embed_frontend` tag 构建
- 后端只提供 API 服务
- 前端需要单独启动开发服务器（通常在 3000 端口）
- 适合开发时热重载

**关键代码：** `backend/server/server_frontend_not_embed.go:17-18`

```go
func embedFrontend(e *echo.Echo) {
    slog.Info("Skip embedding frontend, build with 'embed_frontend' tag if you want embedded frontend.")
```

### 前端构建产物路径

**前端构建产物会自动输出到 `backend/server/dist` 目录**

在 `frontend/package.json:8` 中，`release` 脚本明确指定了输出目录：

```json
"release": "... vite build --mode release --outDir=../backend/server/dist --emptyOutDir"
```

这个命令会：
1. 构建前端项目
2. 将构建产物直接输出到 `../backend/server/dist` 目录
3. `--emptyOutDir` 会先清空目标目录

Go 的 `embed` 指令会在**编译时**自动查找 `backend/server/dist` 目录（相对于 `.go` 文件的位置），并将其内容嵌入到二进制文件中。

**不需要手动复制文件**，整个流程是自动化的。

---

## 构建流程

### 生产环境构建（带嵌入式前端）

```bash
# 1. 构建前端（自动输出到 backend/server/dist）
pnpm --dir frontend run release

# 2. 构建后端（Go embed 会自动找到 backend/server/dist 并嵌入）
go build -tags embed_frontend -ldflags "-w -s" -p=16 -o ./bytebase-build/bytebase ./backend/bin/server/main.go

# 3. 运行（单个二进制文件包含前后端）
./bytebase-build/bytebase --port 8080 --data ./data
```

### 开发环境构建

```bash
# 后端（不嵌入前端）
go build -ldflags "-w -s" -p=16 -o ./bytebase-build/bytebase ./backend/bin/server/main.go
./bytebase-build/bytebase --port 8080 --data ./data

# 前端（独立运行，带热重载）
pnpm --dir frontend dev  # 运行在 3000 端口
```

---

## 运行方式

### 使用嵌入式数据库运行（最简单）

```bash
# 直接运行，不需要任何外部依赖
./bytebase-build/bytebase --port 8080 --data ./data

# 访问管理界面
# http://localhost:8080
```

### 使用外部 PostgreSQL 运行

#### PostgreSQL 连接字符串格式

```
postgresql://[用户名]:[密码]@[主机]:[端口]/[数据库名]
```

#### 完整步骤

假设你的 PostgreSQL 配置如下：
- 用户名：`postgres`（默认）
- 密码：`19630403`
- 主机：`localhost`
- 端口：`5432`（默认）
- 数据库名：`bytebase`（需要提前创建）

```bash
# 1. 创建数据库（如果还没有）
psql -U postgres -c "CREATE DATABASE bytebase;"

# 2. 运行 Bytebase（使用外部 PostgreSQL）
PG_URL=postgresql://postgres:19630403@localhost:5432/bytebase \
  ./bytebase-build/bytebase \
  --port 8080 \
  --data ./data

# 3. 浏览器访问
# http://localhost:8080
```

或者分开写：

```bash
export PG_URL=postgresql://postgres:19630403@localhost:5432/bytebase
./bytebase-build/bytebase --port 8080 --data ./data
```

#### 其他常见配置

```bash
# 不同的用户名
PG_URL=postgresql://myuser:19630403@localhost:5432/bytebase

# 不同的主机和端口
PG_URL=postgresql://postgres:19630403@192.168.1.100:5433/bytebase

# 使用 Unix socket
PG_URL=postgresql://postgres:19630403@/bytebase?host=/var/run/postgresql
```

### 启动成功标志

启动成功后，你会看到类似这样的输出：

```
Starting Bytebase x.x.x(xxxxxxx)...
___________________________________________________________________________________________

██████╗ ██╗   ██╗████████╗███████╗██████╗  █████╗ ███████╗███████╗
██╔══██╗╚██╗ ██╔╝╚══██╔══╝██╔════╝██╔══██╗██╔══██╗██╔════╝██╔════╝
██████╔╝ ╚████╔╝    ██║   █████╗  ██████╔╝███████║███████╗█████╗
██╔══██╗  ╚██╔╝     ██║   ██╔══╝  ██╔══██╗██╔══██║╚════██║██╔══╝
██████╔╝   ██║      ██║   ███████╗██████╔╝██║  ██║███████║███████╗
╚═════╝    ╚═╝      ╚═╝   ╚══════╝╚═════╝ ╚═╝  ╚═╝╚══════╝╚══════╝

Version x.x.x(xxxxxxx) has started on port 8080 🚀
___________________________________________________________________________________________
```

然后在浏览器打开 `http://localhost:8080` 即可看到 Bytebase 的管理界面。

---

## 许可证修改说明

### 修改概述

**文件：** `backend/enterprise/license.go`

**修改目的：** 为本地开发环境强制启用企业版功能，绕过许可证验证。

### 修改内容

在 `LoadSubscription` 方法中，添加了开发模式的强制返回逻辑：

```go
// DEV MODE: Force return ENTERPRISE plan for local development
// This bypasses all license checks and unlocks all features
return &v1pb.Subscription{
    Plan:            v1pb.PlanType_ENTERPRISE,
    Seats:           -1,  // Unlimited seats
    Instances:       -1,  // Unlimited instances
    ActiveInstances: -1,  // Unlimited active instances
    ExpiresTime:     nil, // Never expires
    Trialing:        false,
}
```

### 功能说明

这个修改会：

1. **强制返回企业版订阅** - 无论是否有有效许可证
2. **解锁所有功能** - 企业版的所有高级功能都可用
3. **无限制配额**：
   - 无限座位数（Seats: -1）
   - 无限实例数（Instances: -1）
   - 无限活跃实例数（ActiveInstances: -1）
4. **永不过期** - ExpiresTime 为 nil
5. **非试用模式** - Trialing: false

### 原始逻辑

原始的许可证验证逻辑被注释掉了，包括：
- 缓存检查
- 过期时间验证
- 从数据库加载许可证
- 许可证解析和验证

### 恢复正常许可证检查

如果需要恢复正常的许可证验证流程：

1. 删除强制返回企业版的代码块
2. 取消注释原始的许可证验证逻辑

### ⚠️ 重要提示

**此修改仅用于本地开发和测试环境。**

- ❌ **不要**将此修改提交到生产环境
- ❌ **不要**将此修改推送到公共代码仓库
- ✅ 仅在本地开发时使用
- ✅ 用于测试企业版功能
- ✅ 用于开发需要企业版功能的特性

### Git 操作建议

```bash
# 将此文件添加到 .git/info/exclude（本地忽略）
echo "backend/enterprise/license.go" >> .git/info/exclude

# 或者在提交前撤销此文件的修改
git checkout backend/enterprise/license.go
```

---

## 总结

- **数据库**：Bytebase 可以使用内置 PostgreSQL（默认）或外部 PostgreSQL
- **前端**：生产环境前端会嵌入到二进制文件中，开发环境需要独立运行
- **构建**：前端构建产物自动输出到 `backend/server/dist`，Go 编译时自动嵌入
- **运行**：最简单的方式是直接运行二进制文件，访问 `http://localhost:8080`
- **许可证**：当前修改强制启用企业版功能，仅用于本地开发

---

**文档生成时间：** 2025-10-28
