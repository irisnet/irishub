# IRIShub System Parameters

In IRISnet, there are some special parameters can be modified through on-chain governance. 
All the IRIS holders are able to modify. If the community is not satisfied with certain modifiable parameters, it is available to submit a `parameter-change` proposal in governance module.

## Parameters in Auth

| key |Description | Range|
|----| ---|---|
|`auth/gasPriceThreshold`  |minimum of gas price |(0, 10^18iris-atto]
|`auth/txSizeLimit`  |the limitation of the normal tx size |[500, 1500]

## Parameters in Stake

| key |Description | Range|
|----| ---|---|
|`stake/MaxValidators`|  maximum number of validators|[100, 200]
|`stake/UnbondingTime`|  unbonding time|[2week,)

Details in [distribution](../stake.md)

## Parameters in Distribution

| key |Description | Range|
|----| ---|---|
|`distr/BaseProposerReward` | standard ratio of the block reward| (0, 0.02]
|`distr/BonusProposerReward` | maximum additional bonus ratio| (0, 0.08]
|`distr/CommunityTax`  | proportion of contributions to community funds|(0, 0.2]

Details in [distribution](../distribution.md)

## Parameters in Mint

| key |Description | Range|
|----| ---|---|
|`mint/Inflation` | Inflation coefficient|[0,0.2]

Details in [distribution](../mint.md)

## Parameters in Slashing

| key |Description | Range|
|----| ---|---|
| `slashing/CensorshipJailDuration` | Censorship Jail Duration | (0, 4week)
| `slashing/DoubleSignJailDuration`| DoubleSign Jail Duration | (0, 4week)
| `slashing/DowntimeJailDuration`  | Downtime Jail Duration| (0, 4week)
| `slashing/MaxEvidenceAge`| Acceptable earliest time of the evidence| [1day,)      
| `slashing/MinSignedPerWindow`|Minimum voting ratio in the slash window |[0.5, 0.9]      
| `slashing/SignedBlocksWindow`| The number of blocks for slash statistics|[100, 140000] 
| `slashing/SlashFractionCensorship`| Slash ratio of censorship |  [0.005, 0.1]
| `slashing/SlashFractionDoubleSign`| Slash ratio of double-sign | [0.01, 0.1]
| `slashing/SlashFractionDowntime`  | Slash ratio of downtime     | [0.005, 0.1]   

Details in [slashing](../slashing.md)

## Parameters in Service

| key |Description | Range|
|----| ---|---|
| `service/ArbitrationTimeLimit`|  maximum time of dispute resolution| [5days, 10days]
| `service/ComplaintRetrospect`|    maximum time for submit a dispute| [15days, 30days]
| `service/MaxRequestTimeout`|        maximum number of waiting blocks for service invocation|[20,)
| `service/MinDepositMultiple`|      a multiple of the minimum deposit amount of service binding| [500, 5000]
| `service/ServiceFeeTax`|                tax rate of service fee| (0, 0.2]
| `service/SlashFraction`|                slash fraction|  (0, 0.01]
| `service/TxSizeLimit`|                   the limit of the service tx size| [2000, 6000]

Details in [service](../service.md)
