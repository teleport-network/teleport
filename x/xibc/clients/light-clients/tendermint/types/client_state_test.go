package types_test

import (
	ics23 "github.com/confio/ics23/go"

	"github.com/teleport-network/teleport/x/xibc/clients/light-clients/tendermint/types"
	clienttypes "github.com/teleport-network/teleport/x/xibc/core/client/types"
	commitmenttypes "github.com/teleport-network/teleport/x/xibc/core/commitment/types"
	"github.com/teleport-network/teleport/x/xibc/core/host"
	packettypes "github.com/teleport-network/teleport/x/xibc/core/packet/types"
	"github.com/teleport-network/teleport/x/xibc/exported"
	xibctesting "github.com/teleport-network/teleport/x/xibc/testing"
)

var (
	invalidProof = []byte("invalid proof")
	prefix       = commitmenttypes.MerklePrefix{KeyPrefix: []byte("ibc")}
)

func (suite *TendermintTestSuite) TestValidate() {
	testCases := []struct {
		name        string
		clientState *types.ClientState
		expPass     bool
	}{{
		name:        "invalid chainID",
		clientState: types.NewClientState("  ", types.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, height, commitmenttypes.GetSDKSpecs(), prefix, 0),
		expPass:     false,
	}, {
		name:        "invalid trust level",
		clientState: types.NewClientState(chainID, types.Fraction{Numerator: 0, Denominator: 1}, trustingPeriod, ubdPeriod, maxClockDrift, height, commitmenttypes.GetSDKSpecs(), prefix, 0),
		expPass:     false,
	}, {
		name:        "invalid trusting period",
		clientState: types.NewClientState(chainID, types.DefaultTrustLevel, 0, ubdPeriod, maxClockDrift, height, commitmenttypes.GetSDKSpecs(), prefix, 0),
		expPass:     false,
	}, {
		name:        "invalid unbonding period",
		clientState: types.NewClientState(chainID, types.DefaultTrustLevel, trustingPeriod, 0, maxClockDrift, height, commitmenttypes.GetSDKSpecs(), prefix, 0),
		expPass:     false,
	}, {
		name:        "invalid max clock drift",
		clientState: types.NewClientState(chainID, types.DefaultTrustLevel, trustingPeriod, ubdPeriod, 0, height, commitmenttypes.GetSDKSpecs(), prefix, 0),
		expPass:     false,
	}, {
		name:        "invalid height",
		clientState: types.NewClientState(chainID, types.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, clienttypes.ZeroHeight(), commitmenttypes.GetSDKSpecs(), prefix, 0),
		expPass:     false,
	}, {
		name:        "trusting period not less than unbonding period",
		clientState: types.NewClientState(chainID, types.DefaultTrustLevel, ubdPeriod, ubdPeriod, maxClockDrift, height, commitmenttypes.GetSDKSpecs(), prefix, 0),
		expPass:     false,
	}, {
		name:        "proof specs is nil",
		clientState: types.NewClientState(chainID, types.DefaultTrustLevel, ubdPeriod, ubdPeriod, maxClockDrift, height, nil, prefix, 0),
		expPass:     false,
	}, {
		name:        "proof specs contains nil",
		clientState: types.NewClientState(chainID, types.DefaultTrustLevel, ubdPeriod, ubdPeriod, maxClockDrift, height, []*ics23.ProofSpec{ics23.TendermintSpec, nil}, prefix, 0),
		expPass:     false,
	}}

	for _, tc := range testCases {
		err := tc.clientState.Validate()
		if tc.expPass {
			suite.Require().NoError(err, tc.name)
		} else {
			suite.Require().Error(err, tc.name)
		}
	}
}

func (suite *TendermintTestSuite) TestInitialize() {
	testCases := []struct {
		name           string
		consensusState exported.ConsensusState
		expPass        bool
	}{{
		name:           "valid consensus",
		consensusState: &types.ConsensusState{},
		expPass:        true,
	}}

	path := xibctesting.NewPath(suite.chainA, suite.chainB)
	path.RegisterRelayers()
	err := path.EndpointA.CreateClient()
	suite.Require().NoError(err)

	clientState := path.EndpointA.GetClientState()

	relayers := path.EndpointA.Chain.App.XIBCKeeper.ClientKeeper.GetAllRelayers(path.EndpointA.Chain.GetContext())
	suite.Require().Equal(path.EndpointA.Chain.SenderAcc.String(), relayers[0].Address, "relayer does not match")
	store := path.EndpointA.ClientStore()

	for _, tc := range testCases {
		err := clientState.Initialize(suite.chainA.GetContext(), suite.chainA.Codec, store, tc.consensusState)
		if tc.expPass {
			suite.Require().NoError(err, "valid case returned an error")
		} else {
			suite.Require().Error(err, "invalid case didn't return an error")
		}
	}
}

