package healthchecker

import (
"log"
"net/http"
"time"

"homelab/internal/db"
"homelab/internal/handlers"
"homelab/internal/metrics"
"homelab/internal/models"
)

const pingRetention = 7 * 24 * time.Hour

func Start() {
ticker := time.NewTicker(15 * time.Second)
defer ticker.Stop()
cleanupTicker := time.NewTicker(1 * time.Hour)
defer cleanupTicker.Stop()

check()

for {
select {
case <-ticker.C:
check()
case <-cleanupTicker.C:
cleanupOldPings()
}
}
}

func check() {
var services []models.Service
db.DB.Find(&services)
client := &http.Client{Timeout: 5 * time.Second}
for _, s := range services {
if s.URL == "" || s.CheckType != "http" {
continue
}
go func(s models.Service) {
now := time.Now()
start := time.Now()
status := "down"
ms := 0
resp, err := client.Get(s.URL)
if err == nil {
resp.Body.Close()
status = "up"
ms = int(time.Since(start).Milliseconds())
}
db.DB.Model(&models.Service{}).Where("id = ?", s.ID).Updates(map[string]interface{}{
"status":          status,
"last_checked_at": now,
"response_ms":     ms,
})
db.DB.Create(&models.ServicePing{
ServiceID:  s.ID,
Status:     status,
ResponseMs: ms,
CheckedAt:  now,
})
metrics.HealthChecksTotal.WithLabelValues(status).Inc()
metrics.HealthCheckDuration.Observe(time.Since(start).Seconds())
if status != s.Status {
handlers.FireWebhooks(s.Name, status)
}
}(s)
}
updateGauges()
}

func updateGauges() {
var up, down int64
db.DB.Model(&models.Service{}).Where("status = ?", "up").Count(&up)
db.DB.Model(&models.Service{}).Where("status = ?", "down").Count(&down)
metrics.ServicesUp.Set(float64(up))
metrics.ServicesDown.Set(float64(down))
}

func cleanupOldPings() {
cutoff := time.Now().Add(-pingRetention)
res := db.DB.Where("checked_at < ?", cutoff).Delete(&models.ServicePing{})
if res.RowsAffected > 0 {
log.Printf("cleaned up %d old pings", res.RowsAffected)
}
}