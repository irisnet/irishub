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
| --chain-id      |                            | [string] Chain ID of tendermint node                                                                                                                 |          |
| --height        |                            | [int] block height to query, omit to get most recent provable block                                                                                  |          |
| --help, -h      |                            | Help for query-params                                                                                                                                |          |
| --indent        |                            | Add indent to JSON response                                                                                                                          |          |
| --key           |                            | [string] key name of parameter                                                                                                                       |          |
| --ledger        |                            | Use a connected Ledger device                                                                                                                        |          |
| --module        |                            | [string] module name                                                                                                                                 |          |
| --node          | tcp://localhost:26657      | [string] \<host>:\<port> to tendermint rpc interface for this chain                                                                                  |          |
| --trust-node    | true                       | Don't verify proofs for responses                                                                                                                    |          |

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