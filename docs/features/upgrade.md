# Software Upgrade

## Concepts

### Plan

The `Upgrade` module defines a `Plan` type in which a live upgrade is scheduled to occur. A `Plan` can be scheduled at a specific block height or time, but not both.
A `Plan` is created once a (frozen) release candidate along with an appropriate upgrade `Handler` (see below) is agreed upon, where the `Name` of a `Plan` corresponds to a
specific `Handler`. Typically, a `Plan` is created through a governance proposal process, where if voted upon and passed, will be scheduled. The `Info` of a `Plan`
may contain various metadata about the upgrade, typically application specific upgrade info to be included on-chain such as a git commit that validators could automatically upgrade to.

#### Sidecar Process

If an operator running the application binary also runs a sidecar process to assist in the automatic download and upgrade of the binary, `Info` allows this process to be frictionless. That is, the upgrade module implements the specification specified by [cosmovisor upgradeable binary specification](https://github.com/cosmos/cosmos-sdk/tree/master/cosmovisor#upgradeable-binary-specification), and `cosmovisor` can Site selection is used to fully automate the upgrade process for node operators. By filling the `Info` field with the necessary information, the binary can be downloaded automatically. Refer to [here](https://github.com/cosmos/cosmos-sdk/tree/master/cosmovisor#auto-download).

```go
type Plan struct {
 Name string
 Time time.Time
 Height int64
 Info string
 UpgradedClientState *types.Any
}
```

### Handler

The `Upgrade` module facilitates upgrading from major version X to major version Y. To accomplish this, node operators must first upgrade their current binary to a new binary that has a corresponding `Handler` for the new version Y. It is assumed that this version has fully been tested and approved by the community at large. This `Handler` defines what state migrations need to occur before the new binary Y can successfully run the chain. Naturally, this `Handler` is application specific and not defined on a per-module basis. Registering a `Handler` is done via `Keeper#SetUpgradeHandler` in the application.

```go
type UpgradeHandler func(Context, Plan)
```

During the execution of each `BeginBlock`, the `Upgrade` module checks whether there is a `plan` that should be executed (the block height or time when the `BeginBlock` is scheduled to run). If it exists, execute the corresponding `processor`. If the plan is expected to be executed but the corresponding processor is not registered, or the binary upgrade is too early, the node will gracefully panic and exit.

### StoreLoader

The `Upgrade` module also facilitates store migrations as part of the upgrade. The `StoreLoader` sets the migrations that need to occur before the new binary can successfully run the chain. This `StoreLoader` is also application specific and not defined on a per-module basis. Registering this `StoreLoader` is done via `app#SetStoreLoader` in the application.

```go
func UpgradeStoreLoader (upgradeHeight int64, storeUpgrades *store.StoreUpgrades) baseapp.StoreLoader
```

If there's a planned upgrade and the upgrade height is reached, the old binary writes `UpgradeInfo` to the disk before panic'ing.

```go
type UpgradeInfo struct {
    Name    string
    Height  int64
}
```

This information is critical to ensure the `StoreUpgrades` happens smoothly at correct height and expected upgrade. It eliminates the chances for the new binary to execute `StoreUpgrades` multiple times everytime on restart. Also if there are multiple upgrades planned on same height, the `Name` will ensure these `StoreUpgrades` takes place only in planned upgrade handler.

Currently in the upgrade process, the state transition methods involved mainly support the following three types: `Renamed`, `Deleted`, and `Added`

#### Renamed

The user can specify to migrate all data under oldKey (prefix) to storage under newKey (prefix) during the upgrade process.

#### Deleted

Users can delete all data under the specified key (prefix) during the upgrade process.

#### Added

Users can apply for a new storage area with the specified key as the prefix during the upgrade process.

### Proposal

Typically, a `Plan` is proposed and submitted through governance via a `SoftwareUpgradeProposal`. This proposal prescribes to the standard governance process. If the proposal passes,
the `Plan`, which targets a specific `Handler`, is persisted and scheduled. The upgrade can be delayed or hastened by updating the `Plan.Time` in a new proposal.

```go
type SoftwareUpgradeProposal struct {
    Title       string
    Description string
    Plan        Plan
}
```

#### Governance Process

When a upgrade proposal is accepted, the upgrade process is devided into two steps.

##### Stop network consensus

After the software upgrade proposal is accepted, the system will perform pre-upgrade preparations at the designated height of the BeginBlock stage, including downloading the upgrade plan and suspending network consensus.

###### Download the upgrade plan

In order to be able to upgrade the software smoothly, the information required for the upgrade must be recorded before stopping the network consensus: `plan name`, `upgrade height`.
  
-`Plan name`: When the network restarts, it needs to be routed to the corresponding `UpgradeHandler` and `UpgradeStoreLoader` according to the plan name.
-`Upgrade height`: When the network restarts, check whether the network upgrade plan is required.

###### Suspend network consensus

After the software upgrade proposal is accepted, the system will gracefully suspend the network consensus at the specified height of the `BeginBlock` stage.

##### Restart the new software

The user replaces the software with the specified version and restarts the network. The system will detect whether the `processor` specified by the `plan name` is included. If it is included, the system first executes the `Handler` program, and then starts the network consensus. If not, the system report an error and exit.

#### Cancelling Upgrade Proposals

The upgrade proposal can be cancelled. There is a proposal type of `Cancel Software Upgrade (CancelSoftwareUpgrade)`. When this type of proposal is voted through, the currently ongoing upgrade plan will be removed. Of course, this needs to be voted and executed before the upgrade plan is executed.

If the current upgrade plan has been executed, but there are problems with the upgrade plan, then the proposal of `Cancel Software Upgrade` is invalid (because the network has stopped consensus). At this time, there is another solution to make up for this mistake, which is to use the `--unsafe-skip-upgrades` parameter to skip the specified upgrade height when restarting the network (not really skip the height, but jump via software upgrade `Handler`). Of course, this requires that 2/3 of the validators participating in the consensus perform the same operation, otherwise the network consensus cannot be reached.

## Upgrade process

### Submit an upgrade proposal

The first step in the implementation of the software upgrade process is to initiate a software upgrade proposal by the governance module. The proposal details the upgrade time and upgrade content. For details, see the above [Concept](#Concepts). An example of a command line to initiate a proposal is as follows:

```bash
iris tx gov submit-proposal software-upgrade bifrost-rc2 \
  --deposit 10000000uiris \
  --upgrade-time 2021-02-09T13:00:00Z \
  --title "mainnet software upgrade" \
  --upgrade-info "Commit: 0ef5dd0b4d140a4788f05fc1a0bd409b3c6a0492. After the proposal is approved, please use the commit hash to build and restart your node." \
  --description "Upgrade the mainnet software version from v1.0.0-rc0 to v1.0.0-rc2."
  --from=node0 --chain-id=test 1000000uiris -b block -y
```

### Deposit and vote for the proposal

The execution process of the software upgrade proposal is basically the same as that of other ordinary proposals. Both validators and delegators are required to comment on the proposal. For specific information, please refer to [governance module](./governance.md). An example of the command line to deposit the proposal is as follows:

```bash
iris tx gov deposit 1 1000000000uiris --from=node0 --chain-id=test --fees=1000000uiris -b block -y
```

Once the deposit amount reaches the minimum deposit amount, the proposal will enter the voting period, and the validator or delegator needs to vote on the proposal. An example of the command line to initiate a vote is as follows:

```bash
iris tx gov vote 1 yes  --from=node0 --chain-id=test 1000000uiris -b block -y
```

When the software upgrade proposal is passed, the upgrade module will create an upgrade plan to stop all nodes from the network consensus at a specified height or time, and wait for the new software to restart the network.

### Restart the network

When the upgrade proposal is passed, the node will stop producing blocks, and the user needs to download the source code and compile the new software according to the new version information specified in the first step [Submit Upgrade Proposal](#submit-an-upgrade-proposal), refer to [Install](../get-started/install.md). After the new software is installed, restart the node with the new version, and the node will execute the upgrade logic corresponding to the plan name. Once the voting power of the entire network exceeds 2/3 and restarts the network using the new version, the blockchain network will re-reach a new consensus and continue to produce new blocks.
