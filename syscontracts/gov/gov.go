// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package gov

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

// GovOptionWeight is an auto generated low-level Go binding around an user-defined struct.
type GovOptionWeight struct {
	Option uint32
	Weight uint64
}

// GovMetaData contains all meta data concerning the Gov contract.
var GovMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"voter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"proposalId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"voteOption\",\"type\":\"uint32\"}],\"name\":\"Voted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"voter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"proposalId\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"option\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"weight\",\"type\":\"uint64\"}],\"indexed\":false,\"internalType\":\"structGov.OptionWeight[]\",\"name\":\"options\",\"type\":\"tuple[]\"}],\"name\":\"VotedWeighted\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"proposalId\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"voteOption\",\"type\":\"uint32\"}],\"name\":\"vote\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"proposalId\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"option\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"weight\",\"type\":\"uint64\"}],\"internalType\":\"structGov.OptionWeight[]\",\"name\":\"options\",\"type\":\"tuple[]\"}],\"name\":\"vote\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b506105d1806100206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c8063645566e51461003b578063a653afab14610057575b600080fd5b61005560048036038101906100509190610181565b610073565b005b610071600480360381019061006c919061036f565b6100b2565b005b7fd14efccbe77f09e3139993753f0e0d883c70f9c377196f3a1d57ec20a6d942973383836040516100a69392919061042a565b60405180910390a15050565b7f472a83d63c58b68a77d3fd3476709b328eac4e058c6c4fd88c9446f19782cc7b3383836040516100e59392919061055d565b60405180910390a15050565b6000604051905090565b600080fd5b600080fd5b600067ffffffffffffffff82169050919050565b61012281610105565b811461012d57600080fd5b50565b60008135905061013f81610119565b92915050565b600063ffffffff82169050919050565b61015e81610145565b811461016957600080fd5b50565b60008135905061017b81610155565b92915050565b60008060408385031215610198576101976100fb565b5b60006101a685828601610130565b92505060206101b78582860161016c565b9150509250929050565b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b61020f826101c6565b810181811067ffffffffffffffff8211171561022e5761022d6101d7565b5b80604052505050565b60006102416100f1565b905061024d8282610206565b919050565b600067ffffffffffffffff82111561026d5761026c6101d7565b5b602082029050602081019050919050565b600080fd5b600080fd5b60006040828403121561029e5761029d610283565b5b6102a86040610237565b905060006102b88482850161016c565b60008301525060206102cc84828501610130565b60208301525092915050565b60006102eb6102e684610252565b610237565b9050808382526020820190506040840283018581111561030e5761030d61027e565b5b835b8181101561033757806103238882610288565b845260208401935050604081019050610310565b5050509392505050565b600082601f830112610356576103556101c1565b5b81356103668482602086016102d8565b91505092915050565b60008060408385031215610386576103856100fb565b5b600061039485828601610130565b925050602083013567ffffffffffffffff8111156103b5576103b4610100565b5b6103c185828601610341565b9150509250929050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006103f6826103cb565b9050919050565b610406816103eb565b82525050565b61041581610105565b82525050565b61042481610145565b82525050565b600060608201905061043f60008301866103fd565b61044c602083018561040c565b610459604083018461041b565b949350505050565b600081519050919050565b600082825260208201905092915050565b6000819050602082019050919050565b61049681610145565b82525050565b6104a581610105565b82525050565b6040820160008201516104c1600085018261048d565b5060208201516104d4602085018261049c565b50505050565b60006104e683836104ab565b60408301905092915050565b6000602082019050919050565b600061050a82610461565b610514818561046c565b935061051f8361047d565b8060005b8381101561055057815161053788826104da565b9750610542836104f2565b925050600181019050610523565b5085935050505092915050565b600060608201905061057260008301866103fd565b61057f602083018561040c565b818103604083015261059181846104ff565b905094935050505056fea2646970667358221220144e2bb698566eeacbd08aaf2126f412595dc40c7b22f26297bfca10db12d40764736f6c634300080a0033",
}

// GovABI is the input ABI used to generate the binding from.
// Deprecated: Use GovMetaData.ABI instead.
var GovABI = GovMetaData.ABI

// GovBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use GovMetaData.Bin instead.
var GovBin = GovMetaData.Bin

