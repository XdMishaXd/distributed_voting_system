package router

import (
	"notification_service/internal/http_server/handlers/infrastructure/health"
	metricshandler "notification_service/internal/http_server/handlers/infrastructure/metrics_handler"
	"notification_service/internal/metrics"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func SetupRouter(m *metrics.Metrics) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)

	r.Get("/health", health.New())
	r.Get("/metrics", metricshandler.New(m))

	return r
}
