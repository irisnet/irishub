package e2e

import (
	"github.com/stretchr/testify/suite"

	"mods.irisnet.org/simapp"
)

// TestSuite is a suite of end-to-end tests for the module
type TestSuite struct {
	suite.Suite
    simapp.Network
}

// SetupSuite creates a new network for integration tests
func (s *TestSuite) SetupSuite() {
	depInjectOptions := simapp.DepinjectOptions{
		Config:    AppConfig,
		Providers: []interface{}{
			ProvideEVMKeeper(),
			ProvideICS20Keeper(),
		},
	}

	s.T().Log("setting up e2e test suite")
	s.Network = simapp.SetupNetwork(s.T(),depInjectOptions)
}

// TearDownSuite tears down the integration test suite
func (s *TestSuite) TearDownSuite() {
	s.T().Log("tearing down e2e nft test suite")
	s.Network.Cleanup()
}