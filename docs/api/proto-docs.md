<!-- This file is auto-generated. Please do not modify it yourself. -->
# Protobuf Documentation
<a name="top"></a>

## Table of Contents

- [teleport/aggregate/v1/aggregate.proto](#teleport/aggregate/v1/aggregate.proto)
    - [RegisterCoinProposal](#teleport.aggregate.v1.RegisterCoinProposal)
    - [RegisterERC20Proposal](#teleport.aggregate.v1.RegisterERC20Proposal)
    - [RegisterERC20TraceProposal](#teleport.aggregate.v1.RegisterERC20TraceProposal)
    - [ToggleTokenRelayProposal](#teleport.aggregate.v1.ToggleTokenRelayProposal)
    - [TokenPair](#teleport.aggregate.v1.TokenPair)
    - [UpdateTokenPairERC20Proposal](#teleport.aggregate.v1.UpdateTokenPairERC20Proposal)
  
    - [Owner](#teleport.aggregate.v1.Owner)
  
- [teleport/aggregate/v1/genesis.proto](#teleport/aggregate/v1/genesis.proto)
    - [GenesisState](#teleport.aggregate.v1.GenesisState)
    - [Params](#teleport.aggregate.v1.Params)
  
- [teleport/aggregate/v1/query.proto](#teleport/aggregate/v1/query.proto)
    - [QueryParamsRequest](#teleport.aggregate.v1.QueryParamsRequest)
    - [QueryParamsResponse](#teleport.aggregate.v1.QueryParamsResponse)
    - [QueryTokenPairRequest](#teleport.aggregate.v1.QueryTokenPairRequest)
    - [QueryTokenPairResponse](#teleport.aggregate.v1.QueryTokenPairResponse)
    - [QueryTokenPairsRequest](#teleport.aggregate.v1.QueryTokenPairsRequest)
    - [QueryTokenPairsResponse](#teleport.aggregate.v1.QueryTokenPairsResponse)
  
    - [Query](#teleport.aggregate.v1.Query)
  
- [teleport/aggregate/v1/tx.proto](#teleport/aggregate/v1/tx.proto)
    - [MsgConvertCoin](#teleport.aggregate.v1.MsgConvertCoin)
    - [MsgConvertCoinResponse](#teleport.aggregate.v1.MsgConvertCoinResponse)
    - [MsgConvertERC20](#teleport.aggregate.v1.MsgConvertERC20)
    - [MsgConvertERC20Response](#teleport.aggregate.v1.MsgConvertERC20Response)
  
    - [Msg](#teleport.aggregate.v1.Msg)
  
- [xibc/apps/multicall/v1/multicall.proto](#xibc/apps/multicall/v1/multicall.proto)
- [xibc/apps/rcc/v1/event.proto](#xibc/apps/rcc/v1/event.proto)
    - [EventAckPacket](#xibc.apps.rcc.v1.EventAckPacket)
    - [EventRecvPacket](#xibc.apps.rcc.v1.EventRecvPacket)
  
- [xibc/apps/rcc/v1/rcc.proto](#xibc/apps/rcc/v1/rcc.proto)
    - [RCCPacketData](#xibc.apps.rcc.v1.RCCPacketData)
  
- [xibc/apps/transfer/v1/event.proto](#xibc/apps/transfer/v1/event.proto)
    - [EventAckPacket](#xibc.apps.transfer.v1.EventAckPacket)
    - [EventRecvPacket](#xibc.apps.transfer.v1.EventRecvPacket)
  
- [xibc/apps/transfer/v1/transfer.proto](#xibc/apps/transfer/v1/transfer.proto)
    - [FungibleTokenPacketData](#xibc.apps.transfer.v1.FungibleTokenPacketData)
  
- [xibc/clients/tssclient/v1/tssclient.proto](#xibc/clients/tssclient/v1/tssclient.proto)
    - [ClientState](#xibc.clients.tssclient.v1.ClientState)
    - [ConsensusState](#xibc.clients.tssclient.v1.ConsensusState)
    - [Header](#xibc.clients.tssclient.v1.Header)
  
- [xibc/core/client/v1/client.proto](#xibc/core/client/v1/client.proto)
    - [ClientConsensusStates](#xibc.core.client.v1.ClientConsensusStates)
    - [ConsensusStateWithHeight](#xibc.core.client.v1.ConsensusStateWithHeight)
    - [CreateClientProposal](#xibc.core.client.v1.CreateClientProposal)
    - [Height](#xibc.core.client.v1.Height)
    - [IdentifiedClientState](#xibc.core.client.v1.IdentifiedClientState)
    - [IdentifiedRelayers](#xibc.core.client.v1.IdentifiedRelayers)
    - [RegisterRelayerProposal](#xibc.core.client.v1.RegisterRelayerProposal)
    - [ToggleClientProposal](#xibc.core.client.v1.ToggleClientProposal)
    - [UpgradeClientProposal](#xibc.core.client.v1.UpgradeClientProposal)
  
- [xibc/core/client/v1/event.proto](#xibc/core/client/v1/event.proto)
    - [EventCreateClientProposal](#xibc.core.client.v1.EventCreateClientProposal)
    - [EventRegisterRelayerProposal](#xibc.core.client.v1.EventRegisterRelayerProposal)
    - [EventToggleClientProposal](#xibc.core.client.v1.EventToggleClientProposal)
    - [EventUpdateClient](#xibc.core.client.v1.EventUpdateClient)
    - [EventUpgradeClientProposal](#xibc.core.client.v1.EventUpgradeClientProposal)
  
- [xibc/core/client/v1/genesis.proto](#xibc/core/client/v1/genesis.proto)
    - [GenesisMetadata](#xibc.core.client.v1.GenesisMetadata)
    - [GenesisState](#xibc.core.client.v1.GenesisState)
    - [IdentifiedGenesisMetadata](#xibc.core.client.v1.IdentifiedGenesisMetadata)
  
- [xibc/core/client/v1/query.proto](#xibc/core/client/v1/query.proto)
    - [QueryClientStateRequest](#xibc.core.client.v1.QueryClientStateRequest)
    - [QueryClientStateResponse](#xibc.core.client.v1.QueryClientStateResponse)
    - [QueryClientStatesRequest](#xibc.core.client.v1.QueryClientStatesRequest)
    - [QueryClientStatesResponse](#xibc.core.client.v1.QueryClientStatesResponse)
    - [QueryConsensusStateRequest](#xibc.core.client.v1.QueryConsensusStateRequest)
    - [QueryConsensusStateResponse](#xibc.core.client.v1.QueryConsensusStateResponse)
    - [QueryConsensusStatesRequest](#xibc.core.client.v1.QueryConsensusStatesRequest)
    - [QueryConsensusStatesResponse](#xibc.core.client.v1.QueryConsensusStatesResponse)
    - [QueryRelayersRequest](#xibc.core.client.v1.QueryRelayersRequest)
    - [QueryRelayersResponse](#xibc.core.client.v1.QueryRelayersResponse)
  
    - [Query](#xibc.core.client.v1.Query)
  
- [xibc/core/client/v1/tx.proto](#xibc/core/client/v1/tx.proto)
    - [MsgUpdateClient](#xibc.core.client.v1.MsgUpdateClient)
    - [MsgUpdateClientResponse](#xibc.core.client.v1.MsgUpdateClientResponse)
  
    - [Msg](#xibc.core.client.v1.Msg)
  
- [xibc/core/commitment/v1/commitment.proto](#xibc/core/commitment/v1/commitment.proto)
    - [MerklePath](#xibc.core.commitment.v1.MerklePath)
    - [MerklePrefix](#xibc.core.commitment.v1.MerklePrefix)
    - [MerkleProof](#xibc.core.commitment.v1.MerkleProof)
  
- [xibc/core/packet/v1/event.proto](#xibc/core/packet/v1/event.proto)
    - [EventAcknowledgePacket](#xibc.core.packet.v1.EventAcknowledgePacket)
    - [EventRecvPacket](#xibc.core.packet.v1.EventRecvPacket)
    - [EventSendPacket](#xibc.core.packet.v1.EventSendPacket)
    - [EventWriteAck](#xibc.core.packet.v1.EventWriteAck)
  
- [xibc/core/packet/v1/packet.proto](#xibc/core/packet/v1/packet.proto)
    - [Acknowledgement](#xibc.core.packet.v1.Acknowledgement)
    - [Packet](#xibc.core.packet.v1.Packet)
    - [PacketState](#xibc.core.packet.v1.PacketState)
  
- [xibc/core/packet/v1/genesis.proto](#xibc/core/packet/v1/genesis.proto)
    - [GenesisState](#xibc.core.packet.v1.GenesisState)
    - [PacketSequence](#xibc.core.packet.v1.PacketSequence)
  
- [xibc/core/packet/v1/query.proto](#xibc/core/packet/v1/query.proto)
    - [QueryPacketAcknowledgementRequest](#xibc.core.packet.v1.QueryPacketAcknowledgementRequest)
    - [QueryPacketAcknowledgementResponse](#xibc.core.packet.v1.QueryPacketAcknowledgementResponse)
    - [QueryPacketAcknowledgementsRequest](#xibc.core.packet.v1.QueryPacketAcknowledgementsRequest)
    - [QueryPacketAcknowledgementsResponse](#xibc.core.packet.v1.QueryPacketAcknowledgementsResponse)
    - [QueryPacketCommitmentRequest](#xibc.core.packet.v1.QueryPacketCommitmentRequest)
    - [QueryPacketCommitmentResponse](#xibc.core.packet.v1.QueryPacketCommitmentResponse)
    - [QueryPacketCommitmentsRequest](#xibc.core.packet.v1.QueryPacketCommitmentsRequest)
    - [QueryPacketCommitmentsResponse](#xibc.core.packet.v1.QueryPacketCommitmentsResponse)
    - [QueryPacketReceiptRequest](#xibc.core.packet.v1.QueryPacketReceiptRequest)
    - [QueryPacketReceiptResponse](#xibc.core.packet.v1.QueryPacketReceiptResponse)
    - [QueryUnreceivedAcksRequest](#xibc.core.packet.v1.QueryUnreceivedAcksRequest)
    - [QueryUnreceivedAcksResponse](#xibc.core.packet.v1.QueryUnreceivedAcksResponse)
    - [QueryUnreceivedPacketsRequest](#xibc.core.packet.v1.QueryUnreceivedPacketsRequest)
    - [QueryUnreceivedPacketsResponse](#xibc.core.packet.v1.QueryUnreceivedPacketsResponse)
  
    - [Query](#xibc.core.packet.v1.Query)
  
- [xibc/core/packet/v1/tx.proto](#xibc/core/packet/v1/tx.proto)
    - [MsgAcknowledgement](#xibc.core.packet.v1.MsgAcknowledgement)
    - [MsgAcknowledgementResponse](#xibc.core.packet.v1.MsgAcknowledgementResponse)
    - [MsgRecvPacket](#xibc.core.packet.v1.MsgRecvPacket)
    - [MsgRecvPacketResponse](#xibc.core.packet.v1.MsgRecvPacketResponse)
  
    - [Msg](#xibc.core.packet.v1.Msg)
  
- [xibc/core/types/v1/genesis.proto](#xibc/core/types/v1/genesis.proto)
    - [GenesisState](#xibc.core.types.v1.GenesisState)
  
- [Scalar Value Types](#scalar-value-types)



<a name="teleport/aggregate/v1/aggregate.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## teleport/aggregate/v1/aggregate.proto



<a name="teleport.aggregate.v1.RegisterCoinProposal"></a>

### RegisterCoinProposal
RegisterCoinProposal is a gov Content type to register a token pair


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  | title of the proposal |
| `description` | [string](#string) |  | proposal description |
| `metadata` | [cosmos.bank.v1beta1.Metadata](#cosmos.bank.v1beta1.Metadata) |  | token pair of Cosmos native denom and ERC20 token address |






<a name="teleport.aggregate.v1.RegisterERC20Proposal"></a>

### RegisterERC20Proposal
RegisterCoinProposal is a gov Content type to register a token pair


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  | title of the proposal |
| `description` | [string](#string) |  | proposal description |
| `erc20_address` | [string](#string) |  | contract address of ERC20 token |






<a name="teleport.aggregate.v1.RegisterERC20TraceProposal"></a>

### RegisterERC20TraceProposal
RegisterERC20TraceProposal is a gov Content type to register a ERC20 trace


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  | title of the proposal |
| `description` | [string](#string) |  | proposal description |
| `erc20_address` | [string](#string) |  | contract address of ERC20 token |
| `origin_token` | [string](#string) |  | origin token |
| `origin_chain` | [string](#string) |  | origin chain |
| `scale` | [uint64](#uint64) |  | scale: real_amount = packet_amount * (10 ** scale) |






<a name="teleport.aggregate.v1.ToggleTokenRelayProposal"></a>

### ToggleTokenRelayProposal
ToggleTokenRelayProposal is a gov Content type to toggle the internal
relaying of a token pair


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  | title of the proposal |
| `description` | [string](#string) |  | proposal description |
| `token` | [string](#string) |  | token identifier can be either the hex contract address of the ERC20 or the Cosmos base denomination |






<a name="teleport.aggregate.v1.TokenPair"></a>

### TokenPair
TokenPair defines an instance that records pairing consisting of a Cosmos
native Coin and an ERC20 token address


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `erc20_address` | [string](#string) |  | address of ERC20 contract token |
| `denom` | [string](#string) |  | cosmos base denomination to be mapped to |
| `enabled` | [bool](#bool) |  | shows token mapping enable status |
| `contract_owner` | [Owner](#teleport.aggregate.v1.Owner) |  | ERC20 owner address ENUM (0 invalid, 1 ModuleAccount, 2 external address) |






<a name="teleport.aggregate.v1.UpdateTokenPairERC20Proposal"></a>

### UpdateTokenPairERC20Proposal
UpdateTokenPairERC20Proposal is a gov Content type to update a token pair's
ERC20 contract address


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  | title of the proposal |
| `description` | [string](#string) |  | proposal description |
| `erc20_address` | [string](#string) |  | contract address of ERC20 token |
| `new_erc20_address` | [string](#string) |  | new address of ERC20 token contract |





 <!-- end messages -->


<a name="teleport.aggregate.v1.Owner"></a>

### Owner
Owner enumerates the ownership of a ERC20 contract

| Name | Number | Description |
| ---- | ------ | ----------- |
| OWNER_UNSPECIFIED | 0 | OWNER_UNSPECIFIED defines an invalid/undefined owner |
| OWNER_MODULE | 1 | OWNER_MODULE erc20 is owned by the intrarelayer module account |
| OWNER_EXTERNAL | 2 | EXTERNAL erc20 is owned by an external account |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="teleport/aggregate/v1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## teleport/aggregate/v1/genesis.proto



<a name="teleport.aggregate.v1.GenesisState"></a>

### GenesisState
GenesisState defines the module's genesis state


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#teleport.aggregate.v1.Params) |  | module parameters |
| `token_pairs` | [TokenPair](#teleport.aggregate.v1.TokenPair) | repeated | registered token pairs |






<a name="teleport.aggregate.v1.Params"></a>

### Params
Params defines the aggregate module params


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `enable_aggregate` | [bool](#bool) |  | parameter to enable the intrarelaying of Cosmos coins <--> ERC20 tokens |
| `enable_evm_hook` | [bool](#bool) |  | parameter to enable the EVM hook to convert an ERC20 token to a Cosmos Coin by transferring the Tokens through a MsgEthereumTx to the ModuleAddress Ethereum address. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="teleport/aggregate/v1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## teleport/aggregate/v1/query.proto



<a name="teleport.aggregate.v1.QueryParamsRequest"></a>

### QueryParamsRequest
QueryParamsRequest is the request type for the Query/Params RPC method






<a name="teleport.aggregate.v1.QueryParamsResponse"></a>

### QueryParamsResponse
QueryParamsResponse is the response type for the Query/Params RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#teleport.aggregate.v1.Params) |  |  |






<a name="teleport.aggregate.v1.QueryTokenPairRequest"></a>

### QueryTokenPairRequest
QueryTokenPairRequest is the request type for the Query/TokenPair RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `token` | [string](#string) |  | token identifier can be either the hex contract address of the ERC20 or the Cosmos base denomination |






<a name="teleport.aggregate.v1.QueryTokenPairResponse"></a>

### QueryTokenPairResponse
QueryTokenPairResponse is the response type for the Query/TokenPair RPC
method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `token_pair` | [TokenPair](#teleport.aggregate.v1.TokenPair) |  |  |






<a name="teleport.aggregate.v1.QueryTokenPairsRequest"></a>

### QueryTokenPairsRequest
QueryTokenPairsRequest is the request type for the Query/TokenPairs RPC
method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination for the request |






<a name="teleport.aggregate.v1.QueryTokenPairsResponse"></a>

### QueryTokenPairsResponse
QueryTokenPairsResponse is the response type for the Query/TokenPairs RPC
method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `token_pairs` | [TokenPair](#teleport.aggregate.v1.TokenPair) | repeated |  |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines the pagination in the response |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="teleport.aggregate.v1.Query"></a>

### Query
Query defines the gRPC querier service

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `TokenPairs` | [QueryTokenPairsRequest](#teleport.aggregate.v1.QueryTokenPairsRequest) | [QueryTokenPairsResponse](#teleport.aggregate.v1.QueryTokenPairsResponse) | Retrieves registered token pairs | GET|/teleport/aggregate/v1/token_pairs|
| `TokenPair` | [QueryTokenPairRequest](#teleport.aggregate.v1.QueryTokenPairRequest) | [QueryTokenPairResponse](#teleport.aggregate.v1.QueryTokenPairResponse) | Retrieves a registered token pair | GET|/teleport/aggregate/v1/token_pairs/{token}|
| `Params` | [QueryParamsRequest](#teleport.aggregate.v1.QueryParamsRequest) | [QueryParamsResponse](#teleport.aggregate.v1.QueryParamsResponse) | Params retrieves the aggregate module params | GET|/teleport/aggregate/v1/params|

 <!-- end services -->



<a name="teleport/aggregate/v1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## teleport/aggregate/v1/tx.proto



<a name="teleport.aggregate.v1.MsgConvertCoin"></a>

### MsgConvertCoin
MsgConvertCoin defines a Msg to convert a Cosmos Coin to a ERC20 token


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `coin` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  | Cosmos coin which denomination is registered on aggregate bridge. The coin amount defines the total ERC20 tokens to convert. |
| `receiver` | [string](#string) |  | recipient hex address to receive ERC20 token |
| `sender` | [string](#string) |  | cosmos bech32 address from the owner of the given ERC20 tokens |






<a name="teleport.aggregate.v1.MsgConvertCoinResponse"></a>

### MsgConvertCoinResponse
MsgConvertCoinResponse returns no fields






<a name="teleport.aggregate.v1.MsgConvertERC20"></a>

### MsgConvertERC20
MsgConvertERC20 defines a Msg to convert an ERC20 token to a Cosmos SDK coin.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_address` | [string](#string) |  | ERC20 token contract address registered on aggregate bridge |
| `amount` | [string](#string) |  | amount of ERC20 tokens to mint |
| `receiver` | [string](#string) |  | bech32 address to receive SDK coins. |
| `sender` | [string](#string) |  | sender hex address from the owner of the given ERC20 tokens |






<a name="teleport.aggregate.v1.MsgConvertERC20Response"></a>

### MsgConvertERC20Response
MsgConvertERC20Response returns no fields





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="teleport.aggregate.v1.Msg"></a>

### Msg
Msg defines the aggregate Msg service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `ConvertCoin` | [MsgConvertCoin](#teleport.aggregate.v1.MsgConvertCoin) | [MsgConvertCoinResponse](#teleport.aggregate.v1.MsgConvertCoinResponse) | ConvertCoin mints a ERC20 representation of the SDK Coin denom that is registered on the token mapping. | GET|/teleport/aggregate/v1/tx/convert_coin|
| `ConvertERC20` | [MsgConvertERC20](#teleport.aggregate.v1.MsgConvertERC20) | [MsgConvertERC20Response](#teleport.aggregate.v1.MsgConvertERC20Response) | ConvertERC20 mints a Cosmos coin representation of the ERC20 token contract that is registered on the token mapping. | GET|/teleport/aggregate/v1/tx/convert_erc20|

 <!-- end services -->



<a name="xibc/apps/multicall/v1/multicall.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## xibc/apps/multicall/v1/multicall.proto


 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="xibc/apps/rcc/v1/event.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## xibc/apps/rcc/v1/event.proto



<a name="xibc.apps.rcc.v1.EventAckPacket"></a>

### EventAckPacket
EventAckPacket is emitted on acknowledgement packet


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `src_chain` | [string](#string) |  |  |
| `dest_chain` | [string](#string) |  |  |
| `sender` | [string](#string) |  |  |
| `contract_address` | [string](#string) |  |  |
| `data` | [bytes](#bytes) |  |  |
| `result` | [bytes](#bytes) |  |  |






<a name="xibc.apps.rcc.v1.EventRecvPacket"></a>

### EventRecvPacket
EventRecvPacket is emitted on receive packet


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `src_chain` | [string](#string) |  |  |
| `dest_chain` | [string](#string) |  |  |
| `sender` | [string](#string) |  |  |
| `contract_address` | [string](#string) |  |  |
| `data` | [bytes](#bytes) |  |  |
| `result` | [bytes](#bytes) |  |  |
| `message` | [string](#string) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="xibc/apps/rcc/v1/rcc.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## xibc/apps/rcc/v1/rcc.proto



<a name="xibc.apps.rcc.v1.RCCPacketData"></a>

### RCCPacketData
RCCPacketData defines a struct for the packet payload


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `src_chain` | [string](#string) |  |  |
| `dest_chain` | [string](#string) |  |  |
| `sender` | [string](#string) |  |  |
| `contract_address` | [string](#string) |  |  |
| `data` | [bytes](#bytes) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="xibc/apps/transfer/v1/event.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## xibc/apps/transfer/v1/event.proto



<a name="xibc.apps.transfer.v1.EventAckPacket"></a>

### EventAckPacket
EventAckPacket is emitted on acknowledgement packet


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `src_chain` | [string](#string) |  |  |
| `dest_chain` | [string](#string) |  |  |
| `sender` | [string](#string) |  |  |
| `receiver` | [string](#string) |  |  |
| `token` | [string](#string) |  |  |
| `amount` | [bytes](#bytes) |  |  |
| `result` | [bytes](#bytes) |  |  |






<a name="xibc.apps.transfer.v1.EventRecvPacket"></a>

### EventRecvPacket
EventRecvPacket is emitted on receive packet


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `src_chain` | [string](#string) |  |  |
| `dest_chain` | [string](#string) |  |  |
| `sender` | [string](#string) |  |  |
| `receiver` | [string](#string) |  |  |
| `token` | [string](#string) |  |  |
| `amount` | [bytes](#bytes) |  |  |
| `result` | [bytes](#bytes) |  |  |
| `message` | [string](#string) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="xibc/apps/transfer/v1/transfer.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## xibc/apps/transfer/v1/transfer.proto



<a name="xibc.apps.transfer.v1.FungibleTokenPacketData"></a>

### FungibleTokenPacketData
FungibleTokenPacketData defines a struct for the packet payload


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `src_chain` | [string](#string) |  |  |
| `dest_chain` | [string](#string) |  |  |
| `sender` | [string](#string) |  |  |
| `receiver` | [string](#string) |  |  |
| `amount` | [bytes](#bytes) |  |  |
| `token` | [string](#string) |  |  |
| `ori_token` | [string](#string) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="xibc/clients/tssclient/v1/tssclient.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## xibc/clients/tssclient/v1/tssclient.proto



<a name="xibc.clients.tssclient.v1.ClientState"></a>

### ClientState
ClientState from Tss tracks the current tss address, and a possible frozen
height.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `tss_address` | [string](#string) |  |  |
| `pubkey` | [bytes](#bytes) |  |  |
| `part_pubkeys` | [bytes](#bytes) | repeated |  |






<a name="xibc.clients.tssclient.v1.ConsensusState"></a>

### ConsensusState
ConsensusState defines the consensus state






<a name="xibc.clients.tssclient.v1.Header"></a>

### Header



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `tss_address` | [string](#string) |  |  |
| `pubkey` | [bytes](#bytes) |  |  |
| `part_pubkeys` | [bytes](#bytes) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="xibc/core/client/v1/client.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## xibc/core/client/v1/client.proto



<a name="xibc.core.client.v1.ClientConsensusStates"></a>

### ClientConsensusStates
ClientConsensusStates defines all the stored consensus states for a given
client.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `chain_name` | [string](#string) |  | client identifier |
| `consensus_states` | [ConsensusStateWithHeight](#xibc.core.client.v1.ConsensusStateWithHeight) | repeated | consensus states and their heights associated with the client |






<a name="xibc.core.client.v1.ConsensusStateWithHeight"></a>

### ConsensusStateWithHeight
ConsensusStateWithHeight defines a consensus state with an additional height
field.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `height` | [Height](#xibc.core.client.v1.Height) |  | consensus state height |
| `consensus_state` | [google.protobuf.Any](#google.protobuf.Any) |  | consensus state |






<a name="xibc.core.client.v1.CreateClientProposal"></a>

### CreateClientProposal
CreateClientProposal defines a overnance proposal to create an XIBC client


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  | the title of the update proposal |
| `description` | [string](#string) |  | the description of the proposal |
| `chain_name` | [string](#string) |  | the client identifier for the client to be updated if the proposal passes |
| `client_state` | [google.protobuf.Any](#google.protobuf.Any) |  | client state |
| `consensus_state` | [google.protobuf.Any](#google.protobuf.Any) |  | consensus state associated with the client that corresponds to a given height. |






<a name="xibc.core.client.v1.Height"></a>

### Height
Height is a monotonically increasing data type
that can be compared against another Height for the purposes of updating and
freezing clients

Normally the RevisionHeight is incremented at each height while keeping
RevisionNumber the same. However some consensus algorithms may choose to
reset the height in certain conditions e.g. hard forks, state-machine
breaking changes In these cases, the RevisionNumber is incremented so that
height continues to be monitonically increasing even as the RevisionHeight
gets reset


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `revision_number` | [uint64](#uint64) |  | the revision that the client is currently on |
| `revision_height` | [uint64](#uint64) |  | the height within the given revision |






<a name="xibc.core.client.v1.IdentifiedClientState"></a>

### IdentifiedClientState
IdentifiedClientState defines a client state with an additional client
identifier field.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `chain_name` | [string](#string) |  | client identifier |
| `client_state` | [google.protobuf.Any](#google.protobuf.Any) |  | client state |






<a name="xibc.core.client.v1.IdentifiedRelayers"></a>

### IdentifiedRelayers
IdentifiedRelayer defines a list of authorized relayers for the specified
client.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `chain_name` | [string](#string) |  | client identifier |
| `relayers` | [string](#string) | repeated | authorized relayer list |






<a name="xibc.core.client.v1.RegisterRelayerProposal"></a>

### RegisterRelayerProposal
RegisterRelayerProposal defines a overnance proposal to register some
relayers for updating a client state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  | the title of the update proposal |
| `description` | [string](#string) |  | the description of the proposal |
| `chain_name` | [string](#string) |  | the client identifier for the client to be updated if the proposal passes |
| `relayers` | [string](#string) | repeated | relayer address list |






<a name="xibc.core.client.v1.ToggleClientProposal"></a>

### ToggleClientProposal
ToggleClientProposal defines a overnance proposal to toggle XIBC client type


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  | the title of the toggle client proposal |
| `description` | [string](#string) |  | the description of the proposal |
| `chain_name` | [string](#string) |  | the client identifier for the client to be updated if the proposal passes |
| `client_state` | [google.protobuf.Any](#google.protobuf.Any) |  | client state |
| `consensus_state` | [google.protobuf.Any](#google.protobuf.Any) |  | consensus state associated with the client that corresponds to a given height. |






<a name="xibc.core.client.v1.UpgradeClientProposal"></a>

### UpgradeClientProposal
UpgradeClientProposal defines a overnance proposal to overide an XIBC client
state


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  | the title of the update proposal |
| `description` | [string](#string) |  | the description of the proposal |
| `chain_name` | [string](#string) |  | the client identifier for the client to be updated if the proposal passes |
| `client_state` | [google.protobuf.Any](#google.protobuf.Any) |  | client state |
| `consensus_state` | [google.protobuf.Any](#google.protobuf.Any) |  | consensus state |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="xibc/core/client/v1/event.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## xibc/core/client/v1/event.proto



<a name="xibc.core.client.v1.EventCreateClientProposal"></a>

### EventCreateClientProposal
EventCreateClientProposal is emitted on create client proposal


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `chain_name` | [string](#string) |  |  |
| `client_type` | [string](#string) |  |  |
| `consensus_height` | [string](#string) |  |  |






<a name="xibc.core.client.v1.EventRegisterRelayerProposal"></a>

### EventRegisterRelayerProposal



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `chain_name` | [string](#string) |  |  |
| `relayers` | [string](#string) | repeated |  |






<a name="xibc.core.client.v1.EventToggleClientProposal"></a>

### EventToggleClientProposal
EventToggleClientProposal is emitted on toggle client proposal


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `chain_name` | [string](#string) |  |  |
| `client_type` | [string](#string) |  |  |
| `consensus_height` | [string](#string) |  |  |






<a name="xibc.core.client.v1.EventUpdateClient"></a>

### EventUpdateClient
EventUpdateClient is emitted on update client


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `chain_name` | [string](#string) |  |  |
| `client_type` | [string](#string) |  |  |
| `consensus_height` | [string](#string) |  |  |
| `header` | [string](#string) |  |  |






<a name="xibc.core.client.v1.EventUpgradeClientProposal"></a>

### EventUpgradeClientProposal
EventUpgradeClientProposal is emitted on upgrade client proposal


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `chain_name` | [string](#string) |  |  |
| `client_type` | [string](#string) |  |  |
| `consensus_height` | [string](#string) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="xibc/core/client/v1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## xibc/core/client/v1/genesis.proto



<a name="xibc.core.client.v1.GenesisMetadata"></a>

### GenesisMetadata
GenesisMetadata defines the genesis type for metadata that clients may return
with ExportMetadata


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key` | [bytes](#bytes) |  | store key of metadata without chainName-prefix |
| `value` | [bytes](#bytes) |  | metadata value |






<a name="xibc.core.client.v1.GenesisState"></a>

### GenesisState
GenesisState defines the xibc client submodule's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `clients` | [IdentifiedClientState](#xibc.core.client.v1.IdentifiedClientState) | repeated | client states with their corresponding identifiers |
| `clients_consensus` | [ClientConsensusStates](#xibc.core.client.v1.ClientConsensusStates) | repeated | consensus states from each client |
| `clients_metadata` | [IdentifiedGenesisMetadata](#xibc.core.client.v1.IdentifiedGenesisMetadata) | repeated | metadata from each client |
| `native_chain_name` | [string](#string) |  | the chain name of the current chain |
| `relayers` | [IdentifiedRelayers](#xibc.core.client.v1.IdentifiedRelayers) | repeated | IdentifiedRelayer defines a list of authorized relayers for the specified client. |






<a name="xibc.core.client.v1.IdentifiedGenesisMetadata"></a>

### IdentifiedGenesisMetadata
IdentifiedGenesisMetadata has the client metadata with the corresponding
chain name.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `chain_name` | [string](#string) |  |  |
| `metadata` | [GenesisMetadata](#xibc.core.client.v1.GenesisMetadata) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="xibc/core/client/v1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## xibc/core/client/v1/query.proto



<a name="xibc.core.client.v1.QueryClientStateRequest"></a>

### QueryClientStateRequest
QueryClientStateRequest is the request type for the Query/ClientState RPC
method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `chain_name` | [string](#string) |  | client state unique identifier |






<a name="xibc.core.client.v1.QueryClientStateResponse"></a>

### QueryClientStateResponse
QueryClientStateResponse is the response type for the Query/ClientState RPC
method. Besides the client state, it includes a proof and the height from
which the proof was retrieved.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `client_state` | [google.protobuf.Any](#google.protobuf.Any) |  | client state associated with the request identifier |
| `proof` | [bytes](#bytes) |  | merkle proof of existence |
| `proof_height` | [Height](#xibc.core.client.v1.Height) |  | height at which the proof was retrieved |






<a name="xibc.core.client.v1.QueryClientStatesRequest"></a>

### QueryClientStatesRequest
QueryClientStatesRequest is the request type for the Query/ClientStates RPC
method






<a name="xibc.core.client.v1.QueryClientStatesResponse"></a>

### QueryClientStatesResponse
QueryClientStatesResponse is the response type for the Query/ClientStates RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `client_states` | [IdentifiedClientState](#xibc.core.client.v1.IdentifiedClientState) | repeated | list of stored ClientStates of the chain. |






<a name="xibc.core.client.v1.QueryConsensusStateRequest"></a>

### QueryConsensusStateRequest
QueryConsensusStateRequest is the request type for the Query/ConsensusState
RPC method. Besides the consensus state, it includes a proof and the height
from which the proof was retrieved.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `chain_name` | [string](#string) |  | client identifier |
| `revision_number` | [uint64](#uint64) |  | consensus state revision number |
| `revision_height` | [uint64](#uint64) |  | consensus state revision height |
| `latest_height` | [bool](#bool) |  | latest_height overrrides the height field and queries the latest stored ConsensusState |






<a name="xibc.core.client.v1.QueryConsensusStateResponse"></a>

### QueryConsensusStateResponse
QueryConsensusStateResponse is the response type for the Query/ConsensusState
RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `consensus_state` | [google.protobuf.Any](#google.protobuf.Any) |  | consensus state associated with the client identifier at the given height |
| `proof` | [bytes](#bytes) |  | merkle proof of existence |
| `proof_height` | [Height](#xibc.core.client.v1.Height) |  | height at which the proof was retrieved |






<a name="xibc.core.client.v1.QueryConsensusStatesRequest"></a>

### QueryConsensusStatesRequest
QueryConsensusStatesRequest is the request type for the Query/ConsensusStates
RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `chain_name` | [string](#string) |  | client identifier |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination request |






<a name="xibc.core.client.v1.QueryConsensusStatesResponse"></a>

### QueryConsensusStatesResponse
QueryConsensusStatesResponse is the response type for the
Query/ConsensusStates RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `consensus_states` | [ConsensusStateWithHeight](#xibc.core.client.v1.ConsensusStateWithHeight) | repeated | consensus states associated with the identifier |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination response |






<a name="xibc.core.client.v1.QueryRelayersRequest"></a>

### QueryRelayersRequest
QueryRelayersRequest is the request type for the Query/Relayers RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `chain_name` | [string](#string) |  | client identifier |






<a name="xibc.core.client.v1.QueryRelayersResponse"></a>

### QueryRelayersResponse
QueryConsensusStatesResponse is the response type for the Query/Relayers RPC
method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `relayers` | [string](#string) | repeated | relayers address associated with the client |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="xibc.core.client.v1.Query"></a>

### Query
Query provides defines the gRPC querier service

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `ClientState` | [QueryClientStateRequest](#xibc.core.client.v1.QueryClientStateRequest) | [QueryClientStateResponse](#xibc.core.client.v1.QueryClientStateResponse) | ClientState queries an XIBC client. | GET|/ibc/core/client/v1beta1/client_states/{chain_name}|
| `ClientStates` | [QueryClientStatesRequest](#xibc.core.client.v1.QueryClientStatesRequest) | [QueryClientStatesResponse](#xibc.core.client.v1.QueryClientStatesResponse) | ClientStates queries all the XIBC clients of a chain. | GET|/ibc/core/client/v1beta1/client_states|
| `ConsensusState` | [QueryConsensusStateRequest](#xibc.core.client.v1.QueryConsensusStateRequest) | [QueryConsensusStateResponse](#xibc.core.client.v1.QueryConsensusStateResponse) | ConsensusState queries a consensus state associated with a client state at a given height. | GET|/ibc/core/client/v1beta1/consensus_states/{chain_name}/revision/{revision_number}/height/{revision_height}|
| `ConsensusStates` | [QueryConsensusStatesRequest](#xibc.core.client.v1.QueryConsensusStatesRequest) | [QueryConsensusStatesResponse](#xibc.core.client.v1.QueryConsensusStatesResponse) | ConsensusStates queries all the consensus state associated with a given client. | GET|/ibc/core/client/v1beta1/consensus_states/{chain_name}|
| `Relayers` | [QueryRelayersRequest](#xibc.core.client.v1.QueryRelayersRequest) | [QueryRelayersResponse](#xibc.core.client.v1.QueryRelayersResponse) | Relayers queries all the relayers associated with a given client. | GET|/ibc/core/client/v1beta1/relayers/{chain_name}|

 <!-- end services -->



<a name="xibc/core/client/v1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## xibc/core/client/v1/tx.proto



<a name="xibc.core.client.v1.MsgUpdateClient"></a>

### MsgUpdateClient
MsgUpdateClient defines an sdk.Msg to update a XIBC client state using the
given header.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `chain_name` | [string](#string) |  | client unique identifier |
| `header` | [google.protobuf.Any](#google.protobuf.Any) |  | header to update the client |
| `signer` | [string](#string) |  | signer address |






<a name="xibc.core.client.v1.MsgUpdateClientResponse"></a>

### MsgUpdateClientResponse
MsgUpdateClientResponse defines the Msg/UpdateClient response type.





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="xibc.core.client.v1.Msg"></a>

### Msg
Msg defines the xibc/client Msg service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `UpdateClient` | [MsgUpdateClient](#xibc.core.client.v1.MsgUpdateClient) | [MsgUpdateClientResponse](#xibc.core.client.v1.MsgUpdateClientResponse) | UpdateClient defines a rpc handler method for MsgUpdateClient. | |

 <!-- end services -->



<a name="xibc/core/commitment/v1/commitment.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## xibc/core/commitment/v1/commitment.proto



<a name="xibc.core.commitment.v1.MerklePath"></a>

### MerklePath
MerklePath is the path used to verify commitment proofs, which can be an
arbitrary structured object (defined by a commitment type).
MerklePath is represented from root-to-leaf


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key_path` | [string](#string) | repeated |  |






<a name="xibc.core.commitment.v1.MerklePrefix"></a>

### MerklePrefix
MerklePrefix is merkle path prefixed to the key.
The constructed key from the Path and the key will be append(Path.KeyPath,
append(Path.KeyPrefix, key...))


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key_prefix` | [bytes](#bytes) |  |  |






<a name="xibc.core.commitment.v1.MerkleProof"></a>

### MerkleProof
MerkleProof is a wrapper type over a chain of CommitmentProofs.
It demonstrates membership or non-membership for an element or set of
elements, verifiable in conjunction with a known commitment root. Proofs
should be succinct.
MerkleProofs are ordered from leaf-to-root


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `proofs` | [ics23.CommitmentProof](#ics23.CommitmentProof) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="xibc/core/packet/v1/event.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## xibc/core/packet/v1/event.proto



<a name="xibc.core.packet.v1.EventAcknowledgePacket"></a>

### EventAcknowledgePacket
EventAcknowledgePacket is emitted on acknowledgement packet


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `sequence` | [string](#string) |  |  |
| `src_chain` | [string](#string) |  |  |
| `dst_chain` | [string](#string) |  |  |
| `relay_chain` | [string](#string) |  |  |
| `ports` | [string](#string) | repeated |  |
| `data_list` | [bytes](#bytes) | repeated |  |
| `ack` | [bytes](#bytes) |  |  |






<a name="xibc.core.packet.v1.EventRecvPacket"></a>

### EventRecvPacket
EventRecvPacket is emitted on receive packet


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `sequence` | [string](#string) |  |  |
| `src_chain` | [string](#string) |  |  |
| `dst_chain` | [string](#string) |  |  |
| `relay_chain` | [string](#string) |  |  |
| `ports` | [string](#string) | repeated |  |
| `data_list` | [bytes](#bytes) | repeated |  |






<a name="xibc.core.packet.v1.EventSendPacket"></a>

### EventSendPacket
EventSendPacket is emitted on send packet


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `sequence` | [string](#string) |  |  |
| `src_chain` | [string](#string) |  |  |
| `dst_chain` | [string](#string) |  |  |
| `relay_chain` | [string](#string) |  |  |
| `ports` | [string](#string) | repeated |  |
| `data_list` | [bytes](#bytes) | repeated |  |






<a name="xibc.core.packet.v1.EventWriteAck"></a>

### EventWriteAck
EventWriteAck is emitted on receive packet


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `sequence` | [string](#string) |  |  |
| `src_chain` | [string](#string) |  |  |
| `dst_chain` | [string](#string) |  |  |
| `relay_chain` | [string](#string) |  |  |
| `ports` | [string](#string) | repeated |  |
| `data_list` | [bytes](#bytes) | repeated |  |
| `ack` | [bytes](#bytes) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="xibc/core/packet/v1/packet.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## xibc/core/packet/v1/packet.proto



<a name="xibc.core.packet.v1.Acknowledgement"></a>

### Acknowledgement
Acknowledgement is the recommended acknowledgement format to be used by
app-specific protocols.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `results` | [bytes](#bytes) | repeated | the execution results of the packet data list |
| `message` | [string](#string) |  | error message |






<a name="xibc.core.packet.v1.Packet"></a>

### Packet
Packet defines a type that carries data across different chains through XIBC


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `sequence` | [uint64](#uint64) |  | number corresponds to the order of sends and receives, where a Packet with an earlier sequence number must be sent and received before a Packet with a later sequence number. |
| `source_chain` | [string](#string) |  | identifies the chain id of the sending chain. |
| `destination_chain` | [string](#string) |  | identifies the chain id of the receiving chain. |
| `relay_chain` | [string](#string) |  | identifies the chain id of the relay chain. |
| `ports` | [string](#string) | repeated | identifies the ports on the sending chain and destination chain |
| `data_list` | [bytes](#bytes) | repeated | actual opaque bytes transferred directly to the application module |






<a name="xibc.core.packet.v1.PacketState"></a>

### PacketState
PacketState defines the generic type necessary to retrieve and store
packet commitments, acknowledgements, and receipts.
Caller is responsible for knowing the context necessary to interpret this
state as a commitment, acknowledgement, or a receipt.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `source_chain` | [string](#string) |  | the sending chain identifier. |
| `destination_chain` | [string](#string) |  | the receiving chain identifier. |
| `sequence` | [uint64](#uint64) |  | packet sequence. |
| `data` | [bytes](#bytes) |  | embedded data that represents packet state. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="xibc/core/packet/v1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## xibc/core/packet/v1/genesis.proto



<a name="xibc.core.packet.v1.GenesisState"></a>

### GenesisState
GenesisState defines the xibc channel submodule's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `acknowledgements` | [PacketState](#xibc.core.packet.v1.PacketState) | repeated |  |
| `commitments` | [PacketState](#xibc.core.packet.v1.PacketState) | repeated |  |
| `receipts` | [PacketState](#xibc.core.packet.v1.PacketState) | repeated |  |
| `send_sequences` | [PacketSequence](#xibc.core.packet.v1.PacketSequence) | repeated |  |
| `recv_sequences` | [PacketSequence](#xibc.core.packet.v1.PacketSequence) | repeated |  |
| `ack_sequences` | [PacketSequence](#xibc.core.packet.v1.PacketSequence) | repeated |  |






<a name="xibc.core.packet.v1.PacketSequence"></a>

### PacketSequence
PacketSequence defines the genesis type necessary to retrieve and store
next send and receive sequences.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `source_chain` | [string](#string) |  |  |
| `destination_chain` | [string](#string) |  |  |
| `sequence` | [uint64](#uint64) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="xibc/core/packet/v1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## xibc/core/packet/v1/query.proto



<a name="xibc.core.packet.v1.QueryPacketAcknowledgementRequest"></a>

### QueryPacketAcknowledgementRequest
QueryPacketAcknowledgementRequest is the request type for the
Query/PacketAcknowledgement RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `dest_chain` | [string](#string) |  | dest chain name |
| `source_chain` | [string](#string) |  | source chain name |
| `sequence` | [uint64](#uint64) |  | packet sequence |






<a name="xibc.core.packet.v1.QueryPacketAcknowledgementResponse"></a>

### QueryPacketAcknowledgementResponse
QueryPacketAcknowledgementResponse defines the client query response for a
packet which also includes a proof and the height from which the proof was
retrieved


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `acknowledgement` | [bytes](#bytes) |  | packet associated with the request fields |
| `proof` | [bytes](#bytes) |  | merkle proof of existence |
| `proof_height` | [xibc.core.client.v1.Height](#xibc.core.client.v1.Height) |  | height at which the proof was retrieved |






<a name="xibc.core.packet.v1.QueryPacketAcknowledgementsRequest"></a>

### QueryPacketAcknowledgementsRequest
QueryPacketAcknowledgementsRequest is the request type for the
Query/QueryPacketCommitments RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `dest_chain` | [string](#string) |  | dest chain name |
| `source_chain` | [string](#string) |  | source chain name |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination request |






<a name="xibc.core.packet.v1.QueryPacketAcknowledgementsResponse"></a>

### QueryPacketAcknowledgementsResponse
QueryPacketAcknowledgemetsResponse is the request type for the
Query/QueryPacketAcknowledgements RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `acknowledgements` | [PacketState](#xibc.core.packet.v1.PacketState) | repeated |  |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination response |
| `height` | [xibc.core.client.v1.Height](#xibc.core.client.v1.Height) |  | query block height |






<a name="xibc.core.packet.v1.QueryPacketCommitmentRequest"></a>

### QueryPacketCommitmentRequest
QueryPacketCommitmentRequest is the request type for the
QueryPacketCommitment RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `dest_chain` | [string](#string) |  | dest chain name |
| `source_chain` | [string](#string) |  | source chain name |
| `sequence` | [uint64](#uint64) |  | packet sequence |






<a name="xibc.core.packet.v1.QueryPacketCommitmentResponse"></a>

### QueryPacketCommitmentResponse
QueryPacketCommitmentResponse defines the client query response for a packet
which also includes a proof and the height from which the proof was retrieved


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `commitment` | [bytes](#bytes) |  | packet associated with the request fields |
| `proof` | [bytes](#bytes) |  | merkle proof of existence |
| `proof_height` | [xibc.core.client.v1.Height](#xibc.core.client.v1.Height) |  | height at which the proof was retrieved |






<a name="xibc.core.packet.v1.QueryPacketCommitmentsRequest"></a>

### QueryPacketCommitmentsRequest
QueryPacketCommitmentsRequest is the request type for the
Query/QueryPacketCommitments RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `dest_chain` | [string](#string) |  | dest chain name |
| `source_chain` | [string](#string) |  | source chain name |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination request |






<a name="xibc.core.packet.v1.QueryPacketCommitmentsResponse"></a>

### QueryPacketCommitmentsResponse
QueryPacketCommitmentsResponse is the request type for the
Query/QueryPacketCommitments RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `commitments` | [PacketState](#xibc.core.packet.v1.PacketState) | repeated |  |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination response |
| `height` | [xibc.core.client.v1.Height](#xibc.core.client.v1.Height) |  | query block height |






<a name="xibc.core.packet.v1.QueryPacketReceiptRequest"></a>

### QueryPacketReceiptRequest
QueryPacketReceiptRequest is the request type for the Query/PacketReceipt RPC
method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `dest_chain` | [string](#string) |  | dest chain name |
| `source_chain` | [string](#string) |  | source chain name |
| `sequence` | [uint64](#uint64) |  | packet sequence |






<a name="xibc.core.packet.v1.QueryPacketReceiptResponse"></a>

### QueryPacketReceiptResponse
QueryPacketReceiptResponse defines the client query response for a packet
receipt which also includes a proof, and the height from which the proof was
retrieved


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `received` | [bool](#bool) |  | success flag for if receipt exists |
| `proof` | [bytes](#bytes) |  | merkle proof of existence |
| `proof_height` | [xibc.core.client.v1.Height](#xibc.core.client.v1.Height) |  | height at which the proof was retrieved |






<a name="xibc.core.packet.v1.QueryUnreceivedAcksRequest"></a>

### QueryUnreceivedAcksRequest
QueryUnreceivedAcks is the request type for the Query/UnreceivedAcks RPC
method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `dest_chain` | [string](#string) |  | dest chain name |
| `source_chain` | [string](#string) |  | source chain name |
| `packet_ack_sequences` | [uint64](#uint64) | repeated | list of acknowledgement sequences |






<a name="xibc.core.packet.v1.QueryUnreceivedAcksResponse"></a>

### QueryUnreceivedAcksResponse
QueryUnreceivedAcksResponse is the response type for the Query/UnreceivedAcks
RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `sequences` | [uint64](#uint64) | repeated | list of unreceived acknowledgement sequences |
| `height` | [xibc.core.client.v1.Height](#xibc.core.client.v1.Height) |  | query block height |






<a name="xibc.core.packet.v1.QueryUnreceivedPacketsRequest"></a>

### QueryUnreceivedPacketsRequest
QueryUnreceivedPacketsRequest is the request type for the
Query/UnreceivedPackets RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `dest_chain` | [string](#string) |  | dest chain name |
| `source_chain` | [string](#string) |  | source chain name |
| `packet_commitment_sequences` | [uint64](#uint64) | repeated | list of packet sequences |






<a name="xibc.core.packet.v1.QueryUnreceivedPacketsResponse"></a>

### QueryUnreceivedPacketsResponse
QueryUnreceivedPacketsResponse is the response type for the
Query/UnreceivedPacketCommitments RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `sequences` | [uint64](#uint64) | repeated | list of unreceived packet sequences |
| `height` | [xibc.core.client.v1.Height](#xibc.core.client.v1.Height) |  | query block height |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="xibc.core.packet.v1.Query"></a>

### Query
Query provides defines the gRPC querier service

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `PacketCommitment` | [QueryPacketCommitmentRequest](#xibc.core.packet.v1.QueryPacketCommitmentRequest) | [QueryPacketCommitmentResponse](#xibc.core.packet.v1.QueryPacketCommitmentResponse) | PacketCommitment queries a stored packet commitment hash. | GET|/xibc/core/packet/v1beta1/source_chains/{source_chain}/dest_chains/{dest_chain}/packet_commitments/{sequence}|
| `PacketCommitments` | [QueryPacketCommitmentsRequest](#xibc.core.packet.v1.QueryPacketCommitmentsRequest) | [QueryPacketCommitmentsResponse](#xibc.core.packet.v1.QueryPacketCommitmentsResponse) | PacketCommitments returns all the packet commitments hashes associated | GET|/xibc/core/packet/v1beta1/source_chains/{source_chain}/dest_chains/{dest_chain}/packet_commitments|
| `PacketReceipt` | [QueryPacketReceiptRequest](#xibc.core.packet.v1.QueryPacketReceiptRequest) | [QueryPacketReceiptResponse](#xibc.core.packet.v1.QueryPacketReceiptResponse) | PacketReceipt queries if a given packet sequence has been received on the queried chain | GET|/xibc/core/packet/v1beta1/source_chains/{source_chain}/dest_chains/{dest_chain}/packet_receipts/{sequence}|
| `PacketAcknowledgement` | [QueryPacketAcknowledgementRequest](#xibc.core.packet.v1.QueryPacketAcknowledgementRequest) | [QueryPacketAcknowledgementResponse](#xibc.core.packet.v1.QueryPacketAcknowledgementResponse) | PacketAcknowledgement queries a stored packet acknowledgement hash. | GET|/xibc/core/packet/v1beta1/source_chains/{source_chain}/dest_chains/{dest_chain}/packet_acks/{sequence}|
| `PacketAcknowledgements` | [QueryPacketAcknowledgementsRequest](#xibc.core.packet.v1.QueryPacketAcknowledgementsRequest) | [QueryPacketAcknowledgementsResponse](#xibc.core.packet.v1.QueryPacketAcknowledgementsResponse) | PacketAcknowledgements returns all the packet acknowledgements associated | GET|/xibc/core/packet/v1beta1/source_chains/{source_chain}/dest_chains/{dest_chain}/packet_acknowledgements|
| `UnreceivedPackets` | [QueryUnreceivedPacketsRequest](#xibc.core.packet.v1.QueryUnreceivedPacketsRequest) | [QueryUnreceivedPacketsResponse](#xibc.core.packet.v1.QueryUnreceivedPacketsResponse) | UnreceivedPackets returns all the unreceived XIBC packets associated with sequences. | GET|/xibc/core/packet/v1beta1/source_chains/{source_chain}/dest_chains/{dest_chain}/packet_commitments/{packet_commitment_sequences}/unreceived_packets|
| `UnreceivedAcks` | [QueryUnreceivedAcksRequest](#xibc.core.packet.v1.QueryUnreceivedAcksRequest) | [QueryUnreceivedAcksResponse](#xibc.core.packet.v1.QueryUnreceivedAcksResponse) | UnreceivedAcks returns all the unreceived XIBC acknowledgements associated with sequences. | GET|/xibc/core/packet/v1beta1/source_chains/{source_chain}/dest_chains/{dest_chain}/packet_commitments/{packet_ack_sequences}/unreceived_acks|

 <!-- end services -->



<a name="xibc/core/packet/v1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## xibc/core/packet/v1/tx.proto



<a name="xibc.core.packet.v1.MsgAcknowledgement"></a>

### MsgAcknowledgement
MsgAcknowledgement receives incoming XIBC acknowledgement


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `packet` | [Packet](#xibc.core.packet.v1.Packet) |  |  |
| `acknowledgement` | [bytes](#bytes) |  |  |
| `proof_acked` | [bytes](#bytes) |  |  |
| `proof_height` | [xibc.core.client.v1.Height](#xibc.core.client.v1.Height) |  |  |
| `signer` | [string](#string) |  |  |






<a name="xibc.core.packet.v1.MsgAcknowledgementResponse"></a>

### MsgAcknowledgementResponse
MsgAcknowledgementResponse defines the Msg/Acknowledgement response type.






<a name="xibc.core.packet.v1.MsgRecvPacket"></a>

### MsgRecvPacket
MsgRecvPacket receives incoming XIBC packet


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `packet` | [Packet](#xibc.core.packet.v1.Packet) |  |  |
| `proof_commitment` | [bytes](#bytes) |  |  |
| `proof_height` | [xibc.core.client.v1.Height](#xibc.core.client.v1.Height) |  |  |
| `signer` | [string](#string) |  |  |






<a name="xibc.core.packet.v1.MsgRecvPacketResponse"></a>

### MsgRecvPacketResponse
MsgRecvPacketResponse defines the Msg/RecvPacket response type.





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="xibc.core.packet.v1.Msg"></a>

### Msg
Msg defines the xibc/packet Msg service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `RecvPacket` | [MsgRecvPacket](#xibc.core.packet.v1.MsgRecvPacket) | [MsgRecvPacketResponse](#xibc.core.packet.v1.MsgRecvPacketResponse) | RecvPacket defines a rpc handler method for MsgRecvPacket. | |
| `Acknowledgement` | [MsgAcknowledgement](#xibc.core.packet.v1.MsgAcknowledgement) | [MsgAcknowledgementResponse](#xibc.core.packet.v1.MsgAcknowledgementResponse) | Acknowledgement defines a rpc handler method for MsgAcknowledgement. | |

 <!-- end services -->



<a name="xibc/core/types/v1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## xibc/core/types/v1/genesis.proto



<a name="xibc.core.types.v1.GenesisState"></a>

### GenesisState
GenesisState defines the xibc module's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `client_genesis` | [xibc.core.client.v1.GenesisState](#xibc.core.client.v1.GenesisState) |  | Clients genesis state |
| `packet_genesis` | [xibc.core.packet.v1.GenesisState](#xibc.core.packet.v1.GenesisState) |  | Packet genesis state |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |

