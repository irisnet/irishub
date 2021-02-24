# Migrate IRIShub from v0.16 to v1.0.1

## 1. Export genesis file

Stop irishub v0.16 daemon and use `irishub v0.16.4 (which fixed bugs in export)` to export mainnet state genesis with `--for-zero-height` at the upgrade block height

```bash
iris export --home [v0.16_node_home] --height 9146455 --for-zero-height
```

## 2. Migrate genesis file

Migrate the exported genesis.json with irishub v1.0.1

```bash
iris migrate genesis.json --chain-id irishub-1 > genesis_v1.0.1.json
```

Specify the upgrade height + 1 as the initial height of irishub v1.0.1

```bash
# export initial_height=$[${upgrade block height} + 1]
export initial_height=9146456
jq --arg v "$initial_height" '.initial_height=$v' genesis_v1.0.1.json | sponge genesis_v1.0.1.json
```

Check if sha256sum is correct

```bash
sha256sum genesis_v1.0.1.json
```

## 3. Initialize new node

Initialize the new node with irishub v1.0.1

```bash
iris init [moniker] --home [v1.0.1_node_home]
```

## 4. Migrate privkey file

Migrate privkey file with irishub v1.0.1.

- KMS user
If you are using KMS to deploy node, please upgrade `tmkms` first, and then modify the relevant configuration. Please refer to the [kms](../tools/kms.md) for details

- Not KMS user
If you are not using KMS to deploy node, and the node configuration file exists, you can use either one of the following two ways:

  - Rename file
  
    ```bash
    cp [v0.16_node_home]/config/priv_validator.json [v1.0.1_node_home]/config/priv_validator_key.json
    ```

  - Use script

    ```bash
    go run migrate/scripts/privValUpgrade.go [v0.16_node_home]/config/priv_validator.json [v1.0.1_node_home]/config/priv_validator_key.json [v1.0.1_node_home]/data/priv_validator_state.json
    ```

## 5. Migrate node key file

Migrate node key file with irishub v1.0.1

```bash
cp [v0.16_node_home]/config/node_key.json [v1.0.1_node_home]/config/node_key.json
```

## 6. Copy migrated genesis file

Copy genesis_v1.0.1.json to new node home

```bash
cp genesis_v1.0.1.json [v1.0.1_node_home]/config/genesis.json
```

## 7. Config new node

Configure `minimum-gas-prices` in `[v1.0.1_node_home]/config/app.toml`

```toml

# The minimum gas prices a validator is willing to accept for processing a
# transaction. A transaction's fees must meet the minimum of any denomination
# specified in this config (e.g. 0.25token1;0.0001token2).
minimum-gas-prices = "0.2uiris"

```

Copy `persistent_peers` in `[v0.16_node_home]/config/config.toml` to `[v1.0.1_node_home]/config/config.toml`

```toml

# Comma separated list of nodes to keep persistent connections to
persistent_peers = ""

```

And configure other fields refer to `[v0.16_node_home]/config/config.toml`

## 8. Start new node

Start new node with irishub v1.0.1

```bash
iris unsafe-reset-all --home [v1.0.1_node_home]
iris start --home [v1.0.1_node_home]
```
