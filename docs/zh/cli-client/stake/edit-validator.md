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

## 特有flags

| 名称                | 类型   | 是否必填 | 默认值   | 功能描述         |
| --------------------| -----  | -------- | -------- | ------------------------------------------------------------------- |
| --commission-rate   | string | float    | 0.0      | 佣金比率 |
| --moniker           | string | false    | ""       | 验证人名称 |
| --identity          | string | false    | ""       | 身份签名 |
| --website           | string | false    | ""       | 网址  |
| --details           | string | false    | ""       | 验证人节点详细信息 |


## 示例

```
iriscli stake edit-validator --from=<key name> --chain-id=test-irishub --fee=0.004iris --commission-rate=0.15
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
 })
```