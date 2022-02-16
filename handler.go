package connect

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/bufconnect/connect/codec"
	"github.com/bufconnect/connect/compress"
	"github.com/bufconnect/connect/compress/gzip"
)

type handlerCfg struct {
	Compressors      map[string]compress.Compressor
	Codecs           map[string]codec.Codec
	MaxRequestBytes  int64
	Registrar        *Registrar
	Interceptor      Interceptor
	Procedure        string
	RegistrationName string
	HandleGRPC       bool
	HandleGRPCWeb    bool
}

func newHandlerConfiguration(procedure, registrationName string, opts []HandlerOption) (*handlerCfg, *Error) {
	cfg := handlerCfg{
		Procedure:        procedure,
		RegistrationName: registrationName,
		Compressors: map[string]compress.Compressor{
			gzip.Name: gzip.New(),
		},
		Codecs:        make(map[string]codec.Codec),
		HandleGRPC:    true,
		HandleGRPCWeb: true,
	}
	for _, opt := range opts {
		opt.applyToHandler(&cfg)
	}
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	if reg := cfg.Registrar; reg != nil && cfg.RegistrationName != "" {
		reg.register(cfg.RegistrationName)
	}
	return &cfg, nil
}

func (c *handlerCfg) Validate() *Error {
	if _, ok := c.Codecs[""]; ok {
		return Wrap(
			CodeUnknown,
			errors.New("can't register codec with an empty name"),
		)
	}
	if _, ok := c.Compressors[""]; ok {
		return Wrap(
			CodeUnknown,
			errors.New("can't register compressor with an empty name"),
		)
	}
	return nil
}

func (c *handlerCfg) newSpecification(t StreamType) Specification {
	return Specification{
		Procedure: c.Procedure,
		Type:      t,
	}
}

func (c *handlerCfg) newProtocolHandlers(stype StreamType) ([]protocolHandler, *Error) {
	var protocols []protocol
	if c.HandleGRPC {
		protocols = append(protocols, &grpc{web: false})
	}
	if c.HandleGRPCWeb {
		protocols = append(protocols, &grpc{web: true})
	}
	handlers := make([]protocolHandler, 0, len(protocols))
	codecs := newReadOnlyCodecs(c.Codecs)
	compressors := newReadOnlyCompressors(c.Compressors)
	for _, p := range protocols {
		ph, err := p.NewHandler(&protocolHandlerParams{
			Spec:            c.newSpecification(stype),
			Codecs:          codecs,
			Compressors:     compressors,
			MaxRequestBytes: c.MaxRequestBytes,
		})
		if err != nil {
			return nil, Wrap(CodeUnknown, err)
		}
		handlers = append(handlers, ph)
	}
	return handlers, nil
}

// A HandlerOption configures a Handler.
//
// In addition to any options grouped in the documentation below, remember that
// Registrars and Options are also valid HandlerOptions.
type HandlerOption interface {
	applyToHandler(*handlerCfg)
}

type handleGRPCWebOption struct {
	enable bool
}

// HandleGRPCWeb enables or disables support for the gRPC-Web protocol. By
// default, gRPC-Web is enabled. Note that handlers always support the standard
// HTTP/2 gRPC protocol.
func HandleGRPCWeb(enable bool) HandlerOption {
	return &handleGRPCWebOption{enable}
}

func (o *handleGRPCWebOption) applyToHandler(c *handlerCfg) {
	c.HandleGRPCWeb = o.enable
}

// A Handler is the server-side implementation of a single RPC defined by a
// protocol buffer service. It's the interface between the connect library and
// the code generated by the connect protoc plugin; most users won't ever need
// to deal with it directly.
//
// To see an example of how Handler is used in the generated code, see the
// internal/gen/proto/go-connect/connect/ping/v1test package.
type Handler struct {
	spec             Specification
	interceptor      Interceptor
	implementation   func(context.Context, Sender, Receiver, error /* client-visible */)
	protocolHandlers []protocolHandler
}

// NewUnaryHandler constructs a Handler. The supplied package, service, and
// method names must be protobuf identifiers. For example, a handler for the
// URL "/acme.foo.v1.FooService/Bar" would have package "acme.foo.v1", service
// "FooService", and method "Bar".
//
// Remember that NewUnaryHandler is usually called from generated code - most
// users won't need to deal with protobuf identifiers directly.
func NewUnaryHandler[Req, Res any](
	procedure, registrationName string,
	unary func(context.Context, *Request[Req]) (*Response[Res], error),
	opts ...HandlerOption,
) (*Handler, error) {
	cfg, err := newHandlerConfiguration(procedure, registrationName, opts)
	if err != nil {
		return nil, err
	}
	implementation := func(ctx context.Context, sender Sender, receiver Receiver, clientVisibleError error) {
		defer receiver.Close()

		var req *Request[Req]
		if clientVisibleError != nil {
			// The protocol implementation failed to establish a stream. To make the
			// resulting error visible to the interceptor stack, we still want to
			// call the wrapped unary Func. To do that safely, we need a useful
			// Request struct. (Note that we are *not* actually calling the
			// handler's implementation.)
			req = receiveRequestMetadata[Req](receiver)
		} else {
			var err error
			req, err = ReceiveRequest[Req](receiver)
			if err != nil {
				// Interceptors should see this error too. Just as above, they need a
				// useful Request.
				clientVisibleError = err
				req = receiveRequestMetadata[Req](receiver)
			}
		}

		untyped := Func(func(ctx context.Context, req AnyRequest) (AnyResponse, error) {
			if clientVisibleError != nil {
				// We've already encountered an error, short-circuit before calling the
				// handler's implementation.
				return nil, clientVisibleError
			}
			if err := ctx.Err(); err != nil {
				return nil, err
			}
			typed, ok := req.(*Request[Req])
			if !ok {
				return nil, Errorf(CodeInternal, "unexpected handler request type %T", req)
			}
			return unary(ctx, typed)
		})
		if ic := cfg.Interceptor; ic != nil {
			untyped = ic.Wrap(untyped)
		}

		res, err := untyped(ctx, req)
		if err != nil {
			_ = sender.Close(err)
			return
		}
		mergeHeaders(sender.Header(), res.Header())
		_ = sender.Close(sender.Send(res.Any()))
	}

	protocolHandlers, err := cfg.newProtocolHandlers(StreamTypeUnary)
	if err != nil {
		return nil, err
	}
	return &Handler{
		spec:             cfg.newSpecification(StreamTypeUnary),
		interceptor:      nil, // already applied
		implementation:   implementation,
		protocolHandlers: protocolHandlers,
	}, nil
}

