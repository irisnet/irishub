# iriscli gov query-proposal

## 描述

查询单个提议的详情

## 使用方式

```
iriscli gov query-proposal <flags>
```
打印帮助信息:

```
iriscli gov query-proposal --help
```

## 标志

| 名称, 速记       | 默认值                      | 描述                                                                                                                                                 | 是否必须  |
| --------------- | -------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --proposal-id   |                            | 提议ID                                                                                                        | Yes      |
## 例子

### 查询指定的提议

```shell
iriscli gov query-proposal --chain-id=<chain-id> --proposal-id=<proposal-id>
```

查询指定提议的详情，可以得到结果如下：

```txt
Proposal 94:
  Title:              test proposal
  Type:               TxTaxUsage
  Status:             Rejected
  Submit Time:        2019-05-10 06:37:18.776274942 +0000 UTC
  Deposit End Time:   2019-05-10 06:37:28.776274942 +0000 UTC
  Total Deposit:      1100iris
  Voting Start Time:  2019-05-10 06:37:18.776274942 +0000 UTC
  Voting End Time:    2019-05-10 06:37:28.776274942 +0000 UTC
  Description:        a new text proposal
```
