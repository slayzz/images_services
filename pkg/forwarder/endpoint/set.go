package endpoint

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/slayzz/images_services/pkg/forwarder/service"
	endpoint2 "github.com/slayzz/images_services/pkg/imager/endpoint"
)

type Set struct {
	HandleImageEndpoint endpoint.Endpoint
}

func NewForwarderSet(svc service.ForwarderService, logger log.Logger) Set {
	var handleImageEndpoint endpoint.Endpoint
	{
		handleImageEndpoint = MakeImageHandleEndpoint(svc)
		handleImageEndpoint = LoggingMiddleware(log.With(logger, "method", "HandleImage"))(handleImageEndpoint)
	}
	return Set{
		HandleImageEndpoint: handleImageEndpoint,
	}
}

func MakeImageHandleEndpoint(svc service.ForwarderService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(Image)
		err := svc.HandleImage(ctx, req.Image, req.Msg)
		return endpoint2.ImageResponse{Err: err}, err
	}
}

func (s *Set) HandleImage(ctx context.Context, imageBytes []byte, message string) error {
	resp, err := s.HandleImageEndpoint(ctx, Image{Image: imageBytes, Msg: message})
	if err != nil {
		return err
	}

	response := resp.(ImageResponse)
	if response.Err != nil {
		return response.Err
	}

	return nil
}

type Image struct {
	Image []byte
	Msg   string
}

type ImageResponse struct {
	Err error
}
