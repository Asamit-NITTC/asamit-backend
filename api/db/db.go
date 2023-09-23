package db

import (
	"database/sql"
	"log"
	"os"
	"time"

	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	WaitTime = 1
)

func InitalizeDB() (*gorm.DB, *sql.DB) {
	user := os.Getenv("CLOUD_SQL_USER_NAME")
	pass := os.Getenv("CLOUD_SQL_PASSWORD")
	ip := os.Getenv("CLOUD_SQL_IP")
	port := os.Getenv("CLOUD_SQL_PORT")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/asamit", user, pass, ip, port)
	var err error
	var db *gorm.DB
	for i := 0; i < 10; i++ {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		log.Printf("retrying to connect to db: %d", i)
		time.Sleep(WaitTime * time.Second)
	}
	if err != nil {
		log.Fatal(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}

	return db, sqlDB
}
