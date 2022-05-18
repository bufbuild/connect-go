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
	"bufio"
	"context"
	"errors"
	"fmt"
	"math"
	"net/http"
	"net/textproto"
	"runtime"
	"strconv"
	"strings"
	"time"

	statusv1 "github.com/bufbuild/connect-go/internal/gen/go/connectext/grpc/status/v1"
)

const (
	grpcHeaderCompression       = "Grpc-Encoding"
	grpcHeaderAcceptCompression = "Grpc-Accept-Encoding"
	grpcHeaderTimeout           = "Grpc-Timeout"
	grpcHeaderStatus            = "Grpc-Status"
	grpcHeaderMessage           = "Grpc-Message"
	grpcHeaderDetails           = "Grpc-Status-Details-Bin"

	grpcFlagEnvelopeTrailer = 0b10000000

	grpcTimeoutMaxHours = math.MaxInt64 / int64(time.Hour) // how many hours fit into a time.Duration?
	grpcMaxTimeoutChars = 8                                // from gRPC protocol

	grpcContentTypeDefault    = "application/grpc"
	grpcWebContentTypeDefault = "application/grpc-web"
	grpcContentTypePrefix     = grpcContentTypeDefault + "+"
	grpcWebContentTypePrefix  = grpcWebContentTypeDefault + "+"
)

var (
	grpcTimeoutUnits = []struct {
		size time.Duration
		char byte
	}{
		{time.Nanosecond, 'n'},
		{time.Microsecond, 'u'},
		{time.Millisecond, 'm'},
		{time.Second, 'S'},
		{time.Minute, 'M'},
		{time.Hour, 'H'},
	}
	grpcTimeoutUnitLookup = make(map[byte]time.Duration)
)

func init() {
	for _, pair := range grpcTimeoutUnits {
		grpcTimeoutUnitLookup[pair.char] = pair.size
	}
}

type protocolGRPC struct {
	web bool
}

// NewHandler implements protocol, so it must return an interface.
func (g *protocolGRPC) NewHandler(params *protocolHandlerParams) protocolHandler {
	bare, prefix := grpcContentTypeDefault, grpcContentTypePrefix
	if g.web {
		bare, prefix = grpcWebContentTypeDefault, grpcWebContentTypePrefix
	}
	contentTypes := make(map[string]struct{})
	for _, name := range params.Codecs.Names() {
		contentTypes[prefix+name] = struct{}{}
	}
	if params.Codecs.Get(codecNameProto) != nil {
		contentTypes[bare] = struct{}{}
	}
	return &grpcHandler{
		protocolHandlerParams: *params,
		web:                   g.web,
		accept:                contentTypes,
	}
}

// NewClient implements protocol, so it must return an interface.
func (g *protocolGRPC) NewClient(params *protocolClientParams) (protocolClient, error) {
	if err := validateRequestURL(params.URL); err != nil {
		return nil, err
	}
	return &grpcClient{
		protocolClientParams: *params,
		web:                  g.web,
	}, nil
}

type grpcHandler struct {
	protocolHandlerParams

	web    bool
	accept map[string]struct{}
}

func (g *grpcHandler) ContentTypes() map[string]struct{} {
	return g.accept
}

func (*grpcHandler) SetTimeout(request *http.Request) (context.Context, context.CancelFunc, error) {
	timeout, err := grpcParseTimeout(request.Header.Get(grpcHeaderTimeout))
	if err != nil && !errors.Is(err, errNoTimeout) {
		// Errors here indicate that the client sent an invalid timeout header, so
		// the error text is safe to send back.
		return nil, nil, NewError(CodeInvalidArgument, err)
	} else if err != nil {
		// err wraps errNoTimeout, nothing to do.
		return request.Context(), nil, nil
	}
	ctx, cancel := context.WithTimeout(request.Context(), timeout)
	return ctx, cancel, nil
}

