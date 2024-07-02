package e2e

import (
	"github.com/cosmos/cosmos-sdk/testutil/network"
	"github.com/stretchr/testify/suite"

	"mods.irisnet.org/modules/token/keeper"
	"mods.irisnet.org/simapp"
)

// ModifyConfigFn is used to modify the config for a test suite
type ModifyConfigFn = func(cfg *network.Config)

// TestSuite is a suite of end-to-end tests for the module
type TestSuite struct {
	suite.Suite
	simapp.Network
	modifyConfigFn ModifyConfigFn
}

// SetupSuite creates a new network for integration tests
func (s *TestSuite) SetupSuite() {
	s.T().Log("setting up e2e test suite")

	depInjectOptions := simapp.DepinjectOptions{
		Config: AppConfig,
		Providers: []interface{}{
			keeper.ProvideMockEVM(),
			keeper.ProvideMockICS20(),
		},
	}
	if s.modifyConfigFn == nil {
		s.Network = simapp.SetupNetwork(s.T(), depInjectOptions)
		return
	}

	cfg, err := simapp.NewConfig(depInjectOptions)
	s.Require().NoError(err)

	s.modifyConfigFn(&cfg)
	s.Network = simapp.SetupNetworkWithConfig(s.T(), cfg)
}

// TearDownSuite tears down the integration test suite
func (s *TestSuite) TearDownSuite() {
	s.T().Log("tearing down e2e test suite")
	s.Network.Cleanup()
}

// SetModifyConfigFn sets the modify config function
func (s *TestSuite) SetModifyConfigFn(fn ModifyConfigFn) {
	s.modifyConfigFn = fn
}
