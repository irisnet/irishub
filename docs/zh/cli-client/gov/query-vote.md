# iriscli gov query-vote

## 描述

查询指定提议、指定投票者的投票情况

## 使用方式

```
iriscli gov query-vote <flags>
```
打印帮助信息:

```
iriscli gov query-vote --help
```
## 标志

| 名称, 速记       | 默认值                      | 描述                                                                                                                                                 | 是否必须  |
| --------------- | -------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --proposal-id   |                            | 提议ID                                                                                                        | Yes      |
| --voter         |                            | bech32编码的投票人地址                                                                                                                        | Yes      |

## 例子

### 查询投票

通过指定提议、指定投票者查询投票情况。

```shell
iriscli gov query-vote --chain-id=<chain-id> --proposal-id=<proposal-id> --voter=<voter_address>
```

```txt
Voter iaa14q5rf9sl2dqd2uxrxykafxq3nu3lj2fpascegs voted with option Abstain on proposal 99
```

