package keeper_test

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/tharsis/ethermint/tests"

	"github.com/teleport-network/teleport/x/aggregate/types"
)

const (
	contractMinterBurner = iota + 1
	contractDirectBalanceManipulation
	contractMaliciousDelayed
)

const (
	erc20Name          = "Coin Token"
	erc20Symbol        = "CTKN"
	erc20Decimals      = uint8(18)
	cosmosTokenBase    = "acoin"
	cosmosTokenDisplay = "coin"
	cosmosDecimals     = uint8(6)
	defaultExponent    = uint32(18)
	zeroExponent       = uint32(0)
)

func (suite *KeeperTestSuite) setupRegisterERC20Pair(contractType int) common.Address {
	var contract common.Address
	// Deploy contract
	switch contractType {
	case contractDirectBalanceManipulation:
		contract = suite.DeployContractDirectBalanceManipulation(erc20Name, erc20Symbol)
	case contractMaliciousDelayed:
		contract = suite.DeployContractMaliciousDelayed(erc20Name, erc20Symbol)
	default:
		contract = suite.DeployContract(erc20Name, erc20Symbol, erc20Decimals)
	}
	suite.Commit()

	_, err := suite.app.AggregateKeeper.RegisterERC20(suite.ctx, contract)
	suite.Require().NoError(err)
	return contract
}

func (suite *KeeperTestSuite) setupRegisterCoin() (banktypes.Metadata, *types.TokenPair) {
	validMetadata := banktypes.Metadata{
		Description: "description of the token",
		Base:        cosmosTokenBase,
		// NOTE: Denom units MUST be increasing
		DenomUnits: []*banktypes.DenomUnit{
			{
				Denom:    cosmosTokenBase,
				Exponent: 0,
			},
			{
				Denom:    cosmosTokenBase[1:],
				Exponent: uint32(18),
			},
		},
		Name:    cosmosTokenBase,
		Symbol:  erc20Symbol,
		Display: cosmosTokenBase,
	}

	err := suite.app.BankKeeper.MintCoins(suite.ctx, types.ModuleName, sdk.Coins{sdk.NewInt64Coin(validMetadata.Base, 1)})
	suite.Require().NoError(err)

	// pair := types.NewTokenPair(contractAddr, cosmosTokenBase, true, types.OWNER_MODULE)
	pair, err := suite.app.AggregateKeeper.RegisterCoin(suite.ctx, validMetadata)
	suite.Require().NoError(err)
	suite.Commit()
	return validMetadata, pair
}

