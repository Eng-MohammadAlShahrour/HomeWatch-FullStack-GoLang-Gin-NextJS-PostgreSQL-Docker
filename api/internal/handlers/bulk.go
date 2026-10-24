package handlers

import (
"net/http"

"homelab/internal/db"
"homelab/internal/models"

"github.com/gin-gonic/gin"
)

type bulkRequest struct {
IDs []uint `json:"ids" binding:"required"`
}

func BulkDeleteServices(c *gin.Context) {
var req bulkRequest
if err := c.ShouldBindJSON(&req); err != nil || len(req.IDs) == 0 {
c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or empty ids"})
return
}
res := db.DB.Delete(&models.Service{}, req.IDs)
c.JSON(http.StatusOK, gin.H{"deleted": res.RowsAffected})
}

func BulkDeleteDevices(c *gin.Context) {
var req bulkRequest
if err := c.ShouldBindJSON(&req); err != nil || len(req.IDs) == 0 {
c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or empty ids"})
return
}
db.DB.Where("device_id IN ?", req.IDs).Delete(&models.Service{})
res := db.DB.Delete(&models.Device{}, req.IDs)
c.JSON(http.StatusOK, gin.H{"deleted": res.RowsAffected})
}

func BulkPingNow(c *gin.Context) {
var req bulkRequest
if err := c.ShouldBindJSON(&req); err != nil || len(req.IDs) == 0 {
c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or empty ids"})
return
}
db.DB.Model(&models.Service{}).Where("id IN ?", req.IDs).Update("last_checked_at", nil)
c.JSON(http.StatusOK, gin.H{"scheduled": len(req.IDs)})
}