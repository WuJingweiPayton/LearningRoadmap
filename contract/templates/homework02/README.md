# Contract Homework 02 (Foundry)

对应题目见 [homework02.md](../../homework02.md)。

## 合约

- `src/BeggingContract.sol`：捐赠记录、`donate` / `withdraw` / `getDonation`、`onlyOwner`
- 额外实现：`Donation` 事件、`topDonors`（前 3 名捐赠排行）

## 本地测试

```bash
forge test -vv
```

## 部署到 Sepolia / Goerli

1. 配置环境变量（示例）：

```bash
export SEPOLIA_RPC_URL="https://..."
export PRIVATE_KEY="0x..."
```

2. 部署：

```bash
forge script script/DeployBegging.s.sol:DeployBeggingScript \
  --rpc-url $SEPOLIA_RPC_URL \
  --broadcast \
  -vvvv
```

3. 将部署得到的合约地址填入下方，并附上 Etherscan / Remix 测试截图后一并提交。

## 提交清单

| 项 | 内容 |
|----|------|
| 合约代码 | 本仓库 `src/BeggingContract.sol` |
| 合约地址 | _（部署后填写）_ |
| 测试截图 | _（部署与交互后补充）_ |
