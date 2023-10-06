package main

import (
	"fmt"
	"os"

	"github.com/Asamit-NITTC/asamit-backend-test/db"
	"github.com/Asamit-NITTC/asamit-backend-test/models"
	"github.com/Asamit-NITTC/asamit-backend-test/router"
	"github.com/Asamit-NITTC/asamit-backend-test/webstorage"
)

func main() {
	if os.Getenv("MODE") == "DEBUG" {
		db, sqlDB := db.InitializeLocalDB()
		defer sqlDB.Close()
		models.MigrateDB(db)
		models.InsertDummyData(db)

		ctx, bucket := webstorage.InitializeCloudStorage()

		r := router.NewRouter(db, ctx, bucket)
		port := fmt.Sprintf(":%s", os.Getenv("PORT"))
		r.Run(port)
	} else {
		db, sqlDB := db.InitializeDB()
		defer sqlDB.Close()

		ctx, bucket := webstorage.InitializeCloudStorage()
		r := router.NewRouter(db, ctx, bucket)
		port := fmt.Sprintf(":%s", os.Getenv("PORT"))
		r.Run(port)
	}

}
