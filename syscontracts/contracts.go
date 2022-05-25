package syscontracts

import (
	_ "embed" // embed compiled smart contract
)

const (
	// system contracts
	StakingContractAddress = "0x0000000000000000000000000000000010000001"
	GovContractAddress     = "0x0000000000000000000000000000000010000002"

	// xibc core contracts
	PacketContractAddress = "0x0000000000000000000000000000000020000001"
	CrossChainAddress = "0x0000000000000000000000000000000020000002"

	// app contracts
	AgentContractAddress = "0x0000000000000000000000000000000040000001"

	// token contracts
	WTELEContractAddress = "0x0000000000000000000000000000000050000001"
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
