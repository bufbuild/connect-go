// Code generated by protoc-gen-go-rerpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-rerpc v0.0.1
// - protoc             v3.17.3
// source: internal/crosstest/v1test/cross.proto

package crosspb

import (
	context "context"
	errors "errors"
	rerpc "github.com/rerpc/rerpc"
	proto "google.golang.org/protobuf/proto"
	http "net/http"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the
// rerpc package are compatible. If you get a compiler error that this constant
// isn't defined, this code was generated with a version of rerpc newer than the
// one compiled into your binary. You can fix the problem by either regenerating
// this code with an older version of rerpc or updating the rerpc version
// compiled into your binary.
const _ = rerpc.SupportsCodeGenV0 // requires reRPC v0.0.1 or later

// CrossServiceClientReRPC is a client for the
// internal.crosstest.v1test.CrossService service.
type CrossServiceClientReRPC interface {
	Ping(ctx context.Context, req *PingRequest, opts ...rerpc.CallOption) (*PingResponse, error)
	Fail(ctx context.Context, req *FailRequest, opts ...rerpc.CallOption) (*FailResponse, error)
	Sum(ctx context.Context, opts ...rerpc.CallOption) *CrossServiceClientReRPC_Sum
	CountUp(ctx context.Context, req *CountUpRequest, opts ...rerpc.CallOption) (*CrossServiceClientReRPC_CountUp, error)
	CumSum(ctx context.Context, opts ...rerpc.CallOption) *CrossServiceClientReRPC_CumSum
}

type crossServiceClientReRPC struct {
	doer    rerpc.Doer
	baseURL string
	options []rerpc.CallOption
}

// NewCrossServiceClientReRPC constructs a client for the
// internal.crosstest.v1test.CrossService service. Call options passed here
// apply to all calls made with this client.
//
// The URL supplied here should be the base URL for the gRPC server (e.g.,
// https://api.acme.com or https://acme.com/grpc).
func NewCrossServiceClientReRPC(baseURL string, doer rerpc.Doer, opts ...rerpc.CallOption) CrossServiceClientReRPC {
	return &crossServiceClientReRPC{
		baseURL: strings.TrimRight(baseURL, "/"),
		doer:    doer,
		options: opts,
	}
}

func (c *crossServiceClientReRPC) mergeOptions(opts []rerpc.CallOption) []rerpc.CallOption {
	merged := make([]rerpc.CallOption, 0, len(c.options)+len(opts))
	for _, o := range c.options {
		merged = append(merged, o)
	}
	for _, o := range opts {
		merged = append(merged, o)
	}
	return merged
}

// Ping calls internal.crosstest.v1test.CrossService.Ping. Call options passed
// here apply only to this call.
func (c *crossServiceClientReRPC) Ping(ctx context.Context, req *PingRequest, opts ...rerpc.CallOption) (*PingResponse, error) {
	merged := c.mergeOptions(opts)
	ic := rerpc.ConfiguredCallInterceptor(merged...)
	ctx, call := rerpc.NewCall(
		ctx,
		c.doer,
		rerpc.StreamTypeUnary,
		c.baseURL,
		"internal.crosstest.v1test", // protobuf package
		"CrossService",              // protobuf service
		"Ping",                      // protobuf method
		merged...,
	)
	wrapped := rerpc.Func(func(ctx context.Context, msg proto.Message) (proto.Message, error) {
		stream := call(ctx)
		if err := stream.Send(req); err != nil {
			_ = stream.CloseSend(err)
			_ = stream.CloseReceive()
			return nil, err
		}
		if err := stream.CloseSend(nil); err != nil {
			_ = stream.CloseReceive()
			return nil, err
		}
		var res PingResponse
		if err := stream.Receive(&res); err != nil {
			_ = stream.CloseReceive()
			return nil, err
		}
		return &res, stream.CloseReceive()
	})
	if ic != nil {
		wrapped = ic.Wrap(wrapped)
	}
	res, err := wrapped(ctx, req)
	if err != nil {
		return nil, err
	}
	typed, ok := res.(*PingResponse)
	if !ok {
		return nil, rerpc.Errorf(rerpc.CodeInternal, "expected response to be internal.crosstest.v1test.PingResponse, got %v", res.ProtoReflect().Descriptor().FullName())
	}
	return typed, nil
}

// Fail calls internal.crosstest.v1test.CrossService.Fail. Call options passed
// here apply only to this call.
func (c *crossServiceClientReRPC) Fail(ctx context.Context, req *FailRequest, opts ...rerpc.CallOption) (*FailResponse, error) {
	merged := c.mergeOptions(opts)
	ic := rerpc.ConfiguredCallInterceptor(merged...)
	ctx, call := rerpc.NewCall(
		ctx,
		c.doer,
		rerpc.StreamTypeUnary,
		c.baseURL,
		"internal.crosstest.v1test", // protobuf package
		"CrossService",              // protobuf service
		"Fail",                      // protobuf method
		merged...,
	)
	wrapped := rerpc.Func(func(ctx context.Context, msg proto.Message) (proto.Message, error) {
		stream := call(ctx)
		if err := stream.Send(req); err != nil {
			_ = stream.CloseSend(err)
			_ = stream.CloseReceive()
			return nil, err
		}
		if err := stream.CloseSend(nil); err != nil {
			_ = stream.CloseReceive()
			return nil, err
		}
		var res FailResponse
		if err := stream.Receive(&res); err != nil {
			_ = stream.CloseReceive()
			return nil, err
		}
		return &res, stream.CloseReceive()
	})
	if ic != nil {
		wrapped = ic.Wrap(wrapped)
	}
	res, err := wrapped(ctx, req)
	if err != nil {
		return nil, err
	}
	typed, ok := res.(*FailResponse)
	if !ok {
		return nil, rerpc.Errorf(rerpc.CodeInternal, "expected response to be internal.crosstest.v1test.FailResponse, got %v", res.ProtoReflect().Descriptor().FullName())
	}
	return typed, nil
}

// Sum calls internal.crosstest.v1test.CrossService.Sum. Call options passed
// here apply only to this call.
func (c *crossServiceClientReRPC) Sum(ctx context.Context, opts ...rerpc.CallOption) *CrossServiceClientReRPC_Sum {
	merged := c.mergeOptions(opts)
	ic := rerpc.ConfiguredCallInterceptor(merged...)
	ctx, call := rerpc.NewCall(
		ctx,
		c.doer,
		rerpc.StreamTypeClient,
		c.baseURL,
		"internal.crosstest.v1test", // protobuf package
		"CrossService",              // protobuf service
		"Sum",                       // protobuf method
		merged...,
	)
	if ic != nil {
		call = ic.WrapStream(call)
	}
	stream := call(ctx)
	return NewCrossServiceClientReRPC_Sum(stream)
}

// CountUp calls internal.crosstest.v1test.CrossService.CountUp. Call options
// passed here apply only to this call.
func (c *crossServiceClientReRPC) CountUp(ctx context.Context, req *CountUpRequest, opts ...rerpc.CallOption) (*CrossServiceClientReRPC_CountUp, error) {
	merged := c.mergeOptions(opts)
	ic := rerpc.ConfiguredCallInterceptor(merged...)
	ctx, call := rerpc.NewCall(
		ctx,
		c.doer,
		rerpc.StreamTypeServer,
		c.baseURL,
		"internal.crosstest.v1test", // protobuf package
		"CrossService",              // protobuf service
		"CountUp",                   // protobuf method
		merged...,
	)
	if ic != nil {
		call = ic.WrapStream(call)
	}
	stream := call(ctx)
	if err := stream.Send(req); err != nil {
		_ = stream.CloseSend(err)
		_ = stream.CloseReceive()
		return nil, err
	}
	if err := stream.CloseSend(nil); err != nil {
		_ = stream.CloseReceive()
		return nil, err
	}
	return NewCrossServiceClientReRPC_CountUp(stream), nil
}

// CumSum calls internal.crosstest.v1test.CrossService.CumSum. Call options
// passed here apply only to this call.
func (c *crossServiceClientReRPC) CumSum(ctx context.Context, opts ...rerpc.CallOption) *CrossServiceClientReRPC_CumSum {
	merged := c.mergeOptions(opts)
	ic := rerpc.ConfiguredCallInterceptor(merged...)
	ctx, call := rerpc.NewCall(
		ctx,
		c.doer,
		rerpc.StreamTypeBidirectional,
		c.baseURL,
		"internal.crosstest.v1test", // protobuf package
		"CrossService",              // protobuf service
		"CumSum",                    // protobuf method
		merged...,
	)
	if ic != nil {
		call = ic.WrapStream(call)
	}
	stream := call(ctx)
	return NewCrossServiceClientReRPC_CumSum(stream)
}

// CrossServiceReRPC is a server for the internal.crosstest.v1test.CrossService
// service. To make sure that adding methods to this protobuf service doesn't
// break all implementations of this interface, all implementations must embed
// UnimplementedCrossServiceReRPC.
//
// By default, recent versions of grpc-go have a similar forward compatibility
// requirement. See https://github.com/grpc/grpc-go/issues/3794 for a longer
// discussion.
type CrossServiceReRPC interface {
	Ping(context.Context, *PingRequest) (*PingResponse, error)
	Fail(context.Context, *FailRequest) (*FailResponse, error)
	Sum(context.Context, *CrossServiceReRPC_Sum) error
	CountUp(context.Context, *CountUpRequest, *CrossServiceReRPC_CountUp) error
	CumSum(context.Context, *CrossServiceReRPC_CumSum) error
	mustEmbedUnimplementedCrossServiceReRPC()
}

// NewCrossServiceHandlerReRPC wraps the service implementation in an HTTP
// handler. It returns the handler and the path on which to mount it.
func NewCrossServiceHandlerReRPC(svc CrossServiceReRPC, opts ...rerpc.HandlerOption) (string, *http.ServeMux) {
	mux := http.NewServeMux()
	ic := rerpc.ConfiguredHandlerInterceptor(opts...)

	pingFunc := rerpc.Func(func(ctx context.Context, req proto.Message) (proto.Message, error) {
		typed, ok := req.(*PingRequest)
		if !ok {
			return nil, rerpc.Errorf(
				rerpc.CodeInternal,
				"can't call internal.crosstest.v1test.CrossService.Ping with a %v",
				req.ProtoReflect().Descriptor().FullName(),
			)
		}
		return svc.Ping(ctx, typed)
	})
	if ic != nil {
		pingFunc = ic.Wrap(pingFunc)
	}
	ping := rerpc.NewHandler(
		rerpc.StreamTypeUnary,
		"internal.crosstest.v1test", // protobuf package
		"CrossService",              // protobuf service
		"Ping",                      // protobuf method
		func(ctx context.Context, sf rerpc.StreamFunc) {
			stream := sf(ctx)
			defer stream.CloseReceive()
			if err := ctx.Err(); err != nil {
				if errors.Is(err, context.Canceled) {
					_ = stream.CloseSend(rerpc.Wrap(rerpc.CodeCanceled, err))
					return
				}
				if errors.Is(err, context.DeadlineExceeded) {
					_ = stream.CloseSend(rerpc.Wrap(rerpc.CodeDeadlineExceeded, err))
					return
				}
				_ = stream.CloseSend(err) // unreachable per context docs
			}
			var req PingRequest
			if err := stream.Receive(&req); err != nil {
				_ = stream.CloseSend(err)
				return
			}
			res, err := pingFunc(ctx, &req)
			if err != nil {
				if _, ok := rerpc.AsError(err); !ok {
					if errors.Is(err, context.Canceled) {
						err = rerpc.Wrap(rerpc.CodeCanceled, err)
					}
					if errors.Is(err, context.DeadlineExceeded) {
						err = rerpc.Wrap(rerpc.CodeDeadlineExceeded, err)
					}
				}
				_ = stream.CloseSend(err)
				return
			}
			_ = stream.CloseSend(stream.Send(res))
		},
		opts...,
	)
	mux.Handle(ping.Path(), ping)

	failFunc := rerpc.Func(func(ctx context.Context, req proto.Message) (proto.Message, error) {
		typed, ok := req.(*FailRequest)
		if !ok {
			return nil, rerpc.Errorf(
				rerpc.CodeInternal,
				"can't call internal.crosstest.v1test.CrossService.Fail with a %v",
				req.ProtoReflect().Descriptor().FullName(),
			)
		}
		return svc.Fail(ctx, typed)
	})
	if ic != nil {
		failFunc = ic.Wrap(failFunc)
	}
	fail := rerpc.NewHandler(
		rerpc.StreamTypeUnary,
		"internal.crosstest.v1test", // protobuf package
		"CrossService",              // protobuf service
		"Fail",                      // protobuf method
		func(ctx context.Context, sf rerpc.StreamFunc) {
			stream := sf(ctx)
			defer stream.CloseReceive()
			if err := ctx.Err(); err != nil {
				if errors.Is(err, context.Canceled) {
					_ = stream.CloseSend(rerpc.Wrap(rerpc.CodeCanceled, err))
					return
				}
				if errors.Is(err, context.DeadlineExceeded) {
					_ = stream.CloseSend(rerpc.Wrap(rerpc.CodeDeadlineExceeded, err))
					return
				}
				_ = stream.CloseSend(err) // unreachable per context docs
			}
			var req FailRequest
			if err := stream.Receive(&req); err != nil {
				_ = stream.CloseSend(err)
				return
			}
			res, err := failFunc(ctx, &req)
			if err != nil {
				if _, ok := rerpc.AsError(err); !ok {
					if errors.Is(err, context.Canceled) {
						err = rerpc.Wrap(rerpc.CodeCanceled, err)
					}
					if errors.Is(err, context.DeadlineExceeded) {
						err = rerpc.Wrap(rerpc.CodeDeadlineExceeded, err)
					}
				}
				_ = stream.CloseSend(err)
				return
			}
			_ = stream.CloseSend(stream.Send(res))
		},
		opts...,
	)
	mux.Handle(fail.Path(), fail)

	sum := rerpc.NewHandler(
		rerpc.StreamTypeClient,
		"internal.crosstest.v1test", // protobuf package
		"CrossService",              // protobuf service
		"Sum",                       // protobuf method
		func(ctx context.Context, sf rerpc.StreamFunc) {
			if ic != nil {
				sf = ic.WrapStream(sf)
			}
			stream := sf(ctx)
			typed := NewCrossServiceReRPC_Sum(stream)
			err := svc.Sum(stream.Context(), typed)
			_ = stream.CloseReceive()
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
			_ = stream.CloseSend(err)
		},
		opts...,
	)
	mux.Handle(sum.Path(), sum)

	countUp := rerpc.NewHandler(
		rerpc.StreamTypeServer,
		"internal.crosstest.v1test", // protobuf package
		"CrossService",              // protobuf service
		"CountUp",                   // protobuf method
		func(ctx context.Context, sf rerpc.StreamFunc) {
			if ic != nil {
				sf = ic.WrapStream(sf)
			}
			stream := sf(ctx)
			typed := NewCrossServiceReRPC_CountUp(stream)
			var req CountUpRequest
			if err := stream.Receive(&req); err != nil {
				_ = stream.CloseReceive()
				_ = stream.CloseSend(err)
				return
			}
			if err := stream.CloseReceive(); err != nil {
				_ = stream.CloseSend(err)
				return
			}
			err := svc.CountUp(stream.Context(), &req, typed)
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
			_ = stream.CloseSend(err)
		},
		opts...,
	)
	mux.Handle(countUp.Path(), countUp)

	cumSum := rerpc.NewHandler(
		rerpc.StreamTypeBidirectional,
		"internal.crosstest.v1test", // protobuf package
		"CrossService",              // protobuf service
		"CumSum",                    // protobuf method
		func(ctx context.Context, sf rerpc.StreamFunc) {
			if ic != nil {
				sf = ic.WrapStream(sf)
			}
			stream := sf(ctx)
			typed := NewCrossServiceReRPC_CumSum(stream)
			err := svc.CumSum(stream.Context(), typed)
			_ = stream.CloseReceive()
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
			_ = stream.CloseSend(err)
		},
		opts...,
	)
	mux.Handle(cumSum.Path(), cumSum)

	// Respond to unknown protobuf methods with gRPC and Twirp's 404 equivalents.
	mux.Handle("/", rerpc.NewBadRouteHandler(opts...))

	return cumSum.ServicePath(), mux
}

var _ CrossServiceReRPC = (*UnimplementedCrossServiceReRPC)(nil) // verify interface implementation

// UnimplementedCrossServiceReRPC returns CodeUnimplemented from all methods. To
// maintain forward compatibility, all implementations of CrossServiceReRPC must
// embed UnimplementedCrossServiceReRPC.
type UnimplementedCrossServiceReRPC struct{}

func (UnimplementedCrossServiceReRPC) Ping(context.Context, *PingRequest) (*PingResponse, error) {
	return nil, rerpc.Errorf(rerpc.CodeUnimplemented, "internal.crosstest.v1test.CrossService.Ping isn't implemented")
}

func (UnimplementedCrossServiceReRPC) Fail(context.Context, *FailRequest) (*FailResponse, error) {
	return nil, rerpc.Errorf(rerpc.CodeUnimplemented, "internal.crosstest.v1test.CrossService.Fail isn't implemented")
}

func (UnimplementedCrossServiceReRPC) Sum(context.Context, *CrossServiceReRPC_Sum) error {
	return rerpc.Errorf(rerpc.CodeUnimplemented, "internal.crosstest.v1test.CrossService.Sum isn't implemented")
}

func (UnimplementedCrossServiceReRPC) CountUp(context.Context, *CountUpRequest, *CrossServiceReRPC_CountUp) error {
	return rerpc.Errorf(rerpc.CodeUnimplemented, "internal.crosstest.v1test.CrossService.CountUp isn't implemented")
}

func (UnimplementedCrossServiceReRPC) CumSum(context.Context, *CrossServiceReRPC_CumSum) error {
	return rerpc.Errorf(rerpc.CodeUnimplemented, "internal.crosstest.v1test.CrossService.CumSum isn't implemented")
}

func (UnimplementedCrossServiceReRPC) mustEmbedUnimplementedCrossServiceReRPC() {}

// CrossServiceClientReRPC_Sum is the client-side stream for the
// internal.crosstest.v1test.CrossService.Sum procedure.
type CrossServiceClientReRPC_Sum struct {
	stream rerpc.Stream
}

func NewCrossServiceClientReRPC_Sum(stream rerpc.Stream) *CrossServiceClientReRPC_Sum {
	return &CrossServiceClientReRPC_Sum{stream}
}

func (s *CrossServiceClientReRPC_Sum) Send(msg *SumRequest) error {
	return s.stream.Send(msg)
}

func (s *CrossServiceClientReRPC_Sum) CloseAndReceive() (*SumResponse, error) {
	if err := s.stream.CloseSend(nil); err != nil {
		return nil, err
	}
	var res SumResponse
	err := s.stream.Receive(&res)
	return &res, err
}

// CrossServiceClientReRPC_CountUp is the client-side stream for the
// internal.crosstest.v1test.CrossService.CountUp procedure.
type CrossServiceClientReRPC_CountUp struct {
	stream rerpc.Stream
}

func NewCrossServiceClientReRPC_CountUp(stream rerpc.Stream) *CrossServiceClientReRPC_CountUp {
	return &CrossServiceClientReRPC_CountUp{stream}
}

func (s *CrossServiceClientReRPC_CountUp) Receive() (*CountUpResponse, error) {
	var req CountUpResponse
	if err := s.stream.Receive(&req); err != nil {
		return nil, err
	}
	return &req, nil
}

func (s *CrossServiceClientReRPC_CountUp) Close() error {
	return s.stream.CloseReceive()
}

// CrossServiceClientReRPC_CumSum is the client-side stream for the
// internal.crosstest.v1test.CrossService.CumSum procedure.
type CrossServiceClientReRPC_CumSum struct {
	stream rerpc.Stream
}

func NewCrossServiceClientReRPC_CumSum(stream rerpc.Stream) *CrossServiceClientReRPC_CumSum {
	return &CrossServiceClientReRPC_CumSum{stream}
}

func (s *CrossServiceClientReRPC_CumSum) Send(msg *CumSumRequest) error {
	return s.stream.Send(msg)
}

func (s *CrossServiceClientReRPC_CumSum) CloseSend() error {
	return s.stream.CloseSend(nil)
}

func (s *CrossServiceClientReRPC_CumSum) Receive() (*CumSumResponse, error) {
	var req CumSumResponse
	if err := s.stream.Receive(&req); err != nil {
		return nil, err
	}
	return &req, nil
}

func (s *CrossServiceClientReRPC_CumSum) CloseReceive() error {
	return s.stream.CloseReceive()
}

// CrossServiceReRPC_Sum is the server-side stream for the
// internal.crosstest.v1test.CrossService.Sum procedure.
type CrossServiceReRPC_Sum struct {
	stream rerpc.Stream
}

func NewCrossServiceReRPC_Sum(stream rerpc.Stream) *CrossServiceReRPC_Sum {
	return &CrossServiceReRPC_Sum{stream}
}

func (s *CrossServiceReRPC_Sum) Receive() (*SumRequest, error) {
	var req SumRequest
	if err := s.stream.Receive(&req); err != nil {
		return nil, err
	}
	return &req, nil
}

func (s *CrossServiceReRPC_Sum) SendAndClose(msg *SumResponse) error {
	if err := s.stream.CloseReceive(); err != nil {
		return err
	}
	return s.stream.Send(msg)
}

// CrossServiceReRPC_CountUp is the server-side stream for the
// internal.crosstest.v1test.CrossService.CountUp procedure.
type CrossServiceReRPC_CountUp struct {
	stream rerpc.Stream
}

func NewCrossServiceReRPC_CountUp(stream rerpc.Stream) *CrossServiceReRPC_CountUp {
	return &CrossServiceReRPC_CountUp{stream}
}

func (s *CrossServiceReRPC_CountUp) Send(msg *CountUpResponse) error {
	return s.stream.Send(msg)
}

// CrossServiceReRPC_CumSum is the server-side stream for the
// internal.crosstest.v1test.CrossService.CumSum procedure.
type CrossServiceReRPC_CumSum struct {
	stream rerpc.Stream
}

func NewCrossServiceReRPC_CumSum(stream rerpc.Stream) *CrossServiceReRPC_CumSum {
	return &CrossServiceReRPC_CumSum{stream}
}

func (s *CrossServiceReRPC_CumSum) Receive() (*CumSumRequest, error) {
	var req CumSumRequest
	if err := s.stream.Receive(&req); err != nil {
		return nil, err
	}
	return &req, nil
}

func (s *CrossServiceReRPC_CumSum) Send(msg *CumSumResponse) error {
	return s.stream.Send(msg)
}
