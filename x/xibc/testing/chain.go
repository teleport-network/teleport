package xibctesting

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	packetcontract "github.com/teleport-network/teleport/syscontracts/xibc_packet"
	packettypes "github.com/teleport-network/teleport/x/xibc/core/packet/types"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/tmhash"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmprotoversion "github.com/tendermint/tendermint/proto/tendermint/version"
	tmtypes "github.com/tendermint/tendermint/types"
	"github.com/tendermint/tendermint/version"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/staking/teststaking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/ethereum/go-ethereum/common"

	"github.com/tharsis/ethermint/crypto/ethsecp256k1"
	"github.com/tharsis/ethermint/encoding"
	evm "github.com/tharsis/ethermint/x/evm/types"

	"github.com/teleport-network/teleport/app"
	teletypes "github.com/teleport-network/teleport/types"
	xibctmtypes "github.com/teleport-network/teleport/x/xibc/clients/light-clients/tendermint/types"
	clienttypes "github.com/teleport-network/teleport/x/xibc/core/client/types"
	commitmenttypes "github.com/teleport-network/teleport/x/xibc/core/commitment/types"
	"github.com/teleport-network/teleport/x/xibc/core/host"
	"github.com/teleport-network/teleport/x/xibc/exported"
	"github.com/teleport-network/teleport/x/xibc/testing/mock"
	"github.com/teleport-network/teleport/x/xibc/types"
)

const (
	// Default params constants used to create a TM client
	TrustingPeriod     time.Duration = time.Hour * 24 * 7 * 2
	UnbondingPeriod    time.Duration = time.Hour * 24 * 7 * 3
	MaxClockDrift      time.Duration = time.Second * 10
	DefaultDelayPeriod uint64        = 0

	InvalidID = "InvalidID"

	// used for testing UpdateClientProposal
	Title       = "title"
	Description = "description"
)

var (

	// Default params variables used to create a TM client
	DefaultTrustLevel xibctmtypes.Fraction = xibctmtypes.DefaultTrustLevel
	TestHash                               = tmhash.Sum([]byte("TESTING HASH"))
	TestCoin                               = sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100))
	Prefix                                 = commitmenttypes.MerklePrefix{KeyPrefix: []byte("xibc")}
)

// TestChain is a testing struct that wraps a simapp with the last TM Header, the current ABCI
// header and the validators of the TestChain. It also contains a field called ChainID. This
// is the chainName that *other* chains use to refer to this TestChain. The SenderAddress
// is used for delivering transactions through the application state.
// NOTE: the actual application uses an empty chain-id for ease of testing.
type TestChain struct {
	t *testing.T

	Coordinator    *Coordinator
	App            *app.Teleport
	ChainID        string
	LastHeader     *xibctmtypes.Header // header for last block height committed
	CurrentHeader  tmproto.Header      // header for current block height
	QueryServer    types.QueryServer
	QueryClientEvm evm.QueryClient
	TxConfig       client.TxConfig
	Codec          codec.BinaryCodec

	Vals    *tmtypes.ValidatorSet
	Signers []tmtypes.PrivValidator

	SenderPrivKey cryptotypes.PrivKey
	SenderAcc     sdk.AccAddress
	SenderAddress common.Address
}

