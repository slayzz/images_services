package main

import (
	"flag"
	"github.com/go-kit/kit/log"
	"github.com/slayzz/images_services/common/middlewares"
	"github.com/slayzz/images_services/imager/pkg/services/impl"
	"github.com/slayzz/images_services/imager/pkg/transport"
	"google.golang.org/grpc"
	"net/http"
	"os"
)

const (
	httpport           = ":9000"
	grpcport           = ":10000"
	redirectionAddress = "localhost:10001"
)

func main() {
	var (
		httpListen = flag.String("httpport", httpport, "HTTP listen address")
		_          = flag.String("rpcport", grpcport, "gRPC listen address")
	)
	flag.Parse()

	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stdout)
	logger = log.With(
		logger,
		"time", log.DefaultTimestamp,
		"service", "imagery",
		"caller", log.DefaultCaller,
	)

	conn, err := grpc.Dial(redirectionAddress, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		logger.Log("did not connect: ", err)
		os.Exit(1)
	}

	svc := impl.NewImageService(conn, logger)
	svc = middlewares.LoggingMiddleware(logger)(svc)

	handler := transport.MakeHTTPHandler(svc, logger)
	logger.Log("msg", "HTTP", "addr", *httpListen)
	logger.Log("err", http.ListenAndServe(*httpListen, handler))
}
