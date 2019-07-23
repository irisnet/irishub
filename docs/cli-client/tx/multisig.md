# iriscli tx multisign

## Description

Sign the same transaction by multiple accounts. The tx could be broadcast only when the number of signatures is greater than or equal to the minimum number of signatures.


## Usage
```
iriscli tx multisign <file> <key name> <[signature]...> <flags>
```

## Flags
| Name,shorthand | Type   | Required | Default               | Description                                                  |
| -------------- | ------ | -------- | --------------------- | ------------------------------------------------------------ |
| -h, --help       |        | false     |                       |  Help for sign                                        |

## Example

### Create multisig account 

Firstly, you must create a multisig account, please refer to [add](../keys/add.md)

Create a multisig account with 3 sub-accounts，specify the minimum number of signatures，such as 2. The tx could be broadcast only when the number of signatures is greater than or equal to 2.

```  
iriscli keys add <multi_account_keyname> --multisig-threshold=2 --multisig=<signer_keyname_1>,<signer_keyname_2>,<signer_keyname_3>...
```

::: tips
<signer_keyname> could be the type of "local/offline/ledger"， but not "multi" type。

Offline account can be created by "iriscli keys add --pubkey". 
:::

### Generate tx with multisig account

Create Tx and generate Tx-generate.json :
```  
iriscli bank send --amount=1iris --fee=0.3iris --chain-id=<chain-id> --from=<multi_account_keyname> --to=<address> --generate-only > Tx-generate.json
```

### Sign the tx separately

Use `iriscli keys show <multi_account_keyname>` to get `<multi_account_address>`

Specify the threshold to 2， sign and generate Tx-sign.json.

Sign the tx with signer_1:
```  
iriscli tx sign Tx-generate.json --name=<signer_keyname_1> --chain-id=<chain-id> --multisig=<multi_account_address> --signature-only >Tx-sign-user_1.json
```

Sign the tx with signer_2:
```  
iriscli tx sign Tx-generate.json --name=<signer_keyname_2> --chain-id=<chain-id> --multisig=<multi_account_address> --signature-only >Tx-sign-user_2.json
```

### multisign all the signatures

Create the signed tx and generate Tx-signed.json：

```  
iriscli tx multisign --chain-id=<chain-id> Tx-generate.json <multi_account_keyname> Tx-sign-user_1.json Tx-sign-user_2.json > Tx-signed.json
```


### Broadcast the signed tx

After signing a transaction, it could be broadcast to the network with [broadcast command](broadcast.md)

```  
iriscli tx broadcast Tx-signed.json --commit
```
