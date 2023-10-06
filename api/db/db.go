package db

import (
	"database/sql"
	"log"
	"os"

	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitializeDB() (*gorm.DB, *sql.DB) {
	user := os.Getenv("CLOUD_SQL_USER_NAME")
	pass := os.Getenv("CLOUD_SQL_PASSWORD")
	ip := os.Getenv("CLOUD_SQL_IP")
	port := os.Getenv("CLOUD_SQL_PORT")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/asamit?charset=utf8mb4&parseTime=True&loc=Asia%2FTokyo", user, pass, ip, port)
	var err error
	var db *gorm.DB
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}

	return db, sqlDB
}

func InitializeLocalDB() (*gorm.DB, *sql.DB) {
	dsn := "root:password@tcp(mysql)/asamit?charset=utf8mb4&parseTime=True&loc=Asia%2FTokyo"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}

	return db, sqlDB
}
