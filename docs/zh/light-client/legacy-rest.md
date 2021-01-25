# Legacy Amino JSON REST

## API 端口、激活方式、配置方法



## 完整的 API 端点清单
| Legacy REST Endpoint                                                            | Description                                                         | New gGPC-gateway REST Endpoint                                                                        |
| ------------------------------------------------------------------------------- | ------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------- |
| `GET /node_info`                                                       | Get the properties of the connected node                            | `GET /cosmos/base/tendermint/v1beta1/node_info`                                                                  |
| `GET /syncing`                                                               | Get syncing state of node                                    | `GET /cosmos/base/tendermint/v1beta1/syncing`                                                                        |
| `GET /blocks/latest`                                                          | Get the latest block                                                | `GET /cosmos/base/tendermint/v1beta1/blocks/latest`                                                            |
| `GET /blocks/{height}`                                                          | Get a block at a certain height                                  | `GET /cosmos/base/tendermint/v1beta1/blocks/{height}`                                                            |
| `GET /validatorsets/latest`                                                          | Get the latest validator set                                  | `GET /cosmos/base/tendermint/v1beta1/validatorsets/latest`                                                            |
| `GET /validatorsets/{height}`                                                          | Get a validator set a certain height                        | `GET /cosmos/base/tendermint/v1beta1/validatorsets/{height}`                                                            |
| `GET /txs/{hash}`                                                               | Query tx by hash                                                    | `GET /cosmos/tx/v1beta1/txs/{hash}`                                                                   |
| `GET /txs`                                                                      | Query tx by events                                                  | `GET /cosmos/tx/v1beta1/txs`                                                                          |
| `POST /txs`                                                                     | Broadcast tx                                                        | `POST /cosmos/tx/v1beta1/txs`                                                                         |
| `POST /txs/encode`                                                              | Encodes an Amino JSON tx to an Amino binary tx                      | N/A, use Protobuf directly                                                                            |
| `POST /txs/decode`                                                              | Decodes an Amino binary tx into an Amino JSON tx                    | N/A, use Protobuf directly                                                                            |
| `GET /bank/balances/{address}`                                                  | Get the balance of an address                                       | `GET /cosmos/bank/v1beta1/balances/{address}/{denom}`                                                 |
| `POST /bank/accounts/{address}/transfers`                                       | Send coins from one account to another                              | N/A, use Protobuf directly                                                                            |
| `GET auth/accounts/{address}`                                                  | Get the account information on blockchain                            | `GET /cosmos/auth/v1beta1/accounts/{address}`                                                 |
| `GET /staking/delegators/{delegatorAddr}/delegations`                           | Get all delegations from a delegator                                | `GET /cosmos/staking/v1beta1/delegators/{delegatorAddr}/delegations`                                  |
| `POST /staking/delegators/{delegatorAddr}/delegations`                         | Submit delegation                                                    | N/A, use Protobuf directly                                                                            |
| `GET /staking/delegators/{delegatorAddr}/delegations/{validatorAddr}`           | Query a delegation between a delegator and a validator              | `GET /cosmos/staking/v1beta1/delegators/{delegatorAddr}/delegations/{validatorAddr}`                  |
| `GET /staking/delegators/{delegatorAddr}/unbonding_delegations`                 | Get all unbonding delegations from a delegator                      | `GET /cosmos/staking/v1beta1/delegators/{delegatorAddr}/unbonding_delegations`                        |
| `POST /staking/delegators/{delegatorAddr}/unbonding_delegations`                | Submit an unbonding delegation                                       | N/A, use Protobuf directly                                                                            |
| `GET /staking/delegators/{delegatorAddr}/unbonding_delegations/{validatorAddr}` | Query all unbonding delegations between a delegator and a validator | `GET /cosmos/staking/v1beta1/delegators/{delegatorAddr}/unbonding_delegations/{validatorAddr}`        |
| `GET /staking/redelegations`                                                    | Query redelegations                                                 | `GET /cosmos/staking/v1beta1/v1beta1/delegators/{delegator_addr}/redelegations`                       |
| `POST /staking/delegators/{delegatorAddr}/redelegations`                        | Submit a redelegations                                               | N/A, use Protobuf directly                                                                            |
| `GET /staking/delegators/{delegatorAddr}/validators`                            | Query all validators that a delegator is bonded to                  | `GET /cosmos/staking/v1beta1/delegators/{delegatorAddr}/validators`                                   |
| `GET /staking/delegators/{delegatorAddr}/validators/{validatorAddr}`            | Query a validator that a delegator is bonded to                     | `GET /cosmos/staking/v1beta1/delegators/{delegatorAddr}/validators/{validatorAddr}`                   |
| `GET /staking/validators`                                                       | Get all validators                                                  | `GET /cosmos/staking/v1beta1/validators`                                                              |
| `GET /staking/validators/{validatorAddr}`                                       | Get a single validator info                                         | `GET /cosmos/staking/v1beta1/validators/{validatorAddr}`                                              |
| `GET /staking/validators/{validatorAddr}/delegations`                           | Get all delegations to a validator                                  | `GET /cosmos/staking/v1beta1/validators/{validatorAddr}/delegations`                                  |
| `GET /staking/validators/{validatorAddr}/unbonding_delegations`                 | Get all unbonding delegations from a validator                      | `GET /cosmos/staking/v1beta1/validators/{validatorAddr}/unbonding_delegations`                        |
| `GET /staking/pool`                                                             | Get the current state of the staking pool                           | `GET /cosmos/staking/v1beta1/pool`                                                                    |
| `GET /staking/parameters`                                                       | Get the current staking parameter values                            | `GET /cosmos/staking/v1beta1/params`                                                                  |
| `GET /slashing/signing_infos`                                                   | Get all signing infos                                               | `GET /cosmos/slashing/v1beta1/signing_infos`                                                          |
| `POST /slashing/validators/{validatorAddr}/unjail`                             | Unjail a jailed validator                                           | N/A, use Protobuf directly                                                                            |
| `GET /slashing/parameters`                                                      | Get slashing parameters                                             | `GET /cosmos/slashing/v1beta1/params`                                                                 |
| `POST /gov/proposals`                                                                   | Submit a proposal                                                   | N/A, use Protobuf directly                                                                            |
| `GET /gov/proposals`                                                            | Get all proposals                                                   | `GET /cosmos/gov/v1beta1/proposals`                                                                   |
| `POST /gov/proposals/param_change`                                              | Generate a parameter change proposal transactionl                   | N/A, use Protobuf directly                                                                            |
| `GET /gov/proposals/{proposal-id}`                                              | Get proposal by id                                                  | `GET /cosmos/gov/v1beta1/proposals/{proposal-id}`                                                     |
| `GET /gov/proposals/{proposal-id}/proposer`                                     | Get proposer of a proposal                                          | `GET /cosmos/gov/v1beta1/proposals/{proposal-id}` (Get proposer from `Proposal` struct)               |
| `GET /gov/proposals/{proposal-id}/deposits`                                     | Get deposits of a proposal                                          | `GET /cosmos/gov/v1beta1/proposals/{proposal-id}/deposits`                                            |
| `POST /gov/proposals/{proposal-id}/deposits`                                              | Deposit tokens to a proposal                                 | N/A, use Protobuf directly                                                                            |
| `GET /gov/proposals/{proposal-id}/deposits/{depositor}`                         | Get depositor a of deposit                                          | `GET /cosmos/gov/v1beta1/proposals/{proposal-id}/deposits/{depositor}`                                |
| `GET /gov/proposals/{proposal-id}/votes`                                        | Get votes of a proposal                                             | `GET /cosmos/gov/v1beta1/proposals/{proposal-id}/votes`                                               |
| `POST /gov/proposals/{proposal-id}/votes`                                              | Vote a proposal                                                    | N/A, use Protobuf directly                                                                            |
| `GET /gov/proposals/{proposal-id}/votes/{voter}`                                 | Get voted information by voterAddr.                                               | `/cosmos/gov/v1beta1/proposals/{proposal_id}/votes/{voter}`                                        |
| `GET /gov/proposals/{proposal-id}/tally`                                        | Get tally of a proposal                                             | `GET /cosmos/gov/v1beta1/proposals/{proposal-id}/tally`                                               |
| `GET /gov/parameters/deposit`                                                    | Get governance deposit parameters                                 | `GET /cosmos/gov/v1beta1/params/{type}`                                                               |
| `GET /gov/parameters/tallying`                                                    | Query governance tally parameters                                 | `GET /cosmos/gov/v1beta1/params/{type}`                                                               |
| `GET /gov/parameters/voting`                                                    | Get governance voting parameters                                 | `GET /cosmos/gov/v1beta1/params/{type}`                                                               |
| `GET /distribution/delegators/{delegatorAddr}/rewards`                          | Get the total rewards balance from all delegations                  | `GET /cosmos/distribution/v1beta1/v1beta1/delegators/{delegator_address}/rewards`                     |
| `POST /distribution/delegators/{delegatorAddr}/rewards`                         | Withdraw all delegator rewards                                      | N/A, use Protobuf directly                                                                            |
| `GET /distribution/delegators/{delegatorAddr}/rewards/{validatorAddr}`          | Query a delegation reward                                           | `GET /cosmos/distribution/v1beta1/delegators/{delegatorAddr}/rewards/{validatorAddr}`                 |
| `POST /distribution/delegators/{delegatorAddr}/rewards/{validatorAddr}`          | Withdraw a delegation reward                                          | N/A, use Protobuf directly                                                                            |
| `GET /distribution/delegators/{delegatorAddr}/withdraw_address`                 | Get the rewards withdrawal address                                  | `GET /cosmos/distribution/v1beta1/delegators/{delegatorAddr}/withdraw_address`                        |
| `POST /distribution/delegators/{delegatorAddr}/withdraw_address`          | Replace the rewards withdrawal address                                          | N/A, use Protobuf directly                                                                            |
| `GET /distribution/validators/{validatorAddr}`                                  | Validator distribution information                                  | `GET /cosmos/distribution/v1beta1/validators/{validatorAddr}`                                         |
| `GET /distribution/validators/{validatorAddr}/outstanding_rewards`              | Outstanding rewards of a single validator                           | `GET /cosmos/distribution/v1beta1/validators/{validatorAddr}/outstanding_rewards`                     |
| `GET /distribution/validators/{validatorAddr}/rewards`                          | Commission and self-delegation rewards of a single a validator      | `GET /cosmos/distribution/v1beta1/validators/{validatorAddr}/rewards`                                 |
| `POST /distribution/validators/{validatorAddr}/rewards`                      | Withdraw the validator's rewards                                          | N/A, use Protobuf directly                                                                            |
| `GET /distribution/community_pool`                                              | Get the amount held in the community pool                           | `GET /cosmos/distribution/v1beta1/community_pool`                                                     |
| `GET /distribution/parameters`                                                  | Get the current distribution parameter values                       | `GET /cosmos/distribution/v1beta1/params`                                                             |



