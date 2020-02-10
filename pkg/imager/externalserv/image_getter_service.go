package externalserv

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/slayzz/images_services/pkg/forwarder/service"
	"time"
)

type ImageGetterService interface {
	Run(ctx context.Context)
}

type ImageGetterServiceImpl struct {
	grpcClient     service.ForwarderService
	logger         log.Logger
	imageRequester ImageRequester
}

func NewImageGetterServiceImpl(logger log.Logger, grpcForwarderClient service.ForwarderService, imageRequestor ImageRequester) ImageGetterServiceImpl {
	return ImageGetterServiceImpl{grpcClient: grpcForwarderClient, logger: logger, imageRequester: imageRequestor}
}

func (gi *ImageGetterServiceImpl) Run(ctx context.Context) {
	go func() {
		for {
			imageBytes, err := gi.imageRequester.GetImage()
			if err != nil {
				level.Error(gi.logger).Log("method", "Run", "in", "ImageGetterServiceImpl", "err", err)
				continue
			}

			err = gi.grpcClient.HandleImage(ctx, imageBytes, "randomly got")
			if err != nil {
				level.Error(gi.logger).Log("method", "Run", "in", "ImageGetterServiceImpl", "err", err)
			}
			gi.logger.Log("success", "sent randomly image")
			time.Sleep(10 * time.Second)
		}
	}()
}
