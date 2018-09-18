# Changelog

## 0.4.0

*September 6th, 2018*

BREAKING CHANGES:

- [cosmos-sdk] Upgrade to cosmos-sdk v0.23.0
    - Change the address prefix format:
      - cosmosaccaddr --> faa
      - cosmosaccpub --> fap
      - cosmosvaladdr --> fva
      - cosmosvalpub --> fvp
    - Adjust the Route & rootMultiStore Commit for software upgrade
    - Must specify gas and fee in both command line and rest api
    - The fee should be iris token and the token amount should be no less than 2*(10^10)*gas

FEATURES:

- [tendermint] Upgrade to Tendermint v0.22.6
    - Store the pre-state to support the replay function
- [cosmos-sdk] Upgrade to cosmos-sdk v0.23.0
    - Add the paramProposal and softwareUpgradeProposal in gov module
    - Improve fee token mechanism to more reasonably deduct transaction fee and achieve more ability to defent DDOS attack.
    - Introduce the global parameter module

BUG FIXES

- Default account balance in genesis
- Fix iris version issue
- Fix the unit conflict issue in slashing
- Check the voting power when create validator
- Fix evidence amimo register issue


## 0.3.0

*July 30th, 2018*

BREAKING CHANGES:

- [tendermint] Upgrade to Tendermint v0.22.2
    - Default ports changed from 466xx to 266xx
    - ED25519 addresses are the first 20-bytes of the SHA256 of the raw 32-byte pubkey (Instead of RIPEMD160)
- [cosmos-sdk] Upgrade to cosmos-sdk v0.22.0
- [monitor] Move `iriscli monitor` subcommand to `iris monitor`

FEATURES:

- [lcd] /tx/send is now the only endpoint for posing transaction to irishub; aminofied all transaction messages 
- [monitor] Improve the metrics for iris-monitor 

BUG FIXES

- [cli] solve the issue of iriscli stake sign-info

##

## 0.2.0

*July 19th, 2018*

BREAKING CHANGES:

- [tendermint] Upgrade to Tendermint v0.21.0
- [cosmos-sdk] Upgrade to cosmos-sdk v0.19.1-rc1

FEATURES:

- [lcd] code refactor

- [cli] improve sendingand querying the  transactions 

- [monitor]export data which is collected by Prometheus Server

  ​

##  