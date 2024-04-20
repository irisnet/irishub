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

## [v1.9.0]

### Improvements

* [\#387](https://github.com/irisnet/irismod/pull/387) Preparatory work for realizing erc20.
* [\#393](https://github.com/irisnet/irismod/pull/393) Implement DeployERC20 by gov proposal .
* [\#394](https://github.com/irisnet/irismod/pull/394) Support swap native token from erc20 token.
* [\#395](https://github.com/irisnet/irismod/pull/395) Implement SwapToERC20.
* [\#400](https://github.com/irisnet/irismod/pull/400) Implement swap native token from erc20 contract.
* [\#402](https://github.com/irisnet/irismod/pull/402) Add cli for swap from/to ERC20.
* [\#404](https://github.com/irisnet/irismod/pull/404) Implement token Balances api.
* [\#405](https://github.com/irisnet/irismod/pull/405) Rolled back the cosmos-sdk version to v0.47.9.
* [\#409](https://github.com/irisnet/irismod/pull/409) feat: add a switch for enabling erc20 swap function.

### Bug Fixes

* [\#401](https://github.com/irisnet/irismod/pull/401) Add token test & fix error.
* [\#407](https://github.com/irisnet/irismod/pull/407) Fix some bugs about erc20 in token module.

### API Breaking Changes

* [\#403](https://github.com/irisnet/irismod/pull/403) Standardized parameter naming

## [v1.8.0]

### Improvements

* [\#382](https://github.com/irisnet/irismod/pull/382) Fix dependabot alerts.
* [\#381](https://github.com/irisnet/irismod/pull/381) Forbidden to mint nft under ibc class.
* [\#378](https://github.com/irisnet/irismod/pull/378) Support cosmos-sdk `x/nft` api.
* [\#376](https://github.com/irisnet/irismod/pull/376) Bump up cosmos-sdk to v0.47.4.
* [\#374](https://github.com/irisnet/irismod/pull/374) Move proto file to irismod directory.
* [\#369](https://github.com/irisnet/irismod/pull/369) Support cosmos depinject.
* [\#368](https://github.com/irisnet/irismod/pull/368) Remove sdk simapp dep.
* [\#367](https://github.com/irisnet/irismod/pull/367) Lint proto.
* [\#365](https://github.com/irisnet/irismod/pull/365) Remove duplicate events.
* [\#364](https://github.com/irisnet/irismod/pull/364) Migrate token params.
* [\#363](https://github.com/irisnet/irismod/pull/363) Migrate service params.
* [\#362](https://github.com/irisnet/irismod/pull/362) Migrate htlc params.
* [\#361](https://github.com/irisnet/irismod/pull/361) Migrate farm params.
* [\#360](https://github.com/irisnet/irismod/pull/360) Migrate coinswap params.

### Bug Fixes

* [\#380](https://github.com/irisnet/irismod/pull/380) Fix farm genesis validation.
* [\#367](https://github.com/irisnet/irismod/pull/367) Fix `mt` module rest url conflict.
* [\#356](https://github.com/irisnet/irismod/pull/356) Replace base64.StdEncoding with base64.RawStdEncoding.
* [\#351](https://github.com/irisnet/irismod/pull/351) Fix wrong addr length of the service module.

### API Breaking Changes

* [\#353](https://github.com/irisnet/irismod/pull/353) The commands of the token module only supports the main unit coin

## [v1.7.3]

### Improvements

* [\#335](https://github.com/irisnet/irismod/pull/335) Bump up cosmos-sdk to v0.46.9.
* [\#340](https://github.com/irisnet/irismod/pull/340) The token module supports the exchange of two tokens.
* [\#342](https://github.com/irisnet/irismod/pull/342) Refator token module.
* [\#348](https://github.com/irisnet/irismod/pull/348) Adjust the length limit of classID and nftID in nft.

### Bug Fixes

* [\#336](https://github.com/irisnet/irismod/pull/336) Fix farm genesis validate failed.
* [\#327](https://github.com/irisnet/irismod/pull/327) Only export htlc with state=open.
* [\#347](https://github.com/irisnet/irismod/pull/347) Fix service refund address parse error.
* [\#350](https://github.com/irisnet/irismod/pull/350) Fix address parse errors caused by service rest api conflicts.

## [v1.7.1] - 2022-11-18

### Improvements

* [\#321](https://github.com/irisnet/irismod/pull/321) Bump up cosmos-sdk to v0.46.5.

## [v1.7.0] - 2022-11-15

### Bug Fixes

* [\#304](https://github.com/irisnet/irismod/pull/304) Fix nft module import error.
* [\#314](https://github.com/irisnet/irismod/pull/314) Fix `addLiquidity` panic error.

### Improvements

* [\#305](https://github.com/irisnet/irismod/pull/305) Remove ibc-go from project.
* [\#319](https://github.com/irisnet/irismod/pull/319) Bump up cosmos-sdk to v0.46.4.
* [\#307](https://github.com/irisnet/irismod/pull/307) Refactor proto-gen with docker.
* [\#308](https://github.com/irisnet/irismod/pull/308) Coinswap module adds unilateral injection liquidity function.
* [\#309](https://github.com/irisnet/irismod/pull/309) Refactor nft with cosmos-sdk nft module.

### API Breaking Changes

* [\#309](https://github.com/irisnet/irismod/pull/309) GRPC method `Owner` rename to `NFTsOfOwner`, Remove deprecated `Queries` api

## [v1.6.0] - 2022-08-08

### Improvements

* [\#290](https://github.com/irisnet/irismod/pull/290) Bump cosmos-sdk version to [v0.45.1](https://github.com/cosmos/cosmos-sdk/releases/tag/v0.45.1).

### Bug Fixes

* [\#300](https://github.com/irisnet/irismod/pull/300) fix: pagination query limitation
* (modules/mt) [\#288](https://github.com/irisnet/irismod/pull/288), [\#286](https://github.com/irisnet/irismod/pull/286): Some fixes for the MT module

## [v1.5.2] - 2021-03-03

* (modules/farm) [\#247](https://github.com/irisnet/irismod/pull/247) Add farm proposal
* (modules/coinswap) [\#249](https://github.com/irisnet/irismod/pull/249) Add liquidity pool creation fee
* (modules/nft) [\#245](https://github.com/irisnet/irismod/pull/245) Improve nft module
* (modules/mt) [\#269](https://github.com/irisnet/irismod/pull/269) Feature: MT Module

### Improvements

* (modules/token) [\#252](https://github.com/irisnet/irismod/pull/252)  Precision of extension fee 9->18

## [v1.5.1] - 2022-02-15

* (modules/nft) [\#245](https://github.com/irisnet/irismod/pull/245) improve nft module

## [v1.5.0] - 2021-11-01

* [\#218](https://github.com/irisnet/irismod/pull/218) Bump cosmos-sdk version to [v0.44.2](https://github.com/cosmos/cosmos-sdk/releases/tag/v0.44.2).
* (modules/nft) [\#224](https://github.com/irisnet/irismod/pull/224) Remove the check of denomId in the query.
* (modules/coinswap) [\#221](https://github.com/irisnet/irismod/pull/221) Refactor coinswap module grpc APIs.
* (modules/coinswap) [\#219](https://github.com/irisnet/irismod/pull/219) Refactor coinswap module, the liquidity pool is named with lpt-{index}, and the index is incremented from 1.
* (modules/nft) [\#212](https://github.com/irisnet/irismod/pull/212) Make the order of commands parameters consistent.
* (modules/nft) [\#183](https://github.com/irisnet/irismod/issues/183) Enhance NFT module, make it compatible with ERC721 and restrict mint and update operations.
* [\#180](https://github.com/irisnet/irismod/issues/180) Add simulation tests.
* (modules/farm) [\#179](https://github.com/irisnet/irismod/issues/179) Add farm module.

## [v1.4.0] - 2021-03-26

### Bug Fixes

* (modules/coinswap) [\#161](https://github.com/irisnet/irismod/pull/161) Fix coinswap token validation.
* (modules/htlc) [\#157](https://github.com/irisnet/irismod/pull/157) Fix htlc params validation.
* (modules/htlc) [\#156](https://github.com/irisnet/irismod/pull/156) Fix validation when creating HTLT.
* (modules/coinswap) [\#155](https://github.com/irisnet/irismod/pull/155) Fix min liquidity check in add liquidity.
* (modules/coinswap)[\#153](https://github.com/irisnet/irismod/pull/151) Fix query not-existent reserve pool.
* (modules/service) [\#152](https://github.com/irisnet/irismod/pull/152) Fix update service binding.

### Improvements

* [\#167](https://github.com/irisnet/irismod/pull/167) Bump cosmos-sdk version to [v0.42.3](https://github.com/cosmos/cosmos-sdk/releases/tag/v0.42.3).
* (modules/htlc) [\#158](https://github.com/irisnet/irismod/pull/158) Improve htlc.
* (modules/token)[\#154](https://github.com/irisnet/irismod/pull/154) Add reserved token prefix.
* (modules/coinswap)[\#151](https://github.com/irisnet/irismod/pull/151) Replace prefix `swap/` with `swap`.
* (modules/htlc)[\#146](https://github.com/irisnet/irismod/pull/146) Refactor HTLC module to support HTLT.

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
[v1.4.0]: https://github.com/irisnet/irismod/releases/tag/v1.4.0
