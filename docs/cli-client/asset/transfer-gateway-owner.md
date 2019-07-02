# iriscli asset transfer-gateway-owner

## Introduction

Transfer the ownership of a gateway from the current owner to the new owner.

## Usage

```
iriscli asset transfer-gateway-owner [flags]
```

Print help messages:
```
iriscli asset transfer-gateway-owner --help
```

## Unique Flags

| Name, shorthand     | type   | Required | Default   | Description                                                       |
| --------------------| -----  | -------- | --------  |-------------------------------------------------------- |
| --moniker           | string  | true     | ""       | the unique name of the gateway to be transferred       |
| --to                | Address | true     |          | the new owner to which the gateway will be transferred |


## Examples

```
iriscli asset transfer-gateway-owner --moniker=tgw --to=faa1anyffkfjhyvdawv2trlv4eq00c3xrjmqvldwal --chain-id=irishub 
--from=node0 --fee=0.3iris
```

Output:
```json
{
   "code": 0,
   "data": null,
   "log": "Msg 0: ",
   "info": "",
   "gas_wanted": 50000,
   "gas_used": 4654,
   "codespace": "",
   "tags": [
     {
       "key": "action",
       "value": "transfer_gateway_owner"
     },
     {
       "key": "moniker",
       "value": "tgw"
     }
   ]
}
```
