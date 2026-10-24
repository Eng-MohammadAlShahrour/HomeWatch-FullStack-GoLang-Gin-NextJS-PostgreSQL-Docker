package handlers

import (
"net/http"
"runtime"
"time"

"homelab/internal/db"
"homelab/internal/models"

"github.com/gin-gonic/gin"
)

var processStart = time.Now()

func HealthCheck(c *gin.Context) {
start := time.Now()
dbStatus := "ok"
if sqlDB, err := db.DB.DB(); err != nil || sqlDB.Ping() != nil {
dbStatus = "down"
}
dbLatency := time.Since(start).Milliseconds()

var serviceCount, pingCount int64
db.DB.Model(&models.Service{}).Count(&serviceCount)
db.DB.Model(&models.ServicePing{}).Count(&pingCount)

var m runtime.MemStats
runtime.ReadMemStats(&m)

c.JSON(http.StatusOK, gin.H{
"status":         "ok",
"db":             dbStatus,
"db_latency_ms":  dbLatency,
"services":       serviceCount,
"pings":          pingCount,
"uptime_seconds": int(time.Since(processStart).Seconds()),
"go_version":     runtime.Version(),
"goroutines":     runtime.NumGoroutine(),
"memory_mb":      m.Alloc / 1024 / 1024,
})
}