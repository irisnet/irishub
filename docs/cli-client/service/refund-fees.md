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
Committed at block 79 (tx hash: 1E3A690028116E0DF541A840BF5830598EAD4154F4374B2A4042911C27D68C64, response:
 {
   "code": 0,
   "data": null,
   "log": "Msg 0: ",
   "info": "",
   "gas_wanted": 200000,
   "gas_used": 3912,
   "codespace": "",
   "tags": {
     "action": "service_refund_fees"
   }
 })
```