// NewTestChain initializes a new TestChain instance with a single validator set using a
// generated private key. It also creates a sender account to be used for delivering transactions.
//
// The first block height is committed to state in order to allow for client creations on
// counterparty chains. The TestChain will return with a block height starting at 2.
//
// Time management is handled by the Coordinator in order to ensure synchrony between chains.
// Each update of any chain increments the block header time for all chains by 5 seconds.
func NewTestChain(t *testing.T, coord *Coordinator, chainID string) *TestChain {

	// set DefaultPowerReduction to the teleport PowerReduction
	sdk.DefaultPowerReduction = teletypes.PowerReduction

	// generate validator private/public key
	privVal := mock.NewPV()
	pubKey, err := privVal.GetPubKey()
	require.NoError(t, err)

	// create validator set with single validator
	validator := tmtypes.NewValidator(pubKey, 1)
	valSet := tmtypes.NewValidatorSet([]*tmtypes.Validator{validator})
	signers := []tmtypes.PrivValidator{privVal}

	// generate genesis account
	senderPrivKey, _ := ethsecp256k1.GenerateKey()
	senderAddress := common.BytesToAddress(senderPrivKey.PubKey().Address().Bytes())
	acc := authtypes.NewBaseAccount(senderPrivKey.PubKey().Address().Bytes(), senderPrivKey.PubKey(), 0, 0)
	balance := banktypes.Balance{
		Address: acc.GetAddress().String(),
		Coins:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100000000000000))),
	}

	teleport := SetupWithGenesisValSet(t, valSet, []authtypes.GenesisAccount{acc}, balance)

	// create current header and call begin block
	header := tmproto.Header{
		Height:          1,
		ChainID:         chainID,
		Time:            coord.CurrentTime.UTC(),
		ProposerAddress: valSet.Proposer.Address,
	}

	txConfig := encoding.MakeConfig(app.ModuleBasics).TxConfig

	queryHelperEvm := baseapp.NewQueryServerTestHelper(
		teleport.BaseApp.NewContext(false, header),
		teleport.InterfaceRegistry(),
	)
	evm.RegisterQueryServer(queryHelperEvm, teleport.EvmKeeper)

	// create an account to send transactions from
	chain := &TestChain{
		t:              t,
		Coordinator:    coord,
		ChainID:        chainID,
		App:            teleport,
		CurrentHeader:  header,
		QueryServer:    teleport.XIBCKeeper,
		QueryClientEvm: evm.NewQueryClient(queryHelperEvm),
		TxConfig:       txConfig,
		Codec:          teleport.AppCodec(),
		Vals:           valSet,
		Signers:        signers,
		SenderPrivKey:  senderPrivKey,
		SenderAcc:      acc.GetAddress(),
		SenderAddress:  senderAddress,
	}
	chain.App.XIBCKeeper.ClientKeeper.SetChainName(chain.GetContext(), chainID)
	coord.CommitBlock(chain)

	chain.SetPacketChainName()
	return chain
}

// NewTestChain initializes a new TestChain instance with a single validator set using a
// given private key. It also creates a sender account to be used for delivering transactions.
func NewTestChainWithAccount(t *testing.T, coord *Coordinator, chainID string, senderPrivKey *ethsecp256k1.PrivKey) *TestChain {

	// set DefaultPowerReduction to the teleport PowerReduction
	sdk.DefaultPowerReduction = teletypes.PowerReduction

	// generate validator private/public key
	privVal := mock.NewPV()
	pubKey, err := privVal.GetPubKey()
	require.NoError(t, err)

	// create validator set with single validator
	validator := tmtypes.NewValidator(pubKey, 1)
	valSet := tmtypes.NewValidatorSet([]*tmtypes.Validator{validator})
	signers := []tmtypes.PrivValidator{privVal}

	// generate genesis account
	senderAddress := common.BytesToAddress(senderPrivKey.PubKey().Address().Bytes())
	acc := authtypes.NewBaseAccount(senderPrivKey.PubKey().Address().Bytes(), senderPrivKey.PubKey(), 0, 0)
	balance := banktypes.Balance{
		Address: acc.GetAddress().String(),
		Coins:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100000000000000))),
	}

	teleport := SetupWithGenesisValSet(t, valSet, []authtypes.GenesisAccount{acc}, balance)

	// create current header and call begin block
	header := tmproto.Header{
		Height:          1,
		ChainID:         chainID,
		Time:            coord.CurrentTime.UTC(),
		ProposerAddress: valSet.Proposer.Address,
	}

	txConfig := encoding.MakeConfig(app.ModuleBasics).TxConfig

	queryHelperEvm := baseapp.NewQueryServerTestHelper(
		teleport.BaseApp.NewContext(false, header),
		teleport.InterfaceRegistry(),
	)
	evm.RegisterQueryServer(queryHelperEvm, teleport.EvmKeeper)

	// create an account to send transactions from
	chain := &TestChain{
		t:              t,
		Coordinator:    coord,
		ChainID:        chainID,
		App:            teleport,
		CurrentHeader:  header,
		QueryServer:    teleport.XIBCKeeper,
		QueryClientEvm: evm.NewQueryClient(queryHelperEvm),
		TxConfig:       txConfig,
		Codec:          teleport.AppCodec(),
		Vals:           valSet,
		Signers:        signers,
		SenderPrivKey:  senderPrivKey,
		SenderAcc:      acc.GetAddress(),
		SenderAddress:  senderAddress,
	}
	chain.App.XIBCKeeper.ClientKeeper.SetChainName(chain.GetContext(), chainID)
	coord.CommitBlock(chain)
	chain.SetPacketChainName()
	return chain
}

