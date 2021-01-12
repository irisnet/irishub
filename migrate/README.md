## migrate

### Step 1

Stop irishub v0.16.3 daemon and use irishub v0.16.3 to export mainnet state genesis with '--for-zero-height'
```bash
iris export --home <old_node_home> --for-zero-height
```

### Step 2
Use irishub v1.0.0 to migrate the exported genesis.json
```
iris migrate genesis.json --chain-id test > genesis_new.json
```

### Step 3
Use irishub v1.0.0 to initialize the new node
```
iris init [moniker] --home {$new_node_home}
```

### Step 4
Upgrade privkey file
```
go run migrate/scripts/privValUpgrade.go {$old_node_home}/config/priv_validator.json {$new_node_home}/config/priv_validator_key.json {$new_node_home}/data/priv_validator_state.json
```

### Step 5
Migrate node key file
```
cp {$old_node_home}/config/node_key.json {$new_node_home}/config/node_key.json
```

### Step 6
Copy new genesis.json to new node home
```
cp genesis_new.json {$new_node_home}/config/genesis.json
```

### Step 7
Start new node
```
iris unsafe-reset-all --home {new_node_home}
iris start --home {new_node_home}
```