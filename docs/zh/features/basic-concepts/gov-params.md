# IRIShub系统参数

在IRISnet中，存在一些特殊的参数，它们可通过链上治理被修改。持有IRIS通证的用户都可以参与到参数修改的链上治理。
如果社区对某些可修改的参数不满意，完全可以通过发起参数修改提案来完成修改。

##  链上治理可治理参数

* `DepositProcedure`  抵押阶段的参数（最小抵押金额，抵押期）
* `VotingProcedure`   投票阶段的参数（投票期）
* `TallyingProcedure` 统计阶段的参数（投票是否通过的标准）
* 在`DepositProcedure` 抵押阶段, 以下参数可以通过链上治理来修改：
  * 发起提案的最小抵押数，在genesis文件中记录为 `min_deposit` 字段
  * 发起提案的抵押时长，在genesis文件中记录为 `voting_period` 字段
* 在`VotingProcedure` 抵押阶段, 以下参数可以通过链上治理来修改：
   * 对提案的投票时长，在genesis文件中记录为 `voting_period` 
* 在`TallyingProcedure` 抵押阶段, 以下参数可以通过链上治理来修改：
   * 对提案的支持票的最小比例，在genesis文件中记录为`threshold` 
   * 否听提案所要求的vote所占最小比例，在genesis文件中记录为`veto`
   * 对提案的投票的voting power要求的最小比例，在genesis文件中记录为`participation` 
   
详细见[gov](../governance.md)

## Service模块可治理参数

* `MinDepositMultiple`    服务绑定最小抵押金额的倍数
* `MaxRequestTimeout`     服务调用最大等待区块个数
* `ServiceFeeTax`         服务费的税收比例
* `SlashFraction`         惩罚百分比
* `ComplaintRetrospect`   可提起争议最大时长
* `ArbitrationTimeLimit`  争议解决最大时长

详细见[service](../service.md)
