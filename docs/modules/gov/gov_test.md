# Governance Module Test Flow

#####1：初始化区块链运行环境

```
rm -rf .iris
rm -rf .iriscli
iris init gen-tx --name=iris
iris init --gen-txs --chain-id=gov-test
iris start

```

#####2：提议流程

这里以"参数修改提议"为例，其他提议不需要--params参数。比如，当某一个时间点我们发现，系统产生了很多无用的提议，原因可能是当前最小抵押金额太小了，造成很多人都可以提交一个无意义的提议。
这个时候我们可以使用"参数修改提议"来修改系统预设的最小抵押金额这个参数。我们首先需要知道该参数对应的key值，才能修改。可以使用以下命令查看：
```
iriscli params export gov
```
这个命令会导出所有可以使用"参数修改提议"来修改的参数，比如我们得到如下结果：

```
[
  {
    "key": "gov/depositprocedure/deposit",
    "value": "10000000000000000000iris"
  },
  {
    "key": "gov/depositprocedure/maxDepositPeriod",
    "value": "10"
  },
  {
    "key": "gov/feeToken/gasPriceThreshold",
    "value": "20000000000"
  },
  {
    "key": "gov/tallyingprocedure/penalty",
    "value": "1/100"
  },
  {
    "key": "gov/tallyingprocedure/threshold",
    "value": "1/2"
  },
  {
    "key": "gov/tallyingprocedure/veto",
    "value": "1/3"
  },
  {
    "key": "gov/votingprocedure/votingPeriod",
    "value": "20"
  }
]

```
每一个(key，value)对应一组系统预设的可修改参数，具体的意义会在以后的文档中来完善。这里最小抵押金额的key为gov/depositprocedure/deposit
目前value:10000000000000000000iris（该值是经过进度换算以后的值，具体参考fee-token模块）。我们准备将该值修改为:20000000000000000000iris
提高一倍，来限制提交提议的门槛。可以使用以下命令:

```
iriscli gov submit-proposal --title="update MinDeposit" --description="test" --type="ParameterChange"
                            --deposit="9000000000000000000iris" 
                            --params='[{"key":"gov/depositprocedure/deposit","value":"20000000000000000001iris","op":"update"}]' 
                            --proposer=faa1pkunlumfyglqd9dgup0mwp66kjrp6y09twmuvd 
                            --from=iris 
                            --chain-id=gov-test 
                            --fee=400000000000000iris 
                            --gas=20000

```

这里我抵押的金额为9000000000000000000iris，比最小抵押金额少1000000000000000000iris，所以该提议还未被激活，不能进入投票阶段，需要在10（key:gov/depositprocedure/maxDepositPeriod）个区块之前筹齐10000000000000000000iris的抵押金额(如果当前出块时间为5s,意味着你需要在 5 * 10s之前完成抵押任务)，抵押操作命令如下:

```
iriscli gov deposit --proposalID=1 
                    --depositer=faa1pkunlumfyglqd9dgup0mwp66kjrp6y09twmuvd 
                    --deposit=1000000000000000000iris   
                    --from=iris 
                    --chain-id=gov-test  
                    --fee=200000000000000iris 
                    --gas=20000

```
以上的proposalID是由第一步返回结果得到的，这个阶段我们抵押了1000000000000000000iris个代币，正好等于最小抵押金额，所以该提议将进入投票阶段，提议者可以向各个validator发起投票请求(目前只能靠链下通知，以后可以考虑采用链上通知或者监控通知)。然后各个validator可以先查看该提议内容，使用以下命令:
```
iriscli gov query-proposal --proposalID=1 
```
之后投票者可以根据自己的意愿发起投票，这里我投赞成票(option=Yes)：
```
iriscli gov vote --proposalID=1 
                 --voter=faa1pkunlumfyglqd9dgup0mwp66kjrp6y09twmuvd 
                 --option=Yes  
                 --from=iris 
                 --chain-id=gov-test  
                 --fee=400000000000000iris 
                 --gas=20000
```
注意投票期有最大等待时间20个区块高度(key:gov/votingprocedure/votingPeriod)。如果超过这个高度选票的赞成票比例还未达到1/2(key：gov/tallyingprocedure/threshold)。那么系统认为该提议未被通过。不会退还之前抵押的代币(如果有validator未投票，还会被slash，惩罚的比例为当前stake代币总量的1/100(key：gov/tallyingprocedure/penalty)，当前版本还没有这个机制)。假设当前我们只有一个validator，因为我们投的是赞成票，所以赞成的比例为1>1/2 并且 强烈反对票为0<1/3（key:gov/tallyingprocedure/veto）,所以提议通过。等投票期结束，开始自动执行提议内容：将(key:gov/depositprocedure/deposit,value:10000000000000000000iris)修改为(key:gov/depositprocedure/deposit,value:20000000000000000000iris)。接下来我们可以验证这个结果，查询当前系统的最小抵押金额:
```
iriscli iriscli params export gov/depositprocedure/deposit
```
到此，gov治理流程结束。