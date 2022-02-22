package types

import (
	"context"

	gogogrpc "github.com/gogo/protobuf/grpc"
	"google.golang.org/grpc"

	abci "github.com/tendermint/tendermint/abci/types"
)

// ABCIQueryClient is the client API for ABCI query service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ABCIQueryClient interface {
	Info(ctx context.Context, in *abci.RequestInfo, opts ...grpc.CallOption) (*abci.ResponseInfo, error)
	Query(ctx context.Context, in *abci.RequestQuery, opts ...grpc.CallOption) (*abci.ResponseQuery, error)
}

type abciQueryClient struct {
	cc *grpc.ClientConn
}

func NewABCIQueryClient(cc *grpc.ClientConn) ABCIQueryClient {
	return &abciQueryClient{cc}
}

func (c *abciQueryClient) Info(ctx context.Context, in *abci.RequestInfo, opts ...grpc.CallOption) (*abci.ResponseInfo, error) {
	out := new(abci.ResponseInfo)
	if err := c.cc.Invoke(ctx, "/tendermint.abci.ABCIApplication/Info", in, out, opts...); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *abciQueryClient) Query(ctx context.Context, in *abci.RequestQuery, opts ...grpc.CallOption) (*abci.ResponseQuery, error) {
	out := new(abci.ResponseQuery)
	if err := c.cc.Invoke(ctx, "/tendermint.abci.ABCIApplication/Query", in, out, opts...); err != nil {
		return nil, err
	}
	return out, nil
}

// ABCIQueryServer is the server API for ABCI query service.
type ABCIQueryServer interface {
	Info(context.Context, *abci.RequestInfo) (*abci.ResponseInfo, error)
	Query(context.Context, *abci.RequestQuery) (*abci.ResponseQuery, error)
}

func RegisterServer(qrt gogogrpc.Server, srv ABCIQueryServer) {
	qrt.RegisterService(&_ABCIApplication_serviceDesc, srv)
}

func _ABCIApplication_Info_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(abci.RequestInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ABCIQueryServer).Info(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/tendermint.abci.ABCIApplication/Info",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ABCIQueryServer).Info(ctx, req.(*abci.RequestInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _ABCIApplication_Query_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(abci.RequestQuery)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ABCIQueryServer).Query(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/tendermint.abci.ABCIApplication/Query",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ABCIQueryServer).Query(ctx, req.(*abci.RequestQuery))
	}
	return interceptor(ctx, in, info, handler)
}

var _ABCIApplication_serviceDesc = grpc.ServiceDesc{
	ServiceName: "tendermint.abci.ABCIApplication",
	HandlerType: (*ABCIQueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Info",
			Handler:    _ABCIApplication_Info_Handler,
		},
		{
			MethodName: "Query",
			Handler:    _ABCIApplication_Query_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "tendermint/abci/types.proto",
}
