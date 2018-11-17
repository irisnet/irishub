# iriscli service disable 

## Description

Disable a available service binding

## Usage

```
iriscli service disable [flags]
```

## Flags

| Name, shorthand       | Default                 | Description                                                                                                                                           | Required |
| --------------------- | ----------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --def-chain-id        |                         | [string] the ID of the blockchain defined of the service                                                                                              |  Yes     |
| --service-name        |                         | [string] service name                                                                                                                                 |  Yes     |
| -h, --help            |                         | help for disable                                                                                                                                      |          |

## Examples

### Disable a available service binding
```shell
iriscli service disable --chain-id=test  --from=node0 --fee=0.004iris --def-chain-id=test --service-name=test-service
```

After that, you're done with disabling a available service binding.

```txt
Password to sign with 'node0':
Committed at block 537 (tx hash: 66C5EE634339D168A07C073C6BF209D80081762EB8451974ABC33A41914A7158, response: {Code:0 Data:[] Log:Msg 0:  Info: GasWanted:200000 GasUsed:3510 Tags:[{Key:[97 99 116 105 111 110] Value:[115 101 114 118 105 99 101 45 100 105 115 97 98 108 101] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[99 111 109 112 108 101 116 101 67 111 110 115 117 109 101 100 84 120 70 101 101 45 105 114 105 115 45 97 116 116 111] Value:[34 55 48 50 48 48 48 48 48 48 48 48 48 48 48 34] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0}] Codespace: XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0})
{
   "tags": {
     "action": "service-disable",
     "completeConsumedTxFee-iris-atto": "\"70200000000000\""
   }
 }
```