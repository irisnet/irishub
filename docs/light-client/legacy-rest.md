# Legacy Amino JSON REST

The IRISHub versions v1.0.0 (depends on Cosmos-SDK v0.40) and earlier provided REST endpoints to query the state and broadcast transactions. These endpoints are kept in IRISHub v1.0, but they are marked as deprecated, and will be removed in v1.1 We therefore call these endpoints legacy REST endpoints.

Some important information concerning all legacy REST endpoints:

- Most of these endpoints are backwards-comptatible. All breaking changes are described in the next section.
- In particular, these endpoints still output Amino JSON. Cosmos-SDK v0.40 introduced Protobuf as the default encoding library throughout the codebase, but legacy REST endpoints are one of the few places where the encoding is hardcoded to Amino.

## API Port, Activation and Configuration

All routes are configured under the following fields in `~/.iris/config/app.toml`:

- `api.enable = true|false` field defines if the REST server should be enabled. Defaults to `true`.
- `api.address = {string}` field defines the address (really, the port, since the host should be kept at `0.0.0.0`) the server should bind to. Defaults to `tcp://0.0.0.0:1317`.
- some additional API configuration options are defined in `~/.iris/config/app.toml`, along with comments, please refer to that file directly.

## Legacy REST Endpoint

### Breaking Changes in Legacy REST Endpoints

| Legacy REST Endpoint                                                     | Description                                 | Breaking Change                                                                                                                                                                                                                                                                                                                                                                                                                                                                        |
| ------------------------------------------------------------------------ | ------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `POST /txs`                                                              | Broadcast tx                                | Endpoint will error when trying to broadcast transactions that don't support Amino serialization (e.g. IBC txs)<sup>1</sup>.                                                                                                                                                                                                                                                                                                                                                           |
| `POST /txs/encode`, `POST /txs/decode`                                   | Encode/decode Amino txs from JSON to binary | Endpoint will error when trying to encode/decode transactions that don't support Amino serialization (e.g. IBC txs)<sup>1</sup>.                                                                                                                                                                                                                                                                                                                                                       |
| `GET /txs/{hash}`                                                        | Query tx by hash                            | Endpoint will error when trying to output transactions that don't support Amino serialization (e.g. IBC txs)<sup>1</sup>.                                                                                                                                                                                                                                                                                                                                                              |
| `GET /txs`                                                               | Query tx by events                          | Endpoint will error when trying to output transactions that don't support Amino serialization (e.g. IBC txs)<sup>1</sup>.                                                                                                                                                                                                                                                                                                                                                              |
| `GET /gov/proposals/{id}/votes`, `GET /gov/proposals/{id}/votes/{voter}` | Gov endpoints for querying votes            | All gov endpoints which return votes return int32 in the `option` field instead of string: `1=VOTE_OPTION_YES, 2=VOTE_OPTION_ABSTAIN, 3=VOTE_OPTION_NO, 4=VOTE_OPTION_NO_WITH_VETO`.                                                                                                                                                                                                                                                                                                   |
| `GET /staking/*`                                                         | Staking query endpoints                     | All staking endpoints which return validators have two breaking changes. First, the validator's `consensus_pubkey` field returns an Amino-encoded struct representing an `Any` instead of a bech32-encoded string representing the pubkey. The `value` field of the `Any` is the pubkey's raw key as base64-encoded bytes. Second, the validator's `status` field now returns an int32 instead of string: `1=BOND_STATUS_UNBONDED`, `2=BOND_STATUS_UNBONDING`, `3=BOND_STATUS_BONDED`. |
| `GET /staking/validators`                                                | Get all validators                          | BondStatus is now a protobuf enum instead of an int32, and JSON serialized using its protobuf name, so expect query parameters like `?status=BOND_STATUS_{BONDED,UNBONDED,UNBONDING}` as opposed to `?status={bonded,unbonded,unbonding}`.                                                                                                                                                                                                                                             |

<sup>1</sup>: Transactions that don't support Amino serialization are the ones that contain one or more `Msg`s that are not registered with the Amino codec. Currently in the SDK, only IBC `Msg`s fall into this case.

### Migrating to New REST Endpoints

**IRISHub API Endpoints**

