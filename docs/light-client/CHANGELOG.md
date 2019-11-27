# Change Log

## v0.16.0

*Nov 20th, 2019*

### BREAKING CHANGES

- Remove `initial_supply` property in the `TokenAddition` proposal API

### NON-BREAKING CHANGES

- Add Coinswap module APIs
- Add HTLC module APIs

## v0.15.0

*Aug 20th, 2019*

### BREAKING CHANGES

In this version, all POST methods (except '/tx/broadcast') just generate unsigned transactions, but don't broadcast them. Since '/tx/sign' is removed, users are required to sign the tx locally and use '/tx/broadcast' to broadcast the signed tx.

- Remove Key management APIs
- Remove POST /tx/sign
- Remove GET /distribution/community-tax
- Remove GET /gov/params/{module}

### NON-BREAKING CHANGES

- Add Asset module APIs
- Add Rand module APIs
- Add Params module APIs
- Add GET /bank/token-stats/{id}
- Add POST /bank/accounts/{address}/set-memo-regexp

#### Bank module APIs

| [v0.14.1]                    | [v0.15.0]                    | input changed | output changed | notes                                                        |
| ---------------------------- | ---------------------------- | ------------- | -------------- | ------------------------------------------------------------ |
| GET /bank/accounts/{address} | GET /bank/accounts/{address} | No            | Yes            | 1. Add `memo_regexp` in output; <br> 2. Tokens other than iris-atto could show up in output when people start using the newly introduced asset functionality. |

#### Tendermint module APIs

| [v0.14.1]   | [v0.15.0]   | input changed | output changed | notes                     |
| ----------- | ----------- | ------------- | -------------- | ------------------------- |
| /txs/{hash} | /txs/{hash} | No            | Yes            | Add `timestamp` in output |
| /txs        | /txs        | No            | Yes            | Add `timestamp` in output |

## v0.14.1

*May 31th, 2019*

### BREAKING CHANGES

- Change some api interfaces of `IRISLCD`

#### Bank module APIs

| [v0.14.0]                    | [v0.14.1]                    | input changed | output changed |
| ---------------------------- | ---------------------------- | ------------- | -------------- |
| GET /bank/accounts/{address} | GET /bank/accounts/{address} | No            | Yes            |

## v0.14.0

*May 20th, 2019*

### BREAKING CHANGES

- Change/drop some api interfaces of `IRISLCD`
- Add bank module APIs
- Add `community-tax` and `rewards` APIs in distribution module

#### Tendermint APIs

| [v0.13.1]                   | [v0.14.0]                   | input changed | output changed |
| --------------------------- | --------------------------- | ------------- | -------------- |
| GET /node_info              | GET /node-info              | No            | No             |
| GET /blocks-result/latest   | GET /block-results/latest   | No            | No             |
| GET /blocks-result/{height} | GET /block-results/{height} | No            | No             |
| POST /txs                   | N/A                         | /             | /              |

#### Key management APIs

| [v0.13.1]                    | [v0.14.0] | input changed | output changed |
| ---------------------------- | --------- | ------------- | -------------- |
| GET /auth/accounts/{address} | N/A       | /             | /              |

#### Sign and broadcast transactions

| [v0.13.1]      | [v0.14.0] | input changed | output changed |
| -------------- | --------- | ------------- | -------------- |
| POST /txs/send | N/A       | /             | /              |

#### Bank module APIs

| [v0.13.1]                               | [v0.14.0]                          | input changed | output changed |
| --------------------------------------- | ---------------------------------- | ------------- | -------------- |
| GET /bank/coin/{coin-type}              | GET /bank/coins/{type}             | No            | No             |
| GET /bank/token-stats                   | GET /bank/token-stats              | No            | Yes            |
| GET /bank/balances/{address}            | GET /bank/accounts/{address}       | No            | Yes            |
| POST /bank/accounts/{address}/transfers | POST /bank/accounts/{address}/send | Yes           | No             |
| POST /bank/burn                         | POST /bank/accounts/{address}/burn | Yes           | No             |

#### Stake module APIs

