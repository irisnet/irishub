# TIBC

## Introduction

## Features

### Client

#### create client

Submit a client create proposal.

```bash 
iris tx gov submit-proposal client-create [chain-name] [path/to/client_state.json] [path/to/consensus_state.json] [flags]
```

#### update client

Update existing client with a header.

```bash
iris tx tibc update [chain-name] [path/to/header.json]
```

#### upgrade client

Submit a client upgrade proposal.

```bash
iris tx gov submit-proposal client-upgrade [chain-name] [path/to/client_state.json] [path/to/consensus_state.json] [flags]
```

#### register relayer

Submit a relayer register proposal.

```bash
iris tx gov submit-proposal relayer-register [chain-name] [relayers-address] [flags]
```

#### query consensus state

Query the consensus state for a particular light client at a given height.

```bash
iris query tibc client consensus-state [chain-name] [{revision}-{height}] [flags]
```

### Query the specified nft denom

Query nft denom information based on Denom ID.

`CLI`

```bash
iris q nft denom <denom-id>
```

### Query all nft denom information

Query all issued nft denom information.

`CLI`

```bash
iris q nft denoms
```

### Query the total amount of nft in a specified denom

Query the total amount of nft according to Denom ID; accept the optional owner parameter.

`CLI`

```bash
iris q nft supply <denom-id> --owner=<owner>
```

### Query all nft of the specified account

Query all nft owned by an account; you can specify the Denom ID parameter.

`CLI`

```bash
iris q nft owner <address> --denom-id=<denom-id>
```

### Query all nft of a specified denom

Query all nft according to Denom ID.

`CLI`

```bash
iris q nft collection <denom-id>
```

### Query specified nft

Query specific nft based on Denom ID and Token ID.

`CLI`

```bash
iris q nft token <denom-id> <token-id>
```
