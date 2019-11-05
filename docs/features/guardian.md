# Guardian

## Introduction

IRISnet introduce two types of privileged system user controlled by foundations, the profiler and the trustee.

* Profiler privileges
    1. Submit software upgrade/halt proposal by governance.
    2. Invocate a service by `profiling` mode, under which service fees can be exempted.

* Trustee privileges
    1. To be the destination address if the usage type of a `CommunityTaxUsage` proposal is `Distribute` or `grant`.
    2. Send `withdraw-tax` transaction to withdraw coins to an account from system service fee tax pool.

* Genesis Profiler/Genesis Trustee (Defined in genesis.json)
    1. Only Genesis Profiler can add/delete Ordinary Profiler account
    2. Only Genesis Trustee can add/delete Ordinary Trustee account

## Usage Scenario

1. Add Profiler and Trustee

    Add Profiler (Genesis Profiler account only)

    ```bash
    iriscli guardian add-profiler --address=<profiler-address> --description=<profiler-description> --chain-id=<chain-id> --from=<key-name> --fee=0.3iris
    ```

    Add Trustee (Genesis Trustee account only)

    ```bash
    iriscli guardian add-trustee --address=<trustee-address> --description=<trustee-description> --chain-id=<chain-id> --from=<key-name> --fee=0.3iris
    ```

2. Query Profiler and Trustee list

    Query Profiler list

    ```bash
    iriscli guardian profilers
    ```

    Query Trustee list

    ```bash
    iriscli guardian trustees
    ```

3. Profiler submit software upgrade/halt proposal

    Details in [upgrade](upgrade.md)

4. Profiler call a service by `profiling` mode

    Service fee exempted

    ```bash
    iriscli service call --def-chain-id=<def-chain-id> --service-name=<service-name> --method-id=<method-id> --bind-chain-id=<bind-chain-id> --provider=<provider-address> --service-fee=1iris --request-data=<request-data> --chain-id=<chain-id> --from=<key-name> --fee=0.3iris --profiling=true
    ```

5. Trustee to be the destination address of `CommunityTaxUsage` proposal

    Details in [governance](governance.md#proposals-on-community-funds-usage)

6. Trustee withdraw service fee tax

    ```bash
    iriscli service withdraw-tax --dest-address=<destination-address> --withdraw-amount=1iris --chain-id=<chain-id> --from=<key-name> --fee=0.3iris
    ```

7. Delete Profiler and Trustee (Genesis Trustee account only)

    Delete Profiler (Genesis Trustee account only)

    ```bash
    iriscli guardian delete-profiler --chain-id=<chain-id> --from=<key-name> --fee=0.3iris --address=<profiler-address>
    ```

    Delete Trustee (Genesis Trustee account only)

    ```bash
    iriscli guardian delete-trustee --chain-id=<chain-id> --from=<key-name> --fee=0.3iris --address=<trustee-address>
    ```
