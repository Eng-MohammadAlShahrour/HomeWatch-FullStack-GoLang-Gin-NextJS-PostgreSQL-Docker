package handlers

import (
"bytes"
"encoding/csv"
"net/http"
"strconv"
"time"

"homelab/internal/db"
"homelab/internal/models"

"github.com/gin-gonic/gin"
)

func ExportServicesCSV(c *gin.Context) {
var services []models.Service
db.DB.Order("id asc").Find(&services)

var buf bytes.Buffer
w := csv.NewWriter(&buf)
_ = w.Write([]string{"id", "name", "url", "device_id", "check_type", "check_interval", "status", "response_ms", "last_checked_at", "created_at"})
for _, s := range services {
last := ""
if s.LastCheckedAt != nil {
last = s.LastCheckedAt.Format(time.RFC3339)
}
_ = w.Write([]string{
strconv.Itoa(int(s.ID)), s.Name, s.URL,
strconv.Itoa(int(s.DeviceID)), s.CheckType,
strconv.Itoa(s.CheckInterval), s.Status,
strconv.Itoa(s.ResponseMs), last,
s.CreatedAt.Format(time.RFC3339),
})
}
w.Flush()

c.Header("Content-Type", "text/csv")
c.Header("Content-Disposition", "attachment; filename=services.csv")
c.String(http.StatusOK, buf.String())
}

func ExportDevicesCSV(c *gin.Context) {
var devices []models.Device
db.DB.Order("id asc").Find(&devices)

var buf bytes.Buffer
w := csv.NewWriter(&buf)
_ = w.Write([]string{"id", "name", "type", "ip", "location", "tags", "created_at"})
for _, d := range devices {
_ = w.Write([]string{
strconv.Itoa(int(d.ID)), d.Name, d.Type, d.IP, d.Location, d.Tags,
d.CreatedAt.Format(time.RFC3339),
})
}
w.Flush()

c.Header("Content-Type", "text/csv")
c.Header("Content-Disposition", "attachment; filename=devices.csv")
c.String(http.StatusOK, buf.String())
}

func ExportPingsCSV(c *gin.Context) {
id := c.Param("id")
var pings []models.ServicePing
db.DB.Where("service_id = ?", id).Order("checked_at asc").Find(&pings)

var buf bytes.Buffer
w := csv.NewWriter(&buf)
_ = w.Write([]string{"checked_at", "status", "response_ms"})
for _, p := range pings {
_ = w.Write([]string{p.CheckedAt.Format(time.RFC3339), p.Status, strconv.Itoa(p.ResponseMs)})
}
w.Flush()

c.Header("Content-Type", "text/csv")
c.Header("Content-Disposition", "attachment; filename=service_"+id+"_pings.csv")
c.String(http.StatusOK, buf.String())
}