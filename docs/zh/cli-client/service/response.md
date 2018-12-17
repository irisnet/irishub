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
iriscli service response --request-chain-id=test-irishub --request-id=635-535-0
```

运行成功以后，返回的结果如下:
```txt
Committed at block 306 (tx hash: 5A4C6E00F4F6BF795EB05D2D388CBA0E8A6E6CF17669314B1EE6A31729A22450, response: {Code:0 Data:[] Log:Msg 0:  Info: GasWanted:200000 GasUsed:3398 Tags:[{Key:[97 99 116 105 111 110] Value:[115 101 114 118 105 99 101 45 119 105 116 104 100 114 97 119 45 102 101 101 115] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[99 111 109 112 108 101 116 101 67 111 110 115 117 109 101 100 84 120 70 101 101 45 105 114 105 115 45 97 116 116 111] Value:[34 54 55 57 54 48 48 48 48 48 48 48 48 48 48 48 34] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0}] Codespace: XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0})
```

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

