<!--
order: 0
title: Service Overview
parent:
  title: "Service"
-->

# `service`

## Abstract

The service module implements the IRIS Service model, which bridges the gap between the blockchain world and the conventional application world. It formalizes off-chain service definition and binding (provider registration), facilitates invocation and interaction with those services, and mediates service governance process (profiling and dispute resolution).

## Contents

1. **[State](01_state.md)**
    - [ServiceDefinition](01_state.md#servicedefinition)
    - [ServiceBinding](01_state.md#servicebinding)
    - [ServiceInvocation](01_state.md#serviceinvocation)
    - [Params](01_state.md#parameters)
2. **[Messages](02_messages.md)**
    - [MsgDefineService](02_messages.md#MsgDefineService)
    - [MsgBindService](02_messages.md#MsgBindService)
    - [MsgUpdateServiceBinding](02_messages.md#MsgUpdateServiceBinding)
    - [MsgDiableServiceBinding](02_messages.md#MsgDiableServiceBinding)
    - [MsgEnableServiceBinding](02_messages.md#MsgEnableServiceBinding)
    - [MsgRefundServiceDeposit](02_messages.md#MsgRefundServiceDeposit)
    - [MsgSetWithdrawAddress](02_messages.md#MsgSetWithdrawAddress)
    - [MsgCallService](02_messages.md#MsgCallService)
    - [MsgRespondService](02_messages.md#MsgRespondService)
    - [MsgUpdateRequestContext](02_messages.md#MsgUpdateRequestContext)
    - [MsgPauseRequestContext](02_messages.md#MsgPauseRequestContext)
    - [MsgStartRequestContext](02_messages.md#MsgStartRequestContext)
    - [MsgKillRequestContext](02_messages.md#MsgKillRequestContext)
    - [MsgWithdrawEarnedFees](02_messages.md#MsgWithdrawEarnedFees)
3. **[Events](03_events.md)**
    - [EndBlocker](03_events.md#handlers)
    - [Handlers](03_events.md#handlers)
4. **[Parameters](04_params.md)**
