package keeper

import (
	"strconv"
	"strings"

	"github.com/tendermint/tendermint/libs/log"
	db "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/teleport-network/teleport/x/xibc/core/host"
	"github.com/teleport-network/teleport/x/xibc/core/packet/types"
	"github.com/teleport-network/teleport/x/xibc/exported"
)

// Keeper defines the XIBC packet keeper
type Keeper struct {
	storeKey      sdk.StoreKey
	cdc           codec.BinaryCodec
	clientKeeper  types.ClientKeeper
	accountKeeper types.AccountKeeper
	evmKeeper     types.EVMKeeper
}

// NewKeeper creates a new xibc packet Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec,
	key sdk.StoreKey,
	clientKeeper types.ClientKeeper,
	accountKeeper types.AccountKeeper,
	evmKeeper types.EVMKeeper,
) Keeper {
	return Keeper{
		storeKey:      key,
		cdc:           cdc,
		clientKeeper:  clientKeeper,
		accountKeeper: accountKeeper,
		evmKeeper:     evmKeeper,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+host.ModuleName+"/"+types.SubModuleName)
}

func (k Keeper) GetModuleAccount(ctx sdk.Context) authtypes.ModuleAccountI {
	return k.accountKeeper.GetModuleAccount(ctx, types.SubModuleName)
}

// GetNextSequenceSend gets next send sequence from the store
func (k Keeper) GetNextSequenceSend(ctx sdk.Context, srcChain string, dstChain string) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(host.NextSequenceSendKey(srcChain, dstChain))
	if bz == nil {
		return 1
	}
	return sdk.BigEndianToUint64(bz)
}

// SetNextSequenceSend sets next send sequence to the store
func (k Keeper) SetNextSequenceSend(ctx sdk.Context, srcChain string, dstChain string, sequence uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := sdk.Uint64ToBigEndian(sequence)
	store.Set(host.NextSequenceSendKey(srcChain, dstChain), bz)
}

// GetPacketReceipt gets a packet receipt from the store
func (k Keeper) GetPacketReceipt(ctx sdk.Context, srcChain string, dstChain string, sequence uint64) (string, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(host.PacketReceiptKey(srcChain, dstChain, sequence))
	if bz == nil {
		return "", false
	}
	return string(bz), true
}

// SetPacketReceipt sets an empty packet receipt to the store
func (k Keeper) SetPacketReceipt(ctx sdk.Context, srcChain string, dstChain string, sequence uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(host.PacketReceiptKey(srcChain, dstChain, sequence), []byte{byte(1)})
}

// HasPacketAcknowledgement check if the packet ack hash is already on the store
func (k Keeper) HasPacketReceipt(ctx sdk.Context, srcChain string, dstChain string, sequence uint64) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(host.PacketReceiptKey(srcChain, dstChain, sequence))
}

// GetPacketCommitment gets the packet commitment hash from the store
func (k Keeper) GetPacketCommitment(ctx sdk.Context, srcChain string, dstChain string, sequence uint64) []byte {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(host.PacketCommitmentKey(srcChain, dstChain, sequence))
	return bz
}

// HasPacketCommitment returns true if the packet commitment exists
func (k Keeper) HasPacketCommitment(ctx sdk.Context, srcChain string, dstChain string, sequence uint64) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(host.PacketCommitmentKey(srcChain, dstChain, sequence))
}

// SetPacketCommitment sets the packet commitment hash to the store
func (k Keeper) SetPacketCommitment(ctx sdk.Context, srcChain string, dstChain string, sequence uint64, commitmentHash []byte) {
	store := ctx.KVStore(k.storeKey)
	store.Set(host.PacketCommitmentKey(srcChain, dstChain, sequence), commitmentHash)
}

func (k Keeper) deletePacketCommitment(ctx sdk.Context, srcChain string, dstChain string, sequence uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(host.PacketCommitmentKey(srcChain, dstChain, sequence))
}

// GetPacketRelayer gets the packet relayer from the store
func (k Keeper) GetPacketRelayer(ctx sdk.Context, srcChain string, dstChain string, sequence uint64) string {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(host.PacketRelayerKey(srcChain, dstChain, sequence))
	return string(bz)
}

