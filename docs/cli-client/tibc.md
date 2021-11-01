# TIBC

`TIBC` todo
## Client

### Available Commands

| Name                                     | Description                                                                            |
| ---------------------------------------- | -------------------------------------------------------------------------------------- |
| [create](#iris-tx-gov-submit-proposal-client-create)              | Create a new TIBC client with the specified client state and consensus state.            |
| [update](#iris-tx-tibc-update)                |  Update an existing TIBC client with a header.                 |
| [upgrade](#iris-tx-gov-submit-proposal-client-upgrade)                | Upgrade a TIBC client with the specified client state and consensus state.                                      |
| [regesiter](#iris-tx-gov-submit-proposal-relayer-register)        | Submit a relayer register proposal for the specified client.                                 |
| [header](#iris-query-tibc-client-header)        | Query the latest Tendermint header of the running chain.                                 |
| [node-state](#iris-query-tibc-client-node-state)        | Query a node consensus state. This result is feed to the client creation transaction.                                 |
| [relayers](#iris-query-tibc-client-relayers)        | Query all the registered relayers of a client.                                 |
| [client-state](#iris-query-tibc-client-state)        | Query a client state.                                |
| [client-states](#iris-query-tibc-client-states)        | Query all available light clients.                              |
| [consensus-state](#iris-query-consensus-state)        | Query the consensus state for a particular light client at a given height.                                 |
| [consensus-states](#iris-query-consensus-states)        | Query all the consensus states of a client.                                 |


### iris tx gov submit-proposal client-create
Submit a client create proposal.

```bash 
iris tx gov submit-proposal client-create [chain-name] [path/to/client_state.json] [path/to/consensus_state.json] [flags]
```

**Flags:**

| Name, shorthand | Required | Default                         | Description |
| --------------- | -------- | ------------------------------- | ----------- |
| --title          |          |            |          Title of proposal   |
| --description        |          |  |          description of proposal   |
| --deposit        |          |  |           deposit of proposal  |

### iris tx tibc update

Update existing client with a header.

```bash
iris tx tibc update [chain-name] [path/to/header.json]
```

## iris tx gov submit-proposal client-upgrade

Submit a client upgrade proposal.

```bash
iris tx gov submit-proposal client-upgrade [chain-name] [path/to/client_state.json] [path/to/consensus_state.json] [flags]
```
**Flags:**

| Name, shorthand | Required | Default                         | Description |
| --------------- | -------- | ------------------------------- | ----------- |
| --title          |          |            |          Title of proposal   |
| --description        |          |  |          description of proposal   |
| --deposit        |          |  |           deposit of proposal  |

## iris tx gov submit-proposal relayer-register

Submit a relayer register proposal.

```bash
iris tx gov submit-proposal relayer-register [chain-name] [relayers-address] [flags]
```

**Flags:**

| Name, shorthand | Required | Default                         | Description |
| --------------- | -------- | ------------------------------- | ----------- |
| --title          |          |            |          Title of proposal   |
| --description        |          |  |          description of proposal   |
| --deposit        |          |  |           deposit of proposal  |


### iris query consensus state

Query the consensus state for a particular light client at a given height.

```bash
iris query tibc client consensus-state [chain-name] [{revision}-{height}] [flags]
```
**Flags:**

| Name, shorthand | Required | Default                     | Description |
| --------------- | -------- | --------------------------- | ----------- |
| --height          |   int       |  |       Use a specific height to query state at (this can error if the node is pruning state)      |
| --latest-height          |          |             |  Return latest stored consensus state, format: {revision}-{height}           |
| --node          |   string       |       tcp://localhost:26657      | Host:port to Tendermint RPC interface for this chain          |
| --prove          |          |         true    |  Show proofs for the query results       |

### iris query consensus states

Query all the consensus states of a client.
```bash
iris query tibc client consensus-states [chain-name] [flags]
```

**Flags:**

| Name, shorthand | Required | Default                     | Description |
| --------------- | -------- | --------------------------- | ----------- |
| --count-total          |          |  |       Count total number of records in consensus states to query for      |
| --height          |          |             |  Return latest stored consensus state, format: {revision}-{height}           |

### iris query tibc client header
Query the latest Tendermint header of the running chain.
```bash
iris query tibc client header [flags]
```
| Name, shorthand | Required | Default                     | Description |
| --------------- | -------- | --------------------------- | ----------- |
| --node          |   string       |       tcp://localhost:26657      | Host:port to Tendermint RPC interface for this chain          |
| --height          |          |             |  Return latest stored consensus state, format: {revision}-{height}           |


### iris query tibc client node-state
Query a node consensus state. This result is feed to the client creation transaction.
```bash
iris query tibc client node-state [flags]
```

### iris query tibc client relayers
Query all the registered relayers of a client.
```bash
iris query tibc client relayers [chain-name] [flags]
```

### iris query tibc client state
Query a client state.
```bash
iris query tibc client state [chain-name] [flags]
```

### iris query tibc client states 
Query all available light clients.
```bash
iris query tibc client states [flags]
```

## Packet

### Available Commands

| Name                                     | Description                                                                            |
| ---------------------------------------- | -------------------------------------------------------------------------------------- |
| [send-clean-packet](#iris-tx-tibc-packet-send-clean-packet)              | Send a clean packet.            |
| [clean-packet-commitment](#iris-query-tibc-packet-clean-packet-commitment)                |  Query the clean packet commitment.                 |
| [packet-ack](#iris-query-tibc-packet-packet-ack)                | Query a packet acknowledgement.                                      |
| [packet-commitment](#iris-query-tibc-packet-packet-commitment)        | Query a packet commitment.                                                               |
| [packet-commitments](#iris-query-tibc-packet-packet-commitments)                | Query all packet commitments associated with source.                                      |
| [packet-receipt](#iris-query-tibc-packet-packet-receipt)        | Query a packet receipt.                                                              |
| [unreceived-acks](#iris-query-tibc-packet-unreceived-acks)        | Query all the unreceived packets associated with source chain name and destination chain name.                                                           |


### iris tx tibc packet send-clean-packet
Send a clean packet.
```bash
iris tx tibc packet send-clean-packet [dest-chain-name] [sequence] [flags]
```

**Flags:**

| Name, shorthand | Required | Default                     | Description |
| --------------- | -------- | --------------------------- | ----------- |
| --relay-chain-name          |   string       |            | The name of relay chain          |


### iris query tibc packet clean-packet-commitment
Query the clean packet commitment.
```bash
iris tx tibc packet send-clean-packet [dest-chain-name] [sequence] [flags]
```

### iris query tibc packet packet-ack
Query a packet acknowledgement.
```bash
iris query tibc packet packet-ack [source-chain] [dest-chain] [sequence] [flags]
```
### iris query tibc packet packet-commitment
Query a packet commitment.
```bash
iris query tibc packet packet-commitment [source-chain] [dest-chain] [sequence] [flags]
```

### iris query tibc packet packet-commitments
Query all packet commitments associated with source.
```bash
iris query tibc packet packet-commitments [source-chain] [dest-chain] [flags]
```

### iris query tibc packet packet-receipt
Query a packet receipt.
```bash
iris query tibc packet packet-receipt [source-chain] [dest-chain] [sequence] [flags]
```

### iris query tibc packet unreceived-acks
Query all the unreceived acks associated with source chain name and destination chain name.
```bash
iris query tibc packet unreceived-acks [source-chain] [dest-chain]] [flags]
```

**Flags:**

| Name, shorthand | Required | Default                     | Description |
| --------------- | -------- | --------------------------- | ----------- |
| --sequences          |  int64Slice        |            | comma separated list of packet sequence numbers (default [])          |


### iris query tibc packet unreceived-packets
Query all the unreceived packets associated with source chain name and destination chain name.
```bash
iris query tibc packet unreceived-packets [source-chain] [dest-chain] [flags]
```

**Flags:**

| Name, shorthand | Required | Default                     | Description |
| --------------- | -------- | --------------------------- | ----------- |
| --sequences          |  int64Slice        |            | comma separated list of packet sequence numbers (default [])          |


## Routing

### Available Commands

| Name                                     | Description                                                                            |
| ---------------------------------------- | -------------------------------------------------------------------------------------- |
| [set-rules](#iris-tx-gov-submit-proposal-set-rules)              | Submit a rules set proposal.            |
| [routing-rules](#iris-query-tibc-routing-routing-rules)                |  Query routing rules commitment..                 |

### iris tx gov submit-proposal set-rules
Submit a rules set proposal.

```bash
iris tx gov submit-proposal set-rules [path/to/routing_rules.json] [flags]
```
**Flags:**

| Name, shorthand | Required | Default                         | Description |
| --------------- | -------- | ------------------------------- | ----------- |
| --title          |          |            |          Title of proposal   |
| --description        |          |  |          description of proposal   |
| --deposit        |          |  |           deposit of proposal  |

### iris query tibc routing routing-rules
Query routing rules commitment.
```bash
iris query tibc routing routing-rules
```

## Nft-Transfer

### Available Commands

| Name                                     | Description                                                                            |
| ---------------------------------------- | -------------------------------------------------------------------------------------- |
| [nft-transfer](#iris-tx-tibc-nft-transfer-transfer)              | Transfer a non fungible token through TIBC.            |
| [class-trace](#iris-query-tibc-nft-transfer-class-trace)                |  Query the class trace info from a given trace hash.                |
| [class-traces](#iris-query-tibc-nft-transfer-class-traces)                |  Query the trace info for all nft classes.                |


### iris tx tibc-nft-transfer transfer
Transfer a non fungible token through TIBC.
```bash
iris tx tibc-nft-transfer transfer [dest-chain] [receiver] [class] [id] [flags]
```
**Flags:**

| Name, shorthand | Required | Default                     | Description |
| --------------- | -------- | --------------------------- | ----------- |
| --relay-chain          |  string        |            | relay chain used by cross-chain NFT         |
| --dest-contract          |  string        |            | The destination contract address to receive the nft         |


### iris query tibc-nft-transfer class-trace
Query the class trace info from a given trace hash.
```bash
iris query tibc-nft-transfer class-trace [hash]
```

### iris query tibc-nft-transfer class-traces
Query the trace info for all nft classes.
```bash
iris query tibc-nft-transfer class-traces
```