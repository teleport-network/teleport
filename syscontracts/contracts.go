package syscontracts

import (
	_ "embed" // embed compiled smart contract
)

const (
	StakingContractAddress   = "0x0000000000000000000000000000000010000001"
	GovContractAddress       = "0x0000000000000000000000000000000010000002"
	TransferContractAddress  = "0x0000000000000000000000000000000010000003"
	RCCContractAddress       = "0x0000000000000000000000000000000010000004"
	MultiCallContractAddress = "0x0000000000000000000000000000000010000005"
	WTELEContractAddress     = "0x0000000000000000000000000000000010000006"
	AgentContractAddress     = "0x0000000000000000000000000000000010000007"
	PacketContractAddress    = "0x0000000000000000000000000000000010000008"
)

var (
	//go:embed contracts_compiled/Staking.json
	StakingJSON []byte // nolint: golint

	//go:embed contracts_compiled/Gov.json
	GovJSON []byte // nolint: golint

	//go:embed contracts_compiled/ERC20Burnable.json
	ERC20BurnableJSON []byte

	//go:embed contracts_compiled/ERC20DirectBalanceManipulation.json
	ERC20DirectBalanceManipulationJSON []byte // nolint: golint

	//go:embed contracts_compiled/ERC20MaliciousDelayed.json
	ERC20MaliciousDelayedJSON []byte // nolint: golint

	//go:embed contracts_compiled/ERC20MinterBurnerDecimals.json
	ERC20MinterBurnerDecimalsJSON []byte // nolint: golint

	//go:embed contracts_compiled/WTELE.json
	WTELEJSON []byte // nolint: golint
)
