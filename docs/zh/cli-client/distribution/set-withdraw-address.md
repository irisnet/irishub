# iriscli distribution set-withdraw-addr

## 介绍

设置取回收益时的收款地址

## 用法

```
iriscli distribution set-withdraw-addr [withdraw-address] [flags]
```

打印帮助信息:

```
iriscli distribution set-withdraw-addr --help
```

## 示例

```
iriscli distribution set-withdraw-addr [withdraw-address] --from <key name> --fee=0.004iris --chain-id=<chain-id>
```
执行结果示例

```json
{
   "code": 0,
   "data": null,
   "log": "Msg 0: ",
   "info": "",
   "gas_wanted": 200000,
   "gas_used": 3876,
   "codespace": "",
   "tags": [
     {
       "key": "action",
       "value": "set_withdraw_address"
     },
     {
       "key": "delegator",
       "value": "faa1yclscskdtqu9rgufgws293wxp3njsesxtplqxd"
     }
   ]
 })
```