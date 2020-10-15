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

### Bug Fixes

* (modules/service) [\#26](https://github.com/irisnet/irismod/pull/26) Fix deduct service fees and optimized code.
* (modules/coinswap) [\#25](https://github.com/irisnet/irismod/pull/25) Integrate htlc beginblock.
* (modules/coinswap) [\#19](https://github.com/irisnet/irismod/pull/19) Check liquidity reserve by supply.
* (modules/service) [\#18](https://github.com/irisnet/irismod/pull/18) Fix incorrect price in request and fix init request when insufficient banlances.

## [v1.1.0] - 2020-09-30

### Features

* Add modules `token`, `nft`, `htlc`, `coinswap`, `service`, `oracle`, `random`, `record`.

<!-- Release links -->

[v1.1.0]: https://github.com/irisnet/irismod/releases/tag/v1.1.0
