package transport

import (
	"context"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/slayzz/images_services/common/endpoints"
	"github.com/slayzz/images_services/common/proto/redirection"
)

type grpcServer struct {
	imageHandle grpctransport.Handler
}

func (s *grpcServer) ImageHandle(ctx context.Context, request *redirection.Image) (*redirection.ImageResponse, error) {
	_, resp, err := s.imageHandle.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	return resp.(*redirection.ImageResponse), nil
}

func NewGRPCServer(ctx context.Context, e endpoints.RedirectionEndpoints) redirection.RedirectionServer {
	return &grpcServer{
		imageHandle: grpctransport.NewServer(
			e.ImageHandleEndpoint,
			endpoints.DecodeGRPCHandleImageRequest,
			endpoints.EncodeGRPCHandleImageResponse,
		),
	}
}
