# iriscli asset edit-gateway

## Introduction

Edit a gateway with the given moniker

## Usage

```
iriscli asset edit-gateway [flags]
```

Print help messages:
```
iriscli asset edit-gateway --help
```

## Unique Flags

| Name, shorthand     | type   | Required | Default  | Description                                                         |
| --------------------| -----  | -------- | -------- | ------------------------------------------------------------------- |
| --moniker           | string | true     | ""       | the unique name with a size between 3 and 8 letters      |
| --identity          | string | false    | ""       | Optional identity signature with a maximum length of 128 |
| --details           | string | false    | ""       | Optional website with a maximum length of 280            |
| --website           | string | false    | ""       | Optional website with a maximum length of 128            |


## Examples

```
iriscli asset edit-gateway --moniker=tgw --identity=exchange --details=test --website=http://gateway.io --from=node0 --chain-id=irishub --fee=0.4iris --home=iriscli --commit
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
       "value": "edit_gateway"
     },
     {
       "key": "moniker",
       "value": "tgw"
     }
   ]
}
```