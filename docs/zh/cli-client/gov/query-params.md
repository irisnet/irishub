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

| 名称, 速记       | 默认值                      | 描述                                                                                                                                                 | 是否必须  |
| --------------- | -------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --key           |                            | [string] 参数的键名称                                                                                                                       |          |
| --module        |                            | [string] 模块名称                                                                                                                                 |          |

## 例子
 
### 通过module查参数

```shell
iriscli gov query-params --module=stake
```

可以检索得到stake模块的所有键值。

```txt
 stake/MaxValidators=100
 stake/UnbondingTime=504h0m0s
```

### 通过key查参数

```shell
iriscli gov query-params --key=stake/MaxValidators
```

可以得到gov模块中指定键值的参数详情。

```txt
 stake/MaxValidators=100
```

注意：--module和--key参数不能同时为空.
