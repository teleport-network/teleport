package keeper_test

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"

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
	suite.SetupTest()

	var contractAddr common.Address
	// Deploy contract
	switch contractType {
	case contractDirectBalanceManipulation:
		contractAddr = suite.DeployContractDirectBalanceManipulation(erc20Name, erc20Symbol)
	case contractMaliciousDelayed:
		contractAddr = suite.DeployContractMaliciousDelayed(erc20Name, erc20Symbol)
	default:
		contractAddr = suite.DeployContract(erc20Name, erc20Symbol, erc20Decimals)
	}
	suite.Commit()

	_, err := suite.app.AggregateKeeper.RegisterERC20(suite.ctx, contractAddr)
	suite.Require().NoError(err)
	return contractAddr
}

func (suite *KeeperTestSuite) setupRegisterCoin() (banktypes.Metadata, *types.TokenPair) {
	suite.SetupTest()
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

	err := suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, sdk.Coins{sdk.NewInt64Coin(validMetadata.Base, 1)})
	suite.Require().NoError(err)

	// pair := types.NewTokenPair(contractAddr, cosmosTokenBase, true, types.OWNER_MODULE)
	pair, err := suite.app.AggregateKeeper.RegisterCoin(suite.ctx, validMetadata)
	suite.Require().NoError(err)
	suite.Commit()
	return validMetadata, pair
}

func (suite KeeperTestSuite) TestRegisterCoin() {
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
			"intrarelaying is disabled globally",
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
				regPair := types.NewTokenPair(tests.GenerateAddress(), metadata.Base, true, types.OWNER_MODULE)
				suite.app.AggregateKeeper.SetDenomMap(suite.ctx, regPair.Denom, regPair.GetID())
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

				err := suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, sdk.Coins{sdk.NewInt64Coin(validMetadata.Base, 1)})
				suite.Require().NoError(err)
				suite.app.BankKeeper.SetDenomMetaData(suite.ctx, validMetadata)
			},
			false,
		},
		{
			"ok",
			func() {
				err := suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, sdk.Coins{sdk.NewInt64Coin(metadata.Base, 1)})
				suite.Require().NoError(err)
			},
			true,
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
				Denom:         "acoin",
				Enabled:       true,
				ContractOwner: 1,
			}

			if tc.expPass {
				suite.Require().NoError(err, tc.name)
				suite.Require().Equal(expPair, pair)
			} else {
				suite.Require().Error(err, tc.name)
			}
		})
	}
}

