# iriscli gov query-params

## 描述

查询参数型（ParameterChange）提议的配置

## 使用方式

```
iriscli gov query-params [flags]
```
打印帮助信息:

```
iriscli gov query-params --help
```
## 标志

| 名称, 缩写       | 默认值                      | 描述                                                                                                                                                 | 是否必须  |
| --------------- | -------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --key           |                            | [string] 参数的键名称                                                                                                                       |          |
| --module        |                            | [string] 模块名称                                                                                                                                 |          |

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
