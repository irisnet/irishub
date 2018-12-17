# iriscli gov deposit

## 描述
 
充值保证金以激活提议
 
## 使用方式
 
```
iriscli gov deposit [flags]
```

打印帮助信息:

```
iriscli gov deposit --help
```
## 标志
 
| 名称, 速记        | 默认值                      | 描述                                                                                                                                                 | 是否必须  |
| ---------------- | -------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --deposit        |                            | [string] 发起提议的保证金                                                                                                                         | Yes      |
| --proposal-id    |                            | [string] 充值保证金的提议ID                                                                                                        | Yes      |

## 例子

### 充值保证金

```shell
iriscli gov deposit --chain-id=test --proposal-id=1 --deposit=50iris --from=node0 --fee=0.01iris
```

输入正确的密码后，你就充值了50个iris用以激活提议的投票状态。

```txt
Committed at block 7 (tx hash: C1156A7D383492AE5C2EB1BADE0080C3A36BE8AED491DC5B2331056BED5D60DC, response:
 {
   "code": 0,
   "data": null,
   "log": "Msg 0: ",
   "info": "",
   "gas_wanted": 200000,
   "gas_used": 7944,
   "codespace": "",
   "tags": {
     "action": "deposit",
     "depositor": "faa1x25y3ltr4jvp89upymegvfx7n0uduz5kmh5xuz",
     "proposal-id": "1",
     "voting-period-start": "1"
   }
 })
```

如何查询保证金充值明细？

请点击下述链接：

[query-deposit](query-deposit.md)

[query-deposits](query-deposits.md)