// GetContext returns the current context for the application.
func (chain *TestChain) GetContext() sdk.Context {
	return chain.App.BaseApp.NewContext(false, chain.CurrentHeader)
}

// QueryProof performs an abci query with the given key and returns the proto encoded merkle proof
// for the query and the height at which the proof will succeed on a tendermint verifier.
func (chain *TestChain) QueryProof(key []byte) ([]byte, clienttypes.Height) {
	return chain.QueryProofAtHeight(key, chain.App.LastBlockHeight())
}

// QueryProofAtHeight performs an abci query with the given key and returns the proto encoded merkle proof
// for the query and the height at which the proof will succeed on a tendermint verifier.
func (chain *TestChain) QueryProofAtHeight(key []byte, height int64) ([]byte, clienttypes.Height) {
	res := chain.App.Query(abci.RequestQuery{
		Path:   fmt.Sprintf("store/%s/key", host.StoreKey),
		Height: height - 1,
		Data:   key,
		Prove:  true,
	})

	merkleProof, err := commitmenttypes.ConvertProofs(res.ProofOps)
	require.NoError(chain.t, err)

	proof, err := chain.App.AppCodec().Marshal(&merkleProof)
	require.NoError(chain.t, err)

	revision := clienttypes.ParseChainID(chain.ChainID)

	// proof height + 1 is returned as the proof created corresponds to the height the proof
	// was created in the IAVL tree. Tendermint and subsequently the clients that rely on it
	// have heights 1 above the IAVL tree. Thus we return proof height + 1
	return proof, clienttypes.NewHeight(revision, uint64(res.Height)+1)
}

// QueryUpgradeProof performs an abci query with the given key and returns the proto encoded merkle proof
// for the query and the height at which the proof will succeed on a tendermint verifier.
func (chain *TestChain) QueryUpgradeProof(key []byte, height uint64) ([]byte, clienttypes.Height) {
	res := chain.App.Query(abci.RequestQuery{
		Path:   "store/upgrade/key",
		Height: int64(height - 1),
		Data:   key,
		Prove:  true,
	})

	merkleProof, err := commitmenttypes.ConvertProofs(res.ProofOps)
	require.NoError(chain.t, err)

	proof, err := chain.App.AppCodec().Marshal(&merkleProof)
	require.NoError(chain.t, err)

	revision := clienttypes.ParseChainID(chain.ChainID)

	// proof height + 1 is returned as the proof created corresponds to the height the proof
	// was created in the IAVL tree. Tendermint and subsequently the clients that rely on it
	// have heights 1 above the IAVL tree. Thus we return proof height + 1
	return proof, clienttypes.NewHeight(revision, uint64(res.Height+1))
}

