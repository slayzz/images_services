package service

type RedirectionService interface {
	HandleImage(image []byte, message string) error
}
