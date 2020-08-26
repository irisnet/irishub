# iriscli params

Params模块允许查询系统里预设的参数，查询结果中除了Gov模块的参数，其他都可以通过[Gov模块](./gov.md)发起提议来修改。

```bash
iris query params subspace [subspace] [key] [flags]
```

`subspace`目前支持：` auth`，`bank`，` staking`，`mint`，`distribution`，`slashing`，` gov`，` crisis`，` token`，` record`，` htlc`， `coinswap`，`service`。
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
| `UnbondingTime`     | 抵押解绑时间           | 3w        |
| `MaxValidators`     | 最大验证人数量         | 100       |
| `MaxEntries`        | 解绑、转委托的最大数量 | 7         |
| `BondDenom`         | 可抵押的代币           | iris-atto |
| `HistoricalEntries` |                        | 100       |

## mint

| key         | description    | default   |
| ----------- | -------------- | --------- |
| `Inflation` | 代币增发频率   | 0.04      |
| `MintDenom` | 增发的代币名称 | iris-atto |

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

| key             | description            | default |
| --------------- | ---------------------- | ------- |
| `depositparams` | 提议抵押阶段的相关参数 |         |
| `votingparams`  | 提议投票阶段的相关参数 |         |
| `tallyparams`   | 投票统计阶段的相关参数 |         |

## slashing

| key           | description | default |
| ------------- | ----------- | ------- |
| `ConstantFee` | 固定费用    |         |

## token

| key                 | description                | default   |
| ------------------- | -------------------------- | --------- |
| `TokenTaxRate`      | 发行、增发代币的费率       | 0.4       |
| `IssueTokenBaseFee` | 发行代币所需支付的代币数量 | 60000iris |
| `MintTokenFeeRatio` | 增发代币的费率             | 0.1       |

## coinswap

| key             | description          | default   |
| --------------- | -------------------- | --------- |
| `Fee`           | 支付的手续费率       | 0.003     |
| `StandardDenom` | 支付的手续费代币名称 | iris-atto |

##service

| key                    | description                 | default   |
| ---------------------- | --------------------------- | --------- |
| `MaxRequestTimeout`    | 最大请求超时时间            | 100(区块) |
| `MinDepositMultiple`   | 最小抵押倍数                | 200       |
| `MinDeposit`           | 最小抵押数量                | 6000iris  |
| `ServiceFeeTax`        | 服务费率                    | 0.1       |
| `SlashFraction`        | 惩罚系数                    | 0.001     |
| `ComplaintRetrospect`  | 投诉周期                    | 15d       |
| `ArbitrationTimeLimit` | 仲裁周期                    | 5d        |
| `TxSizeLimit`          | 交易最大字节数(service模块) | 4000      |
| `BaseDenom`            | 服务费支持的代币            | iris-atto |


