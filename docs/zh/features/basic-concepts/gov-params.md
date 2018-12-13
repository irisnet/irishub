# Gov Params

在IRISnet中，存在一些特殊的参数，它们可通过链上治理被修改。持有IRIS通证的用户都可以参与到参数修改的链上治理。如果社区对某些可修改的参数不满意，完全可以通过治理模块设置合适的值。

## Gov Module

* `DepositProcedure`  抵押阶段的参数（最小抵押金额，抵押期）
* `VotingProcedure`   投票阶段的参数（投票期）
* `TallyingProcedure` 统计阶段的参数（投票是否通过的标准）

详细见[gov](../governance.md)

## Service Module

* `MaxRequestTimeout`   服务调用最大等待区块个数
* `MinProviderDeposit`  服务绑定最小抵押金额

详细见[service](../service.md)
