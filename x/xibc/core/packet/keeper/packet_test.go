package keeper_test

import (
	"fmt"

	host "github.com/teleport-network/teleport/x/xibc/core/host"
	"github.com/teleport-network/teleport/x/xibc/core/packet/types"
	"github.com/teleport-network/teleport/x/xibc/exported"
	xibctesting "github.com/teleport-network/teleport/x/xibc/testing"
	xibcmock "github.com/teleport-network/teleport/x/xibc/testing/mock"
)

type testCase = struct {
	name     string
	malleate func()
	expPass  bool
}

var (
	validPacketData = []byte("VALID PACKET DATA")
	relayChain      = ""
)

// TestSendPacket tests SendPacket from chainA to chainB
func (suite *KeeperTestSuite) TestSendPacket() {
	var packet exported.PacketI

	testCases := []testCase{{
		"success",
		func() {
			path := xibctesting.NewPath(suite.chainA, suite.chainB)
			suite.coordinator.SetupClients(path)

			packet = types.NewPacket(1, path.EndpointA.ChainName, path.EndpointB.ChainName, relayChain, []string{xibctesting.MockPort}, [][]byte{validPacketData})
		},
		true,
	}, {
		"sending packet out of order ",
		func() {
			path := xibctesting.NewPath(suite.chainA, suite.chainB)
			suite.coordinator.SetupClients(path)
			packet = types.NewPacket(5, path.EndpointA.ChainName, path.EndpointB.ChainName, relayChain, []string{xibctesting.MockPort}, [][]byte{validPacketData})
		},
		false,
	}, {
		"client state not found",
		func() {
			path := xibctesting.NewPath(suite.chainA, suite.chainB)
			suite.coordinator.SetupClients(path)

			packet = types.NewPacket(1, path.EndpointA.ChainName, xibctesting.InvalidID, relayChain, []string{xibctesting.MockPort}, [][]byte{validPacketData})
		},
		false,
	}, {
		"next sequence wrong",
		func() {
			path := xibctesting.NewPath(suite.chainA, suite.chainB)
			suite.coordinator.SetupClients(path)
			packet = types.NewPacket(1, path.EndpointA.ChainName, path.EndpointB.ChainName, relayChain, []string{xibctesting.MockPort}, [][]byte{validPacketData})
			suite.chainA.App.XIBCKeeper.PacketKeeper.SetNextSequenceSend(suite.chainA.GetContext(), path.EndpointA.ChainName, path.EndpointB.ChainName, 5)
		},
		false,
	}}

	for i, tc := range testCases {
		suite.Run(
			fmt.Sprintf("Case %s, %d", tc.name, i),
			func() {
				suite.SetupTest() // reset
				tc.malleate()

				err := suite.chainA.App.XIBCKeeper.PacketKeeper.SendPacket(suite.chainA.GetContext(), packet)
				if tc.expPass {
					suite.Require().NoError(err)
				} else {
					suite.Require().Error(err)
				}
			},
		)
	}
}

