# Wasm

Wasm module allows you to manage wasm

## Available Commands

| Name                                | Description                        |
| ----------------------------------- | ---------------------------------- |
| [clear-contract-admin](#iris-tx-wasm-clear-contract-admin)       | Clears admin for a contract to prevent further migrations                  |
| [execute](#iris-tx-wasm-execute)         | Execute a command on a wasm contract             |
| [instantiate](#iris-tx-wasm-instantiate) | Instantiate a wasm contract  |
| [migrate](#iris-tx-wasm-migrate)         | Migrate a wasm contract to a new code version |
| [set-contract-admin](#iris-tx-wasm-set-contract-admin)    | Set new admin for a contract            |
| [store](#iris-tx-wasm-store)  | Upload a wasm binary              |
| [code](#iris-query-wasm-code)        | Downloads wasm bytecode for given code id       |
| [contract](#iris-query-wasm-contract)  | Prints out metadata of a contract given its address     |
| [contract-history](#iris-query-wasm-contract-history)  | Prints out the code history for a contract given its address     |
| [contract-state](#iris-query-wasm-contract-state)  | Querying commands for the wasm module     |
| [list-code](#iris-query-wasm-list-code)  | List all wasm bytecode on the chain     |
| [list-contract-by-code](#iris-query-wasm-list-contract-by-code)  | List wasm all bytecode on the chain for given code id     |

## iris tx wasm clear-contract-admin

Clears admin for a contract to prevent further migrations

```bash
iris tx wasm clear-contract-admin [contract_addr_bech32] [flags]
```

## iris tx wasm execute

Execute a command on a wasm contract

```bash
iris tx wasm execute [contract_addr_bech32] [json_encoded_send_args] [flags]
```

### Execute the smart contract

```bash
iris tx wasm execute iaa18vd8fpwxzck93qlwghaj6arh4p7c5n89fqcgm9 "$Vote" 
    --from node0 
    --chain-id=test 
    --keyring-backend=file 
    --home ./cschain/node0/cschaincli 
    --fees 6point 
    --gas="auto" -b block
```

## iris tx wasm instantiate

Instantiate a wasm contract

```bash
iris tx wasm instantiate [code_id_int64] [json_encoded_init_args] --label [text] --admin [address,optional] --amount [coins,optional] [flags]
```
**Flags:**

| Name, shorthand    | Type    | Required | Default       | Description                                                                                                                    |
| -------------------| ------- | -------- | ------------- | ------------------------------------------------------------------------------------------------------------------------------ |
| --label            | string  |          |               | A human-readable name for this contract in lists                                   |
| --admin            | string  |          |               | Address of an admin               |
| --amount           | string  |          |               | Coins to send to the contract during instantiation. |


### Initialize the smart contract

```bash
iris tx wasm instantiate "$CODE_ID" "$INIT" 
    --from node0 
    --label "mint iris" 
    --chain-id=test 
    --keyring-backend=file 
    --home ./cschain/node0/cschaincli 
    --fees 6point 
    --gas="auto" -b block
```

## iris tx wasm migrate

Migrate a wasm contract to a new code version

```bash
iris tx wasm migrate [contract_addr_bech32] [new_code_id_int64] [json_encoded_migration_args] [flags]
```
## iris tx wasm set-contract-admin

Set new admin for a contract

```bash
iris tx wasm set-contract-admin [contract_addr_bech32] [new_admin_addr_bech32] [flags]
```

### Set the admin of the smart contract

```bash
iris tx wasm set-contract-admin iaa18vd8fpwxzck93qlwghaj6arh4p7c5n89fqcgm9 iaa18lwh8r66wf2hc278ncu4mlgqcxh5slhudkuler 
    --from node0  
    --chain-id=test 
    --keyring-backend=file 
    --home ./cschain/node0/cschaincli 
    --fees 6point --gas="auto" -b block
```

## iris tx wasm store

Upload a wasm binary

```bash
iris tx wasm store [wasm file] --source [source] --builder [builder] [flags]
```

**Flags:**

| Name, shorthand            | Type    | Required | Default       | Description                                                                                                                    |
| ---------------------------| ------- | -------- | ------------- | ------------------------------------------------------------------------------------------------------------------------------ |
| --source                   | string  |          |               | A valid URI reference to the contract's source code                                   |
| --builder                  | string  |          |               | A valid docker tag for the build system, optional               |
| --instantiate-only-address | string  |          |               | Only this address can instantiate a contract instance from the code.                                    |
| --instantiate-everybody    | string  |          |               | Everybody can instantiate a contract from the code. |


### Upload a wasm binary

```bash
iris tx wasm store election.wasm 
    --from node0 
    --chain-id=test 
    --keyring-backend=file 
    --home ./cschain/node0/cschaincli 
    --fees 6point -b block
    --gas="auto"
```

## iris query wasm code

Downloads wasm bytecode for given code id

```bash
iris query wasm code [code_id] [output filename] [flags]
```

### Downloads wasm bytecode

```bash
iris query wasm code 1 election.wasm
```

## iris query wasm contract

Prints out metadata of a contract given its address

```bash
iris query wasm contract [bech32_address] [flags]
```

## iris query wasm contract-history

Prints out the code history for a contract given its address

```bash
iris query wasm contract-history [bech32_address] [flags]
```

## iris query wasm contract-state

Querying commands for the wasm module

```bash
iris query wasm contract-state [command]
iris query wasm contract-state [flags]
```

## iris query wasm list-code

List all wasm bytecode on the chain

```bash
iris query wasm list-code [flags]
```

## iris query wasm list-contract-by-code

List wasm all bytecode on the chain for given code id

```bash
iris query wasm list-contract-by-code [code_id] [flags]
```
