# iriscli stake edit-validator

## 介绍

修改验证的的参数，包括佣金比率，验证人节点名称以及其他描述信息

## 用法

```
iriscli stake edit-validator [flags]
```
打印帮助信息
```
iriscli stake edit-validator --help
```

## 特有标志

| 名称                | 类型   | 是否必填 | 默认值   | 功能描述         |
| --------------------| -----  | -------- | -------- | ------------------------------------------------------------------- |
| --commission-rate   | string | float    | 0.0      | 佣金比率 |
| --moniker           | string | false    | ""       | 验证人名称 |
| --identity          | string | false    | ""       | 身份签名 |
| --website           | string | false    | ""       | 网址  |
| --details           | string | false    | ""       | 验证人节点详细信息 |
| --commit         | String | 否     | True                  |是否等到交易有明确的返回值，如果是True，则忽略--async的内容|


## 示例

```
iriscli stake edit-validator --from=<key name> --chain-id=test-irishub --fee=0.4iris --commission-rate=0.15
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
   "gas_used": 3482,
   "codespace": "",
   "tags": {
     "action": "edit_validator",
     "destination-validator": "fva106nhdckyf996q69v3qdxwe6y7408pvyvfcwqmd",
     "identity": "",
     "moniker": "test2"
   }
 }
```