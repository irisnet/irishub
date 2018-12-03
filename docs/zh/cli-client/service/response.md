# iriscli service response 

## Description

查询服务响应

## Usage

```
iriscli service response [flags]
```

## Flags

| Name, shorthand       | Default                 | Description                                                                                                                                           | Required |
| --------------------- | ----------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --request-chain-id    |                         | [string] 发起该服务调用的区块链ID                                                                                              |  Yes     |
| --request-id          |                         | [string] 该服务调用的ID                                                                                                                                |  Yes     |
| -h, --help            |                         | 查询响应命令帮助                                                                                                                                         |          |

## Examples

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
    "provider": "faa1f02ext9duk7h3rx9zm7av0pnlegxve8ne5vw6x",
    "consumer": "faa1f02ext9duk7h3rx9zm7av0pnlegxve8ne5vw6x",
    "output": "q80=",
    "error_msg": null
  }
}
```

