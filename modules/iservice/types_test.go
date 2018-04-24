package iservice

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/irisnet/iris-hub/modules/iservice/bind"
	"github.com/irisnet/iris-hub/modules/iservice/def"
	"github.com/irisnet/iris-hub/types"
	"reflect"
	"testing"
)

func TestNewIServiceRoute(t *testing.T) {
	storeKey := sdk.NewKVStoreKey("iris")

	accountMapper := auth.NewAccountMapperSealed(
		sdk.NewKVStoreKey("bank"), // target store
		&types.AppAccount{},       // prototype
	)

	// add handlers
	coinKeeper := bank.NewCoinKeeper(accountMapper)

	mapper := NewIServiceRoute(storeKey, coinKeeper)

	msg := buildMsg()
	service := mapper[reflect.TypeOf(msg).String()]
	fmt.Printf("%s", reflect.TypeOf(msg).Name())

	result := service.CheckTx(sdk.Context{}, msg)

	fmt.Printf("%v+", result)

	msg = buildMsg2()

	service = mapper[reflect.TypeOf(msg).String()]

	result = service.CheckTx(sdk.Context{}, msg)

	fmt.Printf("%v+", result)
}

func buildMsg() sdk.Msg {
	return def.SvcDefMsg{}
}

func buildMsg2() sdk.Msg {
	return bind.SvcBindMsg{}
}