| Legacy REST Endpoint                                                            | Description                                                         | New gGPC-gateway REST Endpoint                                                                              |
| ------------------------------------------------------------------------------- | ------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------- |
| `GET /bank/balances/{address}`                                                  | Get the balance of an address                                       | `GET /cosmos/bank/v1beta1/balances/{address}`                                                               |
| `POST /bank/accounts/{address}/transfers`                                       | Send coins from one account to another                              | N/A, use Protobuf directly                                                                                  |
| `GET auth/accounts/{address}`                                                   | Get the account information on blockchain                           | `GET /cosmos/auth/v1beta1/accounts/{address}`                                                               |
| `GET /staking/delegators/{delegatorAddr}/delegations`                           | Get all delegations from a delegator                                | `GET /cosmos/staking/v1beta1/delegations/{delegator_addr}`                                                  |
| `POST /staking/delegators/{delegatorAddr}/delegations`                          | Submit delegation                                                   | N/A, use Protobuf directly                                                                                  |
| `GET /staking/delegators/{delegatorAddr}/delegations/{validatorAddr}`           | Query a delegation between a delegator and a validator              | `GET /cosmos/staking/v1beta1/validators/{validator_addr}/delegations/{delegator_addr}`                      |
| `GET /staking/delegators/{delegatorAddr}/unbonding_delegations`                 | Get all unbonding delegations from a delegator                      | `GET /cosmos/staking/v1beta1/delegators/{delegator_addr}/unbonding_delegations`                             |
| `POST /staking/delegators/{delegatorAddr}/unbonding_delegations`                | Submit an unbonding delegation                                      | N/A, use Protobuf directly                                                                                  |
| `GET /staking/delegators/{delegatorAddr}/unbonding_delegations/{validatorAddr}` | Query all unbonding delegations between a delegator and a validator | `GET /cosmos/staking/v1beta1/validators/{validator_addr}/delegations/{delegator_addr}/unbonding_delegation` |
| `GET /staking/redelegations`                                                    | Query redelegations                                                 | `GET /cosmos/staking/v1beta1/delegators/{delegator_addr}/redelegations`                                     |
| `POST /staking/delegators/{delegatorAddr}/redelegations`                        | Submit a redelegations                                              | N/A, use Protobuf directly                                                                                  |
| `GET /staking/delegators/{delegatorAddr}/validators`                            | Query all validators that a delegator is bonded to                  | `GET /cosmos/staking/v1beta1/delegators/{delegator_addr}/validators`                                        |
| `GET /staking/delegators/{delegatorAddr}/validators/{validatorAddr}`            | Query a validator that a delegator is bonded to                     | `GET /cosmos/staking/v1beta1/delegators/{delegator_addr}/validators/{validator_addr}`                       |
| `GET /staking/validators`                                                       | Get all validators                                                  | `GET /cosmos/staking/v1beta1/validators`                                                                    |
| `GET /staking/validators/{validatorAddr}`                                       | Get a single validator info                                         | `GET /cosmos/staking/v1beta1/validators/{validator_addr}`                                                   |
| `GET /staking/validators/{validatorAddr}/delegations`                           | Get all delegations to a validator                                  | `GET /cosmos/staking/v1beta1/validators/{validator_addr}/delegations`                                       |
| `GET /staking/validators/{validatorAddr}/unbonding_delegations`                 | Get all unbonding delegations from a validator                      | `GET /cosmos/staking/v1beta1/validators/{validator_addr}/unbonding_delegations`                             |
| `GET /staking/pool`                                                             | Get the current state of the staking pool                           | `GET /cosmos/staking/v1beta1/pool`                                                                          |
| `GET /staking/parameters`                                                       | Get the current staking parameter values                            | `GET /cosmos/staking/v1beta1/params`                                                                        |
| `GET /slashing/signing_infos`                                                   | Get all signing infos                                               | `GET /cosmos/slashing/v1beta1/signing_infos`                                                                |
| `POST /slashing/validators/{validatorAddr}/unjail`                              | Unjail a jailed validator                                           | N/A, use Protobuf directly                                                                                  |
| `GET /slashing/parameters`                                                      | Get slashing parameters                                             | `GET /cosmos/slashing/v1beta1/params`                                                                       |
| `POST /gov/proposals`                                                           | Submit a proposal                                                   | N/A, use Protobuf directly                                                                                  |
| `GET /gov/proposals`                                                            | Get all proposals                                                   | `GET /cosmos/gov/v1beta1/proposals`                                                                         |
| `POST /gov/proposals/param_change`                                              | Generate a parameter change proposal transactionl                   | N/A, use Protobuf directly                                                                                  |
| `GET /gov/proposals/{proposal-id}`                                              | Get proposal by id                                                  | `GET /cosmos/gov/v1beta1/proposals/{proposal_id}`                                                           |
| `GET /gov/proposals/{proposal-id}/proposer`                                     | Get proposer of a proposal                                          | `GET /cosmos/gov/v1beta1/proposals/{proposal_id}` (Get proposer from `Proposal` struct)                     |
| `GET /gov/proposals/{proposal-id}/deposits`                                     | Get deposits of a proposal                                          | `GET /cosmos/gov/v1beta1/proposals/{proposal_id}/deposits`                                                  |
| `POST /gov/proposals/{proposal-id}/deposits`                                    | Deposit tokens to a proposal                                        | N/A, use Protobuf directly                                                                                  |
| `GET /gov/proposals/{proposal-id}/deposits/{depositor}`                         | Get depositor a of deposit                                          | `GET /cosmos/gov/v1beta1/proposals/{proposal_id}/deposits/{depositor}`                                      |
| `GET /gov/proposals/{proposal-id}/votes`                                        | Get votes of a proposal                                             | `GET /cosmos/gov/v1beta1/proposals/{proposal_id}/votes`                                                     |
| `POST /gov/proposals/{proposal-id}/votes`                                       | Vote a proposal                                                     | N/A, use Protobuf directly                                                                                  |
| `GET /gov/proposals/{proposal-id}/votes/{voter}`                                | Get voted information by voterAddr.                                 | `GET /cosmos/gov/v1beta1/proposals/{proposal_id}/votes/{voter}`                                             |
| `GET /gov/proposals/{proposal-id}/tally`                                        | Get tally of a proposal                                             | `GET /cosmos/gov/v1beta1/proposals/{proposal_id}/tally`                                                     |
| `GET /gov/parameters/deposit`                                                   | Get governance deposit parameters                                   | `GET /cosmos/gov/v1beta1/params/{params_type}`                                                              |
| `GET /gov/parameters/tallying`                                                  | Query governance tally parameters                                   | `GET /cosmos/gov/v1beta1/params/{params_type}`                                                              |
| `GET /gov/parameters/voting`                                                    | Get governance voting parameters                                    | `GET /cosmos/gov/v1beta1/params/{params_type}`                                                              |
| `GET /distribution/delegators/{delegatorAddr}/rewards`                          | Get the total rewards balance from all delegations                  | `GET /cosmos/distribution/v1beta1/delegators/{delegator_address}/rewards`                                   |
| `POST /distribution/delegators/{delegatorAddr}/rewards`                         | Withdraw all delegator rewards                                      | N/A, use Protobuf directly                                                                                  |
| `GET /distribution/delegators/{delegatorAddr}/rewards/{validatorAddr}`          | Query a delegation reward                                           | `GET /cosmos/distribution/v1beta1/delegators/{delegator_address}/rewards/{validator_address}`               |
| `POST /distribution/delegators/{delegatorAddr}/rewards/{validatorAddr}`         | Withdraw a delegation reward                                        | N/A, use Protobuf directly                                                                                  |
| `GET /distribution/delegators/{delegatorAddr}/withdraw_address`                 | Get the rewards withdrawal address                                  | `GET /cosmos/distribution/v1beta1/delegators/{delegator_address}/withdraw_address`                          |
| `POST /distribution/delegators/{delegatorAddr}/withdraw_address`                | Replace the rewards withdrawal address                              | N/A, use Protobuf directly                                                                                  |
| `GET /distribution/validators/{validatorAddr}`                                  | Validator distribution information                                  | N/A, use Protobuf directly                                                                                  |
| `GET /distribution/validators/{validatorAddr}/outstanding_rewards`              | Outstanding rewards of a single validator                           | `GET /cosmos/distribution/v1beta1/validators/{validator_address}/outstanding_rewards`                       |
| `GET /distribution/validators/{validatorAddr}/rewards`                          | Commission and self-delegation rewards of a single a validator      | N/A, use Protobuf directly                                                                                  |
| `POST /distribution/validators/{validatorAddr}/rewards`                         | Withdraw the validator's rewards                                    | N/A, use Protobuf directly                                                                                  |
| `GET /distribution/community_pool`                                              | Get the amount held in the community pool                           | `GET /cosmos/distribution/v1beta1/community_pool`                                                           |
| `GET /distribution/parameters`                                                  | Get the current distribution parameter values                       | `GET /cosmos/distribution/v1beta1/params`                                                                   |

