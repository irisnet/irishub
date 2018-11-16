# iriscli gov deposit

## 描述
 
充值保证金以激活提议
 
## 使用方式
 
```
iriscli gov deposit [flags]
```

## 标志
 
| 名称, 速记        | 默认值                      | 描述                                                                                                                                                 | 是否必须  |
| ---------------- | -------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --account-number |                            | [int] AccountNumber number to sign the tx                                                                                                            |          |
| --async          |                            | broadcast transactions asynchronously                                                                                                                |          |
| --chain-id       |                            | [string] Chain ID of tendermint node                                                                                                                 | Yes      |
| --deposit        |                            | [string] deposit of proposal                                                                                                                         | Yes      |
| --dry-run        |                            | ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it                                                              |          |
| --fee            |                            | [string] Fee to pay along with transaction                                                                                                           | Yes      |
| --from           |                            | [string] Name of private key with which to sign                                                                                                      | Yes      |
| --from-addr      |                            | [string] Specify from address in generate-only mode                                                                                                  |          |
| --gas            | 200000                     | [string] gas limit to set per-transaction; set to "simulate" to calculate required gas automatically                                                 |          |
| --gas-adjustment | 1                          | [float] adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored |          |
| --generate-only  |                            | Build an unsigned transaction and write it to STDOUT                                                                                                 |          |
| --help, -h       |                            | help for submit-proposal                                                                                                                             |          |
| --indent         |                            | Add indent to JSON response                                                                                                                          |          |
| --json           |                            | return output in json format                                                                                                                         |          |
| --ledger         |                            | Use a connected Ledger device                                                                                                                        |          |
| --memo           |                            | [string] Memo to send along with transaction                                                                                                         |          |
| --node           | tcp://localhost:26657      | [string] \<host>:\<port> to tendermint rpc interface for this chain                                                                                  |          |
| --print-response |                            | return tx response (only works with async = false)                                                                                                   |          |
| --proposal-id    |                            | [string] proposalID of proposal depositing on                                                                                                        | Yes      |
| --sequence       |                            | [int] Sequence number to sign the tx                                                                                                                 |          |
| --trust-node     | true                       | Don't verify proofs for responses                                                                                                                    |          |

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
     "depositer": "faa14q5rf9sl2dqd2uxrxykafxq3nu3lj2fp9l7pgd",
     "proposal-id": "1"
   }
 }
```

如何查询保证金充值明细？

请点击下述链接：

[query-deposit](query-deposit.md)

[query-deposits](query-deposits.md)