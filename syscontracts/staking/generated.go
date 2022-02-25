// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package staking

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// StakingMetaData contains all meta data concerning the Staking contract.
var StakingMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"validator\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Delegated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"validatorSrc\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"validatorDest\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Redelegated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"validator\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Undelegated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"validator\",\"type\":\"string\"}],\"name\":\"Withdrew\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"validator\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"delegate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"validatorSrc\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"validatorDest\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"redelegate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"validator\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"undelegate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"validator\",\"type\":\"string\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"03f24de1": "delegate(string,uint256)",
		"7dd0209d": "redelegate(string,string,uint256)",
		"8dfc8897": "undelegate(string,uint256)",
		"31fb67c2": "withdraw(string)",
	},
	Bin: "0x608060405234801561001057600080fd5b50610444806100206000396000f3fe608060405234801561001057600080fd5b506004361061004c5760003560e01c806303f24de11461005157806331fb67c2146100665780637dd0209d146100795780638dfc88971461008c575b600080fd5b61006461005f366004610232565b61009f565b005b610064610074366004610277565b6100de565b6100646100873660046102b4565b61011a565b61006461009a366004610232565b61015c565b7fc181d211c1379e7ca130a707f3b1d49177b2c9eaca63f5b1c0fced957bd94d193383836040516100d29392919061036e565b60405180910390a15050565b7f271a84b5abc74645a8af43af4da7b3540bb0ac7603fbae9ff8b2e2df548332d0338260405161010f9291906103a2565b60405180910390a150565b7f1e4f99bac1ee5d1d13ed93a8febbb6730c1760e6b40b62f4971ecd57f184c20b3384848460405161014f94939291906103c6565b60405180910390a1505050565b7fb6be07748cf5d5075de9024d60c8dadb2346190b734318a49a330774ab1effc63383836040516100d29392919061036e565b634e487b7160e01b600052604160045260246000fd5b600082601f8301126101b657600080fd5b813567ffffffffffffffff808211156101d1576101d161018f565b604051601f8301601f19908116603f011681019082821181831017156101f9576101f961018f565b8160405283815286602085880101111561021257600080fd5b836020870160208301376000602085830101528094505050505092915050565b6000806040838503121561024557600080fd5b823567ffffffffffffffff81111561025c57600080fd5b610268858286016101a5565b95602094909401359450505050565b60006020828403121561028957600080fd5b813567ffffffffffffffff8111156102a057600080fd5b6102ac848285016101a5565b949350505050565b6000806000606084860312156102c957600080fd5b833567ffffffffffffffff808211156102e157600080fd5b6102ed878388016101a5565b9450602086013591508082111561030357600080fd5b50610310868287016101a5565b925050604084013590509250925092565b6000815180845260005b818110156103475760208185018101518683018201520161032b565b81811115610359576000602083870101525b50601f01601f19169290920160200192915050565b6001600160a01b038416815260606020820181905260009061039290830185610321565b9050826040830152949350505050565b6001600160a01b03831681526040602082018190526000906102ac90830184610321565b6001600160a01b03851681526080602082018190526000906103ea90830186610321565b82810360408401526103fc8186610321565b9150508260608301529594505050505056fea26469706673582212209d979531f80622188af0ffa8508cf70fd5ee27a13fdd6f0db103c6362053084364736f6c634300080c0033",
}

// StakingABI is the input ABI used to generate the binding from.
// Deprecated: Use StakingMetaData.ABI instead.
var StakingABI = StakingMetaData.ABI

// Deprecated: Use StakingMetaData.Sigs instead.
// StakingFuncSigs maps the 4-byte function signature to its string representation.
var StakingFuncSigs = StakingMetaData.Sigs

// StakingBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use StakingMetaData.Bin instead.
var StakingBin = StakingMetaData.Bin

