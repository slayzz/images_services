package service

import (
	"context"
	"errors"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/slayzz/images_services/pkg/forwarder/service"
)

type ImageService interface {
	HandleImage(ctx context.Context, imageBytes []byte, message string) error
}

func NewImagerService(logger log.Logger, imagesSize metrics.Counter, grpcClient service.ForwarderService) ImageService {
	var svc ImageService
	{
		svc = &imageService{grpcClient: grpcClient}
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

type imageService struct {
	grpcClient service.ForwarderService
}

func (s *imageService) HandleImage(ctx context.Context, imageBytes []byte, message string) error {
	return s.grpcClient.HandleImage(ctx, imageBytes, message)
}
