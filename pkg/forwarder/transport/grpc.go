package transport

import (
	"context"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/tracing/opentracing"
	"github.com/go-kit/kit/transport"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/slayzz/images_services/pb"
	endpointpkg "github.com/slayzz/images_services/pkg/forwarder/endpoint"
	"github.com/slayzz/images_services/pkg/imager/service"
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"time"
)

type grpcServer struct {
	imageHandle grpctransport.Handler
}

func NewGRPCServer(e endpointpkg.Set, otTracer stdopentracing.Tracer, logger log.Logger) pb.RedirectionServer {
	options := []grpctransport.ServerOption{
		grpctransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	}

	return &grpcServer{
		imageHandle: grpctransport.NewServer(
			e.HandleImageEndpoint,
			decodeGRPCImageRequest,
			encodeGRPCImageResponse,
			append(options, grpctransport.ServerBefore(opentracing.GRPCToContext(otTracer, "ImageHandle", logger)))...,
		),
	}
}

func (s *grpcServer) ImageHandle(ctx context.Context, req *pb.Image) (*pb.ImageResponse, error) {
	_, rep, err := s.imageHandle.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.ImageResponse), nil
}

func NewGRPCClient(conn *grpc.ClientConn, otTracer stdopentracing.Tracer, logger log.Logger) service.ImageService {
	limiter := ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 100))

	var options []grpctransport.ClientOption

	var imageHandleEndpoint endpoint.Endpoint
	{
		imageHandleEndpoint = grpctransport.NewClient(
			conn,
			"Redirection",
			"HandleImage",
			encodeGRPCImageRequest,
			decodeGRPCImageResponse,
			pb.ImageResponse{},
			append(options, grpctransport.ClientBefore(opentracing.ContextToGRPC(otTracer, logger)))...,
		).Endpoint()
		imageHandleEndpoint = limiter(imageHandleEndpoint)
		imageHandleEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "ImageHandle",
			Timeout: 30 * time.Second,
		}))(imageHandleEndpoint)
	}

	return &endpointpkg.Set{HandleImageEndpoint: imageHandleEndpoint}
}

func decodeGRPCImageRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.Image)
	return endpointpkg.Image{
		Image: req.Image,
		Msg:   req.Message,
	}, nil
}

func decodeGRPCImageResponse(_ context.Context, _ interface{}) (interface{}, error) {
	return endpointpkg.ImageResponse{Err: nil}, nil
}

func encodeGRPCImageRequest(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(endpointpkg.Image)
	return &pb.Image{
		Image:   resp.Image,
		Message: resp.Msg,
	}, nil
}

func encodeGRPCImageResponse(_ context.Context, response interface{}) (interface{}, error) {
	return &pb.ImageResponse{}, nil
}
