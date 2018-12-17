# iriscli service withdraw-fees 

## 描述

从服务费收入中取回所有费用

## 用法

```
iriscli service withdraw-fees [flags]
```

## 标志

| Name, shorthand       | Default                 | Description                                                                                                                                           | Required |
| --------------------- | ----------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| -h, --help            |                         | 取回命令帮助                                                                                                                                        |          |

## 示例

### 从服务费收入中取回费用 
```shell
iriscli service withdraw-fees --chain-id=test --from=node0 --fee=0.004iris
```

运行成功以后，返回的结果如下:

```txt
Committed at block 87 (tx hash: C0F61A7200F884277E7C36CAB362E893FCF2445D8E4450D93AEF0755BF346EF6, response:
 {
   "code": 0,
   "data": null,
   "log": "Msg 0: ",
   "info": "",
   "gas_wanted": 200000,
   "gas_used": 3891,
   "codespace": "",
   "tags": {
     "action": "service_withdraw_fees"
   }
 })
```

