# iriscli gov query-params

## 描述

查询参数型（ParameterChange）提议的配置

## 使用方式

```
iriscli gov query-params <flags>
```

打印帮助信息:

```
iriscli gov query-params --help
```

## 标志

| 名称, 速记       | 默认值                      | 描述                                                                                                                                                 | 是否必须  |
| --------------- | -------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --module        |                            | 模块名称                                                                                                                                 |          |

## 例子
 
### 通过module查参数

可以检索得到stake模块的所有键值。

```shell
iriscli gov query-params --module=stake
```

```txt
 stake/MaxValidators=100
 stake/UnbondingTime=504h0m0s
```

### 所有可查询参数的模块

```shell
iriscli gov query-params --module=auth
iriscli gov query-params --module=mint
iriscli gov query-params --module=stake
iriscli gov query-params --module=slashing
iriscli gov query-params --module=distr
iriscli gov query-params --module=service
```
