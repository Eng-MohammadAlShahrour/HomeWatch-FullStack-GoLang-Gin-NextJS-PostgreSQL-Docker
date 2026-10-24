package handlers

import (
"net/http"

"homelab/internal/auth"
"homelab/internal/config"
"homelab/internal/db"
"homelab/internal/models"

"github.com/gin-gonic/gin"
)

type authInput struct {
Email    string `json:"email" binding:"required,email"`
Password string `json:"password" binding:"required,min=6"`
}

var cfg = config.Load()

func Register(c *gin.Context) {
var in authInput
if err := c.ShouldBindJSON(&in); err != nil {
c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
return
}
var existing models.User
if err := db.DB.Where("email = ?", in.Email).First(&existing).Error; err == nil {
c.JSON(http.StatusConflict, gin.H{"error": "email already exists"})
return
}
hash, _ := auth.HashPassword(in.Password)
u := models.User{Email: in.Email, PasswordHash: hash}
if err := db.DB.Create(&u).Error; err != nil {
c.JSON(http.StatusInternalServerError, gin.H{"error": "create failed"})
return
}
token, _ := auth.GenerateToken(cfg.JWTSecret, u.ID, u.Email)
c.JSON(http.StatusCreated, gin.H{"token": token, "user": u})
}

func Login(c *gin.Context) {
var in authInput
if err := c.ShouldBindJSON(&in); err != nil {
c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
return
}
var u models.User
if err := db.DB.Where("email = ?", in.Email).First(&u).Error; err != nil {
c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
return
}
if !auth.CheckPassword(u.PasswordHash, in.Password) {
c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
return
}
token, _ := auth.GenerateToken(cfg.JWTSecret, u.ID, u.Email)
c.JSON(http.StatusOK, gin.H{"token": token, "user": u})
}

func Me(c *gin.Context) {
id := c.GetUint("userID")
var u models.User
if err := db.DB.First(&u, id).Error; err != nil {
c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
return
}
c.JSON(http.StatusOK, u)
}