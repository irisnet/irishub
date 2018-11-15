# iriscli keys

## Description

Keys allows you to manage your local keystore for tendermint.

## Usage

```shell
iriscli keys [command]
```

## Available Commands

| Name          | Description                           |
| ------------- | ------------------------------------- |
| [add](add.md) | Create a new key, or import from seed |
| list          | List all keys                         |
| show          | Show key info for the given name      |
| delete        | Delete the given key                  |

## Flags

| Name, shorthand | Default | Description   | Required |
| --------------- | ------- | ------------- | -------- |
| --help, -h      |         | help for keys |          |

## Extended description

These keys may be in any format supported by go-crypto and can be used by light-clients, full nodes, or any other application that needs to sign with a private key.
