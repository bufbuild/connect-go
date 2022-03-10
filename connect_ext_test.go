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
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/bufbuild/connect"
	"github.com/bufbuild/connect/internal/assert"
	"github.com/bufbuild/connect/internal/gen/connect/connect/ping/v1/pingv1connect"
	pingv1 "github.com/bufbuild/connect/internal/gen/go/connect/ping/v1"
)

const errorMessage = "oh no"

const (
	headerValue    = "some header value"
	trailerValue   = "some trailer value"
	clientHeader   = "Connect-Client-Header"
	handlerHeader  = "Connect-Handler-Header"
	handlerTrailer = "Connect-Handler-Trailer"
)

func expectClientHeader(check bool, req connect.AnyRequest) error {
	if !check {
		return nil
	}
	if err := expectMetadata(req.Header(), "header", clientHeader, headerValue); err != nil {
		return err
	}
	return nil
}

func expectMetadata(meta http.Header, metaType, key, value string) error {
	if got := meta.Get(key); got != value {
		return connect.NewError(connect.CodeInvalidArgument, fmt.Errorf(
			"%s %q: got %q, expected %q",
			metaType,
			key,
			got,
			value,
		))
	}
	return nil
}

type pingServer struct {
	pingv1connect.UnimplementedPingServiceHandler

	checkMetadata bool
}

func (p pingServer) Ping(ctx context.Context, req *connect.Request[pingv1.PingRequest]) (*connect.Response[pingv1.PingResponse], error) {
	if err := expectClientHeader(p.checkMetadata, req); err != nil {
		return nil, err
	}
	res := connect.NewResponse(&pingv1.PingResponse{
		Number: req.Msg.Number,
		Text:   req.Msg.Text,
	})
	res.Header().Set(handlerHeader, headerValue)
	res.Trailer().Set(handlerTrailer, trailerValue)
	return res, nil
}

func (p pingServer) Fail(ctx context.Context, req *connect.Request[pingv1.FailRequest]) (*connect.Response[pingv1.FailResponse], error) {
	if err := expectClientHeader(p.checkMetadata, req); err != nil {
		return nil, err
	}
	err := connect.NewError(connect.Code(req.Msg.Code), errors.New(errorMessage))
	err.Header().Set(handlerHeader, headerValue)
	err.Trailer().Set(handlerTrailer, trailerValue)
	return nil, err
}

func (p pingServer) Sum(
	ctx context.Context,
	stream *connect.ClientStream[pingv1.SumRequest, pingv1.SumResponse],
) error {
	if p.checkMetadata {
		if err := expectMetadata(stream.RequestHeader(), "header", clientHeader, headerValue); err != nil {
			return err
		}
	}
	var sum int64
	for stream.Receive() {
		if err := ctx.Err(); err != nil {
			return err
		}
		sum += stream.Msg().Number
	}
	if stream.Err() != nil {
		return stream.Err()
	}
	response := connect.NewResponse(&pingv1.SumResponse{Sum: sum})
	response.Header().Set(handlerHeader, headerValue)
	response.Trailer().Set(handlerTrailer, trailerValue)
	return stream.SendAndClose(response)
}

func (p pingServer) CountUp(
	ctx context.Context,
	req *connect.Request[pingv1.CountUpRequest],
	stream *connect.ServerStream[pingv1.CountUpResponse],
) error {
	if err := expectClientHeader(p.checkMetadata, req); err != nil {
		return err
	}
	if req.Msg.Number <= 0 {
		return connect.NewError(connect.CodeInvalidArgument, fmt.Errorf(
			"number must be positive: got %v",
			req.Msg.Number,
		))
	}
	stream.ResponseHeader().Set(handlerHeader, headerValue)
	stream.ResponseTrailer().Set(handlerTrailer, trailerValue)
	for i := int64(1); i <= req.Msg.Number; i++ {
		if err := ctx.Err(); err != nil {
			return err
		}
		if err := stream.Send(&pingv1.CountUpResponse{Number: i}); err != nil {
			return err
		}
	}
	return nil
}

func (p pingServer) CumSum(
	ctx context.Context,
	stream *connect.BidiStream[pingv1.CumSumRequest, pingv1.CumSumResponse],
) error {
	var sum int64
	if p.checkMetadata {
		if err := expectMetadata(stream.RequestHeader(), "header", clientHeader, headerValue); err != nil {
			return err
		}
	}
	stream.ResponseHeader().Set(handlerHeader, headerValue)
	stream.ResponseTrailer().Set(handlerTrailer, trailerValue)
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		msg, err := stream.Receive()
		if errors.Is(err, io.EOF) {
			return nil
		} else if err != nil {
			return err
		}
		sum += msg.Number
		if err := stream.Send(&pingv1.CumSumResponse{Sum: sum}); err != nil {
			return err
		}
	}
}

