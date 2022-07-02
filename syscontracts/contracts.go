package syscontracts

import (
	_ "embed" // embed compiled smart contract
)

const (
	// system contracts
	StakingContractAddress = "0x0000000000000000000000000000000010000001"
	GovContractAddress     = "0x0000000000000000000000000000000010000002"
)

var (
	//go:embed contracts_compiled/Staking.json
	StakingJSON []byte // nolint: golint

	//go:embed contracts_compiled/Gov.json
	GovJSON []byte // nolint: golint
)
