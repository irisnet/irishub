# IRIS Hub v5.0.0 Software Upgrade

### Background

IRIS Hub v5.0.0 is a state-breaking software upgrade that moves the application to the `github.com/irisnet/irishub/v5` module line and updates its core Cosmos and IBC dependencies. If this proposal passes, validators must install the designated v5.0.0 binary and restart their nodes at the upgrade height encoded in the on-chain proposal. Replacing the binary without executing the `v5.0.0` upgrade plan will not perform the required state migrations.

The upgrade modernizes IRIS Hub's application stack, retires the legacy Liquid Staking Module (LSM) implementation, and migrates IBC state to the formats required by IBC-Go v10. It is intended to preserve existing accounts, balances, normal staking positions, EVM state, and supported cross-chain routes while applying the migrations described below.

### Release Binary

IMPORTANT: Validators and node operators must use the official IRIS Hub v5.0.0 binary designated for the target network. The upgrade plan name is exactly `v5.0.0`; a binary with a different registered plan name will not execute this migration.

The v5.0.0 release will publish its production Linux artifacts under the [official GitHub release](https://github.com/irisnet/irishub/releases/tag/v5.0.0) using these canonical URLs:

- [`iris-v5.0.0-linux-amd64`](https://github.com/irisnet/irishub/releases/download/v5.0.0/iris-v5.0.0-linux-amd64)
- [`iris-v5.0.0-linux-arm64`](https://github.com/irisnet/irishub/releases/download/v5.0.0/iris-v5.0.0-linux-arm64)
- [`SHA256SUMS`](https://github.com/irisnet/irishub/releases/download/v5.0.0/SHA256SUMS)

Release maintainers will generate `SHA256SUMS` from the final builds. Operators must verify the downloaded binary against this file before installation and review the [v5 migration notes](https://github.com/irisnet/irishub/blob/master/docs/migration/cosmos-sdk-v0.53.8.md) and [pull request #3026](https://github.com/irisnet/irishub/pull/3026).

The EVM chain ID is embedded in the binary at build time through `EVM_CHAIN_ID`; setting that environment variable only when starting an already-built binary does not change it. Release builders must preserve the target network's existing EVM chain ID, and validators and EVM RPC operators must verify that they are running binaries built with the same value before the network resumes. Any binary built for another network must use a distinct asset name so that network-specific builds cannot be confused.

### Proposed Release Contents

The release upgrades the principal consensus and interoperability dependencies:

- Cosmos SDK to [`v0.53.8`](https://github.com/cosmos/cosmos-sdk/releases/tag/v0.53.8).
- IBC-Go to [`v10.7.0`](https://github.com/cosmos/ibc-go/releases/tag/v10.7.0).
- CometBFT to [`v0.38.23`](https://github.com/cometbft/cometbft/releases/tag/v0.38.23).

At the scheduled upgrade height, the registered handler performs the following state changes:

- Deletes the obsolete IBC `capability` KV store and removes `capability` from the persisted module version map before running module migrations.
- Runs the registered Cosmos SDK and IBC migrations, including the IBC transfer migration from the legacy `DenomTrace` representation to the v10 `Denom` representation while preserving existing IBC asset denominations.
- Retains the transfer, NFT transfer, and ICA host routes. IBC integrations must use the v10 message and authority model for operations such as client recovery and parameter updates.
- Removes legacy LSM records and indexes under staking prefixes `0x81` through `0x87`, and rewrites staking parameters, validator records, delegations, and historical staking records through the upstream codecs to remove legacy LSM-only protobuf fields.
- Sets validators' stored minimum self-delegation values to zero, preserving the behavior of the previous LSM implementation.
- Preserves the underlying delegations held by `tokenizeshare_N` module accounts, maps their reward ownership into the standard distribution withdrawal-address state, and preserves the corresponding derivative bank balances and supply.

LSM functionality is not available after this upgrade. Existing derivative balances remain backed by the preserved delegations, but v5.0.0 does not provide a transaction for redeeming those derivative coins. Any future redemption process requires separate technical review and governance agreement.

### Testing and Testnets

Repository regression tests cover application genesis initialization, first-block execution, IBC routing, migration of legacy IBC transfer denomination records, and synthetic legacy LSM state. The LSM test verifies that obsolete records and fields are removed while bank balances, supply, backing delegations, and reward-owner relationships remain intact.

The `v5.0.0` upgrade plan was applied on `nyancat-10` at height `30848000`. The network resumed block production, the `capability` module-version entry was removed, the IBC consensus version advanced from 6 to 8, and the transfer version advanced from 5 to 6. The EVM RPC chain ID was also verified as `16688` after operators restored the network's intended build-time value.

A separate single-validator rehearsal migrated three tokenized-share records using a modified mainnet-derived genesis. That exercise was not an unmodified production snapshot rehearsal. Before a production upgrade, operators should rehearse the exact release against a complete target-chain snapshot and verify existing IBC clients and channels, tokenized-share records, EVM behavior, and post-upgrade application queries. A signed EVM transaction and a complete post-upgrade IBC transfer handshake have not yet been established by the test evidence summarized here.

### Potential risk factors

This upgrade changes core dependencies and performs state migrations at the halt height. Migration time, database size, malformed legacy records, operator configuration differences, or an incompatible binary can delay network restart. The LSM cleanup is designed to fail closed if a tokenized-share record, backing delegation, owner, validator, or withdrawal address is inconsistent.

IBC-Go v10 interface and governance-message changes may require coordinated updates by relayers, wallets, exchanges, and other integrations. Holders of preserved LSM derivative balances will also remain unable to redeem those balances until a separately reviewed mechanism is adopted.

Operators should back up their validator and application data before the halt, verify the release checksum and compiled EVM chain ID, and remain available in the established IRIS Hub validator coordination channels. If a serious issue occurs, validators should stop and coordinate on a single recovery procedure. Once migrated state has been committed, independently replacing the binary with an older release is not a safe rollback procedure; recovery must use an agreed snapshot, height, and binary.

### Governance votes

The voting period, veto threshold, and deposit handling are determined by the active IRIS Hub governance parameters.

**YES** - You support upgrading IRIS Hub to v5.0.0 at the height specified in the on-chain proposal.

**NO** - You do not support this protocol upgrade.

**NO WITH VETO** - You believe the proposal is spam, materially harms minority interests, or violates IRIS Hub governance rules. Rejection thresholds and deposit burning follow the chain's active governance parameters.

**ABSTAIN** - You wish to count toward quorum without voting for or against the proposal.
