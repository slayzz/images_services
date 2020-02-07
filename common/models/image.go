package models

type ImageRequest struct {
	Image   []byte `json:"image"`
	Message string `json:"message"`
}

type ImageResponse struct {
	Err string `json:"error"`
}
