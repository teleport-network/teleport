package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/teleport-network/teleport/x/xibc/core/packet/types"
)

const (
	testChain1 = "firstchain"
	testChain2 = "secondchain"
)

func TestValidateGenesis(t *testing.T) {
	testCases := []struct {
		name     string
		genState types.GenesisState
		expPass  bool
	}{{
		name:     "default",
		genState: types.DefaultGenesisState(),
		expPass:  true,
	}, {
		name: "valid genesis",
		genState: types.NewGenesisState(
			[]types.PacketState{
				types.NewPacketState(testChain1, testChain2, 1, []byte("ack")),
			},
			[]types.PacketState{
				types.NewPacketState(testChain1, testChain2, 1, []byte("commit_hash")),
			},
			[]types.PacketState{
				types.NewPacketState(testChain1, testChain2, 1, []byte("")),
			},
			[]types.PacketSequence{
				types.NewPacketSequence(testChain1, testChain2, 1),
			},
			[]types.PacketSequence{
				types.NewPacketSequence(testChain1, testChain2, 1),
			},
			[]types.PacketSequence{
				types.NewPacketSequence(testChain1, testChain2, 1),
			},
		),
		expPass: true,
	}, {
		name: "invalid ack",
		genState: types.GenesisState{
			Acknowledgements: []types.PacketState{
				types.NewPacketState(testChain1, testChain2, 1, nil),
			},
		},
		expPass: false,
	}, {
		name: "invalid commitment",
		genState: types.GenesisState{
			Commitments: []types.PacketState{
				types.NewPacketState(testChain1, testChain2, 1, nil),
			},
		},
		expPass: false,
	}, {
		name: "invalid send seq",
		genState: types.GenesisState{
			SendSequences: []types.PacketSequence{
				types.NewPacketSequence(testChain1, testChain2, 0),
			},
		},
		expPass: false,
	}, {
		name: "invalid recv seq",
		genState: types.GenesisState{
			RecvSequences: []types.PacketSequence{
				types.NewPacketSequence(testChain1, "(testCha1)", 1),
			},
		},
		expPass: false,
	}, {
		name: "invalid recv seq 2",
		genState: types.GenesisState{
			RecvSequences: []types.PacketSequence{
				types.NewPacketSequence("(testChain1)", testChain2, 1),
			},
		},
		expPass: false,
	}, {
		name: "invalid ack seq",
		genState: types.GenesisState{
			AckSequences: []types.PacketSequence{
				types.NewPacketSequence(testChain1, "(testChain2)", 1),
			},
		},
		expPass: false,
	}}

	for _, tc := range testCases {
		if tc.expPass {
			require.NoError(t, tc.genState.Validate(), tc.name)
		} else {
			require.Error(t, tc.genState.Validate(), tc.name)
		}
	}
}
