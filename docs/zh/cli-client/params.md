# iriscli params

Params模块允许查询系统里预设的参数，查询结果中除了Gov模块的参数，其他都可以通过[Gov模块](./gov.md)发起提议来修改。

```bash
 iriscli params <flags>
```

**标志:**

| 名称，速记 | 默认 | 描述     | 必须 |
| ---------- | ---- | -------- | ---- |
| --module   |      | 模块名称 |      |

## 查询所有参数

```bash
iriscli params
```

示例输出:

```bash
Slashing Params:
  slashing/MaxEvidenceAge:           51840
  slashing/SignedBlocksWindow:       34560
  slashing/MinSignedPerWindow:       0.5000000000
  slashing/DoubleSignJailDuration:   48h0m0s
  slashing/DowntimeJailDuration:     24h0m0s
  slashing/CensorshipJailDuration:   48h0m0s
  slashing/SlashFractionDoubleSign:  0.0100000000
  slashing/SlashFractionDowntime:    0.0000000000
  slashing/SlashFractionCensorship:  0.0000000000
Service Params:
  service/MaxRequestTimeout:     100
  service/MinDepositMultiple:    1000
  service/ServiceFeeTax:         0.0100000000
  service/SlashFraction:         0.0010000000
  service/ComplaintRetrospect:   360h0m0s
  service/ArbitrationTimeLimit:  120h0m0s
  service/TxSizeLimit:           4000
Asset Params:
  asset/AssetTaxRate:          0.4000000000
  asset/IssueTokenBaseFee:     300000000000000000000000iris-atto
  asset/MintTokenFeeRatio:     0.1000000000
  asset/CreateGatewayBaseFee:  600000000000000000000000iris-atto
  asset/GatewayAssetFeeRatio:  0.1000000000
Auth Params:
  auth/gasPriceThreshold:  6000000000000
  auth/txSizeLimit:        1000
Stake Params:
  stake/UnbondingTime:  504h0m0s
  stake/MaxValidators:  100
Mint Params:
  mint/Inflation:  0.0400000000
Distribution Params:
  distr/CommunityTax:        0.0200000000
  distr/BaseProposerReward:  0.0100000000
  distr/CommunityTax:        0.0400000000

Gov Params:
System Halt Period:  60
Proposal Parameter:  [Critical]                         [Important]                        [Normal]
  DepositPeriod:     24h0m0s                            24h0m0s                            24h0m0s
  MinDeposit:        4000000000000000000000iris-atto    2000000000000000000000iris-atto    1000000000000000000000iris-atto
  Voting Period:     2m0s                               2m0s                               2m0s
  Max Num:           1                                  5                                  7
  Threshold:         0.7500000000                       0.6700000000                       0.5000000000
  Veto:              0.3300000000                       0.3300000000                       0.3300000000
  Participation:     0.5000000000                       0.5000000000                       0.5000000000
  Penalty:           0.0000000000                       0.0000000000                       0.0000000000
```

## 查询模块参数

可用的模块名称可以通过[查询所有参数](#查询所有参数)查询

```bash
iriscli params --module=slashing
```

示例输出：

```bash
Slashing Params:
  slashing/MaxEvidenceAge:           51840
  slashing/SignedBlocksWindow:       34560
  slashing/MinSignedPerWindow:       0.7000000000
  slashing/DoubleSignJailDuration:   48h0m0s
  slashing/DowntimeJailDuration:     36h0m0s
  slashing/CensorshipJailDuration:   48h0m0s
  slashing/SlashFractionDoubleSign:  0.0100000000
  slashing/SlashFractionDowntime:    0.0003000000
  slashing/SlashFractionCensorship:  0.0000000000
```
