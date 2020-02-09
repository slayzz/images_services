package transport

import (
	"context"
	"encoding/json"
	kitendpoint "github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/examples/addsvc/pkg/addservice"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/tracing/opentracing"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/slayzz/images_services/pkg/imager/endpoint"
	"io/ioutil"

	"net/http"
)

func NewHTTPHandler(endpoints endpoint.Set, otTracer stdopentracing.Tracer, logger log.Logger) http.Handler {
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(errorEncoder),
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	}
	router := mux.NewRouter().StrictSlash(true)

	router.Methods(http.MethodPost).Path("/image/upload").Handler(httptransport.NewServer(
		endpoints.ImageHandleEndpoint,
		decodeHTTPImageRequest,
		encodeHTTPGenericResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(otTracer, "ImageHandle", logger)))...,
	))

	return router
}


func decodeHTTPImageRequest(_ context.Context, r *http.Request) (interface{}, error) {
	file, _, err := r.FormFile("image")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	message := r.FormValue("message")
	return endpoint.ImageRequest{Image: bytes, Message: message}, nil
}

func encodeHTTPImageResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func encodeHTTPGenericResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if f, ok := response.(kitendpoint.Failer); ok && f.Failed() != nil {
		errorEncoder(ctx, f.Failed(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

type errorWrapper struct {
	Error string `json:"error"`
}

func errorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	w.WriteHeader(err2code(err))
	json.NewEncoder(w).Encode(errorWrapper{Error: err.Error()})
}

func err2code(err error) int {
	switch err {
	case addservice.ErrTwoZeroes, addservice.ErrMaxSizeExceeded, addservice.ErrIntOverflow:
		return http.StatusBadRequest
	}
	return http.StatusInternalServerError
}
