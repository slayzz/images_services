package main

import (
	"flag"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/oklog/oklog/pkg/group"
	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/slayzz/images_services/pb"
	"github.com/slayzz/images_services/pkg/forwarder/endpoint"
	"github.com/slayzz/images_services/pkg/forwarder/service"
	"github.com/slayzz/images_services/pkg/forwarder/transport"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	fs := flag.NewFlagSet("redirection", flag.ExitOnError)
	var (
		grpcAddr = fs.String("grpc-addr", ":10000", "gRPC listen address")
	)
	fs.Parse(os.Args[1:])

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	var tracer stdopentracing.Tracer
	{
		tracer = stdopentracing.GlobalTracer() // no-op
	}

	var (
		svc       = service.NewForwarderService()
		endpoints = endpoint.NewForwarderSet(svc, logger)
		server    = transport.NewGRPCServer(endpoints, tracer, logger)
	)

	var g group.Group
	{
		grpcListener, err := net.Listen("tcp", *grpcAddr)
		if err != nil {
			logger.Log("transport", "gRPC", "during", "Listen", "err", err)
			os.Exit(1)
		}

		g.Add(func() error {
			logger.Log("transport", "gRPC", "addr", *grpcAddr)
			baseServer := grpc.NewServer()
			pb.RegisterRedirectionServer(baseServer, server)
			return baseServer.Serve(grpcListener)
		}, func(error) {
			grpcListener.Close()
		})
	}

	{
		// This function just sits and waits for ctrl-C.
		cancelInterrupt := make(chan struct{})
		g.Add(func() error {
			c := make(chan os.Signal, 1)
			signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
			select {
			case sig := <-c:
				return fmt.Errorf("received signal %s", sig)
			case <-cancelInterrupt:
				return nil
			}
		}, func(error) {
			close(cancelInterrupt)
		})
	}
	logger.Log("exit", g.Run())
}
