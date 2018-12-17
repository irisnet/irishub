# iriscli service call 

## 描述

调用服务方法

## 用法

```
iriscli service call [flags]
```

## 标志

| Name, shorthand       | Default                 | Description                          | Required |
| --------------------- | ----------------------- | ------------------------------------ | -------- |
| --def-chain-id        |                         | [string] 定义该服务的区块链ID           | 是       |
| --service-name        |                         | [string] 服务名称                     | 是       |
| --method-id           |                         | [int] 调用的服务方法ID                 | 是       |
| --bind-chain-id       |                         | [string] 绑定该服务的区块链ID           | 是       |
| --provider            |                         | [string] bech32编码的服务提供商账户地址  | 是       |
| --service-fee         |                         | [string] 服务调用支付的服务费            |          |
| --request-data        |                         | [string] hex编码的服务调用请求数据        |          |
| -h, --help            |                         | 调用命令帮助                             |          |

## 示例

### 发起一个服务调用请求
```shell
iriscli service call --chain-id=test-irishub --from=node0 --fee=0.004iris --def-chain-id=test-irishub --service-name=test-service --method-id=1 --bind-chain-id=test-irishub --provider=faa1qm54q9ta97kwqaedz9wzd90cacdsp6mq54cwda --service-fee=1iris --request-data=434355
```

运行成功以后，返回的结果如下:
```txt
Committed at block 306 (tx hash: 5A4C6E00F4F6BF795EB05D2D388CBA0E8A6E6CF17669314B1EE6A31729A22450, response: {Code:0 Data:[] Log:Msg 0:  Info: GasWanted:200000 GasUsed:3398 Tags:[{Key:[97 99 116 105 111 110] Value:[115 101 114 118 105 99 101 45 119 105 116 104 100 114 97 119 45 102 101 101 115] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[99 111 109 112 108 101 116 101 67 111 110 115 117 109 101 100 84 120 70 101 101 45 105 114 105 115 45 97 116 116 111] Value:[34 54 55 57 54 48 48 48 48 48 48 48 48 48 48 48 34] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0}] Codespace: XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0})
```

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

