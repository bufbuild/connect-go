package rerpc_test

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/rerpc/rerpc"
	pingrpc "github.com/rerpc/rerpc/internal/gen/proto/go-rerpc/rerpc/ping/v1test"
	pingpb "github.com/rerpc/rerpc/internal/gen/proto/go/rerpc/ping/v1test"
)

func ExampleInterceptor() {
	logger := log.New(os.Stdout, "" /* prefix */, 0 /* flags */)
	logProcedure := rerpc.UnaryInterceptorFunc(func(next rerpc.Func) rerpc.Func {
		return rerpc.Func(func(ctx context.Context, req rerpc.AnyRequest) (rerpc.AnyResponse, error) {
			fmt.Println("calling", req.Spec().Procedure)
			return next(ctx, req)
		})
	})
	// This interceptor prevents the client from making network requests in
	// examples. Leave it out in real code!
	short := ShortCircuit(rerpc.Errorf(rerpc.CodeUnimplemented, "no networking in examples"))
	client, err := pingrpc.NewPingServiceClient(
		"https://invalid-test-url",
		http.DefaultClient,
		rerpc.Interceptors(logProcedure, short),
	)
	if err != nil {
		logger.Print("Error: ", err)
		return
	}
	client.Ping(context.Background(), &pingpb.PingRequest{})

	// Output:
	// calling rerpc.ping.v1test.PingService/Ping
}

func ExampleChain() {
	logger := log.New(os.Stdout, "" /* prefix */, 0 /* flags */)
	outer := rerpc.UnaryInterceptorFunc(func(next rerpc.Func) rerpc.Func {
		return rerpc.Func(func(ctx context.Context, req rerpc.AnyRequest) (rerpc.AnyResponse, error) {
			fmt.Println("outer interceptor: before call")
			res, err := next(ctx, req)
			fmt.Println("outer interceptor: after call")
			return res, err
		})
	})
	inner := rerpc.UnaryInterceptorFunc(func(next rerpc.Func) rerpc.Func {
		return rerpc.Func(func(ctx context.Context, req rerpc.AnyRequest) (rerpc.AnyResponse, error) {
			fmt.Println("inner interceptor: before call")
			res, err := next(ctx, req)
			fmt.Println("inner interceptor: after call")
			return res, err
		})
	})
	// This interceptor prevents the client from making network requests in
	// examples. Leave it out in real code!
	short := ShortCircuit(rerpc.Errorf(rerpc.CodeUnimplemented, "no networking in examples"))
	client, err := pingrpc.NewPingServiceClient(
		"https://invalid-test-url",
		http.DefaultClient,
		rerpc.Interceptors(outer, inner, short),
	)
	if err != nil {
		logger.Print("Error: ", err)
		return
	}
	client.Ping(context.Background(), &pingpb.PingRequest{})

	// Output:
	// outer interceptor: before call
	// inner interceptor: before call
	// inner interceptor: after call
	// outer interceptor: after call
}
