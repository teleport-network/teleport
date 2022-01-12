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
	evmkeeper "github.com/tharsis/ethermint/x/evm/keeper"

	adcommon "github.com/teleport-network/teleport/adapter/common"
	"github.com/teleport-network/teleport/syscontracts"
	"github.com/teleport-network/teleport/syscontracts/gov"
)

type HookAdapter struct {
	accountKeeper *authkeeper.AccountKeeper
	evmKeeper     *evmkeeper.Keeper
	router        *baseapp.MsgServiceRouter
	abi           *ethabi.ABI
	govContract   common.Address
	handlers      map[common.Hash]adcommon.EvmEventHandler
}

func (h HookAdapter) InitGenesis(ctx sdk.Context) error {
	codeHash := crypto.Keccak256Hash(common.FromHex(syscontracts.GovContractCode))
	h.evmKeeper.SetCode(ctx, codeHash.Bytes(), common.FromHex(syscontracts.GovContractCode))

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
	evmKeeper *evmkeeper.Keeper,
	router *baseapp.MsgServiceRouter,
) *HookAdapter {
	parsed, err := ethabi.JSON(strings.NewReader(gov.GovMetaData.ABI))
	if err != nil {
		panic(err)
	}
	handlers := make(map[common.Hash]adcommon.EvmEventHandler, len(parsed.Events))

	hook := &HookAdapter{
		accountKeeper: accountKeeper,
		evmKeeper:     evmKeeper,
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