func (suite KeeperTestSuite) TestRegisterCoin() { //nolint
	metadata := banktypes.Metadata{
		Description: "description",
		Base:        cosmosTokenBase,
		// NOTE: Denom units MUST be increasing
		DenomUnits: []*banktypes.DenomUnit{
			{
				Denom:    cosmosTokenBase,
				Exponent: 0,
			},
			{
				Denom:    cosmosTokenDisplay,
				Exponent: defaultExponent,
			},
		},
		Name:    cosmosTokenBase,
		Symbol:  erc20Symbol,
		Display: cosmosTokenDisplay,
	}

	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			"ok",
			func() {
				metadata.Base = cosmosTokenBase
				err := suite.app.BankKeeper.MintCoins(suite.ctx, types.ModuleName, sdk.Coins{sdk.NewInt64Coin(metadata.Base, 1)})
				suite.Require().NoError(err)
			},
			true,
		},
		{
			"conversion is disabled globally",
			func() {
				params := types.DefaultParams()
				params.EnableAggregate = false
				suite.app.AggregateKeeper.SetParams(suite.ctx, params)
			},
			false,
		},
		{
			"denom already registered",
			func() {
				regPair := types.NewTokenPair(tests.GenerateAddress(), []string{metadata.Base}, true, types.OWNER_MODULE)
				suite.app.AggregateKeeper.SetDenomsMap(suite.ctx, regPair.Denoms, regPair.GetID())
				suite.Commit()
			},
			false,
		},
		{
			"token doesn't have supply",
			func() {
			},
			false,
		},
		{
			"metadata different that stored",
			func() {
				metadata.Base = cosmosTokenBase
				validMetadata := banktypes.Metadata{
					Description: "description",
					Base:        cosmosTokenBase,
					// NOTE: Denom units MUST be increasing
					DenomUnits: []*banktypes.DenomUnit{
						{
							Denom:    cosmosTokenBase,
							Exponent: 0,
						},
						{
							Denom:    cosmosTokenDisplay,
							Exponent: uint32(18),
						},
					},
					Name:    erc20Name,
					Symbol:  erc20Symbol,
					Display: cosmosTokenDisplay,
				}

				err := suite.app.BankKeeper.MintCoins(suite.ctx, types.ModuleName, sdk.Coins{sdk.NewInt64Coin(validMetadata.Base, 1)})
				suite.Require().NoError(err)
				suite.app.BankKeeper.SetDenomMetaData(suite.ctx, validMetadata)
			},
			false,
		},
		{
			"tele denom registration - tele",
			func() {
				metadata.Base = "tele"
				err := suite.app.BankKeeper.MintCoins(suite.ctx, types.ModuleName, sdk.Coins{sdk.NewInt64Coin(metadata.Base, 1)})
				suite.Require().NoError(err)
			},
			false,
		},
		{
			"tele denom registration - teleport",
			func() {
				metadata.Base = "teleport"
				err := suite.app.BankKeeper.MintCoins(suite.ctx, types.ModuleName, sdk.Coins{sdk.NewInt64Coin(metadata.Base, 1)})
				suite.Require().NoError(err)
			},
			false,
		},
		{
			"tele denom registration - atele",
			func() {
				metadata.Base = "atele"
				err := suite.app.BankKeeper.MintCoins(suite.ctx, types.ModuleName, sdk.Coins{sdk.NewInt64Coin(metadata.Base, 1)})
				suite.Require().NoError(err)
			},
			false,
		},
		{
			"tele denom registration - wtele",
			func() {
				metadata.Base = "wtele"
				err := suite.app.BankKeeper.MintCoins(suite.ctx, types.ModuleName, sdk.Coins{sdk.NewInt64Coin(metadata.Base, 1)})
				suite.Require().NoError(err)
			},
			false,
		},
		{
			"force delete module account evm",
			func() {
				metadata.Base = cosmosTokenBase
				err := suite.app.BankKeeper.MintCoins(suite.ctx, types.ModuleName, sdk.Coins{sdk.NewInt64Coin(metadata.Base, 1)})
				suite.Require().NoError(err)

				acc := suite.app.AccountKeeper.GetAccount(suite.ctx, types.ModuleAddress.Bytes())
				suite.app.AccountKeeper.RemoveAccount(suite.ctx, acc)
			},
			false,
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset

			tc.malleate()

			pair, err := suite.app.AggregateKeeper.RegisterCoin(suite.ctx, metadata)
			suite.Commit()

			expPair := &types.TokenPair{
				ERC20Address:  "0x90d3e9B208998d1048467bFDcbE3661322373712",
				Denoms:        []string{"acoin"},
				Enabled:       true,
				ContractOwner: 1,
			}

			if tc.expPass {
				suite.Require().NoError(err, tc.name)
				suite.Require().Equal(pair, expPair)
			} else {
				suite.Require().Error(err, tc.name)
			}
		})
	}
}

