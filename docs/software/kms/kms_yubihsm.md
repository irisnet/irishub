# Setting up Tendermint KMS + YubiHSM

::: danger Warning
The following instructions are a brief walkthrough and not a comprehensive guideline.
:::

::: danger Warning
KMS is currently work in progress. Details may vary. Use with care under your own risk.
:::

## YubiHSM
[YubiHSM](https://www.yubico.com/products/yubihsm/): hardware security module providing root of trust for servers and computing devices.

## KMS configuration

In this section, we will configure a KMS to use YubiHSM. 

#### Config file

You can find other configuration examples [here](https://github.com/irisnet/kms/blob/master/tmkms.toml.example)

- Create a `~/.tmkms/tmkms.toml` file with the following content:

```toml
# Example KMS configuration file
[[validator]]
addr = "tcp://localhost:26658"    # or "unix:///path/to/socket"
chain_id = "fuxi"
reconnect = true # true is the default
secret_key = "~/.tmkms/secret_connection.key"
[[providers.yubihsm]]
adapter = { type = "usb" }
auth = { key = 1, password = "password" } # Default YubiHSM admin credentials. Change ASAP!
keys = [{ id = "test", key = 1 }]
#serial_number = "0123456789" # identify serial number of a specific YubiHSM to connect to
```

- Edit `addr` to point to your `iris` instance.
- Adjust `chain-id` to match your `.iris/config/config.toml` settings.
- Edit `auth` to authorize access to your yubihsm.
- Edit `keys` to determine which pubkey you will be using. [How to import key?](#import-private-key))

#### Import private key to yubihsm

```bash
tmkms yubihsm keys import  -p ~/.iris/config/priv_validator.json [id]
```

#### Generate connection secret key

Now you need to generate secret_key:

```bash
tmkms keygen ~/.tmkms/secret_connection.key
```

#### Start the KMS

The last step is to retrieve the validator key that you will use in `iris`.

Start the KMS:

```bash
tmkms start -c ~/.tmkms/tmkms.toml
```

The output should look similar to:

```text
07:28:24 [INFO] tmkms 0.3.0 starting up...
07:28:24 [INFO] [keyring:ledgertm:ledgertm] added validator key icp1zcjduepqa9y67dqgug4u4stf5sf0arvjrnty8eenlfj22vnh78cmejd8qdss8t6ljg
07:28:24 [INFO] KMS node ID: 1BC12314E2E1C29015B66017A397F170C6ECDE4A
```