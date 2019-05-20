# IRISLCD Change Log

## v0.14.0

*May 20th, 2019*

### BREAKING CHANGES:

- Change/drop some api interfaces of `IRISLCD`
- Add bank module APIs
- Add `community-tax` and `rewards` APIs in distribution module


#### Tendermint APIs

| [v0.13.1]      | [v0.14.0]        | 
| --------------- | --------------- |
| GET /node_info  | GET /node-info   | 
| GET /blocks-result/latest | GET /blocks-results/latest | 
| GET /blocks-result/{height}  | GET /blocks-results/{height}  | 
| POST /txs  | N/A | 

#### Key management APIs

| [v0.13.1]      | [v0.14.0]        | 
| --------------- | --------------- |
| GET /auth/accounts/{address}  |  N/A | 
    
#### Sign and broadcast transactions

| [v0.13.1]      | [v0.14.0]        | 
| --------------- | --------------- |
| POST /txs/send  |  N/A | 
| GET /bank/coin/{coin-type} |  N/A | 
| GET /bank/token-stats  |  N/A | 
| GET /bank/balances/{address}  |  N/A | 
| POST /bank/accounts/{address}/transfers  |  N/A | 
| POST /bank/burn  |  N/A | 

#### Bank module APIs

| [v0.13.1]      | [v0.14.0]        | 
| --------------- | --------------- |
| N/A | GET /bank/coins/{type} | 
| N/A | GET /bank/token-stats | 
| N/A | GET /bank/accounts/{address} | 
| N/A | POST /bank/accounts/{address}/send | 
| N/A | POST /bank/accounts/{address}/burn | 


#### Stake module APIs

| [v0.13.1]      | [v0.14.0]        | 
| --------------- | --------------- |
| POST /stake/delegators/{delegatorAddr}/delegate | POST /stake/delegators/{delegatorAddr}/delegations | 
| POST /stake/delegators/{delegatorAddr}/redelegate | POST /stake/delegators/{delegatorAddr}/redelegations | 
| POST /stake/delegators/{delegatorAddr}/unbond | POST /stake/delegators/{delegatorAddr}/unbonding-delegations | 
| GET /stake/delegators/{delegatorAddr}/unbonding_delegations | GET /stake/delegators/{delegatorAddr}/unbonding-delegations | 
| GET /stake/delegators/{delegatorAddr}/unbonding_delegations/{validatorAddr} | GET /stake/delegators/{delegatorAddr}/unbonding-delegations/{validatorAddr} | 
| GET /stake/validators/{validatorAddr}/unbonding_delegations | GET /stake/validators/{validatorAddr}/unbonding-delegations | 

#### Slash module APIs

| [v0.13.1]      | [v0.14.0]        | 
| --------------- | --------------- |
| GET /slashing/validators/{validatorPubKey}/signing_info | GET /slashing/validators/{validatorPubKey}/signing-info |
 
    
#### Distribution module APIs

| [v0.13.1]      | [v0.14.0]        | 
| --------------- | --------------- |
| POST /distribution/{delegatorAddr}/withdrawAddress | POST /distribution/{delegatorAddr}/withdraw-address |
| GET /distribution/{delegatorAddr}/withdrawAddress | GET /distribution/{delegatorAddr}/withdraw-address |
| POST /distribution/{delegatorAddr}/withdrawReward | POST /distribution/{delegatorAddr}/rewards/withdraw |
| GET /distribution/{delegatorAddr}/distrInfo/{validatorAddr} | N/A |
| GET /distribution/{delegatorAddr}/distrInfos | N/A |
| GET /distribution/{validatorAddr}/valDistrInfo | N/A |
| N/A | GET /distribution/{address}/rewards |
| N/A | GET /distribution/community-tax |

#### Service module APIs

| [v0.13.1]      | [v0.14.0]        | 
| --------------- | --------------- |
| POST /service/definition | POST /service/definitions |
| GET /service/definition/{defChainId}/{serviceName} | GET /service/definitions/{defChainId}/{serviceName} |
| POST /service/binding | POST /service/bindings |
| GET /service/binding/{defChainId}/{serviceName}/{bindChainId}/{provider} | GET /service/bindings/{defChainId}/{serviceName}/{bindChainId}/{provider} |
| PUT /service/binding/{defChainId}/{serviceName}/{provider} | PUT /service/bindings/{defChainId}/{serviceName}/{provider} |
| PUT /service/binding/{defChainId}/{serviceName}/{provider}/disable | PUT /service/bindings/{defChainId}/{serviceName}/{provider}/disable |
| PUT /service/binding/{defChainId}/{serviceName}/{provider}/enable | PUT /service/bindings/{defChainId}/{serviceName}/{provider}/enable |
| PUT /service/binding/{defChainId}/{serviceName}/{provider}/deposit/refund| PUT /service/bindings/{defChainId}/{serviceName}/{provider}/deposit/refund |
| POST /service/request | POST /service/requests |
| POST /service/response | POST /service/responses |
| GET /service/response/{reqChainId}/{reqId} | GET /service/responses/{reqChainId}/{reqId} |

#### Query app version

| [v0.13.1]      | [v0.14.0]        | 
| --------------- | --------------- |
| GET /node_version | GET /node-version |