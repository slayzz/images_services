package middlewares

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/slayzz/images_services/imager/pkg/services"
	"time"
)

type ServiceMiddleware func(next services.ImageService) services.ImageService

type loggerMiddleware struct {
	logger log.Logger
	services.ImageService
}

func LoggingMiddleware(logger log.Logger) ServiceMiddleware {
	return func(next services.ImageService) services.ImageService {
		return &loggerMiddleware{
			logger:       logger,
			ImageService: next,
		}
	}
}

func (l *loggerMiddleware) logFunc(begin time.Time, method string, err error) {
	level.Info(l.logger).Log(
		"method", method,
		"err", err,
		"took", time.Since(begin),
	)
}

func (l *loggerMiddleware) HandleImage(ctx context.Context, imageBytes []byte, message string) error {
	defer l.logFunc(time.Now(), "HandleImage", nil)
	return l.ImageService.HandleImage(ctx, imageBytes, message)
}