**Tendermint API Endpoints**

| Legacy REST Endpoint          | Description                                      | New gGPC-gateway REST Endpoint                               |
| ----------------------------- | ------------------------------------------------ | ------------------------------------------------------------ |
| `GET /node_info`              | Get the properties of the connected node         | `GET /cosmos/base/tendermint/v1beta1/node_info`              |
| `GET /syncing`                | Get syncing state of node                        | `GET /cosmos/base/tendermint/v1beta1/syncing`                |
| `GET /blocks/latest`          | Get the latest block                             | `GET /cosmos/base/tendermint/v1beta1/blocks/latest`          |
| `GET /blocks/{height}`        | Get a block at a certain height                  | `GET /cosmos/base/tendermint/v1beta1/blocks/{height}`        |
| `GET /validatorsets/latest`   | Get the latest validator set                     | `GET /cosmos/base/tendermint/v1beta1/validatorsets/latest`   |
| `GET /validatorsets/{height}` | Get a validator set a certain height             | `GET /cosmos/base/tendermint/v1beta1/validatorsets/{height}` |
| `GET /txs/{hash}`             | Query tx by hash                                 | `GET /cosmos/tx/v1beta1/txs/{hash}`                          |
| `GET /txs`                    | Query tx by events                               | `GET /cosmos/tx/v1beta1/txs`                                 |
| `POST /txs`                   | Broadcast tx                                     | `POST /cosmos/tx/v1beta1/txs`                                |
| `POST /txs/encode`            | Encodes an Amino JSON tx to an Amino binary tx   | N/A, use Protobuf directly                                   |
| `POST /txs/decode`            | Decodes an Amino binary tx into an Amino JSON tx | N/A, use Protobuf directly                                   |

## 高优先级查询端点的 breaking changes

