package connect

import (
	"net/http"
)

// ClientStreamForClient is the client's view of a client streaming RPC.
type ClientStreamForClient[Req, Res any] struct {
	sender   Sender
	receiver Receiver
}

// NewClientStreamForClient constructs the client's view of a client streaming RPC.
func NewClientStreamForClient[Req, Res any](s Sender, r Receiver) *ClientStreamForClient[Req, Res] {
	return &ClientStreamForClient[Req, Res]{sender: s, receiver: r}
}

// RequestHeader returns the request headers. Headers are sent to the server with the
// first call to Send.
func (c *ClientStreamForClient[Req, Res]) RequestHeader() http.Header {
	return c.sender.Header()
}

// Send a message to the server. The first call to Send also sends the request
// headers.
func (c *ClientStreamForClient[Req, Res]) Send(msg *Req) error {
	return c.sender.Send(msg)
}

// CloseAndReceive closes the send side of the stream and waits for the
// response.
func (c *ClientStreamForClient[Req, Res]) CloseAndReceive() (*Envelope[Res], error) {
	if err := c.sender.Close(nil); err != nil {
		return nil, err
	}
	res, err := ReceiveUnaryEnvelope[Res](c.receiver)
	if err != nil {
		_ = c.receiver.Close()
		return nil, err
	}
	if err := c.receiver.Close(); err != nil {
		return nil, err
	}
	return res, nil
}

// ServerStreamForClient is the client's view of a server streaming RPC.
type ServerStreamForClient[Res any] struct {
	receiver Receiver
}

// NewServerStreamForClient constructs the client's view of a server streaming
// RPC.
func NewServerStreamForClient[Res any](r Receiver) *ServerStreamForClient[Res] {
	return &ServerStreamForClient[Res]{receiver: r}
}

// Receive a message. When the server is done sending messages, Receive will
// return an error that wraps io.EOF.
func (s *ServerStreamForClient[Res]) Receive() (*Res, error) {
	var res Res
	if err := s.receiver.Receive(&res); err != nil {
		return nil, err
	}
	return &res, nil
}

// ResponseHeader returns the headers received from the server. It blocks until
// the first call to Receive returns.
func (s *ServerStreamForClient[Res]) ResponseHeader() http.Header {
	return s.receiver.Header()
}

// ResponseTrailer returns the trailers received from the server. Trailers
// aren't fully populated until Receive() returns an error wrapping io.EOF.
func (s *ServerStreamForClient[Res]) ResponseTrailer() http.Header {
	return s.receiver.Trailer()
}

// Close the receive side of the stream.
func (s *ServerStreamForClient[Res]) Close() error {
	return s.receiver.Close()
}

// BidiStreamForClient is the client's view of a bidirectional streaming RPC.
type BidiStreamForClient[Req, Res any] struct {
	sender   Sender
	receiver Receiver
}

// NewBidiStreamForClient constructs the client's view of a bidirectional streaming RPC.
func NewBidiStreamForClient[Req, Res any](s Sender, r Receiver) *BidiStreamForClient[Req, Res] {
	return &BidiStreamForClient[Req, Res]{sender: s, receiver: r}
}

// RequestHeader returns the request headers. Headers are sent with the first
// call to Send.
func (b *BidiStreamForClient[Req, Res]) RequestHeader() http.Header {
	return b.sender.Header()
}

// Send a message to the server. The first call to Send also sends the request
// headers.
func (b *BidiStreamForClient[Req, Res]) Send(msg *Req) error {
	return b.sender.Send(msg)
}

// CloseSend closes the send side of the stream.
func (b *BidiStreamForClient[Req, Res]) CloseSend() error {
	return b.sender.Close(nil)
}

// Receive a message. When the server is done sending messages, Receive will
// return an error that wraps io.EOF.
func (b *BidiStreamForClient[Req, Res]) Receive() (*Res, error) {
	var res Res
	if err := b.receiver.Receive(&res); err != nil {
		return nil, err
	}
	return &res, nil
}

// CloseReceive closes the receive side of the stream.
func (b *BidiStreamForClient[Req, Res]) CloseReceive() error {
	return b.receiver.Close()
}

// ResponseHeader returns the headers received from the server. It blocks until
// the first call to Receive returns.
func (b *BidiStreamForClient[Req, Res]) ResponseHeader() http.Header {
	return b.receiver.Header()
}

// ResponseTrailer returns the trailers received from the server. Trailers
// aren't fully populated until Receive() returns an error wrapping io.EOF.
func (b *BidiStreamForClient[Req, Res]) ResponseTrailer() http.Header {
	return b.receiver.Trailer()
}
