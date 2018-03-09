package iservice

import (
	"fmt"
	"strconv"

	"github.com/tendermint/tmlibs/log"

	"github.com/cosmos/cosmos-sdk"
	"github.com/cosmos/cosmos-sdk/errors"
	"github.com/cosmos/cosmos-sdk/modules/auth"
	"github.com/cosmos/cosmos-sdk/stack"
	"github.com/cosmos/cosmos-sdk/state"
)

// nolint
const iserviceModuleName = "iservice"

// Name is the name of the modules.
func Name() string {
	return iserviceModuleName
}

//_______________________________________________________________________

// Handler - the transaction processing handler
type Handler struct {
	stack.PassInitValidate
}

var _ stack.Dispatchable = Handler{} // enforce interface at compile time

// NewHandler returns a new Handler with the default Params
func NewHandler() Handler {
	return Handler{}
}

// Name - return iservice namespace
func (Handler) Name() string {
	return iserviceModuleName
}

// AssertDispatcher - placeholder for stack.Dispatchable
func (Handler) AssertDispatcher() {}

// InitState - set genesis parameters for iservice
func (h Handler) InitState(l log.Logger, store state.SimpleDB,
	module, key, value string, cb sdk.InitStater) (log string, err error) {
	return "", h.initState(module, key, value, store)
}

// separated for testing
func (Handler) initState(module, key, value string, store state.SimpleDB) error {
	if module != iserviceModuleName {
		return errors.ErrUnknownModule(module)
	}

	params := loadParams(store)
	switch key {
	case "gas_define_service":

		// TODO: enforce non-negative integers in input
		i, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("input must be integer, Error: %v", err.Error())
		}
		params.GasDefineService = int64(i)
	default:
		return errors.ErrUnknownKey(key)
	}

	saveParams(store, params)
	return nil
}

// CheckTx checks if the tx is properly structured
func (h Handler) CheckTx(ctx sdk.Context, store state.SimpleDB,
	tx sdk.Tx, _ sdk.Checker) (res sdk.CheckResult, err error) {

	err = tx.ValidateBasic()
	if err != nil {
		return res, err
	}

	// get the sender
	sender, err := getTxSender(ctx)
	if err != nil {
		return res, err
	}

	params := loadParams(store)

	// create the new checker object to
	checker := check{
		store:  store,
		sender: sender,
	}

	// return the fee for each tx type
	switch txInner := tx.Unwrap().(type) {
	case TxDefineService:
		return sdk.NewCheck(params.GasDefineService, ""),
			checker.checkService(txInner)
	}

	return res, errors.ErrUnknownTxType(tx)
}

// DeliverTx executes the tx if valid
func (h Handler) DeliverTx(ctx sdk.Context, store state.SimpleDB,
	tx sdk.Tx, dispatch sdk.Deliver) (res sdk.DeliverResult, err error) {

	// TODO: remove redundancy
	// also we don't need to check the res - gas is already deducted in sdk
	_, err = h.CheckTx(ctx, store, tx, nil)
	if err != nil {
		return
	}

	sender, err := getTxSender(ctx)
	if err != nil {
		return
	}

	params := loadParams(store)
	deliverer := deliver{
		store:  store,
		sender: sender,
		params: params,
	}

	// Run the transaction
	switch _tx := tx.Unwrap().(type) {
	case TxDefineService:
		res.GasUsed = params.GasDefineService
		return res, deliverer.defineService(_tx)
	}
	return
}

// get the sender from the ctx and ensure it matches the tx pubkey
func getTxSender(ctx sdk.Context) (sender sdk.Actor, err error) {
	senders := ctx.GetPermissions("", auth.NameSigs)
	if len(senders) != 1 {
		return sender, ErrMissingSignature()
	}
	return senders[0], nil
}

//_______________________________________________________________________

type check struct {
	store  state.SimpleDB
	sender sdk.Actor
}

func (c check) checkService(tx TxDefineService) error {
	if c.store.Has([]byte(tx.Name)) {
		return ErrServiceExists()
	}
	return nil
}

type deliver struct {
	store  state.SimpleDB
	sender sdk.Actor
	params Params
}

func (d deliver) defineService(tx TxDefineService) error {
	saveService(d.store, &tx)
	return nil
}

func (d deliver) getService(name string) *ServiceDefinition {

	return loadService(d.store, name)
}
