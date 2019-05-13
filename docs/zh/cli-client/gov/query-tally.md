# iriscli gov query-tally

## 描述

查询指定提议的投票统计
 
## 使用方式

```
iriscli gov query-tally <flags>
```

打印帮助信息:

```
iriscli gov query-tally --help
```

## 标志
| 名称, 速记       | 默认值                      | 描述                                                                                                                                                 | 是否必须  |
| --------------- | -------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --proposal-id   |                            | 提议ID                                                                                                        | Yes      |

## 例子

### 查询投票统计

查询指定提议的投票统计

```shell
iriscli gov query-tally --chain-id=<chain-id> --proposal-id=<proposal-id>
```

```txt
Tally Result:
  Yes:        0
  Abstain:    100.0000000000
  No:         0
  NoWithVeto: 200.0000000000
```
