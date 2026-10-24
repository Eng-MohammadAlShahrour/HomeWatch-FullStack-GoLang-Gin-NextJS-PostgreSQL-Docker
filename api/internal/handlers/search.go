package handlers

import (
"net/http"

"homelab/internal/db"
"homelab/internal/models"

"github.com/gin-gonic/gin"
)

func SearchServices(c *gin.Context) {
q := c.Query("q")
if q == "" {
c.JSON(http.StatusOK, gin.H{"services": []models.Service{}, "count": 0})
return
}
pattern := "%" + q + "%"
var services []models.Service
db.DB.Where("LOWER(name) LIKE LOWER(?) OR LOWER(url) LIKE LOWER(?)", pattern, pattern).
Order("id desc").Limit(50).Find(&services)
c.JSON(http.StatusOK, gin.H{"services": services, "count": len(services)})
}

func SearchDevices(c *gin.Context) {
q := c.Query("q")
if q == "" {
c.JSON(http.StatusOK, gin.H{"devices": []models.Device{}, "count": 0})
return
}
pattern := "%" + q + "%"
var devices []models.Device
db.DB.Where("LOWER(name) LIKE LOWER(?) OR LOWER(ip) LIKE LOWER(?) OR LOWER(location) LIKE LOWER(?)",
pattern, pattern, pattern).Order("id desc").Limit(50).Find(&devices)
c.JSON(http.StatusOK, gin.H{"devices": devices, "count": len(devices)})
}