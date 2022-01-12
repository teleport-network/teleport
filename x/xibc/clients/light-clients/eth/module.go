package eth

import "github.com/teleport-network/teleport/x/xibc/clients/light-clients/eth/types"

// Name returns the BIBC eth client name
func Name() string {
	return types.SubModuleName
}
