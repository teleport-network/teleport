package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"

	"github.com/teleport-network/teleport/x/xibc/exported"
)

// RegisterInterfaces register the xibc packet submodule interfaces to protobuf Any
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterInterface(
		"xibc.core.packet.v1.PacketI",
		(*exported.PacketI)(nil),
	)
	registry.RegisterImplementations(
		(*exported.PacketI)(nil),
		&Packet{},
	)
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgRecvPacket{},
		&MsgAcknowledgement{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var SubModuleCdc = codec.NewProtoCodec(codectypes.NewInterfaceRegistry())
