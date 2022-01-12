package types_test

import (
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"

	xibctsstypes "github.com/teleport-network/teleport/x/xibc/clients/tss-client/types"
	clienttypes "github.com/teleport-network/teleport/x/xibc/core/client/types"
	"github.com/teleport-network/teleport/x/xibc/exported"
)

var chainName = "tss"

func (suite *TSSTestSuite) TestCheckHeaderAndUpdateState() {
	number := clienttypes.NewHeight(0, 0)
	initTssAddress := randomAddr()
	updatedTssAddress := randomAddr()

	clientState := exported.ClientState(&xibctsstypes.ClientState{
		ChainId:    "tss",
		TssAddress: initTssAddress,
	})

	updateHeader := xibctsstypes.Header{
		TssAddress: updatedTssAddress,
	}

	var consensusState exported.ConsensusState

	suite.app.XIBCKeeper.ClientKeeper.SetClientConsensusState(suite.ctx, chainName, number, consensusState)

	clientState, _, err := clientState.CheckHeaderAndUpdateState(
		suite.ctx,
		suite.app.AppCodec(),
		suite.app.XIBCKeeper.ClientKeeper.ClientStore(suite.ctx, chainName), // pass in chainName prefixed clientStore
		&updateHeader,
	)

	suite.Require().NoError(err)

	suite.Require().Equal(updatedTssAddress, clientState.(*xibctsstypes.ClientState).TssAddress)
}

func randomAddr() string {
	address := make([]byte, 20)
	rand.Read(address)
	return sdk.AccAddress(address).String()
}
