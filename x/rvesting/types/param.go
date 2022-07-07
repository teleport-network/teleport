package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/bitdao-io/bitnetwork/types"
)

// Parameter keys
var (
	KeyEnableVesting  = []byte("EnableVesting")
	KeyPerBlockReward = []byte("PerBlockReward")
)

var _ paramtypes.ParamSet = &Params{}

func validatePerBlockReward(r interface{}) error {
	reward, ok := r.(sdk.Coins)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", r)
	}
	if len(reward) == 0 {
		return fmt.Errorf("invalid per block reward: %v", reward)
	}
	for _, rr := range reward {
		if len(rr.Denom) == 0 {
			return fmt.Errorf("denom of per block reward can not be empty")
		}
		if rr.IsNegative() {
			return fmt.Errorf("invalid per block reward: %v", rr)
		}
	}
	return nil
}

// ParamSetPairs implements the ParamSet interface and returns all the key/value pairs
// pairs of auth module's parameters.
func (m *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(
			KeyEnableVesting,
			&m.EnableVesting,
			func(value interface{}) error { return nil },
		),
		paramtypes.NewParamSetPair(KeyPerBlockReward, &m.PerBlockReward, validatePerBlockReward),
	}
}

func (m *Params) validate() error {
	if m.EnableVesting {
		return validatePerBlockReward(m.PerBlockReward)
	}
	return nil
}

func DefaultParams() Params {
	return Params{
		EnableVesting:  false,
		PerBlockReward: sdk.NewCoins(types.NewBitCoin(sdk.NewIntWithDecimal(1, 17))),
	}
}

// ParamKeyTable returns the parameter key table.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}
