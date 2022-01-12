package staking

import (
	"errors"
	"strings"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"

	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	ethermint "github.com/tharsis/ethermint/types"
	evmkeeper "github.com/tharsis/ethermint/x/evm/keeper"

	adcommon "github.com/teleport-network/teleport/adapter/common"
	"github.com/teleport-network/teleport/syscontracts"
	"github.com/teleport-network/teleport/syscontracts/staking"
)

type HookAdapter struct {
	accountKeeper   *authkeeper.AccountKeeper
	stakingkeeper   *stakingkeeper.Keeper
	evmKeeper       *evmkeeper.Keeper
	router          *baseapp.MsgServiceRouter
	abi             *ethabi.ABI
	stakingContract common.Address
	handlers        map[common.Hash]adcommon.EvmEventHandler
}

func (h HookAdapter) InitGenesis(ctx sdk.Context) error {
	codeHash := crypto.Keccak256Hash(common.FromHex(syscontracts.StakingContractCode))
	h.evmKeeper.SetCode(ctx, codeHash.Bytes(), common.FromHex(syscontracts.StakingContractCode))

	account := h.accountKeeper.NewAccountWithAddress(ctx, common.HexToAddress(syscontracts.StakingContractAddress).Bytes())
	ethAccount := account.(*ethermint.EthAccount)

	ethAccount.CodeHash = codeHash.Hex()
	h.accountKeeper.SetAccount(ctx, ethAccount)

	return nil
}

func (h HookAdapter) Name() string {
	return "staking"
}

func NewHookAdapter(
	accountKeeper *authkeeper.AccountKeeper,
	stakingKeeper *stakingkeeper.Keeper,
	evmKeeper *evmkeeper.Keeper,
	router *baseapp.MsgServiceRouter,
) *HookAdapter {
	parsed, err := ethabi.JSON(strings.NewReader(staking.StakingMetaData.ABI))
	if err != nil {
		panic(err)
	}
	handlers := make(map[common.Hash]adcommon.EvmEventHandler, len(parsed.Events))

	hook := &HookAdapter{
		accountKeeper:   accountKeeper,
		stakingkeeper:   stakingKeeper,
		evmKeeper:       evmKeeper,
		router:          router,
		handlers:        handlers,
		abi:             &parsed,
		stakingContract: common.HexToAddress(syscontracts.StakingContractAddress),
	}
	for name, event := range parsed.Events {
		switch name {
		case "Delegated":
			handlers[event.ID] = hook.HandleDelegated
		case "Undelegated":
			handlers[event.ID] = hook.HandleUndelegated
		case "Redelegated":
			handlers[event.ID] = hook.HandleRedelegated
		case "Withdrew":
			handlers[event.ID] = hook.HandleWithdrew
		default:
			panic(errors.New("unknown topic"))
		}
	}
	return hook
}
