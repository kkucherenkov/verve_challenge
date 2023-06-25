package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-kit/kit/log"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/go-kit/log/level"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"syscall"
	"verve_challenge_storage/endpoints"
	"verve_challenge_storage/implementation"
	"verve_challenge_storage/middleware"
	"verve_challenge_storage/pb"
	"verve_challenge_storage/pkg/config"
	"verve_challenge_storage/service"
	"verve_challenge_storage/transport"
)

func main() {
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	//declare metrics
	fieldKeys := []string{"method"}
	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "verve_challenge",
		Subsystem: "storage",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)
	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "verve_challenge",
		Subsystem: "storage",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)

	var cfg config.Config
	err := config.ReadFile(&cfg)
	if err != nil {
		return
	}
	err = config.ReadEnv(&cfg)
	if err != nil {
		return
	}

	level.Info(logger).Log("msg", "config loaded")

	gRPCAddr := cfg.Storage.Host + ":" + cfg.Storage.Port

	var svc service.Service
	svc = implementation.CreateStorage()
	svc = middleware.LoggingMiddleware(logger)(svc)
	svc = middleware.MetricsMiddleware(requestCount, requestLatency)(svc)

	errChan := make(chan error)
	ctx := context.Background()
	level.Info(logger).Log("msg", "service created")
	err = svc.Reload(ctx, cfg.FileName)
	if err != nil {
		level.Error(logger).Log("msg", "data load failed", "err", err)
		return
	}
	level.Info(logger).Log("msg", "data loaded")

	// creating Endpoints struct
	e := endpoints.Endpoints{
		GetSizeEndpoint:  endpoints.MakeGetSizeEndpoint(svc),
		GetItemEndpoint:  endpoints.MakeGetItemEndpoint(svc),
		ReloadDbEndpoint: endpoints.MakeReloadDbEndpoint(svc),
	}
	go func() {
		listener, err := net.Listen("tcp", gRPCAddr)
		if err != nil {
			errChan <- err
			return
		}
		handler := transport.NewGRPCServer(ctx, e)
		gRPCServer := grpc.NewServer()
		pb.RegisterStorageServer(gRPCServer, handler)
		level.Info(logger).Log("msg", "service started")
		errChan <- gRPCServer.Serve(listener)

	}()

	router := gin.Default()

	// Register your endpoint with Gin
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	go func() {
		errChan <- router.Run(cfg.Storage.Host + ":" + cfg.Storage.MetricsPort)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()
	level.Error(logger).Log("msg", "service error", "err", <-errChan)
}