func (g *grpcHandler) NewStream(
	responseWriter http.ResponseWriter,
	request *http.Request,
) (Sender, Receiver, error) {
	// We need to parse metadata before entering the interceptor stack, but we'd
	// like to report errors to the client in a format they understand (if
	// possible). We'll collect any such errors here; once we return, the Handler
	// will send the error to the client.
	requestCompression, responseCompression, failed := negotiateCompression(
		g.CompressionPools,
		request.Header.Get(grpcHeaderCompression),
		request.Header.Get(grpcHeaderAcceptCompression),
	)

	// We must write any remaining headers here:
	// (1) any writes to the stream will implicitly send the headers, so we
	// should get all of gRPC's required response headers ready.
	// (2) interceptors should be able to see these headers.
	//
	// Since we know that these header keys are already in canonical form, we can
	// skip the normalization in Header.Set.
	responseWriter.Header()[headerContentType] = []string{request.Header.Get(headerContentType)}
	responseWriter.Header()[grpcHeaderAcceptCompression] = []string{g.CompressionPools.CommaSeparatedNames()}
	if responseCompression != compressionIdentity {
		responseWriter.Header()[grpcHeaderCompression] = []string{responseCompression}
	}

	codecName := grpcCodecFromContentType(g.web, request.Header.Get(headerContentType))
	sender, receiver := wrapHandlerStreamWithCodedErrors(newGRPCHandlerStream(
		g.Spec,
		g.web,
		responseWriter,
		request,
		g.CompressMinBytes,
		g.Codecs.Get(codecName), // handler.go guarantees that this is not nil
		g.Codecs.Protobuf(),     // for errors
		g.CompressionPools.Get(requestCompression),
		g.CompressionPools.Get(responseCompression),
		g.BufferPool,
	))
	// We can't return failed as-is: a nil *Error is non-nil when returned as an
	// error interface.
	if failed != nil {
		// Negotiation failed, so we can't establish a stream. To make the
		// request's HTTP trailers visible to interceptors, we should try to read
		// the body to EOF.
		_ = discard(request.Body)
		return sender, receiver, failed
	}
	return sender, receiver, nil
}

type grpcClient struct {
	protocolClientParams

	web bool
}

func (g *grpcClient) WriteRequestHeader(_ StreamType, header http.Header) {
	// We know these header keys are in canonical form, so we can bypass all the
	// checks in Header.Set.
	header[headerUserAgent] = []string{grpcUserAgent()}
	header[headerContentType] = []string{grpcContentTypeFromCodecName(g.web, g.Codec.Name())}
	if g.CompressionName != "" && g.CompressionName != compressionIdentity {
		header[grpcHeaderCompression] = []string{g.CompressionName}
	}
	if acceptCompression := g.CompressionPools.CommaSeparatedNames(); acceptCompression != "" {
		header[grpcHeaderAcceptCompression] = []string{acceptCompression}
	}
	if !g.web {
		// No HTTP trailers in gRPC-Web.
		header["Te"] = []string{"trailers"}
	}
}

func (g *grpcClient) NewStream(
	ctx context.Context,
	spec Spec,
	header http.Header,
) (Sender, Receiver) {
	if deadline, ok := ctx.Deadline(); ok {
		if encodedDeadline, err := grpcEncodeTimeout(time.Until(deadline)); err == nil {
			// Tests verify that the error in encodeTimeout is unreachable, so we
			// should be safe without observability for the error case.
			header[grpcHeaderTimeout] = []string{encodedDeadline}
		}
	}
	duplexCall := newDuplexHTTPCall(
		ctx,
		g.HTTPClient,
		g.URL,
		spec,
		header,
	)
	sender := &grpcClientSender{
		spec:       spec,
		duplexCall: duplexCall,
		marshaler: grpcMarshaler{
			envelopeWriter: envelopeWriter{
				writer:           duplexCall,
				compressionPool:  g.CompressionPools.Get(g.CompressionName),
				codec:            g.Codec,
				compressMinBytes: g.CompressMinBytes,
				bufferPool:       g.BufferPool,
			},
		},
	}
	var receiver Receiver
	if g.web {
		webReceiver := &grpcWebClientReceiver{
			spec:             spec,
			bufferPool:       g.BufferPool,
			compressionPools: g.CompressionPools,
			codec:            g.Codec,
			protobuf:         g.Protobuf,
			header:           make(http.Header),
			trailer:          make(http.Header),
			duplexCall:       duplexCall,
		}
		receiver = webReceiver
		duplexCall.SetValidateResponse(webReceiver.validateResponse)
	} else {
		grpcReceiver := &grpcClientReceiver{
			spec:             spec,
			bufferPool:       g.BufferPool,
			compressionPools: g.CompressionPools,
			codec:            g.Codec,
			protobuf:         g.Protobuf,
			header:           make(http.Header),
			trailer:          make(http.Header),
			duplexCall:       duplexCall,
		}
		receiver = grpcReceiver
		duplexCall.SetValidateResponse(grpcReceiver.validateResponse)
	}
	return wrapClientStreamWithCodedErrors(sender, receiver)
}

