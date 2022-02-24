// Code generated by protoc-gen-go-connect. DO NOT EDIT.
// versions:
// - protoc-gen-go-connect v0.0.1
// - protoc              v3.17.3
// source: cross/v1test/cross.proto

package crossv1test

import (
	context "context"
	errors "errors"
	connect "github.com/bufbuild/connect"
	clientstream "github.com/bufbuild/connect/clientstream"
	protobuf "github.com/bufbuild/connect/codec/protobuf"
	protojson "github.com/bufbuild/connect/codec/protojson"
	gzip "github.com/bufbuild/connect/compress/gzip"
	handlerstream "github.com/bufbuild/connect/handlerstream"
	v1test "github.com/bufbuild/connect/internal/crosstest/gen/proto/go/cross/v1test"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the
// connect package are compatible. If you get a compiler error that this
// constant isn't defined, this code was generated with a version of connect
// newer than the one compiled into your binary. You can fix the problem by
// either regenerating this code with an older version of connect or updating
// the connect version compiled into your binary.
const _ = connect.IsAtLeastVersion0_0_1

// CrossServiceClient is a client for the cross.v1test.CrossService service.
type CrossServiceClient interface {
	Ping(context.Context, *connect.Request[v1test.PingRequest]) (*connect.Response[v1test.PingResponse], error)
	Fail(context.Context, *connect.Request[v1test.FailRequest]) (*connect.Response[v1test.FailResponse], error)
	Sum(context.Context) *clientstream.Client[v1test.SumRequest, v1test.SumResponse]
	CountUp(context.Context, *connect.Request[v1test.CountUpRequest]) (*clientstream.Server[v1test.CountUpResponse], error)
	CumSum(context.Context) *clientstream.Bidirectional[v1test.CumSumRequest, v1test.CumSumResponse]
}

// NewCrossServiceClient constructs a client for the cross.v1test.CrossService
// service. By default, it uses the binary protobuf codec.
//
// The URL supplied here should be the base URL for the gRPC server (e.g.,
// https://api.acme.com or https://acme.com/grpc).
func NewCrossServiceClient(baseURL string, doer connect.Doer, opts ...connect.ClientOption) (CrossServiceClient, error) {
	baseURL = strings.TrimRight(baseURL, "/")
	opts = append([]connect.ClientOption{
		connect.WithCodec(protobuf.Name, protobuf.New()),
		connect.WithCompressor(gzip.Name, gzip.New()),
	}, opts...)
	var (
		client crossServiceClient
		err    error
	)
	client.ping, err = connect.NewUnaryClientImplementation[v1test.PingRequest, v1test.PingResponse](
		doer,
		baseURL,
		"cross.v1test.CrossService/Ping",
		opts...,
	)
	if err != nil {
		return nil, err
	}
	client.fail, err = connect.NewUnaryClientImplementation[v1test.FailRequest, v1test.FailResponse](
		doer,
		baseURL,
		"cross.v1test.CrossService/Fail",
		opts...,
	)
	if err != nil {
		return nil, err
	}
	client.sum, err = connect.NewStreamClientImplementation(
		doer,
		baseURL,
		"cross.v1test.CrossService/Sum",
		connect.StreamTypeClient,
		opts...,
	)
	if err != nil {
		return nil, err
	}
	client.countUp, err = connect.NewStreamClientImplementation(
		doer,
		baseURL,
		"cross.v1test.CrossService/CountUp",
		connect.StreamTypeServer,
		opts...,
	)
	if err != nil {
		return nil, err
	}
	client.cumSum, err = connect.NewStreamClientImplementation(
		doer,
		baseURL,
		"cross.v1test.CrossService/CumSum",
		connect.StreamTypeBidirectional,
		opts...,
	)
	if err != nil {
		return nil, err
	}
	return &client, nil
}

// crossServiceClient implements CrossServiceClient.
type crossServiceClient struct {
	ping    func(context.Context, *connect.Request[v1test.PingRequest]) (*connect.Response[v1test.PingResponse], error)
	fail    func(context.Context, *connect.Request[v1test.FailRequest]) (*connect.Response[v1test.FailResponse], error)
	sum     func(context.Context) (connect.Sender, connect.Receiver)
	countUp func(context.Context) (connect.Sender, connect.Receiver)
	cumSum  func(context.Context) (connect.Sender, connect.Receiver)
}

var _ CrossServiceClient = (*crossServiceClient)(nil) // verify interface implementation

// Ping calls cross.v1test.CrossService.Ping.
func (c *crossServiceClient) Ping(ctx context.Context, req *connect.Request[v1test.PingRequest]) (*connect.Response[v1test.PingResponse], error) {
	return c.ping(ctx, req)
}

// Fail calls cross.v1test.CrossService.Fail.
func (c *crossServiceClient) Fail(ctx context.Context, req *connect.Request[v1test.FailRequest]) (*connect.Response[v1test.FailResponse], error) {
	return c.fail(ctx, req)
}

// Sum calls cross.v1test.CrossService.Sum.
func (c *crossServiceClient) Sum(ctx context.Context) *clientstream.Client[v1test.SumRequest, v1test.SumResponse] {
	sender, receiver := c.sum(ctx)
	return clientstream.NewClient[v1test.SumRequest, v1test.SumResponse](sender, receiver)
}

// CountUp calls cross.v1test.CrossService.CountUp.
func (c *crossServiceClient) CountUp(ctx context.Context, req *connect.Request[v1test.CountUpRequest]) (*clientstream.Server[v1test.CountUpResponse], error) {
	sender, receiver := c.countUp(ctx)
	for key, values := range req.Header() {
		sender.Header()[key] = append(sender.Header()[key], values...)
	}
	for key, values := range req.Trailer() {
		sender.Trailer()[key] = append(sender.Trailer()[key], values...)
	}
	if err := sender.Send(req.Msg); err != nil {
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
func (c *crossServiceClient) CumSum(ctx context.Context) *clientstream.Bidirectional[v1test.CumSumRequest, v1test.CumSumResponse] {
	sender, receiver := c.cumSum(ctx)
	return clientstream.NewBidirectional[v1test.CumSumRequest, v1test.CumSumResponse](sender, receiver)
}

// CrossServiceHandler is an implementation of the cross.v1test.CrossService
// service.
type CrossServiceHandler interface {
	Ping(context.Context, *connect.Request[v1test.PingRequest]) (*connect.Response[v1test.PingResponse], error)
	Fail(context.Context, *connect.Request[v1test.FailRequest]) (*connect.Response[v1test.FailResponse], error)
	Sum(context.Context, *handlerstream.Client[v1test.SumRequest, v1test.SumResponse]) error
	CountUp(context.Context, *connect.Request[v1test.CountUpRequest], *handlerstream.Server[v1test.CountUpResponse]) error
	CumSum(context.Context, *handlerstream.Bidirectional[v1test.CumSumRequest, v1test.CumSumResponse]) error
}

// WithCrossServiceHandler wraps the service implementation in a
// connect.MuxOption, which can then be passed to connect.NewServeMux.
//
// By default, services support the gRPC and gRPC-Web protocols with the binary
// protobuf and JSON codecs.
func WithCrossServiceHandler(svc CrossServiceHandler, opts ...connect.HandlerOption) connect.MuxOption {
	handlers := make([]connect.Handler, 0, 5)
	opts = append([]connect.HandlerOption{
		connect.WithCodec(protobuf.Name, protobuf.New()),
		connect.WithCodec(protojson.Name, protojson.New()),
		connect.WithCompressor(gzip.Name, gzip.New()),
	}, opts...)

	ping, err := connect.NewUnaryHandler(
		"cross.v1test.CrossService/Ping", // procedure name
		"cross.v1test.CrossService",      // reflection name
		svc.Ping,
		opts...,
	)
	if err != nil {
		return connect.WithHandlers(nil, err)
	}
	handlers = append(handlers, *ping)

	fail, err := connect.NewUnaryHandler(
		"cross.v1test.CrossService/Fail", // procedure name
		"cross.v1test.CrossService",      // reflection name
		svc.Fail,
		opts...,
	)
	if err != nil {
		return connect.WithHandlers(nil, err)
	}
	handlers = append(handlers, *fail)

	sum, err := connect.NewStreamHandler(
		"cross.v1test.CrossService/Sum", // procedure name
		"cross.v1test.CrossService",     // reflection name
		connect.StreamTypeClient,
		func(ctx context.Context, sender connect.Sender, receiver connect.Receiver) {
			typed := handlerstream.NewClient[v1test.SumRequest, v1test.SumResponse](sender, receiver)
			err := svc.Sum(ctx, typed)
			_ = receiver.Close()
			_ = sender.Close(err)
		},
		opts...,
	)
	if err != nil {
		return connect.WithHandlers(nil, err)
	}
	handlers = append(handlers, *sum)

	countUp, err := connect.NewStreamHandler(
		"cross.v1test.CrossService/CountUp", // procedure name
		"cross.v1test.CrossService",         // reflection name
		connect.StreamTypeServer,
		func(ctx context.Context, sender connect.Sender, receiver connect.Receiver) {
			typed := handlerstream.NewServer[v1test.CountUpResponse](sender)
			req, err := connect.ReceiveRequest[v1test.CountUpRequest](receiver)
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
			_ = sender.Close(err)
		},
		opts...,
	)
	if err != nil {
		return connect.WithHandlers(nil, err)
	}
	handlers = append(handlers, *countUp)

	cumSum, err := connect.NewStreamHandler(
		"cross.v1test.CrossService/CumSum", // procedure name
		"cross.v1test.CrossService",        // reflection name
		connect.StreamTypeBidirectional,
		func(ctx context.Context, sender connect.Sender, receiver connect.Receiver) {
			typed := handlerstream.NewBidirectional[v1test.CumSumRequest, v1test.CumSumResponse](sender, receiver)
			err := svc.CumSum(ctx, typed)
			_ = receiver.Close()
			_ = sender.Close(err)
		},
		opts...,
	)
	if err != nil {
		return connect.WithHandlers(nil, err)
	}
	handlers = append(handlers, *cumSum)

	return connect.WithHandlers(handlers, nil)
}

// UnimplementedCrossServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedCrossServiceHandler struct{}

var _ CrossServiceHandler = (*UnimplementedCrossServiceHandler)(nil) // verify interface implementation

func (UnimplementedCrossServiceHandler) Ping(context.Context, *connect.Request[v1test.PingRequest]) (*connect.Response[v1test.PingResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cross.v1test.CrossService.Ping isn't implemented"))
}

func (UnimplementedCrossServiceHandler) Fail(context.Context, *connect.Request[v1test.FailRequest]) (*connect.Response[v1test.FailResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cross.v1test.CrossService.Fail isn't implemented"))
}

func (UnimplementedCrossServiceHandler) Sum(context.Context, *handlerstream.Client[v1test.SumRequest, v1test.SumResponse]) error {
	return connect.NewError(connect.CodeUnimplemented, errors.New("cross.v1test.CrossService.Sum isn't implemented"))
}

func (UnimplementedCrossServiceHandler) CountUp(context.Context, *connect.Request[v1test.CountUpRequest], *handlerstream.Server[v1test.CountUpResponse]) error {
	return connect.NewError(connect.CodeUnimplemented, errors.New("cross.v1test.CrossService.CountUp isn't implemented"))
}

func (UnimplementedCrossServiceHandler) CumSum(context.Context, *handlerstream.Bidirectional[v1test.CumSumRequest, v1test.CumSumResponse]) error {
	return connect.NewError(connect.CodeUnimplemented, errors.New("cross.v1test.CrossService.CumSum isn't implemented"))
}
