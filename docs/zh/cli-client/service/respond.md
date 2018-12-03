# iriscli service respond 

## Description

响应服务调用

## Usage

```
iriscli service respond [flags]
```

## Flags

| Name, shorthand       | Default                 | Description                                                                                                                                           | Required |
| --------------------- | ----------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --request-chain-id    |                         | [string] 发起该服务调用的区块链ID                                                                                              |  Yes     |
| --request-id          |                         | [string] 该服务调用的ID                                                                                                                                |  Yes     |
| --response-data       |                         | [string] hex编码的服务调用响应数据                                                                       |  Yes     |
| -h, --help            |                         | 响应命令帮助                                                                                                                                         |          |

## Examples

### 响应一个服务调用 
```shell
iriscli service respond --chain-id=test --from=node0 --fee=0.004iris --request-chain-id=test --request-id=230-130-0 --response-data=abcd
```

运行成功以后，返回的结果如下:

```txt
Committed at block 130 (tx hash: DB40CE593198FC1B112337C613934F4E325F0718770D40616473369090327994, response: {Code:0 Data:[] Log:Msg 0:  Info: GasWanted:200000 GasUsed:8144 Tags:[{Key:[97 99 116 105 111 110] Value:[115 101 114 118 105 99 101 45 99 97 108 108] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[114 101 113 117 101 115 116 45 105 100] Value:[50 51 48 45 49 51 48 45 48] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[112 114 111 118 105 100 101 114] Value:[102 97 97 49 102 48 50 101 120 116 57 100 117 107 55 104 51 114 120 57 122 109 55 97 118 48 112 110 108 101 103 120 118 101 56 110 101 53 118 119 54 120] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[99 111 110 115 117 109 101 114] Value:[102 97 97 49 102 48 50 101 120 116 57 100 117 107 55 104 51 114 120 57 122 109 55 97 118 48 112 110 108 101 103 120 118 101 56 110 101 53 118 119 54 120] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[99 111 109 112 108 101 116 101 67 111 110 115 117 109 101 100 84 120 70 101 101 45 105 114 105 115 45 97 116 116 111] Value:[34 49 54 50 56 56 48 48 48 48 48 48 48 48 48 48 34] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0}] Codespace: XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0})
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