// grpcClientSender works for both gRPC and gRPC-Web. From our perspective, the
// protocols differ only in how trailers are sent, and clients aren't allowed
// to send trailers.
type grpcClientSender struct {
	spec       Spec
	duplexCall *duplexHTTPCall
	marshaler  grpcMarshaler
}

func (s *grpcClientSender) Spec() Spec {
	return s.spec
}

func (s *grpcClientSender) Header() http.Header {
	return s.duplexCall.Header()
}

func (s *grpcClientSender) Trailer() (http.Header, bool) {
	return nil, false
}

func (s *grpcClientSender) Send(message any) error {
	// Don't return typed nils.
	if err := s.marshaler.Marshal(message); err != nil {
		return err
	}
	return nil
}

func (s *grpcClientSender) Close(_ error) error {
	return s.duplexCall.CloseWrite()
}

type grpcClientReceiver struct {
	spec             Spec
	bufferPool       *bufferPool
	compressionPools readOnlyCompressionPools
	codec            Codec
	protobuf         Codec // for errors
	header           http.Header
	trailer          http.Header
	duplexCall       *duplexHTTPCall
}

func (r *grpcClientReceiver) Spec() Spec {
	return r.spec
}

func (r *grpcClientReceiver) Header() http.Header {
	r.duplexCall.BlockUntilResponseReady()
	return r.header
}

func (r *grpcClientReceiver) Trailer() (http.Header, bool) {
	r.duplexCall.BlockUntilResponseReady()
	return r.trailer, true
}

func (r *grpcClientReceiver) Receive(message any) error {
	unmarshaler := r.unmarshaler()
	err := (&unmarshaler).Unmarshal(message)
	if err == nil {
		return nil
	}
	// See if the server sent an explicit error in the HTTP trailers. First, we
	// need to read the body to EOF.
	_ = discard(r.duplexCall)
	mergeHeaders(r.trailer, r.duplexCall.ResponseTrailer())
	if serverErr := grpcErrorFromTrailer(r.bufferPool, r.protobuf, r.trailer); serverErr != nil {
		// This is expected from a protocol perspective, but receiving trailers
		// means that we're _not_ getting a message. For users to realize that
		// the stream has ended, Receive must return an error.
		serverErr.meta = r.Header().Clone()
		mergeHeaders(serverErr.meta, r.trailer)
		r.duplexCall.SetError(serverErr)
		return serverErr
	}
	// There's no error in the trailers, so this was probably an error
	// converting the bytes to a message, an error reading from the network, or
	// just an EOF. We're going to return it to the user, but we also want to
	// setResponseError so Send errors out.
	r.duplexCall.SetError(err)
	return err
}

func (r *grpcClientReceiver) Close() error {
	return r.duplexCall.CloseRead()
}

