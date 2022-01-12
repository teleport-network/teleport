package gabci

import (
	"context"
	"strings"

	"github.com/teleport-network/teleport/grpc_abci/types"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	gogogrpc "github.com/gogo/protobuf/grpc"

	abci "github.com/tendermint/tendermint/abci/types"
)

func RegisterGRPCABCIQuery(qrt gogogrpc.Server, app abci.Application) {
	types.RegisterServer(qrt, NewGRPCABCIQuery(app))
}

// GRPCABCIQuery is a GRPC wrapper for Application
type GRPCABCIQuery struct {
	app abci.Application
}

func NewGRPCABCIQuery(app abci.Application) *GRPCABCIQuery {
	return &GRPCABCIQuery{app}
}

func (app *GRPCABCIQuery) Info(ctx context.Context, req *abci.RequestInfo) (*abci.ResponseInfo, error) {
	areq := abci.RequestInfo{
		Version:      req.Version,
		BlockVersion: req.BlockVersion,
		P2PVersion:   req.P2PVersion,
	}
	res := app.app.Info(areq)
	return &res, nil
}

func (app *GRPCABCIQuery) Query(ctx context.Context, req *abci.RequestQuery) (*abci.ResponseQuery, error) {
	var res abci.ResponseQuery
	path := splitPath(req.Path)
	if len(path) == 0 {
		res = sdkerrors.QueryResult(sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "no query path provided"))
		return &res, nil
	}
	if path[0] != "app" && path[0] != "store" && path[0] != "p2p" && path[0] != "custom" {
		res = sdkerrors.QueryResult(sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown query path"))
		return &res, nil
	}
	res = app.app.Query(*req)
	return &res, nil
}

// splitPath splits a string path using the delimiter '/'.
//
// e.g. "this/is/funny" becomes []string{"this", "is", "funny"}
func splitPath(requestPath string) (path []string) {
	path = strings.Split(requestPath, "/")

	// first element is empty string
	if len(path) > 0 && path[0] == "" {
		path = path[1:]
	}

	return path
}
