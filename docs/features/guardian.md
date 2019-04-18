# Guardian User Guide

## Introduction
IRISnet introduce two types of privileged system user controlled by foundations, the profiler and the trustee. 

* Profiler privileges
    1. Submit software upgrade/halt proposal by governance.
    2. Invocate a service by `profiling` mode, under which service fees can be exempted.
    
* Trustee privileges
    1. To be the destination address if the usage type of a `TxTaxUsage` proposal is `Distribute` or `grant`.
    2. Send `withdraw-tax` transaction to withdraw coins to an account from system service fee tax pool.
    
* Genesis Profiler (Defined in genesis.json)
    1. Only Genesis Profiler can add/delete Ordinary Profiler account
    2. Only Genesis Profiler can add/delete Trustee account
    
## Usage Scenario
1. Add Profiler and Trustee (Genesis Profiler account only)

    Add Profiler
    ```shell
    iriscli guardian add-profiler --address=<profiler_address> --description=<profiler_description> --chain-id=<chain-id> --from=<key_name> --fee=0.3iris 
    ```

    Add Trustee
    ```shell
    iriscli guardian add-trustee --address=<trustee_address> --description=<trustee_description> --chain-id=<chain-id> --from=<key_name> --fee=0.3iris 
    ```
    
2. Query Profiler and Trustee list

    Query Profiler list
    ```shell
    iriscli guardian profilers
    ```
    Query Trustee list
    ```shell
    iriscli guardian trustees
    ```
    
3. Profiler submit software upgrade/halt proposal

    Details in [upgrade](upgrade.md)

4. Profiler call a service by `profiling` mode

    Service fee exempted
    ```shell
    iriscli service call --def-chain-id=<def-chain-id> --service-name=<service_name> --method-id=<method_id> --bind-chain-id=<bind-chain-id> --provider=<provider_address> --service-fee=1iris --request-data=<request_data> --chain-id=<chain-id> --from=<key_name> --fee=0.3iris --profiling=true
    ```
    
5. Trustee to be the destination address of `TxTaxUsage` proposal

    Details in [governance](governance.md#proposals-on-community-funds-usage)
    
6. Trustee withdraw service fee tax

    ```shell
    iriscli service withdraw-tax --dest-address=<destination_address> --withdraw-amount=1iris --chain-id=<chain-id> --from=<key_name> --fee=0.3iris
    ```
    
7. Delete Profiler and Trustee (Genesis Profiler account only)

    Delete Profiler
    ```shell
    iriscli guardian delete-profiler --chain-id=<chain-id> --from=<key_name> --fee=0.3iris --address=<profiler_address>
    ```
    
    Delete Trustee
    ```shell
    iriscli guardian delete-trustee --chain-id=<chain-id> --from=<key_name> --fee=0.3iris --address=<trustee_address>
    ```
