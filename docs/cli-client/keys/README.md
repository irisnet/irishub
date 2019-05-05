# iriscli keys

## Description

Keys allows you to manage your local keystore for tendermint.

## Usage

```shell
iriscli keys <command>
```

## Available Commands

| Name                    | Description                                                                                  |
| ----------------------- | -------------------------------------------------------------------------------------------- |
| [mnemonic](mnemonic.md) | Create a bip39 mnemonic, sometimes called a seed phrase, by reading from the system entropy. |
| [new](new.md)           | Derive a new private key using an interactive command that will prompt you for each input.   |
| [add](add.md)           | Create a new key, or import from seed                                                        |
| [list](list.md)         | List all keys                                                                                |
| [show](show.md)         | Show key info for the given name                                                             |
| [delete](delete.md)     | Delete the given key                                                                         |
| [update](update.md)     | Change the password used to protect private key                                              |

## Flags

| Name, shorthand | Default | Description   | Required |
| --------------- | ------- | ------------- | -------- |
| --help, -h      |         | Help for keys |          |

## Global Flags

| Name, shorthand | Default        | Description                            | Required |
| --------------- | -------------- | -------------------------------------- | -------- |
| --encoding, -e  | hex            | Binary encoding (hex/b64/btc) |          |
| --home          | $HOME/.iriscli | Directory for config and data |          |
| --output, -o    | text           | Output format (text/json)     |          |
| --trace         |                | Print out full stack trace on errors   |          |

## Extended description

These keys may be in any format supported by go-crypto and can be used by light-clients, full nodes, or any other application that needs to sign with a private key.
