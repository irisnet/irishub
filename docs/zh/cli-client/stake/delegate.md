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

## 特有的flags

| 名称                | 类型   | 是否必填 | 默认值   | 功能描述         |
| --------------------| -----  | -------- | -------- | ---------------- |
| --address-delegator | string | true     | ""       | 验证人地址 |
| --amount            | string | true     | ""       | 委托的token数量 |

## 示例

在chain-id为test的链上执行委托10iris的命令：
```
iriscli stake delegate --chain-id=test-irishub --from=KeyName --fee=0.04iris --amount=10iris --address-validator=fva106nhdckyf996q69v3qdxwe6y7408pvyvfcwqmd
```
输出信息：
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
     "delegator": "faa106nhdckyf996q69v3qdxwe6y7408pvyvufy0x2",
     "destination-validator": "fva106nhdckyf996q69v3qdxwe6y7408pvyvfcwqmd"
   }
 })
```