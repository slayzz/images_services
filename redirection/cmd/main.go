package main

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/slayzz/images_services/common/endpoints"
	"github.com/slayzz/images_services/common/proto/redirection"
	"github.com/slayzz/images_services/redirection/pkg/service/impl"
	"github.com/slayzz/images_services/redirection/pkg/transport"
	"google.golang.org/grpc"
	"net"
	"os"
)

const (
	port = ":10001"
)

func main() {
	ctx := context.Background()

	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stdout)
	logger = log.With(
		logger,
		"time", log.DefaultTimestamp,
		"service", "redirection",
		"caller", log.DefaultCaller,
	)

	svc := impl.NewRedirectionService(logger)
	ep := endpoints.NewRedirectionEndpoints(svc)

	listener, err := net.Listen("tcp", port)
	if err != nil {
		logger.Log("error", err)
		os.Exit(1)
	}

	handler := transport.NewGRPCServer(ctx, ep)
	gRPCServer := grpc.NewServer()
	redirection.RegisterRedirectionServer(gRPCServer, handler)
	logger.Log("msg", "HTTP", "addr", port)
	logger.Log("err", gRPCServer.Serve(listener))
}
