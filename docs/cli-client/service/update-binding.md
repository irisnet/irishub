# iriscli service update-binding 

## Description

Update a service binding

## Usage

```
iriscli service update-binding [flags]
```

## Flags
| Name, shorthand       | Default                 | Description                                                                                                                                           | Required |
| --------------------- | ----------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --avg-rsp-time        |                         | [int] the average service response time in milliseconds                                                                                               |          |
| --bind-type           |                         | [string] type of binding, valid values can be Local and Global                                                                                        |          |
| --def-chain-id        |                         | [string] the ID of the blockchain defined of the service                                                                                              |  Yes     |
| --deposit             |                         | [string] deposit of binding, will add to the current deposit balance                                                                                  |          |
| --prices              |                         | [strings] prices of binding, will contains all method                                                                                                 |          |
| --service-name        |                         | [string] service name                                                                                                                                 |  Yes     |
| --usable-time         |                         | [int] an integer represents the number of usable service invocations per 10,000                                                                       |          |
| -h, --help            |                         | help for update-binding                                                                                                                               |          |

## Examples

### Update an existing service binding
```shell
iriscli service update-binding --chain-id=test --from=node0 --fee=0.004iris --service-name=test-service --def-chain-id=test --bind-type=Local --deposit=1iris --prices=1iris --avg-rsp-time=10000 --usable-time=100
```

After that, you're done with updating an existing service binding.

```txt
Password to sign with 'node0':
Committed at block 417 (tx hash: 8C9969A2BF3F7A8C13C2E0B57CE4FD7BE43454280559831D7E39B0FD3C1FCD28, response: {Code:0 Data:[] Log:Msg 0:  Info: GasWanted:200000 GasUsed:5042 Tags:[{Key:[97 99 116 105 111 110] Value:[115 101 114 118 105 99 101 45 117 112 100 97 116 101 45 98 105 110 100 105 110 103] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[99 111 109 112 108 101 116 101 67 111 110 115 117 109 101 100 84 120 70 101 101 45 105 114 105 115 45 97 116 116 111] Value:[34 49 48 48 56 52 48 48 48 48 48 48 48 48 48 48 34] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0}] Codespace: XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0})
{
   "tags": {
     "action": "service-update-binding",
     "completeConsumedTxFee-iris-atto": "\"100840000000000\""
   }
 }
```

