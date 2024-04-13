package types

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
)

const (
	// DefaultGasCap is the default gas cap for eth_call
	DefaultGasCap uint64 = 25000000
)

// TransactionArgs represents the arguments to construct a new transaction
// or a message call using JSON-RPC.
type TransactionArgs struct {
	From                 *common.Address `json:"from"`
	To                   *common.Address `json:"to"`
	Gas                  *hexutil.Uint64 `json:"gas"`
	GasPrice             *hexutil.Big    `json:"gasPrice"`
	MaxFeePerGas         *hexutil.Big    `json:"maxFeePerGas"`
	MaxPriorityFeePerGas *hexutil.Big    `json:"maxPriorityFeePerGas"`
	Value                *hexutil.Big    `json:"value"`
	Nonce                *hexutil.Uint64 `json:"nonce"`

	// We accept "data" and "input" for backwards-compatibility reasons.
	// "input" is the newer name and should be preferred by clients.
	Data  *hexutil.Bytes `json:"data"`
	Input *hexutil.Bytes `json:"input"`

	// Introduced by AccessListTxType transaction.
	AccessList *ethtypes.AccessList `json:"accessList,omitempty"`
	ChainID    *hexutil.Big         `json:"chainId,omitempty"`
}

// EthCallRequest represents the arguments to the eth_call RPC
type EthCallRequest struct {
	// args uses the same json format as the json rpc api.
	Args []byte `json:"args,omitempty"`
	// gas_cap defines the default gas cap to be used
	GasCap uint64 `json:"gas_cap,omitempty"`
	// proposer_address of the requested block in hex format
	ProposerAddress sdk.ConsAddress `json:"proposer_address,omitempty"`
	// chain_id is the eip155 chain id parsed from the requested block header
	ChainID int64 `json:"chain_id,omitempty"`
}

// Result represents the result of a contract execution
type Result struct {
	// hash of the ethereum transaction in hex format. This hash differs from the
	// Tendermint sha256 hash of the transaction bytes. See
	// https://github.com/tendermint/tendermint/issues/6539 for reference
	Hash string
	// logs contains the transaction hash and the proto-compatible ethereum
	// logs.
	Logs []*ethtypes.Log
	// ret is the returned data from evm function (result or data supplied with revert
	// opcode)
	Ret []byte
	// vm_error is the error returned by vm execution
	VMError string
	// gas_used specifies how much gas was consumed by the transaction
	GasUsed uint64
}

// Failed returns if the contract execution failed in vm errors
func (r *Result) Failed() bool {
	return len(r.VMError) > 0
}

// Return is a helper function to help caller distinguish between revert reason
// and function return. Return returns the data after execution if no error occurs.
func (r *Result) Return() []byte {
	if r.Failed() {
		return nil
	}
	return common.CopyBytes(r.Ret)
}

// Revert returns the concrete revert reason if the execution is aborted by `REVERT`
// opcode. Note the reason can be nil if no data supplied with revert opcode.
func (r *Result) Revert() []byte {
	if r.VMError != vm.ErrExecutionReverted.Error() {
		return nil
	}
	return common.CopyBytes(r.Ret)
}

var _ vm.EVMLogger = &NoOpTracer{}

// NoOpTracer is an empty implementation of vm.Tracer interface
type NoOpTracer struct{}

// NewNoOpTracer creates a no-op vm.Tracer
func NewNoOpTracer() *NoOpTracer {
	return &NoOpTracer{}
}

// CaptureStart implements vm.Tracer interface
func (dt NoOpTracer) CaptureStart(env *vm.EVM,
	from common.Address,
	to common.Address,
	create bool,
	input []byte,
	gas uint64,
	value *big.Int) {
}

// CaptureState implements vm.Tracer interface
func (dt NoOpTracer) CaptureState(pc uint64, op vm.OpCode, gas, cost uint64, scope *vm.ScopeContext, rData []byte, depth int, err error) {
}

// CaptureFault implements vm.Tracer interface
func (dt NoOpTracer) CaptureFault(pc uint64, op vm.OpCode, gas, cost uint64, scope *vm.ScopeContext, depth int, err error) {
}

// CaptureEnd implements vm.Tracer interface
func (dt NoOpTracer) CaptureEnd(output []byte, gasUsed uint64, tm time.Duration, err error) {}

// CaptureEnter implements vm.Tracer interface
func (dt NoOpTracer) CaptureEnter(typ vm.OpCode, from common.Address, to common.Address, input []byte, gas uint64, value *big.Int) {
}

// CaptureExit implements vm.Tracer interface
func (dt NoOpTracer) CaptureExit(output []byte, gasUsed uint64, err error) {}

// CaptureTxStart implements vm.Tracer interface
func (dt NoOpTracer) CaptureTxStart(gasLimit uint64) {}

// CaptureTxEnd implements vm.Tracer interface
func (dt NoOpTracer) CaptureTxEnd(restGas uint64) {}

// HexString is a byte array that serializes to hex
type HexString []byte

// MarshalJSON serializes ByteArray to hex
func (s HexString) MarshalJSON() ([]byte, error) {
	return json.Marshal(fmt.Sprintf("%x", string(s)))
}

// UnmarshalJSON deserializes ByteArray to hex
func (s *HexString) UnmarshalJSON(data []byte) error {
	var x string
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	str, err := hex.DecodeString(x)
	if err != nil {
		return err
	}
	*s = str
	return nil
}

// CompiledContract contains compiled bytecode and abi
type CompiledContract struct {
	ABI abi.ABI
	Bin HexString 
}

