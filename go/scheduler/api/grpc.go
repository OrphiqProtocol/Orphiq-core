package api

import (
	"context"

	"google.golang.org/grpc"

	cmnGrpc "github.com/oasisprotocol/oasis-core/go/common/grpc"
	"github.com/oasisprotocol/oasis-core/go/common/pubsub"
)

var (
	// serviceName is the gRPC service name.
	serviceName = cmnGrpc.NewServiceName("Scheduler")

	// methodGetValidators is the GetValidators method.
	methodGetValidators = serviceName.NewMethod("GetValidators", int64(0))
	// methodGetCommittees is the GetCommittees method.
	methodGetCommittees = serviceName.NewMethod("GetCommittees", GetCommitteesRequest{})
	// methodStateToGenesis is the StateToGenesis method.
	methodStateToGenesis = serviceName.NewMethod("StateToGenesis", int64(0))
	// methodConsensusParameters is the ConsensusParameters method.
	methodConsensusParameters = serviceName.NewMethod("ConsensusParameters", int64(0))

	// methodWatchCommittees is the WatchCommittees method.
	methodWatchCommittees = serviceName.NewMethod("WatchCommittees", nil)

	// serviceDesc is the gRPC service descriptor.
	serviceDesc = grpc.ServiceDesc{
		ServiceName: string(serviceName),
		HandlerType: (*Backend)(nil),
		Methods: []grpc.MethodDesc{
			{
				MethodName: methodGetValidators.ShortName(),
				Handler:    handlerGetValidators,
			},
			{
				MethodName: methodGetCommittees.ShortName(),
				Handler:    handlerGetCommittees,
			},
			{
				MethodName: methodStateToGenesis.ShortName(),
				Handler:    handlerStateToGenesis,
			},
			{
				MethodName: methodConsensusParameters.ShortName(),
				Handler:    handlerConsensusParameters,
			},
		},
		Streams: []grpc.StreamDesc{
			{
				StreamName:    methodWatchCommittees.ShortName(),
				Handler:       handlerWatchCommittees,
				ServerStreams: true,
			},
		},
	}
)

func handlerGetValidators(
	srv any,
	ctx context.Context,
	dec func(any) error,
	interceptor grpc.UnaryServerInterceptor,
) (any, error) {
	var height int64
	if err := dec(&height); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(Backend).GetValidators(ctx, height)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: methodGetValidators.FullName(),
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(Backend).GetValidators(ctx, req.(int64))
	}
	return interceptor(ctx, height, info, handler)
}

func handlerGetCommittees(
	srv any,
	ctx context.Context,
	dec func(any) error,
	interceptor grpc.UnaryServerInterceptor,
) (any, error) {
	var req GetCommitteesRequest
	if err := dec(&req); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(Backend).GetCommittees(ctx, &req)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: methodGetCommittees.FullName(),
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(Backend).GetCommittees(ctx, req.(*GetCommitteesRequest))
	}
	return interceptor(ctx, &req, info, handler)
}

func handlerStateToGenesis(
	srv any,
	ctx context.Context,
	dec func(any) error,
	interceptor grpc.UnaryServerInterceptor,
) (any, error) {
	var height int64
	if err := dec(&height); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(Backend).StateToGenesis(ctx, height)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: methodStateToGenesis.FullName(),
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(Backend).StateToGenesis(ctx, req.(int64))
	}
	return interceptor(ctx, height, info, handler)
}

func handlerConsensusParameters(
	srv any,
	ctx context.Context,
	dec func(any) error,
	interceptor grpc.UnaryServerInterceptor,
) (any, error) {
	var height int64
	if err := dec(&height); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(Backend).ConsensusParameters(ctx, height)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: methodConsensusParameters.FullName(),
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(Backend).ConsensusParameters(ctx, req.(int64))
	}
	return interceptor(ctx, height, info, handler)
}

func handlerWatchCommittees(srv any, stream grpc.ServerStream) error {
	if err := stream.RecvMsg(nil); err != nil {
		return err
	}

	ctx := stream.Context()
	ch, sub, err := srv.(Backend).WatchCommittees(ctx)
	if err != nil {
		return err
	}
	defer sub.Close()

	for {
		select {
		case c, ok := <-ch:
			if !ok {
				return nil
			}

			if err := stream.SendMsg(c); err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// RegisterService registers a new scheduler service with the given gRPC server.
func RegisterService(server *grpc.Server, service Backend) {
	server.RegisterService(&serviceDesc, service)
}

// Client is a gRPC scheduler client.
type Client struct {
	conn *grpc.ClientConn
}

// NewClient creates a new gRPC scheduler client.
func NewClient(c *grpc.ClientConn) *Client {
	return &Client{
		conn: c,
	}
}

func (c *Client) GetValidators(ctx context.Context, height int64) ([]*Validator, error) {
	var rsp []*Validator
	if err := c.conn.Invoke(ctx, methodGetValidators.FullName(), height, &rsp); err != nil {
		return nil, err
	}
	return rsp, nil
}

func (c *Client) GetCommittees(ctx context.Context, request *GetCommitteesRequest) ([]*Committee, error) {
	var rsp []*Committee
	if err := c.conn.Invoke(ctx, methodGetCommittees.FullName(), request, &rsp); err != nil {
		return nil, err
	}
	return rsp, nil
}

func (c *Client) StateToGenesis(ctx context.Context, height int64) (*Genesis, error) {
	var rsp Genesis
	if err := c.conn.Invoke(ctx, methodStateToGenesis.FullName(), height, &rsp); err != nil {
		return nil, err
	}
	return &rsp, nil
}

func (c *Client) ConsensusParameters(ctx context.Context, height int64) (*ConsensusParameters, error) {
	var rsp ConsensusParameters
	if err := c.conn.Invoke(ctx, methodConsensusParameters.FullName(), height, &rsp); err != nil {
		return nil, err
	}
	return &rsp, nil
}

func (c *Client) WatchCommittees(ctx context.Context) (<-chan *Committee, pubsub.ClosableSubscription, error) {
	ctx, sub := pubsub.NewContextSubscription(ctx)

	stream, err := c.conn.NewStream(ctx, &serviceDesc.Streams[0], methodWatchCommittees.FullName())
	if err != nil {
		return nil, nil, err
	}
	if err = stream.SendMsg(nil); err != nil {
		return nil, nil, err
	}
	if err = stream.CloseSend(); err != nil {
		return nil, nil, err
	}

	ch := make(chan *Committee)
	go func() {
		defer close(ch)

		for {
			var ev Committee
			if serr := stream.RecvMsg(&ev); serr != nil {
				return
			}

			select {
			case ch <- &ev:
			case <-ctx.Done():
				return
			}
		}
	}()

	return ch, sub, nil
}

func (c *Client) Cleanup() {
}
