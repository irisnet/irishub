# iriscli gov pull-params

## 描述

生成param.json文件

## 使用方式

```
iriscli gov pull-params [flags]
```

## 标志

| 名称, 速记       | 默认值                      | 描述                                                                                                                                                 | 是否必须  |
| --------------- | -------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --chain-id      |                            | [string] tendermint节点的链ID                                                                                                                 |          |
| --height        |                            | [int] 查询的区块高度                                                                                  |          |
| --help, -h      |                            | 查询命令帮助                                                                                                                                 |          |
| --indent        |                            | 在JSON响应中添加缩进                                                                                                                          |          |
| --ledger        |                            | 使用连接的硬件记账设备                                                                                                                        |          |
| --node          | tcp://localhost:26657      | [string] tendermint节点开启的远程过程调用接口\<主机>:\<端口>                                                                                  |          |
| --path          | $HOME/.iris                | [string] iris home目录                                                                                                                      |          |
| --trust-node    | true                       | 关闭响应结果校验                                                                                                                    |          |

## 例子

### 生成参数文件

```shell
iriscli gov pull-params
```

执行该指令，你会收到如下提示信息：

```txt
Save the parameter config file in  /Users/trevorfu/.iris/config/params.json
```

打开--path/config目录下的params.json文件，你可以看到json格式的文件内容。

```txt
{
  "gov": {
    "Gov/govDepositProcedure": {
      "min_deposit": [
        {
          "denom": "iris-atto",
          "amount": "1000000000000000000000"
        }
      ],
      "max_deposit_period": "172800000000000"
    },
    "Gov/govVotingProcedure": {
:  "gov": {
    "Gov/govDepositProcedure": {
      "min_deposit": [
        {
          "denom": "iris-atto",
          "amount": "1000000000000000000000"
        }
      ],
      "max_deposit_period": "172800000000000"
    },
    "Gov/govVotingProcedure": {
      "voting_period": "172800000000000"
:  "gov": {
    "Gov/govDepositProcedure": {
      "min_deposit": [
        {
          "denom": "iris-atto",
          "amount": "1000000000000000000000"
        }
      ],
      "max_deposit_period": "172800000000000"
    },
    "Gov/govVotingProcedure": {
      "voting_period": "172800000000000"
    },
    "Gov/govTallyingProcedure": {
      "threshold": "0.5000000000",
      "veto": "0.3340000000",
      "participation": "0.6670000000"
    }
  }
}
```