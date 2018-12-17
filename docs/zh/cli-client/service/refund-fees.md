# iriscli service refund-fees 

## 描述

从服务费退款中退还所有费用

## 用法

```
iriscli service refund-fees [flags]
```

## 标志

| Name, shorthand       | Default                 | Description                                                                                                                                           | Required |
| --------------------- | ----------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| -h, --help            |                         | 退款命令帮助                                                                                                                                         |          |

## 示例

### 从服务费退款中退还费用 
```shell
iriscli service refund-fees --chain-id=test-irishub --from=node0 --fee=0.004iris
```

运行成功以后，返回的结果如下:

```json
{
   "tags": {
     "action": "service-refund-fees",
     "completeConsumedTxFee-iris-atto": "\"679600000000000\""
   }
 }
```

