# Homework 03 测试报告

命令在 `LearningRoadmap/contract/templates/homework03` 下执行。

```bash
forge test -vv
forge coverage --report summary
```

Solc 0.8.38。覆盖率统计时 Foundry 会关掉 optimizer，这是 `forge coverage` 的默认行为。

## 结果

12 个测试，0 失败，0 跳过。

| 测试 | 覆盖点 |
|------|--------|
| `test_mintAndTransfer` | ERC721 铸造和 `transferFrom` |
| `test_quoteUsd_ethAndErc20` | ETH/USD 与 ERC20/USD 换成 18 位美元 |
| `test_feeBpsTiers` | 100 美元、1000 美元两档手续费边界 |
| `test_ethBidAndSettle_chargesDynamicFee` | ETH 出价、到期结算、2.5% 手续费 |
| `test_erc20OutbidsEth_thenSettles` | ERC20 以更高美元价压过 ETH，原 ETH 退回，低价再出价失败，0.5% 手续费 |
| `test_noBids_returnsNftToSeller` | 无人出价时 NFT 退回卖家 |
| `test_settleBeforeEndReverts` | 未到期不能结算 |
| `test_missingPriceFeedReverts` | 未配置预言机的代币不能报价 |
| `test_upgrade_preservesAuctionAndAllowsV2` | 非 owner 不能升级；升级后拍卖与报价仍在；v2 新函数可用；升级后仍能结算 |
| `test_adminGuardsAndBadPrice` | 非 owner 不能改 feed；手续费接收人；价格为 0 时拒绝 |
| `test_zeroDurationReverts` | 拍卖时长为 0 时拒绝 |
| `test_cannotReinitialize` | 代理和实现合约都不能再次 `initialize` |

## 覆盖率

| 文件 | Lines | Statements | Branches | Funcs |
|------|-------|------------|----------|-------|
| `src/AuctionHouseCore.sol` | 98.92% (92/93) | 90.24% (111/123) | 57.14% (16/28) | 100% (14/14) |
| `src/AuctionHouseUUPS.sol` | 100% (6/6) | 100% (3/3) | — | 100% (3/3) |
| `src/AuctionHouseV2UUPS.sol` | 100% (4/4) | 100% (2/2) | — | 100% (2/2) |
| `src/AuctionNFT.sol` | 100% (4/4) | 100% (2/2) | — | 100% (2/2) |
| `script/DeployAuctionUUPS.s.sol` | 0% | 0% | — | 0% |
| `script/UpgradeAuctionUUPS.s.sol` | 0% | 0% | — | 0% |

部署脚本没有纳入单元测试：它们只在带 RPC 和私钥的 `forge script --broadcast` 里运行。业务合约函数覆盖率是 100%。

未覆盖的分支主要是失败路径：ETH 转账失败、错误的 `msg.value`、价格小数位超过 18。这些路径在本地 mock 里不会自然发生。
