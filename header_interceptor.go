package rerpc

import (
	"context"
)

// HeaderInterceptor makes it easier to write interceptors that inspect or
// mutate HTTP headers. It applies the same logic to unary and streaming
// procedures, wrapping the send or receive side of the stream as appropriate.
type HeaderInterceptor struct {
	inspectRequestHeader  func(Specification, Header)
	inspectResponseHeader func(Specification, Header)
}

var _ Interceptor = (*HeaderInterceptor)(nil)

// NewHeaderInterceptor constructs a HeaderInterceptor. Nil function pointers
// are treated as no-ops.
func NewHeaderInterceptor(
	inspectRequestHeader func(Specification, Header),
	inspectResponseHeader func(Specification, Header),
) *HeaderInterceptor {
	h := HeaderInterceptor{
		inspectRequestHeader:  inspectRequestHeader,
		inspectResponseHeader: inspectResponseHeader,
	}
	if h.inspectRequestHeader == nil {
		h.inspectRequestHeader = func(_ Specification, _ Header) {}
	}
	if h.inspectResponseHeader == nil {
		h.inspectResponseHeader = func(_ Specification, _ Header) {}
	}
	return &h
}

// Wrap implements Interceptor.
func (h *HeaderInterceptor) Wrap(next Func) Func {
	return Func(func(ctx context.Context, req AnyRequest) (AnyResponse, error) {
		h.inspectRequestHeader(req.Spec(), req.Header())
		res, err := next(ctx, req)
		if err != nil {
			return nil, err
		}
		h.inspectResponseHeader(req.Spec(), res.Header())
		return res, nil
	})
}

// WrapContext implements Interceptor with a no-op.
func (h *HeaderInterceptor) WrapContext(ctx context.Context) context.Context {
	return ctx
}

// WrapSender implements Interceptor. Depending on whether it's operating on a
// client or handler, it wraps the sender with the request- or
// response-inspecting function.
func (h *HeaderInterceptor) WrapSender(ctx context.Context, sender Sender) Sender {
	if sender.Spec().IsClient {
		return &headerMutatingSender{Sender: sender, mutate: h.inspectRequestHeader}
	}
	return &headerMutatingSender{Sender: sender, mutate: h.inspectResponseHeader}
}

// WrapReceiver implements Interceptor. Depending on whether it's operating on a
// client or handler, it wraps the sender with the response- or
// request-inspecting function.
func (h *HeaderInterceptor) WrapReceiver(ctx context.Context, receiver Receiver) Receiver {
	if receiver.Spec().IsClient {
		return &headerMutatingReceiver{Receiver: receiver, mutate: h.inspectResponseHeader}
	}
	return &headerMutatingReceiver{Receiver: receiver, mutate: h.inspectRequestHeader}
}

type headerMutatingSender struct {
	Sender

	called bool // senders don't need to be thread-safe
	mutate func(Specification, Header)
}

func (s *headerMutatingSender) Send(m any) error {
	if !s.called {
		s.mutate(s.Spec(), s.Header())
		s.called = true
	}
	return s.Sender.Send(m)
}

type headerMutatingReceiver struct {
	Receiver

	called bool // receivers don't need to be thread-safe
	mutate func(Specification, Header)
}

func (r *headerMutatingReceiver) Receive(m any) error {
	if !r.called {
		r.mutate(r.Spec(), r.Header())
		r.called = true
	}
	if err := r.Receiver.Receive(m); err != nil {
		return err
	}
	return nil
}
