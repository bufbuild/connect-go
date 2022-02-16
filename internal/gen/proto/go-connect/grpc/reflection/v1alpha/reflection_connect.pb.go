// Code generated by protoc-gen-go-connect. DO NOT EDIT.
// versions:
// - protoc-gen-go-connect v0.0.1
// - protoc              v3.17.3
// source: grpc/reflection/v1alpha/reflection.proto

package reflectionv1alpha1

import (
	context "context"
	connect "github.com/bufconnect/connect"
	clientstream "github.com/bufconnect/connect/clientstream"
	protobuf "github.com/bufconnect/connect/codec/protobuf"
	gzip "github.com/bufconnect/connect/compress/gzip"
	handlerstream "github.com/bufconnect/connect/handlerstream"
	v1alpha "github.com/bufconnect/connect/internal/gen/proto/go/grpc/reflection/v1alpha"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the
// connect package are compatible. If you get a compiler error that this
// constant isn't defined, this code was generated with a version of connect
// newer than the one compiled into your binary. You can fix the problem by
// either regenerating this code with an older version of connect or updating
// the connect version compiled into your binary.
const _ = connect.IsAtLeastVersion0_0_1

// ServerReflectionClient is a client for the
// internal.reflection.v1alpha1.ServerReflection service.
type ServerReflectionClient interface {
	// The reflection service is structured as a bidirectional stream, ensuring
	// all related requests go to a single server.
	ServerReflectionInfo(context.Context) *clientstream.Bidirectional[v1alpha.ServerReflectionRequest, v1alpha.ServerReflectionResponse]
}

// NewServerReflectionClient constructs a client for the
// internal.reflection.v1alpha1.ServerReflection service. By default, it uses
// the binary protobuf codec.
//
// The URL supplied here should be the base URL for the gRPC server (e.g.,
// https://api.acme.com or https://acme.com/grpc).
func NewServerReflectionClient(baseURL string, doer connect.Doer, opts ...connect.ClientOption) (ServerReflectionClient, error) {
	baseURL = strings.TrimRight(baseURL, "/")
	opts = append([]connect.ClientOption{
		connect.WithGRPC(true),
		connect.WithCodec(protobuf.NameBinary, protobuf.NewBinary()),
		connect.WithCompressor(gzip.Name, gzip.New()),
	}, opts...)
	var (
		client serverReflectionClient
		err    error
	)
	client.serverReflectionInfo, err = connect.NewClientStream(
		doer,
		connect.StreamTypeBidirectional,
		baseURL,
		"internal.reflection.v1alpha1.ServerReflection/ServerReflectionInfo",
		opts...,
	)
	if err != nil {
		return nil, err
	}
	return &client, nil
}

// serverReflectionClient implements ServerReflectionClient.
type serverReflectionClient struct {
	serverReflectionInfo func(context.Context) (connect.Sender, connect.Receiver)
}

var _ ServerReflectionClient = (*serverReflectionClient)(nil) // verify interface implementation

// ServerReflectionInfo calls
// internal.reflection.v1alpha1.ServerReflection.ServerReflectionInfo.
func (c *serverReflectionClient) ServerReflectionInfo(ctx context.Context) *clientstream.Bidirectional[v1alpha.ServerReflectionRequest, v1alpha.ServerReflectionResponse] {
	sender, receiver := c.serverReflectionInfo(ctx)
	return clientstream.NewBidirectional[v1alpha.ServerReflectionRequest, v1alpha.ServerReflectionResponse](sender, receiver)
}

// ServerReflectionHandler is an implementation of the
// internal.reflection.v1alpha1.ServerReflection service.
type ServerReflectionHandler interface {
	// The reflection service is structured as a bidirectional stream, ensuring
	// all related requests go to a single server.
	ServerReflectionInfo(context.Context, *handlerstream.Bidirectional[v1alpha.ServerReflectionRequest, v1alpha.ServerReflectionResponse]) error
}

// WithServerReflectionHandler wraps the service implementation in a
// connect.MuxOption, which can then be passed to connect.NewServeMux.
//
// By default, services support the gRPC and gRPC-Web protocols with the binary
// protobuf and JSON codecs.
func WithServerReflectionHandler(svc ServerReflectionHandler, opts ...connect.HandlerOption) connect.MuxOption {
	handlers := make([]connect.Handler, 0, 1)
	opts = append([]connect.HandlerOption{
		connect.WithGRPC(true),
		connect.WithGRPCWeb(true),
		connect.WithCodec(protobuf.NameBinary, protobuf.NewBinary()),
		connect.WithCodec(protobuf.NameJSON, protobuf.NewJSON()),
		connect.WithCompressor(gzip.Name, gzip.New()),
	}, opts...)

	serverReflectionInfo, err := connect.NewStreamingHandler(
		connect.StreamTypeBidirectional,
		"internal.reflection.v1alpha1.ServerReflection/ServerReflectionInfo", // procedure name
		"internal.reflection.v1alpha1.ServerReflection",                      // reflection name
		func(ctx context.Context, sender connect.Sender, receiver connect.Receiver) {
			typed := handlerstream.NewBidirectional[v1alpha.ServerReflectionRequest, v1alpha.ServerReflectionResponse](sender, receiver)
			err := svc.ServerReflectionInfo(ctx, typed)
			_ = receiver.Close()
			_ = sender.Close(err)
		},
		opts...,
	)
	if err != nil {
		return connect.WithHandlers(nil, err)
	}
	handlers = append(handlers, *serverReflectionInfo)

	return connect.WithHandlers(handlers, nil)
}

// UnimplementedServerReflectionHandler returns CodeUnimplemented from all
// methods.
type UnimplementedServerReflectionHandler struct{}

var _ ServerReflectionHandler = (*UnimplementedServerReflectionHandler)(nil) // verify interface implementation

func (UnimplementedServerReflectionHandler) ServerReflectionInfo(context.Context, *handlerstream.Bidirectional[v1alpha.ServerReflectionRequest, v1alpha.ServerReflectionResponse]) error {
	return connect.Errorf(connect.CodeUnimplemented, "internal.reflection.v1alpha1.ServerReflection.ServerReflectionInfo isn't implemented")
}
