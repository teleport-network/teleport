package keeper

import (
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/teleport-network/teleport/x/rvesting/types"
)

type Keeper struct {
	paramSubspace paramtypes.Subspace
	BankKeeper    BankKeeper
	AccountKeeper AccountKeeper

	FeeCollectorName string // name of the FeeCollector ModuleAccount
}

func NewKeeper(paramSubspace paramtypes.Subspace, bk BankKeeper, ak AccountKeeper, feeCollectorName string) Keeper {
	// set KeyTable if it has not already been set
	if !paramSubspace.HasKeyTable() {
		paramSubspace = paramSubspace.WithKeyTable(types.ParamKeyTable())
	}
	return Keeper{
		paramSubspace:    paramSubspace,
		BankKeeper:       bk,
		AccountKeeper:    ak,
		FeeCollectorName: feeCollectorName,
	}
}
