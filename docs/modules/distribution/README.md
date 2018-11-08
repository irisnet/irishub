# Introduction 

This module is in charge of distributing collected transaction fee and inflated token to all validators and delegators. To reduce computation stress, a lazy distribution strategy is brought in. `lazy` means that the benefit won't be paid directly to contributors. The contributors are required to explicitly send transactions to withdraw their benefit, otherwise, their benefit will be kept in the global pool. 

# Command Line Interface

1. Set withdraw address

    Validator operators or delegators can specify other address as their withdraw address. If no other address has been specified, the delegator address or validator self-delegator address will be used as default address.
    For example: 
    ```bash
    iriscli distribution set-withdraw-addr faa1syva9fvh8m6dc6wjnjedah64mmpq7rwwz6nj0k --from mykey --fee=0.004iris --chain-id=irishub-test
    ```

2. Query withdraw address

    For example:
    ```bash
    iriscli distribution withdraw-address faa1vm068fnjx28zv7k9kd9j85wrwhjn8vfsxfmcrz
    ```
    If the given delegator doesn't specify other withdraw address, the query result will be empty.

3. Query delegation distribution information

    For example:
    ```bash
    iriscli distribution delegation-distr-info --address-delegator=faa1ezzh0humhy3329xg4avhcjtay985nll0zswc5j --address-validator=fva1ezzh0humhy3329xg4avhcjtay985nll0hpyhf4
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

4. Query delegator distribution information

    For example: 
    ```bash
    iriscli distribution delegator-distr-info faa1ezzh0humhy3329xg4avhcjtay985nll0zswc5j
    ```
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

5. Query validator distribution information
    ```bash
    
    ```

6. withdraw rewards 
    ```bash
    
    ```

# Rest API Interface
