package syscontracts

import (
	"fmt"

	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

func ParseLog(out interface{}, abi *ethabi.ABI, log *ethtypes.Log, event string) (err error) {
	if log.Topics[0] != abi.Events[event].ID {
		return fmt.Errorf("event signature mismatch")
	}
	if len(log.Data) > 0 {
		if err = abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return err
		}
	}
	var indexed ethabi.Arguments
	for _, arg := range abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	return ethabi.ParseTopics(out, indexed, log.Topics[1:])
}