// SetPacketRelayer sets the packet relayer  to the store
func (k Keeper) SetPacketRelayer(ctx sdk.Context, srcChain string, dstChain string, sequence uint64, relayer string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(host.PacketRelayerKey(srcChain, dstChain, sequence), []byte(relayer))
}

// SetPacketAcknowledgement sets the packet ack hash to the store
func (k Keeper) SetPacketAcknowledgement(ctx sdk.Context, srcChain string, dstChain string, sequence uint64, ackHash []byte) {
	store := ctx.KVStore(k.storeKey)
	store.Set(host.PacketAcknowledgementKey(srcChain, dstChain, sequence), ackHash)
}

// GetPacketAcknowledgement gets the packet ack hash from the store
func (k Keeper) GetPacketAcknowledgement(ctx sdk.Context, srcChain string, dstChain string, sequence uint64) ([]byte, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(host.PacketAcknowledgementKey(srcChain, dstChain, sequence))
	if bz == nil {
		return nil, false
	}
	return bz, true
}

// HasPacketAcknowledgement check if the packet ack hash is already on the store
func (k Keeper) HasPacketAcknowledgement(ctx sdk.Context, srcChain string, dstChain string, sequence uint64) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(host.PacketAcknowledgementKey(srcChain, dstChain, sequence))
}

// IteratePacketSequence provides an iterator over all send, receive or ack sequences.
// For each sequence, cb will be called. If the cb returns true, the iterator
// will close and stop.
func (k Keeper) IteratePacketSequence(ctx sdk.Context, iterator db.Iterator, cb func(srcChain string, dstChain string, sequence uint64) bool) {
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		srcChain, dstChain, err := host.ParsePath(string(iterator.Key()))
		if err != nil {
			// return if the key is invalid
			return
		}

		sequence := sdk.BigEndianToUint64(iterator.Value())

		if cb(srcChain, dstChain, sequence) {
			break
		}
	}
}

// GetAllPacketSendSeqs returns all stored next send sequences.
func (k Keeper) GetAllPacketSendSeqs(ctx sdk.Context) (seqs []types.PacketSequence) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte(host.KeyNextSeqSendPrefix))
	k.IteratePacketSequence(ctx, iterator, func(srcChain string, dstChain string, nextSendSeq uint64) bool {
		ps := types.NewPacketSequence(srcChain, dstChain, nextSendSeq)
		seqs = append(seqs, ps)
		return false
	})
	return seqs
}

// GetAllPacketRecvSeqs returns all stored next recv sequences.
func (k Keeper) GetAllPacketRecvSeqs(ctx sdk.Context) (seqs []types.PacketSequence) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte(host.KeyNextSeqRecvPrefix))
	k.IteratePacketSequence(ctx, iterator, func(srcChain string, dstChain string, nextRecvSeq uint64) bool {
		ps := types.NewPacketSequence(srcChain, dstChain, nextRecvSeq)
		seqs = append(seqs, ps)
		return false
	})
	return seqs
}

// GetAllPacketAckSeqs returns all stored next acknowledgements sequences.
func (k Keeper) GetAllPacketAckSeqs(ctx sdk.Context) (seqs []types.PacketSequence) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte(host.KeyNextSeqAckPrefix))
	k.IteratePacketSequence(ctx, iterator, func(srcChain string, dstChain string, nextAckSeq uint64) bool {
		ps := types.NewPacketSequence(srcChain, dstChain, nextAckSeq)
		seqs = append(seqs, ps)
		return false
	})
	return seqs
}

// IteratePacketCommitment provides an iterator over all PacketCommitment objects. For each
// packet commitment, cb will be called. If the cb returns true, the iterator will close
// and stop.
func (k Keeper) IteratePacketCommitment(ctx sdk.Context, cb func(srcChain string, dstChain string, sequence uint64, hash []byte) bool) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte(host.KeyPacketCommitmentPrefix))
	k.iterateHashes(ctx, iterator, cb)
}