func (suite KeeperTestSuite) TestRegisterERC20() { //nolint
	var (
		contractAddr common.Address
		pair         types.TokenPair
	)
	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			"conversion is disabled globally",
			func() {
				params := types.DefaultParams()
				params.EnableAggregate = false
				suite.app.AggregateKeeper.SetParams(suite.ctx, params)
			},
			false,
		},
		{
			"token ERC20 already registered",
			func() {
				suite.app.AggregateKeeper.SetERC20Map(suite.ctx, pair.GetERC20Contract(), pair.GetID())
			},
			false,
		},
		{
			"denom already registered",
			func() {
				suite.app.AggregateKeeper.SetDenomsMap(suite.ctx, pair.Denoms, pair.GetID())
			},
			false,
		},
		{
			"meta data already stored",
			func() {
				_, _ = suite.app.AggregateKeeper.CreateCoinMetadata(suite.ctx, contractAddr)
			},
			false,
		},
		{
			"ok",
			func() {},
			true,
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset

			contractAddr = suite.DeployContract(erc20Name, erc20Symbol, cosmosDecimals)
			suite.Commit()
			coinName := types.CreateDenom(contractAddr.String())
			pair = types.NewTokenPair(contractAddr, []string{coinName}, true, types.OWNER_EXTERNAL)

			tc.malleate()

			_, err := suite.app.AggregateKeeper.RegisterERC20(suite.ctx, contractAddr)
			metadata, found := suite.app.BankKeeper.GetDenomMetaData(suite.ctx, coinName)
			if tc.expPass {
				suite.Require().NoError(err, tc.name)
				// Metadata variables
				suite.Require().True(found)
				suite.Require().Equal(coinName, metadata.Base)
				suite.Require().Equal(coinName, metadata.Name)
				suite.Require().Equal(types.SanitizeERC20Name(erc20Name), metadata.Display)
				suite.Require().Equal(erc20Symbol, metadata.Symbol)
				// Denom units
				suite.Require().Equal(len(metadata.DenomUnits), 2)
				suite.Require().Equal(coinName, metadata.DenomUnits[0].Denom)
				suite.Require().Equal(uint32(zeroExponent), metadata.DenomUnits[0].Exponent)
				suite.Require().Equal(types.SanitizeERC20Name(erc20Name), metadata.DenomUnits[1].Denom)
				// Custom exponent at contract creation matches coin with token
				suite.Require().Equal(metadata.DenomUnits[1].Exponent, uint32(cosmosDecimals))
			} else {
				suite.Require().Error(err, tc.name)
			}
		})
	}
}

func (suite KeeperTestSuite) TestToggleConverision() { //nolint
	var (
		contractAddr common.Address
		id           []byte
		pair         types.TokenPair
	)

	testCases := []struct {
		name              string
		malleate          func()
		expPass           bool
		conversionEnabled bool
	}{
		{
			"token not registered",
			func() {
				contractAddr = suite.DeployContract(erc20Name, erc20Symbol, erc20Decimals)
				suite.Commit()
				pair = types.NewTokenPair(contractAddr, []string{cosmosTokenBase}, true, types.OWNER_MODULE)
			},
			false,
			false,
		},
		{
			"token not registered - pair not found",
			func() {
				contractAddr = suite.DeployContract(erc20Name, erc20Symbol, erc20Decimals)
				suite.Commit()
				pair = types.NewTokenPair(contractAddr, []string{cosmosTokenBase}, true, types.OWNER_MODULE)
				suite.app.AggregateKeeper.SetERC20Map(suite.ctx, common.HexToAddress(pair.ERC20Address), pair.GetID())
			},
			false,
			false,
		},
		{
			"disable conversion",
			func() {
				contractAddr = suite.setupRegisterERC20Pair(contractMinterBurner)
				id = suite.app.AggregateKeeper.GetTokenPairID(suite.ctx, contractAddr.String())
				pair, _ = suite.app.AggregateKeeper.GetTokenPair(suite.ctx, id)
			},
			true,
			false,
		},
		{
			"disable and enable conversion",
			func() {
				contractAddr = suite.setupRegisterERC20Pair(contractMinterBurner)
				id = suite.app.AggregateKeeper.GetTokenPairID(suite.ctx, contractAddr.String())
				pair, _ = suite.app.AggregateKeeper.GetTokenPair(suite.ctx, id)
				pair, _ = suite.app.AggregateKeeper.ToggleConversion(suite.ctx, contractAddr.String())
			},
			true,
			true,
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset

			tc.malleate()

			var err error
			pair, err = suite.app.AggregateKeeper.ToggleConversion(suite.ctx, contractAddr.String())
			// Request the pair using the GetPairToken func to make sure that is updated on the db
			pair, _ = suite.app.AggregateKeeper.GetTokenPair(suite.ctx, id)
			if tc.expPass {
				suite.Require().NoError(err, tc.name)
				if tc.conversionEnabled {
					suite.Require().True(pair.Enabled)
				} else {
					suite.Require().False(pair.Enabled)
				}
			} else {
				suite.Require().Error(err, tc.name)
			}
		})
	}
}
