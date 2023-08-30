package db

import (
	"database/sql"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitalizeDB() (*gorm.DB, *sql.DB) {
	//name := os.Getenv("PLANET_SCALE_USER_NAME")
	//password := os.Getenv("PLANET_SCALE_USER_PASSWORD")
	//ip := os.Getenv("PLANET_SCALE_IP")
	//dsn := fmt.Sprintf("%s:%s@tcp(%s)/koyofes2023-reception?tls=True", name, password, ip)
	dsn := "root:password@tcp(mysql:3306)/asamit"
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
