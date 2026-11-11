# Go 作业 04 - Gin 博客后端

对应 [homework04.md](../../homework04.md)：Gin + GORM 个人博客 API（注册登录 JWT、文章 CRUD、评论、统一错误与日志）。

## 环境

- Go 1.20+
- 本机已有 **mysql-local**（数算协同 §5）或 SQLite

## 快速开始

### 方式一：本机 mysql-local（推荐）

```bash
docker start mysql-local

docker exec -it mysql-local mysql -u root -proot123 -e \
  "CREATE DATABASE IF NOT EXISTS blog_api CHARACTER SET utf8mb4;"

cd LearningRoadmap/backend/templates/homework04
go mod tidy
go run ./cmd/server
```

默认 `http://127.0.0.1:8080`，数据库参数见 `.env.example`。

### 方式二：SQLite（无 MySQL 时）

```bash
export DB_DRIVER=sqlite
export DB_DSN=blog.db
go run ./cmd/server
```

## API 概览

| 方法 | 路径 | 认证 | 说明 |
|------|------|------|------|
| POST | `/api/register` | 否 | 注册（bcrypt 存密码） |
| POST | `/api/login` | 否 | 登录，返回 JWT |
| GET | `/api/posts` | 否 | 文章列表 |
| GET | `/api/posts/:id` | 否 | 文章详情 |
| POST | `/api/posts` | Bearer JWT | 创建文章 |
| PUT | `/api/posts/:id` | Bearer JWT | 更新（仅作者） |
| DELETE | `/api/posts/:id` | Bearer JWT | 删除（仅作者） |
| GET | `/api/posts/:id/comments` | 否 | 评论列表 |
| POST | `/api/posts/:id/comments` | Bearer JWT | 发表评论 |

请求头：`Authorization: Bearer <token>`

## 测试

```bash
go test -v ./...
```

接口测试使用内存 SQLite，**不依赖**本机 MySQL。

### Postman 示例

1. `POST /api/register`  
   `{"username":"alice","password":"secret1","email":"alice@example.com"}`
2. `POST /api/login`  
   `{"username":"alice","password":"secret1"}` → 复制 `token`
3. `POST /api/posts` + Bearer  
   `{"title":"Hello","content":"world"}`
4. `GET /api/posts/1/comments`
5. `POST /api/posts/1/comments`  
   `{"content":"nice"}`

## 目录结构

- `cmd/server`：启动入口
- `internal/models`：User / Post / Comment
- `internal/handler`：HTTP 处理
- `internal/middleware`：JWT、请求日志
- `internal/auth`：JWT 签发与校验
- `internal/database`：连接与迁移
- `internal/router`：路由注册
