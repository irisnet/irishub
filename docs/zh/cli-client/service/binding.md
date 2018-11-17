# iriscli service binding

## 描述

查询服务绑定

## 用法

```
iriscli service binding [flags]
```

## 标志

| Name, shorthand | Default                    | Description                                            | Required |
| --------------- | -------------------------- | ----------------------------------------------------   | -------- |
| --bind-chain-id |                            | [string] 绑定该服务的区块链ID                             | 是        |
| --def-chain-id  |                            | [string] 定义该服务的区块链ID                             | 是        |
| --provider      |                            | [string] 服务提供者的区块链地址(bech32编码)                 | 是        |
| --service-name  |                            | [string] 服务名称                                        | 是        |
| --help, -h      |                            | 查询绑定命令帮助                                           |          |

## 例子

### 查询服务绑定

```shell
iriscli service binding --def-chain-id=test --service-name=test-service --bind-chain-id=test --provider=faa1ydhmma8l4m9dygsh7l08fgrwka6yczs0gkfnvd
```

运行成功以后，返回的结果如下:

```json
{
  "type": "iris-hub/service/SvcBinding",
  "value": {
    "def_name": "test-service",
    "def_chain_id": "test",
    "bind_chain_id": "test",
    "provider": "faa1ydhmma8l4m9dygsh7l08fgrwka6yczs0gkfnvd",
    "binding_type": "Local",
    "deposit": [
      {
        "denom": "iris-atto",
        "amount": "1000000000000000000000"
      }
    ],
    "price": [
      {
        "denom": "iris-atto",
        "amount": "1000000000000000000"
      }
    ],
    "level": {
      "avg_rsp_time": "10000",
      "usable_time": "100"
    },
    "available": true,
    "disable_height": "0"
  }
}
```