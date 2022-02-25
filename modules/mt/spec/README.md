<!--
order: 0
title: MT Overview
parent:
  title: "MT"
-->

# MT Specification

## Overview

A structure that manage multiple token types. A single deployed contract may include any combination of fungible tokens, non-fungible tokens or other configurations (e.g. semi-fungible tokens).

This standard was first developed on Ethereum within the ERC-1155 and the subsequent EIP of the same name.The ERC-1155 standard addressed some of the restrictions of Ethereum regarding storage costs and semi-fungible assets.

## Contents

1. **[State](./01_state.md)**
   - [MT](./01_state.md#MT)
   - [Collection](./01_state.md#Collection)
   - [Balance](./01_state.md#Balance)
1. **[Messages](./02_messages.md)**
   - [Issue Denom](./02_messages.md#MsgIssueDenom)
   - [Transfer Denom](./02_messages.md#MsgTransferDenom)
   - [Mint MT](./02_messages.md#MsgMintMT)
   - [Edit MT](./02_messages.md#MsgEditMT)
   - [Transfer MT](./02_messages.md#MsgTransferMT)
   - [Burn MT](./02_messages.md#MsgBurnMT)
1. **[Events](./03_events.md)**
   - [Handlers](03_events.md#Handlers)
1. **[Future Improvements](./04_future_improvements.md)**
