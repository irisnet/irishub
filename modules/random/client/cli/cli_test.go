package cli_test

// import (
// 	"context"
// 	"encoding/hex"
// 	"encoding/json"
// 	"fmt"
// 	"testing"
// 	"time"

// 	"github.com/stretchr/testify/suite"
// 	"github.com/tidwall/gjson"

// 	"github.com/cosmos/cosmos-sdk/client/flags"
// 	sdk "github.com/cosmos/cosmos-sdk/types"

// 	servicecli "github.com/irisnet/irismod/modules/service/client/cli"
// 	servicetestutil "github.com/irisnet/irismod/modules/service/client/testutil"
// 	servicetypes "github.com/irisnet/irismod/modules/service/types"
// 	"github.com/irisnet/irismod/simapp"
// 	randomcli "irismod.io/random/client/cli"
// 	randomtestutil "irismod.io/random/client/testutil"
// 	randomtypes "irismod.io/random/types"
// )

// type IntegrationTestSuite struct {
// 	suite.Suite

// 	network simapp.Network
// }

// func (s *IntegrationTestSuite) SetupSuite() {
// 	s.T().Log("setting up integration test suite")

// 	s.network = simapp.SetupNetwork(s.T())
// }

// func (s *IntegrationTestSuite) TearDownSuite() {
// 	s.T().Log("tearing down integration test suite")
// 	s.network.Cleanup()
// }

// func TestIntegrationTestSuite(t *testing.T) {
// 	suite.Run(t, new(IntegrationTestSuite))
// }

// func (s *IntegrationTestSuite) TestRandom() {
// 	val := s.network.Validators[0]
// 	clientCtx := val.ClientCtx
// 	expectedCode := uint32(0)

// 	// ---------------------------------------------------------------------------
// 	serviceDeposit := fmt.Sprintf("50000%s", s.network.BondDenom)
// 	servicePrices := fmt.Sprintf(`{"price": "50%s"}`, s.network.BondDenom)
// 	qos := int64(3)
// 	options := "{}"
// 	provider := val.Address

// 	from := val.Address
// 	blockInterval := 4
// 	oracle := true
// 	serviceFeeCap := fmt.Sprintf("50%s", s.network.BondDenom)

// 	respResult := `{"code":200,"message":""}`
// 	seedStr := "ABCDEF12ABCDEF12ABCDEF12ABCDEF12ABCDEF12ABCDEF12ABCDEF12ABCDEF12"
// 	respOutput := fmt.Sprintf(`{"header":{},"body":{"seed":"%s"}}`, seedStr)

// 	// ------bind random service-------------
// 	args := []string{
// 		fmt.Sprintf("--%s=%s", servicecli.FlagServiceName, randomtypes.ServiceName),
// 		fmt.Sprintf("--%s=%s", servicecli.FlagDeposit, serviceDeposit),
// 		fmt.Sprintf("--%s=%s", servicecli.FlagPricing, servicePrices),
// 		fmt.Sprintf("--%s=%d", servicecli.FlagQoS, qos),
// 		fmt.Sprintf("--%s=%s", servicecli.FlagOptions, options),
// 		fmt.Sprintf("--%s=%s", servicecli.FlagProvider, provider),

// 		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
// 		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
// 		fmt.Sprintf(
// 			"--%s=%s",
// 			flags.FlagFees,
// 			sdk.NewCoins(sdk.NewCoin(s.network.BondDenom, sdk.NewInt(10))).String(),
// 		),
// 	}

// 	txResult := servicetestutil.BindServiceExec(
// 		s.T(),
// 		s.network,
// 		clientCtx,
// 		provider.String(),
// 		args...)
// 	s.Require().Equal(expectedCode, txResult.Code)

// 	// ------test GetCmdRequestRandom()-------------
// 	args = []string{
// 		fmt.Sprintf("--%s=%s", randomcli.FlagServiceFeeCap, serviceFeeCap),
// 		fmt.Sprintf("--%s=%t", randomcli.FlagOracle, oracle),
// 		fmt.Sprintf("--%s=%d", randomcli.FlagBlockInterval, blockInterval),

