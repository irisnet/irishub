---
order: 3
---

# Key Management System

## Introduction

KMS is short for Key Management System, please refer to the [Tendermint KMS](https://github.com/tendermint/kms) for more details.

## Building

Detailed build instructions can be found [here](https://github.com/tendermint/kms#installation).

:::tip
When compiling the KMS, ensure you have enabled the applicable features:
:::

| Backend               | Recommended Command line              |
|-----------------------|---------------------------------------|
| YubiHSM               | ```cargo build --features yubihsm```  |
| Ledger+Tendermint App | ```cargo build --features ledgertm``` |

## Configuration

[Tendermint KMS](https://github.com/tendermint/kms) supports all blockchains built on the [tendermint](https://github.com/tendermint/tendermint) consensus engine, including IRIShub.

To enable KMS, you need to edit the `priv_validator_laddr` in your `<iris-home>/config/config.toml` file first. E.g.:

```toml
# TCP or UNIX socket address for Tendermint to listen on for
# connections from an external PrivValidator process
priv_validator_laddr = "localhost:26658"
```

You can download the [example config file](https://github.com/tendermint/kms/blob/master/tmkms.toml.example) which supports IRIShub, you just need to edit it as follows:

- Edit `addr` to point to your `iris` instance.
- Adjust `chain-id` to match your `<iris-home>/config/genesis.json` settings.
- Edit `auth` to authorize access to your yubihsm.
- Edit `keys` to determine which pubkey you will be using.

Then start tmkms:

```bash
tmkms start
```

A KMS can be configured in various ways:

### Using a YubiHSM

Detailed information on how to setup a KMS with YubiHSM2 can be found [here](https://github.com/tendermint/kms/blob/master/README.yubihsm.md).

If you want to import an existing IRIShub private_key:

```bash
tmkms yubihsm keys import <iris_home>/config/priv_validator.json -i <id>
```

### Using a Ledger device running the Tendermint app

- [Using a Ledger device running the Tendermint Validator app](kms/kms_ledger.md)
