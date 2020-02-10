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
	"github.com/slayzz/images_services/pkg/forwarder/service"
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"time"
)

type grpcServer struct {
	handleImage grpctransport.Handler
}

func NewGRPCForwarderServer(e endpointpkg.Set, otTracer stdopentracing.Tracer, logger log.Logger) forwarder.ForwarderServer {
	options := []grpctransport.ServerOption{
		grpctransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	}

	return &grpcServer{
		handleImage: grpctransport.NewServer(
			e.HandleImageEndpoint,
			decodeGRPCImageRequest,
			encodeGRPCImageResponse,
			append(options, grpctransport.ServerBefore(opentracing.GRPCToContext(otTracer, "ImageHandle", logger)))...,
		),
	}
}

func (s *grpcServer) HandleImage(ctx context.Context, req *forwarder.Image) (*forwarder.ImageResponse, error) {
	_, rep, err := s.handleImage.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*forwarder.ImageResponse), nil
}

func NewGRPCForwarderClient(conn *grpc.ClientConn, otTracer stdopentracing.Tracer, logger log.Logger) service.ForwarderService {
	limiter := ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 100))

	var options []grpctransport.ClientOption

	var imageHandleEndpoint endpoint.Endpoint
	{
		imageHandleEndpoint = grpctransport.NewClient(
			conn,
			"Forwarder",
			"HandleImage",
			encodeGRPCImageRequest,
			decodeGRPCImageResponse,
			forwarder.ImageResponse{},
			append(options, grpctransport.ClientBefore(opentracing.ContextToGRPC(otTracer, logger)))...,
		).Endpoint()
		imageHandleEndpoint = limiter(imageHandleEndpoint)
		imageHandleEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "HandleImage",
			Timeout: 30 * time.Second,
		}))(imageHandleEndpoint)
	}

	return &endpointpkg.Set{HandleImageEndpoint: imageHandleEndpoint}
}

func decodeGRPCImageRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*forwarder.Image)
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
	return &forwarder.Image{
		Image:   resp.Image,
		Message: resp.Msg,
	}, nil
}

func encodeGRPCImageResponse(_ context.Context, response interface{}) (interface{}, error) {
	return &forwarder.ImageResponse{}, nil
}
