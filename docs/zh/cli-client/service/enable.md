# iriscli service enable 

## 描述

启用一个不可用的服务绑定

## 用法

```
iriscli service enable [flags]
```

## 标志

| Name, shorthand       | Default                 | Description                                                                       | Required |
| --------------------- | ----------------------- | --------------------------------------------------------------------------------- | -------- |
| --def-chain-id        |                         | [string] 定义该服务的区块链ID                                                         |  Yes     |
| --deposit string      |                         | [string] 绑定押金, 将会增加当前服务绑定押金                                             |          |
| --service-name        |                         | [string] 服务名称                                                                   |  Yes     |
| -h, --help            |                         | 启用命令帮助                                                                         |          |

## 例子

### 启用一个不可用的服务绑定
```shell
iriscli service enable --chain-id=test  --from=node0 --fee=0.004iris --def-chain-id=test --service-name=test-service
```

运行成功以后，返回的结果如下:

```txt
Committed at block 654 (tx hash: CF74E7629F0098AC3295F454F5C15BD5846A1F77C4E6C6FBA551606672B364DD, response: {Code:0 Data:[] Log:Msg 0:  Info: GasWanted:200000 GasUsed:5036 Tags:[{Key:[97 99 116 105 111 110] Value:[115 101 114 118 105 99 101 45 101 110 97 98 108 101] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[99 111 109 112 108 101 116 101 67 111 110 115 117 109 101 100 84 120 70 101 101 45 105 114 105 115 45 97 116 116 111] Value:[34 49 48 48 55 50 48 48 48 48 48 48 48 48 48 48 34] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0}] Codespace: XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0})
{
   "tags": {
     "action": "service-enable",
     "completeConsumedTxFee-iris-atto": "\"100720000000000\""
   }
 }
```