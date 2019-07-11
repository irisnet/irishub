# iriscli params

## Description

Params allows you to query system parameters, some parameters can be governed by gov module

## Usage:

```
 iriscli params --module=<module name>
```

## Flags

| Name,shorthand | Default | Description   | Required |
| -------------- | ------- | ------------- | -------- |
| -h, --help     |         | Help for params |          |
| --module       |    ""   | name for module|    false  |
## Examples

```
○ → iriscli params
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
System Halt Period:     60
Proposal Parameter:    [Critical]         [Important]        [Normal]
  DepositPeriod:        24h0m0s         24h0m0s        24h0m0s
  MinDeposit:           4000000000000000000000iris-atto         2000000000000000000000iris-atto        1000000000000000000000iris-atto
  Voting Period:        2m0s         2m0s        2m0s
  Max Num:              1         5        7
  Threshold:            0.7500000000         0.6700000000        0.5000000000
  Veto:                 0.3300000000         0.3300000000        0.3300000000
  Participation:        0.5000000000         0.5000000000        0.5000000000
  Penalty:              0.0000000000         0.0000000000        0.0000000000
```

