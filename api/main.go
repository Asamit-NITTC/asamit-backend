package main

import (
	"fmt"
	"os"

	"github.com/Asamit-NITTC/asamit-backend-test/db"
	"github.com/Asamit-NITTC/asamit-backend-test/models"
	"github.com/Asamit-NITTC/asamit-backend-test/router"
)

func main() {
	db, sqlDB := db.InitalizeDB()
	defer sqlDB.Close()
	r := router.NewRouter(db)
	models.MigrateDB(db)
	models.InsertDummyData(db)
	port := fmt.Sprintf(":%s", os.Getenv("PORT"))
	r.Run(port)
}