// TestRecvPacket test RecvPacket on chainB. Since packet commitment verification will always
// occur last (resource instensive), only tests expected to succeed and packet commitment
// verification tests need to simulate sending a packet from chainA to chainB.
func (suite *KeeperTestSuite) TestRecvPacket() {
	var packet types.Packet

	testCases := []testCase{{
		"success",
		func() {
			path := xibctesting.NewPath(suite.chainA, suite.chainB)
			suite.coordinator.SetupClients(path)

			packet = types.NewPacket(1, path.EndpointA.ChainName, path.EndpointB.ChainName, relayChain, []string{xibctesting.MockPort}, [][]byte{validPacketData})
			err := path.EndpointA.SendPacket(packet)
			suite.Require().NoError(err)
		},
		true,
	}, {
		"success with out of order packet",
		func() {
			path := xibctesting.NewPath(suite.chainA, suite.chainB)
			suite.coordinator.SetupClients(path)
			packet = types.NewPacket(1, path.EndpointA.ChainName, path.EndpointB.ChainName, relayChain, []string{xibctesting.MockPort}, [][]byte{validPacketData})

			// send 2 packets
			err := path.EndpointA.SendPacket(packet)
			suite.Require().NoError(err)

			// set sequence to 2
			packet = types.NewPacket(2, path.EndpointA.ChainName, path.EndpointB.ChainName, relayChain, []string{xibctesting.MockPort}, [][]byte{validPacketData})
			err = path.EndpointA.SendPacket(packet)
			suite.Require().NoError(err)
		},
		true,
	}, {
		"port not found",
		func() {
			// use wrong port
			path := xibctesting.NewPath(suite.chainA, suite.chainB)
			suite.coordinator.SetupClients(path)
			packet = types.NewPacket(1, path.EndpointA.ChainName, path.EndpointB.ChainName, relayChain, []string{xibctesting.InvalidID}, [][]byte{validPacketData})
		},
		false,
	}, {
		"receipt already stored",
		func() {
			path := xibctesting.NewPath(suite.chainA, suite.chainB)
			suite.coordinator.SetupClients(path)
			packet = types.NewPacket(1, path.EndpointA.ChainName, path.EndpointB.ChainName, relayChain, []string{xibctesting.MockPort}, [][]byte{validPacketData})
			err := path.EndpointA.SendPacket(packet)
			suite.Require().NoError(err)
			suite.chainB.App.XIBCKeeper.PacketKeeper.SetPacketReceipt(suite.chainB.GetContext(), path.EndpointA.ChainName, path.EndpointB.ChainName, 1)
		},
		false,
	}, {
		"validation failed",
		func() {
			// packet commitment not set resulting in invalid proof
			path := xibctesting.NewPath(suite.chainA, suite.chainB)
			suite.coordinator.SetupClients(path)
			packet = types.NewPacket(1, path.EndpointA.ChainName, path.EndpointB.ChainName, relayChain, []string{xibctesting.InvalidID}, [][]byte{validPacketData})
		},
		false,
	}}

	for i, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s, %d/%d tests", tc.name, i, len(testCases)), func() {
			suite.SetupTest() // reset
			tc.malleate()

			// get proof of packet commitment from chainA
			packetKey := host.PacketCommitmentKey(packet.GetSourceChain(), packet.GetDestChain(), packet.GetSequence())
			proof, proofHeight := suite.chainA.QueryProof(packetKey)

			msg := &types.MsgRecvPacket{
				Packet:          packet,
				ProofCommitment: proof,
				ProofHeight:     proofHeight,
			}
			err := suite.chainB.App.XIBCKeeper.PacketKeeper.RecvPacket(suite.chainB.GetContext(), msg)

			if tc.expPass {
				suite.Require().NoError(err)

				receipt, receiptStored := suite.chainB.App.XIBCKeeper.PacketKeeper.GetPacketReceipt(
					suite.chainB.GetContext(), packet.GetSourceChain(), packet.GetDestChain(), packet.GetSequence(),
				)

				suite.Require().True(receiptStored, "packet receipt not stored after RecvPacket")
				suite.Require().Equal(string([]byte{byte(1)}), receipt, "packet receipt is not empty string")

			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestWriteAcknowledgement() {
	var (
		ack    []byte
		packet exported.PacketI
	)

	testCases := []testCase{{
		"success",
		func() {
			path := xibctesting.NewPath(suite.chainA, suite.chainB)
			suite.coordinator.SetupClients(path)
			packet = types.NewPacket(1, path.EndpointA.ChainName, path.EndpointB.ChainName, relayChain, []string{xibctesting.MockPort}, [][]byte{validPacketData})
			ack = xibctesting.TestHash
		},
		true,
	}, {
		"no-op, already acked",
		func() {
			path := xibctesting.NewPath(suite.chainA, suite.chainB)
			suite.coordinator.SetupClients(path)
			packet = types.NewPacket(1, path.EndpointA.ChainName, path.EndpointB.ChainName, relayChain, []string{xibctesting.MockPort}, [][]byte{validPacketData})
			ack = xibctesting.TestHash
			suite.chainB.App.XIBCKeeper.PacketKeeper.SetPacketAcknowledgement(suite.chainB.GetContext(), packet.GetSourceChain(), packet.GetDestChain(), packet.GetSequence(), ack)
		},
		false,
	}, {
		"empty acknowledgement",
		func() {
			path := xibctesting.NewPath(suite.chainA, suite.chainB)
			suite.coordinator.SetupClients(path)
			packet = types.NewPacket(1, path.EndpointA.ChainName, path.EndpointB.ChainName, relayChain, []string{xibctesting.MockPort}, [][]byte{validPacketData})
			ack = nil
		},
		false,
	}}

	for i, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s, %d/%d tests", tc.name, i, len(testCases)), func() {
			suite.SetupTest() // reset
			tc.malleate()

			err := suite.chainB.App.XIBCKeeper.PacketKeeper.WriteAcknowledgement(suite.chainB.GetContext(), packet, ack)
			if tc.expPass {
				suite.Require().NoError(err, "Invalid Case %d passed: %s", i, tc.name)
			} else {
				suite.Require().Error(err, "Case %d failed: %s", i, tc.name)
			}
		})
	}
}

