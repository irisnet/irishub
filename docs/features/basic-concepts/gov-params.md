# System Parameters

In IRISnet, there are some special parameters can be modified through on-chain governance. 
All the IRIS holders are able to modify. If the community is not satisfied with certain modifiable 
parameters, it is available to put up a `parameter-change` proposal in governance module.

## Parameters in Auth

* `auth/gasPriceThreshold`  minimum of gas price
* `auth/txSizeLimit`  the limit of the normal txsize

## Parameters in Stake

* `stake/MaxValidators`  maximum number of validators
* `stake/UnbondingTime`  unbonding time

## Parameters in Distribution

* `distr/BaseProposerReward`  benchmark ratio of the block reward
* `distr/BonusProposerReward`  maximum additional bonus ratio
* `distr/CommunityTax`  proportion of contributions to community funds

Details in [distribution](../distribution.md)

## Parameters in Mint

* `mint/Inflation`  Inflation coefficient

## Parameters in Slashing

* `slashing/CensorshipJailDuration` 
* `slashing/DoubleSignJailDuration`
* `slashing/DowntimeJailDuration`  
* `slashing/MaxEvidenceAge`         
* `slashing/MinSignedPerWindow`     
* `slashing/SignedBlocksWindow`      
* `slashing/SlashFractionCensorship` 
* `slashing/SlashFractionDoubleSign` 
* `slashing/SlashFractionDowntime`   

Details in [slashing](../slashing.md)

## Parameters in Service

* `service/ArbitrationTimeLimit` maximum time of dispute resolution
* `service/ComplaintRetrospect`   maximum time for submit a dispute
* `service/MaxRequestTimeout`       maximum number of waiting blocks for service invocation
* `service/MinDepositMultiple`     a multiple of the minimum deposit amount of service binding
* `service/ServiceFeeTax`               tax rate of service fee
* `service/SlashFraction`               slash fraction
* `service/TxSizeLimit`                  the limit of the service txsize

Details in [service](../service.md)

