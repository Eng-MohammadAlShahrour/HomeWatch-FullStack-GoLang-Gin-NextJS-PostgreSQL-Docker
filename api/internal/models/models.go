package models

import "time"

type User struct {
ID           uint      `gorm:"primaryKey" json:"id"`
Email        string    `gorm:"uniqueIndex;not null" json:"email"`
PasswordHash string    `gorm:"not null" json:"-"`
CreatedAt    time.Time `json:"created_at"`
}

type Device struct {
ID        uint      `gorm:"primaryKey" json:"id"`
Name      string    `gorm:"not null" json:"name"`
Type      string    `json:"type"`
IP        string    `json:"ip"`
Location  string    `json:"location"`
Tags      string    `json:"tags"`
CreatedAt time.Time `json:"created_at"`
Services  []Service `json:"services,omitempty"`
}

type Service struct {
ID            uint       `gorm:"primaryKey" json:"id"`
DeviceID      uint       `json:"device_id"`
Name          string     `gorm:"not null" json:"name"`
URL           string     `json:"url"`
CheckType     string     `gorm:"default:'http'" json:"check_type"`
CheckInterval int        `gorm:"default:30" json:"check_interval"`
Status        string     `gorm:"default:'unknown'" json:"status"`
LastCheckedAt *time.Time `json:"last_checked_at"`
ResponseMs    int        `json:"response_ms"`
CreatedAt     time.Time  `json:"created_at"`
}

// ServicePing stores each health check result for uptime history.
type ServicePing struct {
ID         uint      `gorm:"primaryKey" json:"id"`
ServiceID  uint      `gorm:"index:idx_service_time,priority:1" json:"service_id"`
Status     string    `json:"status"`
ResponseMs int       `json:"response_ms"`
CheckedAt  time.Time `gorm:"index:idx_service_time,priority:2" json:"checked_at"`
}
type Webhook struct {
ID        uint      `gorm:"primaryKey" json:"id"`
Name      string    `gorm:"not null" json:"name"`
URL       string    `gorm:"not null" json:"url"`
Events    string    `gorm:"default:'down'" json:"events"`
Enabled   bool      `gorm:"default:true" json:"enabled"`
CreatedAt time.Time `json:"created_at"`
}