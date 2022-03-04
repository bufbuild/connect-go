// Copyright 2021-2022 Buf Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package connect

import (
	"context"
	"net/http"
)

type clientConfiguration struct {
	Protocol               protocol
	Procedure              string
	MaxResponseBytes       int64
	CompressMinBytes       int
	Interceptor            Interceptor
	CompressionPools       map[string]compressionPool
	Codec                  Codec
	RequestCompressionName string
}

func newClientConfiguration(procedure string, options []ClientOption) (*clientConfiguration, *Error) {
	config := clientConfiguration{
		Protocol:         &protocolGRPC{web: false}, // default to HTTP/2 gRPC
		Procedure:        procedure,
		CompressionPools: make(map[string]compressionPool),
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
	if c.Codec == nil || c.Codec.Name() == "" {
		return errorf(CodeUnknown, "no codec configured")
	}
	if c.RequestCompressionName != "" && c.RequestCompressionName != compressionIdentity {
		if _, ok := c.CompressionPools[c.RequestCompressionName]; !ok {
			return errorf(CodeUnknown, "unknown compression %q", c.RequestCompressionName)
		}
	}
	if c.Protocol == nil {
		return errorf(CodeUnknown, "no protocol configured")
	}
	return nil
}

func (c *clientConfiguration) Protobuf() Codec {
	if c.Codec.Name() == codecNameProto {
		return c.Codec
	}
	return &protoBinaryCodec{}
}

func (c *clientConfiguration) newSpecification(t StreamType) Specification {
	return Specification{
		StreamType: t,
		Procedure:  c.Procedure,
		IsClient:   true,
	}
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
		CompressionName:  config.RequestCompressionName,
		CompressionPools: newReadOnlyCompressionPools(config.CompressionPools),
		Codec:            config.Codec,
		Protobuf:         config.Protobuf(),
		MaxResponseBytes: config.MaxResponseBytes,
		CompressMinBytes: config.CompressMinBytes,
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
) (func(context.Context, *Envelope[Req]) (*Envelope[Res], error), error) {
	config, err := newClientConfiguration(procedure, options)
	if err != nil {
		return nil, err
	}
	spec := config.newSpecification(StreamTypeUnary)
	protocolClient, protocolErr := config.Protocol.NewClient(&protocolClientParams{
		Spec:             spec,
		CompressionName:  config.RequestCompressionName,
		CompressionPools: newReadOnlyCompressionPools(config.CompressionPools),
		Codec:            config.Codec,
		Protobuf:         config.Protobuf(),
		MaxResponseBytes: config.MaxResponseBytes,
		CompressMinBytes: config.CompressMinBytes,
		Doer:             doer,
		BaseURL:          baseURL,
	})
	if protocolErr != nil {
		return nil, NewError(CodeUnknown, protocolErr)
	}
	send := UnaryFunc(func(ctx context.Context, request AnyEnvelope) (AnyEnvelope, error) {
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
		response, err := ReceiveUnaryEnvelope[Res](receiver)
		if err != nil {
			_ = receiver.Close()
			return nil, err
		}
		return response, receiver.Close()
	})
	if ic := config.Interceptor; ic != nil {
		send = ic.WrapUnary(send)
	}
	return func(ctx context.Context, request *Envelope[Req]) (*Envelope[Res], error) {
		// To make the specification and RPC headers visible to the full interceptor
		// chain (as though they were supplied by the caller), we'll add them here.
		request.spec = spec
		protocolClient.WriteRequestHeader(request.Header())
		response, err := send(ctx, request)
		if err != nil {
			return nil, err
		}
		typed, ok := response.(*Envelope[Res])
		if !ok {
			return nil, errorf(CodeInternal, "unexpected client response type %T", response)
		}
		return typed, nil
	}, nil
}