// GetAllPacketCommitments returns all stored PacketCommitments objects.
func (k Keeper) GetAllPacketCommitments(ctx sdk.Context) (commitments []types.PacketState) {
	k.IteratePacketCommitment(ctx, func(srcChain string, dstChain string, sequence uint64, hash []byte) bool {
		pc := types.NewPacketState(srcChain, dstChain, sequence, hash)
		commitments = append(commitments, pc)
		return false
	})
	return commitments
}

// IteratePacketCommitmentByPath provides an iterator over all PacketCommmitment objects
func (k Keeper) IteratePacketCommitmentByPath(ctx sdk.Context, srcChain string, dstChain string, cb func(_, _ string, sequence uint64, hash []byte) bool) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte(host.PacketCommitmentPrefixPath(srcChain, dstChain)))
	k.iterateHashes(ctx, iterator, cb)
}

// GetAllPacketCommitmentsByPath returns all stored PacketCommitments objects
func (k Keeper) GetAllPacketCommitmentsByPath(ctx sdk.Context, srcChain string, dstChain string) (commitments []types.PacketState) {
	k.IteratePacketCommitmentByPath(ctx, srcChain, dstChain, func(_, _ string, sequence uint64, hash []byte) bool {
		pc := types.NewPacketState(srcChain, dstChain, sequence, hash)
		commitments = append(commitments, pc)
		return false
	})
	return commitments
}

// IteratePacketReceipt provides an iterator over all PacketReceipt objects. For each
// receipt, cb will be called. If the cb returns true, the iterator will close
// and stop.
func (k Keeper) IteratePacketReceipt(ctx sdk.Context, cb func(srcChain string, dstChain string, sequence uint64, receipt []byte) bool) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte(host.KeyPacketReceiptPrefix))
	k.iterateHashes(ctx, iterator, cb)
}

// GetAllPacketReceipts returns all stored PacketReceipt objects.
func (k Keeper) GetAllPacketReceipts(ctx sdk.Context) (receipts []types.PacketState) {
	k.IteratePacketReceipt(ctx, func(srcChain string, dstChain string, sequence uint64, receipt []byte) bool {
		packetReceipt := types.NewPacketState(srcChain, dstChain, sequence, receipt)
		receipts = append(receipts, packetReceipt)
		return false
	})
	return receipts
}

// IteratePacketAcknowledgement provides an iterator over all PacketAcknowledgement objects. For each
// acknowledgement, cb will be called. If the cb returns true, the iterator will close
// and stop.
func (k Keeper) IteratePacketAcknowledgement(ctx sdk.Context, cb func(srcChain string, dstChain string, sequence uint64, hash []byte) bool) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte(host.KeyPacketAckPrefix))
	k.iterateHashes(ctx, iterator, cb)
}

// GetAllPacketAcks returns all stored PacketAcknowledgements objects.
func (k Keeper) GetAllPacketAcks(ctx sdk.Context) (acks []types.PacketState) {
	k.IteratePacketAcknowledgement(ctx, func(srcChain string, dstChain string, sequence uint64, ack []byte) bool {
		packetAck := types.NewPacketState(srcChain, dstChain, sequence, ack)
		acks = append(acks, packetAck)
		return false
	})
	return acks
}

// common functionality for IteratePacketCommitment and IteratePacketAcknowledgement
func (k Keeper) iterateHashes(_ sdk.Context, iterator db.Iterator, cb func(srcChain string, dstChain string, sequence uint64, hash []byte) bool) {
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		keySplit := strings.Split(string(iterator.Key()), "/")
		srcChain := keySplit[1]
		dstChain := keySplit[2]

		sequence, err := strconv.ParseUint(keySplit[len(keySplit)-1], 10, 64)
		if err != nil {
			panic(err)
		}

		if cb(srcChain, dstChain, sequence, iterator.Value()) {
			break
		}
	}
}

// ValidatePacket validates packet sequence
func (k Keeper) ValidatePacket(ctx sdk.Context, packet exported.PacketI) error {
	if err := packet.ValidateBasic(); err != nil {
		return err
	}
	chainName := k.clientKeeper.GetChainName(ctx)
	if packet.GetDstChain() != chainName && packet.GetSrcChain() != chainName {
		return sdkerrors.Wrap(types.ErrInvalidPacket, "packet/ack illegal!")
	}
	return nil
}
