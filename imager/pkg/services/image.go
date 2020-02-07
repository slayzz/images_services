package services

import "context"

type ImageService interface {
	HandleImage(ctx context.Context, imageBytes []byte, message string) error
}
