package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/teleport-network/teleport/x/xibc/core/routing/types"
	xibctesting "github.com/teleport-network/teleport/x/xibc/testing"
)

type KeeperTestSuite struct {
	suite.Suite
	coordinator *xibctesting.Coordinator
	chain       *xibctesting.TestChain
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.coordinator = xibctesting.NewCoordinator(suite.T(), 1)
	suite.chain = suite.coordinator.GetChain(xibctesting.GetChainID(0))
	suite.coordinator.CommitNBlocks(suite.chain, 2)
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) TestSetRouter() {
	// There is a sealed router in the RoutingKeeper
	suite.Require().True(suite.chain.App.XIBCKeeper.RoutingKeeper.Router.Sealed())
	suite.chain.App.XIBCKeeper.RoutingKeeper.Router = nil
	router := types.NewRouter()
	router.AddRoute("1", nil)
	router.Sealed()
	suite.chain.App.XIBCKeeper.RoutingKeeper.SetRouter(router)
	suite.Require().Equal(router, suite.chain.App.XIBCKeeper.RoutingKeeper.Router)
}
