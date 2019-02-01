# iriscli bank send

## Description

Sending tokens to another address. 

## Usage:

```
iriscli bank send --to=<account address> --from <key name> --fee=0.4iris --chain-id=<chain-id> --amount=10iris
```

 

## Flags

| Name,shorthand   | Type   | Required | Default               | Description                                                  |
| ---------------- | ------ | -------- | --------------------- | ------------------------------------------------------------ |
| --amount         | String | True     |                       | Amount of coins to send, for instance: 10iris                |
| --to             | String |          |                       | Bech32 encoding address to receive coins                     |


## Examples

### Send token to a address 

```
 iriscli bank send --to=faa19aamjx3xszzxgqhrh0yqd4hkurkea7f6d429yx  --from=test  --fee=0.4iris --chain-id=test-irishub --amount=10iris
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
      "recipient": "faa1893x4l2rdshytfzvfpduecpswz7qtpstpr9x4h",
      "sender": "faa106nhdckyf996q69v3qdxwe6y7408pvyvufy0x2"
    }
  })
```
### Common Issues

* Wrong password

```$xslt
ERROR: Ciphertext decryption failed
```
