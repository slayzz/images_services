package service

import (
	"context"
	"errors"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
)

type ImageService interface {
	HandleImage(ctx context.Context, imageBytes []byte, message string) error
}

func New(logger log.Logger, imagesSize metrics.Counter) ImageService {
	var svc ImageService
	{
		svc = NewImageService()
		svc = LoggingMiddleware(logger)(svc)
		svc = InstrumentingMiddleware(imagesSize)(svc)
	}
	return svc
}

var (
	ErrTwoZeroes       = errors.New("can't sum two zeroes")
	ErrIntOverflow     = errors.New("integer overflow")
	ErrMaxSizeExceeded = errors.New("result exceeds maximum size")
)

func NewImageService() ImageService {
	return &imageService{}
}

type imageService struct{}

func (s *imageService) HandleImage(ctx context.Context, imageBytes []byte, message string) error {
	//request := models.ImageRequest{
	//	Image:   imageBytes,
	//	Message: message,
	//}

	//_, err := i.grpcHandleImage(ctx, request)
	//if err != nil {
	//	return err
	//}
	return nil
}
