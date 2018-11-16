# iriscli gov query-params

## 描述

查询参数型（ParameterChange）提议的配置

## 使用方式

```
iriscli gov query-params [flags]
```

## 标志

| 名称, 速记       | 默认值                      | 描述                                                                                                                                                 | 是否必须  |
| --------------- | -------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --chain-id      |                            | [string] tendermint节点的链ID                                                                                                                 |          |
| --height        |                            | [int] 查询的区块高度                                                                                  |          |
| --help, -h      |                            | 查询命令帮助                                                                                                                                |          |
| --indent        |                            | 在JSON响应中添加缩进                                                                                                                          |          |
| --key           |                            | [string] 参数的键名称                                                                                                                       |          |
| --ledger        |                            | 使用连接的硬件记账设备                                                                                                                        |          |
| --module        |                            | [string] 模块名称                                                                                                                                 |          |
| --node          | tcp://localhost:26657      | [string] tendermint节点开启的远程过程调用接口\<主机>:\<端口>                                                                                  |          |
| --trust-node    | true                       | 关闭响应结果校验                                                                                                                    |          |

## 例子
 
### 通过module查参数

```shell
iriscli gov query-params --module=gov
```

可以检索得到gov模块的所有键值。

```txt
[
 "Gov/govDepositProcedure",
 "Gov/govTallyingProcedure",
 "Gov/govVotingProcedure"
]
```

### 通过key查参数

```shell
iriscli gov query-params --key=Gov/govDepositProcedure
```

可以得到gov模块中指定键值的参数详情。

```txt
{"key":"Gov/govDepositProcedure","value":"{\"min_deposit\":[{\"denom\":\"iris-atto\",\"amount\":\"1000000000000000000000\"}],\"max_deposit_period\":172800000000000}","op":""}
```

注意：--module和--key参数不能同时为空.