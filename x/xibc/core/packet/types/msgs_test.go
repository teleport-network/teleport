package types_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/store/iavl"
	"github.com/cosmos/cosmos-sdk/store/rootmulti"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/teleport-network/teleport/app"
	clienttypes "github.com/teleport-network/teleport/x/xibc/core/client/types"
	commitmenttypes "github.com/teleport-network/teleport/x/xibc/core/commitment/types"
	"github.com/teleport-network/teleport/x/xibc/core/packet/types"
)

const (
	// valid constatns used for testing
	sourceChain = "source-chain"
	destChain   = "dest-chain"
	relayChain  = ""
	port        = "mock"
)

// define variables used for testing
var (
	height            = clienttypes.NewHeight(0, 1)
	validPacketData   = []byte("testdata")
	unknownPacketData = []byte("unknown")

	packet        = types.NewPacket(1, sourceChain, destChain, relayChain, []string{port}, [][]byte{validPacketData})
	invalidPacket = types.NewPacket(0, sourceChain, destChain, relayChain, []string{port}, [][]byte{unknownPacketData})

	addr      = sdk.AccAddress("testaddr111111111111")
	emptyAddr sdk.AccAddress
)

type TypesTestSuite struct {
	suite.Suite

	proof []byte
}

func (suite *TypesTestSuite) SetupTest() {
	teleport := app.Setup(false, nil)
	db := dbm.NewMemDB()
	store := rootmulti.NewStore(db)
	storeKey := storetypes.NewKVStoreKey("iavlStoreKey")

	store.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, nil)
	_ = store.LoadVersion(0)
	iavlStore := store.GetCommitStore(storeKey).(*iavl.Store)

	iavlStore.Set([]byte("KEY"), []byte("VALUE"))
	_ = store.Commit()

	res := store.Query(abci.RequestQuery{
		Path:  fmt.Sprintf("/%s/key", storeKey.Name()), // required path to get key/value+proof
		Data:  []byte("KEY"),
		Prove: true,
	})

	merkleProof, err := commitmenttypes.ConvertProofs(res.ProofOps)
	suite.Require().NoError(err)
	proof, err := teleport.AppCodec().Marshal(&merkleProof)
	suite.Require().NoError(err)

	suite.proof = proof
}

func TestTypesTestSuite(t *testing.T) {
	suite.Run(t, new(TypesTestSuite))
}

func (suite *TypesTestSuite) TestMsgRecvPacketValidateBasic() {
	testCases := []struct {
		name    string
		msg     *types.MsgRecvPacket
		expPass bool
	}{
		{"success", types.NewMsgRecvPacket(packet, suite.proof, height, addr), true},
		{"proof height is zero", types.NewMsgRecvPacket(packet, suite.proof, clienttypes.ZeroHeight(), addr), false},
		{"missing signer address", types.NewMsgRecvPacket(packet, suite.proof, height, emptyAddr), false},
		{"invalid packet", types.NewMsgRecvPacket(invalidPacket, suite.proof, height, addr), false},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			if tc.expPass {
				suite.Require().NoError(tc.msg.ValidateBasic())
			} else {
				suite.Require().Error(tc.msg.ValidateBasic())
			}
		})
	}
}

func (suite *TypesTestSuite) TestMsgRecvPacketGetSigners() {
	msg := types.NewMsgRecvPacket(packet, suite.proof, height, addr)
	res := msg.GetSigners()

	expected := "[7465737461646472313131313131313131313131]"
	suite.Require().Equal(expected, fmt.Sprintf("%v", res))
}

func (suite *TypesTestSuite) TestMsgAcknowledgementValidateBasic() {
	testCases := []struct {
		name    string
		msg     *types.MsgAcknowledgement
		expPass bool
	}{
		{"success", types.NewMsgAcknowledgement(packet, packet.GetDataList()[0], suite.proof, height, addr), true},
		{"proof height must be > 0", types.NewMsgAcknowledgement(packet, packet.GetDataList()[0], suite.proof, clienttypes.ZeroHeight(), addr), false},
		{"empty ack", types.NewMsgAcknowledgement(packet, nil, suite.proof, height, addr), false},
		{"missing signer address", types.NewMsgAcknowledgement(packet, packet.GetDataList()[0], suite.proof, height, emptyAddr), false},
		{"invalid packet", types.NewMsgAcknowledgement(invalidPacket, packet.GetDataList()[0], suite.proof, height, addr), false},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			if tc.expPass {
				suite.Require().NoError(tc.msg.ValidateBasic())
			} else {
				suite.Require().Error(tc.msg.ValidateBasic())
			}
		})
	}
}
