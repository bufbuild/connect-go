package rerpc

import (
	"errors"
	"fmt"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

// An Error captures three pieces of information: a Code, a human-readable
// message, and an optional collection of arbitrary protobuf messages called
// "details" (more on those below). Servers send the code, message, and details
// over the wire to clients. reRPC's Error wraps a standard Go error, using the
// underlying error's Error() string as the message. Take care not to leak
// sensitive information from public APIs!
//
// Protobuf service implementations and Interceptors should return Errors
// (using the Wrap or Errorf functions) rather than plain Go errors. If service
// implementations or Interceptors instead return a plain Go error, reRPC will
// use AsError to find an Error to send over the wire. If no Error can be
// found, reRPC will use CodeUnknown and the returned error's message.
//
// Error codes and messages are explained in the gRPC documentation linked
// below. Unfortunately, error details were introduced before gRPC adopted a
// formal proposal process, so they're not clearly documented anywhere and
// may differ slightly between implementations. Roughly, they're an optional
// mechanism for servers, middleware, and proxies to send strongly-typed errors
// and localized messages to clients.
//
// See https://github.com/grpc/grpc/blob/master/doc/PROTOCOL-HTTP2.md and
// https://github.com/grpc/grpc/blob/master/doc/statuscodes.md for further
// details.
type Error struct {
	code    Code
	err     error
	details []*anypb.Any
}

// Wrap annotates any error with a status code and error details. If the code
// is CodeOK, the returned error is nil.
func Wrap(c Code, err error, details ...proto.Message) *Error {
	if c == CodeOK {
		return nil
	}
	e := &Error{
		code: c,
		err:  err,
	}
	if len(details) > 0 {
		e.SetDetails(details...)
	}
	return e
}

// Errorf calls fmt.Errorf with the supplied template and arguments, then wraps
// the resulting error. If the code is CodeOK, the returned error is nil.
func Errorf(c Code, template string, args ...any) *Error {
	return Wrap(c, fmt.Errorf(template, args...))
}

// AsError uses errors.As to unwrap any error and look for a reRPC *Error.
func AsError(err error) (*Error, bool) {
	var re *Error
	ok := errors.As(err, &re)
	return re, ok
}

func (e *Error) Error() string {
	text := fmt.Sprintf("%v", e.err)
	if text == "" {
		return e.code.String()
	}
	return e.code.String() + ": " + text
}

// Unwrap implements errors.Wrapper, which allows errors.Is and errors.As
// access to the underlying error.
func (e *Error) Unwrap() error {
	return e.err
}

// Code returns the error's status code.
func (e *Error) Code() Code {
	if e == nil {
		return CodeOK
	}
	return e.code
}

// Details returns a deep copy of the error's details.
func (e *Error) Details() []*anypb.Any {
	if len(e.details) == 0 {
		return nil
	}
	ds := make([]*anypb.Any, len(e.details))
	for i, d := range e.details {
		ds[i] = proto.Clone(d).(*anypb.Any)
	}
	return ds
}

// AddDetail appends a message to the error's details.
func (e *Error) AddDetail(m proto.Message) error {
	if d, ok := m.(*anypb.Any); ok {
		e.details = append(e.details, proto.Clone(d).(*anypb.Any))
		return nil
	}
	detail, err := anypb.New(m)
	if err != nil {
		return fmt.Errorf("can't add message to error details: %w", err)
	}
	e.details = append(e.details, detail)
	return nil
}

// SetDetails overwrites the error's details.
func (e *Error) SetDetails(details ...proto.Message) error {
	e.details = make([]*anypb.Any, 0, len(details))
	for _, d := range details {
		if err := e.AddDetail(d); err != nil {
			return err
		}
	}
	return nil
}

// CodeOf returns the error's status code if it is or wraps a *rerpc.Error,
// CodeOK if the error is nil, and CodeUnknown otherwise.
func CodeOf(err error) Code {
	if err == nil {
		return CodeOK
	}
	if rerr, ok := AsError(err); ok {
		return rerr.Code()
	}
	return CodeUnknown
}