// DeployGov deploys a new Ethereum contract, binding an instance of Gov to it.
func DeployGov(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Gov, error) {
	parsed, err := GovMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(GovBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Gov{GovCaller: GovCaller{contract: contract}, GovTransactor: GovTransactor{contract: contract}, GovFilterer: GovFilterer{contract: contract}}, nil
}

// Gov is an auto generated Go binding around an Ethereum contract.
type Gov struct {
	GovCaller     // Read-only binding to the contract
	GovTransactor // Write-only binding to the contract
	GovFilterer   // Log filterer for contract events
}

// GovCaller is an auto generated read-only Go binding around an Ethereum contract.
type GovCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GovTransactor is an auto generated write-only Go binding around an Ethereum contract.
type GovTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GovFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type GovFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GovSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type GovSession struct {
	Contract     *Gov              // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// GovCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type GovCallerSession struct {
	Contract *GovCaller    // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// GovTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type GovTransactorSession struct {
	Contract     *GovTransactor    // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// GovRaw is an auto generated low-level Go binding around an Ethereum contract.
type GovRaw struct {
	Contract *Gov // Generic contract binding to access the raw methods on
}

// GovCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type GovCallerRaw struct {
	Contract *GovCaller // Generic read-only contract binding to access the raw methods on
}

// GovTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type GovTransactorRaw struct {
	Contract *GovTransactor // Generic write-only contract binding to access the raw methods on
}

