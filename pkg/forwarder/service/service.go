package service

import (
	"context"
	"fmt"
)

type ForwarderService interface {
	HandleImage(ctx context.Context, image []byte, message string) error
}

func NewForwarderService() ForwarderService {
	return &forwarderService{}
}

type forwarderService struct{}

func (s *forwarderService) HandleImage(ctx context.Context, imageBytes []byte, message string) error {
	//request := models.ImageRequest{
	//	Image:   imageBytes,
	//	Message: message,
	//}

	//_, err := i.grpcHandleImage(ctx, request)
	//if err != nil {
	//	return err
	//}
	fmt.Println("HELLO WE ARE HERER")
	return nil
}
