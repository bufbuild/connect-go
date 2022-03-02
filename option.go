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

package connect

import (
	"compress/gzip"
	"io/ioutil"
	"strings"
)

// Option implements both ClientOption and HandlerOption, so it can be applied
// both client-side and server-side.
type Option interface {
	ClientOption
	HandlerOption
}

type compressMinBytesOption struct {
	min int
}

// WithCompressMinBytes sets a minimum size threshold for compression:
// regardless of compressor configuration, messages smaller than the configured
// minimum are sent uncompressed.
//
// The default minimum is zero. Setting a minimum compression threshold may
// improve overall performance, because the CPU cost of compressing very small
// messages usually isn't worth the small reduction in network I/O.
func WithCompressMinBytes(min int) Option {
	return &compressMinBytesOption{min: min}
}

func (o *compressMinBytesOption) applyToClient(config *clientConfiguration) {
	config.CompressMinBytes = o.min
}

func (o *compressMinBytesOption) applyToHandler(config *handlerConfiguration) {
	config.CompressMinBytes = o.min
}

type replaceProcedurePrefixOption struct {
	prefix      string
	replacement string
}

// WithReplaceProcedurePrefix changes the URL used to call a procedure.
// Typically, generated code sets the procedure name: for example, a protobuf
// procedure's name and URL is composed from the fully-qualified protobuf
// package name, the service name, and the method name. This option replaces a
// prefix of the procedure name with another static string. Using this option
// is usually a bad idea, but it's occasionally necessary to prevent protobuf
// package collisions. (For example, connect uses this option to serve the
// health and reflection APIs without generating runtime conflicts with
// grpc-go.)
//
// WithReplaceProcedurePrefix doesn't change the data exposed by the reflection
// API. To prevent inconsistencies between the reflection data and the actual
// service URL, using this option disables reflection for the modified service
// (though other services can still be introspected).
func WithReplaceProcedurePrefix(prefix, replacement string) Option {
	return &replaceProcedurePrefixOption{
		prefix:      prefix,
		replacement: replacement,
	}
}

func (o *replaceProcedurePrefixOption) applyToClient(config *clientConfiguration) {
	config.Procedure = o.transform(config.Procedure)
}

func (o *replaceProcedurePrefixOption) applyToHandler(config *handlerConfiguration) {
	config.Procedure = o.transform(config.Procedure)
	config.RegistrationName = "" // disable reflection
}

func (o *replaceProcedurePrefixOption) transform(name string) string {
	if !strings.HasPrefix(name, o.prefix) {
		return name
	}
	return o.replacement + strings.TrimPrefix(name, o.prefix)
}

type readMaxBytesOption struct {
	Max int64
}

// WithReadMaxBytes limits the performance impact of pathologically large
// messages sent by the other party. For handlers, WithReadMaxBytes limits the size
// of message that the client can send. For clients, WithReadMaxBytes limits the
// size of message that the server can respond with. Limits are applied before
// decompression and apply to each protobuf message, not to the stream as a
// whole.
//
// Setting WithReadMaxBytes to zero allows any message size. Both clients and
// handlers default to allowing any request size.
func WithReadMaxBytes(n int64) Option {
	return &readMaxBytesOption{n}
}

func (o *readMaxBytesOption) applyToClient(config *clientConfiguration) {
	config.MaxResponseBytes = o.Max
}

func (o *readMaxBytesOption) applyToHandler(config *handlerConfiguration) {
	config.MaxRequestBytes = o.Max
}

type codecOption struct {
	Codec Codec
}

// WithCodec registers a serialization method with a client or handler.
// Registering a codec with an empty name is a no-op.
//
// Typically, generated code automatically supplies this option with the
// appropriate codec(s). For example, handlers generated from protobuf schemas
// using protoc-gen-connect-go automatically register binary and JSON codecs.
// Users with more specialized needs may override the default codecs by
// registering a new codec under the same name.
//
// Handlers may have multiple codecs registered, and use whichever the client
// chooses. Clients may only have a single codec.
func WithCodec(c Codec) Option {
	return &codecOption{Codec: c}
}

// WithProtoBinaryCodec registers a binary protocol buffer codec that uses
// google.golang.org/protobuf/proto.
//
// Handlers and clients generated by protoc-gen-connect-go have
// WithProtoBinaryCodec applied by default. To replace the default binary protobuf
// codec (with vtprotobuf, for example), apply WithCodec with a Codec whose
// name is "proto".
func WithProtoBinaryCodec() Option {
	return WithCodec(&protoBinaryCodec{})
}

// WithProtoJSONCodec registers a codec that serializes protocol buffer
// messages as JSON. It uses the standard protobuf JSON mapping as implemented
// by google.golang.org/protobuf/encoding/protojson: fields are named using
// lowerCamelCase, zero values are omitted, missing required fields are errors,
// enums are emitted as strings, etc.
//
// Handlers generated by protoc-gen-connect-go have WithProtoJSONCodec
// applied by default.
func WithProtoJSONCodec() Option {
	return WithCodec(&protoJSONCodec{})
}