// DeployStaking deploys a new Ethereum contract, binding an instance of Staking to it.
func DeployStaking(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Staking, error) {
	parsed, err := StakingMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(StakingBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Staking{StakingCaller: StakingCaller{contract: contract}, StakingTransactor: StakingTransactor{contract: contract}, StakingFilterer: StakingFilterer{contract: contract}}, nil
}

// Staking is an auto generated Go binding around an Ethereum contract.
type Staking struct {
	StakingCaller     // Read-only binding to the contract
	StakingTransactor // Write-only binding to the contract
	StakingFilterer   // Log filterer for contract events
}

// StakingCaller is an auto generated read-only Go binding around an Ethereum contract.
type StakingCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakingTransactor is an auto generated write-only Go binding around an Ethereum contract.
type StakingTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakingFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type StakingFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakingSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type StakingSession struct {
	Contract     *Staking          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StakingCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type StakingCallerSession struct {
	Contract *StakingCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// StakingTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type StakingTransactorSession struct {
	Contract     *StakingTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// StakingRaw is an auto generated low-level Go binding around an Ethereum contract.
type StakingRaw struct {
	Contract *Staking // Generic contract binding to access the raw methods on
}

// StakingCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type StakingCallerRaw struct {
	Contract *StakingCaller // Generic read-only contract binding to access the raw methods on
}

// StakingTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type StakingTransactorRaw struct {
	Contract *StakingTransactor // Generic write-only contract binding to access the raw methods on
}

// NewStaking creates a new instance of Staking, bound to a specific deployed contract.
func NewStaking(address common.Address, backend bind.ContractBackend) (*Staking, error) {
	contract, err := bindStaking(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Staking{StakingCaller: StakingCaller{contract: contract}, StakingTransactor: StakingTransactor{contract: contract}, StakingFilterer: StakingFilterer{contract: contract}}, nil
}

// NewStakingCaller creates a new read-only instance of Staking, bound to a specific deployed contract.
func NewStakingCaller(address common.Address, caller bind.ContractCaller) (*StakingCaller, error) {
	contract, err := bindStaking(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StakingCaller{contract: contract}, nil
}

// NewStakingTransactor creates a new write-only instance of Staking, bound to a specific deployed contract.
func NewStakingTransactor(address common.Address, transactor bind.ContractTransactor) (*StakingTransactor, error) {
	contract, err := bindStaking(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StakingTransactor{contract: contract}, nil
}

// NewStakingFilterer creates a new log filterer instance of Staking, bound to a specific deployed contract.
func NewStakingFilterer(address common.Address, filterer bind.ContractFilterer) (*StakingFilterer, error) {
	contract, err := bindStaking(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StakingFilterer{contract: contract}, nil
}

// bindStaking binds a generic wrapper to an already deployed contract.
func bindStaking(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(StakingABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Staking *StakingRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Staking.Contract.StakingCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Staking *StakingRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Staking.Contract.StakingTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Staking *StakingRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Staking.Contract.StakingTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Staking *StakingCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Staking.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Staking *StakingTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Staking.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Staking *StakingTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Staking.Contract.contract.Transact(opts, method, params...)
}

// Delegate is a paid mutator transaction binding the contract method 0x03f24de1.
//
// Solidity: function delegate(string validator, uint256 amount) returns()
func (_Staking *StakingTransactor) Delegate(opts *bind.TransactOpts, validator string, amount *big.Int) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "delegate", validator, amount)
}

// Delegate is a paid mutator transaction binding the contract method 0x03f24de1.
//
// Solidity: function delegate(string validator, uint256 amount) returns()
func (_Staking *StakingSession) Delegate(validator string, amount *big.Int) (*types.Transaction, error) {
	return _Staking.Contract.Delegate(&_Staking.TransactOpts, validator, amount)
}

// Delegate is a paid mutator transaction binding the contract method 0x03f24de1.
//
// Solidity: function delegate(string validator, uint256 amount) returns()
func (_Staking *StakingTransactorSession) Delegate(validator string, amount *big.Int) (*types.Transaction, error) {
	return _Staking.Contract.Delegate(&_Staking.TransactOpts, validator, amount)
}

// Redelegate is a paid mutator transaction binding the contract method 0x7dd0209d.
//
// Solidity: function redelegate(string validatorSrc, string validatorDest, uint256 amount) returns()
func (_Staking *StakingTransactor) Redelegate(opts *bind.TransactOpts, validatorSrc string, validatorDest string, amount *big.Int) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "redelegate", validatorSrc, validatorDest, amount)
}

// Redelegate is a paid mutator transaction binding the contract method 0x7dd0209d.
//
// Solidity: function redelegate(string validatorSrc, string validatorDest, uint256 amount) returns()
func (_Staking *StakingSession) Redelegate(validatorSrc string, validatorDest string, amount *big.Int) (*types.Transaction, error) {
	return _Staking.Contract.Redelegate(&_Staking.TransactOpts, validatorSrc, validatorDest, amount)
}

// Redelegate is a paid mutator transaction binding the contract method 0x7dd0209d.
//
// Solidity: function redelegate(string validatorSrc, string validatorDest, uint256 amount) returns()
func (_Staking *StakingTransactorSession) Redelegate(validatorSrc string, validatorDest string, amount *big.Int) (*types.Transaction, error) {
	return _Staking.Contract.Redelegate(&_Staking.TransactOpts, validatorSrc, validatorDest, amount)
}

// Undelegate is a paid mutator transaction binding the contract method 0x8dfc8897.
//
// Solidity: function undelegate(string validator, uint256 amount) returns()
func (_Staking *StakingTransactor) Undelegate(opts *bind.TransactOpts, validator string, amount *big.Int) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "undelegate", validator, amount)
}

