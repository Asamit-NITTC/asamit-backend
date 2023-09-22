package db

import (
	"database/sql"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	WaitTime = 1
)

func InitalizeDB() (*gorm.DB, *sql.DB) {
	//name := os.Getenv("PLANET_SCALE_USER_NAME")
	//password := os.Getenv("PLANET_SCALE_USER_PASSWORD")
	//ip := os.Getenv("PLANET_SCALE_IP")
	//dsn := fmt.Sprintf("%s:%s@tcp(%s)/koyofes2023-reception?tls=True", name, password, ip)
	dsn := "root:password@tcp(mysql)/asamit"
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
