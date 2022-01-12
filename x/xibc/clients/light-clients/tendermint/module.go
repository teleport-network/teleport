package tendermint

import (
	"github.com/teleport-network/teleport/x/xibc/clients/light-clients/tendermint/types"
)

// Name returns the XIBC tendermint client name
func Name() string {
	return types.SubModuleName
}
