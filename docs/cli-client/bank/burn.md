# iriscli bank burn

## Description

This command is used for burning tokens from specified address. Anyone can send this transaction type. 
## Usage:

```
iriscli bank burn --from <key name> --fee=0.3iris --chain-id=<chain-id> --amount=10iris
```

## Flags

| Name,shorthand   | Type   | Required | Default               | Description                                                  |
| ---------------- | ------ | -------- | --------------------- | ------------------------------------------------------------ |
| --amount         | String | True     |                       | Amount of coins to burn, for instance: 10iris                |


## Examples

### Burn token 

```
 iriscli bank burn --from=test  --fee=0.3iris --chain-id=<chain-id> --amount=10iris --commit
```

After that, you will get the detail info for the send

```
[Committed at block 87 (tx hash: AEA8E49C1BC9A81CAFEE8ACA3D0D96DA7B5DC43B44C06BACEC7DCA2F9C4D89FC, response:
  {
    "code": 0,
    "data": null,
    "log": "Msg 0: ",
    "info": "",
    "gas_wanted": 200000,
    "gas_used": 3839,
    "codespace": "",
    "tags": {
      "action": "burn",
      "burnFrom": "iaa106nhdckyf996q69v3qdxwe6y7408pvyvufy0x2"
      "burnAmount": "10iris"
    }
  })
```
### Common Issues

* Wrong password

```$xslt
ERROR: Ciphertext decryption failed
```
