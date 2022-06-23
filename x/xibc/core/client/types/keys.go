package types

const (
	// SubModuleName defines the XIBC client name
	SubModuleName string = "client"

	// RouterKey is the gov route for XIBC client
	RouterKey string = "xibcclient"

	// KeyClientName is the key used to store the chain name in the keeper.
	KeyClientName = "chainName"

	// KeyRelayers is the key used to store the relayers address in the keeper.
	KeyRelayers = "relayers"
)
