# iriscli service bind 

## Description

Create a new service binding

## Usage

```
iriscli service bind [flags]
```

## Flags

| Name, shorthand       | Default                 | Description                                                                                                                                           | Required |
| --------------------- | ----------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --avg-rsp-time        |                         | [int] the average service response time in milliseconds                                                                                               |  Yes     |
| --bind-type           |                         | [string] type of binding, valid values can be Local and Global                                                                                        |  Yes     |
| --def-chain-id        |                         | [string] the ID of the blockchain defined of the service                                                                                              |  Yes     |
| --deposit             |                         | [string] deposit of binding                                                                                                                           |  Yes     |
| --prices              |                         | [strings] prices of binding, will contains all method                                                                                                 |          |
| --service-name        |                         | [string] service name                                                                                                                                 |  Yes     |
| --usable-time         |                         | [int] an integer represents the number of usable service invocations per 10,000                                                                       |  Yes     |
| -h, --help            |                         | help for bind                                                                                                                                         |          |

## Examples

### Add a binding to an existing service definition
```shell
iriscli service bind --chain-id=test --from=node0 --fee=0.004iris --service-name=test-service --def-chain-id=test --bind-type=Local --deposit=1iris --prices=1iris --avg-rsp-time=10000 --usable-time=100
```

After that, you're done with adding a binding to an existing service definition.

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