func TestServerProtoGRPC(t *testing.T) {
	mux := http.NewServeMux()
	mux.Handle(pingv1connect.NewPingServiceHandler(
		pingServer{checkMetadata: true},
	))

	testPing := func(t *testing.T, client pingv1connect.PingServiceClient) {
		t.Run("ping", func(t *testing.T) {
			num := rand.Int63()
			req := connect.NewRequest(&pingv1.PingRequest{Number: num})
			req.Header().Set(clientHeader, headerValue)
			expect := &pingv1.PingResponse{Number: num}
			res, err := client.Ping(context.Background(), req)
			assert.Nil(t, err)
			assert.Equal(t, res.Msg, expect)
			assert.Equal(t, res.Header().Get(handlerHeader), headerValue)
			assert.Equal(t, res.Trailer().Get(handlerTrailer), trailerValue)
		})
		t.Run("large ping", func(t *testing.T) {
			// Using a large payload splits the request and response over multiple
			// packets, ensuring that we're managing HTTP readers and writers
			// correctly.
			hellos := strings.Repeat("hello", 1024*1024) // ~5mb
			req := connect.NewRequest(&pingv1.PingRequest{Text: hellos})
			req.Header().Set(clientHeader, headerValue)
			res, err := client.Ping(context.Background(), req)
			assert.Nil(t, err)
			assert.Equal(t, res.Msg.Text, hellos)
			assert.Equal(t, res.Header().Get(handlerHeader), headerValue)
			assert.Equal(t, res.Trailer().Get(handlerTrailer), trailerValue)
		})
	}
	testSum := func(t *testing.T, client pingv1connect.PingServiceClient) {
		t.Run("sum", func(t *testing.T) {
			const (
				upTo   = 10
				expect = 55 // 1+10 + 2+9 + ... + 5+6 = 55
			)
			stream := client.Sum(context.Background())
			stream.RequestHeader().Set(clientHeader, headerValue)
			for i := int64(1); i <= upTo; i++ {
				err := stream.Send(&pingv1.SumRequest{Number: i})
				assert.Nil(t, err, assert.Sprintf("send %d", i))
			}
			res, err := stream.CloseAndReceive()
			assert.Nil(t, err)
			assert.Equal(t, res.Msg.Sum, expect)
			assert.Equal(t, res.Header().Get(handlerHeader), headerValue)
			assert.Equal(t, res.Trailer().Get(handlerTrailer), trailerValue)
		})
	}
	testCountUp := func(t *testing.T, client pingv1connect.PingServiceClient) {
		t.Run("count_up", func(t *testing.T) {
			const n = 5
			got := make([]int64, 0, n)
			expect := make([]int64, 0, n)
			for i := 1; i <= n; i++ {
				expect = append(expect, int64(i))
			}
			request := connect.NewRequest(&pingv1.CountUpRequest{Number: n})
			request.Header().Set(clientHeader, headerValue)
			stream, err := client.CountUp(context.Background(), request)
			assert.Nil(t, err)
			for stream.Receive() {
				got = append(got, stream.Msg().Number)
			}
			assert.Nil(t, stream.Err())
			assert.Nil(t, stream.Close())
			assert.Equal(t, got, expect)
		})
	}
	testCumSum := func(t *testing.T, client pingv1connect.PingServiceClient, expectSuccess bool) {
		t.Run("cumsum", func(t *testing.T) {
			send := []int64{3, 5, 1}
			expect := []int64{3, 8, 9}
			var got []int64
			stream := client.CumSum(context.Background())
			stream.RequestHeader().Set(clientHeader, headerValue)
			if !expectSuccess { // server doesn't support HTTP/2
				err := stream.Send(&pingv1.CumSumRequest{})
				assert.Nil(t, err) // haven't gotten response back yet
				assert.Nil(t, stream.CloseSend())
				_, err = stream.Receive()
				assert.NotNil(t, err) // should be 505
				assert.True(
					t,
					strings.Contains(err.Error(), "HTTP status 505"),
					assert.Sprintf("expected 505, got %v", err),
				)
				assert.Nil(t, stream.CloseReceive())
				return
			}
			var wg sync.WaitGroup
			wg.Add(2)
			go func() {
				defer wg.Done()
				for i, n := range send {
					err := stream.Send(&pingv1.CumSumRequest{Number: n})
					assert.Nil(t, err, assert.Sprintf("send error #%d", i))
				}
				assert.Nil(t, stream.CloseSend())
			}()
			go func() {
				defer wg.Done()
				for {
					msg, err := stream.Receive()
					if errors.Is(err, io.EOF) {
						break
					}
					assert.Nil(t, err)
					got = append(got, msg.Sum)
				}
				assert.Nil(t, stream.CloseReceive())
			}()
			wg.Wait()
			assert.Equal(t, got, expect)
			assert.Equal(t, stream.ResponseHeader().Get(handlerHeader), headerValue)
			assert.Equal(t, stream.ResponseTrailer().Get(handlerTrailer), trailerValue)
		})
	}
	testErrors := func(t *testing.T, client pingv1connect.PingServiceClient) {
		t.Run("errors", func(t *testing.T) {
			request := connect.NewRequest(&pingv1.FailRequest{
				Code: int32(connect.CodeResourceExhausted),
			})
			request.Header().Set(clientHeader, headerValue)

			response, err := client.Fail(context.Background(), request)
			assert.Nil(t, response)
			assert.NotNil(t, err)
			var connectErr *connect.Error
			ok := errors.As(err, &connectErr)
			assert.True(t, ok, assert.Sprintf("conversion to *connect.Error"))
			assert.Equal(t, connectErr.Code(), connect.CodeResourceExhausted)
			assert.Equal(t, connectErr.Error(), "ResourceExhausted: "+errorMessage)
			assert.Zero(t, connectErr.Details())
			assert.Equal(t, connectErr.Header().Get(handlerHeader), headerValue)
			assert.Equal(t, connectErr.Trailer().Get(handlerTrailer), trailerValue)
		})
	}
	testMatrix := func(t *testing.T, server *httptest.Server, bidi bool) {
		run := func(t *testing.T, opts ...connect.ClientOption) {
			client, err := pingv1connect.NewPingServiceClient(server.Client(), server.URL, opts...)
			assert.Nil(t, err)
			testPing(t, client)
			testSum(t, client)
			testCountUp(t, client)
			testCumSum(t, client, bidi)
			testErrors(t, client)
		}
		t.Run("identity", func(t *testing.T) {
			run(t, connect.WithGRPC())
		})
		t.Run("gzip", func(t *testing.T) {
			run(t, connect.WithGRPC(), connect.WithGzipRequests())
		})
		t.Run("json_gzip", func(t *testing.T) {
			run(
				t,
				connect.WithGRPC(),
				connect.WithProtoJSONCodec(),
				connect.WithGzipRequests(),
			)
		})
		t.Run("web", func(t *testing.T) {
			run(t, connect.WithGRPCWeb())
		})
		t.Run("web_json_gzip", func(t *testing.T) {
			run(
				t,
				connect.WithGRPCWeb(),
				connect.WithProtoJSONCodec(),
				connect.WithGzipRequests(),
			)
		})
	}

	t.Run("http1", func(t *testing.T) {
		server := httptest.NewServer(mux)
		defer server.Close()
		testMatrix(t, server, false /* bidi */)
	})
	t.Run("http2", func(t *testing.T) {
		server := httptest.NewUnstartedServer(mux)
		server.EnableHTTP2 = true
		server.StartTLS()
		defer server.Close()
		testMatrix(t, server, true /* bidi */)
	})
}

