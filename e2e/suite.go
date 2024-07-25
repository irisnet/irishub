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
	
}

// SetupSuiteWithModifyConfigFn sets up the end-to-end test suite with the given modifyConfigFn.
//
// Parameters:
// - modifyConfigFn: A function that modifies the config for the test suite.
//
// Return type: None.
func (s *TestSuite) SetupSuiteWithModifyConfigFn(modifyConfigFn ModifyConfigFn) {
	s.T().Log("setting up e2e test suite")

	cfg, err := simapp.NewConfig(s.DepinjectOptions())
	s.Require().NoError(err)

	modifyConfigFn(&cfg)
	s.Network = simapp.SetupNetworkWithConfig(s.T(), cfg)
}

// SetupSuite creates a new network for integration tests
func (s *TestSuite) SetupSuite() {
	s.T().Log("setting up e2e test suite")
	s.Network = simapp.SetupNetwork(s.T(), s.DepinjectOptions())
}

// TearDownSuite tears down the integration test suite
func (s *TestSuite) TearDownSuite() {
	s.T().Log("tearing down e2e test suite")
	s.Network.Cleanup()
}

// DepinjectOptions returns the depinject options for the test suite
func (s *TestSuite) DepinjectOptions() simapp.DepinjectOptions {
	return defaultDepinjectOptions()
}

func defaultDepinjectOptions() simapp.DepinjectOptions {
	return simapp.DepinjectOptions{
		Config: AppConfig,
		Providers: []interface{}{
			keeper.ProvideMockEVM(),
			keeper.ProvideMockICS20(),
		},
	}
}