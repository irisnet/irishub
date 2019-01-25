# iriscli distribution set-withdraw-addr

## Description

Set withdraw address for delegator

## Usage

```
iriscli distribution set-withdraw-addr [withdraw-address] [flags]
```

Print help messages:
```
iriscli distribution set-withdraw-addr --help
```


## Examples

```
iriscli distribution set-withdraw-addr [withdraw address] --from <key name> --fee=0.004iris --chain-id=<chain-id>
```
Example response:
```json
{
   "code": 0,
   "data": null,
   "log": "Msg 0: ",
   "info": "",
   "gas_wanted": 200000,
   "gas_used": 3876,
   "codespace": "",
   "tags": [
     {
       "key": "action",
       "value": "set_withdraw_address"
     },
     {
       "key": "delegator",
       "value": "faa1yclscskdtqu9rgufgws293wxp3njsesxtplqxd"
     }
   ]
 })
```