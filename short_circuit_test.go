package rerpc_test

import (
	"context"

	"github.com/rerpc/rerpc"
)

type shortCircuit struct {
	err error
}

var _ rerpc.Interceptor = (*shortCircuit)(nil)

func (sc *shortCircuit) Wrap(next rerpc.Func) rerpc.Func {
	return rerpc.Func(func(_ context.Context, _ rerpc.AnyRequest) (rerpc.AnyResponse, error) {
		return nil, sc.err
	})
}

func (sc *shortCircuit) WrapStream(next rerpc.StreamFunc) rerpc.StreamFunc {
	return rerpc.StreamFunc(func(ctx context.Context) (context.Context, rerpc.Sender, rerpc.Receiver) {
		ctx, s, r := next(ctx)
		return ctx, &errSender{Sender: s, err: sc.err}, &errReceiver{Receiver: r, err: sc.err}
	})
}

// ShortCircuit builds an interceptor that doesn't call the wrapped RPC at all.
// Instead, it returns the supplied Error immediately. ShortCircuit works for
// unary and streaming RPCs.
//
// This is primarily useful when testing error handling. It's also used
// throughout reRPC's examples to avoid making network requests.
func ShortCircuit(err error) rerpc.Interceptor {
	return &shortCircuit{err}
}

type errSender struct {
	rerpc.Sender

	err error
}

var _ rerpc.Sender = (*errSender)(nil)

func (s *errSender) Send(_ any) error    { return s.err }
func (s *errSender) Close(_ error) error { return s.err }

type errReceiver struct {
	rerpc.Receiver

	err error
}

var _ rerpc.Receiver = (*errReceiver)(nil)

func (r *errReceiver) Receive(_ any) error { return r.err }
func (r *errReceiver) Close() error        { return r.err }
