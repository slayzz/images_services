package service

import (
	"context"
	"fmt"
	"github.com/slayzz/images_services/pkg/forwarder/service"
)

type telegramService struct{

}

func NewTelegramService() service.ForwarderService {
	return &telegramService{}
}

func (*telegramService) HandleImage(ctx context.Context, image []byte, message string) error {
	fmt.Println("This is telegram")
	return nil
}
