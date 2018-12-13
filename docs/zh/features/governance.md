# Gov User Guide

## 基本功能描述

1. 文本提议的链上治理
2. 参数修改提议的链上治理
3. 软件升级提议的链上治理

## 交互流程

### 治理流程

1. 任何用户可以发起提议，并抵押一部分token，如果超过`min_deposit`,提议进入投票，否则留在抵押期。其他人可以对在抵押期的提议进行抵押token，如果提议的抵押token总和超过`min_deposit`,则进入投票期。但若提议在抵押期停留的出块数目超过`max_deposit_period`，则提议被关闭。
2. 进入投票期的提议，只有验证人和委托人可以进行投票。如果委托人没投票，则他继承他委托的验证人的投票选项。如果委托人投票了，则覆盖他委托的验证人的投票选项。当提议到达`voting_perid`,统计投票结果。
3. 我们统计结果有参与度的限制，其他逻辑细节见[CosmosSDK-Gov-spec](https://github.com/cosmos/cosmos-sdk/blob/v0.26.0/docs/spec/governance/overview.md)

## 使用场景

### 参数修改的使用场景

场景一：通过命令行带入参数修改信息进行参数修改

```
# 根据gov模块名查询的可修改的参数
iriscli gov query-params --module=gov --trust-node

# 结果
[
"Gov/govDepositProcedure",
"Gov/govTallyingProcedure",
"Gov/govVotingProcedure"
]

# 根据Key查询可修改参数的内容
iriscli gov query-params --key=Gov/govDepositProcedure --trust-node

# 结果
{"key":"Gov/govDepositProcedure","value":"{\"min_deposit\":[{\"denom\":\"iris-atto\",\"amount\":\"10000000000000000000\"}],\"max_deposit_period\":172800000000000}","op":""}

# 发送提议，返回参数修改的内容
iriscli gov submit-proposal --title="update MinDeposit" --description="test" --type="ParameterChange" --deposit="10iris"  --param='{"key":"Gov/govDepositProcedure","value":"{\"min_deposit\":[{\"denom\":\"iris-atto\",\"amount\":\"20000000000000000000\"}],\"max_deposit_period\":172800000000000}","op":"update"}}' --from=x --chain-id=gov-test --fee=0.05iris --gas=200000

# 对提议进行抵押
iriscli gov deposit --proposal-id=1 --deposit=1iris --from=x --chain-id=gov-test --fee=0.05iris --gas=200000

# 对提议投票
iriscli gov vote --proposal-id=1 --option=Yes  --from=x --chain-id=gov-test --fee=0.05iris --gas=200000

# 查询提议情况
iriscli gov query-proposal --proposal-id=1 --trust-node

```

场景二，通过文件修改参数

```
# 导出配置文件
iriscli gov pull-params --path=iris --trust-node

# 查询配置文件信息
cat iris/config/params.json                                              {
"gov": {
"Gov/govDepositProcedure": {
"min_deposit": [
{
"denom": "iris-atto",
"amount": "10000000000000000000"
}
],
"max_deposit_period": "172800000000000"
},
"Gov/govVotingProcedure": {
"voting_period": "10000000000"
},
"Gov/govTallyingProcedure": {
"threshold": "0.5000000000",
"veto": "0.3340000000",
"participation": "0.6670000000"
}
}

# 修改配置文件 (TallyingProcedure的governance_penalty)
vi iris/config/params.json                                               {
"gov": {
"Gov/govDepositProcedure": {
"min_deposit": [
{
"denom": "iris-atto",
"amount": "10000000000000000000"
}
],
"max_deposit_period": "172800000000000"
},
"Gov/govVotingProcedure": {
"voting_period": "10000000000"
},
"Gov/govTallyingProcedure": {
"threshold": "0.5000000000",
"veto": "0.3340000000",
"participation": "0.4990000000"
}
}

# 通过文件修改参数的命令，返回参数修改的内容
iriscli gov submit-proposal --title="update MinDeposit" --description="test" --type="ParameterChange" --deposit="10iris"  --path=iris --key=Gov/govTallyingProcedure --op=update --from=x --chain-id=gov-test --fee=0.05iris --gas=20000

# 对提议进行抵押
iriscli gov deposit --proposal-id=1 --deposit=1iris --from=x --chain-id=gov-test --fee=0.05iris --gas=20000

# 对提议投票
iriscli gov vote --proposal-id=1 --option=Yes  --from=x --chain-id=gov-test --fee=0.05iris --gas=20000

# 查询提议情况
iriscli gov query-proposal --proposal-id=1 --trust-node
```

### 软件升级提议部分

详细参考[Upgrade](upgrade.md)

## 基本参数


```
# DepositProcedure（抵押阶段的参数）
"Gov/govDepositProcedure": {
"min_deposit": [
{
"denom": "iris-atto",
"amount": "10000000000000000000"
}
],
"max_deposit_period": "172800000000000"
}
```

* 可修改参数
* 参数的key:"Gov/gov/DepositProcedure"
* `min_deposit[0].denom`  最小抵押token只能是单位是iris-atto的iris通证。
* `min_deposit[0].amount` 最小抵押token数量,默认:10iris,范围（1iris，200iris）
* `max_deposit_period`    补交抵押token的窗口期,默认:172800000000000纳秒==2天,范围（20秒，3天）

```
# VotingProcedure（投票阶段的参数）
"Gov/govVotingProcedure": {
"voting_period": "10000000000"
}
```
* 可修改参数
* `voting_perid` 投票的窗口期,默认:172800000000000纳秒==2天,范围（20秒，3天）

```
# TallyingProcedure (统计阶段段参数)
"Gov/govTallyingProcedure": {
"threshold": "0.5000000000",
"veto": "0.3340000000",
"participation": "0.6670000000"
}
```
* 可修改参数
* `veto` 默认:0.334,范围（0，1）
* `threshold` 默认:0.500,范围（0，1）
* `participation` 默认:0.667,范围（0，1）
*  投票统计逻辑：如果所有投票者的`voting_power`占系统总的`voting_power`的比例没有超过participation，投票不通过。如果强烈反对的`voting_power`占所有投票者的`voting_power` 超过 veto,提议不通过。然后再看赞同的`voting_power`占排除投弃权以外的投票者的总`voting_power` 是否超过 threshold, 超过则提议通过,不超过则不通过。

### 社区基金使用提议
有三种使用方式, `Burn`，`Distribute` and `Grant`。 `Burn`表示从社区基金中销毁代币。`Distribute` and `Grant` 将从社区基金中向目标受托人账户转移代币，然后受托人将这些代币分发或赠给其他账户。
```shell
# Submit Burn usage proposal
iriscli gov submit-proposal --title="burn tokens 5%" --description="test" --type="TxTaxUsage" --usage="Burn" --deposit="10iris"  --percent=0.05 --from=x --chain-id=gov-test --fee=0.05iris --gas=200000

# Submit Distribute usage proposal
iriscli gov submit-proposal --title="distribute tokens 5%" --description="test" --type="TxTaxUsage" --usage="Distribute" --deposit="10iris"  --percent=0.05 --dest-address=[destnation-address] --from=x --chain-id=gov-test --fee=0.05iris --gas=200000

# Submit Grant usage proposal
iriscli gov submit-proposal --title="grant tokens 5%" --description="test" --type="TxTaxUsage" --usage="Grant" --deposit="10iris"  --percent=0.05 --dest-address=[destnation-address] --from=x --chain-id=gov-test --fee=0.05iris --gas=200000

# Deposit for a proposal
iriscli gov deposit --proposal-id=1 --deposit=1iris --from=x --chain-id=gov-test --fee=0.05iris --gas=200000

# Vote for a proposal
iriscli gov vote --proposal-id=1 --option=Yes  --from=x --chain-id=gov-test --fee=0.05iris --gas=200000

# Query the state of a proposal
iriscli gov query-proposal --proposal-id=1 --trust-node
```