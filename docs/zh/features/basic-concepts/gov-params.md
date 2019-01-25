# IRIShub系统参数

在IRISnet中，存在一些特殊的参数，它们可通过链上治理被修改。持有IRIS通证的用户都可以参与到参数修改的链上治理。
如果社区对某些可修改的参数不满意，完全可以通过发起参数修改提案来完成修改。

## Auth模块可治理参数

| key |描述 | 有效范围|
|----| ---|---|---|
|`auth/gasPriceThreshold`  |最小的gas单价|(0, 10^18iris-atto]
|`auth/txSizeLimit`  |常交易大小限制 |[500, 1500]

## Stake模块可治理参数

| key |描述 | 有效范围|
|----| ---|---|---|
|`stake/MaxValidators`|  最多验证人数目|[100, 200]
|`stake/UnbondingTime`|  解绑时间|[2week,)

## Distribution模块可治理参数
| key |描述 | 有效范围|
|----| ---|---|---|
|`distr/BaseProposerReward` | 出块奖励的基准比例| (0, 0.02]
|`distr/BonusProposerReward` | 最大额外奖励比例| (0, 0.08]
|`distr/CommunityTax`  | 贡献给社区基金的比例|(0, 0.2]

Details in [distribution](../distribution.md)

## Mint模块可治理参数

| key |描述 | 有效范围|
|----| ---|---|---|
|`mint/Inflation` | 通胀系数 |[0,0.2]

## Slashing模块可治理参数

| key |描述 | 有效范围|
|----| ---|---|
| `slashing/CensorshipJailDuration` | Censorship后Jail的时间 | (0, 4week)
| `slashing/DoubleSignJailDuration`| DoubleSign后Jail的时间| (0, 4week)
| `slashing/DowntimeJailDuration`  | Downtime后Jail的时间 | (0, 4week)
| `slashing/MaxEvidenceAge`|可接受的最早的作恶证据时间 | [1day,)          
| `slashing/MinSignedPerWindow`|slash窗口中最小投票比例 |[0.5, 0.9]      
| `slashing/SignedBlocksWindow`|slash统计窗口区块数 |[100, 140000]      
| `slashing/SlashFractionCensorship`|Censorship后slash的比例 |  [0.005, 0.1]
| `slashing/SlashFractionDoubleSign`|DoubleSign后slash的比例 | [0.01, 0.1]
| `slashing/SlashFractionDowntime`|Downtime后slash的比例 | [0.005, 0.1]   

Details in [slashing](../slashing.md)

## Service模块可治理参数

| key |描述 | 有效范围|
|----| ---|---|
| `service/ArbitrationTimeLimit`|争议解决最大时长 | [5days, 10days]
| `service/ComplaintRetrospect`|可提起争议最大时长| [15days, 30days]
| `service/MaxRequestTimeout`| 服务调用最大等待区块个数|[20,)
| `service/MinDepositMultiple`|服务绑定最小抵押金额的倍数| [500, 5000]
| `service/ServiceFeeTax`|服务费的税收比例| (0, 0.2]
| `service/SlashFraction`|惩罚百分比|  (0, 0.01]
| `service/TxSizeLimit`|service类型交易的大小限制| [2000, 6000]

Details in [service](../service.md)
