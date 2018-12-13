# Guardian User Guide

## Introduction
IRISnet introduce two types of privileged system user controlled by foundations, the profiler and the trustee. 

* Profiler privileges[TODO]
    1. Submit software upgrade/halt proposal by governance.
    2. Invocate a service by profiling mode, under which service fees can be exempted.
    
* Trustee privileges[TODO]
    1. To be the destination address if the usage type of a `TxTaxUsage` proposal is `Distribute` or `grant`.
    2. Send transaction to withdraw coins to an account from system service fee tax pool.
    
## Usage Scenario
1. Add profiler

    Only a profiler can add a new one.
    ```shell
    iriscli guardian add-profiler --profiler-address=[profiler address] --profiler-name=[name] --chain-id=[chain-id] --from=[key name] --fee=0.004iris 
    ```
    
2. Query profiler and trustee list

    Query profiler list
    ```shell
    iriscli guardian profilers
    ```
    Query trustee list
    ```shell
    iriscli guardian trustees
    ```