# iriscli service respond 

## 描述

响应服务调用

## 用法

```
iriscli service respond [flags]
```

## 标志

| Name, shorthand       | Default                 | Description                                                                                                                                           | Required |
| --------------------- | ----------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --request-chain-id    |                         | [string] 发起该服务调用的区块链ID                                                                                              | 是       |
| --request-id          |                         | [string] 该服务调用的ID                                                                                                                                | 是       |
| --response-data       |                         | [string] hex编码的服务调用响应数据                                                                       |        |
| -h, --help            |                         | 响应命令帮助                                                                                                                                         |          |

## 示例

### 响应一个服务调用 
```shell
iriscli service respond --chain-id=test-irishub --from=node0 --fee=0.004iris --request-chain-id=test-irishub --request-id=230-130-0 --response-data=abcd
```

运行成功以后，返回的结果如下:

```json
{
   "tags": {
     "action": "service-call",
     "completeConsumedTxFee-iris-atto": "\"162880000000000\"",
     "consumer": "faa1f02ext9duk7h3rx9zm7av0pnlegxve8ne5vw6x",
     "provider": "faa1f02ext9duk7h3rx9zm7av0pnlegxve8ne5vw6x",
     "request-id": "230-130-0"
   }
 }
```

