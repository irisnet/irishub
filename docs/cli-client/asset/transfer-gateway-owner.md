# iriscli asset transfer-gateway-owner

## Introduction

Transfer the ownership of a gateway from the current owner to the new owner. The command is only used to generate the unsigned tx.

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
| --------------------| -----  | -------- | --------  -------------------------------------------------------- |
| --moniker           | string  | true     | ""       | the unique name of the gateway to be transferred       |
| --to                | Address | true     |          | the new owner to which the gateway will be transferred |


## Examples

```
iriscli asset transfer-gateway-owner --moniker=tgw --to=faa1anyffkfjhyvdawv2trlv4eq00c3xrjmqvldwal --chain-id=irishub 
--from=node0 --fee=0.3iris
```

Output:
```txt
{"type":"irishub/bank/StdTx","value":{"msg":[{"type":"irishub/asset/MsgTransferGatewayOwner","value":{"owner":"faa1an4wfvsnxrp97lug5fngct6melhgcuvdv2qye3","moniker":"testgw","to":"faa1anyffkfjhyvdawv2trlv4eq00c3xrjmqvldwal"}}],"fee":{"amount":[{"denom":"iris-atto","amount":"600000000000000000"}],"gas":"50000"},"signatures":null,"memo":""}}
```
