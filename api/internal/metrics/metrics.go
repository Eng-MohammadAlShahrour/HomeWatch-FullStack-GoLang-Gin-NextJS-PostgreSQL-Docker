package metrics

import (
"github.com/prometheus/client_golang/prometheus"
"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
HTTPRequestsTotal = promauto.NewCounterVec(
prometheus.CounterOpts{Name: "homewatch_http_requests_total", Help: "Total HTTP requests"},
[]string{"method", "path", "status"},
)
HTTPDuration = promauto.NewHistogramVec(
prometheus.HistogramOpts{Name: "homewatch_http_request_duration_seconds", Help: "HTTP request duration", Buckets: prometheus.DefBuckets},
[]string{"method", "path"},
)
ServicesUp = promauto.NewGauge(prometheus.GaugeOpts{Name: "homewatch_services_up", Help: "Services up"})
ServicesDown = promauto.NewGauge(prometheus.GaugeOpts{Name: "homewatch_services_down", Help: "Services down"})
HealthChecksTotal = promauto.NewCounterVec(
prometheus.CounterOpts{Name: "homewatch_health_checks_total", Help: "Total health checks"},
[]string{"status"},
)
HealthCheckDuration = promauto.NewHistogram(prometheus.HistogramOpts{
Name: "homewatch_health_check_duration_seconds", Help: "Health check duration", Buckets: prometheus.DefBuckets,
})
WebhooksSent = promauto.NewCounterVec(
prometheus.CounterOpts{Name: "homewatch_webhooks_sent_total", Help: "Webhooks sent"},
[]string{"status"},
)
)