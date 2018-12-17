# iriscli service requests 

## 描述

查询服务请求列表

## 用法

```
iriscli service requests [flags]
```

## 标志

| Name, shorthand       | Default                 | Description                                                                                                                                           | Required |
| --------------------- | ----------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --def-chain-id        |                         | [string] 定义该服务的区块链ID                                                                                              | 是       |
| --service-name        |                         | [string] 服务名称                                                                                                                                 | 是       |
| --bind-chain-id       |                         | [string] 绑定该服务的区块链ID                                                                                                                                 | 是       |
| --provider            |                         | [string] bech32编码的服务提供商账户地址                                                                       | 是       |
| -h, --help            |                         | 查询请求列表命令帮助                                                                                                                                         |          |

## 示例

### Query service request list
```shell
iriscli service requests --def-chain-id=test-irishub --service-name=test-service --bind-chain-id=test-irishub --provider=faa1f02ext9duk7h3rx9zm7av0pnlegxve8ne5vw6x
```

运行成功以后，返回的结果如下:
```txt
Committed at block 306 (tx hash: 5A4C6E00F4F6BF795EB05D2D388CBA0E8A6E6CF17669314B1EE6A31729A22450, response: {Code:0 Data:[] Log:Msg 0:  Info: GasWanted:200000 GasUsed:3398 Tags:[{Key:[97 99 116 105 111 110] Value:[115 101 114 118 105 99 101 45 119 105 116 104 100 114 97 119 45 102 101 101 115] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[99 111 109 112 108 101 116 101 67 111 110 115 117 109 101 100 84 120 70 101 101 45 105 114 105 115 45 97 116 116 111] Value:[34 54 55 57 54 48 48 48 48 48 48 48 48 48 48 48 34] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0}] Codespace: XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0})
```

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

