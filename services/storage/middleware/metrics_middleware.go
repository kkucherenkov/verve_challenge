package middleware

import (
	"context"
	"github.com/go-kit/kit/metrics"
	"time"
	"verve_challenge_storage/pkg/model"
	"verve_challenge_storage/service"
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

func (mw metricsMiddleware) Size(ctx context.Context) (err error, size int) {
	defer func(begin time.Time) {
		lvs := []string{"method", "Size"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	err, size = mw.Service.Size(ctx)
	return
}

func (mw metricsMiddleware) Get(ctx context.Context, id string) (err error, item model.Item) {
	defer func(begin time.Time) {
		lvs := []string{"method", "GetItem"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	err, item = mw.Service.Get(ctx, id)
	return
}

func (mw metricsMiddleware) Reload(ctx context.Context, path string) (err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "ReloadDb"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	err = mw.Service.Reload(ctx, path)
	return
}
