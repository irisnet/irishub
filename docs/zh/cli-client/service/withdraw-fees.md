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
iriscli service withdraw-fees --chain-id=test-irishub --from=node0 --fee=0.004iris
```

运行成功以后，返回的结果如下:

```json
{
   "tags": {
     "action": "service-withdraw-fees",
     "completeConsumedTxFee-iris-atto": "\"679600000000000\""
   }
 }
```

