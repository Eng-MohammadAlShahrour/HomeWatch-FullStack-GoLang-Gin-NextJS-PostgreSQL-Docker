package handlers

import (
"net/http"

"homelab/internal/db"
"homelab/internal/models"

"github.com/gin-gonic/gin"
)

func Dashboard(c *gin.Context) {
var devices, services, up, down int64
db.DB.Model(&models.Device{}).Count(&devices)
db.DB.Model(&models.Service{}).Count(&services)
db.DB.Model(&models.Service{}).Where("status = ?", "up").Count(&up)
db.DB.Model(&models.Service{}).Where("status = ?", "down").Count(&down)
c.JSON(http.StatusOK, gin.H{
"devices":  devices,
"services": services,
"up":       up,
"down":     down,
})
}