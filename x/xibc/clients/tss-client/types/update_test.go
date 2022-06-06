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

	initPubkey := randomPubkey()
	updatedPubkey := randomPubkey()

	initPartPubkeys := [][]byte{randomPubkey()}
	updatedPartPubkeys := [][]byte{randomPubkey()}

	initThreshold := rand.Uint64()
	updatedThreshold := rand.Uint64()

	clientState := exported.ClientState(&xibctsstypes.ClientState{
		TssAddress:  initTssAddress,
		Pubkey:      initPubkey,
		PartPubkeys: initPartPubkeys,
		Threshold:   initThreshold,
	})

	updateHeader := xibctsstypes.Header{
		TssAddress:  updatedTssAddress,
		Pubkey:      updatedPubkey,
		PartPubkeys: updatedPartPubkeys,
		Threshold:   updatedThreshold,
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
	suite.Require().Equal(updatedPubkey, clientState.(*xibctsstypes.ClientState).Pubkey)
	suite.Require().Equal(updatedPartPubkeys, clientState.(*xibctsstypes.ClientState).PartPubkeys)
	suite.Require().Equal(updatedThreshold, clientState.(*xibctsstypes.ClientState).Threshold)
}

func randomPubkey() []byte {
	pubkey := make([]byte, 20)
	rand.Read(pubkey)
	return pubkey
}

func randomAddr() string {
	address := make([]byte, 20)
	rand.Read(address)
	return sdk.AccAddress(address).String()
}
