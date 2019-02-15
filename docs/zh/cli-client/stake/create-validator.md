# iriscli stake create-validator

## 介绍

发送交易申请成为验证人，并在在此验证人上委托一定数额的token

## 用法

```
iriscli stake create-validator [flags]
```

打印帮助信息
```
iriscli stake create-validator --help
```

## 特有的标志

| 名称                         | 类型   | 是否必填 | 默认值   | 功能描述         |
| ---------------------------- | -----  | -------- | -------- | ------------------------------------ |
| --amount                     | string | true     | ""       | 委托token的数量 |
| --commission-max-change-rate | float  | true     | 0.0      | 佣金比率每天变化的最大数量 |
| --commission-max-rate        | float  | true     | 0.0      | 最大佣金比例 |
| --commission-rate            | float  | true     | 0.0      | 初始佣金比例 |
| --details                    | string | false    | ""       | 验证人节点的详细信息 |
| --genesis-format             | bool   | false    | false    | 是否已genesis transaction的方式倒出 |
| --identity                   | string | false    | ""       | 身份信息的签名 |
| --ip                         | string | false    | ""       | 验证人节点的IP |
| --moniker                    | string | true     | ""       | 验证人节点名称 |
| --pubkey                     | string | true     | ""       | Amino编码的验证人公钥 |
| --website                    | string | false    | ""       | 验证人节点的网址 |

## 示例

```
iriscli stake create-validator --chain-id=test-irishub--from=<key name> --fee=0.4iris --pubkey=<Validator PubKey> --commission-max-change-rate=0.01 --commission-max-rate=0.2 --commission-rate=0.1 --amount=100iris --moniker=<validator name>
```

返回信息：
```txt
Committed at block 306 (tx hash: 5A4C6E00F4F6BF795EB05D2D388CBA0E8A6E6CF17669314B1EE6A31729A22450, response: {Code:0 Data:[] Log:Msg 0:  Info: GasWanted:200000 GasUsed:3398 Tags:[{Key:[97 99 116 105 111 110] Value:[115 101 114 118 105 99 101 45 119 105 116 104 100 114 97 119 45 102 101 101 115] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[99 111 109 112 108 101 116 101 67 111 110 115 117 109 101 100 84 120 70 101 101 45 105 114 105 115 45 97 116 116 111] Value:[34 54 55 57 54 48 48 48 48 48 48 48 48 48 48 48 34] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0}] Codespace: XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0})
```

```json
 {
   "code": 0,
   "data": null,
   "log": "Msg 0: ",
   "info": "",
   "gas_wanted": 200000,
   "gas_used": 12050,
   "codespace": "",
   "tags": {
     "action": "create_validator",
     "destination-validator": "fva1xpqw0kq0ktt3we5gq43vjphh7xcjfy6sfqamll",
     "identity": "",
     "moniker": "test"
   }
 }
```
