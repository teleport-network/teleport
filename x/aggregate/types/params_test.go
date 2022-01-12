package types

import (
	"testing"

	"github.com/stretchr/testify/suite"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

type ParamsTestSuite struct {
	suite.Suite
}

func TestParamsTestSuite(t *testing.T) {
	suite.Run(t, new(ParamsTestSuite))
}

func (suite *ParamsTestSuite) TestParamKeyTable() {
	suite.Require().IsType(paramtypes.KeyTable{}, ParamKeyTable())
}

func (suite *ParamsTestSuite) TestParamsValidate() {
	testCases := []struct {
		name     string
		params   Params
		expError bool
	}{
		{"default", DefaultParams(), false},
		{"valid", NewParams(true, true), false},
		{"empty", Params{}, false},
	}

	for _, tc := range testCases {
		if tc.expError {
			suite.Require().Error(tc.params.Validate(), tc.name)
		} else {
			suite.Require().NoError(tc.params.Validate(), tc.name)
		}
	}
}

func (suite *ParamsTestSuite) TestParamsValidatePriv() {
	suite.Require().Error(validateBool(1))
	suite.Require().NoError(validateBool(true))
}
