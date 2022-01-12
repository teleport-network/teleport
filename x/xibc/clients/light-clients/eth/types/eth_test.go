package types_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/teleport-network/teleport/app"
)

type ETHTestSuite struct {
	suite.Suite

	ctx sdk.Context
	app *app.Teleport
}

func (suite *ETHTestSuite) SetupTest() {
	bitos := app.Setup(false, nil)

	suite.ctx = bitos.BaseApp.NewContext(false, tmproto.Header{Time: time.Now()})
	suite.app = bitos
}

func TestETHTestSuite(t *testing.T) {
	suite.Run(t, new(ETHTestSuite))
}
