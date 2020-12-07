# Wasm

wasm模块用于管理你在IRISHub的智能合约。

## 可用命令

| 名称                                | 描述                     |
| ----------------------------------- | ------------------------ |
| [clear-contract-admin](#iris-tx-wasm-clear-contract-admin)       | 清除管理员的合约，以防止进一步的迁移                  |
| [execute](#iris-tx-wasm-execute)         | 在wasm合约上执行命令             |
| [instantiate](#iris-tx-wasm-instantiate) | 实例化wasm合约  |
| [migrate](#iris-tx-wasm-migrate)         | 将wasm合同迁移到新的代码版本 |
| [set-contract-admin](#iris-tx-wasm-set-contract-admin)    | 为合约设置新管理员            |
| [store](#iris-tx-wasm-store)  | 上传wasm二进制文件            |
| [code](#iris-query-wasm-code)        | 下载wasm字节码以获得给定的代码ID       |
| [contract](#iris-query-wasm-contract)  | 根据合约的地址打印出合约的元数据     |
| [contract-history](#iris-query-wasm-contract-history)  | 根据合约的地址打印出合约的代码历史记录     |
| [contract-state](#iris-query-wasm-contract-state)  | 查询wasm模块的命令     |
| [list-code](#iris-query-wasm-list-code)  | 列出链上的所有wasm字节码     |
| [list-contract-by-code](#iris-query-wasm-list-contract-by-code)  | 列出wasm链上所有给定代码ID的字节码    |

## iris tx wasm clear-contract-admin

清除管理员的合约，以防止进一步的迁移 

```bash
iris tx wasm clear-contract-admin [contract_addr_bech32] [flags]
```

## iris tx wasm execute

执行智能合约中的方法

```bash
iris tx wasm execute [contract_addr_bech32] [json_encoded_send_args] [flags]
```

### 执行合约

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

实例化wasm合约

```bash
iris tx wasm instantiate [code_id_int64] [json_encoded_init_args] --label [text] --admin [address,optional] --amount [coins,optional] [flags]
```

### 初始化合约

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

将wasm合同迁移到新的代码版本

```bash
iris tx wasm migrate [contract_addr_bech32] [new_code_id_int64] [json_encoded_migration_args] [flags]
```
## iris tx wasm set-contract-admin

为合约设置新管理员

```bash
iris tx wasm set-contract-admin [contract_addr_bech32] [new_admin_addr_bech32] [flags]
```

### 变更管理员

```bash
iris tx wasm set-contract-admin iaa18vd8fpwxzck93qlwghaj6arh4p7c5n89fqcgm9 iaa18lwh8r66wf2hc278ncu4mlgqcxh5slhudkuler 
    --from node0  
    --chain-id=test 
    --keyring-backend=file 
    --home ./cschain/node0/cschaincli 
    --fees 6point --gas="auto" -b block
```

## iris tx wasm store

上传wasm二进制文件

```bash
iris tx wasm store [wasm file] --source [source] --builder [builder] [flags]
```

**标识：**

| 名称，速记       | 类型    | 必须 | 默认          | 描述                                                               |
| ---------------- | ------- | ---- | ------------- | ------------------------------------------------------------------ |
| --source                   | string  |          |               | 智能合约源码地址                                   |
| --builder                  | string  |          |               | 合法的docker标签               |
| --instantiate-only-address | string  |          |               | 只有该地址可以初始化合约                                    |
| --instantiate-everybody    | string  |          |               | 任何人都可以初始化合约 |


### 上传合约

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

下载wasm字节码以获得给定的代码ID

```bash
iris query wasm code [code_id] [output filename] [flags]
```

### 下载wasm字节码

```bash
iris query wasm code 1 election.wasm
```

## iris query wasm contract

根据合约的地址打印出合约的元数据

```bash
iris query wasm contract [bech32_address] [flags]
```

## iris query wasm contract-history

根据合约的地址打印出合约的代码历史记录

```bash
iris query wasm contract-history [bech32_address] [flags]
```

## iris query wasm contract-state

查询wasm模块的命令

```bash
iris query wasm contract-state [command]
iris query wasm contract-state [flags]
```

## iris query wasm list-code

列出链上的所有wasm字节码 

```bash
iris query wasm list-code [flags]
```

## iris query wasm list-contract-by-code

列出wasm链上所有给定代码ID的字节码

```bash
iris query wasm list-contract-by-code [code_id] [flags]
```
