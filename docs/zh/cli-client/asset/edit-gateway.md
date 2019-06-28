# iriscli asset edit-gateway

## 描述

编辑指定名字的网关信息

## 使用方式

```
iriscli asset edit-gateway [flags]
```

打印帮助信息:
```
iriscli asset edit-gateway --help
```

## 特有的标志

| 命令, 速记     | 类型   | 是否必须 | 默认值  | 描述                                                         |
| --------------------| -----  | -------- | -------- | ------------------------------------------------------------------- |
| --moniker           | string | 是    | ""       | 唯一的网关名字, 长度为3-8个英文字符|
| --identity          | string | 否    | ""       | 可选的身份 (例如UPort or Keybase), 最大128个字符 |
| --details           | string | 否    | ""       | 可选的描述, 最大280个字符|
| --website           | string | 否    | ""       | 可选的外部网址, 最大128个字符|


## 示例

```
iriscli asset edit-gateway --moniker=tgw --identity=exchange --details=test --website=http://gateway.io --from=node0 --chain-id=irishub --fee=0.4iris --home=iriscli --commit
```

输出信息:
```json
{
   "code": 0,
   "data": null,
   "log": "Msg 0: ",
   "info": "",
   "gas_wanted": 50000,
   "gas_used": 4654,
   "codespace": "",
   "tags": [
     {
       "key": "action",
       "value": "edit_gateway"
     },
     {
       "key": "moniker",
       "value": "tgw"
     }
   ]
}
```