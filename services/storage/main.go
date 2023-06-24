package main

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/log/level"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"syscall"
	"verve_challenge_storage/endpoints"
	"verve_challenge_storage/implementation"
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
	svc = service.LoggingMiddleware(logger)(svc)
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