// QueryClientStateProof performs and abci query for a client state
// stored with a given chainName and returns the ClientState along with the proof
func (chain *TestChain) QueryClientStateProof(chainName string) (exported.ClientState, []byte) {
	// retrieve client state to provide proof for
	clientState, found := chain.App.XIBCKeeper.ClientKeeper.GetClientState(chain.GetContext(), chainName)
	require.True(chain.t, found)

	clientKey := host.FullClientStateKey(chainName)
	proofClient, _ := chain.QueryProof(clientKey)

	return clientState, proofClient
}

// QueryConsensusStateProof performs an abci query for a consensus state
// stored on the given chainName. The proof and consensusHeight are returned.
func (chain *TestChain) QueryConsensusStateProof(chainName string) ([]byte, clienttypes.Height) {
	clientState := chain.GetClientState(chainName)

	consensusHeight := clientState.GetLatestHeight().(clienttypes.Height)
	consensusKey := host.FullConsensusStateKey(chainName, consensusHeight)
	proofConsensus, _ := chain.QueryProof(consensusKey)

	return proofConsensus, consensusHeight
}

// NextBlock sets the last header to the current header and increments the current header to be
// at the next block height. It does not update the time as that is handled by the Coordinator.
//
// CONTRACT: this function must only be called after app.Commit() occurs
func (chain *TestChain) NextBlock() {
	// set the last header to the current header
	// use nil trusted fields
	chain.LastHeader = chain.CurrentTMClientHeader()

	// increment the current header
	chain.CurrentHeader = tmproto.Header{
		ChainID: chain.ChainID,
		Height:  chain.App.LastBlockHeight() + 1,
		AppHash: chain.App.LastCommitID().Hash,
		// NOTE: the time is increased by the coordinator to maintain time synchrony amongst chains.
		Time:               chain.CurrentHeader.Time,
		ValidatorsHash:     chain.Vals.Hash(),
		NextValidatorsHash: chain.Vals.Hash(),
		ProposerAddress:    chain.Vals.Proposer.Address,
	}

	chain.App.BeginBlock(abci.RequestBeginBlock{Header: chain.CurrentHeader})
}

// sendMsgs delivers a transaction through the application without returning the result.
func (chain *TestChain) sendMsgs(msgs ...sdk.Msg) error {
	_, err := chain.SendMsgs(msgs...)
	return err
}

// SendMsgs delivers a transaction through the application. It updates the senders sequence
// number and updates the TestChain's headers. It returns the result and error if one
// occurred.
func (chain *TestChain) SendMsgs(msgs ...sdk.Msg) (*sdk.Result, error) {
	// ensure the chain has the latest time
	chain.Coordinator.UpdateTimeForChain(chain)

	account := chain.App.AccountKeeper.GetAccount(chain.GetContext(), chain.SenderAcc)

	_, r, err := app.SignAndDeliver(
		chain.t,
		chain.TxConfig,
		chain.App.BaseApp,
		chain.GetContext().BlockHeader(),
		msgs,
		chain.ChainID,
		[]uint64{account.GetAccountNumber()},
		[]uint64{account.GetSequence()},
		true,
		true,
		chain.SenderPrivKey,
	)
	if err != nil {
		return nil, err
	}

	// SignAndDeliver calls app.Commit()
	chain.NextBlock()

	chain.Coordinator.IncrementTime()

	return r, nil
}

// GetClientState retrieves the client state for the provided chainName. The client is
// expected to exist otherwise testing will fail.
func (chain *TestChain) GetClientState(chainName string) exported.ClientState {
	clientState, found := chain.App.XIBCKeeper.ClientKeeper.GetClientState(chain.GetContext(), chainName)
	require.True(chain.t, found)

	return clientState
}

// GetConsensusState retrieves the consensus state for the provided chainName and height.
// It will return a success boolean depending on if consensus state exists or not.
func (chain *TestChain) GetConsensusState(chainName string, height exported.Height) (exported.ConsensusState, bool) {
	return chain.App.XIBCKeeper.ClientKeeper.GetClientConsensusState(chain.GetContext(), chainName, height)
}

