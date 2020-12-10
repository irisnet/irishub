# Wasm

wasm模块用于管理你在IRISHub的智能合约。

## 可用命令

| 名称                                | 描述                     |
| ----------------------------------- | ------------------------ |
| [store](#iris-tx-wasm-store)  | 上传wasm二进制文件            |
| [instantiate](#iris-tx-wasm-instantiate) | 实例化wasm合约  |
| [execute](#iris-tx-wasm-execute)         | 在wasm合约上执行命令             |
| [migrate](#iris-tx-wasm-migrate)         | 将wasm合同迁移到新的代码版本 |
| [set-contract-admin](#iris-tx-wasm-set-contract-admin)    | 为合约设置新管理员            |
| [clear-contract-admin](#iris-tx-wasm-clear-contract-admin)       | 清除管理员的合约，以防止进一步的迁移                  |
| [code](#iris-query-wasm-code)        | 下载wasm字节码以获得给定的代码ID       |
| [contract](#iris-query-wasm-contract)  | 根据合约的地址打印出合约的元数据     |
| [contract-history](#iris-query-wasm-contract-history)  | 根据合约的地址打印出合约的代码历史记录     |
| [list-code](#iris-query-wasm-list-code)  | 列出链上的所有wasm字节码     |
| [list-contract-by-code](#iris-query-wasm-list-contract-by-code)  | 列出wasm链上所有给定代码ID的字节码    |
| [contract-state](#iris-query-wasm-contract-state)  | 查询wasm模块的命令     |

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


## iris tx wasm instantiate

实例化wasm合约

```bash
iris tx wasm instantiate [code_id_int64] [json_encoded_init_args] --label [text] --admin [address,optional] --amount [coins,optional] [flags]
```

**标识：**

| 名称，速记       | 类型    | 必须 | 默认          | 描述                                                               |
| -------------------| ------- | -------- | ------------- | ------------------------------------------------------------------------------------------------------------------------------ |
| --label            | string  |          |               | 列表中此合约的易读名称                                  |
| --admin            | string  |          |               | admin地址               |
| --amount           | string  |          |               | 在实例化期间发送给合约的币数 |


### 初始化智能合约

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

## iris tx wasm clear-contract-admin

清除管理员的合约，以防止进一步的迁移

```bash
iris tx wasm clear-contract-admin [contract_addr_bech32] [flags]
```

### 清空合约管理人权限

```bash
iris tx wasm clear-contract-admin iaa18vd8fpwxzck93qlwghaj6arh4p7c5n89fqcgm9 
    --from node0  
    --chain-id=test 
    --keyring-backend=file 
    --home ./cschain/node0/cschaincli 
    --fees 6point --gas="auto" -b block

```

## iris query wasm code

下载智能合约的二进制数据(上传的.wasm数据)

```bash
iris query wasm code [code_id] [output filename] [flags]
```

### 下载智能合约

```bash
iris query wasm code 1 election.wasm
```

## iris query wasm contract

查询合约的信息，包括合约地址、codeID，label等信息。

```bash
iris query wasm contract [bech32_address] [flags]
```

### 查询合约信息

```bash
iris query wasm contract iaa18vd8fpwxzck93qlwghaj6arh4p7c5n89fqcgm9
```

## iris query wasm contract-history

使用合约地址查询合约的migrate历史信息。

```bash
iris query wasm contract-history [bech32_address] [flags]
```

### 查询migrate历史信息
```bash
iris query wasm contract-history iaa18vd8fpwxzck93qlwghaj6arh4p7c5n89fqcgm9
```

## iris query wasm list-code

查询链上所有的合约code_id。

```bash
iris query wasm list-code [flags]
```

## iris query wasm list-contract-by-code

使用code_id查询合约的基本信息。

```bash
iris query wasm list-contract-by-code [code_id] [flags]
```

## iris query wasm contract-state

根据合约地址查询合约内部存储的状态信息

```bash
iris query wasm contract-state [command]
iris query wasm contract-state [flags]
```
### 查询所有状态信息

```bash
iris query wasm contract-state all iaa18vd8fpwxzck93qlwghaj6arh4p7c5n89fqcgm9
```

### 根据用户指定的合约地址以及合约编写时指定的key，查询当前合约的状态信息

```bash
iris query wasm contract-state raw iaa18vd8fpwxzck93qlwghaj6arh4p7c5n89fqcgm9 0006636F6E666967
```

