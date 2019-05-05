# iriscli gov deposit

## 描述
 
抵押保证金以激活提议
 
## 使用方式
 
```
iriscli gov deposit <flags>
```

打印帮助信息:

```
iriscli gov deposit --help
```


## 特殊标志
 
| 名称, 速记        | 默认值                      | 描述                                                                                                                                                 | 是否必须  |
| ---------------- | -------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --deposit        |                            | 发起提议的保证金                                                                                                                         | Yes      |
| --proposal-id    |                            | 抵押保证金的提议ID                                                                                                        | Yes      |

## 例子

### 抵押保证金

```shell
iriscli gov deposit --chain-id=<chain-id> --proposal-id=<proposal-id> --deposit=50iris --from=<key_name> --fee=0.3iris
```

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
     "depositor": "iaa1x25y3ltr4jvp89upymegvfx7n0uduz5krcj7ul",
     "proposal-id": "1",
     "voting-period-start": "1"
   }
 })
```

当所有的抵押金额超过该提议类型的最小抵押额`MinDeposit`，提议将进入投票阶段

| GovParams | Critical | Important | Normal |
| ------ | ------ | ------ | ------|
| MinDeposit | 4000 iris | 2000 iris | 1000 iris |


### 如何查询抵押保证金

[query-deposit](query-deposit.md)

[query-deposits](query-deposits.md)
