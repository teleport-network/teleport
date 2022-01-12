package types_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/teleport-network/teleport/app"
)

type TSSTestSuite struct {
	suite.Suite

	ctx sdk.Context
	app *app.Teleport
}

func (suite *TSSTestSuite) SetupTest() {
	teleport := app.Setup(false, nil)

	suite.ctx = teleport.BaseApp.NewContext(false, tmproto.Header{Time: time.Now()})
	suite.app = teleport
}

func TestTSSTestSuite(t *testing.T) {
	suite.Run(t, new(TSSTestSuite))
}