- `GET /blocks/latest`&&`GET /blocks/{height}`

  - Specific changes:
    - The block_meta field is no longer used, the block_id in the original block_meta field structure is moved to the first level, and the header field is moved to the block
    - The num_txs and total_txs fields are no longer used, and the transaction number can be obtained by traversing the txs field

  - json example:

    ```json
    {
        "block_id": {
            "hash": "DC2EEC73C327BD338EE5667827A4EECA2A2A4752B38D5669CD17EDE07CFB6F30",
            "parts": {
                "total": 1,
                "hash": "1F794EA5185AE489A3D53FD6E19A690373CAB510B981B23E672668CBE0B668E5"
            }
        },
        "block": {
            "header": {
            "version": {
                "block": "11"
            },
            "chain_id": "irishub-1",
            "height": "5",
            "time": "2021-01-18T07:29:21.918234Z",
            "last_block_id": {
                "hash": "889428AAB1975F94C39F44F1BD9C94B2A46E0BC0EFF9AD625939DCB763E82D1F",
                "parts": {
                    "total": 1,
                    "hash": "A9FEDF247839148459E12CD0D8A495A9EE663F7C9E719D85133B27F1D810D52B"
                }
            },
            "last_commit_hash": "35084203D333AF835637F7D9F1FEAA03554AECAB9786EECC6EB5F236156A19F1",
            "data_hash": "A2853A6749C904A8C26C5DB8CB0DD731C44EEAF2AD6AEE2E633DF8F8FD0CA04F",
            "validators_hash": "79E608E448B5B9D0784FDE890506DDE025E0E079CE10B7E01687B9C0E2DFC124",
            "next_validators_hash": "79E608E448B5B9D0784FDE890506DDE025E0E079CE10B7E01687B9C0E2DFC124",
            "consensus_hash": "048091BC7DDC283F77BFBF91D73C44DA58C3DF8A9CBC867405D8B7F3DAADA22F",
            "app_hash": "8F668541D8D565B40373E1492ED6729674539FCB1705437E309522DD491E46DC",
            "last_results_hash": "E3B0C44298FC1C149AFBF4C8996FB92427AE41E4649B934CA495991B7852B855",
            "evidence_hash": "E3B0C44298FC1C149AFBF4C8996FB92427AE41E4649B934CA495991B7852B855",
            "proposer_address": "86C89798CEA6D07FB8550AFDD8DEEA0DA52BFEF4"
            },
            "data": {
                "txs": [
                    "CowBCokBChwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kEmkKKmlhYTE4YXduM2s3MHUwNXRsY3VsOHcycW5sNjRnMDAydWo0a2puOTNybhIqaWFhMXc5NzZhNWpyaHNqMDZkcW1yaDJ4OXF4emVsNzRxdGNtYXBrbHhjGg8KBXN0YWtlEgYxMDAwMDASZQpQCkYKHy9jb3Ntb3MuY3J5cHRvLnNlY3AyNTZrMS5QdWJLZXkSIwohA01sYgbsLpw+B9M+p6vyJCh1wfigTWbLpnhNfeDKxKIlEgQKAggBGAESEQoLCgVzdGFrZRICMTAQwJoMGkC+qpaBhJ20qboyU0HWBL0zVlW4klBXGZGsa8n2W1rxIVbq39DZUskGSI8WKNl1stM7QGhycu7YLU30z8vsg8N5"
                ]
            },
            "evidence": {
                "evidence": []
            },
            "last_commit": {
                "height": "4",
                "round": 0,
                "block_id": {
                    "hash": "889428AAB1975F94C39F44F1BD9C94B2A46E0BC0EFF9AD625939DCB763E82D1F",
                    "parts": {
                        "total": 1,
                        "hash": "A9FEDF247839148459E12CD0D8A495A9EE663F7C9E719D85133B27F1D810D52B"
                    }
                },
                "signatures": [
                    {
                        "block_id_flag": 2,
                        "validator_address": "86C89798CEA6D07FB8550AFDD8DEEA0DA52BFEF4",
                        "timestamp": "2021-01-18T07:29:21.918234Z",
                        "signature": "UGqkOIqSSzaW/2YzkzfoobcUIizzPcl9BVECl+jwhyMJSkrMnD3DdPUlS2Vd1IAU3u8qQOmP09+m/r5R0gp6CQ=="
                    }
                ]
            }
        }
    }
    ```

