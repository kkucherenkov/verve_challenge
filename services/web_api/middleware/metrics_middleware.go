package middleware

import (
	"github.com/go-kit/kit/metrics"
	"time"
	"verve_challenge_web_api/pkg/model"
	"verve_challenge_web_api/service"
)

type metricsMiddleware struct {
	service.Service
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
}

func MetricsMiddleware(requestCount metrics.Counter,
	requestLatency metrics.Histogram) ServiceMiddleware {
	return func(next service.Service) service.Service {
		return metricsMiddleware{
			next,
			requestCount,
			requestLatency,
		}
	}
}

func (mw metricsMiddleware) GetItem(id string) (err error, item model.Item) {
	defer func(begin time.Time) {
		lvs := []string{"method", "GetItem"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	err, item = mw.Service.GetItem(id)
	return
}

func (mw metricsMiddleware) Reload(path string) (err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "ReloadDb"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	err = mw.Service.Reload(path)
	return
}