// GetValsAtHeight will return the validator set of the chain at a given height. It will return
// a success boolean depending on if the validator set exists or not at that height.
func (chain *TestChain) GetValsAtHeight(height int64) (*tmtypes.ValidatorSet, bool) {
	histInfo, ok := chain.App.StakingKeeper.GetHistoricalInfo(chain.GetContext(), height)
	if !ok {
		return nil, false
	}

	valSet := stakingtypes.Validators(histInfo.Valset)

	tmValidators, err := teststaking.ToTmValidators(valSet, sdk.DefaultPowerReduction)
	if err != nil {
		panic(err)
	}
	return tmtypes.NewValidatorSet(tmValidators), true
}

// GetAcknowledgement retrieves an acknowledgement for the provided packet. If the
// acknowledgement does not exist then testing will fail.
func (chain *TestChain) GetAcknowledgement(packet exported.PacketI) []byte {
	ack, found := chain.App.XIBCKeeper.PacketKeeper.GetPacketAcknowledgement(
		chain.GetContext(),
		packet.GetSrcChain(),
		packet.GetDstChain(),
		packet.GetSequence(),
	)
	require.True(chain.t, found)

	return ack
}

// GetPrefix returns the prefix for used by a chain
func (chain *TestChain) GetPrefix() commitmenttypes.MerklePrefix {
	return commitmenttypes.NewMerklePrefix([]byte(""))
}

// ConstructMsgCreateClient constructs a message to create a new client state (tendermint or solomachine).
// NOTE: a solo machine client will be created with an empty diversifier.
func (chain *TestChain) ConstructMsgCreateClient(counterparty *TestChain, chainName string, clientType string) error {
	var (
		clientState    exported.ClientState
		consensusState exported.ConsensusState
	)

	switch clientType {
	case exported.Tendermint:
		height := counterparty.LastHeader.GetHeight().(clienttypes.Height)
		clientState = xibctmtypes.NewClientState(
			counterparty.ChainID, DefaultTrustLevel,
			TrustingPeriod, UnbondingPeriod, MaxClockDrift,
			height, commitmenttypes.GetSDKSpecs(), Prefix, 0,
		)
		consensusState = counterparty.LastHeader.ConsensusState()
	default:
		chain.t.Fatalf("unsupported client state type %s", clientType)
	}

	err := chain.App.XIBCKeeper.ClientKeeper.CreateClient(
		chain.GetContext(),
		chainName,
		clientState, consensusState,
	)

	require.NoError(chain.t, err)
	return err
}

// CreateTMClient will construct and execute a tendermint MsgCreateClient. A counterparty
// client will be created on the (target) chain.
func (chain *TestChain) CreateTMClient(counterparty *TestChain, chainName string) error {
	// construct MsgCreateClient using counterparty
	return chain.ConstructMsgCreateClient(counterparty, chainName, exported.Tendermint)
}

// UpdateTMClient will construct and execute a tendermint MsgUpdateClient. The counterparty
// client will be updated on the (target) chain. UpdateTMClient mocks the relayer flow
// necessary for updating a Tendermint client.
func (chain *TestChain) UpdateTMClient(counterparty *TestChain, chainName string) error {
	header, err := chain.ConstructUpdateTMClientHeader(counterparty, chainName)
	require.NoError(chain.t, err)

	msg, err := clienttypes.NewMsgUpdateClient(
		chainName, header,
		chain.SenderAcc,
	)
	require.NoError(chain.t, err)

	return chain.sendMsgs(msg)
}

// ConstructUpdateTMClientHeader will construct a valid tendermint Header to update the
// light client on the source chain.
func (chain *TestChain) ConstructUpdateTMClientHeader(counterparty *TestChain, chainName string) (*xibctmtypes.Header, error) {
	return chain.ConstructUpdateTMClientHeaderWithTrustedHeight(counterparty, chainName, clienttypes.ZeroHeight())
}