func (r *grpcClientReceiver) unmarshaler() grpcUnmarshaler {
	// Blocks until response is ready.
	compression := r.duplexCall.ResponseHeader().Get(grpcHeaderCompression)
	// On the hot path, pass values up the stack.
	return grpcUnmarshaler{
		web: false,
		envelopeReader: envelopeReader{
			reader:          r.duplexCall,
			codec:           r.codec,
			compressionPool: r.compressionPools.Get(compression),
			bufferPool:      r.bufferPool,
		},
	}
}

// validateResponse is called by duplexHTTPCall in a separate goroutine.
func (r *grpcClientReceiver) validateResponse(response *http.Response) *Error {
	return grpcValidateResponse(
		response,
		r.header,
		r.trailer,
		r.compressionPools,
		r.bufferPool,
		r.protobuf,
	)
}

type grpcWebClientReceiver struct {
	spec             Spec
	bufferPool       *bufferPool
	compressionPools readOnlyCompressionPools
	codec            Codec
	protobuf         Codec // for errors
	header           http.Header
	trailer          http.Header
	duplexCall       *duplexHTTPCall
}

func (r *grpcWebClientReceiver) Spec() Spec {
	return r.spec
}

func (r *grpcWebClientReceiver) Header() http.Header {
	return r.header
}

func (r *grpcWebClientReceiver) Trailer() (http.Header, bool) {
	return r.trailer, true
}

func (r *grpcWebClientReceiver) Receive(message any) error {
	unmarshaler := r.unmarshaler()
	err := unmarshaler.Unmarshal(message)
	if err == nil {
		return nil
	}
	// See if the server sent an explicit error in the gRPC-Web trailers.
	mergeHeaders(r.trailer, unmarshaler.WebTrailer())
	if serverErr := grpcErrorFromTrailer(r.bufferPool, r.protobuf, r.trailer); serverErr != nil {
		// This is expected from a protocol perspective, but receiving a block of
		// trailers means that we're _not_ getting a standard message. For users to
		// realize that the stream has ended, Receive must return an error.
		serverErr.meta = r.Header().Clone()
		mergeHeaders(serverErr.meta, r.trailer)
		r.duplexCall.SetError(serverErr)
		return serverErr
	}
	r.duplexCall.SetError(err)
	return err
}

func (r *grpcWebClientReceiver) Close() error {
	return r.duplexCall.CloseRead()
}

func (r *grpcWebClientReceiver) unmarshaler() grpcUnmarshaler {
	// Blocks until response is ready.
	compression := r.duplexCall.ResponseHeader().Get(grpcHeaderCompression)
	// On the hot path, pass values up the stack.
	return grpcUnmarshaler{
		web: true,
		envelopeReader: envelopeReader{
			reader:          r.duplexCall,
			codec:           r.codec,
			compressionPool: r.compressionPools.Get(compression),
			bufferPool:      r.bufferPool,
		},
	}
}

// validateResponse is called by duplexHTTPCall in a separate goroutine.
func (r *grpcWebClientReceiver) validateResponse(response *http.Response) *Error {
	return grpcValidateResponse(
		response,
		r.header,
		r.trailer,
		r.compressionPools,
		r.bufferPool,
		r.protobuf,
	)
}

type grpcHandlerSender struct {
	spec        Spec
	web         bool
	marshaler   grpcMarshaler
	protobuf    Codec // for errors
	writer      http.ResponseWriter
	header      http.Header
	trailer     http.Header
	wroteToBody bool
	bufferPool  *bufferPool
}

func (hs *grpcHandlerSender) Send(message any) error {
	defer hs.flush()
	if !hs.wroteToBody {
		mergeHeaders(hs.writer.Header(), hs.header)
		hs.wroteToBody = true
	}
	if err := hs.marshaler.Marshal(message); err != nil {
		return err // already coded
	}
	// don't return typed nils
	return nil
}

