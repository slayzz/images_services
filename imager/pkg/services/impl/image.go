package impl

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/slayzz/images_services/common/endpoints"
	"github.com/slayzz/images_services/common/models"
	"github.com/slayzz/images_services/common/proto/redirection"
	"github.com/slayzz/images_services/imager/pkg/services"
	"google.golang.org/grpc"
)

type imageService struct {
	grpcHandleImage endpoint.Endpoint
	logger          log.Logger
}

func (i *imageService) HandleImage(ctx context.Context, imageBytes []byte, message string) error {
	request := models.ImageRequest{
		Image:   imageBytes,
		Message: message,
	}

	_, err := i.grpcHandleImage(ctx, request)
	if err != nil {
		return err
	}

	i.logger.Log("msg", "all-ok")
	return nil
}

func NewImageService(conn *grpc.ClientConn, logger log.Logger) services.ImageService {
	return &imageService{
		grpcHandleImage: grpctransport.NewClient(
			conn,
			"Redirection",
			"ImageHandle",
			endpoints.EncodeGRPCHandleImageResponse,
			endpoints.DecodeGRPCHandleImageResponse,
			redirection.ImageResponse{},
		).Endpoint(),
		logger: logger,
	}
}
