<!--
order: 0
title: Token Overview
parent:
  title: "Token"
-->

# Token Specification

## Abstract

This specification describes the management of Fungible Token on the chain, Anyone could issue a new token on on the chain, or propose pegging an existing token from any other blockchains via On-Chain Governance.

## Contents

1. **[State](01_state.md)**
   - [Token](01_state.md#token)
   - [Params](01_state.md#params)
1. **[Messages](02_messages.md)**
   - [MsgIssueToken](02_messages.md#msgissuetoken)
   - [MsgEditToken](02_messages.md#msgedittoken)
   - [MsgMintToken](02_messages.md#msgminttoken)
   - [MsgBurnToken](02_messages.md#msgburntoken)
   - [MsgTransferTokenOwner](02_messages.md#msgtransfertokenowner)
1. **[Events](03_events.md)**
   - [Handlers](03_events.md#handlers)
1. **[Parameters](04_params.md)**
