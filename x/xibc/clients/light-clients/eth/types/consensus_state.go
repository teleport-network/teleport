package types

import (
	"github.com/teleport-network/teleport/x/xibc/exported"
)

var _ exported.ConsensusState = (*ConsensusState)(nil)

func (m *ConsensusState) ClientType() string {
	return exported.BSC
}

func (m *ConsensusState) GetRoot() []byte {
	return m.Root
}

func (m *ConsensusState) GetTimestamp() uint64 {
	return m.Timestamp
}

func (m *ConsensusState) ValidateBasic() error {
	return nil
}
