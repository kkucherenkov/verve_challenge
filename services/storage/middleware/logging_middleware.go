package middleware

import (
	"context"
	"encoding/json"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"time"
	"verve_challenge_storage/pkg/model"
	"verve_challenge_storage/service"
)

type loggingMiddleware struct {
	service.Service
	logger log.Logger
}

func LoggingMiddleware(logger log.Logger) ServiceMiddleware {
	return func(next service.Service) service.Service {
		return loggingMiddleware{next, logger}
	}
}

func (mw loggingMiddleware) Size(ctx context.Context) (err error, size int) {
	defer func(begin time.Time) {
		level.Info(mw.logger).Log(
			"function", "Size",
			"result", size,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	err, size = mw.Service.Size(ctx)
	return
}

func (mw loggingMiddleware) Get(ctx context.Context, id string) (err error, item model.Item) {
	defer func(begin time.Time) {
		res, err := json.Marshal(item)
		result := string(res)
		level.Info(mw.logger).Log(
			"function", "GetItem",
			"id", id,
			"result", result,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	err, item = mw.Service.Get(ctx, id)
	return
}

func (mw loggingMiddleware) Reload(ctx context.Context, path string) (err error) {
	defer func(begin time.Time) {
		level.Info(mw.logger).Log(
			"function", "ReloadDb",
			"path", path,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	err = mw.Service.Reload(ctx, path)
	return
}
