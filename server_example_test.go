package rerpc_test

import (
	"context"
	"net/http"
	"time"

	"github.com/rerpc/rerpc"
	"github.com/rerpc/rerpc/health"
	pingrpc "github.com/rerpc/rerpc/internal/gen/proto/go-rerpc/rerpc/ping/v1test"
	pingpb "github.com/rerpc/rerpc/internal/gen/proto/go/rerpc/ping/v1test"
	"github.com/rerpc/rerpc/reflection"
)

// ExamplePingServer implements some trivial business logic. The protobuf
// definition for this API is in internal/ping/v1test/ping.proto.
type ExamplePingServer struct {
	pingrpc.UnimplementedPingServiceServer
}

// Ping implements pingpb.PingServiceReRPC.
func (*ExamplePingServer) Ping(ctx context.Context, req *pingpb.PingRequest) (*pingpb.PingResponse, error) {
	return &pingpb.PingResponse{
		Number: req.Number,
		Msg:    req.Msg,
	}, nil
}

func Example() {
	// The business logic here is trivial, but the rest of the example is meant
	// to be somewhat realistic. This server has basic timeouts configured, and
	// it also exposes gRPC's server reflection and health check APIs.
	ping := &ExamplePingServer{}             // our business logic
	reg := rerpc.NewRegistrar()              // for gRPC reflection
	checker := health.NewChecker(reg)        // basic health checks
	limit := rerpc.ReadMaxBytes(1024 * 1024) // limit request size

	// Next, we convert our implementation of the PingService into a slice
	// of net/http Handlers.
	pingHandler, err := pingrpc.NewPingServiceHandler(ping, reg, limit)
	if err != nil {
		panic(err)
	}

	// NewServeMux returns a plain net/http *ServeMux. Since a mux is an
	// http.Handler, reRPC works with any Go HTTP middleware (e.g., net/http's
	// StripPrefix).
	mux := rerpc.NewServeMux(
		rerpc.NewNotFoundHandler(), // fallback handler
		pingHandler,                // business logic
		reflection.NewHandler(reg), // server reflection
		health.NewHandler(checker), // health checks
	)

	// Timeouts, connection handling, TLS configuration, and other low-level
	// transport details are handled by net/http. Everything you already know (or
	// anything you learn) about hardening net/http Servers applies to reRPC
	// too. Keep in mind that any timeouts you set will also apply to streaming
	// RPCs!
	//
	// If you're not familiar with the many timeouts exposed by net/http, start with
	// https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/.
	srv := &http.Server{
		Addr:           ":http",
		Handler:        mux,
		ReadTimeout:    2500 * time.Millisecond,
		WriteTimeout:   5 * time.Second,
		MaxHeaderBytes: rerpc.MaxHeaderBytes,
	}
	// You could also use golang.org/x/net/http2/h2c to serve gRPC requests
	// without TLS.
	srv.ListenAndServeTLS("testdata/server.crt", "testdata/server.key")
}