// ConstructUpdateTMClientHeaderWithTrustedHeight will construct a valid tendermint Header to update the
// light client on the source chain.
func (chain *TestChain) ConstructUpdateTMClientHeaderWithTrustedHeight(
	counterparty *TestChain,
	chainName string,
	trustedHeight clienttypes.Height,
) (
	*xibctmtypes.Header, error,
) {
	header := counterparty.LastHeader
	// Relayer must query for LatestHeight on client to get TrustedHeight if the trusted height is not set
	if trustedHeight.IsZero() {
		trustedHeight = chain.GetClientState(chainName).GetLatestHeight().(clienttypes.Height)
	}
	var (
		tmTrustedVals *tmtypes.ValidatorSet
		ok            bool
	)
	// Once we get TrustedHeight from client, we must query the validators from the counterparty chain
	// If the LatestHeight == LastHeader.Height, then TrustedValidators are current validators
	// If LatestHeight < LastHeader.Height, we can query the historical validator set from HistoricalInfo
	if trustedHeight == counterparty.LastHeader.GetHeight() {
		tmTrustedVals = counterparty.Vals
	} else {
		// NOTE: We need to get validators from counterparty at height: trustedHeight+1
		// since the last trusted validators for a header at height h
		// is the NextValidators at h+1 committed to in header h by
		// NextValidatorsHash
		tmTrustedVals, ok = counterparty.GetValsAtHeight(int64(trustedHeight.RevisionHeight + 1))
		if !ok {
			return nil, sdkerrors.Wrapf(
				xibctmtypes.ErrInvalidHeaderHeight,
				"could not retrieve trusted validators at trustedHeight: %d",
				trustedHeight,
			)
		}
	}
	// inject trusted fields into last header
	// for now assume revision number is 0
	header.TrustedHeight = trustedHeight

	trustedVals, err := tmTrustedVals.ToProto()
	if err != nil {
		return nil, err
	}
	header.TrustedValidators = trustedVals

	return header, nil

}

// ExpireClient fast forwards the chain's block time by the provided amount of time which will
// expire any clients with a trusting period less than or equal to this amount of time.
func (chain *TestChain) ExpireClient(amount time.Duration) {
	chain.Coordinator.IncrementTimeBy(amount)
}

// CurrentTMClientHeader creates a TM header using the current header parameters
// on the chain. The trusted fields in the header are set to nil.
func (chain *TestChain) CurrentTMClientHeader() *xibctmtypes.Header {
	return chain.CreateTMClientHeader(
		chain.ChainID, chain.CurrentHeader.Height, clienttypes.Height{},
		chain.CurrentHeader.Time, chain.Vals, nil, chain.Signers,
	)
}

