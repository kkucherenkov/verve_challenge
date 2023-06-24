package service

import (
	"encoding/json"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"time"
	"verve_challenge_web_api/pkg/model"
)

type ServiceMiddleware func(Service) Service

type loggingMiddleware struct {
	Service
	logger log.Logger
}

func LoggingMiddleware(logger log.Logger) ServiceMiddleware {
	return func(next Service) Service {
		return loggingMiddleware{next, logger}
	}
}

func (mw loggingMiddleware) GetItem(id string) (err error, item model.Item) {
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
	err, item = mw.Service.GetItem(id)
	return
}

func (mw loggingMiddleware) Reload(path string) (err error) {
	defer func(begin time.Time) {
		level.Info(mw.logger).Log(
			"function", "ReloadDb",
			"path", path,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	err = mw.Service.Reload(path)
	return
}
