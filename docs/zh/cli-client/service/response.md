# iriscli service response 

## 描述

查询服务响应

## 用法

```
iriscli service response [flags]
```

## 标志

| Name, shorthand       | Default                 | Description                                                                                                                                           | Required |
| --------------------- | ----------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --request-chain-id    |                         | [string] 发起该服务调用的区块链ID                                                                                              | 是       |
| --request-id          |                         | [string] 该服务调用的ID                                                                                                                                | 是       |
| -h, --help            |                         | 查询响应命令帮助                                                                                                                                         |          |

## 示例

### 查询服务响应
```shell
iriscli service response --request-chain-id=test --request-id=635-535-0
```

运行成功以后，返回的结果如下:

```json
{
  "type": "iris-hub/service/SvcResponse",
  "value": {
    "req_chain_id": "test",
    "request_height": "535",
    "request_intra_tx_counter": 0,
    "expiration_height": "635",
    "provider": "iaa1f02ext9duk7h3rx9zm7av0pnlegxve8npm2k6m",
    "consumer": "iaa1f02ext9duk7h3rx9zm7av0pnlegxve8npm2k6m",
    "output": "q80=",
    "error_msg": null
  }
}
```

