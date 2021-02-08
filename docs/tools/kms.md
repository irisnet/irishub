---
order: 3
---

# Key Management System

## Introduction

[Tendermint KMS](https://github.com/iqlusioninc/tmkms) is a key management service that allows separating key management from Tendermint nodes. In addition it provides other advantages such as:

- Improved security and risk management policies
- Unified API and support for various HSM (hardware security modules)
- Double signing protection (software or hardware based)

## Building

Detailed build instructions can be found [here](https://github.com/iqlusioninc/tmkms#installation).

:::tip
When compiling the KMS, ensure you have enabled the `yubihsm` features:
:::

```bash
cargo install tmkms --features=yubihsm --version=0.10.0-beta2
```

## Initialization

Initialize configuration files for IRISHub

```bash
tmkms init -n irishub /path/to/kms/home
```

## Configuration

To enable KMS, you need to edit the `priv_validator_laddr` in your `<iris-home>/config/config.toml` file first. E.g.:

```toml
# TCP or UNIX socket address for Tendermint to listen on for
# connections from an external PrivValidator process
priv_validator_laddr = "localhost:26658"
```

Then, downLoad [priv_validator_state.json example](https://github.com/irisnet/irishub/blob/master/docs/tools/priv_validator_state.json) and modify all field values to match your `<iris-home>/data/priv_validator_state.json` values.

Next, you just need to edit the configuration file `/path/to/kms/home/tmkms.toml` as follows:

- Configure `state_file` as the `priv_validator_state.json` completed in the previous step.
- Write your Yubihsm password to file `yubihsm-password.txt` and configure `password_file` as it.
- Edit `addr` to point to your `iris` instance(note: no need to specify the connection id, just like tcp://localhost:26658).
- Adjust `chain_id` to match your `<iris-home>/config/genesis.json` settings.
- Edit `auth` to authorize access to your Yubihsm.
- Edit `keys` to determine which pubkey you will be using.
- Edit `protocol_version` to v0.34.

Then start tmkms:

```bash
tmkms start -c /path/to/kms/home/tmkms.toml
```

### Using a YubiHSM

Detailed information on how to setup a KMS with YubiHSM2 can be found [here](https://github.com/iqlusioninc/tmkms/blob/master/README.yubihsm.md).

If you want to import an existing IRIShub private_key:

```bash
tmkms yubihsm keys import <iris_home>/config/priv_validator.json -i <id> -t json -c /path/to/kms/home/tmkms.toml
```
