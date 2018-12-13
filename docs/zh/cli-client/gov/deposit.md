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
 
| 名称, 缩写        | 默认值                      | 描述                                                                                                                                                 | 是否必须  |
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
Password to sign with 'node0':
Committed at block 473 (tx hash: 0309E969589F19A9D9E4BFC9479720487FBF929ED6A88824414C5E7E91709206, response: {Code:0 Data:[] Log:Msg 0:  Info: GasWanted:200000 GasUsed:6710 Tags:[{Key:[97 99 116 105 111 110] Value:[100 101 112 111 115 105 116] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[100 101 112 111 115 105 116 101 114] Value:[102 97 97 49 52 113 53 114 102 57 115 108 50 100 113 100 50 117 120 114 120 121 107 97 102 120 113 51 110 117 51 108 106 50 102 112 57 108 55 112 103 100] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[112 114 111 112 111 115 97 108 45 105 100] Value:[49] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[99 111 109 112 108 101 116 101 67 111 110 115 117 109 101 100 84 120 70 101 101 45 105 114 105 115 45 97 116 116 111] Value:[34 51 51 53 53 48 48 48 48 48 48 48 48 48 48 48 34] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0}] Codespace: XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0})
{
   "tags": {
     "action": "deposit",
     "completeConsumedTxFee-iris-atto": "\"335500000000000\"",
     "depositor": "faa14q5rf9sl2dqd2uxrxykafxq3nu3lj2fp9l7pgd",
     "proposal-id": "1"
   }
 }
```

如何查询保证金充值明细？

请点击下述链接：

[query-deposit](query-deposit.md)

[query-deposits](query-deposits.md)