// CreateTMClientHeader creates a TM header to update the TM client. Args are passed in to allow
// caller flexibility to use params that differ from the chain.
func (chain *TestChain) CreateTMClientHeader(
	chainID string,
	blockHeight int64,
	trustedHeight clienttypes.Height,
	timestamp time.Time,
	tmValSet *tmtypes.ValidatorSet,
	tmTrustedVals *tmtypes.ValidatorSet,
	signers []tmtypes.PrivValidator,
) *xibctmtypes.Header {
	var (
		valSet      *tmproto.ValidatorSet
		trustedVals *tmproto.ValidatorSet
	)
	require.NotNil(chain.t, tmValSet)

	vsetHash := tmValSet.Hash()

	tmHeader := tmtypes.Header{
		Version:            tmprotoversion.Consensus{Block: version.BlockProtocol, App: 2},
		ChainID:            chainID,
		Height:             blockHeight,
		Time:               timestamp,
		LastBlockID:        MakeBlockID(make([]byte, tmhash.Size), 10_000, make([]byte, tmhash.Size)),
		LastCommitHash:     chain.App.LastCommitID().Hash,
		DataHash:           tmhash.Sum([]byte("data_hash")),
		ValidatorsHash:     vsetHash,
		NextValidatorsHash: vsetHash,
		ConsensusHash:      tmhash.Sum([]byte("consensus_hash")),
		AppHash:            chain.CurrentHeader.AppHash,
		LastResultsHash:    tmhash.Sum([]byte("last_results_hash")),
		EvidenceHash:       tmhash.Sum([]byte("evidence_hash")),
		ProposerAddress:    tmValSet.Proposer.Address, //nolint:staticcheck
	}
	hhash := tmHeader.Hash()
	blockID := MakeBlockID(hhash, 3, tmhash.Sum([]byte("part_set")))
	voteSet := tmtypes.NewVoteSet(chainID, blockHeight, 1, tmproto.PrecommitType, tmValSet)

	commit, err := tmtypes.MakeCommit(blockID, blockHeight, 1, voteSet, signers, timestamp)
	require.NoError(chain.t, err)

	signedHeader := &tmproto.SignedHeader{
		Header: tmHeader.ToProto(),
		Commit: commit.ToProto(),
	}

	if tmValSet != nil { // nolint
		valSet, err = tmValSet.ToProto()
		if err != nil {
			panic(err)
		}
	}

	if tmTrustedVals != nil {
		trustedVals, err = tmTrustedVals.ToProto()
		if err != nil {
			panic(err)
		}
	}

	// The trusted fields may be nil. They may be filled before relaying messages to a client.
	// The relayer is responsible for querying client and injecting appropriate trusted fields.
	return &xibctmtypes.Header{
		SignedHeader:      signedHeader,
		ValidatorSet:      valSet,
		TrustedHeight:     trustedHeight,
		TrustedValidators: trustedVals,
	}
}

// MakeBlockID copied unimported test functions from tmtypes to use them here
func MakeBlockID(hash []byte, partSetSize uint32, partSetHash []byte) tmtypes.BlockID {
	return tmtypes.BlockID{
		Hash: hash,
		PartSetHeader: tmtypes.PartSetHeader{
			Total: partSetSize,
			Hash:  partSetHash,
		},
	}
}

// CreateSortedSignerArray takes two PrivValidators, and the corresponding Validator structs
// (including voting power). It returns a signer array of PrivValidators that matches the
// sorting of ValidatorSet.
// The sorting is first by .VotingPower (descending), with secondary index of .Address (ascending).
func CreateSortedSignerArray(
	altPrivVal, suitePrivVal tmtypes.PrivValidator, altVal, suiteVal *tmtypes.Validator,
) []tmtypes.PrivValidator {
	switch {
	case altVal.VotingPower > suiteVal.VotingPower:
		return []tmtypes.PrivValidator{altPrivVal, suitePrivVal}
	case altVal.VotingPower < suiteVal.VotingPower:
		return []tmtypes.PrivValidator{suitePrivVal, altPrivVal}
	default:
		if bytes.Compare(altVal.Address, suiteVal.Address) == -1 {
			return []tmtypes.PrivValidator{altPrivVal, suitePrivVal}
		}
		return []tmtypes.PrivValidator{suitePrivVal, altPrivVal}
	}
}

func (chain *TestChain) RegisterRelayer(chains []string, addresses []string) {
	chain.App.XIBCKeeper.ClientKeeper.RegisterRelayers(
		chain.GetContext(),
		chain.SenderAcc.String(),
		chains,
		addresses,
	)
}

func (chain *TestChain) SetPacketChainName() {
	packetContractAbi := packetcontract.PacketContract.ABI
	if _, err := chain.App.XIBCKeeper.PacketKeeper.CallEVM(
		chain.GetContext(),
		packetContractAbi,
		packettypes.ModuleAddress,
		packetcontract.PacketContractAddress,
		"setChainName",
		chain.App.XIBCKeeper.ClientKeeper.GetChainName(chain.GetContext()),
	); err != nil {
		panic(err)
	}
}
