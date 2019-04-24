# iriscli gov query-proposals

## 描述

通过可选项过滤查询满足条件的提议

## 使用方式

```
iriscli gov query-proposals <flags>
```

打印帮助信息:

```
iriscli gov query-proposals --help
```

## 标志

| 名称, 速记       | 默认值                      | 描述                                                                                                                                                 | 是否必须  |
| --------------- | -------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --depositor     |                            |按抵押人过滤 |     否     |
| --limit         |                            |限制返回最新[数量]提议。 默认为所有提议  |     否     |
| --status        |                            |按提议状态过滤提议       |     否     |
| --voter         |                            |按投票人过滤      |     否     |

## 例子

### 查询提议

```shell
iriscli gov query-proposals --chain-id=<chain-id>
```

默认查询所有的提议。

```txt
  1 - test proposal
  2 - new proposal
```

当然这里可以查询指定条件的提议。

```shell
gov query-proposals --chain-id=<chain-id> --depositor=iaa14q5rf9sl2dqd2uxrxykafxq3nu3lj2fpascegs
```

可以得到存款人是iaa14q5rf9sl2dqd2uxrxykafxq3nu3lj2fpascegs地址的提议。
```txt
  2 - new proposal
```

查询最新的3条提议
```shell
iriscli gov query-proposals --chain-id=<chain-id> --limit=3
```