// NewStreamingHandler constructs a Handler. The supplied package, service, and
// method names must be protobuf identifiers. For example, a handler for the
// URL "/acme.foo.v1.FooService/Bar" would have package "acme.foo.v1", service
// "FooService", and method "Bar".
//
// Remember that NewStreamingHandler is usually called from generated code -
// most users won't need to deal with protobuf identifiers directly.
func NewStreamingHandler(
	stype StreamType,
	procedure, registrationName string,
	implementation func(context.Context, Sender, Receiver),
	opts ...HandlerOption,
) (*Handler, error) {
	cfg, err := newHandlerConfiguration(procedure, registrationName, opts)
	if err != nil {
		return nil, err
	}
	protocolHandlers, err := cfg.newProtocolHandlers(stype)
	if err != nil {
		return nil, err
	}
	return &Handler{
		spec:        cfg.newSpecification(stype),
		interceptor: cfg.Interceptor,
		implementation: func(ctx context.Context, s Sender, r Receiver, clientVisibleErr error) {
			if clientVisibleErr != nil {
				_ = r.Close()
				_ = s.Close(clientVisibleErr)
				return
			}
			implementation(ctx, s, r)
		},
		protocolHandlers: protocolHandlers,
	}, nil
}

// ServeHTTP implements http.Handler.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// We don't need to defer functions  to close the request body or read to
	// EOF: the stream we construct later on already does that, and we only
	// return early when dealing with misbehaving clients. In those cases, it's
	// okay if we can't re-use the connection.
	isBidi := (h.spec.Type & StreamTypeBidirectional) == StreamTypeBidirectional
	if isBidi && r.ProtoMajor < 2 {
		h.failNegotiation(w, http.StatusHTTPVersionNotSupported)
		return
	}

	methodHandlers := make([]protocolHandler, 0, len(h.protocolHandlers))
	for _, ph := range h.protocolHandlers {
		if ph.ShouldHandleMethod(r.Method) {
			methodHandlers = append(methodHandlers, ph)
		}
	}
	if len(methodHandlers) == 0 {
		// grpc-go returns a 500 here, but interoperability with non-gRPC HTTP
		// clients is better if we return a 405.
		h.failNegotiation(w, http.StatusMethodNotAllowed)
		return
	}

	// TODO: for GETs, we should parse the Accept header and offer each handler
	// each content-type.
	ctype := r.Header.Get("Content-Type")
	for _, ph := range methodHandlers {
		if !ph.ShouldHandleContentType(ctype) {
			continue
		}
		ctx := r.Context()
		if ic := h.interceptor; ic != nil {
			ctx = ic.WrapContext(ctx)
		}
		// Most errors returned from ph.NewStream are caused by invalid requests.
		// For example, the client may have specified an invalid timeout or an
		// unavailable codec. We'd like those errors to be visible to the
		// interceptor chain, so we're going to capture them here and pass them to
		// the implementation.
		sender, receiver, clientVisibleError := ph.NewStream(w, r.WithContext(ctx))
		// If NewStream errored and the protocol doesn't want the error sent to
		// the client, sender and/or receiver may be nil. We still want the
		// error to be seen by interceptors, so we provide no-op Sender and
		// Receiver implementations.
		if clientVisibleError != nil && sender == nil {
			sender = newNopSender(h.spec, w.Header())
		}
		if clientVisibleError != nil && receiver == nil {
			receiver = newNopReceiver(h.spec, r.Header)
		}
		if ic := h.interceptor; ic != nil {
			// Unary interceptors were handled in NewUnaryHandler.
			sender = ic.WrapSender(ctx, sender)
			receiver = ic.WrapReceiver(ctx, receiver)
		}
		h.implementation(ctx, sender, receiver, clientVisibleError)
		return
	}
	h.failNegotiation(w, http.StatusUnsupportedMediaType)
}

// Path returns the URL pattern to use when registering this handler.
func (h *Handler) path() string {
	return fmt.Sprintf("/" + h.spec.Procedure)
}

func (h *Handler) failNegotiation(w http.ResponseWriter, code int) {
	// None of the registered protocols is able to serve the request.
	for _, ph := range h.protocolHandlers {
		ph.WriteAccept(w.Header())
	}
	w.WriteHeader(code)
}