- `GET /block-results/latest`&&`GET /block-results/{height}`

  - Specific changes:
    - These two interfaces have been cancelled
    - Tendermint RPC's block_results interface can be directly called to obtain relevant data (key and value in the events field of the query result are all base64 encoded)

  - json example:

    ```json
    {
        "jsonrpc":"2.0",
        "id":-1,
        "result":{
            "height":"5",
            "txs_results":[
                {
                    "code":0,
                    "data":"CgYKBHNlbmQ=",
                    "log":"[{\"events\":[{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"send\"},{\"key\":\"sender\",\"value\":\"iaa18awn3k70u05tlcul8w2qnl64g002uj4kjn93rn\"},{\"key\":\"module\",\"value\":\"bank\"}]},{\"type\":\"transfer\",\"attributes\":[{\"key\":\"recipient\",\"value\":\"iaa1w976a5jrhsj06dqmrh2x9qxzel74qtcmapklxc\"},{\"key\":\"sender\",\"value\":\"iaa18awn3k70u05tlcul8w2qnl64g002uj4kjn93rn\"},{\"key\":\"amount\",\"value\":\"1000000uiris\"}]}]}]",
                    "info":"",
                    "gas_wanted":"200000",
                    "gas_used":"69256",
                    "events":[
                        {
                            "type":"transfer",
                            "attributes":[
                                {
                                    "key":"cmVjaXBpZW50",
                                    "value":"aWFhMTd4cGZ2YWttMmFtZzk2MnlsczZmODR6M2tlbGw4YzVsOW1yM2Z2",
                                    "index":true
                                },
                                {
                                    "key":"c2VuZGVy",
                                    "value":"aWFhMThhd24zazcwdTA1dGxjdWw4dzJxbmw2NGcwMDJ1ajRram45M3Ju",
                                    "index":true
                                },
                                {
                                    "key":"YW1vdW50",
                                    "value":"MTBzdGFrZQ==",
                                    "index":true
                                }
                            ]
                        },
                        {
                            "type":"message",
                            "attributes":[
                                {
                                    "key":"c2VuZGVy",
                                    "value":"aWFhMThhd24zazcwdTA1dGxjdWw4dzJxbmw2NGcwMDJ1ajRram45M3Ju",
                                    "index":true
                                }
                            ]
                        },
                        {
                            "type":"message",
                            "attributes":[
                                {
                                    "key":"YWN0aW9u",
                                    "value":"c2VuZA==",
                                    "index":true
                                }
                            ]
                        },
                        {
                            "type":"transfer",
                            "attributes":[
                                {
                                    "key":"cmVjaXBpZW50",
                                    "value":"aWFhMXc5NzZhNWpyaHNqMDZkcW1yaDJ4OXF4emVsNzRxdGNtYXBrbHhj",
                                    "index":true
                                },
                                {
                                    "key":"c2VuZGVy",
                                    "value":"aWFhMThhd24zazcwdTA1dGxjdWw4dzJxbmw2NGcwMDJ1ajRram45M3Ju",
                                    "index":true
                                },
                                {
                                    "key":"YW1vdW50",
                                    "value":"MTAwMDAwc3Rha2U=",
                                    "index":true
                                }
                            ]
                        },
                        {
                            "type":"message",
                            "attributes":[
                                {
                                    "key":"c2VuZGVy",
                                    "value":"aWFhMThhd24zazcwdTA1dGxjdWw4dzJxbmw2NGcwMDJ1ajRram45M3Ju",
                                    "index":true
                                }
                            ]
                        },
                        {
                            "type":"message",
                            "attributes":[
                                {
                                    "key":"bW9kdWxl",
                                    "value":"YmFuaw==",
                                    "index":true
                                }
                            ]
                        }
                    ],
                    "codespace":""
                }
            ],
            "begin_block_events":[
                {
                    "type":"transfer",
                    "attributes":[
                        {
                            "key":"cmVjaXBpZW50",
                            "value":"aWFhMTd4cGZ2YWttMmFtZzk2MnlsczZmODR6M2tlbGw4YzVsOW1yM2Z2",
                            "index":true
                        },
                        {
                            "key":"c2VuZGVy",
                            "value":"aWFhMW0zaDMwd2x2c2Y4bGxydXh0cHVrZHZzeTBrbTJrdW04YW44Zjkz",
                            "index":true
                        },
                        {
                            "key":"YW1vdW50",
                            "value":"MTI2NzUyMzVzdGFrZQ==",
                            "index":true
                        }
                    ]
                },
                {
                    "type":"message",
                    "attributes":[
                        {
                            "key":"c2VuZGVy",
                            "value":"aWFhMW0zaDMwd2x2c2Y4bGxydXh0cHVrZHZzeTBrbTJrdW04YW44Zjkz",
                            "index":true
                        }
                    ]
                },
                {
                    "type":"mint",
                    "attributes":[
                        {
                            "key":"bGFzdF9pbmZsYXRpb25fdGltZQ==",
                            "value":"MjAyMS0wMS0xOCAwNzoyOToxNi43NjIyMSArMDAwMCBVVEM=",
                            "index":true
                        },
                        {
                            "key":"aW5mbGF0aW9uX3RpbWU=",
                            "value":"MjAyMS0wMS0xOCAwNzoyOToyMS45MTgyMzQgKzAwMDAgVVRD",
                            "index":true
                        },
                        {
                            "key":"bWludF9jb2lu",
                            "value":"MTI2NzUyMzU=",
                            "index":true
                        }
                    ]
                },
                {
                    "type":"transfer",
                    "attributes":[
                        {
                            "key":"cmVjaXBpZW50",
                            "value":"aWFhMWp2NjVzM2dycWY2djZqbDNkcDR0NmM5dDlyazk5Y2Q4amF5ZHR3",
                            "index":true
                        },
                        {
                            "key":"c2VuZGVy",
                            "value":"aWFhMTd4cGZ2YWttMmFtZzk2MnlsczZmODR6M2tlbGw4YzVsOW1yM2Z2",
                            "index":true
                        },
                        {
                            "key":"YW1vdW50",
                            "value":"MTI2NzUyMzVzdGFrZQ==",
                            "index":true
                        }
                    ]
                },
                {
                    "type":"message",
                    "attributes":[
                        {
                            "key":"c2VuZGVy",
                            "value":"aWFhMTd4cGZ2YWttMmFtZzk2MnlsczZmODR6M2tlbGw4YzVsOW1yM2Z2",
                            "index":true
                        }
                    ]
                },
                {
                    "type":"proposer_reward",
                    "attributes":[
                        {
                            "key":"YW1vdW50",
                            "value":"NjMzNzYxLjc1MDAwMDAwMDAwMDAwMDAwMHN0YWtl",
                            "index":true
                        },
                        {
                            "key":"dmFsaWRhdG9y",
                            "value":"aXZhMThhd24zazcwdTA1dGxjdWw4dzJxbmw2NGcwMDJ1ajRrOHowNzc1",
                            "index":true
                        }
                    ]
                },
                {
                    "type":"commission",
                    "attributes":[
                        {
                            "key":"YW1vdW50",
                            "value":"NjMzNzYxLjc1MDAwMDAwMDAwMDAwMDAwMHN0YWtl",
                            "index":true
                        },
                        {
                            "key":"dmFsaWRhdG9y",
                            "value":"aXZhMThhd24zazcwdTA1dGxjdWw4dzJxbmw2NGcwMDJ1ajRrOHowNzc1",
                            "index":true
                        }
                    ]
                },
                {
                    "type":"rewards",
                    "attributes":[
                        {
                            "key":"YW1vdW50",
                            "value":"NjMzNzYxLjc1MDAwMDAwMDAwMDAwMDAwMHN0YWtl",
                            "index":true
                        },
                        {
                            "key":"dmFsaWRhdG9y",
                            "value":"aXZhMThhd24zazcwdTA1dGxjdWw4dzJxbmw2NGcwMDJ1ajRrOHowNzc1",
                            "index":true
                        }
                    ]
                },
                {
                    "type":"commission",
                    "attributes":[
                        {
                            "key":"YW1vdW50",
                            "value":"MTE3ODc5NjguNTUwMDAwMDAwMDAwMDAwMDAwc3Rha2U=",
                            "index":true
                        },
                        {
                            "key":"dmFsaWRhdG9y",
                            "value":"aXZhMThhd24zazcwdTA1dGxjdWw4dzJxbmw2NGcwMDJ1ajRrOHowNzc1",
                            "index":true
                        }
                    ]
                },
                {
                    "type":"rewards",
                    "attributes":[
                        {
                            "key":"YW1vdW50",
                            "value":"MTE3ODc5NjguNTUwMDAwMDAwMDAwMDAwMDAwc3Rha2U=",
                            "index":true
                        },
                        {
                            "key":"dmFsaWRhdG9y",
                            "value":"aXZhMThhd24zazcwdTA1dGxjdWw4dzJxbmw2NGcwMDJ1ajRrOHowNzc1",
                            "index":true
                        }
                    ]
                }
            ],
            "end_block_events":null,
            "validator_updates":null,
            "consensus_param_updates":{
                "block":{
                    "max_bytes":"22020096",
                    "max_gas":"-1"
                },
                "evidence":{
                    "max_age_num_blocks":"100000",
                    "max_age_duration":"172800000000000",
                    "max_bytes":"1048576"
                },
                "validator":{
                    "pub_key_types":[
                        "ed25519"
                    ]
                }
            }
        }
    }
    ```

