package xibctesting

import (
	"bytes"
	"fmt"

	packettypes "github.com/teleport-network/teleport/x/xibc/core/packet/types"
)

// Path contains two endpoints representing two chains connected over XIBC
type Path struct {
	EndpointA *Endpoint
	EndpointB *Endpoint
}

// NewPath constructs an endpoint for each chain using the default values
// for the endpoints. Each endpoint is updated to have a pointer to the
// counterparty endpoint.
func NewPath(chainA, chainB *TestChain) *Path {
	endpointA := NewDefaultEndpoint(chainA)
	endpointB := NewDefaultEndpoint(chainB)

	endpointA.Counterparty = endpointB
	endpointB.Counterparty = endpointA

	return &Path{
		EndpointA: endpointA,
		EndpointB: endpointB,
	}
}

// RelayPacket attempts to relay the packet first on EndpointA and then on EndpointB
// if EndpointA does not contain a packet commitment for that packet. An error is returned
// if a relay step fails or the packet commitment does not exist on either endpoint.
func (path *Path) RelayPacket(packet packettypes.Packet, ack []byte) error {
	if bytes.Equal(
		packettypes.CommitPacket(packet),
		path.EndpointA.Chain.App.XIBCKeeper.PacketKeeper.GetPacketCommitment(
			path.EndpointA.Chain.GetContext(),
			packet.GetSourceChain(),
			packet.GetDestChain(),
			packet.GetSequence(),
		),
	) {
		// packet found, relay from A to B
		if err := path.EndpointB.UpdateClient(); err != nil {
			return err
		}
		if err := path.EndpointB.RecvPacket(packet); err != nil {
			return err
		}
		if path.EndpointB.ChainName != packet.DestinationChain {
			return nil
		}
		return path.EndpointA.AcknowledgePacket(packet, ack)
	}

	if bytes.Equal(
		packettypes.CommitPacket(packet),
		path.EndpointB.Chain.App.XIBCKeeper.PacketKeeper.GetPacketCommitment(
			path.EndpointB.Chain.GetContext(),
			packet.GetSourceChain(),
			packet.GetDestChain(),
			packet.GetSequence(),
		),
	) {
		// packet found, relay B to A
		if err := path.EndpointA.UpdateClient(); err != nil {
			return err
		}
		if err := path.EndpointA.RecvPacket(packet); err != nil {
			return err
		}
		if path.EndpointA.ChainName != packet.DestinationChain {
			return nil
		}
		return path.EndpointB.AcknowledgePacket(packet, ack)
	}

	return fmt.Errorf("packet commitment does not exist on either endpoint for provided packet")
}