## 高优先级查询端点的 breaking changes

- `GET /blocks/latest`&&`GET /blocks/{height}`

    - 具体变化：

        - 不再使用block_meta字段，原block_meta字段结构中block_id移出至第一层，header字段移到block中
        - 不再使用num_txs、total_txs字段，交易数可遍历txs字段获得

    - json示例：

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

    - 具体变化：

        - 这两个接口已被取消
        - 可直接调用 Tendermint RPC 的 block_results 接口获得相关数据(查询结果 events 字段中的 key、value 皆为 base64 编码)

    - json示例：

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
                      "log":"[{"events":[{"type":"message","attributes":[{"key":"action","value":"send"},{"key":"sender","value":"iaa18awn3k70u05tlcul8w2qnl64g002uj4kjn93rn"},{"key":"module","value":"bank"}]},{"type":"transfer","attributes":[{"key":"recipient","value":"iaa1w976a5jrhsj06dqmrh2x9qxzel74qtcmapklxc"},{"key":"sender","value":"iaa18awn3k70u05tlcul8w2qnl64g002uj4kjn93rn"},{"key":"amount","value":"1000000uiris"}]}]}]",
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

    - 具体变化：

        - 不再使用tags，取而代之的是events字段
        - 不再使用result字段，原result中字段移出至第一层
        - 不再使用coin_flow字段

    - json示例：

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

    - 具体变化：

        - 该接口已被取消
        - 原接口中coins字段可通过/bank/balances/{address}接口查询
        - 原接口中其他字段可通过/auth/accounts/{address}接口查询

    - json示例：
      ```json
      // /bank/balances/{address}
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
      
      // /auth/accounts/{address}
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

与 cosmoshub-3 主网对接的代码一样，交易结构如下:

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

其中，IRISHub地址前缀使用 `iaa1`，影响字段：
- value.msg.value.from_adress
- value.msg.value.to_address

Denom 使用 `uiris` （1iris = 10<sup>6</sup>uiris），影响字段：
- value.msg.value.amount.denom
- value.fee.amount.denom



## 广播交易（完全向后兼容）

与 cosmoshub-3 主网对接的代码一样，调用接口 `POST /txs` 发送交易，示例：

```bash
curl -X POST "http://localhost:1317/txs" -H "accept: application/json" -H "Content-Type: application/json" -d "{ \"tx\": {\"msg\":[{\"type\":\"cosmos-sdk/MsgSend\",\"value\":{\"from_address\":\"iaa1rkgdpj6fyyyu7pnhmc3v7gw9uls4mnajvzdwkt\",\"to_address\":\"iaa1q6t5439f0rkvkzl38m0f43e0kpv3mx7x2shlq8\",\"amount\":[{\"denom\":\"uiris\",\"amount\":\"1000000\"}]}}],\"fee\":{\"amount\":[{\"denom\":\"uiris\",\"amount\":\"30000\"}],\"gas\":\"200000\"},\"signatures\":[{\"pub_key\":{\"type\":\"tendermint/PubKeySecp256k1\",\"value\":\"AxGagdsRTKni/h1+vCFzTpNltwoiU7SwIR2dg6Jl5a//\"},\"signature\":\"Pu8yiRVO8oB2YDDHyB047dXNArbVImasmKBrm8Kr+6B08y8QQ7YG1eVgHi5OIYYclccCf3Ju/BQ78qsMWMniNQ==\"}],\"memo\":\"Sent via irishub client\"}, \"mode\": \"block\"}"
```


## 查询交易（breaking changes）

##### 1. 查询并解析最新高度的区块信息
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
##### 2. 对区块包含的转账交易进行处理
第 1 步，从 [5.1.2.1. 查询并解析最新高度的区块信息](#5121-查询并解析最新高度的区块信息) 获取的区块信息中解析得到tx，对其进行base64解码
第 2 步，将解码得到的结果使用SHA256进行hash，得到txhash
```go
tx := "CowBCokBChwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kEmkKKmlhYTE4YXduM2s3MHUwNXRsY3VsOHcycW5sNjRnMDAydWo0a2puOTNybhIqaWFhMXc5NzZhNWpyaHNqMDZkcW1yaDJ4OXF4emVsNzRxdGNtYXBrbHhjGg8KBXN0YWtlEgYxMDAwMDASZQpQCkYKHy9jb3Ntb3MuY3J5cHRvLnNlY3AyNTZrMS5QdWJLZXkSIwohA01sYgbsLpw+B9M+p6vyJCh1wfigTWbLpnhNfeDKxKIlEgQKAggBGAESEQoLCgVzdGFrZRICMTAQwJoMGkC+qpaBhJ20qboyU0HWBL0zVlW4klBXGZGsa8n2W1rxIVbq39DZUskGSI8WKNl1stM7QGhycu7YLU30z8vsg8N5"
txbytes, _ := base64.StdEncoding.DecodeString(tx)
txhash := sha256.Sum256(txbytes)
```
第 3 步，使用得到的txhash查询对应的tx信息（**注意：IRISHub中有两种转账消息类型，`cosmos-sdk/MsgSend` 和 `cosmos-sdk/MsgMultiSend`**）
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