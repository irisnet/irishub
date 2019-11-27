# Governance

## Basic Function Description

1. On-chain governance proposals on plain text
2. On-chain governance proposals on parameter change
3. On-chain governance proposals on software upgrade
4. On-chain governance proposals on software halt
5. On-chain governance proposals on tax usage
6. On-chain governance proposals on token addition

## Interactive process

### Proposal Level

Specific Proposal for different levels:

- Critical：`SoftwareUpgrade`, `SystemHalt`
- Important：`Parameter`,`TokenAddition`
- Normal：`CommunityTaxUsage`,`PlainText`

`SoftwareUpgrade Proposal` and `SystemHalt Proposal` can only be submitted by the profiler.

Different levels correspond to different parameters：

| GovParams     | Critical  | Important | Normal      | Range                  |
| ------------- | --------- | --------- | ----------- | ---------------------- |
| MinDeposit    | 4000 iris | 2000 iris | 1000 iris   | [10iris,10000iris]     |
| DepositPeriod | 24 hours  | 24 hours  | 24 hours    | [20s,3d]               |
| VotingPeriod  | 120 hours | 120 hours | 120 hours   | [20s,7d]               |
| MaxNum        | 1         | 5         | 7           | Critical==1, other(1,) |
| Participation | 0.5       | 0.5       | 0.5 |(0,1)  |                        |
| Threshold     | 0.75      | 0.67      | 0.5 |(0,1)  |                        |
| Veto          | 0.33      | 0.33      | 0.33 |(0,1) |                        |
| Penalty       | 0         | 0         | 0 |(0,1)    |                        |

- `MinDeposit`  The minimum of  deposit
- `DepositPeriod`  Window period for deposit
- `VotingPeriod` Window period for voting
- `MaxNum` The maximum number of proposal that can exist at the same time
- `Penalty`   The proportion of the slash
- `Veto`  the power of Veto / all voted power
- `Threshold` the power of Yes / all voted power
- `Participation` all voted power / total voting power

### Deposit Procedure

The proposer at least deposit more the 30% amount of `MinDeposit` to submit a proposal, when the total deposit amount exceeds `MinDeposit`, the proposal enter the voting procedure. If the time exceeds `MaxDepositPeriod` and the total deposit has not yet exceeded `MinDeposit`, the proposal will be deleted and the full deposit won't be refunded. It is not allowed to deposit a proposal which is in voting procedure.

### Voting Procedure

Only the validator and delegator can vote , and they can't vote twice for one proposal. The voting options are `Yes` , `Abstain` , `No` , `NoWithVeto` .

### Tallying Procedure

There are three tallying results: `PASS`, `REJECT`, `REJECTVETO`.

On the premise that the `voting_power of all voters` / `total voting_power of the system` exceeds `participation`, if the ratio of `NoWithVeto` voting power to all voters' voting power over `veto`, the result is `REJECTVETO`. Then if the ratio of `Yes` voting power to all voter's voting power over `threshold`, the result is `PASS`. Otherwise, the result is `REJECT`.

### Burning Mechanism

Whether the proposal is passed or not, 20% `Deposit` will be burned for the cost of governance. The remaining `Deposit` will be returned. But if the result of proposal is `REJECTVETO`,  all `Deposit` will be burned.

### Slashing Mechanism

The validator should be slashed according to the proportion of `Penalty` if he fails to vote for a proposal.

If the validator quit from validator set in voting procedure, then he wont be slashed.

## Usage Scenario

### Usage scenario of parameter change

Change the parameters through the command lines

```bash
# Query parameters can be changed by the modules'name in gov
iriscli params --module=mint

# Result
Mint Params:
  mint/Inflation=0.0400000000

# Send proposal for parameters change
iriscli gov submit-proposal --title=<title> --description=<description> --type=Parameter --deposit=8iris  --param="mint/Inflation=0.0000000000" --from=<key-name> --chain-id=<chain-id> --fee=0.3iris --commit

# Deposit for a proposal
iriscli gov deposit --proposal-id=<proposal-id> --deposit=1000iris --from=<key-name> --chain-id=<chain-id> --fee=0.3iris --commit

# Vote for a proposal
iriscli gov vote --proposal-id=<proposal-id> --option=Yes --from=<key-name> --chain-id=<chain-id> --fee=0.3iris --commit

# Query the state of a proposal
iriscli gov query-proposal --proposal-id=<proposal-id>
```

### Proposals on community funds usage

There are three usages, `Burn`, `Distribute` and `Grant`. `Burn` means burning tokens from community funds. `Distribute` and `Grant` will transfer tokens to the destination trustee's account from community funds.

```bash
# Submit Burn usage proposal
iriscli gov submit-proposal --title="burn tokens 5%" --description=<description> --type="CommunityTaxUsage" --usage="Burn" --deposit="10iris"  --percent=0.05 --from=<key-name> --chain-id=<chain-id> --fee=0.3iris --commit

# Submit Distribute usage proposal
iriscli gov submit-proposal --title="distribute tokens 5%" --description="test" --type="CommunityTaxUsage" --usage="Distribute" --deposit="10iris"  --percent=0.05 --dest-address=<dest-address (only trustees)> --from=<key-name> --chain-id=<chain-id> --fee=0.3iris --commit

# Submit Grant usage proposal
iriscli gov submit-proposal --title="grant tokens 5%" --description="test" --type="CommunityTaxUsage" --usage="Grant" --deposit="10iris"  --percent=0.05 --dest-address=<dest-address (only trustees)> --from=<key-name> --chain-id=<chain-id> --fee=0.3iris --commit
```

### Proposals on system halting

Sending this proposal which can terminate the system, the node will be closed after systemHaltHeight (= proposal height + systemHaltPeriod), and only `query-only` mode is available after re-starting the node.

```bash
# submit the SystemHaltProposal
iriscli gov submit-proposal --title=<title> --description=<description> --type=SystemHalt --deposit=10iris --fee=0.3iris --from=<key-name> --chain-id=<chain-id> --commit
```

### Proposals on software upgrade

Detail in [Upgrade](upgrade.md)
