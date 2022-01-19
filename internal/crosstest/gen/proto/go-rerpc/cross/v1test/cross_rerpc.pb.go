// Code generated by protoc-gen-go-rerpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-rerpc v0.0.1
// - protoc              v3.17.3
// source: cross/v1test/cross.proto

package crossv1test

import (
	context "context"
	errors "errors"
	rerpc "github.com/rerpc/rerpc"
	clientstream "github.com/rerpc/rerpc/clientstream"
	handlerstream "github.com/rerpc/rerpc/handlerstream"
	v1test "github.com/rerpc/rerpc/internal/crosstest/gen/proto/go/cross/v1test"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the
// rerpc package are compatible. If you get a compiler error that this constant
// isn't defined, this code was generated with a version of rerpc newer than the
// one compiled into your binary. You can fix the problem by either regenerating
// this code with an older version of rerpc or updating the rerpc version
// compiled into your binary.
const _ = rerpc.SupportsCodeGenV0 // requires reRPC v0.0.1 or later

// SimpleCrossServiceClient is a client for the cross.v1test.CrossService
// service.
type SimpleCrossServiceClient interface {
	Ping(context.Context, *v1test.PingRequest) (*v1test.PingResponse, error)
	Fail(context.Context, *v1test.FailRequest) (*v1test.FailResponse, error)
	Sum(context.Context) *clientstream.Client[v1test.SumRequest, v1test.SumResponse]
	CountUp(context.Context, *v1test.CountUpRequest) (*clientstream.Server[v1test.CountUpResponse], error)
	CumSum(context.Context) *clientstream.Bidirectional[v1test.CumSumRequest, v1test.CumSumResponse]
}

// FullCrossServiceClient is a client for the cross.v1test.CrossService service.
// It's more complex than SimpleCrossServiceClient, but it gives callers more
// fine-grained control (e.g., sending and receiving headers).
type FullCrossServiceClient interface {
	Ping(context.Context, *rerpc.Request[v1test.PingRequest]) (*rerpc.Response[v1test.PingResponse], error)
	Fail(context.Context, *rerpc.Request[v1test.FailRequest]) (*rerpc.Response[v1test.FailResponse], error)
	Sum(context.Context) *clientstream.Client[v1test.SumRequest, v1test.SumResponse]
	CountUp(context.Context, *rerpc.Request[v1test.CountUpRequest]) (*clientstream.Server[v1test.CountUpResponse], error)
	CumSum(context.Context) *clientstream.Bidirectional[v1test.CumSumRequest, v1test.CumSumResponse]
}

// CrossServiceClient is a client for the cross.v1test.CrossService service.
type CrossServiceClient struct {
	client fullCrossServiceClient
}

var _ SimpleCrossServiceClient = (*CrossServiceClient)(nil)

// NewCrossServiceClient constructs a client for the cross.v1test.CrossService
// service.
//
// The URL supplied here should be the base URL for the gRPC server (e.g.,
// https://api.acme.com or https://acme.com/grpc).
func NewCrossServiceClient(baseURL string, doer rerpc.Doer, opts ...rerpc.ClientOption) (*CrossServiceClient, error) {
	baseURL = strings.TrimRight(baseURL, "/")
	pingFunc, err := rerpc.NewClientFunc[v1test.PingRequest, v1test.PingResponse](
		doer,
		baseURL,
		"cross.v1test", // protobuf package
		"CrossService", // protobuf service
		"Ping",         // protobuf method
		opts...,
	)
	if err != nil {
		return nil, err
	}
	failFunc, err := rerpc.NewClientFunc[v1test.FailRequest, v1test.FailResponse](
		doer,
		baseURL,
		"cross.v1test", // protobuf package
		"CrossService", // protobuf service
		"Fail",         // protobuf method
		opts...,
	)
	if err != nil {
		return nil, err
	}
	sumFunc, err := rerpc.NewClientStream(
		doer,
		rerpc.StreamTypeClient,
		baseURL,
		"cross.v1test", // protobuf package
		"CrossService", // protobuf service
		"Sum",          // protobuf method
		opts...,
	)
	if err != nil {
		return nil, err
	}
	countUpFunc, err := rerpc.NewClientStream(
		doer,
		rerpc.StreamTypeServer,
		baseURL,
		"cross.v1test", // protobuf package
		"CrossService", // protobuf service
		"CountUp",      // protobuf method
		opts...,
	)
	if err != nil {
		return nil, err
	}
	cumSumFunc, err := rerpc.NewClientStream(
		doer,
		rerpc.StreamTypeBidirectional,
		baseURL,
		"cross.v1test", // protobuf package
		"CrossService", // protobuf service
		"CumSum",       // protobuf method
		opts...,
	)
	if err != nil {
		return nil, err
	}
	return &CrossServiceClient{client: fullCrossServiceClient{
		ping:    pingFunc,
		fail:    failFunc,
		sum:     sumFunc,
		countUp: countUpFunc,
		cumSum:  cumSumFunc,
	}}, nil
}

// Ping calls cross.v1test.CrossService.Ping.
func (c *CrossServiceClient) Ping(ctx context.Context, req *v1test.PingRequest) (*v1test.PingResponse, error) {
	res, err := c.client.Ping(ctx, rerpc.NewRequest(req))
	if err != nil {
		return nil, err
	}
	return res.Msg, nil
}

// Fail calls cross.v1test.CrossService.Fail.
func (c *CrossServiceClient) Fail(ctx context.Context, req *v1test.FailRequest) (*v1test.FailResponse, error) {
	res, err := c.client.Fail(ctx, rerpc.NewRequest(req))
	if err != nil {
		return nil, err
	}
	return res.Msg, nil
}

// Sum calls cross.v1test.CrossService.Sum.
func (c *CrossServiceClient) Sum(ctx context.Context) *clientstream.Client[v1test.SumRequest, v1test.SumResponse] {
	return c.client.Sum(ctx)
}

// CountUp calls cross.v1test.CrossService.CountUp.
func (c *CrossServiceClient) CountUp(ctx context.Context, req *v1test.CountUpRequest) (*clientstream.Server[v1test.CountUpResponse], error) {
	return c.client.CountUp(ctx, rerpc.NewRequest(req))
}

// CumSum calls cross.v1test.CrossService.CumSum.
func (c *CrossServiceClient) CumSum(ctx context.Context) *clientstream.Bidirectional[v1test.CumSumRequest, v1test.CumSumResponse] {
	return c.client.CumSum(ctx)
}

// Full exposes the underlying generic client. Use it if you need finer control
// (e.g., sending and receiving headers).
func (c *CrossServiceClient) Full() FullCrossServiceClient {
	return &c.client
}

type fullCrossServiceClient struct {
	ping    func(context.Context, *rerpc.Request[v1test.PingRequest]) (*rerpc.Response[v1test.PingResponse], error)
	fail    func(context.Context, *rerpc.Request[v1test.FailRequest]) (*rerpc.Response[v1test.FailResponse], error)
	sum     rerpc.StreamFunc
	countUp rerpc.StreamFunc
	cumSum  rerpc.StreamFunc
}

var _ FullCrossServiceClient = (*fullCrossServiceClient)(nil)

// Ping calls cross.v1test.CrossService.Ping.
func (c *fullCrossServiceClient) Ping(ctx context.Context, req *rerpc.Request[v1test.PingRequest]) (*rerpc.Response[v1test.PingResponse], error) {
	return c.ping(ctx, req)
}

// Fail calls cross.v1test.CrossService.Fail.
func (c *fullCrossServiceClient) Fail(ctx context.Context, req *rerpc.Request[v1test.FailRequest]) (*rerpc.Response[v1test.FailResponse], error) {
	return c.fail(ctx, req)
}

// Sum calls cross.v1test.CrossService.Sum.
func (c *fullCrossServiceClient) Sum(ctx context.Context) *clientstream.Client[v1test.SumRequest, v1test.SumResponse] {
	_, sender, receiver := c.sum(ctx)
	return clientstream.NewClient[v1test.SumRequest, v1test.SumResponse](sender, receiver)
}

// CountUp calls cross.v1test.CrossService.CountUp.
func (c *fullCrossServiceClient) CountUp(ctx context.Context, req *rerpc.Request[v1test.CountUpRequest]) (*clientstream.Server[v1test.CountUpResponse], error) {
	_, sender, receiver := c.countUp(ctx)
	if err := sender.Send(req.Any()); err != nil {
		_ = sender.Close(err)
		_ = receiver.Close()
		return nil, err
	}
	if err := sender.Close(nil); err != nil {
		_ = receiver.Close()
		return nil, err
	}
	return clientstream.NewServer[v1test.CountUpResponse](receiver), nil
}

// CumSum calls cross.v1test.CrossService.CumSum.
func (c *fullCrossServiceClient) CumSum(ctx context.Context) *clientstream.Bidirectional[v1test.CumSumRequest, v1test.CumSumResponse] {
	_, sender, receiver := c.cumSum(ctx)
	return clientstream.NewBidirectional[v1test.CumSumRequest, v1test.CumSumResponse](sender, receiver)
}

// FullCrossServiceServer is a server for the cross.v1test.CrossService service.
type FullCrossServiceServer interface {
	Ping(context.Context, *rerpc.Request[v1test.PingRequest]) (*rerpc.Response[v1test.PingResponse], error)
	Fail(context.Context, *rerpc.Request[v1test.FailRequest]) (*rerpc.Response[v1test.FailResponse], error)
	Sum(context.Context, *handlerstream.Client[v1test.SumRequest, v1test.SumResponse]) error
	CountUp(context.Context, *rerpc.Request[v1test.CountUpRequest], *handlerstream.Server[v1test.CountUpResponse]) error
	CumSum(context.Context, *handlerstream.Bidirectional[v1test.CumSumRequest, v1test.CumSumResponse]) error
}

// SimpleCrossServiceServer is a server for the cross.v1test.CrossService
// service. It's a simpler interface than FullCrossServiceServer but doesn't
// provide header access.
type SimpleCrossServiceServer interface {
	Ping(context.Context, *v1test.PingRequest) (*v1test.PingResponse, error)
	Fail(context.Context, *v1test.FailRequest) (*v1test.FailResponse, error)
	Sum(context.Context, *handlerstream.Client[v1test.SumRequest, v1test.SumResponse]) error
	CountUp(context.Context, *v1test.CountUpRequest, *handlerstream.Server[v1test.CountUpResponse]) error
	CumSum(context.Context, *handlerstream.Bidirectional[v1test.CumSumRequest, v1test.CumSumResponse]) error
}

// NewFullCrossServiceHandler wraps each method on the service implementation in
// a rerpc.Handler. The returned slice can be passed to rerpc.NewServeMux.
func NewFullCrossServiceHandler(svc FullCrossServiceServer, opts ...rerpc.HandlerOption) []rerpc.Handler {
	handlers := make([]rerpc.Handler, 0, 5)

	ping := rerpc.NewUnaryHandler(
		"cross.v1test", // protobuf package
		"CrossService", // protobuf service
		"Ping",         // protobuf method
		svc.Ping,
		opts...,
	)
	handlers = append(handlers, *ping)

	fail := rerpc.NewUnaryHandler(
		"cross.v1test", // protobuf package
		"CrossService", // protobuf service
		"Fail",         // protobuf method
		svc.Fail,
		opts...,
	)
	handlers = append(handlers, *fail)

	sum := rerpc.NewStreamingHandler(
		rerpc.StreamTypeClient,
		"cross.v1test", // protobuf package
		"CrossService", // protobuf service
		"Sum",          // protobuf method
		func(ctx context.Context, sender rerpc.Sender, receiver rerpc.Receiver) {
			typed := handlerstream.NewClient[v1test.SumRequest, v1test.SumResponse](sender, receiver)
			err := svc.Sum(ctx, typed)
			_ = receiver.Close()
			if err != nil {
				if _, ok := rerpc.AsError(err); !ok {
					if errors.Is(err, context.Canceled) {
						err = rerpc.Wrap(rerpc.CodeCanceled, err)
					}
					if errors.Is(err, context.DeadlineExceeded) {
						err = rerpc.Wrap(rerpc.CodeDeadlineExceeded, err)
					}
				}
			}
			_ = sender.Close(err)
		},
		opts...,
	)
	handlers = append(handlers, *sum)

	countUp := rerpc.NewStreamingHandler(
		rerpc.StreamTypeServer,
		"cross.v1test", // protobuf package
		"CrossService", // protobuf service
		"CountUp",      // protobuf method
		func(ctx context.Context, sender rerpc.Sender, receiver rerpc.Receiver) {
			typed := handlerstream.NewServer[v1test.CountUpResponse](sender)
			req, err := rerpc.ReceiveRequest[v1test.CountUpRequest](receiver)
			if err != nil {
				_ = receiver.Close()
				_ = sender.Close(err)
				return
			}
			if err = receiver.Close(); err != nil {
				_ = sender.Close(err)
				return
			}
			err = svc.CountUp(ctx, req, typed)
			if err != nil {
				if _, ok := rerpc.AsError(err); !ok {
					if errors.Is(err, context.Canceled) {
						err = rerpc.Wrap(rerpc.CodeCanceled, err)
					}
					if errors.Is(err, context.DeadlineExceeded) {
						err = rerpc.Wrap(rerpc.CodeDeadlineExceeded, err)
					}
				}
			}
			_ = sender.Close(err)
		},
		opts...,
	)
	handlers = append(handlers, *countUp)

	cumSum := rerpc.NewStreamingHandler(
		rerpc.StreamTypeBidirectional,
		"cross.v1test", // protobuf package
		"CrossService", // protobuf service
		"CumSum",       // protobuf method
		func(ctx context.Context, sender rerpc.Sender, receiver rerpc.Receiver) {
			typed := handlerstream.NewBidirectional[v1test.CumSumRequest, v1test.CumSumResponse](sender, receiver)
			err := svc.CumSum(ctx, typed)
			_ = receiver.Close()
			if err != nil {
				if _, ok := rerpc.AsError(err); !ok {
					if errors.Is(err, context.Canceled) {
						err = rerpc.Wrap(rerpc.CodeCanceled, err)
					}
					if errors.Is(err, context.DeadlineExceeded) {
						err = rerpc.Wrap(rerpc.CodeDeadlineExceeded, err)
					}
				}
			}
			_ = sender.Close(err)
		},
		opts...,
	)
	handlers = append(handlers, *cumSum)

	return handlers
}

type pluggableCrossServiceServer struct {
	ping    func(context.Context, *rerpc.Request[v1test.PingRequest]) (*rerpc.Response[v1test.PingResponse], error)
	fail    func(context.Context, *rerpc.Request[v1test.FailRequest]) (*rerpc.Response[v1test.FailResponse], error)
	sum     func(context.Context, *handlerstream.Client[v1test.SumRequest, v1test.SumResponse]) error
	countUp func(context.Context, *rerpc.Request[v1test.CountUpRequest], *handlerstream.Server[v1test.CountUpResponse]) error
	cumSum  func(context.Context, *handlerstream.Bidirectional[v1test.CumSumRequest, v1test.CumSumResponse]) error
}

func (s *pluggableCrossServiceServer) Ping(ctx context.Context, req *rerpc.Request[v1test.PingRequest]) (*rerpc.Response[v1test.PingResponse], error) {
	return s.ping(ctx, req)
}

func (s *pluggableCrossServiceServer) Fail(ctx context.Context, req *rerpc.Request[v1test.FailRequest]) (*rerpc.Response[v1test.FailResponse], error) {
	return s.fail(ctx, req)
}

func (s *pluggableCrossServiceServer) Sum(ctx context.Context, stream *handlerstream.Client[v1test.SumRequest, v1test.SumResponse]) error {
	return s.sum(ctx, stream)
}

func (s *pluggableCrossServiceServer) CountUp(ctx context.Context, req *rerpc.Request[v1test.CountUpRequest], stream *handlerstream.Server[v1test.CountUpResponse]) error {
	return s.countUp(ctx, req, stream)
}

func (s *pluggableCrossServiceServer) CumSum(ctx context.Context, stream *handlerstream.Bidirectional[v1test.CumSumRequest, v1test.CumSumResponse]) error {
	return s.cumSum(ctx, stream)
}

// NewCrossServiceHandler wraps each method on the service implementation in a
// rerpc.Handler. The returned slice can be passed to rerpc.NewServeMux.
//
// Unlike NewFullCrossServiceHandler, it allows the service to mix and match the
// signatures of FullCrossServiceServer and SimpleCrossServiceServer. For each
// method, it first tries to find a SimpleCrossServiceServer-style
// implementation. If a simple implementation isn't available, it falls back to
// the more complex FullCrossServiceServer-style implementation. If neither is
// available, it returns an error.
//
// Taken together, this approach lets implementations embed
// UnimplementedCrossServiceServer and implement each method using whichever
// signature is most convenient.
func NewCrossServiceHandler(svc any, opts ...rerpc.HandlerOption) ([]rerpc.Handler, error) {
	var impl pluggableCrossServiceServer

	// Find an implementation of Ping
	if pinger, ok := svc.(interface {
		Ping(context.Context, *v1test.PingRequest) (*v1test.PingResponse, error)
	}); ok {
		impl.ping = func(ctx context.Context, req *rerpc.Request[v1test.PingRequest]) (*rerpc.Response[v1test.PingResponse], error) {
			res, err := pinger.Ping(ctx, req.Msg)
			if err != nil {
				return nil, err
			}
			return rerpc.NewResponse(res), nil
		}
	} else if pinger, ok := svc.(interface {
		Ping(context.Context, *rerpc.Request[v1test.PingRequest]) (*rerpc.Response[v1test.PingResponse], error)
	}); ok {
		impl.ping = pinger.Ping
	} else {
		return nil, errors.New("no Ping implementation found")
	}

	// Find an implementation of Fail
	if failer, ok := svc.(interface {
		Fail(context.Context, *v1test.FailRequest) (*v1test.FailResponse, error)
	}); ok {
		impl.fail = func(ctx context.Context, req *rerpc.Request[v1test.FailRequest]) (*rerpc.Response[v1test.FailResponse], error) {
			res, err := failer.Fail(ctx, req.Msg)
			if err != nil {
				return nil, err
			}
			return rerpc.NewResponse(res), nil
		}
	} else if failer, ok := svc.(interface {
		Fail(context.Context, *rerpc.Request[v1test.FailRequest]) (*rerpc.Response[v1test.FailResponse], error)
	}); ok {
		impl.fail = failer.Fail
	} else {
		return nil, errors.New("no Fail implementation found")
	}

	// Find an implementation of Sum
	if sumer, ok := svc.(interface {
		Sum(context.Context, *handlerstream.Client[v1test.SumRequest, v1test.SumResponse]) error
	}); ok {
		impl.sum = sumer.Sum
	} else {
		return nil, errors.New("no Sum implementation found")
	}

	// Find an implementation of CountUp
	if countUper, ok := svc.(interface {
		CountUp(context.Context, *v1test.CountUpRequest, *handlerstream.Server[v1test.CountUpResponse]) error
	}); ok {
		impl.countUp = func(ctx context.Context, req *rerpc.Request[v1test.CountUpRequest], stream *handlerstream.Server[v1test.CountUpResponse]) error {
			return countUper.CountUp(ctx, req.Msg, stream)
		}
	} else if countUper, ok := svc.(interface {
		CountUp(context.Context, *rerpc.Request[v1test.CountUpRequest], *handlerstream.Server[v1test.CountUpResponse]) error
	}); ok {
		impl.countUp = countUper.CountUp
	} else {
		return nil, errors.New("no CountUp implementation found")
	}

	// Find an implementation of CumSum
	if cumSumer, ok := svc.(interface {
		CumSum(context.Context, *handlerstream.Bidirectional[v1test.CumSumRequest, v1test.CumSumResponse]) error
	}); ok {
		impl.cumSum = cumSumer.CumSum
	} else {
		return nil, errors.New("no CumSum implementation found")
	}

	return NewFullCrossServiceHandler(&impl, opts...), nil
}

var _ FullCrossServiceServer = (*UnimplementedCrossServiceServer)(nil) // verify interface implementation

// UnimplementedCrossServiceServer returns CodeUnimplemented from all methods.
type UnimplementedCrossServiceServer struct{}

func (UnimplementedCrossServiceServer) Ping(context.Context, *rerpc.Request[v1test.PingRequest]) (*rerpc.Response[v1test.PingResponse], error) {
	return nil, rerpc.Errorf(rerpc.CodeUnimplemented, "cross.v1test.CrossService.Ping isn't implemented")
}

func (UnimplementedCrossServiceServer) Fail(context.Context, *rerpc.Request[v1test.FailRequest]) (*rerpc.Response[v1test.FailResponse], error) {
	return nil, rerpc.Errorf(rerpc.CodeUnimplemented, "cross.v1test.CrossService.Fail isn't implemented")
}

func (UnimplementedCrossServiceServer) Sum(context.Context, *handlerstream.Client[v1test.SumRequest, v1test.SumResponse]) error {
	return rerpc.Errorf(rerpc.CodeUnimplemented, "cross.v1test.CrossService.Sum isn't implemented")
}

func (UnimplementedCrossServiceServer) CountUp(context.Context, *rerpc.Request[v1test.CountUpRequest], *handlerstream.Server[v1test.CountUpResponse]) error {
	return rerpc.Errorf(rerpc.CodeUnimplemented, "cross.v1test.CrossService.CountUp isn't implemented")
}

func (UnimplementedCrossServiceServer) CumSum(context.Context, *handlerstream.Bidirectional[v1test.CumSumRequest, v1test.CumSumResponse]) error {
	return rerpc.Errorf(rerpc.CodeUnimplemented, "cross.v1test.CrossService.CumSum isn't implemented")
}
