package gov

import (
	"errors"
	"strings"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"

	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	ethermint "github.com/tharsis/ethermint/types"
	evmtypes "github.com/tharsis/ethermint/x/evm/types"

	adcommon "github.com/teleport-network/teleport/adapter/common"
	"github.com/teleport-network/teleport/syscontracts"
	govcontract "github.com/teleport-network/teleport/syscontracts/gov"
)

var _ evmtypes.EvmHooks = &HookAdapter{}

type HookAdapter struct {
	accountKeeper *authkeeper.AccountKeeper
	router        *baseapp.MsgServiceRouter
	abi           *ethabi.ABI
	govContract   common.Address
	handlers      map[common.Hash]adcommon.EvmEventHandler
}

func (h HookAdapter) InitGenesis(ctx sdk.Context) error {
	codeHash := crypto.Keccak256Hash(govcontract.GovContract.Bin)

	account := h.accountKeeper.NewAccountWithAddress(ctx, common.HexToAddress(syscontracts.GovContractAddress).Bytes())
	ethAccount := account.(*ethermint.EthAccount)

	ethAccount.CodeHash = codeHash.Hex()
	h.accountKeeper.SetAccount(ctx, ethAccount)

	return nil
}

func (h HookAdapter) Name() string {
	return "gov"
}

func NewHookAdapter(
	accountKeeper *authkeeper.AccountKeeper,
	router *baseapp.MsgServiceRouter,
) *HookAdapter {
	parsed, err := ethabi.JSON(strings.NewReader(govcontract.GovMetaData.ABI))
	if err != nil {
		panic(err)
	}
	handlers := make(map[common.Hash]adcommon.EvmEventHandler, len(parsed.Events))

	hook := &HookAdapter{
		accountKeeper: accountKeeper,
		router:        router,
		abi:           &parsed,
		govContract:   common.HexToAddress(syscontracts.GovContractAddress),
		handlers:      handlers,
	}
	for name, event := range parsed.Events {
		switch name {
		case "Voted":
			handlers[event.ID] = hook.HandleVoted
		case "VotedWeighted":
			handlers[event.ID] = hook.HandleVotedWeighted
		default:
			panic(errors.New("unknown topic"))
		}
	}
	return hook
}