func (hs *grpcHandlerSender) Close(err error) error {
	defer hs.flush()
	// If we haven't written the headers yet, do so.
	if !hs.wroteToBody {
		mergeHeaders(hs.writer.Header(), hs.header)
	}
	// gRPC always sends the error's code, message, details, and metadata as
	// trailers. Future protocols may not do this, though, so we don't want to
	// mutate the trailers map that the user sees.
	mergedTrailers := make(http.Header, len(hs.trailer)+2) // always make space for status & message
	mergeHeaders(mergedTrailers, hs.trailer)
	grpcErrorToTrailer(hs.bufferPool, mergedTrailers, hs.protobuf, err)
	if hs.web && !hs.wroteToBody {
		// We're using gRPC-Web and we haven't yet written to the body. Since we're
		// not sending any response messages, the gRPC specification calls this a
		// "trailers-only" response. Under those circumstances, the gRPC-Web spec
		// says that implementations _may_ send trailing metadata as HTTP headers
		// instead. The gRPC-Web spec is explicitly a description of the reference
		// implementations rather than a proper specification, so we're going to
		// emulate Envoy and put the trailing metadata in the HTTP headers.
		mergeHeaders(hs.writer.Header(), mergedTrailers)
		return nil
	}
	if hs.web {
		// We're using gRPC-Web and we've already sent the headers, so we write
		// trailing metadata to the HTTP body.
		if err := hs.marshaler.MarshalWebTrailers(mergedTrailers); err != nil {
			return err
		}
		// Don't return typed nils.
		return nil
	}
	// We're using standard gRPC. Even if we haven't written to the body and
	// we're sending a "trailers-only" response, we must send trailing metadata
	// as HTTP trailers. (If we had frame-level control of the HTTP/2 layer, we
	// could send a single HEADER frame and no DATA frames, but net/http doesn't
	// expose APIs that low-level.) In net/http's ResponseWriter API, we do that
	// by writing to the headers map with a special prefix. This is purely an
	// implementation detail, so we should hide it and _not_ mutate the
	// user-visible headers.
	//
	// Note that this is _very_ finicky, and it's impossible to test with a
	// net/http client. Breaking this logic breaks Envoy's gRPC-Web translation.
	for key, values := range mergedTrailers {
		for _, value := range values {
			hs.writer.Header().Add(http.TrailerPrefix+key, value)
		}
	}
	return nil
}

func (hs *grpcHandlerSender) Spec() Spec {
	return hs.spec
}

func (hs *grpcHandlerSender) Header() http.Header {
	return hs.header
}

func (hs *grpcHandlerSender) Trailer() (http.Header, bool) {
	return hs.trailer, true
}

func (hs *grpcHandlerSender) flush() {
	if f, ok := hs.writer.(http.Flusher); ok {
		f.Flush()
	}
}

type grpcHandlerReceiver struct {
	spec        Spec
	unmarshaler grpcUnmarshaler
	request     *http.Request
}

func (hr *grpcHandlerReceiver) Receive(message any) error {
	if err := hr.unmarshaler.Unmarshal(message); err != nil {
		if errors.Is(err, errSpecialEnvelope) {
			if hr.request.Trailer == nil {
				hr.request.Trailer = hr.unmarshaler.WebTrailer()
			} else {
				mergeHeaders(hr.request.Trailer, hr.unmarshaler.WebTrailer())
			}
		}
		return err // already coded
	}
	// don't return typed nils
	return nil
}

func (hr *grpcHandlerReceiver) Close() error {
	// We don't want to copy unread portions of the body to /dev/null here: if
	// the client hasn't closed the request body, we'll block until the server
	// timeout kicks in. This could happen because the client is malicious, but
	// a well-intentioned client may just not expect the server to be returning
	// an error for a streaming RPC. Better to accept that we can't always reuse
	// TCP connections.
	return hr.request.Body.Close()
}

func (hr *grpcHandlerReceiver) Spec() Spec {
	return hr.spec
}

func (hr *grpcHandlerReceiver) Header() http.Header {
	return hr.request.Header
}

func (hr *grpcHandlerReceiver) Trailer() (http.Header, bool) {
	return nil, false
}

type grpcMarshaler struct {
	envelopeWriter
}

