package rerpc

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"google.golang.org/protobuf/proto"

	"github.com/akshayjshah/rerpc/internal/statuspb/v0"
)

var (
	// Always advertise that reRPC accepts gzip compression.
	acceptEncodingValue = strings.Join([]string{CompressionGzip, CompressionIdentity}, ", ")
	acceptPostValue     = strings.Join(
		[]string{TypeDefaultGRPC, TypeProtoGRPC, TypeJSON},
		", ",
	)
)

// A Handler is the server-side implementation of a single RPC defined by a
// protocol buffer service. It's the interface between the reRPC library and
// the code generated by the reRPC protoc plugin; most users won't ever need to
// deal with it directly.
//
// To see an example of how Handler is used in the generated code, see the
// internal/pingpb/v0 package.
type Handler struct {
	Implementation func(context.Context, proto.Message) (proto.Message, error)
}

// Serve executes the handler, much like the standard library's http.Handler.
// Unlike http.Handler, it requires a pointer to the protoc-generated request
// struct. See the internal/pingpb/v0 package for an example of how this code
// is used in reRPC's generated code.
//
// As long as the caller allocates a new request struct for each call, this
// method is safe to use concurrently.
func (h *Handler) Serve(w http.ResponseWriter, r *http.Request, msg proto.Message) {
	// To ensure that we can re-use connections, always consume and close the
	// request body.
	defer r.Body.Close()
	defer io.Copy(ioutil.Discard, r.Body)

	if r.Method != http.MethodPost {
		// grpc-go returns a 500 here, but interoperability with non-gRPC HTTP
		// clients is much better if we return a 405.
		w.Header().Set("Allow", http.MethodPost)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	ctype := r.Header.Get("Content-Type")
	if ctype != TypeDefaultGRPC && ctype != TypeProtoGRPC && ctype != TypeJSON {
		// grpc-go returns 500, but the spec recommends 415.
		// https://github.com/grpc/grpc/blob/master/doc/PROTOCOL-HTTP2.md#requests
		w.Header().Set("Accept-Post", acceptPostValue)
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	// We're always going to respond with the same content type as the request.
	w.Header().Set("Content-Type", ctype)
	if ctype == TypeJSON {
		h.serveJSON(w, r, msg)
	} else {
		h.serveGRPC(w, r, msg)
	}
}

func (h *Handler) serveJSON(w http.ResponseWriter, r *http.Request, msg proto.Message) {
	var closeWriter func()
	w, closeWriter = maybeGzipWriter(w, r)
	defer closeWriter()

	r, cancel, err := applyTimeout(r)
	if err != nil {
		// Errors here indicate that the client sent an invalid timeout header, so
		// the exact error is safe to send back.
		writeErrorJSON(w, wrap(CodeInvalidArgument, err))
		return
	}
	defer cancel()

	body, closeReader, err := maybeGzipReader(r)
	if err != nil {
		// TODO: observability
		writeErrorJSON(w, errorf(CodeUnknown, "can't read gzipped body"))
		return
	}
	defer closeReader()

	if err := unmarshalJSON(body, msg); err != nil {
		// TODO: observability
		writeErrorJSON(w, errorf(CodeInvalidArgument, "can't unmarshal JSON body"))
		return
	}

	res, implErr := h.Implementation(r.Context(), msg)
	if implErr != nil {
		// It's the user's job to sanitize the error string.
		writeErrorJSON(w, implErr)
		return
	}

	if err := marshalJSON(w, res); err != nil {
		// TODO: observability
		return
	}

	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
}

func (h *Handler) serveGRPC(w http.ResponseWriter, r *http.Request, msg proto.Message) {
	// We always send grpc-accept-encoding. Set it here so it's ready to go in
	// future error cases.
	w.Header().Set("Grpc-Accept-Encoding", acceptEncodingValue)
	// From here on, every gRPC response will have these trailers.
	w.Header().Add("Trailer", "Grpc-Status")
	w.Header().Add("Trailer", "Grpc-Message")
	w.Header().Add("Trailer", "Grpc-Status-Details-Bin")

	requestCompression := CompressionIdentity
	if me := r.Header.Get("Grpc-Encoding"); me != "" {
		switch me {
		case CompressionIdentity:
			requestCompression = CompressionIdentity
		case CompressionGzip:
			requestCompression = CompressionGzip
		default:
			// Per https://github.com/grpc/grpc/blob/master/doc/compression.md, we
			// should return CodeUnimplemented.
			writeErrorGRPC(w, errorf(CodeUnimplemented, "compression %q isn't supported", me))
			return
		}
	}

	// Follow https://github.com/grpc/grpc/blob/master/doc/compression.md.
	// (The grpc-go implementation doesn't read the "grpc-accept-encoding" header
	// and doesn't support compression method asymmetry.)
	responseCompression := requestCompression
	if mae := r.Header.Get("Grpc-Accept-Encoding"); mae != "" {
		for _, enc := range strings.FieldsFunc(mae, splitOnCommasAndSpaces) {
			switch enc {
			case CompressionIdentity:
				responseCompression = CompressionIdentity
				break
			case CompressionGzip:
				responseCompression = CompressionGzip
				break
			}
		}
	}
	w.Header().Set("Grpc-Encoding", responseCompression)

	r, cancel, err := applyTimeout(r)
	if err != nil {
		// Errors here indicate that the client sent an invalid timeout header, so
		// the exact error is safe to send back.
		writeErrorGRPC(w, wrap(CodeInvalidArgument, err))
		return
	}
	defer cancel()

	if err := unmarshalLPM(r.Body, msg, requestCompression); err != nil {
		// TODO: observability
		writeErrorGRPC(w, errorf(CodeInvalidArgument, "can't unmarshal protobuf request"))
		return
	}

	res, implErr := h.Implementation(r.Context(), msg)
	if implErr != nil {
		// It's the user's job to sanitize the error string.
		writeErrorGRPC(w, implErr)
		return
	}

	if err := marshalLPM(w, res, responseCompression); err != nil {
		// It's safe to write gRPC errors even after we've started writing the
		// body.
		// TODO: observability
		writeErrorGRPC(w, errorf(CodeUnknown, "can't marshal protobuf response"))
		return
	}

	writeErrorGRPC(w, nil)
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
}

func splitOnCommasAndSpaces(c rune) bool {
	return c == ',' || c == ' '
}

func writeErrorJSON(w http.ResponseWriter, err error) {
	s := statusFromError(err)
	bs, err := jsonpbMarshaler.Marshal(s)
	if err != nil {
		// TODO: observability
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"code": %d, "message": "error marshaling status with code %d"}`, CodeInternal, s.Code)
		return
	}
	w.WriteHeader(Code(s.Code).http())
	if _, err := w.Write(bs); err != nil {
		// TODO: observability
	}
}

func writeErrorGRPC(w http.ResponseWriter, err error) {
	if err == nil {
		w.Header().Set("Grpc-Status", strconv.Itoa(int(CodeOK)))
		w.Header().Set("Grpc-Message", "")
		w.Header().Set("Grpc-Status-Details-Bin", "")
		return
	}
	// gRPC errors are successes at the HTTP level and net/http automatically
	// sends a 200 if we don't set a status code. Leaving the HTTP status
	// implicit lets us use this function when we hit an error partway through
	// writing the body.
	s := statusFromError(err)
	code := strconv.Itoa(int(s.Code))
	// If we ever need to send more trailers, make sure to declare them in the headers
	// above.
	if bin, err := proto.Marshal(s); err != nil {
		w.Header().Set("Grpc-Status", strconv.Itoa(int(CodeInternal)))
		w.Header().Set("Grpc-Message", percentEncode("error marshaling protobuf status with code "+code))
	} else {
		w.Header().Set("Grpc-Status", code)
		w.Header().Set("Grpc-Message", percentEncode(s.Message))
		w.Header().Set("Grpc-Status-Details-Bin", encodeBinaryHeader(bin))
	}
}

func statusFromError(err error) *statuspb.Status {
	s := &statuspb.Status{
		Code:    int32(CodeUnknown),
		Message: err.Error(),
	}
	if re, ok := AsError(err); ok {
		s.Code = int32(re.Code())
		s.Details = re.Details()
		if e := re.Unwrap(); e != nil {
			s.Message = e.Error() // don't repeat code
		}
	}
	return s
}