| [v0.13.1]                                                    | [v0.14.0]                                                    | input changed | output changed |
| ------------------------------------------------------------ | ------------------------------------------------------------ | ------------- | -------------- |
| POST /stake/delegators/{delegatorAddr}/delegate              | POST /stake/delegators/{delegatorAddr}/delegations           | No            | No             |
| POST /stake/delegators/{delegatorAddr}/redelegate            | POST /stake/delegators/{delegatorAddr}/redelegations         | No            | No             |
| POST /stake/delegators/{delegatorAddr}/unbond                | POST /stake/delegators/{delegatorAddr}/unbonding-delegations | No            | No             |
| GET /stake/delegators/{delegatorAddr}/unbonding_delegations  | GET /stake/delegators/{delegatorAddr}/unbonding-delegations  | No            | No             |
| GET /stake/delegators/{delegatorAddr}/unbonding_delegations/{validatorAddr} | GET /stake/delegators/{delegatorAddr}/unbonding-delegations/{validatorAddr} | No            | No             |
| GET /stake/validators/{validatorAddr}/unbonding_delegations  | GET /stake/validators/{validatorAddr}/unbonding-delegations  | No            | No             |

#### Slash module APIs

| [v0.13.1]                                               | [v0.14.0]                                               | input changed | output changed |
| ------------------------------------------------------- | ------------------------------------------------------- | ------------- | -------------- |
| GET /slashing/validators/{validatorPubKey}/signing_info | GET /slashing/validators/{validatorPubKey}/signing-info | No            | No             |

#### Distribution module APIs

| [v0.13.1]                                                   | [v0.14.0]                                           | input changed | output changed |
| ----------------------------------------------------------- | --------------------------------------------------- | ------------- | -------------- |
| POST /distribution/{delegatorAddr}/withdrawAddress          | POST /distribution/{delegatorAddr}/withdraw-address | No            | No             |
| GET /distribution/{delegatorAddr}/withdrawAddress           | GET /distribution/{delegatorAddr}/withdraw-address  | No            | No             |
| POST /distribution/{delegatorAddr}/withdrawReward           | POST /distribution/{delegatorAddr}/rewards/withdraw | No            | No             |
| GET /distribution/{delegatorAddr}/distrInfo/{validatorAddr} | N/A                                                 | /             | /              |
| GET /distribution/{delegatorAddr}/distrInfos                | N/A                                                 | /             | /              |
| GET /distribution/{validatorAddr}/valDistrInfo              | N/A                                                 | /             | /              |
| N/A                                                         | GET /distribution/{address}/rewards                 | /             | /              |
| N/A                                                         | GET /distribution/community-tax                     | /             | /              |

#### Service module APIs

| [v0.13.1]                                                    | [v0.14.0]                                                    | input changed | output changed |
| ------------------------------------------------------------ | ------------------------------------------------------------ | ------------- | -------------- |
| POST /service/definition                                     | POST /service/definitions                                    | No            | No             |
| GET /service/definition/{defChainId}/{serviceName}           | GET /service/definitions/{defChainId}/{serviceName}          | No            | No             |
| POST /service/binding                                        | POST /service/bindings                                       | No            | No             |
| GET /service/binding/{defChainId}/{serviceName}/{bindChainId}/{provider} | GET /service/bindings/{defChainId}/{serviceName}/{bindChainId}/{provider} | No            | No             |
| PUT /service/binding/{defChainId}/{serviceName}/{provider}   | PUT /service/bindings/{defChainId}/{serviceName}/{provider}  | No            | No             |
| PUT /service/binding/{defChainId}/{serviceName}/{provider}/disable | PUT /service/bindings/{defChainId}/{serviceName}/{provider}/disable | No            | No             |
| PUT /service/binding/{defChainId}/{serviceName}/{provider}/enable | PUT /service/bindings/{defChainId}/{serviceName}/{provider}/enable | No            | No             |
| PUT /service/binding/{defChainId}/{serviceName}/{provider}/deposit/refund | PUT /service/bindings/{defChainId}/{serviceName}/{provider}/deposit/refund | No            | No             |
| POST /service/request                                        | POST /service/requests                                       | No            | No             |
| POST /service/response                                       | POST /service/responses                                      | No            | No             |
| GET /service/response/{reqChainId}/{reqId}                   | GET /service/responses/{reqChainId}/{reqId}                  | No            | No             |

#### Query app version

| [v0.13.1]         | [v0.14.0]         | input changed | output changed |
| ----------------- | ----------------- | ------------- | -------------- |
| GET /node_version | GET /node-version | No            | No             |
