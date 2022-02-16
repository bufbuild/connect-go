package connect

import (
	"fmt"
	"io"
	"net/http"
	"runtime"
	"strconv"
	"strings"

	"github.com/bufconnect/connect/codec"
	"github.com/bufconnect/connect/codec/protobuf"
	statuspb "github.com/bufconnect/connect/internal/gen/proto/go/grpc/status/v1"
	"google.golang.org/protobuf/types/known/anypb"
)

const (
	typeDefaultGRPC       = "application/grpc"
	typeWebGRPC           = "application/grpc-web"
	typeDefaultGRPCPrefix = typeDefaultGRPC + "+"
	typeWebGRPCPrefix     = typeWebGRPC + "+"
	grpcNameProto         = "proto" // gRPC protocols use "proto" instead of "protobuf"
)

// Follows https://github.com/grpc/grpc/blob/master/doc/PROTOCOL-HTTP2.md#user-agents
var userAgent = []string{fmt.Sprintf("grpc-go-connect/%s (%s)", Version, runtime.Version())}

func isCommaOrSpace(c rune) bool {
	return c == ',' || c == ' '
}

func acceptPostValue(web bool, codecs readOnlyCodecs) string {
	bare, prefix := typeDefaultGRPC, typeDefaultGRPCPrefix
	if web {
		bare, prefix = typeWebGRPC, typeWebGRPCPrefix
	}
	names := codecs.Names()
	for i, name := range names {
		if name == protobuf.NameBinary {
			name = grpcNameProto
		}
		names[i] = prefix + name
	}
	if codecs.Get(protobuf.NameBinary) != nil {
		names = append(names, bare)
	}
	return strings.Join(names, ",")
}

func codecFromContentType(web bool, contentType string) string {
	if (!web && contentType == typeDefaultGRPC) || (web && contentType == typeWebGRPC) {
		// implicitly protobuf
		return protobuf.NameBinary
	}
	prefix := typeDefaultGRPCPrefix
	if web {
		prefix = typeWebGRPCPrefix
	}
	if !strings.HasPrefix(contentType, prefix) {
		return ""
	}
	name := strings.TrimPrefix(contentType, prefix)
	if name == grpcNameProto {
		// normalize to our "protobuf"
		return protobuf.NameBinary
	}
	return name
}

func contentTypeFromCodecName(web bool, name string) string {
	// translate back to gRPC's "proto"
	if name == protobuf.NameBinary {
		name = grpcNameProto
	}
	if web {
		return typeWebGRPCPrefix + name
	}
	return typeDefaultGRPCPrefix + name
}

func grpcErrorToTrailer(trailer http.Header, protobuf codec.Codec, err error) error {
	if CodeOf(err) == CodeOK { // safe for nil errors
		trailer.Set("Grpc-Status", strconv.Itoa(int(CodeOK)))
		trailer.Set("Grpc-Message", "")
		trailer.Set("Grpc-Status-Details-Bin", "")
		return nil
	}
	status, statusErr := statusFromError(err)
	if statusErr != nil {
		return statusErr
	}
	code := strconv.Itoa(int(status.Code))
	bin, err := protobuf.Marshal(status)
	if err != nil {
		trailer.Set("Grpc-Status", strconv.Itoa(int(CodeInternal)))
		trailer.Set("Grpc-Message", percentEncode("error marshaling protobuf status with code "+code))
		return Errorf(CodeInternal, "couldn't marshal protobuf status: %w", err)
	}
	trailer.Set("Grpc-Status", code)
	trailer.Set("Grpc-Message", percentEncode(status.Message))
	trailer.Set("Grpc-Status-Details-Bin", EncodeBinaryHeader(bin))
	return nil
}

func statusFromError(err error) (*statuspb.Status, *Error) {
	status := &statuspb.Status{
		Code:    int32(CodeUnknown),
		Message: err.Error(),
	}
	if connectErr, ok := AsError(err); ok {
		status.Code = int32(connectErr.Code())
		for _, detail := range connectErr.details {
			// If the detail is already a protobuf Any, we're golden.
			if anyProtoDetail, ok := detail.(*anypb.Any); ok {
				status.Details = append(status.Details, anyProtoDetail)
				continue
			}
			// Otherwise, we convert it to an Any.
			// TODO: Should we also attempt to delegate this to the detail by
			// attempting an upcast to interface{ AsAny() *anypb.Any }?
			anyProtoDetail, err := anypb.New(detail)
			if err != nil {
				return nil, Errorf(
					CodeInternal,
					"can't create an *anypb.Any from %v (type %T): %v",
					detail, detail, err,
				)
			}
			status.Details = append(status.Details, anyProtoDetail)
		}
		if underlyingErr := connectErr.Unwrap(); underlyingErr != nil {
			status.Message = underlyingErr.Error() // don't repeat code
		}
	}
	return status, nil
}

func discard(r io.Reader) {
	if lr, ok := r.(*io.LimitedReader); ok {
		io.Copy(io.Discard, lr)
		return
	}
	// We don't want to get stuck throwing data away forever, so limit how much
	// we're willing to do here: at most, we'll copy 4 MiB.
	lr := &io.LimitedReader{R: r, N: 1024 * 1024 * 4}
	io.Copy(io.Discard, lr)
}