// NewGov creates a new instance of Gov, bound to a specific deployed contract.
func NewGov(address common.Address, backend bind.ContractBackend) (*Gov, error) {
	contract, err := bindGov(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Gov{GovCaller: GovCaller{contract: contract}, GovTransactor: GovTransactor{contract: contract}, GovFilterer: GovFilterer{contract: contract}}, nil
}

// NewGovCaller creates a new read-only instance of Gov, bound to a specific deployed contract.
func NewGovCaller(address common.Address, caller bind.ContractCaller) (*GovCaller, error) {
	contract, err := bindGov(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &GovCaller{contract: contract}, nil
}

// NewGovTransactor creates a new write-only instance of Gov, bound to a specific deployed contract.
func NewGovTransactor(address common.Address, transactor bind.ContractTransactor) (*GovTransactor, error) {
	contract, err := bindGov(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &GovTransactor{contract: contract}, nil
}

// NewGovFilterer creates a new log filterer instance of Gov, bound to a specific deployed contract.
func NewGovFilterer(address common.Address, filterer bind.ContractFilterer) (*GovFilterer, error) {
	contract, err := bindGov(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &GovFilterer{contract: contract}, nil
}

// bindGov binds a generic wrapper to an already deployed contract.
func bindGov(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(GovABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Gov *GovRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Gov.Contract.GovCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Gov *GovRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Gov.Contract.GovTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Gov *GovRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Gov.Contract.GovTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Gov *GovCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Gov.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Gov *GovTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Gov.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Gov *GovTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Gov.Contract.contract.Transact(opts, method, params...)
}

// Vote is a paid mutator transaction binding the contract method 0x645566e5.
//
// Solidity: function vote(uint64 proposalId, uint32 voteOption) returns()
func (_Gov *GovTransactor) Vote(opts *bind.TransactOpts, proposalId uint64, voteOption uint32) (*types.Transaction, error) {
	return _Gov.contract.Transact(opts, "vote", proposalId, voteOption)
}

// Vote is a paid mutator transaction binding the contract method 0x645566e5.
//
// Solidity: function vote(uint64 proposalId, uint32 voteOption) returns()
func (_Gov *GovSession) Vote(proposalId uint64, voteOption uint32) (*types.Transaction, error) {
	return _Gov.Contract.Vote(&_Gov.TransactOpts, proposalId, voteOption)
}

// Vote is a paid mutator transaction binding the contract method 0x645566e5.
//
// Solidity: function vote(uint64 proposalId, uint32 voteOption) returns()
func (_Gov *GovTransactorSession) Vote(proposalId uint64, voteOption uint32) (*types.Transaction, error) {
	return _Gov.Contract.Vote(&_Gov.TransactOpts, proposalId, voteOption)
}

// Vote0 is a paid mutator transaction binding the contract method 0xa653afab.
//
// Solidity: function vote(uint64 proposalId, (uint32,uint64)[] options) returns()
func (_Gov *GovTransactor) Vote0(opts *bind.TransactOpts, proposalId uint64, options []GovOptionWeight) (*types.Transaction, error) {
	return _Gov.contract.Transact(opts, "vote0", proposalId, options)
}

// Vote0 is a paid mutator transaction binding the contract method 0xa653afab.
//
// Solidity: function vote(uint64 proposalId, (uint32,uint64)[] options) returns()
func (_Gov *GovSession) Vote0(proposalId uint64, options []GovOptionWeight) (*types.Transaction, error) {
	return _Gov.Contract.Vote0(&_Gov.TransactOpts, proposalId, options)
}

// Vote0 is a paid mutator transaction binding the contract method 0xa653afab.
//
// Solidity: function vote(uint64 proposalId, (uint32,uint64)[] options) returns()
func (_Gov *GovTransactorSession) Vote0(proposalId uint64, options []GovOptionWeight) (*types.Transaction, error) {
	return _Gov.Contract.Vote0(&_Gov.TransactOpts, proposalId, options)
}

// GovVotedIterator is returned from FilterVoted and is used to iterate over the raw logs and unpacked data for Voted events raised by the Gov contract.
type GovVotedIterator struct {
	Event *GovVoted // Event containing the contract specifics and raw log

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
func (it *GovVotedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GovVoted)
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
		it.Event = new(GovVoted)
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
func (it *GovVotedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GovVotedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GovVoted represents a Voted event raised by the Gov contract.
type GovVoted struct {
	Voter      common.Address
	ProposalId uint64
	VoteOption uint32
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterVoted is a free log retrieval operation binding the contract event 0xd14efccbe77f09e3139993753f0e0d883c70f9c377196f3a1d57ec20a6d94297.
//
// Solidity: event Voted(address voter, uint64 proposalId, uint32 voteOption)
func (_Gov *GovFilterer) FilterVoted(opts *bind.FilterOpts) (*GovVotedIterator, error) {

	logs, sub, err := _Gov.contract.FilterLogs(opts, "Voted")
	if err != nil {
		return nil, err
	}
	return &GovVotedIterator{contract: _Gov.contract, event: "Voted", logs: logs, sub: sub}, nil
}

// WatchVoted is a free log subscription operation binding the contract event 0xd14efccbe77f09e3139993753f0e0d883c70f9c377196f3a1d57ec20a6d94297.
//
// Solidity: event Voted(address voter, uint64 proposalId, uint32 voteOption)
func (_Gov *GovFilterer) WatchVoted(opts *bind.WatchOpts, sink chan<- *GovVoted) (event.Subscription, error) {

	logs, sub, err := _Gov.contract.WatchLogs(opts, "Voted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GovVoted)
				if err := _Gov.contract.UnpackLog(event, "Voted", log); err != nil {
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

// ParseVoted is a log parse operation binding the contract event 0xd14efccbe77f09e3139993753f0e0d883c70f9c377196f3a1d57ec20a6d94297.
//
// Solidity: event Voted(address voter, uint64 proposalId, uint32 voteOption)
func (_Gov *GovFilterer) ParseVoted(log types.Log) (*GovVoted, error) {
	event := new(GovVoted)
	if err := _Gov.contract.UnpackLog(event, "Voted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GovVotedWeightedIterator is returned from FilterVotedWeighted and is used to iterate over the raw logs and unpacked data for VotedWeighted events raised by the Gov contract.
type GovVotedWeightedIterator struct {
	Event *GovVotedWeighted // Event containing the contract specifics and raw log

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
func (it *GovVotedWeightedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GovVotedWeighted)
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
		it.Event = new(GovVotedWeighted)
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
func (it *GovVotedWeightedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GovVotedWeightedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GovVotedWeighted represents a VotedWeighted event raised by the Gov contract.
type GovVotedWeighted struct {
	Voter      common.Address
	ProposalId uint64
	Options    []GovOptionWeight
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterVotedWeighted is a free log retrieval operation binding the contract event 0x472a83d63c58b68a77d3fd3476709b328eac4e058c6c4fd88c9446f19782cc7b.
//
// Solidity: event VotedWeighted(address voter, uint64 proposalId, (uint32,uint64)[] options)
func (_Gov *GovFilterer) FilterVotedWeighted(opts *bind.FilterOpts) (*GovVotedWeightedIterator, error) {

	logs, sub, err := _Gov.contract.FilterLogs(opts, "VotedWeighted")
	if err != nil {
		return nil, err
	}
	return &GovVotedWeightedIterator{contract: _Gov.contract, event: "VotedWeighted", logs: logs, sub: sub}, nil
}

// WatchVotedWeighted is a free log subscription operation binding the contract event 0x472a83d63c58b68a77d3fd3476709b328eac4e058c6c4fd88c9446f19782cc7b.
//
// Solidity: event VotedWeighted(address voter, uint64 proposalId, (uint32,uint64)[] options)
func (_Gov *GovFilterer) WatchVotedWeighted(opts *bind.WatchOpts, sink chan<- *GovVotedWeighted) (event.Subscription, error) {

	logs, sub, err := _Gov.contract.WatchLogs(opts, "VotedWeighted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GovVotedWeighted)
				if err := _Gov.contract.UnpackLog(event, "VotedWeighted", log); err != nil {
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

// ParseVotedWeighted is a log parse operation binding the contract event 0x472a83d63c58b68a77d3fd3476709b328eac4e058c6c4fd88c9446f19782cc7b.
//
// Solidity: event VotedWeighted(address voter, uint64 proposalId, (uint32,uint64)[] options)
func (_Gov *GovFilterer) ParseVotedWeighted(log types.Log) (*GovVotedWeighted, error) {
	event := new(GovVotedWeighted)
	if err := _Gov.contract.UnpackLog(event, "VotedWeighted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
