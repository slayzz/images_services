package transport

import (
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	endpoints2 "github.com/slayzz/images_services/common/endpoints"
	"github.com/slayzz/images_services/imager/pkg/services"
	"net/http"
)

func MakeHTTPHandler(s services.ImageService, logger log.Logger) http.Handler {
	router := mux.NewRouter().StrictSlash(true)
	e := endpoints2.MakeEndpoints(s)
	router.Methods(http.MethodPost).Path("/image/upload").Handler(httptransport.NewServer(
		e.ImageHandleEndpoint,
		endpoints2.DecodeImageRequest,
		endpoints2.EncodeJsonResponse,
	))

	return router
}
