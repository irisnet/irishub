## migrate

### Step 1

use irishub v0.16.3 to export mainnet state genesis with '--for-zero-height'
```bash
iris export --home <your_home> --for-zero-height
```

### Step 2
use irishub v1.0.0 to migrate the exported genesis.json by Step 1
```
iris migrate genesis.json --chain-id test > genesis_new.json
```

### Step 3
Upgrade privkey file
```
go run migrate/scripts/privValUpgrade.go {$node_home}/config/priv_validator.json {$node_home}/config/priv_validator.json {$node_home}/data/priv_validator_state.json
```

### Step 4
Copy new genesis.json to node home
```
cp genesis_new.json {$node_home}/config/genesis.json
```

### Step 5
reset node and start
```
iris unsafe-reset-all --home {node_home}
iris start --home {node_home}
```