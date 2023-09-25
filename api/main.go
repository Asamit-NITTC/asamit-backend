package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Asamit-NITTC/asamit-backend-test/db"
	"github.com/Asamit-NITTC/asamit-backend-test/models"
	"github.com/Asamit-NITTC/asamit-backend-test/router"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	sqlDB, err := db.ConnectWithConnector()
	if err != nil {
		log.Fatal(err)
	}
	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()
	models.MigrateDB(db)
	models.InsertDummyData(db)
	r := router.NewRouter(db)
	port := fmt.Sprintf(":%s", os.Getenv("PORT"))
	r.Run(port)
}
