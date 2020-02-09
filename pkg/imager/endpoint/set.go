package endpoint

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/slayzz/images_services/pkg/imager/service"
)

type Set struct {
	ImageHandleEndpoint endpoint.Endpoint
}

func NewImagerSet(svc service.ImageService, logger log.Logger, duration metrics.Histogram) Set {
	var imageHandleEndpoint endpoint.Endpoint
	{
		imageHandleEndpoint = MakeImageHandleEndpoint(svc)
		imageHandleEndpoint = LoggingMiddleware(log.With(logger, "method", "HandleImage"))(imageHandleEndpoint)
		imageHandleEndpoint = InstrumentingMiddleware(duration.With("method", "HandleImage"))(imageHandleEndpoint)
	}
	return Set{
		ImageHandleEndpoint: imageHandleEndpoint,
	}
}

func MakeImageHandleEndpoint(svc service.ImageService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ImageRequest)
		err := svc.HandleImage(ctx, req.Image, req.Message)
		return ImageResponse{Err: err}, err
	}
}

func (s *Set) HandleImage(ctx context.Context, imageBytes []byte, message string) error {
	resp, err := s.ImageHandleEndpoint(ctx, ImageRequest{Image: imageBytes, Message: message})
	if err != nil {
		return err
	}

	response := resp.(ImageResponse)
	if response.Err != nil {
		return response.Err
	}

	return nil
}

// compile time assertions for our response types implementing endpoint.Failer.
var (
	_ endpoint.Failer = ImageResponse{}
)

type ImageRequest struct {
	Image   []byte `json:"image"`
	Message string `json:"message"`
}

type ImageResponse struct {
	Err error `json:"error"`
}

func (r ImageResponse) Failed() error {
	return r.Err
}
