package main

import (
	"flag"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/prometheus"
	"github.com/oklog/oklog/pkg/group"
	stdopentracing "github.com/opentracing/opentracing-go"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/slayzz/images_services/pkg/imager/endpoint"
	"github.com/slayzz/images_services/pkg/imager/service"
	"github.com/slayzz/images_services/pkg/imager/transport"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	fs := flag.NewFlagSet("imager", flag.ExitOnError)

	var (
		debugPort = fs.String("debug-addr", ":9000", "debug-http-listen-address")
		httpPort  = fs.String("http-addr", ":9001", "http-listen-address")
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
		tracer = stdopentracing.GlobalTracer()
	}

	var imagesSize metrics.Counter
	{
		// business-level metrics
		imagesSize = prometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "images_services",
			Subsystem: "imager",
			Name:      "images_size",
			Help:      "Total size of bytes summed via the HandleImage method.",
		}, []string{})
	}

	var duration metrics.Histogram
	{
		// endpoint-level metrics
		duration = prometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "images_services",
			Subsystem: "imager",
			Name:      "request_duration_seconds",
			Help:      "Request duration in seconds.",
		}, []string{"method", "success"})
	}
	http.DefaultServeMux.Handle("/metrics", promhttp.Handler())

	svc := service.New(logger, imagesSize)
	endpoints := endpoint.NewImagerSet(svc, logger, duration)
	httpHandler := transport.NewHTTPHandler(endpoints, tracer, logger)

	var g group.Group
	{
		debugListener, err := net.Listen("tcp", *debugPort)
		if err != nil {
			logger.Log("transport", "debug/HTTP", "during", "Listen", "err", err)
			os.Exit(1)
		}
		g.Add(func() error {
			logger.Log("transport", "debug/HTTP", "addr", *debugPort)
			return http.Serve(debugListener, http.DefaultServeMux)
		}, func(error) {
			debugListener.Close()
		})
	}

	{
		httpListener, err := net.Listen("tcp", *httpPort)
		if err != nil {
			logger.Log("transport", "HTTP", "during", "Listen", "err", err)
			os.Exit(1)
		}
		g.Add(func() error {
			logger.Log("transport", "HTTP", "addr", *httpPort)
			return http.Serve(httpListener, httpHandler)
		}, func(error) {
			httpListener.Close()
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
