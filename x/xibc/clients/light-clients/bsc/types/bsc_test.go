package types_test

import (
	"testing"
	"time"

	"github.com/teleport-network/teleport/app"

	"github.com/stretchr/testify/suite"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type BSCTestSuite struct {
	suite.Suite

	ctx sdk.Context
	app *app.Teleport
}

func (suite *BSCTestSuite) SetupTest() {
	teleprot := app.Setup(false, nil)

	suite.ctx = teleprot.BaseApp.NewContext(false, tmproto.Header{Time: time.Now()})
	suite.app = teleprot
}

func TestBSCTestSuite(t *testing.T) {
	suite.Run(t, new(BSCTestSuite))
}
