# Nft

`NFT` provides the ability to digitize assets. Through this module, each off-chain asset will be modeled as a unique on-chain nft.

## Available Commands

| Name                                   | Description                 |
| -------------------------------------- | --------------------------- |
| [issue](#iris-tx-nft-issue)        | Specify the nft Denom (nft category) and metadata JSON Schema to issue nft.    |
| [mint](#iris-tx-nft-mint)          | Additional issuance (create) of specific nft of this type can be made.    |
| [edit](#iris-tx-nft-edit)          | The metadata of the specified nft can be updated.   |
| [transfer](#iris-tx-nft-transfer)  | Transfer designated nft.   |
| [burn](#iris-tx-nft-burn)          | Destroy the created nft.   | 
| [supply](#iris-query-nft-supply)   | Query the total amount of nft according to Denom; accept the optional owner parameter. | 
| [owner](#iris-query-nft-owner)    | Query all nft owned by an account; you can specify the Denom parameter. | 
| [collection](#iris-query-nft-collection)   | Query all nft according to Denom. | 
| [denom](#iris-query-nft-denom)   | Query nft denom information based on Denom. | 
| [denoms](#iris-query-nft-denoms)   | Query the total amount of nft according to Denom; accept the optional owner parameter. | 
| [token](#iris-query-nft-token)   | Query specific nft based on Denom and ID. |                                    

## iris tx nft issue

Specify the nft Denom (nft category) and metadata JSON Schema to issue nft.   

```bash
iris tx nft issue [denom] [flags]
```

**Flags:**

| Name, shorthand      | Required | Default | Description     |
| -------------------- | --------- | -------------------------------------------------- | ---- |
| --name               |           | The name of the denom                |          |
| --schema             |         | Denom data structure definition         |          |


## iris tx nft mint

Additional issuance (create) of specific nft of this type can be made.  

```bash
iris tx nft mint [denomID] [tokenID] [flags]
```

**Flags:**

| Name, shorthand      | Required | Default | Description     |
| ----------- | ----- | -------------------------------------- | ---- |
| --uri       |  |  URI of off-line token data              |      |
| --recipient  | |  Receiver of the nft   |  |
| --name       |  | The name of nft     |      |


## iris tx nft edit

The metadata of the specified nft can be updated. 

```bash
iris tx nft edit [denomID] [tokenID] [flags]
```

**Flags:**

| Name, shorthand      | Required | Default | Description     |
| ----------- | ----- | -------------------------------------- | ---- |
| --uri       |  |   URI of off-line token data                   |      |
| --name       |  | The name of nft       |      |


## iris tx nft transfer 

Transfer designated nft. 

```bash
iris tx nft transfer [recipient] [denomID] [tokenID] [flags]
```

**Flags:**

| Name, shorthand      | Required | Default | Description     |
| ----------- | ----- | -------------------------------------- | ---- |
| --uri       |  |  URI of off-line token data                |      |
| --name       |  | The name of nft       |      |


## iris tx nft burn 

Destroy the created nft. 

```bash
iris tx nft burn [denomID] [tokenID] [flags]
```


## iris query nft 

Query nft 

### iris query nft supply 

```bash
iris query nft supply [denomID] 
iris query nft supply [denomID] --owner=<owner address>
```

### iris query nft owner 

```bash
iris query nft owner [owner address]
iris query nft owner --denom=<denomID> [owner address]
```

### iris query nft collection 

```bash
iris query nft collection [denomID] 
```

### iris query nft denom 

```bash
iris query nft denom [denomID]
```

### iris query nft denoms 

```bash
iris query nft denoms 
```

### iris query nft token 

```bash
iris query nft token [denomID] [tokenID]
```