func (m *grpcMarshaler) MarshalWebTrailers(trailer http.Header) *Error {
	raw := m.envelopeWriter.bufferPool.Get()
	defer m.envelopeWriter.bufferPool.Put(raw)
	if err := trailer.Write(raw); err != nil {
		return errorf(CodeInternal, "format trailers: %w", err)
	}
	return m.Write(&envelope{
		Data:  raw,
		Flags: grpcFlagEnvelopeTrailer,
	})
}

type grpcUnmarshaler struct {
	envelopeReader envelopeReader
	web            bool
	webTrailer     http.Header
}

func (u *grpcUnmarshaler) Unmarshal(message any) *Error {
	err := u.envelopeReader.Unmarshal(message)
	if err == nil {
		return nil
	}
	if !errors.Is(err, errSpecialEnvelope) {
		return err
	}
	env := u.envelopeReader.last
	if !u.web || !env.IsSet(grpcFlagEnvelopeTrailer) {
		return errorf(CodeUnknown, "protocol error: invalid envelope flags %d", env.Flags)
	}

	// Per the gRPC-Web specification, trailers should be encoded as an HTTP/1
	// headers block _without_ the terminating newline. To make the headers
	// parseable by net/textproto, we need to add the newline.
	if err := env.Data.WriteByte('\n'); err != nil {
		return errorf(CodeUnknown, "unmarshal web trailers: %w", err)
	}
	bufferedReader := bufio.NewReader(env.Data)
	mimeReader := textproto.NewReader(bufferedReader)
	mimeHeader, mimeErr := mimeReader.ReadMIMEHeader()
	if mimeErr != nil {
		return errorf(
			CodeInvalidArgument,
			"gRPC-Web protocol error: received invalid trailers: %w",
			mimeErr,
		)
	}
	u.webTrailer = http.Header(mimeHeader)
	return errSpecialEnvelope
}

func (u *grpcUnmarshaler) WebTrailer() http.Header {
	return u.webTrailer
}

func newGRPCHandlerStream(
	spec Spec,
	web bool,
	responseWriter http.ResponseWriter,
	request *http.Request,
	compressMinBytes int,
	codec Codec,
	protobuf Codec, // for errors
	requestCompressionPools *compressionPool,
	responseCompressionPools *compressionPool,
	bufferPool *bufferPool,
) (*grpcHandlerSender, *grpcHandlerReceiver) {
	sender := &grpcHandlerSender{
		spec: spec,
		web:  web,
		marshaler: grpcMarshaler{
			envelopeWriter: envelopeWriter{
				writer:           responseWriter,
				compressionPool:  responseCompressionPools,
				codec:            codec,
				compressMinBytes: compressMinBytes,
				bufferPool:       bufferPool,
			},
		},
		protobuf:   protobuf,
		writer:     responseWriter,
		header:     make(http.Header),
		trailer:    make(http.Header),
		bufferPool: bufferPool,
	}
	receiver := &grpcHandlerReceiver{
		spec: spec,
		unmarshaler: grpcUnmarshaler{
			envelopeReader: envelopeReader{
				reader:          request.Body,
				codec:           codec,
				compressionPool: requestCompressionPools,
				bufferPool:      bufferPool,
			},
			web: web,
		},
		request: request,
	}
	return sender, receiver
}

