# 治理参数

在IRISnet中，存在一些特殊的参数，它们可通过链上治理被修改。持有IRIS通证的用户都可以参与到参数修改的链上治理。
如果社区对某些可修改的参数不满意，可以发起[参数修改提案](../features/governance.md#usage-scenario-of-parameter-change)，社区投票通过后即可在线自动完成修改。

## Auth 模块可治理参数

| 字段                          | 描述                               | 有效范围                  | 当前值 |
| ----------------------------- | ---------------------------------- | ------------------------- | ------ |
| `auth/MaxMemoCharacters`      | 交易的memo字段的最大字符数         | (0, 18446744073709551615] | 256    |
| `auth/TxSigLimit`             | 每个交易的最大签名数               | (0, 18446744073709551615] | 7      |
| `auth/TxSizeCostPerByte`      | 交易每个字节消耗的gas量            | (0, 18446744073709551615] | 10     |
| `auth/SigVerifyCostED25519`   | 在ED25519算法签名验证上花费的gas   | (0, 18446744073709551615] | 590    |
| `auth/SigVerifyCostSecp256k1` | 在Secp256k1算法签名验证上花费的gas | (0, 18446744073709551615] | 1000   |

## Bank 模块可治理参数

| 字段                      | 描述                    | 有效范围     | 当前值 |
| ------------------------- | ----------------------- | ------------ | ------ |
| `bank/SendEnabled`        | 支持transfer的token集合 |              | []     |
| `bank/DefaultSendEnabled` | 默认是否支持转账功能    | {true,false} | true   |

详见 [Bank](../features/bank.md)

## Coinswap  模块可治理参数

| 字段           | 描述     | 有效范围 | 当前值               |
| -------------- | -------- | -------- | -------------------- |
| `coinswap/Fee` | 交换费用 | (0,1)    | 0.003000000000000000 |

详见 [Coinswap](../features/coinswap.md)

## Distribution 模块可治理参数

| 字段                               | 描述                   | 有效范围     | 当前值 |
| ---------------------------------- | ---------------------- | ------------ | ------ |
| `distribution/communitytax`        | 提现收取的手续费率     | [0, 1]       | 0.02   |
| `distribution/baseproposerreward`  | 区块提议者的基础奖励率 | [0, 1]       | 0.01   |
| `distribution/bonusproposerreward` | 区块提议者的奖励率     | [0, 1]       | 0.04   |
| `distribution/withdrawaddrenabled` | 是否支持设置提现地址   | {true,false} | true   |

详见 [Distribution](../features/distribution.md)

## Gov 模块可治理参数

| 字段                | 描述                   | 有效范围                                                 | 当前值                                                                                                         |
| ------------------- | ---------------------- | -------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------- |
| `gov/depositparams` | 提议抵押阶段的相关参数 | max_deposit_period:(0, 9223372036854775807]              | {"min_deposit": [{"denom": "uiris", "amount": "1000000000"}], "max_deposit_period": "86400s" },                |
| `gov/votingparams`  | 提议投票阶段的相关参数 | voting_period:(0, 9223372036854775807]                   | {"voting_period": "432000s"}                                                                                   |
| `gov/tallyparams`   | 投票统计阶段的相关参数 | quorum:[0,1]<br>threshold:(0,1]<br/>veto_threshold:(0,1] | {"quorum":"0.500000000000000000","threshold": "0.500000000000000000","veto_threshold": "0.330000000000000000"} |

详见 [Governance](../features/governance.md)

## IBC 模块可治理参数

| 字段                      | 描述             | 有效范围     | 当前值                             |
| ------------------------- | ---------------- | ------------ | ---------------------------------- |
| `ibc/AllowedClients`      | 支持ibc的客户端  |              | ["06-solomachine","07-tendermint"] |
| `transfer/SendEnabled`    | 是否支持transfer | {true,false} | false                              |
| `transfer/ReceiveEnabled` | 是否支持Receive  | {true,false} | false                              |

## Mint 模块可治理参数

| 字段             | 描述           | 有效范围 | 当前值 |
| ---------------- | -------------- | -------- | ------ |
| `mint/Inflation` | 代币增发频率   | [0, 0.2] | 0.04   |
| `mint/MintDenom` | 增发的代币名称 |          | uiris  |

详见 [Mint](../features/mint.md)

## Service 模块可治理参数

| 字段                           | 描述                        | 有效范围                  | 当前值                                      |
| ------------------------------ | --------------------------- | ------------------------- | ------------------------------------------- |
| `service/ArbitrationTimeLimit` | 仲裁周期                    | (0, 9223372036854775807]  | 120h0m0s                                    |
| `service/ComplaintRetrospect`  | 投诉周期                    | (0, 9223372036854775807]  | 360h0m0s                                    |
| `service/MaxRequestTimeout`    | 最大请求超时时间            | (0, 9223372036854775807]  | 100                                         |
| `service/MinDepositMultiple`   | 最小抵押倍数                | (0, 9223372036854775807]  | 1000                                        |
| `service/ServiceFeeTax`        | 服务费率                    | [0, 1)                    | 0.05                                        |
| `service/SlashFraction`        | 惩罚系数                    | [0, 1]                    | 0.001                                       |
| `service/TxSizeLimit`          | 交易最大字节数(service模块) | (0, 18446744073709551615] | 4000                                        |
| `service/MinDeposit`           | 最小抵押数量                | amount: (0, +∞)           | [{"denom": "uiris","amount": "5000000000"}] |
| `service/BaseDenom`            | 必须用于抵押的代币          |                           | uiris                                       |

详见 [Service](../features/service.md)

## Slashing 模块可治理参数

| 字段                               | 描述                     | 有效范围                  | 当前值 |
| ---------------------------------- | ------------------------ | ------------------------- | ------ |
| `slashing/DowntimeJailDuration`    | 验证人最大的下线时间     | (0, 9223372036854775807]  | 10m0s  |
| `slashing/MinSignedPerWindow`      | 每个窗口最低投票率       | [0, 1]                    | 0.7    |
| `slashing/SignedBlocksWindow`      | 验证人下线的滑动窗口大小 | (0, 18446744073709551615] | 34560  |
| `slashing/SlashFractionDoubleSign` | 双重签名的惩罚系数       | [0, 1]                    | 0.01   |
| `slashing/SlashFractionDowntime`   | 验证人下线的惩罚系数     | [0, 1]                    | 0.0003 |

详见 [Slashing](../features/slashing.md)

## Staking 模块可治理参数

| 字段                        | 描述                   | 有效范围                 | 当前值   |
| --------------------------- | ---------------------- | ------------------------ | -------- |
| `staking/UnbondingTime`     | 抵押解绑时间           | (0, 9223372036854775807] | 1814400s |
| `staking/MaxValidators`     | 最大验证人数量         | (0, 4294967295]          | 100      |
| `staking/MaxEntries`        | 解绑、转委托的最大数量 | (0, 4294967295]          | 7        |
| `staking/BondDenom`         | 可抵押的代币           |                          | uiris    |
| `staking/HistoricalEntries` | 历史条目               | [0, 4294967295]          | 10000    |

详见 [Staking](../features/staking.md)

## Token 模块可治理参数

| 字段                      | 描述                       | 有效范围        | 当前值                              |
| ------------------------- | -------------------------- | --------------- | ----------------------------------- |
| `token/TokenTaxRate`      | 发行、增发代币的费率       | [0, 1]          | 0.4                                 |
| `token/IssueTokenBaseFee` | 发行代币所需支付的基准费用 | amount: (0, +∞) | {"denom": "iris","amount": "60000"} |
| `token/MintTokenFeeRatio` | 增发代币的费率             | [0, 1]          | 0.1                                 |

详见  [Token](../features/token.md)
