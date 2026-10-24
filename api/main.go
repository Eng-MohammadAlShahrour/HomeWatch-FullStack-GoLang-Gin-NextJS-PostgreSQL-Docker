package main

import (
"log"
"os"

"homelab/internal/config"
"homelab/internal/db"
"homelab/internal/handlers"
"homelab/internal/healthchecker"
"homelab/internal/middleware"

"github.com/gin-contrib/cors"
"github.com/gin-gonic/gin"
"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
cfg := config.Load()
db.Init(cfg.DatabaseURL)
go healthchecker.Start()

gin.SetMode(gin.ReleaseMode)
r := gin.Default()

r.Use(cors.New(cors.Config{
AllowOrigins:     []string{"*"},
AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
ExposeHeaders:    []string{"Content-Length"},
AllowCredentials: false,
}))
r.Use(middleware.MetricsMiddleware())
r.Use(middleware.RateLimit(middleware.NewRateLimiter(300, 50)))

// Public endpoints
r.GET("/metrics", gin.WrapH(promhttp.Handler()))
r.GET("/health", handlers.HealthCheck)

api := r.Group("/api")
{
api.POST("/register", handlers.Register)
api.POST("/login", handlers.Login)

auth := api.Group("")
auth.Use(middleware.AuthRequired(cfg.JWTSecret))
{
auth.GET("/me", handlers.Me)
auth.GET("/devices", handlers.ListDevices)
auth.POST("/devices", handlers.CreateDevice)
auth.PUT("/devices/:id", handlers.UpdateDevice)
auth.DELETE("/devices/:id", handlers.DeleteDevice)
auth.GET("/services", handlers.ListServices)
auth.POST("/services", handlers.CreateService)
auth.PUT("/services/:id", handlers.UpdateService)
auth.DELETE("/services/:id", handlers.DeleteService)
auth.GET("/dashboard", handlers.Dashboard)
auth.GET("/uptime/summary", handlers.GetUptimeSummary)
auth.GET("/services/:id/uptime", handlers.GetServiceUptime)

// Export
auth.GET("/export/services.csv", handlers.ExportServicesCSV)
auth.GET("/export/devices.csv", handlers.ExportDevicesCSV)
auth.GET("/export/services/:id/pings.csv", handlers.ExportPingsCSV)

// Search
auth.GET("/search/services", handlers.SearchServices)
auth.GET("/search/devices", handlers.SearchDevices)

// Bulk
auth.POST("/bulk/services/delete", handlers.BulkDeleteServices)
auth.POST("/bulk/devices/delete", handlers.BulkDeleteDevices)
auth.POST("/bulk/services/ping", handlers.BulkPingNow)

// Percentiles
auth.GET("/services/:id/percentiles", handlers.GetLatencyPercentiles)

// Webhooks
auth.GET("/webhooks", handlers.ListWebhooks)
auth.POST("/webhooks", handlers.CreateWebhook)
auth.DELETE("/webhooks/:id", handlers.DeleteWebhook)
}
}

port := os.Getenv("PORT")
if port == "" {
port = "8080"
}
log.Printf("API listening on :%s", port)
if err := r.Run(":" + port); err != nil {
log.Fatal(err)
}
}