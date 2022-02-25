package connect

import (
	"context"
	"net/http"

	"github.com/bufbuild/connect/codec"
	"github.com/bufbuild/connect/codec/protobuf"
	"github.com/bufbuild/connect/compress"
)

type clientConfiguration struct {
	Protocol          protocol
	Procedure         string
	MaxResponseBytes  int64
	Interceptor       Interceptor
	Compressors       map[string]compress.Compressor
	Codec             codec.Codec
	CodecName         string
	RequestCompressor string
}

func newClientConfiguration(procedure string, options []ClientOption) (*clientConfiguration, *Error) {
	config := clientConfiguration{
		Protocol:    &protocolGRPC{web: false}, // default to HTTP/2 gRPC
		Procedure:   procedure,
		Compressors: make(map[string]compress.Compressor),
	}
	for _, opt := range options {
		opt.applyToClient(&config)
	}
	if err := config.Validate(); err != nil {
		return nil, err
	}
	return &config, nil
}

func (c *clientConfiguration) Validate() *Error {
	if c.Codec == nil || c.CodecName == "" {
		return errorf(CodeUnknown, "no codec configured")
	}
	if c.RequestCompressor != "" && c.RequestCompressor != compress.NameIdentity {
		if _, ok := c.Compressors[c.RequestCompressor]; !ok {
			return errorf(CodeUnknown, "no registered compressor for %q", c.RequestCompressor)
		}
	}
	if c.Protocol == nil {
		return errorf(CodeUnknown, "no protocol configured")
	}
	return nil
}

func (c *clientConfiguration) Protobuf() codec.Codec {
	if c.CodecName == protobuf.Name {
		return c.Codec
	}
	return protobuf.New()
}

func (c *clientConfiguration) newSpecification(t StreamType) Specification {
	return Specification{
		StreamType: t,
		Procedure:  c.Procedure,
		IsClient:   true,
	}
}

// A ClientOption configures a connect client.
//
// In addition to any options grouped in the documentation below, remember that
// Options are also valid ClientOptions.
type ClientOption interface {
	applyToClient(*clientConfiguration)
}

type requestCompressorOption struct {
	Name string
}

// WithRequestCompressor configures the client to use the specified algorithm
// to compress request messages. If the algorithm has not been registered using
// WithCompressor, the generated client constructor will return an error.
//
// Because some servers don't support compression, clients default to sending
// uncompressed requests.
func WithRequestCompressor(name string) ClientOption {
	return &requestCompressorOption{Name: name}
}

func (o *requestCompressorOption) applyToClient(config *clientConfiguration) {
	config.RequestCompressor = o.Name
}

type enableGRPCWebOption struct{}

// WithGRPCWeb switches clients to the gRPC-Web protocol. Clients generated by
// protoc-gen-go-connect default to using gRPC's HTTP/2 variant.
func WithGRPCWeb() ClientOption {
	return &enableGRPCWebOption{}
}

func (o *enableGRPCWebOption) applyToClient(config *clientConfiguration) {
	config.Protocol = &protocolGRPC{web: true}
}

// NewStreamClientImplementation is used by generated code - most users will
// never need to use it directly. It returns a stream constructor for a
// client-, server-, or bidirectional streaming remote procedure.
func NewStreamClientImplementation(
	doer Doer,
	baseURL, procedure string,
	stype StreamType,
	options ...ClientOption,
) (func(context.Context) (Sender, Receiver), error) {
	config, err := newClientConfiguration(procedure, options)
	if err != nil {
		return nil, err
	}
	protocolClient, protocolErr := config.Protocol.NewClient(&protocolClientParams{
		Spec:             config.newSpecification(stype),
		CompressorName:   config.RequestCompressor,
		Compressors:      newReadOnlyCompressors(config.Compressors),
		CodecName:        config.CodecName,
		Codec:            config.Codec,
		Protobuf:         config.Protobuf(),
		MaxResponseBytes: config.MaxResponseBytes,
		Doer:             doer,
		BaseURL:          baseURL,
	})
	if protocolErr != nil {
		return nil, NewError(CodeUnknown, protocolErr)
	}
	return func(ctx context.Context) (Sender, Receiver) {
		if ic := config.Interceptor; ic != nil {
			ctx = ic.WrapStreamContext(ctx)
		}
		header := make(http.Header, 8) // arbitrary power of two, prevent immediate resizing
		protocolClient.WriteRequestHeader(header)
		sender, receiver := protocolClient.NewStream(ctx, header)
		if ic := config.Interceptor; ic != nil {
			sender = ic.WrapStreamSender(ctx, sender)
			receiver = ic.WrapStreamReceiver(ctx, receiver)
		}
		return sender, receiver
	}, nil
}

// NewUnaryClientImplementation is used by generated code - most users will
// never need to use it directly. It returns a strongly-typed function to call
// a unary procedure.
func NewUnaryClientImplementation[Req, Res any](
	doer Doer,
	baseURL, procedure string,
	options ...ClientOption,
) (func(context.Context, *Message[Req]) (*Message[Res], error), error) {
	config, err := newClientConfiguration(procedure, options)
	if err != nil {
		return nil, err
	}
	spec := config.newSpecification(StreamTypeUnary)
	protocolClient, protocolErr := config.Protocol.NewClient(&protocolClientParams{
		Spec:             spec,
		CompressorName:   config.RequestCompressor,
		Compressors:      newReadOnlyCompressors(config.Compressors),
		CodecName:        config.CodecName,
		Codec:            config.Codec,
		Protobuf:         config.Protobuf(),
		MaxResponseBytes: config.MaxResponseBytes,
		Doer:             doer,
		BaseURL:          baseURL,
	})
	if protocolErr != nil {
		return nil, NewError(CodeUnknown, protocolErr)
	}
	send := Func(func(ctx context.Context, request AnyMessage) (AnyMessage, error) {
		sender, receiver := protocolClient.NewStream(ctx, request.Header())
		mergeHeaders(sender.Trailer(), request.Trailer())
		if err := sender.Send(request.Any()); err != nil {
			_ = sender.Close(err)
			_ = receiver.Close()
			return nil, err
		}
		if err := sender.Close(nil); err != nil {
			_ = receiver.Close()
			return nil, err
		}
		response, err := ReceiveUnaryMessage[Res](receiver)
		if err != nil {
			_ = receiver.Close()
			return nil, err
		}
		return response, receiver.Close()
	})
	if ic := config.Interceptor; ic != nil {
		send = ic.WrapUnary(send)
	}
	return func(ctx context.Context, request *Message[Req]) (*Message[Res], error) {
		// To make the specification and RPC headers visible to the full interceptor
		// chain (as though they were supplied by the caller), we'll add them here.
		request.spec = spec
		protocolClient.WriteRequestHeader(request.Header())
		response, err := send(ctx, request)
		if err != nil {
			return nil, err
		}
		typed, ok := response.(*Message[Res])
		if !ok {
			return nil, errorf(CodeInternal, "unexpected client response type %T", response)
		}
		return typed, nil
	}, nil
}
