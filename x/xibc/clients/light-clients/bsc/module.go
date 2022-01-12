package bsc

import "github.com/teleport-network/teleport/x/xibc/clients/light-clients/bsc/types"

// Name returns the BIBC bsc client name
func Name() string {
	return types.SubModuleName
}
