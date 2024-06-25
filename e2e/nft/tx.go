package nft

import (
	"fmt"

	"github.com/cometbft/cometbft/crypto"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"mods.irisnet.org/e2e"
	nftcli "mods.irisnet.org/modules/nft/client/cli"
	"mods.irisnet.org/simapp"
)

// TxTestSuite is a suite of end-to-end tests for the nft module
type TxTestSuite struct {
	suite.Suite

	network simapp.Network
}

// SetupSuite creates a new network for integration tests
func (s *TxTestSuite) SetupSuite() {
	depInjectOptions := simapp.DepinjectOptions{
		Config:    e2e.AppConfig,
		Providers: []interface{}{
			e2e.ProvideEVMKeeper(),
			e2e.ProvideICS20Keeper(),
		},
	}

	s.T().Log("setting up e2e test suite")
	s.network = simapp.SetupNetwork(s.T(),depInjectOptions)
}

// TearDownSuite tears down the integration test suite
func (s *TxTestSuite) TearDownSuite() {
	s.T().Log("tearing down e2e nft test suite")
	s.network.Cleanup()
}


// TestTxCmd tests all tx command in the nft module
func (s *TxTestSuite) TestTxCmd() {
	val := s.network.Validators[0]
	val2 := s.network.Validators[1]
	clientCtx := val.ClientCtx
	expectedCode := uint32(0)

	// ---------------------------------------------------------------------------

	from := val.Address
	tokenName := "Kitty Token"
	uri := "uri"
	uriHash := "uriHash"
	description := "description"
	data := "{\"key1\":\"value1\",\"key2\":\"value2\"}"
	tokenID := "kitty"
	//owner     := "owner"
	denomName := "name"
	denomID := "denom"
	schema := "schema"
	symbol := "symbol"
	mintRestricted := true
	updateRestricted := false

	//------test GetCmdIssueDenom()-------------
	args := []string{
		fmt.Sprintf("--%s=%s", nftcli.FlagDenomName, denomName),
		fmt.Sprintf("--%s=%s", nftcli.FlagSchema, schema),
		fmt.Sprintf("--%s=%s", nftcli.FlagSymbol, symbol),
		fmt.Sprintf("--%s=%s", nftcli.FlagURI, uri),
		fmt.Sprintf("--%s=%s", nftcli.FlagURIHash, uriHash),
		fmt.Sprintf("--%s=%s", nftcli.FlagDescription, description),
		fmt.Sprintf("--%s=%s", nftcli.FlagData, data),
		fmt.Sprintf("--%s=%t", nftcli.FlagMintRestricted, mintRestricted),
		fmt.Sprintf("--%s=%t", nftcli.FlagUpdateRestricted, updateRestricted),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf(
			"--%s=%s",
			flags.FlagFees,
			sdk.NewCoins(sdk.NewCoin(s.network.BondDenom, sdk.NewInt(10))).String(),
		),
	}

	txResult := IssueDenomExec(s.T(),
		s.network,
		clientCtx, from.String(), denomID, args...)
	s.Require().Equal(expectedCode, txResult.Code)

	//------test GetCmdQueryDenom()-------------
	queryDenomResponse := QueryDenomExec(s.T(), s.network, clientCtx, denomID)
	s.Require().Equal(denomName, queryDenomResponse.Name)
	s.Require().Equal(schema, queryDenomResponse.Schema)
	s.Require().Equal(symbol, queryDenomResponse.Symbol)
	s.Require().Equal(uri, queryDenomResponse.Uri)
	s.Require().Equal(uriHash, queryDenomResponse.UriHash)
	s.Require().Equal(description, queryDenomResponse.Description)
	s.Require().Equal(data, queryDenomResponse.Data)
	s.Require().Equal(mintRestricted, queryDenomResponse.MintRestricted)
	s.Require().Equal(updateRestricted, queryDenomResponse.UpdateRestricted)

	//------test GetCmdQueryDenoms()-------------
	queryDenomsResponse := QueryDenomsExec(s.T(), s.network, clientCtx)
	s.Require().Equal(1, len(queryDenomsResponse.Denoms))
	s.Require().Equal(denomID, queryDenomsResponse.Denoms[0].Id)

	//------test GetCmdMintNFT()-------------
	args = []string{
		fmt.Sprintf("--%s=%s", nftcli.FlagData, data),
		fmt.Sprintf("--%s=%s", nftcli.FlagRecipient, from.String()),
		fmt.Sprintf("--%s=%s", nftcli.FlagURI, uri),
		fmt.Sprintf("--%s=%s", nftcli.FlagURIHash, uriHash),
		fmt.Sprintf("--%s=%s", nftcli.FlagTokenName, tokenName),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf(
			"--%s=%s",
			flags.FlagFees,
			sdk.NewCoins(sdk.NewCoin(s.network.BondDenom, sdk.NewInt(10))).String(),
		),
	}

	txResult = MintNFTExec(s.T(),
		s.network,
		clientCtx, from.String(), denomID, tokenID, args...)
	s.Require().Equal(expectedCode, txResult.Code)

	//------test GetCmdQuerySupply()-------------
	querySupplyResponse := QuerySupplyExec(s.T(), s.network, clientCtx, denomID)
	s.Require().Equal(uint64(1), querySupplyResponse.Amount)

	//------test GetCmdQueryNFT()-------------
	queryNFTResponse := QueryNFTExec(s.T(), s.network, clientCtx, denomID, tokenID)
	s.Require().Equal(tokenID, queryNFTResponse.Id)
	s.Require().Equal(tokenName, queryNFTResponse.Name)
	s.Require().Equal(uri, queryNFTResponse.URI)
	s.Require().Equal(uriHash, queryNFTResponse.UriHash)
	s.Require().Equal(data, queryNFTResponse.Data)
	s.Require().Equal(from.String(), queryNFTResponse.Owner)

	//------test GetCmdQueryOwner()-------------
	queryNFTsOfOwnerResponse := QueryOwnerExec(
		s.T(),
		s.network,
		clientCtx,
		from.String(),
	)
	s.Require().Equal(from.String(), queryNFTsOfOwnerResponse.Owner.Address)
	s.Require().Equal(denomID, queryNFTsOfOwnerResponse.Owner.IDCollections[0].DenomId)
	s.Require().Equal(tokenID, queryNFTsOfOwnerResponse.Owner.IDCollections[0].TokenIds[0])

	//------test GetCmdQueryCollection()-------------
	queryCollectionResponse := QueryCollectionExec(s.T(), s.network, clientCtx, denomID)
	s.Require().Equal(1, len(queryCollectionResponse.Collection.NFTs))

	//------test GetCmdEditNFT()-------------
	newTokenData := "{\"key1\":\"value1\",\"key2\":\"value2\"}"
	newTokenURI := "newuri"
	newTokenURIHash := "newuriHash"
	newTokenName := "new Kitty Token"
	args = []string{
		fmt.Sprintf("--%s=%s", nftcli.FlagData, newTokenData),
		fmt.Sprintf("--%s=%s", nftcli.FlagURI, newTokenURI),
		fmt.Sprintf("--%s=%s", nftcli.FlagURIHash, newTokenURIHash),
		fmt.Sprintf("--%s=%s", nftcli.FlagTokenName, newTokenName),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf(
			"--%s=%s",
			flags.FlagFees,
			sdk.NewCoins(sdk.NewCoin(s.network.BondDenom, sdk.NewInt(10))).String(),
		),
	}

	txResult = EditNFTExec(s.T(),
		s.network,
		clientCtx, from.String(), denomID, tokenID, args...)
	s.Require().Equal(expectedCode, txResult.Code)

	queryNFTResponse = QueryNFTExec(s.T(), s.network, clientCtx, denomID, tokenID)
	s.Require().Equal(newTokenName, queryNFTResponse.Name)
	s.Require().Equal(newTokenURI, queryNFTResponse.URI)
	s.Require().Equal(newTokenURIHash, queryNFTResponse.UriHash)
	s.Require().Equal(newTokenData, queryNFTResponse.Data)

	//------test GetCmdTransferNFT()-------------
	recipient := sdk.AccAddress(crypto.AddressHash([]byte("dgsbl")))

	args = []string{
		fmt.Sprintf("--%s=%s", nftcli.FlagData, data),
		fmt.Sprintf("--%s=%s", nftcli.FlagURI, uri),
		fmt.Sprintf("--%s=%s", nftcli.FlagURIHash, uriHash),
		fmt.Sprintf("--%s=%s", nftcli.FlagTokenName, tokenName),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf(
			"--%s=%s",
			flags.FlagFees,
			sdk.NewCoins(sdk.NewCoin(s.network.BondDenom, sdk.NewInt(10))).String(),
		),
	}

	txResult = TransferNFTExec(s.T(),
		s.network,
		clientCtx, from.String(), recipient.String(), denomID, tokenID, args...)
	s.Require().Equal(expectedCode, txResult.Code)

	queryNFTResponse = QueryNFTExec(s.T(), s.network, clientCtx, denomID, tokenID)
	s.Require().Equal(tokenID, queryNFTResponse.Id)
	s.Require().Equal(tokenName, queryNFTResponse.Name)
	s.Require().Equal(uri, queryNFTResponse.URI)
	s.Require().Equal(uriHash, queryNFTResponse.UriHash)
	s.Require().Equal(data, queryNFTResponse.Data)
	s.Require().Equal(recipient.String(), queryNFTResponse.Owner)

	//------test GetCmdBurnNFT()-------------
	newTokenID := "dgsbl"
	args = []string{
		fmt.Sprintf("--%s=%s", nftcli.FlagData, newTokenData),
		fmt.Sprintf("--%s=%s", nftcli.FlagRecipient, from.String()),
		fmt.Sprintf("--%s=%s", nftcli.FlagURI, newTokenURI),
		fmt.Sprintf("--%s=%s", nftcli.FlagTokenName, newTokenName),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf(
			"--%s=%s",
			flags.FlagFees,
			sdk.NewCoins(sdk.NewCoin(s.network.BondDenom, sdk.NewInt(10))).String(),
		),
	}

	txResult = MintNFTExec(s.T(),
		s.network,
		clientCtx, from.String(), denomID, newTokenID, args...)
	s.Require().Equal(expectedCode, txResult.Code)

	querySupplyResponse = QuerySupplyExec(s.T(), s.network, clientCtx, denomID)
	s.Require().Equal(uint64(2), querySupplyResponse.Amount)

	args = []string{
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf(
			"--%s=%s",
			flags.FlagFees,
			sdk.NewCoins(sdk.NewCoin(s.network.BondDenom, sdk.NewInt(10))).String(),
		),
	}
	txResult = BurnNFTExec(s.T(),
		s.network,
		clientCtx, from.String(), denomID, newTokenID, args...)
	s.Require().Equal(expectedCode, txResult.Code)

	querySupplyResponse = QuerySupplyExec(s.T(), s.network, clientCtx, denomID)
	s.Require().Equal(uint64(1), querySupplyResponse.Amount)

	//------test GetCmdTransferDenom()-------------
	args = []string{
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf(
			"--%s=%s",
			flags.FlagFees,
			sdk.NewCoins(sdk.NewCoin(s.network.BondDenom, sdk.NewInt(10))).String(),
		),
	}

	txResult = TransferDenomExec(s.T(),
		s.network,
		clientCtx, from.String(), val2.Address.String(), denomID, args...)
	s.Require().Equal(expectedCode, txResult.Code)

	queryDenomResponse = QueryDenomExec(s.T(), s.network, clientCtx, denomID)
	s.Require().Equal(val2.Address.String(), queryDenomResponse.Creator)
	s.Require().Equal(denomName, queryDenomResponse.Name)
	s.Require().Equal(schema, queryDenomResponse.Schema)
	s.Require().Equal(symbol, queryDenomResponse.Symbol)
	s.Require().Equal(mintRestricted, queryDenomResponse.MintRestricted)
	s.Require().Equal(updateRestricted, queryDenomResponse.UpdateRestricted)
}
