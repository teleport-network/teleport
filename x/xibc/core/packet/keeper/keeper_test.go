package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/teleport-network/teleport/x/xibc/core/packet/types"
	xibctesting "github.com/teleport-network/teleport/x/xibc/testing"
)

// KeeperTestSuite is a testing suite to test keeper functions.
type KeeperTestSuite struct {
	suite.Suite

	coordinator *xibctesting.Coordinator

	// testing chains used for convenience and readability
	chainA *xibctesting.TestChain
	chainB *xibctesting.TestChain
}

// TestKeeperTestSuite runs all the tests within this package.
func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

// SetupTest creates a coordinator with 2 test chains.
func (suite *KeeperTestSuite) SetupTest() {
	suite.coordinator = xibctesting.NewCoordinator(suite.T(), 2)
	suite.chainA = suite.coordinator.GetChain(xibctesting.GetChainID(0))
	suite.chainB = suite.coordinator.GetChain(xibctesting.GetChainID(1))
	// commit some blocks so that QueryProof returns valid proof (cannot return valid query if height <= 1)
	suite.coordinator.CommitNBlocks(suite.chainA, 2)
	suite.coordinator.CommitNBlocks(suite.chainB, 2)
}

// TestGetAllSequences sets all packet sequences
func (suite KeeperTestSuite) TestGetAllSequences() {
	path := xibctesting.NewPath(suite.chainA, suite.chainB)
	suite.coordinator.SetupClients(path)

	seq1 := types.NewPacketSequence(path.EndpointA.ChainName, path.EndpointB.ChainName, 1)
	seq2 := types.NewPacketSequence(path.EndpointA.ChainName, path.EndpointB.ChainName, 2)

	// seq1 should be overwritten by seq2
	expSeqs := []types.PacketSequence{seq2}

	ctxA := suite.chainA.GetContext()

	for _, seq := range []types.PacketSequence{seq1, seq2} {
		suite.chainA.App.XIBCKeeper.PacketKeeper.SetNextSequenceSend(ctxA, path.EndpointA.ChainName, path.EndpointB.ChainName, seq.Sequence)
	}

	sendSeqs := suite.chainA.App.XIBCKeeper.PacketKeeper.GetAllPacketSendSeqs(ctxA)
	suite.Len(sendSeqs, 1)

	suite.Require().Equal(expSeqs, sendSeqs)
}

// TestGetAllPacketState creates a set of acks, packet commitments, and receipts
func (suite KeeperTestSuite) TestGetAllPacketState() {
	path := xibctesting.NewPath(suite.chainA, suite.chainB)
	suite.coordinator.SetupClients(path)

	ack1 := types.NewPacketState(path.EndpointA.ChainName, path.EndpointB.ChainName, 1, []byte("ack"))
	ack2 := types.NewPacketState(path.EndpointA.ChainName, path.EndpointB.ChainName, 2, []byte("ack"))

	ack2dup := types.NewPacketState(path.EndpointA.ChainName, path.EndpointB.ChainName, 2, []byte("ack"))

	receipt := string([]byte{byte(1)})
	rec1 := types.NewPacketState(path.EndpointA.ChainName, path.EndpointB.ChainName, 1, []byte(receipt))
	rec2 := types.NewPacketState(path.EndpointA.ChainName, path.EndpointB.ChainName, 2, []byte(receipt))

	comm1 := types.NewPacketState(path.EndpointA.ChainName, path.EndpointB.ChainName, 1, []byte("hash"))
	comm2 := types.NewPacketState(path.EndpointA.ChainName, path.EndpointB.ChainName, 2, []byte("hash"))

	expAcks := []types.PacketState{ack1, ack2}
	expReceipts := []types.PacketState{rec1, rec2}
	expCommitments := []types.PacketState{comm1, comm2}

	ctxA := suite.chainA.GetContext()

	// set acknowledgements
	for _, ack := range []types.PacketState{ack1, ack2, ack2dup} {
		suite.chainA.App.XIBCKeeper.PacketKeeper.SetPacketAcknowledgement(ctxA, path.EndpointA.ChainName, path.EndpointB.ChainName, ack.Sequence, ack.Data)
	}

	// set packet receipts
	for _, rec := range expReceipts {
		suite.chainA.App.XIBCKeeper.PacketKeeper.SetPacketReceipt(ctxA, path.EndpointA.ChainName, path.EndpointB.ChainName, rec.Sequence)
	}

	// set packet commitments
	for _, comm := range expCommitments {
		suite.chainA.App.XIBCKeeper.PacketKeeper.SetPacketCommitment(ctxA, path.EndpointA.ChainName, path.EndpointB.ChainName, comm.Sequence, comm.Data)
	}

	acks := suite.chainA.App.XIBCKeeper.PacketKeeper.GetAllPacketAcks(ctxA)
	receipts := suite.chainA.App.XIBCKeeper.PacketKeeper.GetAllPacketReceipts(ctxA)
	commitments := suite.chainA.App.XIBCKeeper.PacketKeeper.GetAllPacketCommitments(ctxA)

	suite.Require().Len(acks, len(expAcks))
	suite.Require().Len(commitments, len(expCommitments))
	suite.Require().Len(receipts, len(expReceipts))

	suite.Require().Equal(expAcks, acks)
	suite.Require().Equal(expReceipts, receipts)
	suite.Require().Equal(expCommitments, commitments)
}

