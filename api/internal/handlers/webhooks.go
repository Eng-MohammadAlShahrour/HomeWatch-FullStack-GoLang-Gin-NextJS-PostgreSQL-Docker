package handlers

import (
"bytes"
"encoding/json"
"net/http"
"strings"
"time"

"homelab/internal/db"
"homelab/internal/metrics"
"homelab/internal/models"

"github.com/gin-gonic/gin"
)

type webhookInput struct {
Name   string `json:"name" binding:"required"`
URL    string `json:"url" binding:"required,url"`
Events string `json:"events"`
}

func ListWebhooks(c *gin.Context) {
var hooks []models.Webhook
db.DB.Order("id desc").Find(&hooks)
c.JSON(http.StatusOK, hooks)
}

func CreateWebhook(c *gin.Context) {
var in webhookInput
if err := c.ShouldBindJSON(&in); err != nil {
c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
return
}
if in.Events == "" {
in.Events = "down"
}
w := models.Webhook{Name: in.Name, URL: in.URL, Events: in.Events, Enabled: true}
if err := db.DB.Create(&w).Error; err != nil {
c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
return
}
c.JSON(http.StatusCreated, w)
}

func DeleteWebhook(c *gin.Context) {
id := c.Param("id")
db.DB.Delete(&models.Webhook{}, id)
c.JSON(http.StatusOK, gin.H{"ok": true})
}

func FireWebhooks(serviceName, status string) {
var hooks []models.Webhook
db.DB.Where("enabled = ?", true).Find(&hooks)
for _, h := range hooks {
if !eventMatches(h.Events, status) {
continue
}
go sendWebhook(h, serviceName, status)
}
}

func eventMatches(events, status string) bool {
for _, e := range strings.Split(events, ",") {
if strings.TrimSpace(e) == status {
return true
}
}
return false
}

func sendWebhook(h models.Webhook, serviceName, status string) {
payload := map[string]interface{}{
"service": serviceName,
"status":  status,
"time":    time.Now().Format(time.RFC3339),
"text":    "HomeWatch: " + serviceName + " is " + status,
}
b, _ := json.Marshal(payload)
client := &http.Client{Timeout: 10 * time.Second}
resp, err := client.Post(h.URL, "application/json", bytes.NewBuffer(b))
if err != nil {
metrics.WebhooksSent.WithLabelValues("error").Inc()
return
}
resp.Body.Close()
metrics.WebhooksSent.WithLabelValues("ok").Inc()
}