package xibc

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"

	"github.com/teleport-network/teleport/x/xibc"
	"github.com/teleport-network/teleport/x/xibc/client/cli"
	clienttypes "github.com/teleport-network/teleport/x/xibc/core/client/types"
	"github.com/teleport-network/teleport/x/xibc/core/host"
	packettypes "github.com/teleport-network/teleport/x/xibc/core/packet/types"
	"github.com/teleport-network/teleport/x/xibc/keeper"
	"github.com/teleport-network/teleport/x/xibc/simulation"
	"github.com/teleport-network/teleport/x/xibc/types"
)

var (
	_ module.AppModule           = AppModule{}
	_ module.AppModuleBasic      = AppModuleBasic{}
	_ module.AppModuleSimulation = AppModule{}
)

// AppModuleBasic defines the basic application module used by the xibc module.
type AppModuleBasic struct{}

var _ module.AppModuleBasic = AppModuleBasic{}

// Name returns the xibc module's name.
func (AppModuleBasic) Name() string {
	return host.ModuleName
}

// RegisterLegacyAminoCodec does nothing. XIBC does not support amino.
func (AppModuleBasic) RegisterLegacyAminoCodec(*codec.LegacyAmino) {}

// DefaultGenesis returns default genesis state as raw bytes for the xibc module.
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(types.DefaultGenesisState())
}

// ValidateGenesis performs genesis state validation for the xibc module.
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, config client.TxEncodingConfig, bz json.RawMessage) error {
	var gs types.GenesisState
	if err := cdc.UnmarshalJSON(bz, &gs); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", host.ModuleName, err)
	}

	return gs.Validate()
}

// RegisterRESTRoutes does nothing. XIBC does not support legacy REST routes.
func (AppModuleBasic) RegisterRESTRoutes(client.Context, *mux.Router) {}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the xibc module.
func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	_ = clienttypes.RegisterQueryHandlerClient(context.Background(), mux, clienttypes.NewQueryClient(clientCtx))
	_ = packettypes.RegisterQueryHandlerClient(context.Background(), mux, packettypes.NewQueryClient(clientCtx))
}

// GetTxCmd returns the root tx command for the xibc module.
func (AppModuleBasic) GetTxCmd() *cobra.Command {
	return cli.GetTxCmd()
}

// GetQueryCmd returns no root query command for the xibc module.
func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	return cli.GetQueryCmd()
}

// RegisterInterfaces registers module concrete types into protobuf Any.
func (AppModuleBasic) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	types.RegisterInterfaces(registry)
}

// AppModule implements an application module for the xibc module.
type AppModule struct {
	AppModuleBasic
	keeper *keeper.Keeper

	// create localhost by default
	createLocalhost bool
}

// NewAppModule creates a new AppModule object
func NewAppModule(k *keeper.Keeper) AppModule {
	return AppModule{
		keeper: k,
	}
}

// Name returns the xibc module's name.
func (AppModule) Name() string {
	return host.ModuleName
}

// RegisterInvariants registers the xibc module invariants.
func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {
	// TODO:
}

// Route returns the message routing key for the xibc module.
func (am AppModule) Route() sdk.Route {
	return sdk.Route{}
}

// QuerierRoute returns the xibc module's querier route name.
func (AppModule) QuerierRoute() string {
	return host.QuerierRoute
}

// LegacyQuerierHandler returns nil. XIBC does not support the legacy querier.
func (am AppModule) LegacyQuerierHandler(legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return nil
}

// RegisterServices registers module services.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	clienttypes.RegisterMsgServer(cfg.MsgServer(), am.keeper)
	packettypes.RegisterMsgServer(cfg.MsgServer(), am.keeper)
	types.RegisterQueryService(cfg.QueryServer(), am.keeper)
}

// InitGenesis performs genesis initialization for the xibc module. It returns
// no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, bz json.RawMessage) []abci.ValidatorUpdate {
	var gs types.GenesisState
	if err := cdc.UnmarshalJSON(bz, &gs); err != nil {
		panic(fmt.Sprintf("failed to unmarshal %s genesis state: %s", host.ModuleName, err))
	}
	xibc.InitGenesis(ctx, *am.keeper, am.createLocalhost, &gs)
	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns the exported genesis state as raw bytes for the xibc
// module.
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(xibc.ExportGenesis(ctx, *am.keeper))
}

// ConsensusVersion implements AppModule/ConsensusVersion.
func (AppModule) ConsensusVersion() uint64 { return 1 }

// BeginBlock returns the begin blocker for the xibc module.
func (am AppModule) BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock) {
}

// EndBlock returns the end blocker for the xibc module. It returns no validator
// updates.
func (am AppModule) EndBlock(ctx sdk.Context, req abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

// ____________________________________________________________________________

// AppModuleSimulation functions

// GenerateGenesisState creates a randomized GenState of the xibc module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	simulation.RandomizedGenState(simState)
}

// ProposalContents doesn't return any content functions for governance proposals.
func (AppModule) ProposalContents(_ module.SimulationState) []simtypes.WeightedProposalContent {
	return nil
}

// RandomizedParams returns nil since XIBC doesn't register parameter changes.
func (AppModule) RandomizedParams(_ *rand.Rand) []simtypes.ParamChange {
	return nil
}

// RegisterStoreDecoder registers a decoder for xibc module's types
func (am AppModule) RegisterStoreDecoder(sdr sdk.StoreDecoderRegistry) {
	sdr[host.StoreKey] = simulation.NewDecodeStore(*am.keeper)
}

// WeightedOperations returns the all the xibc module operations with their respective weights.
func (am AppModule) WeightedOperations(_ module.SimulationState) []simtypes.WeightedOperation {
	return nil
}