// TestSetSequence verifies that the keeper correctly sets the sequence counters.
func (suite *KeeperTestSuite) TestSetSequence() {
	path := xibctesting.NewPath(suite.chainA, suite.chainB)
	suite.coordinator.SetupClients(path)

	ctxA := suite.chainA.GetContext()
	one := uint64(1)

	seq := suite.chainA.App.XIBCKeeper.PacketKeeper.GetNextSequenceSend(ctxA, path.EndpointA.ChainName, path.EndpointB.ChainName)
	suite.Require().Equal(one, seq)

	nextSeqSend := uint64(10)
	suite.chainA.App.XIBCKeeper.PacketKeeper.SetNextSequenceSend(ctxA, path.EndpointA.ChainName, path.EndpointB.ChainName, nextSeqSend)

	storedNextSeqSend := suite.chainA.App.XIBCKeeper.PacketKeeper.GetNextSequenceSend(ctxA, path.EndpointA.ChainName, path.EndpointB.ChainName)
	suite.Require().Equal(nextSeqSend, storedNextSeqSend)
}

func (suite *KeeperTestSuite) TestGetAllPacketCommitmentsByPath() {
	path := xibctesting.NewPath(suite.chainA, suite.chainB)
	suite.coordinator.SetupClients(path)

	ctxA := suite.chainA.GetContext()
	expectedSeqs := make(map[uint64]bool)
	hash := []byte("commitment")

	seq := uint64(15)
	maxSeq := uint64(25)
	suite.Require().Greater(maxSeq, seq)

	// create consecutive commitments
	for i := uint64(1); i < seq; i++ {
		suite.chainA.App.XIBCKeeper.PacketKeeper.SetPacketCommitment(ctxA, path.EndpointA.ChainName, path.EndpointB.ChainName, i, hash)
		expectedSeqs[i] = true
	}

	// add non-consecutive commitments
	for i := seq; i < maxSeq; i += 2 {
		suite.chainA.App.XIBCKeeper.PacketKeeper.SetPacketCommitment(ctxA, path.EndpointA.ChainName, path.EndpointB.ChainName, i, hash)
		expectedSeqs[i] = true
	}

	suite.chainA.App.XIBCKeeper.PacketKeeper.SetPacketCommitment(ctxA, path.EndpointA.ChainName, "EndpointBChainName", maxSeq+1, hash)

	commitments := suite.chainA.App.XIBCKeeper.PacketKeeper.GetAllPacketCommitmentsByPath(ctxA, path.EndpointA.ChainName, path.EndpointB.ChainName)

	suite.Require().Equal(len(expectedSeqs), len(commitments))
	// ensure above for loops occurred
	suite.Require().NotEqual(0, len(commitments))

	// verify that all the packet commitments were stored
	for _, packet := range commitments {
		suite.Require().True(expectedSeqs[packet.Sequence])
		suite.Require().Equal(path.EndpointA.ChainName, packet.SourceChain)
		suite.Require().Equal(path.EndpointB.ChainName, packet.DestinationChain)
		suite.Require().Equal(hash, packet.Data)

		// prevent duplicates from passing checks
		expectedSeqs[packet.Sequence] = false
	}
}

// TestSetPacketAcknowledgement verifies that packet acknowledgements are correctly set in the keeper
func (suite *KeeperTestSuite) TestSetPacketAcknowledgement() {
	path := xibctesting.NewPath(suite.chainA, suite.chainB)
	suite.coordinator.SetupClients(path)

	ctxA := suite.chainA.GetContext()
	seq := uint64(10)

	storedAckHash, found := suite.chainA.App.XIBCKeeper.PacketKeeper.GetPacketAcknowledgement(ctxA, path.EndpointA.ChainName, path.EndpointB.ChainName, seq)
	suite.Require().False(found)
	suite.Require().Nil(storedAckHash)

	ackHash := []byte("ackhash")
	suite.chainA.App.XIBCKeeper.PacketKeeper.SetPacketAcknowledgement(ctxA, path.EndpointA.ChainName, path.EndpointB.ChainName, seq, ackHash)

	storedAckHash, found = suite.chainA.App.XIBCKeeper.PacketKeeper.GetPacketAcknowledgement(ctxA, path.EndpointA.ChainName, path.EndpointB.ChainName, seq)
	suite.Require().True(found)
	suite.Require().Equal(ackHash, storedAckHash)
	suite.Require().True(suite.chainA.App.XIBCKeeper.PacketKeeper.HasPacketAcknowledgement(ctxA, path.EndpointA.ChainName, path.EndpointB.ChainName, seq))
}
