<!--
order: 0
title: Farm Overview
parent:
  title: "Farm"
-->

# Farm Specification

## Abstract

This specification describes the creation, destruction, additional bonuses, stake `lptoken`, unstake `lptoken`, retrieving rewards functions of the farm pool

## Contents

1. **[State](01_state.md)**
   - [Params](01_state.md#params)
   - [FarmPool](01_state.md#farmPool)
   - [RewardRule](01_state.md#rewardRule)
   - [FarmInfo](01_state.md#farmInfo)
2. **[Messages](02_messages.md)**
   - [MsgCreatePool](02_messages.md#msgCreatePool)
   - [MsgDestroyPool](02_messages.md#msgDestroyPool)
   - [MsgAdjustPool](02_messages.md#msgAdjustPool)
   - [MsgStake](02_messages.md#msgStake)
   - [MsgUnstake](02_messages.md#msgUnstake)
   - [MsgHarvest](02_messages.md#msgHarvest)
3. **[Events](03_events.md)**
   - [Handlers](03_events.md#handlers)
4. **[Parameters](04_params.md)**