func (o *codecOption) applyToClient(config *clientConfiguration) {
	if o.Codec == nil || o.Codec.Name() == "" {
		return
	}
	config.Codec = o.Codec
}

func (o *codecOption) applyToHandler(config *handlerConfiguration) {
	if o.Codec == nil || o.Codec.Name() == "" {
		return
	}
	config.Codecs[o.Codec.Name()] = o.Codec
}

type compressionOption struct {
	Name            string
	CompressionPool compressionPool
}

// WithCompression configures client and server compression strategies. The
// Compressors and Decompressors produced by the supplied constructors must use
// the same algorithm.
//
// For handlers, WithCompression registers a compression algorithm. Clients may
// send messages compressed with that algorithm and/or request compressed
// responses.
//
// For clients, WithCompression serves two purposes. First, the client
// asks servers to compress responses using any of the registered algorithms.
// (gRPC's compression negotiation is complex, but most of Google's gRPC server
// implementations won't compress responses unless the request is compressed.)
// Second, it makes all the registered algorithms available for use with
// WithRequestCompression. Note that actually compressing requests requires
// using both WithCompression and WithRequestCompression.
//
// Calling WithCompression with an empty name or nil constructors is a no-op.
func WithCompression[D Decompressor, C Compressor](
	name string,
	newDecompressor func() D,
	newCompressor func() C,
) Option {
	return &compressionOption{
		Name:            name,
		CompressionPool: newCompressionPool(newDecompressor, newCompressor),
	}
}

// WithGzip registers a gzip compressor backed by the standard library's gzip
// package with the default compression level.
//
// Handlers with this option applied accept gzipped requests and can send
// gzipped responses. Clients with this option applied request gzipped
// responses, but don't automatically send gzipped requests (since the server
// may not support them). Use WithGzipRequests to gzip requests.
//
// Handlers and clients generated by protoc-gen-connect-go apply WithGzip by
// default.
func WithGzip() Option {
	return WithCompression(
		compressionGzip,
		func() *gzip.Reader { return &gzip.Reader{} },
		func() *gzip.Writer { return gzip.NewWriter(ioutil.Discard) },
	)
}

func (o *compressionOption) applyToClient(config *clientConfiguration) {
	o.apply(config.CompressionPools)
}

func (o *compressionOption) applyToHandler(config *handlerConfiguration) {
	o.apply(config.CompressionPools)
}

func (o *compressionOption) apply(m map[string]compressionPool) {
	if o.Name == "" || o.CompressionPool == nil {
		return
	}
	m[o.Name] = o.CompressionPool
}

type interceptOption struct {
	interceptors []Interceptor
}

// WithInterceptors configures a client or handler's interceptor stack. Repeated
// WithInterceptors options are applied in order, so
//
//   WithInterceptors(A) + WithInterceptors(B, C) == WithInterceptors(A, B, C)
//
// Unary interceptors compose like an onion. The first interceptor provided is
// the outermost layer of the onion: it acts first on the context and request,
// and last on the response and error.
//
// Stream interceptors also behave like an onion: the first interceptor
// provided is the first to wrap the context and is the outermost wrapper for
// the (Sender, Receiver) pair. It's the first to see sent messages and the
// last to see received messages.
//
// Applied to client and handler, WithInterceptors(A, B, ..., Y, Z) produces:
//
//        client.Send()     client.Receive()
//              |                 ^
//              v                 |
//           A ---               --- A
//           B ---               --- B
//             ...               ...
//           Y ---               --- Y
//           Z ---               --- Z
//              |                 ^
//              v                 |
//           network            network
//              |                 ^
//              v                 |
//           A ---               --- A
//           B ---               --- B
//             ...               ...
//           Y ---               --- Y
//           Z ---               --- Z
//              |                 ^
//              v                 |
//       handler.Receive() handler.Send()
//              |                 ^
//              |                 |
//              -> handler logic --
//
// Note that in clients, the Sender handles the request message(s) and the
// Receiver handles the response message(s). For handlers, it's the reverse.
// Depending on your interceptor's logic, you may need to wrap one side of the
// stream on the clients and the other side on handlers.
func WithInterceptors(interceptors ...Interceptor) Option {
	return &interceptOption{interceptors}
}

func (o *interceptOption) applyToClient(config *clientConfiguration) {
	config.Interceptor = o.chainWith(config.Interceptor)
}

func (o *interceptOption) applyToHandler(config *handlerConfiguration) {
	config.Interceptor = o.chainWith(config.Interceptor)
}

func (o *interceptOption) chainWith(current Interceptor) Interceptor {
	if len(o.interceptors) == 0 {
		return current
	}
	if current == nil && len(o.interceptors) == 1 {
		return o.interceptors[0]
	}
	if current == nil && len(o.interceptors) > 1 {
		return newChain(o.interceptors)
	}
	return newChain(append([]Interceptor{current}, o.interceptors...))
}
