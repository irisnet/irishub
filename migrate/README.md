# Migrate IRISHub to v1.0.0 from v0.16.3

## 1. Export genesis file

Stop irishub v0.16.3 daemon and use irishub v0.16.3 to export mainnet state genesis with `--for-zero-height` at the upgrade block height

```bash
iris export --home [v0.16.3_node_home] --height [upgrade-height] --for-zero-height
```

## 2. Migrate genesis file

Migrate the exported genesis.json with irishub v1.0.0

```bash
iris migrate genesis.json --chain-id [chain-id] > genesis_v1.0.0.json
```

Check if md5 is correct

```bash
md5sum genesis_new.json
```

## 3. Initialize new node

Initialize the new node with irishub v1.0.0

```bash
iris init [moniker] --home [v1.0.0_node_home]
```

## 4. Migrate privkey file

Migrate privkey file with irishub v1.0.0

```bash
go run migrate/scripts/privValUpgrade.go [v0.16.3_node_home]/config/priv_validator.json [v1.0.0_node_home]/config/priv_validator_key.json [v1.0.0_node_home]/data/priv_validator_state.json
```

## 5. Migrate node key file

Migrate node key file with irishub v1.0.0

```bash
cp [v0.16.3_node_home]/config/node_key.json [v1.0.0_node_home]/config/node_key.json
```

## 6. Copy migrated genesis file

Copy new genesis.json to new node home

```bash
cp genesis_v1.0.0.json [v1.0.0_node_home]/config/genesis.json
```

## 7. Config new node

Add one step to modify minimum-gas-prices and migrate config.toml (eg.persistent_peers)

```bash

```

## 8. Start new node

Start new node with irishub v1.0.0

```bash
iris unsafe-reset-all --home [v1.0.0_node_home]
iris start --home [v1.0.0_node_home]
```
