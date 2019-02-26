# KMS - Key Management System

....

## What is a KMS?

...

## Building

Detailed build instructions can be found [here](https://github.com/irisnet/kms#installation).

::: tip
When compiling the KMS, ensure you have enabled the applicable features:
:::

| Backend               | Recommended Command line              |
|-----------------------|---------------------------------------|
| YubiHSM               | ```cargo build --features yubihsm```  |
| Ledger+Tendermint App | ```cargo build --features ledgertm``` |

## Configuration

If you want to enable KMS, you need to edit `priv_validator_laddr` in your `~/.iris/config/config.toml` file first. E.g.：
```text
# TCP or UNIX socket address for Tendermint to listen on for
# connections from an external PrivValidator process
priv_validator_laddr = "localhost:26658"
```

The KMS provides different alternatives

- [Using a CPU-based signer](kms_cpu.md)
- [Using a YubiHSM](kms_yubihsm.md)
- [Using a Ledger device running the Tendermint Validator app](kms_ledger.md)