// test verification of the packet commitment on chainB being represented
// in the light client on chainA. A send from chainB to chainA is simulated.
func (suite *TendermintTestSuite) TestVerifyPacketCommitment() {
	var (
		clientState *types.ClientState
		proof       []byte
		proofHeight exported.Height
	)

	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{{
		"successful verification",
		func() {},
		true,
	}, {
		"latest client height < height",
		func() {
			proofHeight = clientState.LatestHeight.Increment()
		},
		false,
	}, {
		"proof verification failed",
		func() {
			proof = invalidProof
		},
		false,
	}}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest() // reset

			// setup testing conditions
			path := xibctesting.NewPath(suite.chainA, suite.chainB)

			suite.coordinator.SetupClients(path)

			relayerA := path.EndpointA.Chain.App.XIBCKeeper.ClientKeeper.GetAllRelayers(path.EndpointA.Chain.GetContext())
			suite.Require().Equal(path.EndpointA.Chain.SenderAcc.String(), relayerA[0].Address, "relayer does not match")

			relayerB := path.EndpointB.Chain.App.XIBCKeeper.ClientKeeper.GetAllRelayers(path.EndpointB.Chain.GetContext())
			suite.Require().Equal(path.EndpointB.Chain.SenderAcc.String(), relayerB[0].Address, "relayer does not match")

			// setup testing conditions
			packet := packettypes.NewPacket(path.EndpointA.ChainName, path.EndpointB.ChainName, "", 1, []byte("mock transfer"), []byte("mock rcc"), "", 0)

			err := path.EndpointA.SendPacket(packet)
			suite.Require().NoError(err)

			var ok bool
			clientStateI := path.EndpointB.GetClientState()
			clientState, ok = clientStateI.(*types.ClientState)
			suite.Require().True(ok)

			// make packet commitment proof
			packetKey := host.PacketCommitmentKey(packet.GetSourceChain(), packet.GetDestChain(), packet.GetSequence())
			proof, proofHeight = suite.chainA.QueryProof(packetKey)

			tc.malleate() // make changes as necessary

			store := path.EndpointB.ClientStore()

			commitment, err := packettypes.CommitPacket(packet)
			suite.Require().NoError(err)

			err = clientState.VerifyPacketCommitment(
				suite.chainB.GetContext(), store, suite.chainB.Codec,
				proofHeight, proof, packet.GetSourceChain(),
				packet.GetDestChain(), packet.GetSequence(), commitment,
			)

			if tc.expPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

// test verification of the acknowledgement on chainB being represented
// in the light client on chainA. A send and ack from chainA to chainB
// is simulated.
func (suite *TendermintTestSuite) TestVerifyPacketAcknowledgement() {
	var (
		clientState *types.ClientState
		proof       []byte
		proofHeight exported.Height
	)

	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			"successful verification",
			func() {},
			true,
		},
		{
			"latest client height < height",
			func() {
				proofHeight = clientState.LatestHeight.Increment()
			},
			false,
		},
		{
			"proof verification failed",
			func() {
				proof = invalidProof
			},
			false,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest() // reset

			// setup testing conditions
			path := xibctesting.NewPath(suite.chainA, suite.chainB)
			suite.coordinator.SetupClients(path)

			packet := packettypes.NewPacket(
				path.EndpointA.ChainName,
				path.EndpointB.ChainName,
				"",
				1,
				[]byte("mock transfer"),
				[]byte("mock rcc"),
				"",
				0,
			)

			// send packet
			err := path.EndpointA.SendPacket(packet)
			suite.Require().NoError(err)

			// write receipt and ack
			err = path.EndpointB.RecvPacket(*packet)
			suite.Require().NoError(err)

			var ok bool
			clientStateI := path.EndpointA.GetClientState()
			clientState, ok = clientStateI.(*types.ClientState)
			suite.Require().True(ok)

			prefix = suite.chainB.GetPrefix()

			// make packet acknowledgement proof
			acknowledgementKey := host.PacketAcknowledgementKey(packet.GetSourceChain(), packet.GetDestChain(), packet.GetSequence())
			proof, proofHeight = suite.chainB.QueryProof(acknowledgementKey)

			// reset time and block delays to 0, malleate may change to a specific non-zero value.
			tc.malleate() // make changes as necessary

			ctx := suite.chainA.GetContext()
			store := path.EndpointA.ClientStore()

			ack, err := packettypes.NewResultAcknowledgement(
				0,
				[]byte(""),
				"",
				suite.chainB.SenderAcc.String(),
			).AbiPack()
			suite.Require().NoError(err)

			err = clientState.VerifyPacketAcknowledgement(
				ctx,
				store,
				suite.chainA.Codec,
				proofHeight,
				proof,
				packet.GetSourceChain(),
				packet.GetDestChain(),
				packet.GetSequence(),
				packettypes.CommitAcknowledgement(ack),
			)

			if tc.expPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}