- `GET /txs`&&`GET /txs/{hash}`

  - Specific changes:
    - Tags are no longer used; use the events field instead
    - The result field is no longer used, and the field in the original result is moved to the first level
    - The coin_flow field is no longer used

  - json examples:

    ```json
    {
        "height": "5",
        "txhash": "E663768B616B1ACD2912E47C36FEBC7DB0E0974D6DB3823D4C656E0EAB8C679D",
        "data": "0A060A0473656E64",
        "raw_log": "[{\"events\":[{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"send\"},{\"key\":\"sender\",\"value\":\"iaa18awn3k70u05tlcul8w2qnl64g002uj4kjn93rn\"},{\"key\":\"module\",\"value\":\"bank\"}]},{\"type\":\"transfer\",\"attributes\":[{\"key\":\"recipient\",\"value\":\"iaa1w976a5jrhsj06dqmrh2x9qxzel74qtcmapklxc\"},{\"key\":\"sender\",\"value\":\"iaa18awn3k70u05tlcul8w2qnl64g002uj4kjn93rn\"},{\"key\":\"amount\",\"value\":\"1000000uiris\"}]}]}]",
        "logs": [
            {
                "events": [
                    {
                        "type": "message",
                        "attributes": [
                            {
                                "key": "action",
                                "value": "send"
                            },
                            {
                                "key": "sender",
                                "value": "iaa18awn3k70u05tlcul8w2qnl64g002uj4kjn93rn"
                            },
                            {
                                "key": "module",
                                "value": "bank"
                            }
                        ]
                    },
                    {
                        "type": "transfer",
                        "attributes": [
                            {
                                "key": "recipient",
                                "value": "iaa1w976a5jrhsj06dqmrh2x9qxzel74qtcmapklxc"
                            },
                            {
                                "key": "sender",
                                "value": "iaa18awn3k70u05tlcul8w2qnl64g002uj4kjn93rn"
                            },
                            {
                                "key": "amount",
                                "value": "1000000uiris"
                            }
                        ]
                    }
                ]
            }
        ],
        "gas_wanted": "200000",
        "gas_used": "69256",
        "tx": {
            "type": "cosmos-sdk/StdTx",
            "value": {
                "msg": [
                    {
                        "type": "cosmos-sdk/MsgSend",
                        "value": {
                            "from_address": "iaa18awn3k70u05tlcul8w2qnl64g002uj4kjn93rn",
                            "to_address": "iaa1w976a5jrhsj06dqmrh2x9qxzel74qtcmapklxc",
                            "amount": [
                                {
                                    "denom": "uiris",
                                    "amount": "1000000"
                                }
                            ]
                        }
                    }
                ],
                "fee": {
                    "amount": [
                        {
                            "denom": "uiris",
                            "amount": "30000"
                        }
                    ],
                    "gas": "200000"
                },
                "signatures": [],
                "memo": "",
                "timeout_height": "0"
            }
        },
        "timestamp": "2021-01-18T07:29:21Z"
    }
    ```

- `GET /bank/accounts/{address}`

  - Specific changes:

    - This interface has been cancelled
    - The coins field in the original interface can be queried through the /bank/balances/{address} interface
    - Other fields in the original interface can be queried through the /auth/accounts/{address} interface

  - json example:
    - /bank/balances/{address}

        ```json
        {
            "height": "98",
            "result": [
                {
                    "denom": "node0token",
                    "amount": "1000000000"
                },
                {
                    "denom": "uiris",
                    "amount": "4999999999999899899990"
                }
            ]
        }
        ```

    - /auth/accounts/{address}

        ```json
        {
            "height": "142",
            "result": {
                "type": "cosmos-sdk/BaseAccount",
                "value": {
                    "address": "iaa18awn3k70u05tlcul8w2qnl64g002uj4kjn93rn",
                    "public_key": {
                        "type": "tendermint/PubKeySecp256k1",
                        "value": "A01sYgbsLpw+B9M+p6vyJCh1wfigTWbLpnhNfeDKxKIl"
                    },
                    "sequence": "2"
                }
            }
        }
        ```

