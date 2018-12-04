# iriscli service refund-fees 

## Description

Refund all fees from service return fees

## Usage

```
iriscli service refund-fees [flags]
```

## Flags

| Name, shorthand       | Default                 | Description                                                                                                                                           | Required |
| --------------------- | ----------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| -h, --help            |                         | help for refund-fees                                                                                                                                         |          |

## Examples

### Refund fees from service return fees 
```shell
iriscli service refund-fees --chain-id=test --from=node0 --fee=0.004iris
```

After that, you're done with refunding fees from service return fees.

```txt
Committed at block 450 (tx hash: B423D2D34710DEAA1D88AC376FEAD7935B28E63EEA9EACE6F7A7E654126CF877, response: {Code:0 Data:[] Log:Msg 0:  Info: GasWanted:200000 GasUsed:3398 Tags:[{Key:[97 99 116 105 111 110] Value:[115 101 114 118 105 99 101 45 114 101 102 117 110 100 45 102 101 101 115] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[99 111 109 112 108 101 116 101 67 111 110 115 117 109 101 100 84 120 70 101 101 45 105 114 105 115 45 97 116 116 111] Value:[34 54 55 57 54 48 48 48 48 48 48 48 48 48 48 48 34] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0}] Codespace: XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0})
{
   "tags": {
     "action": "service-refund-fees",
     "completeConsumedTxFee-iris-atto": "\"679600000000000\""
   }
 }
```