type pluggablePingServer struct {
	pingv1connect.UnimplementedPingServiceHandler

	ping func(context.Context, *connect.Request[pingv1.PingRequest]) (*connect.Response[pingv1.PingResponse], error)
}

func (p *pluggablePingServer) Ping(
	ctx context.Context,
	req *connect.Request[pingv1.PingRequest],
) (*connect.Response[pingv1.PingResponse], error) {
	return p.ping(ctx, req)
}

func TestHeaderBasic(t *testing.T) {
	const (
		key  = "Test-Key"
		cval = "client value"
		hval = "client value"
	)

	pingServer := &pluggablePingServer{
		ping: func(ctx context.Context, req *connect.Request[pingv1.PingRequest]) (*connect.Response[pingv1.PingResponse], error) {
			assert.Equal(t, req.Header().Get(key), cval)
			res := connect.NewResponse(&pingv1.PingResponse{})
			res.Header().Set(key, hval)
			return res, nil
		},
	}
	mux := http.NewServeMux()
	mux.Handle(pingv1connect.NewPingServiceHandler(pingServer))
	server := httptest.NewServer(mux)
	defer server.Close()

	client, err := pingv1connect.NewPingServiceClient(server.Client(), server.URL, connect.WithGRPC())
	assert.Nil(t, err)
	req := connect.NewRequest(&pingv1.PingRequest{})
	req.Header().Set(key, cval)
	res, err := client.Ping(context.Background(), req)
	assert.Nil(t, err)
	assert.Equal(t, res.Header().Get(key), hval)
}
