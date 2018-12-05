# iriscli service refund-deposit 

## 描述

取回所有押金

## 用法

```
iriscli service refund-deposit [flags]
```

## 标志

| Name, shorthand       | Default                 | Description                                                                        | Required |
| --------------------- | ----------------------- | ---------------------------------------------------------------------------------- | -------- |
| --def-chain-id        |                         | [string] 定义该服务的区块链ID                                                         | 是       |
| --service-name        |                         | [string] 服务名称                                                                   | 是       |
| -h, --help            |                         | 取回押金命令帮助                                                                      |          |

## 示例

### 取回所有押金
```shell
iriscli service refund-deposit --chain-id=test  --from=node0 --fee=0.004iris --def-chain-id=test --service-name=test-service
```

运行成功以后，返回的结果如下:

```txt
Password to sign with 'node0':
Committed at block 991 (tx hash: 8A7F0EA61AB73A8B241945C8942EC8593774346B36BB70E36E138A23E7A473EE, response: {Code:0 Data:[] Log:Msg 0:  Info: GasWanted:200000 GasUsed:4614 Tags:[{Key:[97 99 116 105 111 110] Value:[115 101 114 118 105 99 101 45 114 101 102 117 110 100 45 100 101 112 111 115 105 116] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[99 111 109 112 108 101 116 101 67 111 110 115 117 109 101 100 84 120 70 101 101 45 105 114 105 115 45 97 116 116 111] Value:[34 57 50 50 56 48 48 48 48 48 48 48 48 48 48 34] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0}] Codespace: XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0})
{
   "tags": {
     "action": "service-refund-deposit",
     "completeConsumedTxFee-iris-atto": "\"92280000000000\""
   }
 }
```