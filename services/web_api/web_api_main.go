package main

import (
	"context"
	"github.com/gin-gonic/gin"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"verve_challenge_web_api/endpoints"
	"verve_challenge_web_api/model"
	"verve_challenge_web_api/service"

	"os"
	"verve_challenge_web_api/pkg/config"
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

	grpcAddr := cfg.WebApi.StorageAddr

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		level.Error(logger).Log("msg", "error connection closing", "err", err.Error())
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			level.Error(logger).Log("msg", "error connection closing", "err", err.Error())
		}
	}(conn)

	svc := service.NewService(ctx, conn)
	svc = service.LoggingMiddleware(logger)(svc)
	logger.Log("msg", "service created")
	getItemHandler := httptransport.NewServer(
		endpoints.MakeGetItemEndpoint(svc),
		model.DecodeGetItemRequest,
		model.EncodeResponse,
	)
	reloadHandler := httptransport.NewServer(
		endpoints.MakeReloadEndpoint(svc),
		model.DecodeReloadRequest,
		model.EncodeResponse,
	)
	router := gin.Default()

	// Register your endpoint with Gin
	router.GET("/promotions/:id", gin.WrapH(getItemHandler))
	router.POST("/promotions/reload", gin.WrapH(reloadHandler))
	// Start the Gin server
	level.Info(logger).Log("msg", "service starting")
	err = router.Run(cfg.WebApi.Host + ":" + cfg.WebApi.Port)
	if err != nil {
		level.Error(logger).Log("msg", "service start failed", "err", err)
		return
	}

}