// Undelegate is a paid mutator transaction binding the contract method 0x8dfc8897.
//
// Solidity: function undelegate(string validator, uint256 amount) returns()
func (_Staking *StakingSession) Undelegate(validator string, amount *big.Int) (*types.Transaction, error) {
	return _Staking.Contract.Undelegate(&_Staking.TransactOpts, validator, amount)
}

// Undelegate is a paid mutator transaction binding the contract method 0x8dfc8897.
//
// Solidity: function undelegate(string validator, uint256 amount) returns()
func (_Staking *StakingTransactorSession) Undelegate(validator string, amount *big.Int) (*types.Transaction, error) {
	return _Staking.Contract.Undelegate(&_Staking.TransactOpts, validator, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0x31fb67c2.
//
// Solidity: function withdraw(string validator) returns()
func (_Staking *StakingTransactor) Withdraw(opts *bind.TransactOpts, validator string) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "withdraw", validator)
}

// Withdraw is a paid mutator transaction binding the contract method 0x31fb67c2.
//
// Solidity: function withdraw(string validator) returns()
func (_Staking *StakingSession) Withdraw(validator string) (*types.Transaction, error) {
	return _Staking.Contract.Withdraw(&_Staking.TransactOpts, validator)
}

// Withdraw is a paid mutator transaction binding the contract method 0x31fb67c2.
//
// Solidity: function withdraw(string validator) returns()
func (_Staking *StakingTransactorSession) Withdraw(validator string) (*types.Transaction, error) {
	return _Staking.Contract.Withdraw(&_Staking.TransactOpts, validator)
}

