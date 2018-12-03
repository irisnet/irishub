# iriscli service bind 

## 描述

创建一个新的服务绑定

## 用法

```
iriscli service bind [flags]
```

## 标志

| Name, shorthand       | Default                 | Description                                                                        | Required |
| --------------------- | ----------------------- | ---------------------------------------------------------------------------------- | -------- |
| --avg-rsp-time        |                         | [int]  服务平均返回时间的毫秒数表示                                                     | 是       |
| --bind-type           |                         | [string] 对服务是本地还是全局的设置，可选值Local/Global                                  | 是       |
| --def-chain-id        |                         | [string] 定义该服务的区块链ID                                                          | 是       |
| --deposit             |                         | [string] 服务提供者的保证金                                                            | 是       |
| --prices              |                         | [strings] 服务定价,按照服务方法排序的定价列表                                             |          |
| --service-name        |                         | [string] 服务名称                                                                    | 是       |
| --usable-time         |                         | [int] 每一万次服务调用的可用性的整数表示                                                  | 是       |
| -h, --help            |                         | 绑定命令帮助                                                                          |          |

## 示例

### 添加服务绑定到已存在的服务定义
```shell
iriscli service bind --chain-id=test --from=node0 --fee=0.004iris --service-name=test-service --def-chain-id=test --bind-type=Local --deposit=1iris --prices=1iris --avg-rsp-time=10000 --usable-time=100
```

运行成功以后，返回的结果如下:

```txt
Password to sign with 'node0':
Committed at block 6 (tx hash: 87A477AEA41B22F7294084B4794837211C43A297D73EABA2F42F6436F3D975DD, response: {Code:0 Data:[] Log:Msg 0:  Info: GasWanted:200000 GasUsed:5568 Tags:[{Key:[97 99 116 105 111 110] Value:[115 101 114 118 105 99 101 45 98 105 110 100] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[99 111 109 112 108 101 116 101 67 111 110 115 117 109 101 100 84 120 70 101 101 45 105 114 105 115 45 97 116 116 111] Value:[34 49 49 49 51 54 48 48 48 48 48 48 48 48 48 48 34] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0}] Codespace: XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0})
{
   "tags": {
     "action": "service-bind",
     "completeConsumedTxFee-iris-atto": "\"111360000000000\""
   }
 }
```

