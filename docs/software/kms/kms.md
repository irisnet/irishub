# KMS - Key Management System

## What is a KMS?

Please refer to [kms](https://github.com/tendermint/kms).

## Building

Detailed build instructions can be found [here](https://github.com/tendermint/kms#installation).

::: tip
When compiling the KMS, ensure you have enabled the applicable features:
:::

| Backend               | Recommended Command line              |
|-----------------------|---------------------------------------|
| YubiHSM               | ```cargo build --features yubihsm```  |
| Ledger+Tendermint App | ```cargo build --features ledgertm``` |

## Configuration

[tendermint/kms](https://github.com/tendermint/kms) supports all blockchains built on [tendermint](https://github.com/tendermint/tendermint) consensus engine, including IRIShub.

If you want to enable KMS, you need to edit `priv_validator_laddr` in your `<iris_home>/config/config.toml` file first. E.g.:

```text
# TCP or UNIX socket address for Tendermint to listen on for
# connections from an external PrivValidator process
priv_validator_laddr = "localhost:26658"
```

You can download the [example config file](https://github.com/tendermint/kms/blob/master/tmkms.toml.example) with support for IRIShub, you just have to edit it as follows:

- Edit `addr` to point to your `iris` instance.
- Adjust `chain-id` to match your `<iris_home>/config/genesis.json` settings.
- Edit `auth` to authorize access to your yubihsm.
- Edit `keys` to determine which pubkey you will be using.

Then start tmkms:
```bash
tmkms start
```

A KMS can be configured in various ways:

### Using a YubiHSM
Detailed information on how to setup a KMS with YubiHSM2 can be found [here](https://github.com/tendermint/kms/blob/master/README.yubihsm.md).

If you want to import IRIShub private_key that already exists, you can:
```bash
tmkms yubihsm keys import <iris_home>/config/priv_validator.json -i <id>
``` 

### Using a Ledger device running the Tendermint app
- [Using a Ledger device running the Tendermint Validator app](kms_ledger.md)
