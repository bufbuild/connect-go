package connect_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bufconnect/connect"
	"github.com/bufconnect/connect/internal/assert"
	pingrpc "github.com/bufconnect/connect/internal/gen/proto/go-connect/connect/ping/v1test"
	pingpb "github.com/bufconnect/connect/internal/gen/proto/go/connect/ping/v1test"
)

type assertCalledInterceptor struct {
	called *bool
}

func (i *assertCalledInterceptor) Wrap(next connect.Func) connect.Func {
	return connect.Func(func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		*i.called = true
		return next(ctx, req)
	})
}

func (i *assertCalledInterceptor) WrapContext(ctx context.Context) context.Context {
	return ctx
}

func (i *assertCalledInterceptor) WrapSender(_ context.Context, s connect.Sender) connect.Sender {
	*i.called = true
	return s
}

func (i *assertCalledInterceptor) WrapReceiver(_ context.Context, r connect.Receiver) connect.Receiver {
	*i.called = true
	return r
}

func TestClientStreamErrors(t *testing.T) {
	_, err := pingrpc.NewPingServiceClient("INVALID_URL", http.DefaultClient)
	assert.NotNil(t, err, "client construction error")
	// We don't even get to calling methods on the client, so there's no question
	// of interceptors running. Once we're calling methods on the client, all
	// errors are visible to interceptors.
}

func TestHandlerStreamErrors(t *testing.T) {
	// If we receive an HTTP request and send a response, interceptors should
	// fire - even if we can't successfully set up a stream. (This is different
	// from clients, where stream creation fails before any HTTP request is
	// issued.)
	var called bool
	reset := func() {
		called = false
	}
	mux, err := connect.NewServeMux(
		connect.NewNotFoundHandler(),
		pingrpc.NewPingServiceHandler(
			pingServer{},
			connect.Interceptors(&assertCalledInterceptor{&called}),
		),
	)
	assert.Nil(t, err, "mux construction error")
	server := httptest.NewServer(mux)
	defer server.Close()

	t.Run("unary", func(t *testing.T) {
		defer reset()
		request, err := http.NewRequest(
			http.MethodPost,
			server.URL+"/connect.ping.v1test.PingService/Ping",
			strings.NewReader(""),
		)
		assert.Nil(t, err, "error constructing request")
		request.Header.Set("Content-Type", "application/grpc+proto")
		request.Header.Set("Grpc-Timeout", "INVALID")
		res, err := server.Client().Do(request)
		assert.Nil(t, err, "network error sending request")
		assert.Equal(t, res.StatusCode, http.StatusOK, "response HTTP status")
		assert.True(t, called, "expected interceptors to be called")
	})
	t.Run("stream", func(t *testing.T) {
		defer reset()
		request, err := http.NewRequest(
			http.MethodPost,
			server.URL+"/connect.ping.v1test.PingService/CountUp",
			strings.NewReader(""),
		)
		assert.Nil(t, err, "error constructing request")
		request.Header.Set("Content-Type", "application/grpc+proto")
		request.Header.Set("Grpc-Timeout", "INVALID")
		res, err := server.Client().Do(request)
		assert.Nil(t, err, "network error sending request")
		assert.Equal(t, res.StatusCode, http.StatusOK, "response HTTP status")
		assert.True(t, called, "expected interceptors to be called")
	})
}

func TestOnionOrderingEndToEnd(t *testing.T) {
	// Helper function: returns a function that asserts that there's some value
	// set for header "expect", and adds a value for header "add".
	newInspector := func(expect, add string) func(connect.Specification,
		http.Header) {
		return func(spec connect.Specification, h http.Header) {
			if expect != "" {
				assert.NotZero(
					t,
					h.Get(expect),
					"%s (IsClient %v): header %q missing: %v",
					assert.Fmt(spec.Procedure, spec.IsClient, expect, h),
				)
			}
			h.Set(add, "v")
		}
	}
	// Helper function: asserts that there's a value present for header keys
	// "one", "two", "three", and "four".
	assertAllPresent := func(spec connect.Specification, h http.Header) {
		for _, k := range []string{"one", "two", "three", "four"} {
			assert.NotZero(
				t,
				h.Get(k),
				"%s (IsClient %v): checking all headers, %q missing: %v",
				assert.Fmt(spec.Procedure, spec.IsClient, k, h),
			)
		}
	}

	// The client and handler interceptor onions are the meat of the test. The
	// order of interceptor execution must be the same for unary and streaming
	// procedures.
	//
	// Requests should fall through the client onion from top to bottom, traverse
	// the network, and then fall through the handler onion from top to bottom.
	// Responses should climb up the handler onion, traverse the network, and
	// then climb up the client onion.
	//
	// The request and response sides of this onion are numbered to make the
	// intended order clear.
	clientOnion := connect.Interceptors(
		connect.NewHeaderInterceptor(
			// 1 (start). request: should see protocol-related headers
			func(_ connect.Specification, h http.Header) {
				assert.NotZero(t, h.Get("Grpc-Accept-Encoding"), "grpc-accept-encoding missing")
			},
			// 12 (end). response: check "one"-"four"
			assertAllPresent,
		),
		connect.NewHeaderInterceptor(
			newInspector("", "one"),       // 2. request: add header "one"
			newInspector("three", "four"), // 11. response: check "three", add "four"
		),
		connect.NewHeaderInterceptor(
			newInspector("one", "two"),   // 3. request: check "one", add "two"
			newInspector("two", "three"), // 10. response: check "two", add "three"
		),
	)
	handlerOnion := connect.Interceptors(
		connect.NewHeaderInterceptor(
			newInspector("two", "three"), // 4. request: check "two", add "three"
			newInspector("one", "two"),   // 9. response: check "one", add "two"
		),
		connect.NewHeaderInterceptor(
			newInspector("three", "four"), // 5. request: check "three", add "four"
			newInspector("", "one"),       // 8. response: add "one"
		),
		connect.NewHeaderInterceptor(
			assertAllPresent, // 6. request: check "one"-"four"
			nil,              // 7. response: no-op
		),
	)

	mux, err := connect.NewServeMux(
		connect.NewNotFoundHandler(),
		pingrpc.NewPingServiceHandler(
			pingServer{},
			handlerOnion,
		),
	)
	assert.Nil(t, err, "mux construction error")
	server := httptest.NewServer(mux)
	defer server.Close()

	client, err := pingrpc.NewPingServiceClient(
		server.URL,
		server.Client(),
		clientOnion,
	)
	assert.Nil(t, err, "client construction error")

	_, err = client.Ping(context.Background(), connect.NewRequest(&pingpb.PingRequest{Number: 10}))
	assert.Nil(t, err, "error calling Ping")

	_, err = client.CountUp(context.Background(), connect.NewRequest(&pingpb.CountUpRequest{Number: 10}))
	assert.Nil(t, err, "error calling CountUp")
}
