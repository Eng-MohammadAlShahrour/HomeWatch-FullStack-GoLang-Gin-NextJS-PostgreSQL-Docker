package handlers

import (
"net/http"

"homelab/internal/db"
"homelab/internal/models"

"github.com/gin-gonic/gin"
)

func ListDevices(c *gin.Context) {
var devices []models.Device
db.DB.Preload("Services").Order("id desc").Find(&devices)
c.JSON(http.StatusOK, devices)
}

func CreateDevice(c *gin.Context) {
var d models.Device
if err := c.ShouldBindJSON(&d); err != nil {
c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
return
}
if err := db.DB.Create(&d).Error; err != nil {
c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
return
}
c.JSON(http.StatusCreated, d)
}

func UpdateDevice(c *gin.Context) {
id := c.Param("id")
var d models.Device
if err := db.DB.First(&d, id).Error; err != nil {
c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
return
}
var in models.Device
if err := c.ShouldBindJSON(&in); err != nil {
c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
return
}
d.Name, d.Type, d.IP, d.Location, d.Tags = in.Name, in.Type, in.IP, in.Location, in.Tags
db.DB.Save(&d)
c.JSON(http.StatusOK, d)
}

func DeleteDevice(c *gin.Context) {
id := c.Param("id")
db.DB.Where("device_id = ?", id).Delete(&models.Service{})
db.DB.Delete(&models.Device{}, id)
c.JSON(http.StatusOK, gin.H{"ok": true})
}