## 构造未签名交易（完全向后兼容）

The same code as integrating with cosmoshub-3 mainnet. The transaction structure is as follows:

```json
{
    "type": "cosmos-sdk/StdTx",
    "value": {
        "msg": [
            {
                "type": "cosmos-sdk/MsgSend",
                "value": {
                    "from_address": "iaa1rkgdpj6fyyyu7pnhmc3v7gw9uls4mnajvzdwkt",
                    "to_address": "iaa1q6t5439f0rkvkzl38m0f43e0kpv3mx7x2shlq8",
                    "amount": [
                        {
                            "denom": "uiris",
                            "amount": "1000000"
                        }
                    ]
                }
            }
        ],
        "fee": {
            "amount": [
                {
                    "denom": "uiris",
                    "amount": "30000"
                }
            ],
            "gas": "200000"
        },
        "signatures": null,
        "memo": "Sent via irishub client"
    }
}
```

Where the IRISHub address prefix uses `iaa` instead, which affects the fields:

- value.msg.value.from_adress
- value.msg.value.to_address

Denom uses `uiris` instead (1iris = 10<sup>6</sup>uiris), which affects fields:

- value.msg.value.amount.denom
- value.fee.amount.denom

## Broadcasting a transaction（完全向后兼容）

The same code as integrating with cosmoshub-3 mainnet, call `POST /txs` to send a transaction, as the example below:

```bash
curl -X POST "http://localhost:1317/txs" -H "accept: application/json" -H "Content-Type: application/json" -d "{ \"tx\": {\"msg\":[{\"type\":\"cosmos-sdk/MsgSend\",\"value\":{\"from_address\":\"iaa1rkgdpj6fyyyu7pnhmc3v7gw9uls4mnajvzdwkt\",\"to_address\":\"iaa1q6t5439f0rkvkzl38m0f43e0kpv3mx7x2shlq8\",\"amount\":[{\"denom\":\"uiris\",\"amount\":\"1000000\"}]}}],\"fee\":{\"amount\":[{\"denom\":\"uiris\",\"amount\":\"30000\"}],\"gas\":\"200000\"},\"signatures\":[{\"pub_key\":{\"type\":\"tendermint/PubKeySecp256k1\",\"value\":\"AxGagdsRTKni/h1+vCFzTpNltwoiU7SwIR2dg6Jl5a//\"},\"signature\":\"Pu8yiRVO8oB2YDDHyB047dXNArbVImasmKBrm8Kr+6B08y8QQ7YG1eVgHi5OIYYclccCf3Ju/BQ78qsMWMniNQ==\"}],\"memo\":\"Sent via irishub client\"}, \"mode\": \"block\"}"
```

## 查询交易（breaking changes）

### 1. Query and parse the latest block information

```json
{
    "block_id": {
        "hash": "DC2EEC73C327BD338EE5667827A4EECA2A2A4752B38D5669CD17EDE07CFB6F30",
        "parts": {
            "total": 1,
            "hash": "1F794EA5185AE489A3D53FD6E19A690373CAB510B981B23E672668CBE0B668E5"
        }
    },
    "block": {
        "header": {
            "version": {
                "block": "11"
            },
            "chain_id": "irishub-1",
            "height": "5",
            "time": "2021-01-18T07:29:21.918234Z",
            "last_block_id": {
                "hash": "889428AAB1975F94C39F44F1BD9C94B2A46E0BC0EFF9AD625939DCB763E82D1F",
                "parts": {
                    "total": 1,
                    "hash": "A9FEDF247839148459E12CD0D8A495A9EE663F7C9E719D85133B27F1D810D52B"
                }
            },
            "last_commit_hash": "35084203D333AF835637F7D9F1FEAA03554AECAB9786EECC6EB5F236156A19F1",
            "data_hash": "A2853A6749C904A8C26C5DB8CB0DD731C44EEAF2AD6AEE2E633DF8F8FD0CA04F",
            "validators_hash": "79E608E448B5B9D0784FDE890506DDE025E0E079CE10B7E01687B9C0E2DFC124",
            "next_validators_hash": "79E608E448B5B9D0784FDE890506DDE025E0E079CE10B7E01687B9C0E2DFC124",
            "consensus_hash": "048091BC7DDC283F77BFBF91D73C44DA58C3DF8A9CBC867405D8B7F3DAADA22F",
            "app_hash": "8F668541D8D565B40373E1492ED6729674539FCB1705437E309522DD491E46DC",
            "last_results_hash": "E3B0C44298FC1C149AFBF4C8996FB92427AE41E4649B934CA495991B7852B855",
            "evidence_hash": "E3B0C44298FC1C149AFBF4C8996FB92427AE41E4649B934CA495991B7852B855",
            "proposer_address": "86C89798CEA6D07FB8550AFDD8DEEA0DA52BFEF4"
        },
        "data": {
            "txs": [
                "CowBCokBChwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kEmkKKmlhYTE4YXduM2s3MHUwNXRsY3VsOHcycW5sNjRnMDAydWo0a2puOTNybhIqaWFhMXc5NzZhNWpyaHNqMDZkcW1yaDJ4OXF4emVsNzRxdGNtYXBrbHhjGg8KBXN0YWtlEgYxMDAwMDASZQpQCkYKHy9jb3Ntb3MuY3J5cHRvLnNlY3AyNTZrMS5QdWJLZXkSIwohA01sYgbsLpw+B9M+p6vyJCh1wfigTWbLpnhNfeDKxKIlEgQKAggBGAESEQoLCgVzdGFrZRICMTAQwJoMGkC+qpaBhJ20qboyU0HWBL0zVlW4klBXGZGsa8n2W1rxIVbq39DZUskGSI8WKNl1stM7QGhycu7YLU30z8vsg8N5"
            ]
        },
        "evidence": {
            "evidence": []
        },
        "last_commit": {
            "height": "4",
            "round": 0,
            "block_id": {
                "hash": "889428AAB1975F94C39F44F1BD9C94B2A46E0BC0EFF9AD625939DCB763E82D1F",
                "parts": {
                    "total": 1,
                    "hash": "A9FEDF247839148459E12CD0D8A495A9EE663F7C9E719D85133B27F1D810D52B"
                }
            },
            "signatures": [
                {
                    "block_id_flag": 2,
                    "validator_address": "86C89798CEA6D07FB8550AFDD8DEEA0DA52BFEF4",
                    "timestamp": "2021-01-18T07:29:21.918234Z",
                    "signature": "UGqkOIqSSzaW/2YzkzfoobcUIizzPcl9BVECl+jwhyMJSkrMnD3DdPUlS2Vd1IAU3u8qQOmP09+m/r5R0gp6CQ=="
                }
            ]
        }
    }
}
```

