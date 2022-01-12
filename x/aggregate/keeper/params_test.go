package keeper_test

import (
	"github.com/teleport-network/teleport/x/aggregate/types"
)

func (suite *KeeperTestSuite) TestParams() {
	params := suite.app.AggregateKeeper.GetParams(suite.ctx)
	suite.Require().Equal(types.DefaultParams(), params)
	params.EnableAggregate = false
	suite.app.AggregateKeeper.SetParams(suite.ctx, params)
	newParams := suite.app.AggregateKeeper.GetParams(suite.ctx)
	suite.Require().Equal(newParams, params)
}
