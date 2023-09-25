package main

import (
	"fmt"
	"os"

	"github.com/Asamit-NITTC/asamit-backend-test/db"
	"github.com/Asamit-NITTC/asamit-backend-test/router"
)

func main() {
	db, sqlDB := db.InitalizeDB()
	defer sqlDB.Close()
	//models.MigrateDB(db)
	//models.InsertDummyData(db)
	r := router.NewRouter(db)
	port := fmt.Sprintf(":%s", os.Getenv("PORT"))
	r.Run(port)
}
