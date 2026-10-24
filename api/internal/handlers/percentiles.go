package handlers

import (
"net/http"
"sort"
"strconv"
"time"

"homelab/internal/db"
"homelab/internal/models"

"github.com/gin-gonic/gin"
)

func GetLatencyPercentiles(c *gin.Context) {
id := c.Param("id")
hours := 24
if h := c.Query("hours"); h != "" {
if v, err := strconv.Atoi(h); err == nil && v > 0 && v <= 168 {
hours = v
}
}
since := time.Now().Add(-time.Duration(hours) * time.Hour)

var pings []models.ServicePing
db.DB.Where("service_id = ? AND checked_at >= ? AND status = ?", id, since, "up").Find(&pings)

if len(pings) == 0 {
c.JSON(http.StatusOK, gin.H{"service_id": id, "count": 0})
return
}

vals := make([]int, len(pings))
for i, p := range pings {
vals[i] = p.ResponseMs
}
sort.Ints(vals)

pct := func(p float64) int {
idx := int(float64(len(vals)-1) * p)
if idx < 0 {
idx = 0
}
if idx >= len(vals) {
idx = len(vals) - 1
}
return vals[idx]
}

c.JSON(http.StatusOK, gin.H{
"service_id": id,
"hours":      hours,
"count":      len(vals),
"min":        vals[0],
"max":        vals[len(vals)-1],
"p50":        pct(0.50),
"p75":        pct(0.75),
"p95":        pct(0.95),
"p99":        pct(0.99),
})
}