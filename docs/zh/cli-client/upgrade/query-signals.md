# iriscli upgrade query-signals

## 描述

查询软件升级过程中signal的信息

## 标志

| 名称, 速记       | 默认值                 | 描述                                                                                                                                                 | 是否必须  |
| --------------- | --------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --detail     |   false                    | [bool] signals统计的详细信息                                                                                                                   |       |

## 用法

```
iriscli upgrade query-signals
```

打印帮助信息:

```
iriscli upgrade query-signals --help
```

## 例子

查询软件升级过程中signal的信息

```
iriscli upgrade query-signals
```

```
signalsVotingPower/totalVotingPower = 0.5000000000
```

```
iriscli upgrade query-signals --detail
```

```
iva15cv33a67cfey5eze7238hck6yngw36949evplx   100.0000000000
siganalsVotingPower/totalVotingPower = 0.5000000000
```
