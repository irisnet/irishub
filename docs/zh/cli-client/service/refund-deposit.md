# iriscli service refund-deposit 

## 描述

取回所有押金

## 用法

```
iriscli service refund-deposit [flags]
```

## 标志

| Name, shorthand       | Default                 | Description                                                                        | Required |
| --------------------- | ----------------------- | ---------------------------------------------------------------------------------- | -------- |
| --def-chain-id        |                         | [string] 定义该服务的区块链ID                                                         | 是       |
| --service-name        |                         | [string] 服务名称                                                                   | 是       |
| -h, --help            |                         | 取回押金命令帮助                                                                      |          |

## 示例

### 取回所有押金
```shell
iriscli service refund-deposit --chain-id=test  --from=node0 --fee=0.004iris --def-chain-id=test --service-name=test-service
```

运行成功以后，返回的结果如下:

```txt
Committed at block 17 (tx hash: 6C878E864772DE2F29725B743A8B9D40A75B41688F16C278634674653BFD1DFA, response:
 {
   "code": 0,
   "data": null,
   "log": "Msg 0: ",
   "info": "",
   "gas_wanted": 200000,
   "gas_used": 4787,
   "codespace": "",
   "tags": {
     "action": "service_refund_deposit"
   }
 })
```