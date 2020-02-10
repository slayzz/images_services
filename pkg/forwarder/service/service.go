package service

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/log"
)

type ForwarderService interface {
	HandleImage(ctx context.Context, image []byte, message string) error
}

func NewForwarderService(grpcClient ForwarderService, logger log.Logger) ForwarderService {
	var svc ForwarderService
	{
		svc = &forwarderService{telegramGRPCClient: grpcClient}
		svc = LoggingMiddleware(logger)(svc)
	}
	return svc
}

type forwarderService struct {
	telegramGRPCClient ForwarderService
}

func (s *forwarderService) HandleImage(ctx context.Context, imageBytes []byte, message string) error {
	fmt.Println("This is forwarder")
	return s.telegramGRPCClient.HandleImage(ctx, imageBytes, message)
}
