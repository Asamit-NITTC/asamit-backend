package main

import (
	"fmt"
	"os"

	"github.com/Asamit-NITTC/asamit-backend-test/db"
	"github.com/Asamit-NITTC/asamit-backend-test/models"
	"github.com/Asamit-NITTC/asamit-backend-test/router"
)

func main() {
	if os.Getenv("MODE") == "DEBUG" {
		db, sqlDB := db.InitializeLocalDB()
		defer sqlDB.Close()
		models.MigrateDB(db)
		models.InsertDummyData(db)

		r := router.NewRouter(db)
		port := fmt.Sprintf(":%s", os.Getenv("PORT"))
		r.Run(port)
	} else {
		db, sqlDB := db.InitializeDB()
		defer sqlDB.Close()
		r := router.NewRouter(db)
		port := fmt.Sprintf(":%s", os.Getenv("PORT"))
		r.Run(port)
	}

}
