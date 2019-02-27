# iriscli stake unbond

## 介绍

Unbond shares from a validator

## 用法

```
iriscli stake unbond [flags]
```

打印帮助信息

```
iriscli stake unbond --help
```

## 特有标志

| 名称                | 类型   | 是否必填 | 默认值   | 功能描述         |
| --------------------| -----  | -------- | -------- | ------------------------------------------------------------------- |
| --address-validator | string | true     | ""       | 验证人地址 |
| --shares-amount     | float  | false    | 0.0      | 解绑的share数量，正数 |
| --shares-percent    | float  | false    | 0.0      | 解绑的比率，0到1之间的正数 |

用户可以用`--shares-amount`或者`--shares-percent`指定解绑定的token数量，这两个参数不可同时使用。

## 示例

```
iriscli stake unbond --address-validator=<ValidatorAddress> --shares-percent=0.1 --from=<key name> --chain-id=test-irishub --fee=0.3iris
```
输出信息：
```txt
Committed at block 306 (tx hash: 5A4C6E00F4F6BF795EB05D2D388CBA0E8A6E6CF17669314B1EE6A31729A22450, response: {Code:0 Data:[] Log:Msg 0:  Info: GasWanted:200000 GasUsed:3398 Tags:[{Key:[97 99 116 105 111 110] Value:[115 101 114 118 105 99 101 45 119 105 116 104 100 114 97 119 45 102 101 101 115] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[99 111 109 112 108 101 116 101 67 111 110 115 117 109 101 100 84 120 70 101 101 45 105 114 105 115 45 97 116 116 111] Value:[34 54 55 57 54 48 48 48 48 48 48 48 48 48 48 48 34] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0}] Codespace: XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0})
```
```json
 {
   "code": 0,
   "data": "CwiAkrjDmP7///8B",
   "log": "Msg 0: ",
   "info": "",
   "gas_wanted": 200000,
   "gas_used": 18990,
   "codespace": "",
   "tags": {
     "action": "begin_unbonding",
     "delegator": "iaa106nhdckyf996q69v3qdxwe6y7408pvyvufy0x2",
     "end-time": "\u000b\u0008\ufffd\ufffd\ufffdØ\ufffd\ufffd\ufffd\ufffd\u0001",
     "source-validator": "iva1xpqw0kq0ktt3we5gq43vjphh7xcjfy6sfqamll"
   }
 })
```