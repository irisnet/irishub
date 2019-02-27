# iriscli stake delegate

## 介绍

发送委托交易

## 用法

```
iriscli stake delegate [flags]
```

打印帮助信息
```
iriscli stake delegate --help
```

## 特有的标志

| 名称                | 类型   | 是否必填 | 默认值   | 功能描述         |
| --------------------| -----  | -------- | -------- | ---------------- |
| --address-delegator | string | true     | ""       | 验证人地址 |
| --amount            | string | true     | ""       | 委托的token数量 |
| --commit         | String | 否     | True                  |是否等到交易有明确的返回值，如果是True，则忽略--async的内容|

## 示例

在chain-id为test的链上执行委托10iris的命令：
```
iriscli stake delegate --chain-id=<chain-id> --from=KeyName --fee=0.3iris --amount=10iris --address-validator=iva106nhdckyf996q69v3qdxwe6y7408pvyv3hgcms
```
输出信息：
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
   "gas_used": 16462,
   "codespace": "",
   "tags": {
     "action": "delegate",
     "delegator": "iaa106nhdckyf996q69v3qdxwe6y7408pvyvyxzhxh",
     "destination-validator": "iva106nhdckyf996q69v3qdxwe6y7408pvyv3hgcms"
   }
 }
```