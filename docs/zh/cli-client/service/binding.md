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

## 示例

### 查询服务绑定

```shell
iriscli service binding --def-chain-id=test-irishub --service-name=test-service --bind-chain-id=test-irishub --provider=faa1ydhmma8l4m9dygsh7l08fgrwka6yczs0gkfnvd
```

运行成功以后，返回的结果如下:
```txt
Committed at block 306 (tx hash: 5A4C6E00F4F6BF795EB05D2D388CBA0E8A6E6CF17669314B1EE6A31729A22450, response: {Code:0 Data:[] Log:Msg 0:  Info: GasWanted:200000 GasUsed:3398 Tags:[{Key:[97 99 116 105 111 110] Value:[115 101 114 118 105 99 101 45 119 105 116 104 100 114 97 119 45 102 101 101 115] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[99 111 109 112 108 101 116 101 67 111 110 115 117 109 101 100 84 120 70 101 101 45 105 114 105 115 45 97 116 116 111] Value:[34 54 55 57 54 48 48 48 48 48 48 48 48 48 48 48 34] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0}] Codespace: XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0})
```

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