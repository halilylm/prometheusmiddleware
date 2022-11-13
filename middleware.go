package prometheusmiddleware

import (
	"strconv"
	"time"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// PrometheusMiddleware holds collectors, path and registry
// config values
type PrometheusMiddleware struct {
	requestCount    *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
	path            string
	registry        *prometheus.Registry
}

// NewPrometheusMiddleware creates a new instance of prometheus middleware
func NewPrometheusMiddleware(registry *prometheus.Registry, path string) *PrometheusMiddleware {
	reqCount := promauto.With(registry).NewCounterVec(prometheus.CounterOpts{
		Name: "app_http_requests_total",
		Help: "The total number of HTTP requests.",
	}, []string{"method", "path", "code"})

	requestDurations := promauto.With(registry).NewHistogramVec(prometheus.HistogramOpts{
		Name: "app_http_request_duration_seconds",
		Help: "HTTP request durations.",
		Buckets: []float64{
			0.001, // 1ms
			0.002,
			0.005,
			0.01,
			0.02,
			0.05,
			0.1,
			0.2,
			0.5,
			1.0,
			2.0,
			5.0, // 5s
		},
	}, []string{"code"})
	return &PrometheusMiddleware{
		requestCount:    reqCount,
		requestDuration: requestDurations,
		path:            path,
		registry:        registry,
	}
}

// Use is the middleware for the prometheus metrics
func (p *PrometheusMiddleware) Use(app *fiber.App) {
	app.Use(p.middleware())
	app.Get(p.path, adaptor.HTTPHandler(promhttp.HandlerFor(p.registry, promhttp.HandlerOpts{})))
}

// middleware creates a handler to calculate the metrics
func (p *PrometheusMiddleware) middleware() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		if ctx.Path() == p.path {
			return ctx.Next()
		}
		start := time.Now()
		err := ctx.Next()
		code := strconv.Itoa(ctx.Response().StatusCode())
		p.requestCount.WithLabelValues(ctx.Route().Method, ctx.Route().Path, code).Inc()
		p.requestDuration.WithLabelValues(code).Observe(time.Since(start).Seconds())
		return err
	}
}
