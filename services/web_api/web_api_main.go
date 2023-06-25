package main

import (
	"context"
	"github.com/gin-gonic/gin"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"verve_challenge_web_api/endpoints"
	"verve_challenge_web_api/middleware"
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

	//declare metrics
	fieldKeys := []string{"method"}
	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "verve_challenge",
		Subsystem: "webapi",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)
	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "verve_challenge",
		Subsystem: "webapi",
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
	svc = middleware.LoggingMiddleware(logger)(svc)
	svc = middleware.MetricsMiddleware(requestCount, requestLatency)(svc)

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
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	router.POST("/promotions/reload", gin.WrapH(reloadHandler))
	// Start the Gin server
	level.Info(logger).Log("msg", "service starting")
	err = router.Run(cfg.WebApi.Host + ":" + cfg.WebApi.Port)
	if err != nil {
		level.Error(logger).Log("msg", "service start failed", "err", err)
		return
	}

}
