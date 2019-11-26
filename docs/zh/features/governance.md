# 链上治理

## 基本功能描述

1. 纯文本提议的链上治理
2. 参数修改提议的链上治理
3. 软件升级提议的链上治理
4. 网络终止提议的链上治理
5. Tax收入分配提议的链上治理
6. 添加外部资产提议的链上治理

## 交互流程

### 提议级别

不同级别对应的具体 Proposal：

- Critical：`SoftwareUpgrade`, `SystemHalt`
- Important：`Parameter`,`TokenAddition`
- Normal：`CommunityTaxUsage`,`PlainText`

`SoftwareUpgrade Proposal` 和 `SystemHalt Proposal` 只能由profiler发起。

不同级别对应的参数不同：

| 治理参数       | Critical  | Important | Normal    | 范围                    |
| ------------- | --------- | --------- | --------- | ---------------------- |
| MinDeposit    | 4000 iris | 2000 iris | 1000 iris | [10iris,10000iris]     |
| DepositPeriod | 24 hours  | 24 hours  | 24 hours  | [20s,3d]               |
| VotingPeriod  | 120 hours | 120 hours | 120 hours | [20s,7d]               |
| MaxNum        | 1         | 5         | 7         | Critical==1, other(1,) |
| Participation | 0.5       | 0.5       | 0.5       | (0,1)                  |
| Threshold     | 0.75      | 0.67      | 0.5       | (0,1)                  |
| Veto          | 0.33      | 0.33      | 0.33      | (0,1)                  |
| Penalty       | 0         | 0         | 0         | (0,1)                  |

- `MinDeposit`  最小抵押金额
- `DepositPeriod` 抵押阶段的窗口期
- `VotingPeriod` 投票阶段的窗口期
- `MaxNum` 该类型提议在网络中同时能存在的最大个数
- `Penalty`  slash验证人绑定通证的比例
- `Veto`  强烈反对的power占参与投票power的比例， 如果达到这个比例则提议被拒绝， 提议结果为“强烈反对”
- `Threshold`  提议通过所需"赞成"的power占参与投票power的比例， 如果达到这个比例则提议通过， 提议结果为“同意”
- `Participation` 参与投票的power占系统中总voting power的比例， 如果未达到这个比例则提议被拒绝，提议结果为“不同意”

### 抵押阶段

提交提议者至少抵押30%的 `MinDeposit` ，然后其他用户可以继续对该提议进行抵押， 当抵押额超过 `MinDeposit`, 提议才能进入投票阶段。该提议时间超过 `MaxDepositPeriod` ，还未进入投票阶段（总抵押未超过 `MinDeposit`），则提议会被删除，并不会返还抵押金。

不能对进入投票阶段的提议再进行抵押。

### 投票阶段

只有验证人和委托人可以投一次票，不可重复投票。投票选项有：`Yes`同意, `Abstain`弃权,`No`不同意,`NoWithVeto`强烈不同意。

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

```bash
# 根据gov模块名查询的可修改的参数
iriscli params --module=mint

# 结果
Mint Params:
  mint/Inflation=0.0400000000

# 发送提议，返回参数修改的内容
iriscli gov submit-proposal --title=<title> --description=<description> --type=Parameter --deposit=8iris  --param="mint/Inflation=0.0000000000" --from=<key_name> --chain-id=<chain-id> --fee=0.3iris --commit

# 对提议进行抵押
iriscli gov deposit --proposal-id=<proposal-id> --deposit=1000iris --from=<key_name> --chain-id=<chain-id> --fee=0.3iris --commit

# 对提议投票
iriscli gov vote --proposal-id=<proposal-id> --option=Yes --from=<key_name> --chain-id=<chain-id> --fee=0.3iris --commit

# 查询提议情况
iriscli gov query-proposal --proposal-id=<proposal-id>
```

### 社区基金使用提议

有三种使用方式: `Burn`，`Distribute` and `Grant`。 `Burn`表示从社区基金中销毁代币。`Distribute` and `Grant` 将从社区基金中向目标受托人账户转移代币，然后受托人将这些代币分发或赠给其他账户。

```bash
# 提交 Burn 提议
iriscli gov submit-proposal --title="burn tokens 5%" --description=<description> --type="CommunityTaxUsage" --usage="Burn" --deposit="10iris"  --percent=0.05 --from=<key_name> --chain-id=<chain-id> --fee=0.3iris --commit

# 提交 Distribute 提议
iriscli gov submit-proposal --title="distribute tokens 5%" --description="test" --type="CommunityTaxUsage" --usage="Distribute" --deposit="10iris"  --percent=0.05 --dest-address=<dest-address (only trustees)> --from=<key_name> --chain-id=<chain-id> --fee=0.3iris --commit

# 提交 Grant 提议
iriscli gov submit-proposal --title="grant tokens 5%" --description="test" --type="CommunityTaxUsage" --usage="Grant" --deposit="10iris"  --percent=0.05 --dest-address=<dest-address (only trustees)> --from=<key_name> --chain-id=<chain-id> --fee=0.3iris --commit
```

### 系统终止提议

发送该提议可以使系统终止，节点会在提议通过后的SystemHaltHeight（=提议通过高度+systemHaltPeriod）关闭，再启动会进入query-only模式。

```bash
# 发送系统终止提议
iriscli gov submit-proposal --title=<title> --description=<description> --type=SystemHalt --deposit=10iris --fee=0.3iris --from=<key_name> --chain-id=<chain-id> --commit
```

### 软件升级提议

详细参考[Upgrade](upgrade.md)
