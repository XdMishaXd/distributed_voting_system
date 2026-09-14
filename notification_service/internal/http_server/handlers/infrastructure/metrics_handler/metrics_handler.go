package metricshandler

import (
	"net/http"

	"notification_service/internal/metrics"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func New(m *metrics.Metrics) http.HandlerFunc {
	handler := promhttp.HandlerFor(m.Registry, promhttp.HandlerOpts{
		ErrorHandling: promhttp.ContinueOnError,
	})
	return handler.ServeHTTP
}
