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
iriscli service refund-deposit --chain-id=test-irishub  --from=node0 --fee=0.004iris --def-chain-id=test --service-name=test-service
```

运行成功以后，返回的结果如下:

```json
{
   "tags": {
     "action": "service-refund-deposit",
     "completeConsumedTxFee-iris-atto": "\"92280000000000\""
   }
 }
```