func (suite KeeperTestSuite) TestRegisterERC20() {
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
			"intrarelaying is disabled globally",
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
				suite.app.AggregateKeeper.SetDenomMap(suite.ctx, pair.Denom, pair.GetID())
			},
			false,
		},
		{
			"meta data already stored",
			func() {
				_, err := suite.app.AggregateKeeper.CreateCoinMetadata(suite.ctx, contractAddr)
				suite.Require().NoError(err)
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
			pair = types.NewTokenPair(contractAddr, coinName, true, types.OWNER_EXTERNAL)

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

func (suite KeeperTestSuite) TestToggleRelay() {
	var (
		contractAddr common.Address
		id           []byte
		pair         types.TokenPair
	)

	testCases := []struct {
		name         string
		malleate     func()
		expPass      bool
		relayEnabled bool
	}{
		{
			"token not registered",
			func() {
				contractAddr = suite.DeployContract(erc20Name, erc20Symbol, erc20Decimals)
				suite.Commit()
				pair = types.NewTokenPair(contractAddr, cosmosTokenBase, true, types.OWNER_MODULE)
			},
			false,
			false,
		},
		{
			"token not registered - pair not found",
			func() {
				contractAddr = suite.DeployContract(erc20Name, erc20Symbol, erc20Decimals)
				suite.Commit()
				pair = types.NewTokenPair(contractAddr, cosmosTokenBase, true, types.OWNER_MODULE)
				suite.app.AggregateKeeper.SetERC20Map(suite.ctx, common.HexToAddress(pair.ERC20Address), pair.GetID())
			},
			false,
			false,
		},
		{
			"disable relay",
			func() {
				contractAddr = suite.setupRegisterERC20Pair(contractMinterBurner)
				id = suite.app.AggregateKeeper.GetTokenPairID(suite.ctx, contractAddr.String())
				pair, _ = suite.app.AggregateKeeper.GetTokenPair(suite.ctx, id)
			},
			true,
			false,
		},
		{
			"disable and enable relay",
			func() {
				contractAddr = suite.setupRegisterERC20Pair(contractMinterBurner)
				id = suite.app.AggregateKeeper.GetTokenPairID(suite.ctx, contractAddr.String())
				pair, _ = suite.app.AggregateKeeper.GetTokenPair(suite.ctx, id)
				pair, _ = suite.app.AggregateKeeper.ToggleRelay(suite.ctx, contractAddr.String())
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
			pair, err = suite.app.AggregateKeeper.ToggleRelay(suite.ctx, contractAddr.String())
			// Request the pair using the GetPairToken func to make sure that is updated on the db
			pair, _ = suite.app.AggregateKeeper.GetTokenPair(suite.ctx, id)
			if tc.expPass {
				suite.Require().NoError(err, tc.name)
				if tc.relayEnabled {
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

func (suite KeeperTestSuite) TestUpdateTokenPairERC20() {
	var (
		contractAddr    common.Address
		pair            types.TokenPair
		metadata        banktypes.Metadata
		newContractAddr common.Address
	)

	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			"token not registered",
			func() {
				contractAddr = suite.DeployContract(erc20Name, erc20Symbol, erc20Decimals)
				suite.Commit()
				pair = types.NewTokenPair(contractAddr, cosmosTokenBase, true, types.OWNER_MODULE)
			},
			false,
		},
		{
			"token not registered - pair not found",
			func() {
				contractAddr = suite.DeployContract(erc20Name, erc20Symbol, erc20Decimals)
				suite.Commit()
				pair = types.NewTokenPair(contractAddr, cosmosTokenBase, true, types.OWNER_MODULE)

				suite.app.AggregateKeeper.SetERC20Map(suite.ctx, common.HexToAddress(pair.ERC20Address), pair.GetID())
			},
			false,
		},
		{
			"token not registered - Metadata not found",
			func() {
				contractAddr = suite.DeployContract(erc20Name, erc20Symbol, erc20Decimals)
				suite.Commit()
				pair = types.NewTokenPair(contractAddr, cosmosTokenBase, true, types.OWNER_MODULE)

				suite.app.AggregateKeeper.SetTokenPair(suite.ctx, pair)
				suite.app.AggregateKeeper.SetDenomMap(suite.ctx, pair.Denom, pair.GetID())
				suite.app.AggregateKeeper.SetERC20Map(suite.ctx, common.HexToAddress(pair.ERC20Address), pair.GetID())
			},
			false,
		},
		{
			"newErc20 not found",
			func() {
				contractAddr = suite.setupRegisterERC20Pair(contractMinterBurner)
				newContractAddr = common.Address{}
			},
			false,
		},
		{
			"empty denom units",
			func() {
				var found bool
				contractAddr = suite.setupRegisterERC20Pair(contractMinterBurner)
				id := suite.app.AggregateKeeper.GetTokenPairID(suite.ctx, contractAddr.String())
				pair, found = suite.app.AggregateKeeper.GetTokenPair(suite.ctx, id)
				suite.Require().True(found)
				suite.app.BankKeeper.SetDenomMetaData(suite.ctx, banktypes.Metadata{Base: pair.Denom})
				suite.Commit()

				// Deploy a new contract with the same values
				newContractAddr = suite.DeployContract(erc20Name, erc20Symbol, erc20Decimals)
			},
			false,
		},
		{
			"metadata ERC20 details mismatch",
			func() {
				var found bool
				contractAddr = suite.setupRegisterERC20Pair(contractMinterBurner)
				id := suite.app.AggregateKeeper.GetTokenPairID(suite.ctx, contractAddr.String())
				pair, found = suite.app.AggregateKeeper.GetTokenPair(suite.ctx, id)
				suite.Require().True(found)
				metadata := banktypes.Metadata{Base: pair.Denom, DenomUnits: []*banktypes.DenomUnit{{}}}
				suite.app.BankKeeper.SetDenomMetaData(suite.ctx, metadata)
				suite.Commit()

				// Deploy a new contract with the same values
				newContractAddr = suite.DeployContract(erc20Name, erc20Symbol, erc20Decimals)
			},
			false,
		},
		{
			"no denom unit with ERC20 name",
			func() {
				var found bool
				contractAddr = suite.setupRegisterERC20Pair(contractMinterBurner)
				id := suite.app.AggregateKeeper.GetTokenPairID(suite.ctx, contractAddr.String())
				pair, found = suite.app.AggregateKeeper.GetTokenPair(suite.ctx, id)
				suite.Require().True(found)
				metadata := banktypes.Metadata{Base: pair.Denom, Display: erc20Name, Description: types.CreateDenomDescription(contractAddr.String()), Symbol: erc20Symbol, DenomUnits: []*banktypes.DenomUnit{{}}}
				suite.app.BankKeeper.SetDenomMetaData(suite.ctx, metadata)
				suite.Commit()

				// Deploy a new contract with the same values
				newContractAddr = suite.DeployContract(erc20Name, erc20Symbol, erc20Decimals)
			},
			false,
		},
		{
			"denom unit and ERC20 decimals mismatch",
			func() {
				var found bool
				contractAddr = suite.setupRegisterERC20Pair(contractMinterBurner)
				id := suite.app.AggregateKeeper.GetTokenPairID(suite.ctx, contractAddr.String())
				pair, found = suite.app.AggregateKeeper.GetTokenPair(suite.ctx, id)
				suite.Require().True(found)
				metadata := banktypes.Metadata{Base: pair.Denom, Display: erc20Name, Description: types.CreateDenomDescription(contractAddr.String()), Symbol: erc20Symbol, DenomUnits: []*banktypes.DenomUnit{{Denom: erc20Name}}}
				suite.app.BankKeeper.SetDenomMetaData(suite.ctx, metadata)
				suite.Commit()

				// Deploy a new contract with the same values
				newContractAddr = suite.DeployContract(erc20Name, erc20Symbol, erc20Decimals)
			},
			false,
		},
		{
			"ok",
			func() {
				var found bool
				contractAddr = suite.setupRegisterERC20Pair(contractMinterBurner)
				id := suite.app.AggregateKeeper.GetTokenPairID(suite.ctx, contractAddr.String())
				pair, found = suite.app.AggregateKeeper.GetTokenPair(suite.ctx, id)
				suite.Require().True(found)
				metadata := banktypes.Metadata{Base: pair.Denom, Display: erc20Name, Description: types.CreateDenomDescription(contractAddr.String()), Symbol: erc20Symbol, DenomUnits: []*banktypes.DenomUnit{{Denom: erc20Name, Exponent: 18}}}
				suite.app.BankKeeper.SetDenomMetaData(suite.ctx, metadata)
				suite.Commit()

				// Deploy a new contract with the same values
				newContractAddr = suite.DeployContract(erc20Name, erc20Symbol, erc20Decimals)
			},
			true,
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset

			tc.malleate()

			var err error
			pair, err = suite.app.AggregateKeeper.UpdateTokenPairERC20(suite.ctx, contractAddr, newContractAddr)
			metadata, _ = suite.app.BankKeeper.GetDenomMetaData(suite.ctx, types.CreateDenom(contractAddr.String()))

			if tc.expPass {
				suite.Require().NoError(err, tc.name)
				suite.Require().Equal(newContractAddr.String(), pair.ERC20Address)
				suite.Require().Equal(types.CreateDenomDescription(newContractAddr.String()), metadata.Description)
			} else {
				suite.Require().Error(err, tc.name)
				if suite.app.AggregateKeeper.IsTokenPairRegistered(suite.ctx, pair.GetID()) {
					suite.Require().Equal(contractAddr.String(), pair.ERC20Address, "check pair")
					suite.Require().Equal(types.CreateDenomDescription(contractAddr.String()), metadata.Description, "check metadata")
				}
			}
		})
	}
}
