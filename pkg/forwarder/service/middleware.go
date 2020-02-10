package service

import (
	"context"
	"github.com/go-kit/kit/log"
)

type Middleware func(ForwarderService) ForwarderService

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next ForwarderService) ForwarderService {
		return &loggingMiddleware{
			logger: logger,
			next:   next,
		}
	}
}

type loggingMiddleware struct {
	logger log.Logger
	next   ForwarderService
}

func (mw *loggingMiddleware) HandleImage(ctx context.Context, imageBytes []byte, message string) error {
	err := mw.next.HandleImage(ctx, imageBytes, message)
	defer func() {
		mw.logger.Log("method", "HandleImage", "message", message, "err", err)
	}()
	return err
}
