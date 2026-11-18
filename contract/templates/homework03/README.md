# Contract Homework 03 (Foundry)

对应题目见 [homework03.md](../../homework03.md)。

题目原文写的是 Hardhat。这个目录用 **Foundry** 实现同一份功能：ERC721、ETH/ERC20 拍卖、Chainlink 价格换算成美元、**UUPS 与透明代理两种升级路径**、按成交额分档的手续费。

## 项目结构

| 路径 | 作用 |
|------|------|
| `src/AuctionNFT.sol` | ERC721，owner 铸造，转移走 ERC721 标准 |
| `src/AuctionHouseCore.sol` | 拍卖逻辑（UUPS / 透明代理共用） |
| `src/AuctionHouseUUPS.sol` | UUPS 实现合约 v1 |
| `src/AuctionHouseTransparent.sol` | 透明代理实现合约 v1（逻辑在实现里，升级走 ProxyAdmin） |
| `src/AuctionHouseV2UUPS.sol` | UUPS 升级实现 v2 |
| `src/AuctionHouseV2Transparent.sol` | 透明代理升级实现 v2 |
| `src/interfaces/AggregatorV3Interface.sol` | Chainlink `latestRoundData` 接口（与 `@chainlink/contracts` v1.4.0 一致） |
| `script/DeployAuctionUUPS.s.sol` | 部署 NFT + **UUPS**（`ERC1967Proxy`） |
| `script/UpgradeAuctionUUPS.s.sol` | **UUPS** 升级到 v2（`upgradeToAndCall`，需 proxy owner） |
| `script/DeployAuctionTransparent.s.sol` | 部署 NFT + **透明代理**（`ProxyAdmin` + `TransparentUpgradeableProxy`） |
| `script/UpgradeAuctionTransparent.s.sol` | **透明代理** 升级到 v2（`ProxyAdmin.upgradeAndCall`） |
| `test/NftAuction.t.sol` | 单元测试 |

用户交互的是**代理地址**，不是实现合约地址。

## 功能

1. **上架** `createAuction(nft, tokenId, duration)`  
   卖家需先 `approve` 拍卖代理。NFT 会转到拍卖合约。
2. **出价**  
   - `bidWithEth(auctionId)`：附带 ETH  
   - `bidWithErc20(auctionId, token, amount)`：先 `approve`，再划转 ERC20  
   新出价必须比当前最高价的**美元价值**更高。被超过的出价会原路退回。
3. **结算** `settle(auctionId)`  
   到期后任何人都可以调用。有出价时，NFT 给出价最高者，资金（扣除手续费）给卖家。无人出价则 NFT 退回卖家。
4. **美元报价** `quoteUsd(asset, amount)`  
   `address(0)` 表示 ETH。价格来自 Chainlink feed：`amount * price / 10^assetDecimals`，再换到 18 位小数的美元。
5. **动态手续费** `feeBps`  
   - 低于 100 美元：2.5%  
   - 100 美元（含）到低于 1000 美元：1.5%  
   - 1000 美元及以上：0.5%
6. **升级**  
   - **UUPS**：只有 proxy owner 可以对代理调用 `upgradeToAndCall`。  
   - **透明代理**：升级由 `ProxyAdmin` 执行；业务 owner 仍为实现合约里的 `onlyOwner`。  
   两种路径的 v2 都增加 `implementationVersion() == 2` 和 `pingAfterUpgrade()`。

## 本地测试

在本目录执行（`forge` 需在 `PATH` 中，见本机 `~/.foundry/bin`）：

```bash
forge test -vv
forge coverage --report summary
```

最近一次结果见 [TEST_REPORT.md](./TEST_REPORT.md)。`AuctionHouseUUPS.sol`（及 Core）行覆盖率见报告，13 个测试全部通过。

## 部署到 Sepolia

1. 复制环境变量并填入自己的私钥和 RPC：

```bash
cp .env.example .env
```

2. 部署（二选一）。`ETH_USD_FEED` 默认是 Chainlink Sepolia ETH/USD，部署前请到 [Chainlink Data Feeds](https://docs.chain.link/data-feeds/price-feeds/addresses) 核对地址。

**UUPS：**

```bash
source .env
forge script script/DeployAuctionUUPS.s.sol:DeployAuctionUUPSScript \
  --rpc-url "$SEPOLIA_RPC_URL" \
  --private-key "$PRIVATE_KEY" \
  --broadcast
```

**透明代理：**

```bash
forge script script/DeployAuctionTransparent.s.sol:DeployAuctionTransparentScript \
  --rpc-url "$SEPOLIA_RPC_URL" \
  --private-key "$PRIVATE_KEY" \
  --broadcast
```

3. 把日志里的 **proxy** 写入 `.env` 的 `AUCTION_PROXY`；若走透明代理，同时写入 **ProxyAdmin** 到 `PROXY_ADMIN`。升级：

```bash
# UUPS
forge script script/UpgradeAuctionUUPS.s.sol:UpgradeAuctionUUPSScript \
  --rpc-url "$SEPOLIA_RPC_URL" \
  --private-key "$PRIVATE_KEY" \
  --broadcast

# 透明代理
forge script script/UpgradeAuctionTransparent.s.sol:UpgradeAuctionTransparentScript \
  --rpc-url "$SEPOLIA_RPC_URL" \
  --private-key "$PRIVATE_KEY" \
  --broadcast
```

ERC20 出价前，owner 要对那个代币调用 `setPriceFeed(token, aggregator)`。

## 提交时还要你本机补的内容

| 项 | 状态 |
|----|------|
| 合约与测试 | 已完成，本地 `forge test` 通过 |
| 测试报告 | [TEST_REPORT.md](./TEST_REPORT.md) |
| 测试网合约地址 | 部署后把 proxy、NFT 地址填到这里 |
| 测试网截图 | 部署和一次出价/结算后自行补充 |
