package endpoints

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	"github.com/slayzz/images_services/common/models"
	"github.com/slayzz/images_services/imager/pkg/services"
	"io/ioutil"
	"log"
	"net/http"
)

type Endpoints struct {
	ImageHandleEndpoint endpoint.Endpoint
}

func makeImageHandleEndpoint(svc services.ImageService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(models.ImageRequest)

		err := svc.HandleImage(ctx, req.Image, req.Message)

		if err != nil {
			return models.ImageResponse{Err: err.Error()}, err
		}

		return models.ImageResponse{Err: ""}, nil
	}
}

func MakeEndpoints(svc services.ImageService) Endpoints {
	return Endpoints{
		ImageHandleEndpoint: makeImageHandleEndpoint(svc),
	}
}

func DecodeImageRequest(_ context.Context, r *http.Request) (interface{}, error) {
	file, _, err := r.FormFile("image")
	if err != nil {
		log.Println("error on getting a file from multipart", err)
		return nil, err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println("error on reading a file from multipart", err)
		return nil, err
	}
	message := r.FormValue("message")

	return models.ImageRequest{Image: bytes, Message: message}, nil
}

func EncodeJsonResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}
