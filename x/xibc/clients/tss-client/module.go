package tssclient

import (
	"github.com/teleport-network/teleport/x/xibc/clients/tss-client/types"
)

// Name returns the XIBC tendermint client name
func Name() string {
	return types.SubModuleName
}
