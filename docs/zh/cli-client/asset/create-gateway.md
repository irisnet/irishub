# iriscli asset create-gateway

## 描述

创建一个网关。网关用于映射外部资产。

## 使用方式

```
iriscli asset create-gateway [flags]
```

打印帮助信息:
```
iriscli asset create-gateway --help
```

## 特有的标志

| 命令, 速记     | 类型   | 是否必须 | 默认值  | 描述                                                         |
| --------------------| -----  | -------- | -------- | ------------------------------------------------------------------- |
| --moniker           | string | 是    | ""       | 唯一的网关名字, 长度为3-8个英文字符 |
| --identity          | string | 否    | ""       | 可选的身份 (例如UPort or Keybase), 最大128个字符 |
| --details           | string | 否    | ""       | 可选的描述, 最大280个字符|
| --website           | string | 否    | ""       | 可选的外部网址, 最大128个字符|


## 示例

```
iriscli asset create-gateway --moniker=tgw --identity=exchange --details=testgateway --website=http://testgateway.io --from=node0 --chain-id=irishub --fee=0.4iris --home iriscli --commit
```

输出信息:
```json
 Committed at block 985 (tx hash: D5C0DA54046DC6C4E88EBC67D85D1C387EF7BCBEF051D351106F1F408177CC79, response:
 {
   "code": 0,
   "data": null,
   "log": "Msg 0: ",
   "info": "",
   "gas_wanted": 50000,
   "gas_used": 10337,
   "codespace": "",
   "tags": [
     {
       "key": "action",
       "value": "create_gateway"
     },
     {
       "key": "moniker",
       "value": "tgw"
     }
   ]
 })
```