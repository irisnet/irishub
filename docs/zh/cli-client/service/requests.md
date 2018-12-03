# iriscli service requests 

## Description

查询服务请求列表

## Usage

```
iriscli service requests [flags]
```

## Flags

| Name, shorthand       | Default                 | Description                                                                                                                                           | Required |
| --------------------- | ----------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --def-chain-id        |                         | [string] 定义该服务的区块链ID                                                                                              |  Yes     |
| --service-name        |                         | [string] 服务名称                                                                                                                                 |  Yes     |
| --bind-chain-id       |                         | [string] 绑定该服务的区块链ID                                                                                                                                 |  Yes     |
| --provider            |                         | [string] bech32编码的服务提供商账户地址                                                                       |  Yes     |
| -h, --help            |                         | 查询请求列表命令帮助                                                                                                                                         |          |

## Examples

### Query service request list
```shell
iriscli service requests --def-chain-id=test --service-name=test-service --bind-chain-id=test --provider=faa1f02ext9duk7h3rx9zm7av0pnlegxve8ne5vw6x
```

运行成功以后，返回的结果如下:

```json
[
  {
    "def_chain_id": "chain-jsmJQQ",
    "def_name": "test-service",
    "bind_chain_id": "chain-jsmJQQ",
    "req_chain_id": "chain-jsmJQQ",
    "method_id": 1,
    "provider": "faa1f02ext9duk7h3rx9zm7av0pnlegxve8ne5vw6x",
    "consumer": "faa1f02ext9duk7h3rx9zm7av0pnlegxve8ne5vw6x",
    "input": "Q0NV",
    "service_fee": [
      {
        "denom": "iris-atto",
        "amount": "10000000000000000"
      }
    ],
    "profiling": false,
    "request_height": "456",
    "request_intra_tx_counter": 0,
    "expiration_height": "556"
  }
]
```

