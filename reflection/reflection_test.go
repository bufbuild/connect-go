package reflection_test

import (
	"bytes"
	"context"
	"net/http/httptest"
	"testing"

	"github.com/bufconnect/connect"
	"github.com/bufconnect/connect/codec/protobuf"
	"github.com/bufconnect/connect/health"
	"github.com/bufconnect/connect/internal/assert"
	pingrpc "github.com/bufconnect/connect/internal/gen/proto/go-connect/connect/ping/v1test"
	pingpb "github.com/bufconnect/connect/internal/gen/proto/go/connect/ping/v1test"
	reflectionpb "github.com/bufconnect/connect/internal/gen/proto/go/grpc/reflection/v1alpha"
	"github.com/bufconnect/connect/reflection"
)

func TestReflection(t *testing.T) {
	reg := connect.NewRegistrar()
	mux, err := connect.NewServeMux(
		pingrpc.WithPingServiceHandler(
			pingrpc.UnimplementedPingServiceHandler{},
			reg,
		),
		health.WithHandler(health.NewChecker(reg)),
		reflection.WithHandler(reg),
	)
	assert.Nil(t, err, "mux construction error")

	server := httptest.NewUnstartedServer(mux)
	server.EnableHTTP2 = true
	server.StartTLS()
	defer server.Close()

	pingRequestFQN := string((&pingpb.PingRequest{}).ProtoReflect().Descriptor().FullName())
	assert.Equal(t, reg.Services(), []string{
		"connect.ping.v1test.PingService",
	}, "services registered in memory")

	detailed, err := connect.NewClientFunc[
		reflectionpb.ServerReflectionRequest,
		reflectionpb.ServerReflectionResponse,
	](
		server.Client(),
		server.URL,
		"grpc.reflection.v1alpha.ServerReflection/ServerReflectionInfo",
		connect.WithGRPC(true),
		connect.WithCodec(protobuf.NameBinary, protobuf.NewBinary()),
	)
	assert.Nil(t, err, "client construction error")
	call := func(req *reflectionpb.ServerReflectionRequest) (*reflectionpb.ServerReflectionResponse, error) {
		res, err := detailed(context.Background(), connect.NewRequest(req))
		if err != nil {
			return nil, err
		}
		return res.Msg, err
	}

	t.Run("list_services", func(t *testing.T) {
		req := &reflectionpb.ServerReflectionRequest{
			Host: "some-host",
			MessageRequest: &reflectionpb.ServerReflectionRequest_ListServices{
				ListServices: "ignored per proto documentation",
			},
		}
		res, err := call(req)
		assert.Nil(t, err, "reflection RPC error")
		expect := &reflectionpb.ServerReflectionResponse{
			ValidHost:       req.Host,
			OriginalRequest: req,
			MessageResponse: &reflectionpb.ServerReflectionResponse_ListServicesResponse{
				ListServicesResponse: &reflectionpb.ListServiceResponse{
					Service: []*reflectionpb.ServiceResponse{
						{Name: "connect.ping.v1test.PingService"},
					},
				},
			},
		}
		assert.Equal(t, res, expect, "response")
	})
	t.Run("file_by_filename", func(t *testing.T) {
		req := &reflectionpb.ServerReflectionRequest{
			Host: "some-host",
			MessageRequest: &reflectionpb.ServerReflectionRequest_FileByFilename{
				FileByFilename: "connect/ping/v1test/ping.proto",
			},
		}
		res, err := call(req)
		assert.Nil(t, err, "reflection RPC error")
		assert.Nil(t, res.GetErrorResponse(), "error in response")
		fds := res.GetFileDescriptorResponse()
		assert.NotNil(t, fds, "file descriptor response")
		assert.Equal(t, len(fds.FileDescriptorProto), 1, "number of fds returned")
		assert.True(
			t,
			bytes.Contains(fds.FileDescriptorProto[0], []byte(pingRequestFQN)),
			"fd should contain PingRequest struct",
		)
	})
	t.Run("file_containing_symbol", func(t *testing.T) {
		req := &reflectionpb.ServerReflectionRequest{
			Host: "some-host",
			MessageRequest: &reflectionpb.ServerReflectionRequest_FileContainingSymbol{
				FileContainingSymbol: pingRequestFQN,
			},
		}
		res, err := call(req)
		assert.Nil(t, err, "reflection RPC error")
		assert.Nil(t, res.GetErrorResponse(), "error in response")
		fds := res.GetFileDescriptorResponse()
		assert.NotNil(t, fds, "file descriptor response")
		assert.Equal(t, len(fds.FileDescriptorProto), 1, "number of fds returned")
		assert.True(
			t,
			bytes.Contains(fds.FileDescriptorProto[0], []byte(pingRequestFQN)),
			"fd should contain PingRequest struct",
		)
	})
	t.Run("file_containing_extension", func(t *testing.T) {
		req := &reflectionpb.ServerReflectionRequest{
			Host: "some-host",
			MessageRequest: &reflectionpb.ServerReflectionRequest_FileContainingExtension{
				FileContainingExtension: &reflectionpb.ExtensionRequest{
					ContainingType:  pingRequestFQN,
					ExtensionNumber: 42,
				},
			},
		}
		res, err := call(req)
		assert.Nil(t, err, "reflection RPC error")
		msgerr := res.GetErrorResponse()
		assert.NotNil(t, msgerr, "error in response proto")
		assert.Equal(t, msgerr.ErrorCode, int32(connect.CodeNotFound), "error code")
		assert.NotZero(t, msgerr.ErrorMessage, "error message")
	})
	t.Run("all_extension_numbers_of_type", func(t *testing.T) {
		req := &reflectionpb.ServerReflectionRequest{
			Host: "some-host",
			MessageRequest: &reflectionpb.ServerReflectionRequest_AllExtensionNumbersOfType{
				AllExtensionNumbersOfType: pingRequestFQN,
			},
		}
		res, err := call(req)
		assert.Nil(t, err, "reflection RPC error")
		expect := &reflectionpb.ServerReflectionResponse{
			ValidHost:       req.Host,
			OriginalRequest: req,
			MessageResponse: &reflectionpb.ServerReflectionResponse_AllExtensionNumbersResponse{
				AllExtensionNumbersResponse: &reflectionpb.ExtensionNumberResponse{
					BaseTypeName: pingRequestFQN,
				},
			},
		}
		assert.Equal(t, res, expect, "response")
	})
}
