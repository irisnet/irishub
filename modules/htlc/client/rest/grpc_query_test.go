package rest_test

//
//import (
//	"fmt"
//	"testing"
//
//	"github.com/gogo/protobuf/proto"
//	"github.com/stretchr/testify/suite"
//
//	"github.com/cosmos/cosmos-sdk/client/flags"
//	"github.com/cosmos/cosmos-sdk/testutil/network"
//	sdk "github.com/cosmos/cosmos-sdk/types"
//
//	htlccli "github.com/irisnet/irismod/modules/htlc/client/cli"
//	htlctestutil "github.com/irisnet/irismod/modules/htlc/client/testutil"
//	htlctypes "github.com/irisnet/irismod/modules/htlc/types"
//	"github.com/irisnet/irismod/simapp"
//)
//
//type IntegrationTestSuite struct {
//	suite.Suite
//
//	cfg     network.Config
//	network *network.Network
//}
//
//func (s *IntegrationTestSuite) SetupSuite() {
//	s.T().Log("setting up integration test suite")
//
//	cfg := network.DefaultConfig()
//	cfg.AppConstructor = simapp.SimAppConstructor
//	cfg.NumValidators = 2
//
//	s.cfg = cfg
//	s.network = network.New(s.T(), cfg)
//
//	_, err := s.network.WaitForHeight(1)
//	s.Require().NoError(err)
//}
//
//func (s *IntegrationTestSuite) TearDownSuite() {
//	s.T().Log("tearing down integration test suite")
//	s.network.Cleanup()
//}
//
//func TestIntegrationTestSuite(t *testing.T) {
//	suite.Run(t, new(IntegrationTestSuite))
//}
//
//func (s *IntegrationTestSuite) TestNft() {
//	val := s.network.Validators[0]
//	val2 := s.network.Validators[1]
//	clientCtx := val.ClientCtx
//
//	// ---------------------------------------------------------------------------
//
//	from := val.Address
//	to := val2.Address
//	amount := "1000" + sdk.DefaultBondDenom
//	receiverOnOtherChain := "0xcd2a3d9f938e13cd947ec05abc7fe734df8dd826"
//	hashLock := "e8d4133e1a82c74e2746e78c19385706ea7958a0ca441a08dacfa10c48ce2561"
//	secretHex := "5f5f5f6162636465666768696a6b6c6d6e6f707172737475767778797a5f5f5f"
//	timeLock := uint64(50)
//	timestamp := uint64(1580000000)
//
//	args := []string{
//		fmt.Sprintf("--%s=%s", htlccli.FlagTo, to),
//		fmt.Sprintf("--%s=%s", htlccli.FlagAmount, amount),
//		fmt.Sprintf("--%s=%s", htlccli.FlagReceiverOnOtherChain, receiverOnOtherChain),
//		fmt.Sprintf("--%s=%s", htlccli.FlagHashLock, hashLock),
//		fmt.Sprintf("--%s=%s", htlccli.FlagSecret, secretHex),
//		fmt.Sprintf("--%s=%d", htlccli.FlagTimeLock, timeLock),
//		fmt.Sprintf("--%s=%d", htlccli.FlagTimestamp, timestamp),
//
//		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
//		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
//		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
//	}
//
//	respType := proto.Message(&sdk.TxResponse{})
//	expectedCode := uint32(0)
//
//	bz, err := htlctestutil.CreateHTLCExec(clientCtx, from.String(), args...)
//	s.Require().NoError(err)
//	s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType), bz.String())
//	txResp := respType.(*sdk.TxResponse)
//	s.Require().Equal(expectedCode, txResp.Code)
//
//	args = []string{
//		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
//		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
//		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
//	}
//
//	respType = proto.Message(&sdk.TxResponse{})
//
//	bz, err = htlctestutil.ClaimHTLCExec(clientCtx, from.String(), hashLock, secretHex, args...)
//	s.Require().NoError(err)
//	s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType), bz.String())
//	txResp = respType.(*sdk.TxResponse)
//	s.Require().Equal(expectedCode, txResp.Code)
//
//	args = []string{
//		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
//		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
//		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
//	}
//
//	respType = proto.Message(&sdk.TxResponse{})
//
//	bz, err = htlctestutil.RefundHTLCExec(clientCtx, from.String(), hashLock, args...)
//	s.Require().NoError(err)
//	s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType), bz.String())
//	txResp = respType.(*sdk.TxResponse)
//	s.Require().Equal(expectedCode, txResp.Code)
//
//	// ---------------------------------------------------------------------------
//
//	respType = proto.Message(&htlctypes.HTLC{})
//	bz, err = htlctestutil.QueryHTLCExec(val.ClientCtx, hashLock)
//	s.Require().NoError(err)
//	s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType))
//	htlcItem := respType.(*htlctypes.HTLC)
//	s.Require().Equal(amount, htlcItem.Amount)
//}
