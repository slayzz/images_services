package impl

import (
	"github.com/go-kit/kit/log"
	"github.com/slayzz/images_services/redirection/pkg/service"
)

type redirectionService struct {
	logger log.Logger
}

func (r *redirectionService) HandleImage(image []byte, message string) error {
	r.logger.Log("msg", "We have that message: "+message)
	return nil
}

func NewRedirectionService(logger log.Logger) service.RedirectionService {
	return &redirectionService{logger: logger}
}
