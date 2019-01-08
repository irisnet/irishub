# Gov User Guide

## 基本功能描述

1. 参数修改提议的链上治理
2. 软件升级提议的链上治理
3. 网络终止提议的链上治理
4. Tax收入分配提议的链上治理

## 交互流程

### 提议级别

不同级别对应的具体 Proposal：
- Critical：`SoftwareUpgrade`, `SystemHalt`
- Important：`ParameterChange`
- Normal：`TxTaxUsage`

`SoftwareUpgradeProposal` 和 `SystemHaltProposal` 只能由profiler发起。`TxTaxUsage`只能由`trustee`发起。

不同级别对应的参数不同：

| Gov参数 | Critical | Important | Normal |Range|
| ------ | ------ | ------ | ------|------| 
| govDepositProcedure/MinDeposit | 4000 iris | 2000 iris | 1000 iris |[10iris,10000iris]|
| govDepositProcedure/MaxDepositPeriod | 24 hours | 24 hours | 24 hours |[20s,3d]|
| govVotingProcedure/VotingPeriod | 72 hours | 60 hours | 48 hours |[20s,3d]|
| govVotingProcedure/MaxProposal | 1 | 2 | 1 |Critial==1, other(1,)|
| govTallyingProcedure/Participation | 6/7 | 5/6 | 3/4 |(0,1)|
| govTallyingProcedure/Threshold | 5/6 | 4/5 | 2/3 |(0,1)|
| govTallyingProcedure/Veto | 1/3 | 1/3 | 1/3 |(0,1)|
| govTallyingProcedure/Penalty | 0.0009 | 0.0007 | 0.0005 |(0,1)|


* `MinDeposit`  最小抵押金额
* `MaxDepositPeriod` 抵押阶段的窗口期
* `VotingPeriod` 投票阶段的窗口期
* `MaxProposal` 该类型提议在网络中同时能存在的最大个数
* `Penalty`  slash验证人绑定通证的比例
* `Veto`  由govTallyingProcedure/Veto定义
* `Threshold` 由govTallyingProcedure/Threshold定义
* `Participation` 由govTallyingProcedure/Participation定义

### 抵押阶段
提交的提议有抵押金，当抵押金超过 `MinDeposit` ,才能进入投票阶段。该提议超过 `MaxDepositPeriod` ，还未进超过 `MinDeposit`，则提议会被删除，并返还全部抵押金。 
不能对进入投票阶段的提议再进行抵押。

### 投票阶段
只有验证人可以投一次票，不可重复投票。投票选项有：`Yes`同意, `Abstain`弃权,`No`不同意,`NoWithVeto`强烈不同意。

### 统计阶段

统计结果有三类：同意，不同意，强烈不同意。

在所有投票者的`voting_power`占系统总的`voting_power`的比例超过participation的前提下,如果强烈反对的`voting_power`占所有投票者的`voting_power` 超过 veto, 结果是强烈不同意。如果没有超过且赞同的`voting_power`占所有投票者的`voting_power` 超过 threshold，提议结果是同意。其他情况皆为不同意。

### 销毁机制

提议通过或未通过，都要销毁Deposit的20%，作为治理的费用，把剩余的 Deposit 按比例原路退回。但如果是强烈不同意，则把Deposit全部销毁。

### 惩罚机制

如果一个账户提议进入投票阶段，他是验证人，然后该提议进入统计阶段，他还是验证人，但是他并没有投票，则会按`Penalty`的比例被惩罚。

## 使用场景

### 参数修改的使用场景

通过命令行带入参数修改信息进行参数修改

```
# 根据gov模块名查询的可修改的参数
iriscli gov query-params --module=mint

# 结果
mint/Inflation=0.0400000000

# 根据Key查询可修改参数的内容
iriscli gov query-params --module=mint --key=mint/Inflation                           

# 结果
mint/Inflation=0.0400000000

# 发送提议，返回参数修改的内容
iriscli gov submit-proposal --title="update MinDeposit" --description="test" --type="ParameterChange" --deposit=8iris  --param mint/Inflation=0.0000000000 --from=x --chain-id=gov-test --fee=0.05iris --gas=200000 --commit

# 对提议进行抵押
iriscli gov deposit --proposal-id=1 --deposit=8iris --from=x --chain-id=gov-test --fee=0.05iris --gas=200000 --commit

# 对提议投票
iriscli gov vote --proposal-id=1 --option=Yes  --from=x --chain-id=gov-test --fee=0.05iris --gas=200000 --commit

# 查询提议情况
iriscli gov query-proposal --proposal-id=1 

```
### 社区基金使用提议
有三种使用方式, `Burn`，`Distribute` and `Grant`。 `Burn`表示从社区基金中销毁代币。`Distribute` and `Grant` 将从社区基金中向目标受托人账户转移代币，然后受托人将这些代币分发或赠给其他账户。
```shell
# Submit Burn usage proposal
iriscli gov submit-proposal --title="burn tokens 5%" --description="test" --type="TxTaxUsage" --usage="Burn" --deposit="10iris"  --percent=0.05 --from=x --chain-id=gov-test --fee=0.05iris --gas=200000 --commit

# Submit Distribute usage proposal
iriscli gov submit-proposal --title="distribute tokens 5%" --description="test" --type="TxTaxUsage" --usage="Distribute" --deposit="10iris"  --percent=0.05 --dest-address=[destnation-address] --from=x --chain-id=gov-test --fee=0.05iris --gas=200000 --commit

# Submit Grant usage proposal
iriscli gov submit-proposal --title="grant tokens 5%" --description="test" --type="TxTaxUsage" --usage="Grant" --deposit="10iris"  --percent=0.05 --dest-address=[destnation-address] --from=x --chain-id=gov-test --fee=0.05iris --gas=200000 --commit

# Deposit for a proposal
iriscli gov deposit --proposal-id=1 --deposit=1iris --from=x --chain-id=gov-test --fee=0.05iris --gas=200000 --commit

# Vote for a proposal
iriscli gov vote --proposal-id=1 --option=Yes  --from=x --chain-id=gov-test --fee=0.05iris --gas=200000 --commit 

# Query the state of a proposal
iriscli gov query-proposal --proposal-id=1
```

### 系统终止提议

发送该提议可以使系统终止，节点会在提议通过后的SystemHaltHeight（=提议通过高度+systemHaltPeriod）关闭，再启动会进入query-only模式。

```
# 发送系统终止提议
iriscli gov submit-proposal  --title=test_title --description=test_description --type=SystemHalt --deposit=10iris --fee=0.005iris --from=x1 --chain-id=gov-test --commit

# Deposit for a proposal
iriscli gov deposit --proposal-id=1 --deposit=1iris --from=x --chain-id=gov-test --fee=0.05iris --gas=200000

# Vote for a proposal
iriscli gov vote --proposal-id=1 --option=Yes  --from=x --chain-id=gov-test --fee=0.05iris --gas=200000

# Query the state of a proposal
iriscli gov query-proposal --proposal-id=1 
```
### 软件升级提议

详细参考[Upgrade](upgrade.md)


