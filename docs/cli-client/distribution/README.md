# Introduction 

This document description how to use the the command line interface of distribution module.

# Query interface

By default, trust-node mode is enable. If you don't trust the connected node, just append --trust-node=false in each query command.

1. Query withdraw address

    For example:
    ```bash
    iriscli distribution withdraw-address faa1vm068fnjx28zv7k9kd9j85wrwhjn8vfsxfmcrz
    ```
    If the given delegator doesn't specify other withdraw address, the query result will be empty.

2. Query delegation distribution information

    For example:
    ```bash
    iriscli distribution delegation-distr-info --address-delegator=faa1ezzh0humhy3329xg4avhcjtay985nll0zswc5j \
    --address-validator=fva1ezzh0humhy3329xg4avhcjtay985nll0hpyhf4
    ```
    Query result:
    ```json
    {
      "delegator_addr": "faa1ezzh0humhy3329xg4avhcjtay985nll0zswc5j",
      "val_operator_addr": "fva1ezzh0humhy3329xg4avhcjtay985nll0hpyhf4",
      "del_pool_withdrawal_height": "4044"
    }
    ```
    The above response means this delegator has send transaction to withdraw reward at height 4044 or the delegation is created on height 4044.

2. Query delegator distribution information

    For example: 
    ```bash
    iriscli distribution delegator-distr-info faa1ezzh0humhy3329xg4avhcjtay985nll0zswc5j
    ```
    Query result:
    ```json
    [
      {
        "delegator_addr": "faa1ezzh0humhy3329xg4avhcjtay985nll0zswc5j",
        "val_operator_addr": "fva14a70gzu0v2w8dlfx462c9sldvja24qaz6vv4sg",
        "del_pool_withdrawal_height": "10859"
      },
      {
        "delegator_addr": "faa1ezzh0humhy3329xg4avhcjtay985nll0zswc5j",
        "val_operator_addr": "fva1ezzh0humhy3329xg4avhcjtay985nll0hpyhf4",
        "del_pool_withdrawal_height": "4044"
      }
    ]
    ```

4. Query validator distribution information

    For example: 
    ```bash
    iriscli distribution delegator-distr-info faa1ezzh0humhy3329xg4avhcjtay985nll0zswc5j
    ```
    Query result:
    ```json
    {
      "operator_addr": "faa1ezzh0humhy3329xg4avhcjtay985nll0zswc5j",
      "fee_pool_withdrawal_height": "416",
      "del_accum": {
        "update_height": "416",
        "accum": "0.0000000000"
      },
      "del_pool": "0.0000000000iris",
      "val_commission": "0.0000000000iris"
    }
    ```

# Send transactions interface

1. Set withdraw address

    Validator operators or delegators can specify other address as their withdraw address. If no other address has been specified, the delegator address or validator self-delegator address will be used as default address.
    For example: 
    ```bash
    iriscli distribution set-withdraw-addr faa1syva9fvh8m6dc6wjnjedah64mmpq7rwwz6nj0k --from mykey --fee=0.004iris --chain-id=irishub-test
    ```
    
2. withdraw rewards 

    1. Only withdraw the delegation reward from a given validator
    ```bash
    iriscli distribution withdraw-rewards --only-from-validator fva134mhjjyyc7mehvaay0f3d4hj8qx3ee3w3eq5nq --from mykey --fee=0.004iris --chain-id=irishub-test
    ```
    2. Withdraw all delegation reward of a delegator
    ```bash
    iriscli distribution withdraw-rewards --from mykey --fee=0.004iris --chain-id=irishub-test
    ```
    3. If the delegator is a onwer of a validator, withdraw all delegation reward and validator reward:
    ```bash
    iriscli distribution withdraw-rewards --is-validator=true --from mykey --fee=0.004iris --chain-id=irishub-test
    ```