<!--
Guiding Principles:

Changelogs are for humans, not machines.
There should be an entry for every single version.
The same types of changes should be grouped.
Versions and sections should be linkable.
The latest version comes first.
The release date of each version is displayed.
Mention whether you follow Semantic Versioning.

Usage:

Change log entries are to be added to the Unreleased section under the
appropriate stanza (see below). Each entry should ideally include a tag and
the Github issue reference in the following format:

* (<tag>) \#<issue-number> message

The issue numbers will later be link-ified during the release process so you do
not have to worry about including a link manually, but you can if you wish.

Types of changes (Stanzas):

"Features" for new features.
"Improvements" for changes in existing functionality.
"Deprecated" for soon-to-be removed features.
"Bug Fixes" for any bug fixes.
"Client Breaking" for breaking CLI commands and REST routes used by end-users.
"API Breaking" for breaking exported APIs used by developers building on SDK.
"State Machine Breaking" for any changes that result in a different AppState given same genesisState and txList.

Ref: https://keepachangelog.com/en/1.0.0/
-->

# Changelog

## [Unreleased]

## [v1.3.1] - 2021-02-09

### Bug Fixes

* (modules/service) [\#138](https://github.com/irisnet/irismod/pull/138) Modify price regex in the service pricing schema to adapt to the token denom.

### Improvements

* [\#136](https://github.com/irisnet/irismod/issues/136) Clean up the code and fix module specifications.

## [v1.3.0] - 2021-02-07

### Bug Fixes

* (modules/token) [\#129](https://github.com/irisnet/irismod/pull/129) Fix incorrect calculation for minting token.
* (modules/service) [\#123](https://github.com/irisnet/irismod/pull/123) Fix the key path for owner service fees.
* (modules/service) [\#120](https://github.com/irisnet/irismod/pull/120) Fix DisableServiceBinding event type.
* [\#116](https://github.com/irisnet/irismod/pull/116) Adjust ante check logic.
* (modules/token) [\#102](https://github.com/irisnet/irismod/issues/102) Return error if token baseFee is not a native token.
* (modules/token) [\#100](https://github.com/irisnet/irismod/issues/100) Mint&Edit&Burn only accept symbol as denom.
* (modules/token) [\#99](https://github.com/irisnet/irismod/issues/99) Fix incorrect calculation of deduction amount for burning token.

### Improvements

* [\#111](https://github.com/irisnet/irismod/pull/111) Bump cosmos-sdk version to [v0.41.0](https://github.com/cosmos/cosmos-sdk/releases/tag/v0.41.0).
* (modules/service) [\#118](https://github.com/irisnet/irismod/pull/118) Remove super mode from service.
* [\#115](https://github.com/irisnet/irismod/pull/115) Normalize key path
* (modules/oracle) [\#109](https://github.com/irisnet/irismod/pull/109) Adjust oracle FeedName and ValueJsonPath validation.
* (modules/token) [\#107](https://github.com/irisnet/irismod/issues/107) Use symbol to calculate the default amount of IssueTokenBaseFee.
* (modules/service) [\#105](https://github.com/irisnet/irismod/pull/105) Add service tax account.

## [v1.2.1] - 2021-01-28

### Bug Fixes

* (modules/htlc) [\#79](https://github.com/irisnet/irismod/pull/79) Fix HTLC hash-lock length check.

### Improvements

* [\#83](https://github.com/irisnet/irismod/pull/83) Bump cosmos-sdk version to [v0.40.1](https://github.com/cosmos/cosmos-sdk/releases/tag/v0.40.1).
* [\#83](https://github.com/irisnet/irismod/pull/83) Bump tendermint verion to [v0.34.3](https://github.com/tendermint/tendermint/releases/tag/v0.34.3).
* (modules/service)[\#96](https://github.com/irisnet/irismod/pull/96) Change the event key `response_service` to `respond_service`.
* [\#92](https://github.com/irisnet/irismod/issues/92) Normalize msg and genesis validation.
* (modules/service)[\#86](https://github.com/irisnet/irismod/pull/86) Update service default params.
* (modules/token)[\#85](https://github.com/irisnet/irismod/pull/85) Register denomMetadata to bank module.
* (modules/nft)[\#78](https://github.com/irisnet/irismod/pull/78) File can be used as schema parameters of `GetCmdIssueDenom`.

## [v1.2.0] - 2021-01-22

### Bug Fixes

* (modules/htlc) [\#71](https://github.com/irisnet/irismod/pull/71) Empty owner is allowed in endpoint `/ndt/collections/{denom-id}/supply`.
* (modules/service) [\#70](https://github.com/irisnet/irismod/pull/70) Fix minimum deposit calculation.
* (modules/nft) [\#53](https://github.com/irisnet/irismod/pull/53) Automatically generate key if not specified.
* (modules/service) [\#41](https://github.com/irisnet/irismod/issues/41) Fix update options in `CmdUpdateServiceBinding`.
* (modules/token) [\#36](https://github.com/irisnet/irismod/pull/36) Fix REST API `GET /token/params` .
* [\#33](https://github.com/irisnet/irismod/pull/33) Fix the type of CLI flags.
* (modules/service) [\#32](https://github.com/irisnet/irismod/pull/32) Fix service response validate.
* (modules/service) [\#30](https://github.com/irisnet/irismod/pull/30) Fix random and oracle processing service response.

### Improvements


* [\#66](https://github.com/irisnet/irismod/pull/66) Bump cosmos-sdk version to [v0.40.0](https://github.com/cosmos/cosmos-sdk/releases/tag/v0.40.0).
* [\#66](https://github.com/irisnet/irismod/pull/66) Bump tendermint verion to [v0.34.1](https://github.com/tendermint/tendermint/releases/tag/v0.34.1).
* (modules/coinswap)[\#68](https://github.com/irisnet/irismod/pull/68) Remove standard denom from params and replace prefix `uni:` with `swap/`
* (modules/token)[\#67](https://github.com/irisnet/irismod/pull/67) Add token burn.
* [\#64](https://github.com/irisnet/irismod/pull/64) Add preprocessing before exporting the app state.
* [\#62](https://github.com/irisnet/irismod/pull/62) Add paginate to modules.
* [\#39](https://github.com/irisnet/irismod/issues/39) Change bytes to string in proto.
* (modules/service) [\#38](https://github.com/irisnet/irismod/pull/38) Replace msg_index by internal_index to generate request_context_id.
* [\#37](https://github.com/irisnet/irismod/issues/37) Refactor gRPC gateway REST endpiont.
* [\#22](https://github.com/irisnet/irismod/issues/22) Refactor viper.GetXXX() to cmd.Flags().GetXXX() in CLI.


## [v1.1.1] - 2020-10-20

### Bug Fixes

* (modules/coinswap) [\#27](https://github.com/irisnet/irismod/issues/27) Get liquidity reserve via total supply.
* (modules/service) [\#26](https://github.com/irisnet/irismod/pull/26) Fix deduct service fees and optimized code.
* (modules/coinswap) [\#25](https://github.com/irisnet/irismod/pull/25) Integrate htlc beginblock.
* (modules/service) [\#18](https://github.com/irisnet/irismod/pull/18) Fix incorrect price in request and fix init request when insufficient banlances.

## [v1.1.0] - 2020-09-30

### Features

* Add modules `token`, `nft`, `htlc`, `coinswap`, `service`, `oracle`, `random`, `record`.

<!-- Release links -->

[v1.1.0]: https://github.com/irisnet/irismod/releases/tag/v1.1.0
[v1.1.1]: https://github.com/irisnet/irismod/releases/tag/v1.1.1
[v1.2.0]: https://github.com/irisnet/irismod/releases/tag/v1.2.0
[v1.2.1]: https://github.com/irisnet/irismod/releases/tag/v1.2.1
[v1.3.0]: https://github.com/irisnet/irismod/releases/tag/v1.3.0
[v1.3.1]: https://github.com/irisnet/irismod/releases/tag/v1.3.1