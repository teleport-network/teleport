package types

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/ethereum/go-ethereum/common"

	"github.com/tendermint/tendermint/crypto/tmhash"

	"github.com/tharsis/ethermint/tests"
)

type TokenPairTestSuite struct {
	suite.Suite
}

func TestTokenPairSuite(t *testing.T) {
	suite.Run(t, new(TokenPairTestSuite))
}

func (suite *TokenPairTestSuite) TestTokenPairNew() {
	testCases := []struct {
		msg          string
		erc20Address common.Address
		denom        string
		enabled      bool
		owner        Owner
		expectPass   bool
	}{
		{msg: "Register token pair - invalid starts with number", erc20Address: tests.GenerateAddress(), denom: "1test", enabled: true, owner: OWNER_MODULE, expectPass: false},
		{msg: "Register token pair - invalid char '('", erc20Address: tests.GenerateAddress(), denom: "(test", enabled: true, owner: OWNER_MODULE, expectPass: false},
		{msg: "Register token pair - invalid char '^'", erc20Address: tests.GenerateAddress(), denom: "^test", enabled: true, owner: OWNER_MODULE, expectPass: false},
		// TODO: (guille) should the "\" be allowed to support unicode names?
		{msg: "Register token pair - invalid char '\\'", erc20Address: tests.GenerateAddress(), denom: "-test", enabled: true, owner: OWNER_MODULE, expectPass: false},
		// Invalid length
		{msg: "Register token pair - invalid length token (0)", erc20Address: tests.GenerateAddress(), denom: "", enabled: true, owner: OWNER_MODULE, expectPass: false},
		{msg: "Register token pair - invalid length token (1)", erc20Address: tests.GenerateAddress(), denom: "a", enabled: true, owner: OWNER_MODULE, expectPass: false},
		{msg: "Register token pair - invalid length token (128)", erc20Address: tests.GenerateAddress(), denom: strings.Repeat("a", 129), enabled: true, owner: OWNER_MODULE, expectPass: false},
		{msg: "Register token pair - pass", erc20Address: tests.GenerateAddress(), denom: "test", enabled: true, owner: OWNER_MODULE, expectPass: true},
	}

	for i, tc := range testCases {
		tp := NewTokenPair(tc.erc20Address, []string{tc.denom}, tc.enabled, tc.owner)
		err := tp.Validate()

		if tc.expectPass {
			suite.Require().NoError(err, "valid test %d failed: %s, %v", i, tc.msg)
		} else {
			suite.Require().Error(err, "invalid test %d passed: %s, %v", i, tc.msg)
		}
	}
}

func (suite *TokenPairTestSuite) TestTokenPair() {
	testCases := []struct {
		msg        string
		pair       TokenPair
		expectPass bool
	}{
		{msg: "Register token pair - invalid address (no hex)", pair: TokenPair{"0x5dCA2483280D9727c80b5518faC4556617fb19ZZ", []string{"test"}, true, OWNER_MODULE}, expectPass: false},
		{msg: "Register token pair - invalid address (invalid length 1)", pair: TokenPair{"0x5dCA2483280D9727c80b5518faC4556617fb19", []string{"test"}, true, OWNER_MODULE}, expectPass: false},
		{msg: "Register token pair - invalid address (invalid length 2)", pair: TokenPair{"0x5dCA2483280D9727c80b5518faC4556617fb194FFF", []string{"test"}, true, OWNER_MODULE}, expectPass: false},
		{msg: "pass", pair: TokenPair{tests.GenerateAddress().String(), []string{"test"}, true, OWNER_MODULE}, expectPass: true},
	}

	for i, tc := range testCases {
		if tc.expectPass {
			suite.Require().NoError(tc.pair.Validate(), "valid test %d failed: %s, %v", i, tc.msg)
		} else {
			suite.Require().Error(tc.pair.Validate(), "invalid test %d passed: %s, %v", i, tc.msg)
		}
	}
}

func (suite *TokenPairTestSuite) TestGetID() {
	addr := tests.GenerateAddress()
	denom := "test"
	pair := NewTokenPair(addr, []string{denom}, true, OWNER_MODULE)
	id := pair.GetID()
	expID := tmhash.Sum([]byte(addr.String() + "|" + denom))
	suite.Require().Equal(expID, id)
}

func (suite *TokenPairTestSuite) TestGetERC20Contract() {
	expAddr := tests.GenerateAddress()
	denom := "test"
	pair := NewTokenPair(expAddr, []string{denom}, true, OWNER_MODULE)
	addr := pair.GetERC20Contract()
	suite.Require().Equal(expAddr, addr)
}

func (suite *TokenPairTestSuite) TestIsNativeCoin() {
	testCases := []struct {
		name       string
		pair       TokenPair
		expectPass bool
	}{{
		"no owner",
		TokenPair{tests.GenerateAddress().String(), []string{"test"}, true, OWNER_UNSPECIFIED},
		false,
	}, {
		"external ERC20 owner",
		TokenPair{tests.GenerateAddress().String(), []string{"test"}, true, OWNER_EXTERNAL},
		false,
	}, {
		"pass",
		TokenPair{tests.GenerateAddress().String(), []string{"test"}, true, OWNER_MODULE},
		true,
	}}

	for _, tc := range testCases {
		if tc.expectPass {
			suite.Require().True(tc.pair.IsNativeCoin(), tc.name)
		} else {
			suite.Require().False(tc.pair.IsNativeCoin(), tc.name)
		}
	}
}

func (suite *TokenPairTestSuite) TestIsNativeERC20() {
	testCases := []struct {
		name       string
		pair       TokenPair
		expectPass bool
	}{{
		"no owner",
		TokenPair{tests.GenerateAddress().String(), []string{"test"}, true, OWNER_UNSPECIFIED},
		false,
	}, {
		"module owner",
		TokenPair{tests.GenerateAddress().String(), []string{"test"}, true, OWNER_MODULE},
		false,
	}, {
		"pass",
		TokenPair{tests.GenerateAddress().String(), []string{"test"}, true, OWNER_EXTERNAL},
		true,
	}}

	for _, tc := range testCases {
		if tc.expectPass {
			suite.Require().True(tc.pair.IsNativeERC20(), tc.name)
		} else {
			suite.Require().False(tc.pair.IsNativeERC20(), tc.name)
		}
	}
}
