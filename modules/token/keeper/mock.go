package keeper

import (
	"context"
	"fmt"
	"math/big"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"

	"mods.irisnet.org/modules/token/contracts"
	tokentypes "mods.irisnet.org/modules/token/types"
)

var (
	_ tokentypes.EVMKeeper   = (*mockEVM)(nil)
	_ tokentypes.ICS20Keeper = (*mockICS20)(nil)
)

// ProvideMockEVM returns an instance of tokentypes.EVMKeeper.
//
// No parameters.
// Returns a tokentypes.EVMKeeper.
func ProvideMockEVM() tokentypes.EVMKeeper {
	return &mockEVM{
		erc20s: make(map[common.Address]*erc20),
	}
}

// ProvideMockICS20 returns an instance of tokentypes.ICS20Keeper.
//
// No parameters.
// Returns a tokentypes.ICS20Keeper.
func ProvideMockICS20() tokentypes.ICS20Keeper {
	return &mockICS20{}
}

type mockEVM struct {
	erc20s map[common.Address]*erc20
}

// ApplyMessage implements types.EVMKeeper.
func (e *mockEVM) ApplyMessage(ctx sdk.Context, msg core.Message, tracer vm.EVMLogger, commit bool) (*tokentypes.Result, error) {
	isCreate := msg.To() == nil
	if isCreate {
		contractAddr := crypto.CreateAddress(msg.From(), msg.Nonce())

		data := msg.Data()[len(contracts.TokenProxyContract.Bin):]
		args, err := contracts.TokenProxyContract.ABI.Constructor.Inputs.Unpack(data)
		if err != nil {
			return nil, err
		}

		data = args[1].([]byte)
		data = data[4:]
		args, err = contracts.ERC20TokenContract.ABI.Methods[contracts.MethodInitialize].Inputs.Unpack(data)
		if err != nil {
			return nil, err
		}

		name, _ := args[0].(string)
		symbol, _ := args[1].(string)
		scale, _ := args[2].(uint8)
		e.erc20s[contractAddr] = &erc20{
			address: contractAddr,
			scale:   scale,
			name:    name,
			symbol:  symbol,
			balance: make(map[common.Address]*big.Int),
		}
		return &tokentypes.Result{
			Hash: contractAddr.Hex(),
		}, nil
	}

	erc20Contract, ok := e.erc20s[*msg.To()]
	if !ok {
		return nil, fmt.Errorf("erc20 contract not found")
	}
	return e.dispatch(erc20Contract, msg.Data())
}

// ChainID implements types.EVMKeeper.
func (e *mockEVM) ChainID() *big.Int {
	return big.NewInt(16688)
}

// EstimateGas implements types.EVMKeeper.
func (e *mockEVM) EstimateGas(ctx context.Context, req *tokentypes.EthCallRequest) (uint64, error) {
	return 3000000, nil
}

// SupportedKey implements types.EVMKeeper.
func (e *mockEVM) SupportedKey(pubKey cryptotypes.PubKey) bool {
	return true
}

func (e *mockEVM) dispatch(contract *erc20, data []byte) (*tokentypes.Result, error) {
	method, err := contracts.ERC20TokenContract.ABI.MethodById(data[0:4])
	if err != nil {
		return nil, err
	}

	ret, err := contract.call(method, data[4:])
	if err != nil {
		return nil, err
	}
	return &tokentypes.Result{
		Hash: contract.address.Hex(),
		Ret:  ret,
	}, nil
}

type erc20 struct {
	address      common.Address
	scale        uint8
	name, symbol string

	balance map[common.Address]*big.Int
}

func (erc20 erc20) call(method *abi.Method, data []byte) ([]byte, error) {
	args, err := method.Inputs.Unpack(data)
	if err != nil {
		return nil, err
	}

	switch method.Name {
	case "name":
		return method.Outputs.Pack(erc20.name)
	case "symbol":
		return method.Outputs.Pack(erc20.symbol)
	case "decimals":
		return method.Outputs.Pack(erc20.scale)
	case "balanceOf":
		balance, ok := erc20.balance[args[0].(common.Address)]
		if !ok {
			return method.Outputs.Pack(big.NewInt(0))
		}
		return method.Outputs.Pack(balance)
	case "mint":
		to := args[0].(common.Address)
		balance, ok := erc20.balance[args[0].(common.Address)]
		if !ok {
			balance = big.NewInt(0)
		}
		erc20.balance[to] = new(big.Int).Add(balance, args[1].(*big.Int))
		return nil, nil
	case "burn":
		from := args[0].(common.Address)
		erc20.balance[from] = new(big.Int).Sub(erc20.balance[from], args[1].(*big.Int))
		return nil, nil
	default:
		return nil, fmt.Errorf("unknown method %s", method.Name)
	}
}

type mockICS20 struct{}

// HasTrace implements types.ICS20Keeper.
func (t *mockICS20) HasTrace(ctx sdk.Context, denom string) bool {
	return true
}
