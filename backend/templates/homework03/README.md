# Go 作业 03 - 进阶 GORM

对应 [homework03.md](../../homework03.md)。

## 题目对应

| 题目 | 文件 |
|------|------|
| 模型与建表 | `models.go` → `MigrateModels` |
| 关联查询 | `queries.go` |
| 钩子函数 | `Post.AfterCreate`、`Comment.AfterDelete` |
| MySQL 连接 | `db.go` → `OpenMySQLFromEnv` |

## 本机 MySQL（复用已有容器）

**不要为作业单独 `docker run` 新容器。** 使用数算协同开发环境里的 `mysql-local`（见  
`ChinaDataValley/shusuanxietong/.../dataxcompute/docs/本地开发环境启动指南/本地开发环境完整指南.md` §5）。

```bash
docker start mysql-local
```

在 `mysql-local` 中创建作业库（只需一次）：

```bash
docker exec -it mysql-local mysql -u root -proot123 -e \
  "CREATE DATABASE IF NOT EXISTS blog_homework CHARACTER SET utf8mb4;"
```

连接参数见 `.env.example`（默认 `appuser` / `apppass123` @ `127.0.0.1:3306`）。账号与本地 `.env` 不一致时，用环境变量或 `MYSQL_DSN` 覆盖。

## 运行

单元测试（内存 SQLite，无需 MySQL）：

```bash
go mod tidy
go test -v ./...
```

连本机 MySQL 的演示：

```bash
docker start mysql-local
go run .
```

根目录 `main.go` 为演示入口；与 `models.go` 等同为 `package main`（Go 要求同目录只能一个包）。

## 提交

推送后 GitHub Actions 会执行 `go test`（见 `.github/workflows/homework03.yml`）。
