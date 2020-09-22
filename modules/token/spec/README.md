<!--
order: 0
title: Token Overview
parent:
  title: "Token"
-->

# `token`

## Abstract

This specification describes the management of Fungible Token on the chain, Anyone could issue a new token on on the chain, or propose pegging an existing token from any other blockchains via On-Chain Governance.

## Contents

1. **[State](01_state.md)**
    - [Token](01_state.md#token)
    - [Params](01_state.md#params)
2. **[Messages](02_messages.md)**
    - [MsgIssueToken](02_messages.md#msgIssueToken)
    - [MsgEditToken](02_messages.md#msgEditToken)
    - [MsgMintToken](02_messages.md#msgMintToken)
    - [MsgTransferTokenOwner](02_messages.md#msgTransferTokenOwner)
    - [MsgBeginRedelegate](02_messages.md#msgbeginredelegate)
3. **[Events](03_events.md)**
    - [Handlers](03_events.md#handlers)
4. **[Parameters](04_params.md)**
