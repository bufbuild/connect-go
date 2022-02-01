package connect

import (
	"context"
	"net/http"

	"github.com/bufconnect/connect/codec"
)

// A Protocol defines the HTTP semantics to use when sending and receiving
// messages. It ties together codecs, compressors, and net/http to produce
// Senders and Receivers.
//
// For example, connect supports the gRPC protocol using this abstraction. Among
// many other things, the protocol implementation is responsible for
// translating timeouts from Go contexts to HTTP and vice versa. For gRPC, it
// converts timeouts to and from strings (e.g., 10*time.Second <-> "10S"), and
// puts those strings into the "grpc-timeout" HTTP header. Other protocols
// might encode durations differently, put them into a different HTTP header,
// or ignore them entirely.
type protocol interface {
	NewHandler(*protocolHandlerParams) (protocolHandler, error)
	NewClient(*protocolClientParams) (protocolClient, error)
}

// HandlerParams are the arguments provided to a Protocol's NewHandler
// method, bundled into a struct to allow backward-compatible argument
// additions. Protocol implementations should take care to use the supplied
// Specification rather than constructing their own, since new fields may have
// been added.
type protocolHandlerParams struct {
	Spec            Specification
	Codecs          roCodecs
	Compressors     roCompressors
	MaxRequestBytes int64
}

// Handler is the server side of a protocol. HTTP handlers typically support
// multiple protocols, codecs, and compressors.
type protocolHandler interface {
	// ShouldHandleMethod and ShouldHandleContentType check whether the protocol
	// can serve requests with a given HTTP method and content-type. NewStream
	// may assume that any checks in ShouldHandleMethod and
	// ShouldHandleContentType have passed.
	ShouldHandleMethod(string) bool
	ShouldHandleContentType(string) bool

	// If no protocol can serve a request, each protocol's WriteAccept method has
	// a chance to write to the response headers. Protocols should write their
	// supported HTTP methods to the Allow header, and they may write their
	// supported content-types to the Accept-Post or Accept-Patch headers.
	WriteAccept(http.Header)

	// NewStream constructs a Sender and Receiver for the message exchange.
	//
	// Implementations may decide whether the returned error should be sent to
	// the client. (For example, it's helpful to send the client a list of
	// supported compressors if they use an unknown compressor.) If the
	// implementation returns a non-nil Sender, its Close method will be called.
	// If the implementation returns a nil Sender, the error won't be sent to the
	// client.
	//
	// In either case, any returned error is passed through the full interceptor
	// stack. If the implementation returns a nil Sender and/or Receiver, the
	// interceptors receive no-op implementations.
	NewStream(http.ResponseWriter, *http.Request) (Sender, Receiver, error)
}

// ClientParams are the arguments provided to a Protocol's NewClient method,
// bundled into a struct to allow backward-compatible argument additions.
// Protocol implementations should take care to use the supplied Specification
// rather than constructing their own, since new fields may have been added.
type protocolClientParams struct {
	Spec             Specification
	CompressorName   string
	Compressors      roCompressors
	CodecName        string
	Codec            codec.Codec
	MaxResponseBytes int64
	Doer             Doer
	BaseURL          string

	// The gRPC family of protocols always needs access to a protobuf codec to
	// marshal and unmarshal errors.
	Protobuf codec.Codec
}

// Client is the client side of a protocol. HTTP clients typically use a single
// protocol, codec, and compressor to send requests.
type protocolClient interface {
	// WriteRequestHeader writes any protocol-specific request headers.
	WriteRequestHeader(http.Header)

	// NewStream constructs a Sender and Receiver for the message exchange.
	//
	// Implementations should assume that the supplied HTTP headers have already
	// been populated by WriteRequestHeader. When constructing a stream for a
	// unary call, implementations may assume that the sender's Send and Close
	// methods return before the receiver's Receive or Close methods are called.
	NewStream(context.Context, http.Header) (Sender, Receiver)
}
