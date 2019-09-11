---
order: 1
---

# General Concepts

## IRIShub Node Types

### Full Node

A full node is a fully functional peer in IRIShub network

- Voting power is 0
- Save a full backup of transaction history
- Could be upgrade to validator node

### Validator Node

- Validator node is a staking pool, which is responsible for signing votes to reach consensus, and verify/execute transactions in blocks. It could be seen as miners in PoW blockchains.
- Participate in on-chain governance, vote for proposals
- Upgrade software to latest version

### Validator Candidate Node

At the start of IRIShub, only top 100 bonded full node will become validator nodes, the rest will become candidates. The situation will change as delegation amount changes.

- Voting power is 0；
- No Block provision

## IRIShub User Types

### Validator

People that operator validator nodes. They must first stake some tokens with certain transaction and is responsible for maintain the validator nodes
and get rewards in return.

### Delegator

People that cannot, or do not want to run validator operations, can still participate in the staking process as delegators.

### Profiler

Profiler is a type of user that they are a special type of user who can submit software upgrade/halt proposals

### Trustee

Trustee is a type of user that  they are a special type of user who will receive funds from CommunityTaxUsage proposals

## IRIS Token

The IRIS hub has its own native token known as *IRIS*.  It is designed to serve three purposes in the network.

- **Staking.** Similar to the ATOM token in the Cosmos Hub, the IRIS token will be used as a staking token to secure the PoS blockchain.
- **Transaction Fee.** The IRIS token will also be used to pay fees for all transactions in the IRIS network.
- **Service Fee.** It is required that service providers in the IRIS network charge service fees denominated in the IRIS token.

It is intended that the IRIS network will eventually support all whitelisted fee tokens from the Cosmos network, which can be used to pay the transaction fees and service fees.

## Staking Rewards

Validator and its delegators can share the following rewards by portion：

- **Block Inflation**
  Block Inflation exists to incentivize IRIS holders to stake. As more IRIS tokens are staked, more secure the network become(Read more about [Staking](../features/stake.md)).
  Block Inflation will be [distributed every block](../features/mint.md). [Inflation rate](../features/mint.md) in IRISnet for the first year will be 4%.  **This ration could be adjusted by `parameter-change` proposals**.
  In this way, loose IRIS will devalue year by year.
- **Block Proposer Reward**
  In IRIShub, the probability for validators is proportional to its bonded tokens. If one's proposed block is finalized, it gets extra rewards for it.  
- **Fee**
  Each transaction needs a [fee](fee.md#fee) for compensating validators' work[Gas](fee.md#gas). These fees can be paid with IRIS and may later in any tokens which are whitelisted by the Hub’s governance. Fees are distributed to validators in proportion to their stake. A minimum fee/gas ration is set in IRISnet.

Each validator receives revenue in proportion to its total stake. However, before this revenue is distributed to its delegators, the validator can apply a commission for providing staking services.

### Staking Rewards Calculation Formula

The following formulas are based on the current [IRIShub Mainnet Params](gov-params.md).

#### Annual Rewards (ignore proposer rewards and fees)

- **AnnualInflation =** `Base * InflationRate` (aka 2 billion * 4% = 80 million iris)
- **ValidatorRewards =** `(AnnualInflation / BondedTokens) * (1 - CommunityTax) * (ValidatorSelfDelegation +  DelegatorsDelegation * ValidatorCommissionRate)`
- **DelegatorRewards =** `(AnnualInflation / BondedTokens) * (1 - CommunityTax) * DelegatorSelfDelegation * (1 - ValidatorCommissionRate)`

#### Block Rewards

- **BlockInflation =** `AnnualInflation / (365*24*60*12)` (aka 12.68 iris)
- **ProposerExtraRewards =** `(BaseProposerReward + BonusProposerReward * PrecommitPower/TotalVotingPower) * (BlockInflation + BlockCollectedFees)`
- **BlockRewards =** `(BlockInflation + BlockCollectedFees) * (1 - CommunityTax) - ProposerExtraRewards`
- **Proposer**
  - **ProposerTotalRewards =** `(BlockRewards / BondedTokens) * ValidatorBondedTokens + ProposerExtraRewards`
  - **Commission =** `ProposerTotalRewards * ValidatorCommissionRate`
  - **ProposerValidatorRewards =** `ProposerTotalRewards * (ValidatorSelfDelegation / ValidatorBondedTokens) + Commission`
  - **ProposerDelegatorRewards =** `(ProposerTotalRewards - Commission) * (DelegatorSelfDelegation / ValidatorBondedTokens)`
- **NonProposer**
  - **NonProposerTotalRewards =** `(BlockRewards / BondedTokens) * ValidatorBondedTokens`
  - **Commission =** `NonProposerTotalRewards * ValidatorCommissionRate`
  - **NonProposerValidatorRewards =** `NonProposerTotalRewards * (ValidatorSelfDelegation / ValidatorBondedTokens) + Commission`
  - **NonProposerDelegatorRewards =** `(NonProposerTotalRewards - Commission) * (DelegatorSelfDelegation / ValidatorBondedTokens)`

## Validator Responsibilities

Validators have two main responsibilities:

- **Be able to constantly run a correct version of the software:** Validators need to make sure that their servers are always online and their private keys are not compromised.
- **Actively participate in governance:** Validators are required to vote on every proposal.

Additionally, validators are expected to be active members of the community. They should always be up-to-date with the current state of the ecosystem so that they can easily adapt to any change.

## Validator Risks

- **Unavailability**: Validators are expected to keep signing votes for making new blocks. If a validator’s signature has not been included in more than 30% of the last 34,560 blocks (which amounts to approximately 24 hours, assuming an average block-generating time of 5 seconds), this validator will get jailed and removed from current validatorset for 1.5 day, and their bonded tokens will get slashed by 0.03%.
- **Double Sign**:If the protocol detects that a validator voted multiple different opinions about the same block (same height/round), or voted for different blocks at the same height/round, this validator will get jailed and removed from current validatorset for 2 days. Their bonded tokens will get slashed by 1%.
- **Censorship**: If the protocol detects that a proposer included invalid transactions in a block, this validator will get jailed and removed from current validatorset for 2 days.

All metrics mentioned can be adjusted by `parameter-change` proposals.
