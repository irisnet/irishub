# Distribution User Guide

## Introduction 

This module is in charge of distributing collected transaction fee and inflated token to all validators and delegators. To reduce computation stress, a lazy distribution strategy is brought in. `lazy` means that the benefit won't be paid directly to contributors automatically. The contributors are required to explicitly send transactions to withdraw their benefit, otherwise, their benefit will be kept in the global pool. 

## Usage Scenario

1. Set withdraw address

    A delegator may have multiple irishub wallet address. Suppose one of the wallets has many iris token and part of these tokens have been delegated to a validator. The delegator may hope the delegation reward can be paid to another wallet, thus the delegator will have explicit idea about how many tokens he/she has earned.However, by default, the reward will be paid to the wallet(marked as `A`) address which send the delegation transaction. To set another wallet(marked as `B`) as the paid address, delegator need to send another transaction from wallet `A`. The referenced command can be:
    ```bash
    iriscli distribution set-withdraw-addr [address of wallet B] --fee=0.004iris --from=[key name of wallet A] --chain-id=[chain-id]
    ```  
    To verify the whether the above operation take effect, delegator can execute the following command.
    ```bash
    iriscli distribution withdraw-address [address of wallet A]
    ```
    If set withdraw operation is success, the query response must equal to the address of wallet B.

2. Withdraw reward 

    According to the introduction section, our delegation reward won't be paid to our wallet automatically, we have to send transactions to withdraw reward.
    Suppose a delegator operate a validator(marked as `validatorA`), besides, it also has delegations on two other validators(marked as `validatorB` and `validatorC`). All delegations are created from wallet A.
    
    1. Only withdraw the self-delegation reward of from validatorA:
        ```bash
        iriscli distribution withdraw-rewards --only-from-validator [address of validatorA] --from [key name of wallet A] --fee=0.004iris --chain-id=[chain-id]
        ```
    2. Withdraw all delegation reward:
        ```bash
        iriscli distribution withdraw-rewards --from [key name of wallet A] --fee=0.004iris --chain-id=[chain-id]
        ```
    3. Withdraw all delegation reward including commission benefit of `validatorA` :
        ```bash
        iriscli distribution withdraw-rewards --is-validator=true --from [key name of wallet A] --fee=0.004iris --chain-id=[chain-id]
        ```
        
3. Query reward token

    Execute the command to get the estimated inflation rewards :
    
    ```bash
    iriscli distribution withdraw-rewards --from=bob  --dry-run --chain-id=test-irishub --fee=0.004iris
    ```
    
    Output is the following，`withdraw-reward-total`is your estimated inflation rewards：
    
    ```bash
    estimated gas = 6032
    simulation code = 0
    simulation log = Msg 0:
    simulation gas wanted = 200000
    simulation gas used = 6032
    simulation fee amount = 0
    simulation fee denom =
    simulation tag action = withdraw-delegator-rewards-all
    simulation tag delegator = faa1yclscskdtqu9rgufgws293wxp3njsesxtplqxd
    simulation tag withdraw-reward-total = 1308135156755646.7440000000iris-atto
    simulation tag withdraw-reward-from-validator-fva1yclscskdtqu9rgufgws293wxp3njsesx7s40m2 = 1308135156755646.7440000000iris-atto
    simulation tag action = withdraw_delegation_rewards_all    
    
    ```