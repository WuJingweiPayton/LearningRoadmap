# Backend Homework 05 — Sepolia 读写与 abigen

对应题目见 [homework05.md](../../homework05.md)。

## 项目结构

| 路径 | 作用 |
|------|------|
| `contract/Counter.sol` | 任务 2 计数器合约 |
| `artifacts/` | Foundry 编译产物 + `Counter.abi` / `Counter.bin` |
| `bindings/counter.go` | **abigen** 生成的 Go 绑定（已提交，可 `scripts/gen-bindings.sh` 再生） |
| `internal/chain/` | ethclient 连接、查块、转账 |
| `cmd/queryblock/` | 任务 1：查询区块 |
| `cmd/sendtx/` | 任务 1：签名并发送 ETH 转账 |
| `cmd/counter/` | 任务 2：部署 / `increment` / 读 `count` |

## 环境

```bash
cp .env.example .env
# SEPOLIA_RPC_URL、PRIVATE_KEY 等
```

需本机：`go 1.22+`、`forge`（改合约时）、`abigen`（`go install github.com/ethereum/go-ethereum/cmd/abigen@latest`）。

## 任务 1：查询区块

```bash
cd backend/templates/homework05
source .env   # 或 export 变量
go run ./cmd/queryblock
# 指定高度：BLOCK_NUMBER=12345678 go run ./cmd/queryblock
```

输出：区块哈希、时间戳、交易数量。

## 任务 1：发送转账

```bash
# 默认 TRANSFER_ETH=0.0001
go run ./cmd/sendtx
```

输出：交易哈希；等待打包后打印区块号与 status。

## 任务 2：合约与 abigen

重新生成绑定（修改 `Counter.sol` 后）：

```bash
./scripts/gen-bindings.sh
```

链上交互：

```bash
# 部署
COUNTER_ACTION=deploy go run ./cmd/counter

# 读 count
COUNTER_ACTION=read COUNTER_ADDRESS=0x... go run ./cmd/counter

# increment
COUNTER_ACTION=increment COUNTER_ADDRESS=0x... go run ./cmd/counter
```

## 测试

```bash
go test ./...
```

`go test` 仅跑本地单元测试；连 Sepolia 需自行用上述 `go run` 验证。