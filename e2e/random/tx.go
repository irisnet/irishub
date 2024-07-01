package random

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tidwall/gjson"

	"mods.irisnet.org/e2e"
	"mods.irisnet.org/e2e/service"
	randomcli "mods.irisnet.org/modules/random/client/cli"
	randomtypes "mods.irisnet.org/modules/random/types"
	servicecli "mods.irisnet.org/modules/service/client/cli"
	servicetypes "mods.irisnet.org/modules/service/types"
)

// TxTestSuite is a suite of end-to-end tests for the nft module
type TxTestSuite struct {
	e2e.TestSuite
}

// SetupSuite sets up test suite
func (s *TxTestSuite) SetupSuite() {
	s.SetModifyConfigFn(func(cfg *network.Config) {
		var serviceGenState servicetypes.GenesisState
		cfg.Codec.MustUnmarshalJSON(cfg.GenesisState[servicetypes.ModuleName], &serviceGenState)

		serviceGenState.Definitions = append(
			serviceGenState.Definitions,
			servicetypes.GenOraclePriceSvcDefinition(),
			servicetypes.GetRandomSvcDefinition(),
		)
		serviceGenState.Bindings = append(
			serviceGenState.Bindings,
			servicetypes.GenOraclePriceSvcBinding(sdk.DefaultBondDenom),
		)
		cfg.GenesisState[servicetypes.ModuleName] = cfg.Codec.MustMarshalJSON(&serviceGenState)
	})
	s.TestSuite.SetupSuite()
}

// TestTxCmd tests all tx command in the nft module
func (s *TxTestSuite) TestTxCmd() {
	val := s.Validators[0]
	clientCtx := val.ClientCtx
	expectedCode := uint32(0)

	// ---------------------------------------------------------------------------
	serviceDeposit := fmt.Sprintf("50000%s", s.BondDenom)
	servicePrices := fmt.Sprintf(`{"price": "50%s"}`, s.BondDenom)
	qos := int64(3)
	options := "{}"
	provider := val.Address

	from := val.Address
	blockInterval := 4
	oracle := true
	serviceFeeCap := fmt.Sprintf("50%s", s.BondDenom)

	respResult := `{"code":200,"message":""}`
	seedStr := "ABCDEF12ABCDEF12ABCDEF12ABCDEF12ABCDEF12ABCDEF12ABCDEF12ABCDEF12"
	respOutput := fmt.Sprintf(`{"header":{},"body":{"seed":"%s"}}`, seedStr)

	// ------bind random service-------------
	args := []string{
		fmt.Sprintf("--%s=%s", servicecli.FlagServiceName, randomtypes.ServiceName),
		fmt.Sprintf("--%s=%s", servicecli.FlagDeposit, serviceDeposit),
		fmt.Sprintf("--%s=%s", servicecli.FlagPricing, servicePrices),
		fmt.Sprintf("--%s=%d", servicecli.FlagQoS, qos),
		fmt.Sprintf("--%s=%s", servicecli.FlagOptions, options),
		fmt.Sprintf("--%s=%s", servicecli.FlagProvider, provider),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf(
			"--%s=%s",
			flags.FlagFees,
			sdk.NewCoins(sdk.NewCoin(s.BondDenom, sdk.NewInt(10))).String(),
		),
	}

	txResult := service.BindServiceExec(
		s.T(),
		s.Network,
		clientCtx,
		provider.String(),
		args...)
	s.Require().Equal(expectedCode, txResult.Code, txResult.Log)

	// ------test GetCmdRequestRandom()-------------
	args = []string{
		fmt.Sprintf("--%s=%s", randomcli.FlagServiceFeeCap, serviceFeeCap),
		fmt.Sprintf("--%s=%t", randomcli.FlagOracle, oracle),
		fmt.Sprintf("--%s=%d", randomcli.FlagBlockInterval, blockInterval),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf(
			"--%s=%s",
			flags.FlagFees,
			sdk.NewCoins(sdk.NewCoin(s.BondDenom, sdk.NewInt(10))).String(),
		),
	}

	txResult = RequestRandomExec(s.T(), s.Network, clientCtx, from.String(), args...)
	s.Require().Equal(expectedCode, txResult.Code)

	requestID := gjson.Get(txResult.Log, "0.events.1.attributes.0.value").String()
	requestHeight := gjson.Get(txResult.Log, "0.events.1.attributes.2.value").Int()

	// ------test GetCmdQueryRandomRequestQueue()-------------
	qrrResp := QueryRandomRequestQueueExec(
		s.T(),
		s.Network,
		clientCtx,
		fmt.Sprintf("%d", requestHeight),
	)
	s.Require().Len(qrrResp.Requests, 1)

	// ------get service request-------------
	requestHeight ++
	_, err := s.WaitForHeightWithTimeout(
		requestHeight,
		time.Duration(int64(blockInterval+5)*int64(s.TimeoutCommit)),
	)
	s.Require().NoError(err)

	blockResult, err := val.RPCClient.BlockResults(context.Background(), &requestHeight)
	s.Require().NoError(err)
	var requestId string
	for _, event := range blockResult.EndBlockEvents {
		if event.Type == servicetypes.EventTypeNewBatchRequestProvider {
			var found bool
			var requestIds []string
			var requestsBz []byte
			for _, attribute := range event.Attributes {
				if attribute.Key == servicetypes.AttributeKeyRequests {
					requestsBz = []byte(attribute.GetValue())
					found = true
				}
			}
			s.Require().True(found)
			if found {
				err := json.Unmarshal(requestsBz, &requestIds)
				s.Require().NoError(err)
			}
			s.Require().Len(requestIds, 1)
			requestId = requestIds[0]
		}
	}
	s.Require().NotNil(requestId)

	// ------respond service request-------------
	args = []string{
		fmt.Sprintf("--%s=%s", servicecli.FlagRequestID, requestId),
		fmt.Sprintf("--%s=%s", servicecli.FlagResult, respResult),
		fmt.Sprintf("--%s=%s", servicecli.FlagData, respOutput),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf(
			"--%s=%s",
			flags.FlagFees,
			sdk.NewCoins(sdk.NewCoin(s.BondDenom, sdk.NewInt(10))).String(),
		),
	}

	txResult = service.RespondServiceExec(
		s.T(),
		s.Network,
		clientCtx,
		provider.String(),
		args...)
	s.Require().Equal(expectedCode, txResult.Code)

	generateHeight := txResult.Height

	// ------test GetCmdQueryRandom()-------------
	randomResp := QueryRandomExec(s.T(), s.Network, clientCtx, requestID)
	s.Require().NotNil(randomResp.Value)

	generateBLock, err := clientCtx.Client.Block(context.Background(), &generateHeight)
	s.Require().NoError(err)
	seed, err := hex.DecodeString(seedStr)
	s.Require().NoError(err)
	random := randomtypes.MakePRNG(generateBLock.Block.LastBlockID.Hash, generateBLock.Block.Header.Time.Unix(), from, seed, true).
		GetRand().
		FloatString(randomtypes.RandPrec)
	s.Require().Equal(random, randomResp.Value)
}
