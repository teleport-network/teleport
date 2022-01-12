package types

const (
	// SubModuleName defines the XIBC client name
	SubModuleName string = "client"

	// GovRouterKey is the gov route for XIBC client
	GovRouterKey string = SubModuleName

	// KeyClientName is the key used to store the chain name in the keeper.
	KeyClientName = "chainName"

	// KeyRelayers is the key used to store the relayers address in the keeper.
	KeyRelayers = "relayers"
)
