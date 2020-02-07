package endpoints

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/slayzz/images_services/common/models"
	"github.com/slayzz/images_services/common/proto/redirection"
	"github.com/slayzz/images_services/redirection/pkg/service"
)

type RedirectionEndpoints struct {
	ImageHandleEndpoint endpoint.Endpoint
}

func NewRedirectionEndpoints(svc service.RedirectionService) RedirectionEndpoints {
	return RedirectionEndpoints{
		ImageHandleEndpoint: MakeImageHandleEndpoint(svc),
	}
}

func MakeImageHandleEndpoint(svc service.RedirectionService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(models.ImageRequest)

		err := svc.HandleImage(req.Image, req.Message)
		if err != nil {
			return nil, err
		}

		return nil, nil
	}
}

func EncodeGRPCHandleImageRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(redirection.Image)
	return models.ImageRequest{
		Image:   req.Image,
		Message: req.Message,
	}, nil
}

func DecodeGRPCHandleImageRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*redirection.Image)
	return models.ImageRequest{
		Image:   req.Image,
		Message: req.Message,
	}, nil
}

func EncodeGRPCHandleImageResponse(_ context.Context, r interface{}) (interface{}, error) {
	return &redirection.ImageResponse{}, nil
}

func DecodeGRPCHandleImageResponse(_ context.Context, r interface{}) (interface{}, error) {
	return nil, nil
}
