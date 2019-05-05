# iriscli bank send

## Description

Sending tokens to another address， this command includes `create/sign/broadcast` transaction.

## Usage:

Send 10iris
```
iriscli bank send --to=<address> --from=<key_name> --fee=0.3iris --chain-id=<chain-id> --amount=10iris
```
## Flags

| Name,shorthand   | Type   | Required | Default               | Description                                                  |
| ---------------- | ------ | -------- | --------------------- | ------------------------------------------------------------ |
| --amount         | String | True     |                       | Amount of coins to send, for instance: 10iris                |
| --to             | String |          |                       | Bech32 encoding address to receive coins                     |


## Examples

### Send token to a address 

```
 iriscli bank send --to=<address> --from=<key_name> --fee=0.3iris --chain-id=<chain-id> --amount=10iris
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
      "action": "send",
      "recipient": "iaa1893x4l2rdshytfzvfpduecpswz7qtpstevr742",
      "sender": "iaa106nhdckyf996q69v3qdxwe6y7408pvyvyxzhxh"
    }
  })
```

