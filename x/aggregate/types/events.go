package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// aggregate events
const (
	EventTypeTokenLock             = "token_lock"
	EventTypeTokenUnlock           = "token_unlock"
	EventTypeMint                  = "mint"
	EventTypeConvertCoin           = "convert_coin"
	EventTypeConvertERC20          = "convert_erc20"
	EventTypeBurn                  = "burn"
	EventTypeRegisterCoin          = "register_coin"
	EventTypeRegisterERC20         = "register_erc20"
	EventTypeToggleTokenConversion = "toggle_token_conversion" // #nosec
	EventTypeRegisterERC20Trace    = "register_erc20_trace"

	AttributeKeyCosmosCoin  = "cosmos_coin"
	AttributeKeyERC20Token  = "erc20_token" // #nosec
	AttributeKeyReceiver    = "receiver"
	AttributeKeyOriginToken = "origin_token"
	AttributeKeyOriginChain = "origin_chain"

	ERC20EventTransfer = "Transfer"
)

// Event type for Transfer(address from, address to, uint256 value)
type LogTransfer struct {
	From   common.Address
	To     common.Address
	Tokens *big.Int
}
