# Cosmos SDK v0.53.8 Upgrade

This branch uses the Go module path `github.com/irisnet/irishub/v5` and upgrades
to Cosmos SDK `v0.53.8` and IBC-Go `v10.7.0`. The minimum Go version is `1.23.8`.

## Pinned Dependencies

Dependencies are pinned to commits from the following branches using pseudo-versions,
so builds do not depend on moving branch heads:

| Dependency | Branch | Commit |
| --- | --- | --- |
| nft-transfer | `feat/cosmos-sdk-v0.53.8` | `87faf8d2b44a` |
| tibc-go | `feat/cosmos-sdk-v0.53.8` | `0b8e7b4af967` |
| mods.irisnet.org/modules/* | `feat/cosmos-sdk-v0.53.8` | `f4182d6e7472` |
| bianjieai/ethermint | `release/v0.24.0` | `4a4e53b81b98` |

IBC-Go uses `v10.7.0` directly, without a `replace` directive to downgrade it.
Following the replacement configuration in the Ethermint branch, the EVM
implementation remains on go-ethereum `v1.10.26`. Run
`go list -m -json github.com/cosmos/ibc-go/v10 github.com/ethereum/go-ethereum`
to inspect the resolved versions, replacement details, and source directories.

## On-chain Upgrade

The registered software upgrade plan is named **`v5.0.0`**. Existing chains must
execute this plan at the agreed upgrade height; replacing the binary and opening
the old database directly does not perform the required state migrations.

The upgrade loader deletes the obsolete `capability` KV store. The upgrade
handler removes its module-version entry and then runs the registered SDK and
IBC module migrations, including the IBC transfer migration from the old
DenomTrace encoding to the new Denom format. The existing transfer, NFT
transfer, and ICA host routes remain registered. Tendermint and solo-machine
clients are explicitly registered as required by IBC v10. Governance operations
such as client recovery and parameter updates use v10 messages, with the `gov`
module as the authority.

The `auth` PreBlocker is ordered after `upgrade`. Unordered transactions,
epochs, and protocolpool are not enabled.

LSM functionality is not retained. At the upgrade height, the handler removes
LSM records, indexes, and storage prefixes `0x81` through `0x87` from staking,
and rewrites staking parameters, validators, and delegation records to remove
the legacy LSM protobuf fields. Validators' minimum self-delegation values are
set to zero, preserving the behavior of the previous LSM version. The
underlying delegations for `tokenizeshare_N` accounts, their distribution
reward-owner associations, and LSM derivative coin balances and supplies in
bank are preserved for future redemption verification. This upgrade does not
provide a transaction for redeeming derivative coins.

## Docker Build

The Dockerfile uses `golang:1.24.9-alpine` and `alpine:3.18` by default. If
Docker Hub is unavailable, use the build arguments to switch to cached images:

```shell
docker build -t irishub:sdk0538 --build-arg EVM_CHAIN_ID=6688 --build-arg BUILDER_IMAGE=mirror.gcr.io/library/golang:1.24.9-alpine --build-arg RUNTIME_IMAGE=mirror.gcr.io/library/alpine:3.18 .
```

## Validation Scope

The added tests cover real application genesis initialization, first-block
execution, IBC routing, and migration of the legacy DenomTrace encoding through
the registered upgrade handler while preserving IBC asset denoms. A synthetic
LSM state test checks that the handler clears legacy staking fields and records
while preserving bank balances, underlying delegations, and reward-owner
associations. The existing repository tests check compatibility of local modules
and the CLI.

The upgrade rehearsal used the mainnet startup `genesis.json` as its base and
successfully upgraded the three tokenize-share records it contained. The file
was modified for the local single-validator rehearsal, including repair of
slashing signing information; this was not an unmodified production snapshot
rehearsal. Before a production upgrade, rehearse at the target height with a
complete target-chain snapshot, including existing IBC channels and clients and
the handling of remaining LSM assets. Unit tests alone cannot verify those
chain-specific states.

References: [SDK upgrade guide](https://github.com/cosmos/cosmos-sdk/blob/v0.53.8/UPGRADING.md),
[IBC v8.1 to v10 migration guide](https://docs.cosmos.network/ibc/next/migrations/v8_1-to-v10).
