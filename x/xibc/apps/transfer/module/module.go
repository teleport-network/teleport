package transfer

import (
	"encoding/json"
	"math/rand"

	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"

	"github.com/teleport-network/teleport/x/xibc/apps/transfer"
	"github.com/teleport-network/teleport/x/xibc/apps/transfer/keeper"
	"github.com/teleport-network/teleport/x/xibc/apps/transfer/types"
	packettypes "github.com/teleport-network/teleport/x/xibc/core/packet/types"
	porttypes "github.com/teleport-network/teleport/x/xibc/core/routing/types"
	"github.com/teleport-network/teleport/x/xibc/simulation"
)

var (
	_ module.AppModule      = AppModule{}
	_ porttypes.XIBCModule  = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// AppModuleBasic is the XIBC Transfer AppModuleBasic
type AppModuleBasic struct{}

func (a AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return nil
}

func (a AppModuleBasic) ValidateGenesis(jsonCodec codec.JSONCodec, config client.TxEncodingConfig, message json.RawMessage) error {
	return nil
}

func (AppModuleBasic) RegisterRESTRoutes(clientCtx client.Context, rtr *mux.Router) {}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the xibc-transfer module.
func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {}

func (a AppModuleBasic) GetTxCmd() *cobra.Command {
	return nil
}

// GetQueryCmd returns no root query command for the xibc-transfer module.
func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	return nil
}

func (a AppModuleBasic) Name() string {
	return types.ModuleName
}

// RegisterLegacyAminoCodec implements AppModuleBasic interface
func (a AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {}

// RegisterInterfaces registers module concrete types into protobuf Any.
func (a AppModuleBasic) RegisterInterfaces(registry codectypes.InterfaceRegistry) {}

// AppModule represents the AppModule for this module
type AppModule struct {
	AppModuleBasic
	keeper keeper.Keeper
	ak     authkeeper.AccountKeeper
}

func (a AppModule) GenerateGenesisState(simState *module.SimulationState) {
	simulation.RandomizedGenState(simState)
}

func (a AppModule) ProposalContents(simState module.SimulationState) []simtypes.WeightedProposalContent {
	return nil
}

func (a AppModule) RandomizedParams(r *rand.Rand) []simtypes.ParamChange {
	return nil
}

func (a AppModule) RegisterStoreDecoder(sdr sdk.StoreDecoderRegistry) {}

func (a AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	panic("implement me")
}

// NewAppModule creates a new xibc-transfer module
func NewAppModule(k keeper.Keeper, ak authkeeper.AccountKeeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         k,
		ak:             ak,
	}
}

func (a AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) []abci.ValidatorUpdate {
	transfer.InitGenesis(ctx, a.ak)
	return []abci.ValidatorUpdate{}
}

func (a AppModule) ExportGenesis(context sdk.Context, jsonCodec codec.JSONCodec) json.RawMessage {
	return nil
}

func (a AppModule) RegisterInvariants(registry sdk.InvariantRegistry) {}

func (a AppModule) Route() sdk.Route {
	return sdk.Route{}
}

func (a AppModule) QuerierRoute() string {
	return ""
}

func (a AppModule) LegacyQuerierHandler(amino *codec.LegacyAmino) sdk.Querier {
	return nil
}

// RegisterServices registers module services.
func (am AppModule) RegisterServices(cfg module.Configurator) {}

// ConsensusVersion implements AppModule/ConsensusVersion.
func (am AppModule) ConsensusVersion() uint64 { return 1 }

// BeginBlock implements the AppModule interface
func (am AppModule) BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock) {}

// EndBlock implements the AppModule interface
func (am AppModule) EndBlock(ctx sdk.Context, req abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

func (a AppModule) OnRecvPacket(ctx sdk.Context, packetData []byte) (*sdk.Result, packettypes.Result, error) {
	var data types.FungibleTokenPacketData
	err := data.DecodeBytes(packetData)
	if err != nil {
		return nil, packettypes.Result{}, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unmarshal transfer packet failed")
	}
	if len(data.String()) == 0 {
		return nil, packettypes.Result{}, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal transfer packet data")
	}
	result, err := a.keeper.OnRecvPacket(ctx, data)
	if err != nil {
		return nil, packettypes.Result{}, err
	}

	_ = ctx.EventManager().EmitTypedEvent(&types.EventRecvPacket{
		SrcChain:  data.SrcChain,
		DestChain: data.DestChain,
		Sender:    data.Sender,
		Receiver:  data.Receiver,
		Token:     data.Token,
		Amount:    data.Amount,
		Result:    result.Result,
		Message:   result.Message,
	})

	// NOTE: acknowledgement will be written synchronously during XIBC handler execution.
	return &sdk.Result{
		Events: ctx.EventManager().Events().ToABCIEvents(),
	}, result, nil
}

func (a AppModule) OnAcknowledgementPacket(ctx sdk.Context, packetData []byte, result []byte) (*sdk.Result, error) {
	var data types.FungibleTokenPacketData
	err := data.DecodeBytes(packetData)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unmarshal transfer packet failed")
	}
	if len(data.String()) == 0 {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal transfer packet data")
	}
	if err := a.keeper.OnAcknowledgementPacket(ctx, data, result); err != nil {
		return nil, err
	}

	_ = ctx.EventManager().EmitTypedEvent(&types.EventAckPacket{
		SrcChain:  data.SrcChain,
		DestChain: data.DestChain,
		Sender:    data.Sender,
		Receiver:  data.Receiver,
		Token:     data.Token,
		Amount:    data.Amount,
		Result:    result,
	})

	return &sdk.Result{
		Events: ctx.EventManager().Events().ToABCIEvents(),
	}, nil
}
