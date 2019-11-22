# iriscli tx

Tx module allows you to sign or broadcast transactions

## Available Commands

| Name                               | Description                                    |
| ---------------------------------- | ---------------------------------------------- |
| [sign](#iriscli-tx-sign)           | Sign transactions generated offline            |
| [broadcast](#iriscli-tx-broadcast) | Broadcast a signed transaction to the network  |
| [multisig](#iriscli-tx-multisign)  | Sign the same transaction by multiple accounts |

## iriscli tx sign

Sign transactions in generated offline file. The file created with the --generate-only flag.

```bash
iriscli tx sign <file> <flags>
```

### Flags

| Name, shorthand | Type   | Required | Default | Description                                                                              |
| --------------- | ------ | -------- | ------- | ---------------------------------------------------------------------------------------- |
| --append        | bool   | true     | true    | Attach a signature to an existing signature.                                             |
| --name          | string | true     |         | Key name for signature                                                                   |
| --offline       | bool   | true     |         | Offline mode.                                                                            |
| --print-sigs    | bool   | true     |         | Print the address where the transaction must be signed and the signed address, then exit |

### Generate an offline tx

:::tip
You can generate any type of txs offline by appending the flag `--generate-only`
:::

We use a transfer tx in the following examples:

```bash
iriscli bank send --to=<address> --amount=10iris --from=<key-name> --fee=0.3iris --chain-id=irishub --generate-only > unsigned.json
```

The `unsigned.json` should look like:

```json
{"type":"auth/StdTx","value":{"msg":[{"type":"cosmos-sdk/Send","value":{"inputs":[{"address":"iaa19aamjx3xszzxgqhrh0yqd4hkurkea7f646vaym","coins":[{"denom":"iris-atto","amount":"10000000000000000000"}]}],"outputs":[{"address":"iaa19aamjx3xszzxgqhrh0yqd4hkurkea7f646vaym","coins":[{"denom":"iris-atto","amount":"10000000000000000000"}]}]}}],"fee":{"amount":[{"denom":"iris-atto","amount":"4000000000000000"}],"gas":"200000"},"signatures":null,"memo":""}}
```

### Sign tx offline

```bash
iriscli tx sign unsigned.json --name=<key-name> > signed.tx
```

The `signed.json` should look like:

```json
{"type":"auth/StdTx","value":{"msg":[{"type":"cosmos-sdk/Send","value":{"inputs":[{"address":"iaa106nhdckyf996q69v3qdxwe6y7408pvyvyxzhxh","coins":[{"denom":"iris-atto","amount":"10000000000000000000"}]}],"outputs":[{"address":"iaa1893x4l2rdshytfzvfpduecpswz7qtpstevr742","coins":[{"denom":"iris-atto","amount":"10000000000000000000"}]}]}}],"fee":{"amount":[{"denom":"iris-atto","amount":"40000000000000000"}],"gas":"200000"},"signatures":[{"pub_key":{"type":"tendermint/PubKeySecp256k1","value":"Auouudrg0P86v2kq2lykdr97AJYGHyD6BJXAQtjR1gzd"},"signature":"sJewd6lKjma49rAiGVfdT+V0YYerKNx6ZksdumVCvuItqGm24bEN9msh7IJ12Sil1lYjqQjdAcjVCX/77FKlIQ==","account_number":"0","sequence":"3"}],"memo":"test"}}
```

Note the `signature` in the `signed.json` should no longer be empty after signing.

Now it's ready to [broadcast the signed tx](#iriscli-tx-broadcast) to the IRIS Hub.

## iriscli tx broadcast

This command is used to broadcast an offline signed transaction to the network.

### Broadcast offline signed transaction

```bash
iriscli tx broadcast signed.json --chain-id=irishub
```

## iriscli tx multisign

Sign a transaction by multiple accounts. The tx could be broadcasted only when the number of signatures meets the multisig-threshold.

```bash
iriscli tx multisign <file> <key-name> <[signature]...> <flags>
```

### Generate an offline tx by multisig key

:::tip
No multisig key? [Create one](keys.md#create-a-multisig-key)
:::

```bash
iriscli bank send --to=<address> --amount=10iris --fee=0.3iris --chain-id=irishub --from=<multisig-keyname> --generate-only > unsigned.json
```

### Sign the multisig tx

#### Query the multisig address

```bash
iriscli keys show <multisig-keyname>
```

#### Sign the `unsigned.json`

Assume the multisig-threshold is 2, here we sign the `unsigned.json` by 2 of the signers

Sign the tx by signer-1:

```bash
iriscli tx sign unsigned.json --name=<signer-keyname-1> --chain-id=irishub --multisig=<multisig-address> --signature-only > signed-1.json
```

Sign the tx by signer-2:

```bash
iriscli tx sign unsigned.json --name=<signer_keyname_2> --chain-id=irishub --multisig=<multisig-address> --signature-only > signed-2.json
```

#### Merge the signatures

Merge all the signatures into `signed.json`

```bash
iriscli tx multisign --chain-id=irishub unsigned.json <multisig-keyname> signed-1.json signed-2.json > signed.json
```

Now you can [broadcast the signed tx](#iriscli-tx-broadcast).
