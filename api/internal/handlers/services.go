package handlers

import (
"net/http"

"homelab/internal/db"
"homelab/internal/models"

"github.com/gin-gonic/gin"
)

func ListServices(c *gin.Context) {
var s []models.Service
db.DB.Order("id desc").Find(&s)
c.JSON(http.StatusOK, s)
}

func CreateService(c *gin.Context) {
var s models.Service
if err := c.ShouldBindJSON(&s); err != nil {
c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
return
}
if s.CheckInterval <= 0 {
s.CheckInterval = 30
}
if s.CheckType == "" {
s.CheckType = "http"
}
s.Status = "unknown"
if err := db.DB.Create(&s).Error; err != nil {
c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
return
}
c.JSON(http.StatusCreated, s)
}

func UpdateService(c *gin.Context) {
id := c.Param("id")
var s models.Service
if err := db.DB.First(&s, id).Error; err != nil {
c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
return
}
var in models.Service
if err := c.ShouldBindJSON(&in); err != nil {
c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
return
}
s.Name, s.URL, s.DeviceID = in.Name, in.URL, in.DeviceID
if in.CheckInterval > 0 {
s.CheckInterval = in.CheckInterval
}
if in.CheckType != "" {
s.CheckType = in.CheckType
}
db.DB.Save(&s)
c.JSON(http.StatusOK, s)
}

func DeleteService(c *gin.Context) {
id := c.Param("id")
db.DB.Delete(&models.Service{}, id)
c.JSON(http.StatusOK, gin.H{"ok": true})
}