func grpcValidateResponse(
	response *http.Response,
	header, trailer http.Header,
	availableCompressors readOnlyCompressionPools,
	bufferPool *bufferPool,
	protobuf Codec,
) *Error {
	if response.StatusCode != http.StatusOK {
		return errorf(grpcHTTPToCode(response.StatusCode), "HTTP status %v", response.Status)
	}
	if compression := response.Header.Get(grpcHeaderCompression); compression != "" &&
		compression != compressionIdentity &&
		!availableCompressors.Contains(compression) {
		// Per https://github.com/grpc/grpc/blob/master/doc/compression.md, we
		// should return CodeInternal and specify acceptable compression(s) (in
		// addition to setting the Grpc-Accept-Encoding header).
		return errorf(
			CodeInternal,
			"unknown encoding %q: accepted grpc-encoding values are %v",
			compression,
			availableCompressors.CommaSeparatedNames(),
		)
	}
	// When there's no body, gRPC and gRPC-Web servers may send error information
	// in the HTTP headers.
	if err := grpcErrorFromTrailer(bufferPool, protobuf, response.Header); err != nil {
		// Per the specification, only the HTTP status code and Content-Type should
		// be treated as headers. The rest should be treated as trailing metadata.
		if contentType := response.Header.Get(headerContentType); contentType != "" {
			header.Set(headerContentType, contentType)
		}
		mergeHeaders(trailer, response.Header)
		trailer.Del(headerContentType)
		// If we get some actual HTTP trailers, treat those as trailing metadata too.
		_ = discard(response.Body)
		mergeHeaders(trailer, response.Trailer)
		// Merge HTTP headers and trailers into the error metadata.
		err.meta = response.Header.Clone()
		mergeHeaders(err.meta, response.Trailer)
		return err
	}
	// The response is valid, so we should expose the headers.
	mergeHeaders(header, response.Header)
	return nil
}

func grpcHTTPToCode(httpCode int) Code {
	// https://github.com/grpc/grpc/blob/master/doc/http-grpc-status-mapping.md
	// Note that this is not just the inverse of the gRPC-to-HTTP mapping.
	switch httpCode {
	case 400:
		return CodeInternal
	case 401:
		return CodeUnauthenticated
	case 403:
		return CodePermissionDenied
	case 404:
		return CodeUnimplemented
	case 429:
		return CodeUnavailable
	case 502, 503, 504:
		return CodeUnavailable
	default:
		return CodeUnknown
	}
}

// The gRPC wire protocol specifies that errors should be serialized using the
// binary Protobuf format, even if the messages in the request/response stream
// use a different codec. Consequently, this function needs a Protobuf codec to
// unmarshal error information in the headers.
func grpcErrorFromTrailer(bufferPool *bufferPool, protobuf Codec, trailer http.Header) *Error {
	codeHeader := trailer.Get(grpcHeaderStatus)
	if codeHeader == "" || codeHeader == "0" {
		return nil
	}

	code, err := strconv.ParseUint(codeHeader, 10 /* base */, 32 /* bitsize */)
	if err != nil {
		return errorf(CodeUnknown, "gRPC protocol error: got invalid error code %q", codeHeader)
	}
	message := percentDecode(bufferPool, trailer.Get(grpcHeaderMessage))
	retErr := NewError(Code(code), errors.New(message))

	detailsBinaryEncoded := trailer.Get(grpcHeaderDetails)
	if len(detailsBinaryEncoded) > 0 {
		detailsBinary, err := DecodeBinaryHeader(detailsBinaryEncoded)
		if err != nil {
			return errorf(CodeUnknown, "server returned invalid grpc-status-details-bin trailer: %w", err)
		}
		var status statusv1.Status
		if err := protobuf.Unmarshal(detailsBinary, &status); err != nil {
			return errorf(CodeUnknown, "server returned invalid protobuf for error details: %w", err)
		}
		for _, d := range status.Details {
			retErr.details = append(retErr.details, d)
		}
		// Prefer the Protobuf-encoded data to the headers (grpc-go does this too).
		retErr.code = Code(status.Code)
		retErr.err = errors.New(status.Message)
	}

	return retErr
}

func grpcParseTimeout(timeout string) (time.Duration, error) {
	if timeout == "" {
		return 0, errNoTimeout
	}
	unit, ok := grpcTimeoutUnitLookup[timeout[len(timeout)-1]]
	if !ok {
		return 0, fmt.Errorf("gRPC protocol error: timeout %q has invalid unit", timeout)
	}
	num, err := strconv.ParseInt(timeout[:len(timeout)-1], 10 /* base */, 64 /* bitsize */)
	if err != nil || num < 0 {
		return 0, fmt.Errorf("gRPC protocol error: invalid timeout %q", timeout)
	}
	if num > 99999999 { // timeout must be ASCII string of at most 8 digits
		return 0, fmt.Errorf("gRPC protocol error: timeout %q is too long", timeout)
	}
	if unit == time.Hour && num > grpcTimeoutMaxHours {
		// Timeout is effectively unbounded, so ignore it. The grpc-go
		// implementation does the same thing.
		return 0, errNoTimeout
	}
	return time.Duration(num) * unit, nil
}

