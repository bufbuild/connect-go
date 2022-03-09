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
	"log"
	"net/http"
	"os"
	"time"

	"github.com/bufbuild/connect"
	"github.com/bufbuild/connect/internal/gen/connect/connect/ping/v1test/pingv1testrpc"
	pingv1test "github.com/bufbuild/connect/internal/gen/go/connect/ping/v1test"
)

func Example_client() {
	logger := log.New(os.Stdout, "" /* prefix */, 0 /* flags */)
	// Timeouts, connection pooling, custom dialers, and other low-level
	// transport details are handled by net/http. Everything you already know
	// (or everything you learn) about hardening net/http Clients applies to
	// connect too.
	//
	// Of course, you can skip this configuration and use http.DefaultClient for
	// quick proof-of-concept code.
	httpClient := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			Proxy: nil,
			// connect handles compression on a per-message basis, so it's a waste to
			// compress the whole response body.
			DisableCompression: true,
			MaxIdleConns:       128,
			// RPC clients tend to make many requests to few hosts, so allow more
			// idle connections per host.
			MaxIdleConnsPerHost:    16,
			IdleConnTimeout:        90 * time.Second,
			MaxResponseHeaderBytes: 8 * 1024, // 8 KiB, gRPC's recommended setting
		},
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			// Don't follow any redirects.
			return http.ErrUseLastResponse
		},
	}
	// Unfortunately, pkg.go.dev can't run examples that actually use the
	// network. To keep this example runnable, we'll use an HTTP server and
	// client that communicate over in-memory pipes. Don't do this in production!
	httpClient = examplePingServer.Client()

	client, err := pingv1testrpc.NewPingServiceClient(
		httpClient,
		examplePingServer.URL(),
		connect.WithGRPC(),
	)
	if err != nil {
		logger.Println("error:", err)
		return
	}
	res, err := client.Ping(
		context.Background(),
		connect.NewEnvelope(&pingv1test.PingRequest{Number: 42}),
	)
	if err != nil {
		logger.Println("error:", err)
		return
	}
	logger.Println("response content-type:", res.Header().Get("Content-Type"))
	logger.Println("response message:", res.Msg)

	// Output:
	// response content-type: application/grpc+proto
	// response message: number:42
}