// StakingDelegatedIterator is returned from FilterDelegated and is used to iterate over the raw logs and unpacked data for Delegated events raised by the Staking contract.
type StakingDelegatedIterator struct {
	Event *StakingDelegated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *StakingDelegatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingDelegated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(StakingDelegated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *StakingDelegatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingDelegatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingDelegated represents a Delegated event raised by the Staking contract.
type StakingDelegated struct {
	Delegator common.Address
	Validator string
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterDelegated is a free log retrieval operation binding the contract event 0xc181d211c1379e7ca130a707f3b1d49177b2c9eaca63f5b1c0fced957bd94d19.
//
// Solidity: event Delegated(address delegator, string validator, uint256 amount)
func (_Staking *StakingFilterer) FilterDelegated(opts *bind.FilterOpts) (*StakingDelegatedIterator, error) {

	logs, sub, err := _Staking.contract.FilterLogs(opts, "Delegated")
	if err != nil {
		return nil, err
	}
	return &StakingDelegatedIterator{contract: _Staking.contract, event: "Delegated", logs: logs, sub: sub}, nil
}

// WatchDelegated is a free log subscription operation binding the contract event 0xc181d211c1379e7ca130a707f3b1d49177b2c9eaca63f5b1c0fced957bd94d19.
//
// Solidity: event Delegated(address delegator, string validator, uint256 amount)
func (_Staking *StakingFilterer) WatchDelegated(opts *bind.WatchOpts, sink chan<- *StakingDelegated) (event.Subscription, error) {

	logs, sub, err := _Staking.contract.WatchLogs(opts, "Delegated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingDelegated)
				if err := _Staking.contract.UnpackLog(event, "Delegated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseDelegated is a log parse operation binding the contract event 0xc181d211c1379e7ca130a707f3b1d49177b2c9eaca63f5b1c0fced957bd94d19.
//
// Solidity: event Delegated(address delegator, string validator, uint256 amount)
func (_Staking *StakingFilterer) ParseDelegated(log types.Log) (*StakingDelegated, error) {
	event := new(StakingDelegated)
	if err := _Staking.contract.UnpackLog(event, "Delegated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingRedelegatedIterator is returned from FilterRedelegated and is used to iterate over the raw logs and unpacked data for Redelegated events raised by the Staking contract.
type StakingRedelegatedIterator struct {
	Event *StakingRedelegated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *StakingRedelegatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingRedelegated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(StakingRedelegated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *StakingRedelegatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingRedelegatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingRedelegated represents a Redelegated event raised by the Staking contract.
type StakingRedelegated struct {
	Delegator     common.Address
	ValidatorSrc  string
	ValidatorDest string
	Amount        *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterRedelegated is a free log retrieval operation binding the contract event 0x1e4f99bac1ee5d1d13ed93a8febbb6730c1760e6b40b62f4971ecd57f184c20b.
//
// Solidity: event Redelegated(address delegator, string validatorSrc, string validatorDest, uint256 amount)
func (_Staking *StakingFilterer) FilterRedelegated(opts *bind.FilterOpts) (*StakingRedelegatedIterator, error) {

	logs, sub, err := _Staking.contract.FilterLogs(opts, "Redelegated")
	if err != nil {
		return nil, err
	}
	return &StakingRedelegatedIterator{contract: _Staking.contract, event: "Redelegated", logs: logs, sub: sub}, nil
}

// WatchRedelegated is a free log subscription operation binding the contract event 0x1e4f99bac1ee5d1d13ed93a8febbb6730c1760e6b40b62f4971ecd57f184c20b.
//
// Solidity: event Redelegated(address delegator, string validatorSrc, string validatorDest, uint256 amount)
func (_Staking *StakingFilterer) WatchRedelegated(opts *bind.WatchOpts, sink chan<- *StakingRedelegated) (event.Subscription, error) {

	logs, sub, err := _Staking.contract.WatchLogs(opts, "Redelegated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingRedelegated)
				if err := _Staking.contract.UnpackLog(event, "Redelegated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRedelegated is a log parse operation binding the contract event 0x1e4f99bac1ee5d1d13ed93a8febbb6730c1760e6b40b62f4971ecd57f184c20b.
//
// Solidity: event Redelegated(address delegator, string validatorSrc, string validatorDest, uint256 amount)
func (_Staking *StakingFilterer) ParseRedelegated(log types.Log) (*StakingRedelegated, error) {
	event := new(StakingRedelegated)
	if err := _Staking.contract.UnpackLog(event, "Redelegated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingUndelegatedIterator is returned from FilterUndelegated and is used to iterate over the raw logs and unpacked data for Undelegated events raised by the Staking contract.
type StakingUndelegatedIterator struct {
	Event *StakingUndelegated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *StakingUndelegatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingUndelegated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(StakingUndelegated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *StakingUndelegatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingUndelegatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingUndelegated represents a Undelegated event raised by the Staking contract.
type StakingUndelegated struct {
	Delegator common.Address
	Validator string
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterUndelegated is a free log retrieval operation binding the contract event 0xb6be07748cf5d5075de9024d60c8dadb2346190b734318a49a330774ab1effc6.
//
// Solidity: event Undelegated(address delegator, string validator, uint256 amount)
func (_Staking *StakingFilterer) FilterUndelegated(opts *bind.FilterOpts) (*StakingUndelegatedIterator, error) {

	logs, sub, err := _Staking.contract.FilterLogs(opts, "Undelegated")
	if err != nil {
		return nil, err
	}
	return &StakingUndelegatedIterator{contract: _Staking.contract, event: "Undelegated", logs: logs, sub: sub}, nil
}

// WatchUndelegated is a free log subscription operation binding the contract event 0xb6be07748cf5d5075de9024d60c8dadb2346190b734318a49a330774ab1effc6.
//
// Solidity: event Undelegated(address delegator, string validator, uint256 amount)
func (_Staking *StakingFilterer) WatchUndelegated(opts *bind.WatchOpts, sink chan<- *StakingUndelegated) (event.Subscription, error) {

	logs, sub, err := _Staking.contract.WatchLogs(opts, "Undelegated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingUndelegated)
				if err := _Staking.contract.UnpackLog(event, "Undelegated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUndelegated is a log parse operation binding the contract event 0xb6be07748cf5d5075de9024d60c8dadb2346190b734318a49a330774ab1effc6.
//
// Solidity: event Undelegated(address delegator, string validator, uint256 amount)
func (_Staking *StakingFilterer) ParseUndelegated(log types.Log) (*StakingUndelegated, error) {
	event := new(StakingUndelegated)
	if err := _Staking.contract.UnpackLog(event, "Undelegated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingWithdrewIterator is returned from FilterWithdrew and is used to iterate over the raw logs and unpacked data for Withdrew events raised by the Staking contract.
type StakingWithdrewIterator struct {
	Event *StakingWithdrew // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *StakingWithdrewIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingWithdrew)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(StakingWithdrew)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *StakingWithdrewIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingWithdrewIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingWithdrew represents a Withdrew event raised by the Staking contract.
type StakingWithdrew struct {
	Delegator common.Address
	Validator string
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterWithdrew is a free log retrieval operation binding the contract event 0x271a84b5abc74645a8af43af4da7b3540bb0ac7603fbae9ff8b2e2df548332d0.
//
// Solidity: event Withdrew(address delegator, string validator)
func (_Staking *StakingFilterer) FilterWithdrew(opts *bind.FilterOpts) (*StakingWithdrewIterator, error) {

	logs, sub, err := _Staking.contract.FilterLogs(opts, "Withdrew")
	if err != nil {
		return nil, err
	}
	return &StakingWithdrewIterator{contract: _Staking.contract, event: "Withdrew", logs: logs, sub: sub}, nil
}

// WatchWithdrew is a free log subscription operation binding the contract event 0x271a84b5abc74645a8af43af4da7b3540bb0ac7603fbae9ff8b2e2df548332d0.
//
// Solidity: event Withdrew(address delegator, string validator)
func (_Staking *StakingFilterer) WatchWithdrew(opts *bind.WatchOpts, sink chan<- *StakingWithdrew) (event.Subscription, error) {

	logs, sub, err := _Staking.contract.WatchLogs(opts, "Withdrew")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingWithdrew)
				if err := _Staking.contract.UnpackLog(event, "Withdrew", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseWithdrew is a log parse operation binding the contract event 0x271a84b5abc74645a8af43af4da7b3540bb0ac7603fbae9ff8b2e2df548332d0.
//
// Solidity: event Withdrew(address delegator, string validator)
func (_Staking *StakingFilterer) ParseWithdrew(log types.Log) (*StakingWithdrew, error) {
	event := new(StakingWithdrew)
	if err := _Staking.contract.UnpackLog(event, "Withdrew", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
