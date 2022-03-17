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

package connect_test

import (
	"context"
	"net/http"
	"time"

	"github.com/bufbuild/connect"
	"github.com/bufbuild/connect/grpchealth"
	"github.com/bufbuild/connect/grpcreflect"
	"github.com/bufbuild/connect/internal/gen/connect/connect/ping/v1/pingv1connect"
	pingv1 "github.com/bufbuild/connect/internal/gen/go/connect/ping/v1"
)

// ExamplePingServer implements some trivial business logic. The protobuf
// definition for this API is in proto/connect/ping/v1/ping.proto.
type ExamplePingServer struct {
	pingv1connect.UnimplementedPingServiceHandler
}

// Ping implements pingv1connect.PingServiceHandler.
func (*ExamplePingServer) Ping(
	_ context.Context,
	req *connect.Request[pingv1.PingRequest],
) (*connect.Response[pingv1.PingResponse], error) {
	return connect.NewResponse(&pingv1.PingResponse{
		Number: req.Msg.Number,
		Text:   req.Msg.Text,
	}), nil
}

func Example_handler() {
	// The business logic here is trivial, but the rest of the example is meant
	// to be somewhat realistic. This server has basic timeouts configured, and
	// it also exposes gRPC's server reflection and health check APIs.

	// protoc-gen-connect-go generates constructors that return plain net/http
	// Handlers, so they're compatible with most Go HTTP routers and middleware
	// (for example, net/http's StripPrefix).
	mux := http.NewServeMux()
	mux.Handle(pingv1connect.NewPingServiceHandler(
		&ExamplePingServer{},                // our business logic
		connect.WithReadMaxBytes(1024*1024), // limit request size
	))
	// The grpchealth and grpcreflection sub-packages offer support for the
	// standard gRPC health checking and server reflection APIs. Serving these
	// APIs makes it easy to integrate your connect server with Kubernetes health
	// checks, CLI tools like grpcurl, and a variety of other systems.
	services := []string{pingv1connect.PingServiceName}
	mux.Handle(grpchealth.NewHandler(services))
	mux.Handle(grpcreflect.NewHandler(services))

	// Timeouts, connection handling, TLS configuration, and other low-level
	// transport details are handled by net/http. Everything you already know (or
	// anything you learn) about hardening net/http Servers applies to connect
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
		MaxHeaderBytes: 8 * 1024, // 8KiB, gRPC's recommendation
	}
	// You could also use golang.org/x/net/http2/h2c to serve gRPC requests
	// without TLS.
	srv.ListenAndServeTLS("testdata/server.crt", "testdata/server.key")
}