func grpcEncodeTimeout(timeout time.Duration) (string, error) {
	if timeout <= 0 {
		return "0n", nil
	}
	for _, pair := range grpcTimeoutUnits {
		digits := strconv.FormatInt(int64(timeout/pair.size), 10 /* base */)
		if len(digits) < grpcMaxTimeoutChars {
			return digits + string(pair.char), nil
		}
	}
	// The max time.Duration is smaller than the maximum expressible gRPC
	// timeout, so we shouldn't ever reach this case.
	return "", errNoTimeout
}

// grpcUserAgent follows
// https://github.com/grpc/grpc/blob/master/doc/PROTOCOL-HTTP2.md#user-agents:
//
//   While the protocol does not require a user-agent to function it is recommended
//   that clients provide a structured user-agent string that provides a basic
//   description of the calling library, version & platform to facilitate issue diagnosis
//   in heterogeneous environments. The following structure is recommended to library developers:
//
//   User-Agent → "grpc-" Language ?("-" Variant) "/" Version ?( " ("  *(AdditionalProperty ";") ")" )
func grpcUserAgent() string {
	return fmt.Sprintf("grpc-go-connect/%s (%s)", Version, runtime.Version())
}

func grpcCodecFromContentType(web bool, contentType string) string {
	if (!web && contentType == grpcContentTypeDefault) || (web && contentType == grpcWebContentTypeDefault) {
		// implicitly protobuf
		return codecNameProto
	}
	prefix := grpcContentTypePrefix
	if web {
		prefix = grpcWebContentTypePrefix
	}
	return strings.TrimPrefix(contentType, prefix)
}

func grpcContentTypeFromCodecName(web bool, name string) string {
	if web {
		return grpcWebContentTypePrefix + name
	}
	return grpcContentTypePrefix + name
}

func grpcErrorToTrailer(bufferPool *bufferPool, trailer http.Header, protobuf Codec, err error) {
	if err == nil {
		trailer.Set(grpcHeaderStatus, "0") // zero is the gRPC OK status
		trailer.Set(grpcHeaderMessage, "")
		return
	}
	status, statusErr := grpcStatusFromError(err)
	if statusErr != nil {
		trailer.Set(
			grpcHeaderStatus,
			strconv.FormatInt(int64(CodeInternal), 10 /* base */),
		)
		trailer.Set(grpcHeaderMessage, statusErr.Error())
		return
	}
	code := strconv.Itoa(int(status.Code))
	bin, binErr := protobuf.Marshal(status)
	if binErr != nil {
		trailer.Set(
			grpcHeaderStatus,
			strconv.FormatInt(int64(CodeInternal), 10 /* base */),
		)
		trailer.Set(
			grpcHeaderMessage,
			fmt.Sprintf("marshal protobuf status: %v", binErr),
		)
		return
	}
	if connectErr, ok := asError(err); ok {
		mergeHeaders(trailer, connectErr.meta)
	}
	trailer.Set(grpcHeaderStatus, code)
	trailer.Set(grpcHeaderMessage, percentEncode(bufferPool, status.Message))
	trailer.Set(grpcHeaderDetails, EncodeBinaryHeader(bin))
}

func grpcStatusFromError(err error) (*statusv1.Status, error) {
	status := &statusv1.Status{
		Code:    int32(CodeUnknown),
		Message: err.Error(),
	}
	if connectErr, ok := asError(err); ok {
		status.Code = int32(connectErr.Code())
		status.Message = connectErr.Message()
		details, err := connectErr.detailsAsAny()
		if err != nil {
			return nil, err
		}
		status.Details = details
	}
	return status, nil
}
