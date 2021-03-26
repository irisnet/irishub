# Params

Params模块允许查询系统里预设的参数，查询结果中除了Gov模块的参数，其他都可以通过[Gov模块](./gov.md)发起提议来修改。

```bash
iris query params subspace [subspace] [key] [flags]
```

`subspace`目前支持：`auth`，`bank`，`staking`，`mint`，`distribution`，`slashing`，`gov`，`crisis`，`token`，`record`，`htlc`， `coinswap`，`service`。
其中，可用于每个`subspace`查询的参数如下：

## auth

| key                      | description                    | default |
| ------------------------ | ------------------------------ | ------- |
| `MaxMemoCharacters`      | 交易字段中备注的最大字符数     | 256     |
| `TxSigLimit`             | 每笔交易的最大签名数量         | 7       |
| `TxSizeCostPerByte`      | 交易的每个字节消耗的Gas        | 10      |
| `SigVerifyCostED25519`   | edd2519算法签名验证消耗的Gas   | 590     |
| `SigVerifyCostSecp256k1` | secp256k1算法签名验证消耗的Gas | 1000    |

## bank

| key                  | description      | default |
| -------------------- | ---------------- | ------- |
| `SendEnabled`        | 支持转账的代币   | {}      |
| `DefaultSendEnabled` | 是否开启转账功能 | true    |

## staking

| key                 | description            | default   |
| ------------------- | ---------------------- | --------- |
| `UnbondingTime`     | 抵押解绑时间           | 3w(weeks) |
| `MaxValidators`     | 最大验证人数量         | 100       |
| `MaxEntries`        | 解绑、转委托的最大数量 | 7         |
| `BondDenom`         | 可抵押的代币           | uiris     |
| `HistoricalEntries` |                        | 100       |

## mint

| key         | description    | default |
| ----------- | -------------- | ------- |
| `Inflation` | 代币增发频率   | 0.04    |
| `MintDenom` | 增发的代币名称 | uiris   |

## distribution

| key                   | description            | default |
| --------------------- | ---------------------- | ------- |
| `communitytax`        | 提现收取的手续费率     | 0.02    |
| `baseproposerreward`  | 区块提议者的基础奖励率 | 0.01    |
| `bonusproposerreward` | 区块提议者的奖励率     | 0.04    |
| `withdrawaddrenabled` | 是否支持设置提现地址   | true    |

## slashing

| key                       | description              | default |
| ------------------------- | ------------------------ | ------- |
| `SignedBlocksWindow`      | 验证人下线的滑动窗口大小 | 100     |
| `MinSignedPerWindow`      | 每个窗口最低投票率       | 0.5     |
| `DowntimeJailDuration`    | 验证人最大的下线时间     | 10m     |
| `SlashFractionDoubleSign` | 双重签名的惩罚系数       | 1/20    |
| `SlashFractionDowntime`   | 验证人下线的惩罚系数     | 1/100   |

## gov

| key             | description            | default                                                      |
| --------------- | ---------------------- | ------------------------------------------------------------ |
| `depositparams` | 提议抵押阶段的相关参数 | `min_deposit`: 10000000uiris; `max_deposit_period`: 2d(days) |
| `votingparams`  | 提议投票阶段的相关参数 | `voting_period`: 2d(days)                                    |
| `tallyparams`   | 投票统计阶段的相关参数 | `quorum`: 0.334; `threshold`: 0.5; `veto_threshold`: 0.334   |

## crisis

| key           | description | default   |
| ------------- | ----------- | --------- |
| `ConstantFee` | 固定费用    | 1000uiris |

## token

| key                 | description                | default            |
| ------------------- | -------------------------- | ------------------ |
| `TokenTaxRate`      | 发行、增发代币的费率       | 0.4                |
| `IssueTokenBaseFee` | 发行代币所需支付的基准费用 | 60000 * 10^6 uiris |
| `MintTokenFeeRatio` | 增发代币的费率             | 0.1                |

## htlc

| key           | description                    | default |
| ------------- | ------------------------------ | ------- |
| `AssetParams` | 支持的资产列表，`[]AssetParam` | None    |

AssetParam参数如下：

| key                        | description                | Example                                      |
| -------------------------- | -------------------------- | -------------------------------------------- |
| `AssetParam.Denom`         | 资产名                     | "htltbcbnb"                                  |
| `AssetParam.SupplyLimit`   | 资产最大供应量             | 100000                                       |
| `AssetParam.Active`        |                            | 是否激活                                     | true |
| `AssetParam.DeputyAddress` | 代理人的IRISHub地址        | “iaa18n3x722r4jpwmshlxnw3ehlpfzywupzefthcz5” |
| `AssetParam.FixedFee`      | 代理人在其他链的固定手续费 | 1000                                         |
| `AssetParam.MinSwapAmount` | 最小交换金额               | 1                                            |
| `AssetParam.MaxSwapAmount` | 最大交换金额               | 10000                                        |
| `AssetParam.MinBlockLock`  | 最小交换到期高度           | 50                                           |
| `AssetParam.MaxBlockLock`  | 最大交换到期高度           | 25480                                        |


## coinswap

| key   | description    | default |
| ----- | -------------- | ------- |
| `Fee` | 支付的手续费率 | 0.003   |

## service

| key                         | description                 | default   |
| --------------------------- | --------------------------- | --------- |
| `MaxRequestTimeout`         | 最大请求超时时间            | 100(区块) |
| `MinDepositMultiple`        | 最小抵押倍数                | 200       |
| `MinDeposit`                | 最小抵押数量                | 6000uiris |
| `ServiceFeeTax`             | 服务费率                    | 0.1       |
| `SlashFraction`             | 惩罚系数                    | 0.001     |
| `ComplaintRetrospect`       | 投诉周期                    | 15d       |
| `ArbitrationTimeLimit`      | 仲裁周期                    | 5d        |
| `TxSizeLimit`               | 交易最大字节数(service模块) | 4000      |
| `BaseDenom`                 | 服务费支持的代币            | uiris     |
| `RestrictedServiceFeeDenom` | 限制服务费token             | false     |
