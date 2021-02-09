<!--
order: 0
title: Service Overview
parent:
title: "Service"
-->

# `service`

## Abstract

The service module implements the IRIS Service model, which bridges the
gap between the blockchain world and the conventional application world.
It formalizes off-chain service definition and binding (provider
registration), facilitates invocation and interaction with those
services, and mediates service governance process (profiling and dispute
resolution).

## Contents

1. **[State](01_state.md)**
   - [ServiceDefinition](01_state.md#servicedefinition)
   - [ServiceBinding](01_state.md#servicebinding)
   - [ServiceInvocation](01_state.md#serviceinvocation)
   - [Params](01_state.md#parameters)
2. **[Messages](02_messages.md)**
   - [MsgDefineService](02_messages.md#msgdefineservice)
   - [MsgBindService](02_messages.md#msgbindservice)
   - [MsgUpdateServiceBinding](02_messages.md#msgupdateservicebinding)
   - [MsgDiableServiceBinding](02_messages.md#msgdisableservicebinding)
   - [MsgEnableServiceBinding](02_messages.md#msgenableservicebinding)
   - [MsgRefundServiceDeposit](02_messages.md#msgrefundservicedeposit)
   - [MsgSetWithdrawAddress](02_messages.md#msgsetwithdrawaddress)
   - [MsgCallService](02_messages.md#msgcallservice)
   - [MsgRespondService](02_messages.md#msgrespondservice)
   - [MsgUpdateRequestContext](02_messages.md#msgupdaterequestcontext)
   - [MsgPauseRequestContext](02_messages.md#msgpauserequestcontext)
   - [MsgStartRequestContext](02_messages.md#msgstartrequestcontext)
   - [MsgKillRequestContext](02_messages.md#msgkillrequestcontext)
   - [MsgWithdrawEarnedFees](02_messages.md#msgwithdrawearnedfees)
3. **[Events](03_events.md)**
   - [Handlers](03_events.md#handlers)
4. **[Parameters](04_params.md)**