// TestAcknowledgePacket tests the call AcknowledgePacket on chainA.
func (suite *KeeperTestSuite) TestAcknowledgePacket() {
	var (
		packet types.Packet
		ack    = xibcmock.MockAcknowledgement
	)

	testCases := []testCase{{
		"success",
		func() {
			// setup
			path := xibctesting.NewPath(suite.chainA, suite.chainB)
			suite.coordinator.SetupClients(path)
			packet = types.NewPacket(1, path.EndpointA.ChainName, path.EndpointB.ChainName, relayChain, []string{xibctesting.MockPort}, [][]byte{validPacketData})

			// create packet commitment
			err := path.EndpointA.SendPacket(packet)
			suite.Require().NoError(err)

			// create packet receipt and acknowledgement
			err = path.EndpointB.RecvPacket(packet)
			suite.Require().NoError(err)
		},
		true,
	}, {
		"port not found",
		func() {
			// use wrong port naming
			path := xibctesting.NewPath(suite.chainA, suite.chainB)
			suite.coordinator.SetupClients(path)
			packet = types.NewPacket(1, path.EndpointA.ChainName, path.EndpointB.ChainName, relayChain, []string{xibctesting.InvalidID}, [][]byte{validPacketData})
		},
		false,
	}, {
		"packet hasn't been sent",
		func() {
			// packet commitment never written
			path := xibctesting.NewPath(suite.chainA, suite.chainB)
			suite.coordinator.SetupClients(path)
			packet = types.NewPacket(1, path.EndpointA.ChainName, path.EndpointB.ChainName, relayChain, []string{xibctesting.MockPort}, [][]byte{validPacketData})
		},
		false,
	}}

	for i, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s, %d/%d tests", tc.name, i, len(testCases)), func() {
			suite.SetupTest() // reset
			tc.malleate()

			packetKey := host.PacketAcknowledgementKey(packet.GetSourceChain(), packet.GetDestChain(), packet.GetSequence())
			proof, proofHeight := suite.chainB.QueryProof(packetKey)

			msg := &types.MsgAcknowledgement{
				Packet:          packet,
				Acknowledgement: ack,
				ProofAcked:      proof,
				ProofHeight:     proofHeight,
			}
			err := suite.chainA.App.XIBCKeeper.PacketKeeper.AcknowledgePacket(suite.chainA.GetContext(), msg)
			pc := suite.chainA.App.XIBCKeeper.PacketKeeper.GetPacketCommitment(suite.chainA.GetContext(), packet.GetSourceChain(), packet.GetDestChain(), packet.GetSequence())

			if tc.expPass {
				suite.Require().NoError(err, "Case %d failed: %s", i, tc.name)
				suite.Nil(pc)
			} else {
				suite.Require().Error(err, "Invalid Case %d passed: %s", i, tc.name)
			}
		})
	}
}
