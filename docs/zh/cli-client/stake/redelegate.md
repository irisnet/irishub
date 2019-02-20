      # iriscli stake redelegate

## 介绍

把某个委托的一部分或者全部从一个验证人转移到另外一个验证人

## 用法

```
iriscli stake redelegate [flags]
```

打印帮助信息
```
iriscli stake redelegate --help
```

## 特有的flags

| 名称                       | 类型   | 是否必填 | 默认值   | 功能描述         |
| -------------------------- | -----  | -------- | -------- | ------------------------------------------------------------------- |
| --address-validator-dest   | string | true     | ""       | 目标验证人地址 |
| --address-validator-source | string | true     | ""       | 源验证人地址 |
| --shares-amount            | float  | false    | 0.0      | 转移的share数量，正数 |
| --shares-percent           | float  | false    | 0.0      | 转移的比率，0到1之间的正数 |

用户可以用`--shares-amount`或者`--shares-percent`指定装委托的token数量。记住，这两个参数不可同时使用。

## 示例

执行如下在chain-id为test的网络上重新委托10iris的命令：

```
iriscli stake redelegate --chain-id=test-irishub --from=fuxi --fee=0.12iris --address-validator-source=fva106nhdckyf996q69v3qdxwe6y7408pvyvfcwqmd --address-validator-dest=fva1xpqw0kq0ktt3we5gq43vjphh7xcjfy6sfqamll  --shares-percent=0.1
```
输出信息：
```$xslt
Committed at block 1089 (tx hash: D9A60074B1E15ED33D1C0591AF7B6AD967B13409D342980DC0C858F811021C41, response:
 {
   "code": 0,
   "data": "DAiX2MzgBRCAtK+tAQ==",
   "log": "Msg 0: ",
   "info": "",
   "gas_wanted": 200000,
   "gas_used": 27011,
   "codespace": "",
   "tags": {
     "action": "begin_redelegate",
     "delegator": "faa106nhdckyf996q69v3qdxwe6y7408pvyvufy0x2",
     "destination-validator": "fva1xpqw0kq0ktt3we5gq43vjphh7xcjfy6sfqamll",
     "end-time": "\u000c\u0008\ufffd\ufffd\ufffd\ufffd\u0005\u0010\ufffd\ufffd\ufffd\ufffd\u0001",
     "source-validator": "fva106nhdckyf996q69v3qdxwe6y7408pvyvfcwqmd"
   }
 })
```