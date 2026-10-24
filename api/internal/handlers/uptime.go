package handlers

import (
"net/http"
"strconv"
"time"

"homelab/internal/db"
"homelab/internal/models"

"github.com/gin-gonic/gin"
)

type bucketPoint struct {
Time     time.Time `json:"time"`
AvgMs    float64   `json:"avg_ms"`
UpCount  int64     `json:"up_count"`
DownCount int64    `json:"down_count"`
}

// GetServiceUptime returns uptime % + time-bucketed series for one service.
// Query: ?hours=24 (default 24, max 168)
func GetServiceUptime(c *gin.Context) {
id, err := strconv.Atoi(c.Param("id"))
if err != nil {
c.JSON(http.StatusBadRequest, gin.H{"error": "invalid service id"})
return
}

hours := 24
if h := c.Query("hours"); h != "" {
if parsed, err := strconv.Atoi(h); err == nil && parsed > 0 && parsed <= 168 {
hours = parsed
}
}
since := time.Now().Add(-time.Duration(hours) * time.Hour)

var total, up int64
db.DB.Model(&models.ServicePing{}).
Where("service_id = ? AND checked_at >= ?", id, since).
Count(&total)
db.DB.Model(&models.ServicePing{}).
Where("service_id = ? AND checked_at >= ? AND status = ?", id, since, "up").
Count(&up)

uptimePct := 0.0
if total > 0 {
uptimePct = float64(up) / float64(total) * 100
}

// 5-min buckets
var points []bucketPoint
rows, err := db.DB.Model(&models.ServicePing{}).
Select("date_trunc('hour', checked_at) + (EXTRACT(minute FROM checked_at)::int / 5) * INTERVAL '5 minutes' AS bucket, AVG(response_ms) AS avg_ms, COUNT(*) FILTER (WHERE status = 'up') AS up_count, COUNT(*) FILTER (WHERE status = 'down') AS down_count").
Where("service_id = ? AND checked_at >= ?", id, since).
Group("bucket").
Order("bucket ASC").
Rows()
if err == nil {
defer rows.Close()
for rows.Next() {
var b bucketPoint
if err := rows.Scan(&b.Time, &b.AvgMs, &b.UpCount, &b.DownCount); err == nil {
points = append(points, b)
}
}
}
if points == nil {
points = []bucketPoint{}
}

c.JSON(http.StatusOK, gin.H{
"service_id":  id,
"hours":       hours,
"total":       total,
"up":          up,
"down":        total - up,
"uptime_pct":  uptimePct,
"series":      points,
})
}

// GetUptimeSummary returns uptime % for every service.
func GetUptimeSummary(c *gin.Context) {
var services []models.Service
db.DB.Find(&services)

since := time.Now().Add(-24 * time.Hour)
type row struct {
ServiceID uint    `json:"service_id"`
Name      string  `json:"name"`
Total     int64   `json:"total"`
Up        int64   `json:"up"`
UptimePct float64 `json:"uptime_pct"`
AvgMs     float64 `json:"avg_ms"`
}
var out []row

for _, s := range services {
var total, up int64
db.DB.Model(&models.ServicePing{}).
Where("service_id = ? AND checked_at >= ?", s.ID, since).Count(&total)
db.DB.Model(&models.ServicePing{}).
Where("service_id = ? AND checked_at >= ? AND status = ?", s.ID, since, "up").Count(&up)

pct := 0.0
if total > 0 {
pct = float64(up) / float64(total) * 100
}

var avg float64
db.DB.Model(&models.ServicePing{}).
Where("service_id = ? AND checked_at >= ? AND status = ?", s.ID, since, "up").
Select("COALESCE(AVG(response_ms), 0)").Scan(&avg)

out = append(out, row{
ServiceID: s.ID,
Name:      s.Name,
Total:     total,
Up:        up,
UptimePct: pct,
AvgMs:     avg,
})
}
if out == nil {
out = []row{}
}

c.JSON(http.StatusOK, out)
}