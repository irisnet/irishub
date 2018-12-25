# Gov Params

In IRISnet, there are some special parameters can be modified through on-chain governance. All the IRIS holders are able to modify. If the community is not satisfied with certain modifiable parameters, it is available to set the proper values in governance module.

## Gov Module

* `DepositProcedure`   Parameters in deposit period (The minimum of deposit, deposit period)
* `VotingProcedure`    Parameters in voting period（Voting period）
* `TallyingProcedure`  Parameters in tallying period（The standards of voting）

Details in [gov](../governance.md)

## Service Module

* `MaxRequestTimeout`   The maximum of waiting blocks for service invocation
* `MinProviderDeposit`  The minimum deposit for service binding

Details in [service](../service.md)


 "terminator_period": "20000",
      "starting_proposalID": "1",
      "deposits": null,
      "votes": null,
      "proposals": null,
      "deposit_period": {
        "min_deposit": [
          {
            "denom": "iris-atto",
            "amount": "1000000000000000000000"
          }
        ],
        "max_deposit_period": "172800000000000"
      },
      "voting_period": {
        "voting_period": "172800000000000"
      },
      "tallying_procedure": {
        "threshold": "0.5000000000",
        "veto": "0.3340000000",
        "participation": "0.6670000000"
      }
    },