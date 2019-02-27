# iriscli gov query-votes

## 描述

查询指定提议的投票情况

## 使用方式

```
iriscli gov query-votes [flags]
```
打印帮助信息:

```
iriscli gov query-votes --help
```

## 标志

| 名称, 速记       | 默认值                      | 描述                                                                                                                                                 | 是否必须  |
| --------------- | -------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --proposal-id   |                            | [string] 提议ID                                                                                                        | Yes      |

## 例子

### Query votes

```shell
iriscli gov query-votes --chain-id=test --proposal-id=1
```

通过指定的提议查询该提议所有投票者的投票详情。
 
```txt
[
  {
    "voter": "iaa14q5rf9sl2dqd2uxrxykafxq3nu3lj2fpascegs",
    "proposal_id": "1",
    "option": "Yes"
  }
]
```
