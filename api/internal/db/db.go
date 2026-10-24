package db

import (
"log"
"time"

"homelab/internal/models"

"gorm.io/driver/postgres"
"gorm.io/gorm"
"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Init(dsn string) {
var err error
for i := 0; i < 20; i++ {
DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
Logger: logger.Default.LogMode(logger.Warn),
})
if err == nil {
break
}
log.Println("waiting for db...", err)
time.Sleep(2 * time.Second)
}
if err != nil {
log.Fatal("db connection failed: ", err)
}
if err := DB.AutoMigrate(
&models.User{},
&models.Device{},
&models.Service{},
&models.ServicePing{},
		&models.Webhook{},
); err != nil {
log.Fatal(err)
}
log.Println("db ready")
}