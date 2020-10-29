package keeper

import (
	tmbytes "github.com/tendermint/tendermint/libs/bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/irisnet/irismod/modules/service/types"
)

// RegisterModuleService registers a module service
func (k Keeper) RegisterModuleService(moduleName string, moduleService *types.ModuleService) error {
	if _, ok := k.moduleServices[moduleName]; ok {
		return sdkerrors.Wrapf(types.ErrModuleServiceRegistered, "%s already registered for module %s", "module service", moduleName)
	}

	k.moduleServices[moduleName] = moduleService

	return nil
}

func (k Keeper) GetModuleServiceByModuleName(moduleName string) (*types.ModuleService, bool) {
	if k.moduleServices[moduleName] == nil {
		return &types.ModuleService{}, false
	}
	return k.moduleServices[moduleName], true
}

func (k Keeper) GetModuleServiceByServiceName(serviceName string) (string, *types.ModuleService, bool) {
	for moduleName, mdouleSvc := range k.moduleServices {
		if mdouleSvc.ServiceName == serviceName {
			return moduleName, mdouleSvc, true
		}
	}
	return "", nil, false
}

func (k Keeper) RequestModuleService(
	ctx sdk.Context,
	moduleService *types.ModuleService,
	reqContextID tmbytes.HexBytes,
	consumer sdk.AccAddress,
	input string,
) error {
	requestContext, found := k.GetRequestContext(ctx, reqContextID)
	if !found {
		return sdkerrors.Wrap(types.ErrUnknownRequestContext, reqContextID.String())
	}

	requestContextConsumer, _ := sdk.AccAddressFromBech32(requestContext.Consumer)

	pds := make([]sdk.AccAddress, len(requestContext.Providers))
	for i, provider := range requestContext.Providers {
		pd, err := sdk.AccAddressFromBech32(provider)
		if err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid provider address: %s", provider)
		}
		pds[i] = pd
	}

	_, totalPrices, _, err := k.FilterServiceProviders(
		ctx,
		requestContext.ServiceName,
		pds,
		requestContext.Timeout,
		requestContext.ServiceFeeCap,
		requestContextConsumer,
	)
	if err != nil {
		return err
	}

	if err := k.DeductServiceFees(ctx, consumer, totalPrices); err != nil {
		return err
	}

	requestIDs := k.InitiateRequests(ctx, reqContextID, []sdk.AccAddress{moduleService.Provider}, make(map[string][]string))

	result, output := moduleService.ReuquestService(ctx, input)
	request, _, err := k.AddResponse(ctx, requestIDs[0], moduleService.Provider, result, output)
	if err != nil {
		return err
	}

	requestContext.State = types.COMPLETED
	k.SetRequestContext(ctx, reqContextID, requestContext)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, moduleService.Provider.String()),
			sdk.NewAttribute(types.AttributeKeyRequestContextID, request.RequestContextId),
			sdk.NewAttribute(types.AttributeKeyRequestID, requestIDs[0].String()),
			sdk.NewAttribute(types.AttributeKeyServiceName, request.ServiceName),
			sdk.NewAttribute(types.AttributeKeyConsumer, request.Consumer),
		),
	})

	return nil
}
