# iriscli service refund-deposit 

## Description

Refund all deposit from a service binding

## Usage

```
iriscli service refund-deposit [flags]
```

## Flags

| Name, shorthand       | Default                 | Description                                                                                                                                           | Required |
| --------------------- | ----------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --def-chain-id        |                         | [string] the ID of the blockchain defined of the service                                                                                              |  Yes     |
| --service-name        |                         |  [string] service name                                                                                                                                |  Yes     |
| -h, --help            |                         |  help for refund-deposit                                                                                                                              |          |

## Examples

### Refund all deposit from an unavailable service binding
```shell
iriscli service refund-deposit --chain-id=test-irishub  --from=node0 --fee=0.004iris --def-chain-id=test-irishub --service-name=test-service
```

After that, you're done with refunding all deposit from a service binding.
```txt
Committed at block 306 (tx hash: 5A4C6E00F4F6BF795EB05D2D388CBA0E8A6E6CF17669314B1EE6A31729A22450, response: {Code:0 Data:[] Log:Msg 0:  Info: GasWanted:200000 GasUsed:3398 Tags:[{Key:[97 99 116 105 111 110] Value:[115 101 114 118 105 99 101 45 119 105 116 104 100 114 97 119 45 102 101 101 115] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[99 111 109 112 108 101 116 101 67 111 110 115 117 109 101 100 84 120 70 101 101 45 105 114 105 115 45 97 116 116 111] Value:[34 54 55 57 54 48 48 48 48 48 48 48 48 48 48 48 34] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0}] Codespace: XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0})
```
```json
{
   "tags": {
     "action": "service-refund-deposit",
     "completeConsumedTxFee-iris-atto": "\"92280000000000\""
   }
 }
```