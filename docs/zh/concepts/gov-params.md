# 治理参数

在IRISnet中，存在一些特殊的参数，它们可通过链上治理被修改。持有IRIS通证的用户都可以参与到参数修改的链上治理。
如果社区对某些可修改的参数不满意，可以发起[参数修改提案](../features/governance.md#usage-scenario-of-parameter-change)，社区投票通过后即可在线自动完成修改。

## Auth 模块可治理参数

| 字段                     | 描述             | 有效范围            | 当前值        |
| ------------------------ | ---------------- | ------------------- | ------------- |
| `auth/gasPriceThreshold` | 最小的gas单价    | (0, 10^18iris-atto] | 6000000000000 |
| `auth/txSizeLimit`       | 普通交易大小限制 | [500, 1500]         | 1000          |

## Stake 模块可治理参数

| 字段                  | 描述             | 有效范围    | 当前值   |
| --------------------- | ---------------- | ----------- | -------- |
| `stake/MaxValidators` | 验证人的最大数量 | [100, 200]  | 100      |
| `stake/UnbondingTime` | 解绑时间         | [2week, +∞) | 504h0m0s |

详见 [Staking](../features/stake.md)

## Distribution 模块可治理参数

| 字段                        | 描述                 | 有效范围  | 当前值 |
| --------------------------- | -------------------- | --------- | ------ |
| `distr/BaseProposerReward`  | 出块奖励的基准比例   | (0, 0.02] | 0.01   |
| `distr/BonusProposerReward` | 最大额外奖励比例     | (0, 0.08] | 0.04   |
| `distr/CommunityTax`        | 贡献给社区基金的比例 | (0, 0.2]  | 0.02   |

详见 [Distribution](../features/distribution.md)

## Mint 模块可治理参数

| 字段             | 描述     | 有效范围 | 当前值 |
| ---------------- | -------- | -------- | ------ |
| `mint/Inflation` | 通胀系数 | [0, 0.2] | 0.04   |

详见 [Mint](../features/mint.md)

## Slashing 模块可治理参数

| 字段                               | 描述                       | 有效范围      | 当前值  |
| ---------------------------------- | -------------------------- | ------------- | ------- |
| `slashing/CensorshipJailDuration`  | Censorship后Jail的时间     | (0, 4week)    | 48h0m0s |
| `slashing/DoubleSignJailDuration`  | DoubleSign后Jail的时间     | (0, 4week)    | 48h0m0s |
| `slashing/DowntimeJailDuration`    | Downtime后Jail的时间       | (0, 4week)    | 36h0m0s |
| `slashing/MaxEvidenceAge`          | 可接受的最早的作恶证据时间 | [1day, +∞)    | 51840   |
| `slashing/MinSignedPerWindow`      | slash窗口中最小投票比例    | [0.5, 0.9]    | 0.7     |
| `slashing/SignedBlocksWindow`      | slash统计窗口区块数        | [100, 140000] | 34560   |
| `slashing/SlashFractionCensorship` | Censorship后罚款的比例     | [0.005, 0.1]  | 0       |
| `slashing/SlashFractionDoubleSign` | DoubleSign后罚款的比例     | [0.01, 0.1]   | 0.01    |
| `slashing/SlashFractionDowntime`   | Downtime后罚款的比例       | [0.005, 0.1]  | 0.0003  |

详见 [Slashing](../features/slashing.md)

## Asset 模块可治理参数

| 字段                         | 描述                                           | 有效范围 | 当前值                            |
| ---------------------------- | ---------------------------------------------- | -------- | --------------------------------- |
| `asset/AssetTaxRate`         | 资产税率，即Community Tax的比例                | [0, 1]   | 0.4                               |
| `asset/IssueTokenBaseFee`    | 发行Token的基准费用                            | [0, +∞)  | 60000000000000000000000iris-atto  |
| `asset/MintTokenFeeRatio`    | 增发Token的费率(相对于发行费用)                | [0, 1]   | 0.1                               |
| `asset/CreateGatewayBaseFee` | 创建网关的基准费用                             | [0, +∞)  | 120000000000000000000000iris-atto |
| `asset/GatewayAssetFeeRatio` | 发行网关资产的费率(相对于native资产的发行费用) | [0, 1]   | 0.1                               |

详见 [Asset](../features/asset.md)

## Service 模块可治理参数

| 字段                           | 描述                       | 有效范围         | 当前值   |
| ------------------------------ | -------------------------- | ---------------- | -------- |
| `service/ArbitrationTimeLimit` | 争议解决最大时长           | [5days, 10days]  | 120h0m0s |
| `service/ComplaintRetrospect`  | 可提起争议最大时长         | [15days, 30days] | 360h0m0s |
| `service/MaxRequestTimeout`    | 服务调用最大等待区块个数   | [20, +∞)         | 100      |
| `service/MinDepositMultiple`   | 服务绑定最小抵押金额的倍数 | [500, 5000]      | 1000     |
| `service/ServiceFeeTax`        | 服务费的税收比例           | (0, 0.2]         | 0.01     |
| `service/SlashFraction`        | 惩罚百分比                 | (0, 0.01]        | 0.001    |
| `service/TxSizeLimit`          | service类型交易的大小限制  | [2000, 6000]     | 4000     |

详见 [Service](../features/service.md)
