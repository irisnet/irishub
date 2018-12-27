# System Parameters

In IRISnet, there are some special parameters can be modified through on-chain governance. 
All the IRIS holders are able to modify. If the community is not satisfied with certain modifiable 
parameters, it is available to put up a `parameter-change` proposal in governance module.

## Parameters in Governance Module

* In `DepositProcedure` step of governance procedure, the following parameters are up to on-chain governance:
  * Minimum of deposit as `min_deposit` in genesis file
  * Deposit period as `voting_period` in genesis file
* In `VotingProcedure`  step of governance procedure, the following parameters are up to on-chain governance:
   * Voting period as `voting_period` in genesis file
* In `TallyingProcedure`  step of governance procedure, the following parameters are up to on-chain governance:
   * Threshold as `threshold` in genesis file to pass a proposal 
   * Veto percentage as `veto`in genesis file to stop a proposal 
   * Participation percentage as `participation` in genesis file to make the results legitimate

Details in [gov](../governance.md)

## Parameters in Service Module

* `MaxRequestTimeout`   The maximum of waiting blocks for service invocation
* `MinProviderDeposit`  The minimum deposit for service binding

Details in [service](../service.md)

