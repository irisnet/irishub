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
iriscli service update-binding --chain-id=test-irishub --from=node0 --fee=0.004iris --service-name=test-service --def-chain-id=test-irishub --bind-type=Local --deposit=1iris --prices=1iris --avg-rsp-time=10000 --usable-time=100
```

After that, you're done with updating an existing service binding.
```txt
Committed at block 306 (tx hash: 5A4C6E00F4F6BF795EB05D2D388CBA0E8A6E6CF17669314B1EE6A31729A22450, response: {Code:0 Data:[] Log:Msg 0:  Info: GasWanted:200000 GasUsed:3398 Tags:[{Key:[97 99 116 105 111 110] Value:[115 101 114 118 105 99 101 45 119 105 116 104 100 114 97 119 45 102 101 101 115] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[99 111 109 112 108 101 116 101 67 111 110 115 117 109 101 100 84 120 70 101 101 45 105 114 105 115 45 97 116 116 111] Value:[34 54 55 57 54 48 48 48 48 48 48 48 48 48 48 48 34] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0}] Codespace: XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0})
```
```json
{
   "tags": {
     "action": "service-update-binding",
     "completeConsumedTxFee-iris-atto": "\"100840000000000\""
   }
 }
```