### 2. Process the transfer transactions included in the block

Step 1: Analyze tx from the block information obtained in [5.1.2.1. Query and parse the latest block information](#5121-Query-and-parse-the-latest-block-information), and decode it in base64

Step 2: Use SHA256 to hash the decoded result to get txhash

```go
tx := "CowBCokBChwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kEmkKKmlhYTE4YXduM2s3MHUwNXRsY3VsOHcycW5sNjRnMDAydWo0a2puOTNybhIqaWFhMXc5NzZhNWpyaHNqMDZkcW1yaDJ4OXF4emVsNzRxdGNtYXBrbHhjGg8KBXN0YWtlEgYxMDAwMDASZQpQCkYKHy9jb3Ntb3MuY3J5cHRvLnNlY3AyNTZrMS5QdWJLZXkSIwohA01sYgbsLpw+B9M+p6vyJCh1wfigTWbLpnhNfeDKxKIlEgQKAggBGAESEQoLCgVzdGFrZRICMTAQwJoMGkC+qpaBhJ20qboyU0HWBL0zVlW4klBXGZGsa8n2W1rxIVbq39DZUskGSI8WKNl1stM7QGhycu7YLU30z8vsg8N5"
txbytes, _ := base64.StdEncoding.DecodeString(tx)
txhash := sha256.Sum256(txbytes)
```

Step 3: Use the txhash to query the corresponding tx information (**Note: There are two types of transfer messages in IRISHub, `cosmos-sdk/MsgSend` and `cosmos-sdk/MsgMultiSend`**)

```json
{
    "height": "5",
    "txhash": "E663768B616B1ACD2912E47C36FEBC7DB0E0974D6DB3823D4C656E0EAB8C679D",
    "data": "0A060A0473656E64",
    "raw_log": "[{\"events\":[{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"send\"},{\"key\":\"sender\",\"value\":\"iaa18awn3k70u05tlcul8w2qnl64g002uj4kjn93rn\"},{\"key\":\"module\",\"value\":\"bank\"}]},{\"type\":\"transfer\",\"attributes\":[{\"key\":\"recipient\",\"value\":\"iaa1w976a5jrhsj06dqmrh2x9qxzel74qtcmapklxc\"},{\"key\":\"sender\",\"value\":\"iaa18awn3k70u05tlcul8w2qnl64g002uj4kjn93rn\"},{\"key\":\"amount\",\"value\":\"1000000uiris\"}]}]}]",
    "logs": [
        {
            "events": [
                {
                    "type": "message",
                    "attributes": [
                        {
                            "key": "action",
                            "value": "send"
                        },
                        {
                            "key": "sender",
                            "value": "iaa18awn3k70u05tlcul8w2qnl64g002uj4kjn93rn"
                        },
                        {
                            "key": "module",
                            "value": "bank"
                        }
                    ]
                },
                {
                    "type": "transfer",
                    "attributes": [
                        {
                            "key": "recipient",
                            "value": "iaa1w976a5jrhsj06dqmrh2x9qxzel74qtcmapklxc"
                        },
                        {
                            "key": "sender",
                            "value": "iaa18awn3k70u05tlcul8w2qnl64g002uj4kjn93rn"
                        },
                        {
                            "key": "amount",
                            "value": "1000000uiris"
                        }
                    ]
                }
            ]
        }
    ],
    "gas_wanted": "200000",
    "gas_used": "69256",
    "tx": {
        "type": "cosmos-sdk/StdTx",
        "value": {
            "msg": [
                {
                    "type": "cosmos-sdk/MsgSend",
                    "value": {
                        "from_address": "iaa18awn3k70u05tlcul8w2qnl64g002uj4kjn93rn",
                        "to_address": "iaa1w976a5jrhsj06dqmrh2x9qxzel74qtcmapklxc",
                        "amount": [
                            {
                                "denom": "uiris",
                                "amount": "1000000"
                            }
                        ]
                    }
                }
            ],
            "fee": {
                "amount": [
                    {
                        "denom": "uiris",
                        "amount": "30000"
                    }
                ],
                "gas": "200000"
            },
            "signatures": [],
            "memo": "",
            "timeout_height": "0"
        }
    },
    "timestamp": "2021-01-18T07:29:21Z"
}
```