// 		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
// 		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
// 		fmt.Sprintf(
// 			"--%s=%s",
// 			flags.FlagFees,
// 			sdk.NewCoins(sdk.NewCoin(s.network.BondDenom, sdk.NewInt(10))).String(),
// 		),
// 	}

// 	txResult = randomtestutil.RequestRandomExec(s.T(), s.network, clientCtx, from.String(), args...)
// 	s.Require().Equal(expectedCode, txResult.Code)

// 	requestID := gjson.Get(txResult.Log, "0.events.1.attributes.0.value").String()
// 	requestHeight := gjson.Get(txResult.Log, "0.events.1.attributes.2.value").Int()

// 	// ------test GetCmdQueryRandomRequestQueue()-------------
// 	qrrResp := randomtestutil.QueryRandomRequestQueueExec(
// 		s.T(),
// 		s.network,
// 		clientCtx,
// 		fmt.Sprintf("%d", requestHeight),
// 	)
// 	s.Require().Len(qrrResp.Requests, 1)

// 	// ------get service request-------------
// 	requestHeight = requestHeight + 1
// 	_, err := s.network.WaitForHeightWithTimeout(
// 		requestHeight,
// 		time.Duration(int64(blockInterval+5)*int64(s.network.TimeoutCommit)),
// 	)
// 	s.Require().NoError(err)

// 	blockResult, err := val.RPCClient.BlockResults(context.Background(), &requestHeight)
// 	s.Require().NoError(err)
// 	var requestId string
// 	for _, event := range blockResult.EndBlockEvents {
// 		if event.Type == servicetypes.EventTypeNewBatchRequestProvider {
// 			var found bool
// 			var requestIds []string
// 			var requestsBz []byte
// 			for _, attribute := range event.Attributes {
// 				if string(attribute.Key) == servicetypes.AttributeKeyRequests {
// 					requestsBz = []byte(attribute.GetValue())
// 					found = true
// 				}
// 			}
// 			s.Require().True(found)
// 			if found {
// 				err := json.Unmarshal(requestsBz, &requestIds)
// 				s.Require().NoError(err)
// 			}
// 			s.Require().Len(requestIds, 1)
// 			requestId = requestIds[0]
// 		}
// 	}
// 	s.Require().NotNil(requestId)

// 	// ------respond service request-------------
// 	args = []string{
// 		fmt.Sprintf("--%s=%s", servicecli.FlagRequestID, requestId),
// 		fmt.Sprintf("--%s=%s", servicecli.FlagResult, respResult),
// 		fmt.Sprintf("--%s=%s", servicecli.FlagData, respOutput),

// 		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
// 		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
// 		fmt.Sprintf(
// 			"--%s=%s",
// 			flags.FlagFees,
// 			sdk.NewCoins(sdk.NewCoin(s.network.BondDenom, sdk.NewInt(10))).String(),
// 		),
// 	}

// 	txResult = servicetestutil.RespondServiceExec(
// 		s.T(),
// 		s.network,
// 		clientCtx,
// 		provider.String(),
// 		args...)
// 	s.Require().Equal(expectedCode, txResult.Code)

// 	generateHeight := txResult.Height

// 	// ------test GetCmdQueryRandom()-------------
// 	randomResp := randomtestutil.QueryRandomExec(s.T(), s.network, clientCtx, requestID)
// 	s.Require().NotNil(randomResp.Value)

// 	generateBLock, err := clientCtx.Client.Block(context.Background(), &generateHeight)
// 	s.Require().NoError(err)
// 	seed, err := hex.DecodeString(seedStr)
// 	s.Require().NoError(err)
// 	random := randomtypes.MakePRNG(generateBLock.Block.LastBlockID.Hash, generateBLock.Block.Header.Time.Unix(), from, seed, true).
// 		GetRand().
// 		FloatString(randomtypes.RandPrec)
// 	s.Require().Equal(random, randomResp.Value)
// }
