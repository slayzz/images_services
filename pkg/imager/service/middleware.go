package service

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
)

type Middleware func(ImageService) ImageService

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next ImageService) ImageService {
		return &loggingMiddleware{
			logger: logger,
			next:   next,
		}
	}
}

type loggingMiddleware struct {
	logger log.Logger
	next   ImageService
}

func (mw *loggingMiddleware) HandleImage(ctx context.Context, imageBytes []byte, message string) error {
	err := mw.next.HandleImage(ctx, imageBytes, message)
	defer func() {
		mw.logger.Log("method", "HandleImage", "message", message, "err", err)
	}()
	return err
}

func InstrumentingMiddleware(imagesSize metrics.Counter) Middleware {
	return func(next ImageService) ImageService {
		return &instrumentingMiddleware{
			imagesSize: imagesSize,
			next:       next,
		}
	}
}

type instrumentingMiddleware struct {
	imagesSize metrics.Counter
	next       ImageService
}

func (mw *instrumentingMiddleware) HandleImage(ctx context.Context, imageBytes []byte, message string) error {
	mw.imagesSize.Add(float64(len(imageBytes)))
	return mw.next.HandleImage(ctx, imageBytes, message)
}
