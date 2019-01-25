# IRIShub系统参数

在IRISnet中，存在一些特殊的参数，它们可通过链上治理被修改。持有IRIS通证的用户都可以参与到参数修改的链上治理。
如果社区对某些可修改的参数不满意，完全可以通过发起参数修改提案来完成修改。

## Auth模块可治理参数

* `auth/gasPriceThreshold` 最小的gas单价
* `auth/txSizeLimit` 正常交易大小限制

## Stake模块可治理参数

* `stake/MaxValidators` 最多验证人数目
* `stake/UnbondingTime` 解绑时间

## Distribution模块可治理参数

* `distr/BaseProposerReward` 出块奖励的基准比例
* `distr/BonusProposerReward` 最大额外奖励比例
* `distr/CommunityTax` 贡献给社区基金的比例

详细见[distribution](../distribution.md)

## Mint模块可治理参数

* `mint/Inflation` 通胀系数

## Slashing模块可治理参数

* `slashing/CensorshipJailDuration` 
* `slashing/DoubleSignJailDuration`
* `slashing/DowntimeJailDuration`  
* `slashing/MaxEvidenceAge`         
* `slashing/MinSignedPerWindow`     
* `slashing/SignedBlocksWindow`      
* `slashing/SlashFractionCensorship` 
* `slashing/SlashFractionDoubleSign` 
* `slashing/SlashFractionDowntime`   

详细见[slashing](../slashing.md)

## Service模块可治理参数

* `service/ArbitrationTimeLimit` 争议解决最大时长
* `service/ComplaintRetrospect` 可提起争议最大时长
* `service/MaxRequestTimeout` 服务调用最大等待区块个数
* `service/MinDepositMultiple` 服务绑定最小抵押金额的倍数
* `service/ServiceFeeTax`  服务费的税收比例
* `service/SlashFraction` 惩罚百分比
* `service/TxSizeLimit` service类型交易的大小限制

详细见[service](../service.md)
