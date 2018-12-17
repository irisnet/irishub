# iriscli gov vote

## 描述

给VotingPeriod状态的提议投票, 可选项包括: Yes/No/NoWithVeto/Abstain

## 使用方式

```
iriscli gov vote [flags]
```

打印帮助信息:

```
iriscli gov vote --help
```
## 标志

| 名称, 速记        | 默认值                      | 描述                                                                                                                                                 | 是否必须 |
| ---------------- | -------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --option         |                            | [string] 投票选项 {Yes, No, NoWithVeto, Abstain}                                                                                                  | Yes      |
| --proposal-id    |                            | [string] 投票的提议ID                                                                                                            | Yes      |

## 例子

### 给提议投票

```shell
iriscli gov vote --chain-id=test --proposal-id=1 --option=Yes --from node0 --fee=0.01iris
```

输入正确的密码之后，你就完成了对于所指定的提议投票。
注意：验证人和委托人才能对指定提议投票，并且可投票的提议必须是'VotingPeriod'状态。

```txt
Committed at block 43 (tx hash: 01C4C3B00C6048A12AE2CF2294F63C55A69011381B819C35F11B04C921DB81CC, response:
 {
   "code": 0,
   "data": null,
   "log": "Msg 0: ",
   "info": "",
   "gas_wanted": 200000,
   "gas_used": 2048,
   "codespace": "",
   "tags": {
     "action": "vote",
     "proposal-id": "2",
     "voter": "faa1x25y3ltr4jvp89upymegvfx7n0uduz5kmh5xuz"
   }
 })
```

如何查询投票详情？

请点击下述链接：

[query-vote](query-vote.md)

[query-votes](query-votes.md)

[query-tally](query-tally.md)
