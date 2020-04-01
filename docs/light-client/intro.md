---
order: 1
---

# Introduction

The IRIShub API Server is also called LCD(Light Client Daemon). An IRISLCD instance is a light node of IRIShub. Unlike IRIShub full node, it won't store all blocks and execute all transactions, which means it only requires minimal bandwidth, computing and storage resource. In distrust mode, it will track the evolution of validator set change and require full nodes to return consensus proof and merkle proof. Unless validators with more than 2/3 voting power do byzantine behavior, then IRISLCD proof verification algorithm can detect all potential malicious data, which means an IRISLCD instance can provide the same security as full nodes.

The default home folder of irislcd is `$HOME/.irislcd`. Once an IRISLCD is started, it will create two directories: `keys` and `trust-base.db`.The keys store db locates in `keys`. `trust-base.db` stores all trusted validator set and other verification related files.

When IRISLCD is started in distrust mode, it will check whether `trust-base.db` is empty. If true, it will fetch the latest block as its trust basis and save it under `trust-base.db`. The IRISLCD instance always trust the basis. All query proof will be verified based on the trust basis, for detailed proof verification algorithm please refer to [tendermint lite](https://github.com/tendermint/tendermint/blob/master/docs/tendermint-core/light-client-protocol.md).

## Basic Functionalities

- Provide restful APIs and swagger-ui to show these APIs
- Verify query proof

## Start

IRISLCD has two subcommands:

| Subcommand | Description               |
| ---------- | ------------------------- |
| version    | Print the IRISLCD version |
| start      | Start an IRISLCD instance |

The `start` subcommand has the following flags:

| Flag       | Type   | Default                 | Required | Description                                                     |
| ---------- | ------ | ----------------------- | -------- | --------------------------------------------------------------- |
| chain-id   | string |                         | Yes      | Chain ID of Tendermint node                                     |
| home       | string | "$HOME/.irislcd"        |          | Directory for config and data, such as key and checkpoint       |
| node       | string | "tcp://localhost:26657" |          | Full node to connect to                                         |
| laddr      | string | "tcp://localhost:1317"  |          | Address for server to listen on                                 |
| trust-node | bool   | false                   |          | Trust connected  full nodes (Don't verify proofs for responses) |
| max-open   | int    | 1000                    |          | The number of maximum open connections                          |
| cors       | string |                         |          | Set the domains that can make CORS requests                     |

By default, IRISLCD doesn't trust the connected full node. But if you are sure about that the connected full node is trustable, then you should run IRISLCD with `--trust-node` flag:

```bash
irislcd start --node=tcp://localhost:26657 --chain-id=irishub --trust-node
```

To access your IRISLCD instance publicly, you need to specify `--laddr`:

```bash
irislcd start --node=tcp://localhost:26657 --chain-id=irishub --laddr=tcp://0.0.0.0:1317 --trust-node
```

## REST APIs

Once IRISLCD is started, you can open <http://localhost:1317/swagger-ui/> in your browser and all available restful APIs will be shown. The `swagger-ui` page has detailed description about the APIs' functionality and required parameters. Here we just list all APIs and introduce their functionality briefly.

:::tip
**NOTE**

The `POST` apis(except [/tx/broadcast](#broadcast-transactions)) can only be used to generate unsigned transactions, you'll need to sign them in other ways before [broadcasting](#broadcast-transactions).
:::

### Tendermint APIs

such as query blocks, transactions and validator set

1. `GET /node-info`: The properties of the connected node
2. `GET /syncing`: Syncing state of node
3. `GET /blocks/latest`: Get the latest block
4. `GET /blocks/{height}`: Get a block at a certain height
5. `GET /block-results/latest`: Get the latest block result
6. `GET /block-results/{height}`: Get a block result at a certain height
7. `GET /validatorsets/latest`: Get the latest validator set
8. `GET /validatorsets/{height}`: Get a validator set at a certain height
9. `GET /txs/{hash}`: Get a Tx by hash
10. `GET /txs`: Search transactions

### Broadcast transactions API

1. `POST /tx/broadcast`: Broadcast a signed StdTx which is amino or json encoded

This api supports the following special parameters. By default, their values are all false. And each parameter has its unique priority( Here `0` is the top priority). If multiple parameters are specified to true, then the parameters with lower priority will be ignored. For instance, if `simulate` is true, then `commit` and `async` will be ignored.  

| parameter name | Type | Default | Priority | Description                                                                            |
| -------------- | ---- | ------- | -------- | -------------------------------------------------------------------------------------- |
| simulate       | bool | false   | 0        | Ignore the gas field and perform a simulation of a transaction, but don’t broadcast it |
| commit         | bool | false   | 1        | Wait for transaction being included in a block                                         |
| async          | bool | false   | 2        | Broadcast transaction asynchronously                                                   |

### Bank module APIs

1. `GET /bank/coins/{type}`: Query coin type
2. `GET /bank/token-stats`: Query token statistic
3. `GET /bank/token-stats/{symbol}`: Query the specified token statistic
4. `GET /bank/accounts/{address}`: Query the account information on blockchain
5. `POST /bank/accounts/{address}/send`: Send coins (build -> sign -> send)
6. `POST /bank/accounts/{address}/burn`: Burn coins

### Stake module APIs

1. `POST /stake/delegators/{delegatorAddr}/delegations`: Submit delegation transaction
2. `POST /stake/delegators/{delegatorAddr}/redelegations`: Submit redelegation transaction
3. `POST /stake/delegators/{delegatorAddr}/unbonding-delegations`: Submit unbonding transaction
4. `GET /stake/delegators/{delegatorAddr}/delegations`: Get all delegations from a delegator
5. `GET /stake/delegators/{delegatorAddr}/unbonding-delegations`: Get all unbonding delegations from a delegator
6. `GET /stake/delegators/{delegatorAddr}/redelegations`: Get all redelegations from a delegator
7. `GET /stake/delegators/{delegatorAddr}/validators`: Query all validators that a delegator is bonded to
8. `GET /stake/delegators/{delegatorAddr}/validators/{validatorAddr}`: Query a validator that a delegator is bonded to
9. `GET /stake/delegators/{delegatorAddr}/txs` :Get all staking txs from a delegator
10. `GET /stake/delegators/{delegatorAddr}/delegations/{validatorAddr}`: Query the current delegation between a delegator and a validator
11. `GET /stake/delegators/{delegatorAddr}/unbonding-delegations/{validatorAddr}`: Query all unbonding delegations between a delegator and a validator
12. `GET /stake/validators`: Get all validator candidates
13. `GET /stake/validators/{validatorAddr}`: Query the information from a single validator
14. `GET /stake/validators/{validatorAddr}/delegations`:  Get all delegations from a validator
15. `GET /stake/validators/{validatorAddr}/unbonding-delegations`: Get all unbonding delegations from a validator
16. `GET /stake/validators/{validatorAddr}/redelegations`: Get all outgoing redelegations from a validator
17. `GET /stake/pool`: Get the current state of the staking pool
18. `GET /stake/parameters`: Get the current staking parameter values

### Slashing module APIs

1. `GET /slashing/validators/{validatorPubKey}/signing-info`: Get sign info of given validator
2. `POST /slashing/validators/{validatorAddr}/unjail`: Unjail a jailed validator

### Distribution module APIs

1. `POST /distribution/{delegatorAddr}/withdraw-address`: Set withdraw address
2. `GET /distribution/{delegatorAddr}/withdraw-address`: Query withdraw address
3. `POST /distribution/{delegatorAddr}/rewards/withdraw`: Withdraw reward
4. `GET /distribution/{address}/rewards`: Query rewards
5. `GET /distribution/community-tax`: Query community tax

### Governance module APIs

1. `POST /gov/proposals`: Submit a proposal
2. `GET /gov/proposals`: Query proposals
3. `POST /gov/proposals/{proposalId}/deposits`: Deposit tokens to a proposal
4. `GET /gov/proposals/{proposalId}/deposits`: Query deposits
5. `POST /gov/proposals/{proposalId}/votes`: Vote a proposal
6. `GET /gov/proposals/{proposalId}/votes`: Query voters
7. `GET /gov/proposals/{proposalId}`: Query a proposal
8. `GET /gov/proposals/{proposalId}/deposits/{depositor}`: Query deposit
9. `GET /gov/proposals/{proposalId}/votes/{voter}`: Query vote

### Asset module APIs

1. `POST /asset/tokens`: Issue a token
2. `PUT /asset/tokens/{symbol}`: Edit an existing token
3. `POST /asset/tokens/{symbol}/mint`: Mint tokens to a specified address
4. `POST /asset/tokens/{symbol}/transfer`: Transfer the owner of a token to a new owner
5. `GET /asset/tokens/{symbol}`: Query a token by symbol
6. `GET /asset/tokens`: Query tokens by owner
7. `GET /asset/tokens/{symbol}/fee`: Query the fees for issuing and minting the specified token

### Coinswap module APIs

1. `POST /coinswap/liquidities/{voucher-coin-name}/deposit`: add liquidities
2. `POST /coinswap/liquidities/{voucher-coin-name}/withdraw`: withdraw liquidities
3. `POST /coinswap/liquidities/buy`: swap token(buy a fixed number of tokens)
4. `POST /coinswap/liquidities/sell`: swap token(sell a fixed number of tokens)
5. `GET /coinswap/liquidities/{voucher-coin-name}`: query liquidity by the voucher coin name

### HTLC module APIs

1. `POST /htlc/htlcs`: Create an HTLC
2. `GET /htlc/htlcs/{hash-lock}`: Query an HTLC by hash-lock
3. `POST /htlc/htlcs/{hash-lock}/claim`: Claim tokens locked in an OPEN HTLC to the recipient address
4. `POST /htlc/htlcs/{hash-lock}/refund`: Refund from an expired HTLC

### Service module APIs

1. `POST /service/definitions`: Define a new service
2. `GET /service/definitions/{service-name}`: Query a service definition
3. `POST /service/bindings`: Bind a service
4. `GET /service/bindings/{service-name}/{provider}`: Query a service binding
5. `GET /service/bindings{service-name}`: Query all bindings of a service definition
6. `POST /service/providers/{provider}/withdraw-address`: Set a withdrawal address for the provider
7. `GET /service/providers/{provider}/withdraw-address`: Query the withdrawal address of a provider
8. `PUT /service/bindings/{service-name}/{provider}`: Update a service binding
9. `POST /service/bindings/{service-name}/{provider}/disable`: Disable an available service binding
10. `POST /service/bindings/{service-name}/{provider}/enable`: Enable an unavailable service binding
11. `POST /service/bindings/{service-name}/{provider}/refund-deposit`: Refund all deposit from a service binding
12. `POST /service/contexts`: Initiate a service call
13. `GET /service/requests/{request-id}`: Query a request by the request ID
14. `GET /service/requests/{service-name}/{provider}`: Query active requests of a service binding
15. `GET /service/requests/{request-context-id}/{batch-counter}`: Query active requests by the request context ID and batch counter
16. `POST /service/responses`: Respond to a service request
17. `GET /service/responses/{request-id}`: Query a response by the request ID
18. `GET /service/responses/{request-context-id}/{batch-counter}`: Query active responses by the request context ID and batch counter
19. `GET /service/contexts/{request-context-id}`: Query a request context
20. `PUT /service/contexts/{request-context-id}`: Update a request context
21. `POST /service/contexts/{request-context-id}/pause`: Pause a running request context
22. `POST /service/contexts/{request-context-id}/start`: Start a paused request context
23. `POST /service/contexts/{request-context-id}/kill`: Terminate a request context
24. `GET /service/fees/{provider}`: Query the earned fees of a provider
25. `POST /service/fees/{provider}/withdraw`: Withdraw the earned fees of a provider

### Oracle module APIs

1. `POST /oracle/feeds`: Define a new Feed with the initial state of paused.
2. `POST /oracle/feeds/<feed-name>/start`: Start a paused Feed created by the tx signer.
3. `POST /oracle/feeds/<feed-name>/pause`: Pause a running Feed created by the tx signer.
4. `PUT /oracle/feeds/<feed-name>`: Update a Feed definition created by the tx signer.
5. `GET /oracle/feeds/<feed-name>`: Query Feed information by its name.
6. `GET /oracle/feeds?state=<state>`: Query Feed list by Feed state.
7. `GET /oracle/feeds/<feed-name>/values`: Query the result of the feed by its name, in descending order of timestamp.


### Rand module APIs

1. `POST /rand/rands`: Request a randon number
2. `GET /rand/rands/{request-id}`: Query a random number by the specified request id
3. `GET /rand/queue`: Query the pending requests with an optional height

### Params module APIs

1. `GET /params`: Query system params

### Query app version

1. `GET /version`: Version of IRISLCD
2. `GET /node-version`: Version of the connected node
