package types

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/ethereum/go-ethereum/common"
	"github.com/tharsis/ethermint/tests"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type MsgsTestSuite struct {
	suite.Suite
}

func TestMsgsTestSuite(t *testing.T) {
	suite.Run(t, new(MsgsTestSuite))
}

func (suite *MsgsTestSuite) TestMsgConvertCoinGetters() {
	msgInvalid := MsgConvertCoin{}
	msg := NewMsgConvertCoin(
		sdk.NewCoin("test", sdk.NewInt(100)),
		tests.GenerateAddress(),
		sdk.AccAddress(tests.GenerateAddress().Bytes()),
	)
	suite.Require().Nil(msgInvalid.GetSigners())
	suite.Require().NotNil(msg.GetSigners())
}

func (suite *MsgsTestSuite) TestMsgConvertCoinNew() {
	testCases := []struct {
		msg        string
		coin       sdk.Coin
		receiver   common.Address
		sender     sdk.AccAddress
		expectPass bool
	}{
		{
			"msg convert coin - pass",
			sdk.NewCoin("test", sdk.NewInt(100)),
			tests.GenerateAddress(),
			sdk.AccAddress(tests.GenerateAddress().Bytes()),
			true,
		},
	}

	for i, tc := range testCases {
		tx := NewMsgConvertCoin(tc.coin, tc.receiver, tc.sender)
		err := tx.ValidateBasic()

		if tc.expectPass {
			suite.Require().NoError(err, "valid test %d failed: %s, %v", i, tc.msg)
		} else {
			suite.Require().Error(err, "invalid test %d passed: %s, %v", i, tc.msg)
		}
	}
}

func (suite *MsgsTestSuite) TestMsgConvertCoin() {
	testCases := []struct {
		msg        string
		coin       sdk.Coin
		receiver   string
		sender     string
		expectPass bool
	}{
		{
			"invalid denom",
			sdk.Coin{
				Denom:  "",
				Amount: sdk.NewInt(100),
			},
			"0x0000",
			tests.GenerateAddress().String(),
			false,
		},
		{
			"negative coin amount",
			sdk.Coin{
				Denom:  "coin",
				Amount: sdk.NewInt(-100),
			},
			"0x0000",
			tests.GenerateAddress().String(),
			false,
		},
		{
			"msg convert coin - invalid sender",
			sdk.NewCoin("coin", sdk.NewInt(100)),
			tests.GenerateAddress().String(),
			"teleportinvalid",
			false,
		},
		{
			"msg convert coin - invalid receiver",
			sdk.NewCoin("coin", sdk.NewInt(100)),
			"0x0000",
			sdk.AccAddress(tests.GenerateAddress().Bytes()).String(),
			false,
		},
		{
			"msg convert coin - pass",
			sdk.NewCoin("coin", sdk.NewInt(100)),
			tests.GenerateAddress().String(),
			sdk.AccAddress(tests.GenerateAddress().Bytes()).String(),
			true,
		},
		{
			"msg convert coin - pass with `aggregate/` denom",
			sdk.NewCoin("aggregate/0xdac17f958d2ee523a2206206994597c13d831ec7", sdk.NewInt(100)),
			tests.GenerateAddress().String(),
			sdk.AccAddress(tests.GenerateAddress().Bytes()).String(),
			true,
		},
		{
			"msg convert coin - pass with `ibc/{hash}` denom",
			sdk.NewCoin("ibc/7F1D3FCF4AE79E1554D670D1AD949A9BA4E4A3C76C63093E17E446A46061A7A2", sdk.NewInt(100)),
			tests.GenerateAddress().String(),
			sdk.AccAddress(tests.GenerateAddress().Bytes()).String(),
			true,
		},
	}

	for i, tc := range testCases {
		tx := MsgConvertCoin{tc.coin, tc.receiver, tc.sender}
		err := tx.ValidateBasic()

		if tc.expectPass {
			suite.Require().NoError(err, "valid test %d failed: %s, %v", i, tc.msg)
		} else {
			suite.Require().Error(err, "invalid test %d passed: %s, %v", i, tc.msg)
		}
	}
}

func (suite *MsgsTestSuite) TestMsgConvertERC20Getters() {
	msg := NewMsgConvertERC20(
		sdk.NewInt(100),
		sdk.AccAddress(tests.GenerateAddress().Bytes()),
		tests.GenerateAddress(),
		tests.GenerateAddress(),
		"test",
	)
	suite.Require().NotNil(msg.GetSigners())
}

func (suite *MsgsTestSuite) TestMsgConvertERC20New() {
	testCases := []struct {
		msg        string
		amount     sdk.Int
		receiver   sdk.AccAddress
		contract   common.Address
		sender     common.Address
		denom      string
		expectPass bool
	}{
		{
			"msg convert erc20 - pass",
			sdk.NewInt(100),
			sdk.AccAddress(tests.GenerateAddress().Bytes()),
			tests.GenerateAddress(),
			tests.GenerateAddress(),
			"denom",
			true,
		},
	}

	for i, tc := range testCases {
		tx := NewMsgConvertERC20(tc.amount, tc.receiver, tc.contract, tc.sender, tc.denom)
		err := tx.ValidateBasic()

		if tc.expectPass {
			suite.Require().NoError(err, "valid test %d failed: %s, %v", i, tc.msg)
		} else {
			suite.Require().Error(err, "invalid test %d passed: %s, %v", i, tc.msg)
		}
	}
}

func (suite *MsgsTestSuite) TestMsgConvertERC20() {
	testCases := []struct {
		msg        string
		amount     sdk.Int
		receiver   string
		contract   string
		sender     string
		denom      string
		expectPass bool
	}{
		{
			"invalid contract hex address",
			sdk.NewInt(100),
			sdk.AccAddress(tests.GenerateAddress().Bytes()).String(),
			sdk.AccAddress{}.String(),
			tests.GenerateAddress().String(),
			"test",
			false,
		},
		{
			"negative coin amount",
			sdk.NewInt(-100),
			sdk.AccAddress(tests.GenerateAddress().Bytes()).String(),
			tests.GenerateAddress().String(),
			tests.GenerateAddress().String(),
			"test",
			false,
		},
		{
			"invalid receiver address",
			sdk.NewInt(100),
			sdk.AccAddress{}.String(),
			tests.GenerateAddress().String(),
			tests.GenerateAddress().String(),
			"test",
			false,
		},
		{
			"invalid sender address",
			sdk.NewInt(100),
			sdk.AccAddress(tests.GenerateAddress().Bytes()).String(),
			tests.GenerateAddress().String(),
			sdk.AccAddress{}.String(),
			"test",
			false,
		},
		{
			"msg convert erc20 - pass",
			sdk.NewInt(100),
			sdk.AccAddress(tests.GenerateAddress().Bytes()).String(),
			tests.GenerateAddress().String(),
			tests.GenerateAddress().String(),
			"test",
			true,
		},
	}

	for i, tc := range testCases {
		tx := MsgConvertERC20{tc.contract, tc.amount, tc.receiver, tc.sender, tc.denom}
		err := tx.ValidateBasic()

		if tc.expectPass {
			suite.Require().NoError(err, "valid test %d failed: %s, %v", i, tc.msg)
		} else {
			suite.Require().Error(err, "invalid test %d passed: %s, %v", i, tc.msg)
		}
	}
}
