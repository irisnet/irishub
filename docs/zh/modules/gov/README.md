# Gov/Iparam User Guide

## 基本功能描述

1. 文本提议的链上治理
2. 参数修改提议的链上治理
3. 软件升级提议的链上治理（不可用）

## 交互流程

### 治理流程

1. 任何用户可以发起提议，并抵押一部分资金，如果超过`min_deposit`,提议进入投票，否则留在抵押期。其他人可以对在抵押期的提议进行抵押资金，如果提议的抵押资金总和超过`min_deposit`,则进入投票期。但是提议在抵押期停留的区块数目超过`max_deposit_period`，则提议被关闭。
2. 进入投票期的提议，只有验证人和委托人可以进行投票，委托人如果没投票，则他继承他委托的验证人的投票选项，如果委托人投票了，则覆盖他委托的验证人的投票选项，当提议到达`voting_perid`,统计投票结果。
3. 具体提议投票逻辑细节见[CosmosSDK-Gov-spec](https://github.com/cosmos/cosmos-sdk/blob/develop/docs/spec/governance/overview.md)

## 使用场景
### 创建使用环境

```
rm -rf iris                                                                         
rm -rf .iriscli
iris init gen-tx --name=x --home=iris
iris init --gen-txs --chain-id=gov-test -o --home=iris
iris start --home=iris
```

### 参数修改的使用场景

场景一：通过命令行带入参数修改信息进行参数修改

```
# 根据gov模块名查询的可修改的参数
iriscli gov query-params --module=gov --trust-node

# 结果
[
 "Gov/gov/DepositProcedure",
 "Gov/gov/TallyingProcedure",
 "Gov/gov/VotingProcedure"
]

# 根据Key查询可修改参数的内容
iriscli gov query-params --key=Gov/gov/DepositProcedure --trust-node

# 结果
{"key":"Gov/gov/DepositProcedure","value":"{\"min_deposit\":[{\"denom\":\"iris-atto\",\"amount\":\"10000000000000000000\"}],\"max_deposit_period\":10}","op":""}

# 发送提议，返回参数修改的内容
echo 1234567890 | iriscli gov submit-proposal --title="update MinDeposit" --description="test" --type="ParameterChange" --deposit="10iris"  --param='{"key":"Gov/gov/DepositProcedure","value":"{\"min_deposit\":[{\"denom\":\"iris-atto\",\"amount\":\"10000000000000000000\"}],\"max_deposit_period\":20}","op":"update"}' --from=x --chain-id=gov-test --fee=0.05iris --gas=20000

# 对提议进行抵押
echo 1234567890 | iriscli gov deposit --proposal-id=1 --deposit=1iris --from=x --chain-id=gov-test --fee=0.05iris --gas=20000

# 对提议投票
echo 1234567890 | iriscli gov vote --proposal-id=1 --option=Yes  --from=x --chain-id=gov-test --fee=0.05iris --gas=20000

# 查询提议情况
iriscli gov query-proposal --proposal-id=1 --trust-node

```

场景二，通过文件修改参数

```
# 导出配置文件
iriscli gov pull-params --path=iris --trust-node

# 查询配置文件信息
cat iris/config/params.json                                                         
{
  "gov": {
    "Gov/gov/DepositProcedure": {
      "min_deposit": [
        {
          "denom": "iris-atto",
          "amount": "10000000000000000000"
        }
      ],
      "max_deposit_period": "10"
    },
    "Gov/gov/VotingProcedure": {
      "voting_period": "10"
    },
    "Gov/gov/TallyingProcedure": {
      "threshold": "1/2",
      "veto": "1/3",
      "governance_penalty": "1/100"
    }
  }
}
# 修改配置文件(TallyingProcedure的governance_penalty)
vi iris/config/params.json                                                            
{
  "gov": {
    "Gov/gov/DepositProcedure": {
      "min_deposit": [
        {
          "denom": "iris-atto",
          "amount": "10000000000000000000"
        }
      ],
      "max_deposit_period": "10"
    },
    "Gov/gov/VotingProcedure": {
      "voting_period": "10"
    },
    "Gov/gov/TallyingProcedure": {
      "threshold": "1/2",
      "veto": "1/3",
      "governance_penalty": "20/100"
    }
  }
}

# 通过文件修改参数的命令，返回参数修改的内容
echo 1234567890 | iriscli gov submit-proposal --title="update MinDeposit" --description="test" --type="ParameterChange" --deposit="10iris"  --path=iris --key=Gov/gov/TallyingProcedure --op=update --from=x --chain-id=gov-test --fee=0.05iris --gas=20000

# 对提议进行抵押
echo 1234567890 | iriscli gov deposit --proposal-id=1 --deposit=1iris --from=x --chain-id=gov-test --fee=0.05iris --gas=20000

# 对提议投票
echo 1234567890 | iriscli gov vote --proposal-id=1 --option=Yes  --from=x --chain-id=gov-test --fee=0.05iris --gas=20000

# 查询提议情况
iriscli gov query-proposal --proposal-id=1 --trust-node
```

## CLI命令详情

### 治理模块基础方法

```
# 文本类提议
iriscli gov submit-proposal --title="update MinDeposit" --description="test" --type="Text" --deposit="10iris" --from=x --chain-id=gov-test --fee=0.05iris --gas=20000
```

* `--title`       提议的标题
* `--description` 提议的描述
* `--type`        提议的类型 {'Text','ParameterChange','SoftwareUpgrade'}
* `--deposit`     抵押贷币的数量
* 上面就是典型的文本类提议

```
iriscli gov deposit --proposal-id=1 --deposit=1iris --from=x --chain-id=gov-test --fee=0.05iris --gas=20000
```

* `--propsal-id` 抵押提议ID
* `--deposit`    抵押的贷币数目

```
iriscli gov vote --proposal-id=1 --option=Yes  --from=x --chain-id=gov-test --fee=0.05iris --gas=20000
```

* `--proposal-id` 投票提议ID
* `--option`      投票的选项{'Yes'赞同,'Abstain'弃权,'No'不同意,'nowithVeto'强烈不同意}


```
# 查询提议情况
iriscli gov query-proposal --proposal-id=1 --trust-node
```

* `--proposal-id` 查询提议ID



### 参数修改提议部分

```
# 根据gov模块名查询的可修改的参数
iriscli gov query-params --module=gov --trust-node
```

* `--module` 查询module可修改参数的key的列表


```
# 根据Key查询可修改参数的内容
iriscli gov query-params --key=Gov/gov/DepositProcedure --trust-node
```

* `--key` 查询key对应的参数值

```
# 导出配置文件
iriscli gov pull-params --path=iris --trust-node
```

* `--path` 节点初始化的文件夹



```
# 通过命令行带入参数修改信息进行参数修改 
iriscli gov submit-proposal --title="update MinDeposit" --description="test" --type="ParameterChange" --deposit="10iris"  --param='{"key":"Gov/gov/DepositProcedure","value":"{\"min_deposit\":[{\"denom\":\"iris-atto\",\"amount\":\"10000000000000000000\"}],\"max_deposit_period\":20}","op":"update"}' --from=x --chain-id=gov-test --fee=0.05iris --gas=20000
```

* `--param` 参数修改的具体内容（通过query-params得到参数内容，然后直接对其修改，并在"op"上添上update,具体可见使用场景）
* 其他字段与文本提议类似

```
# 通过文件修改参数的命令，返回参数修改的内容
echo 1234567890 | iriscli gov submit-proposal --title="update MinDeposit" --description="test" --type="ParameterChange" --deposit="10iris"  --path=iris --key=Gov/gov/TallyingProcedure --op=update --from=x --chain-id=gov-test --fee=0.05iris --gas=20000
```

* `--path` 节点初始化的文件夹
* `--key`  要修改参数的key
* `--op`   参数修改类型，目前只实现了'update'
* 其他字段与文本提议类似

### 软件升级提议部分

## 基本参数

```
# DepositProcedure（抵押阶段的参数）
"Gov/gov/DepositProcedure": {
    "min_deposit": [
    {
        "denom": "iris-atto",
        "amount": "10000000000000000000"
    }
    ],
    "max_deposit_period": "10"
}
```

* 可修改参数
* 参数的key:"Gov/gov/DepositProcedure"
* `min_deposit[0].denom`  最小抵押贷币的token只能是单位是iris-atto的iris通证。
* `min_deposit[0].amount` 最小抵押贷币的数量,默认范围:10iris（1iris，200iris）
* `max_deposit_period`    补交抵押的窗口期,默认:10 范围（0，1）     

```
# VotingProcedure（投票阶段的参数）
"Gov/gov/VotingProcedure": {
    "voting_period": "10"
},
```
    
* `voting_perid` 投票的窗口期,默认10,范围（20，20000）
   
```
# TallyingProcedure (统计阶段段参数)    
"Gov/gov/TallyingProcedure": {
    "threshold": "1/2",
    "veto": "1/3",
    "governance_penalty": "1/100"
}
```   
* `veto` 默认1/3,范围（0，1）
* `threshold` 默认1/2,范围（0，1）
* `governance_penalty` 未投票的验证人惩罚贷币的比例 默认1/100,范围（0，1）
*  投票统计逻辑：如果强烈反对的voting_power占总的voting_power 超过 veto,提议不通过。然后再看赞同的voting_power占总的投票的voting_power 是否超过 veto,超过则提议不通过,不超过则不通过。


