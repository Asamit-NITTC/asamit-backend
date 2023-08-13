package main

import (
	"github.com/Asamit-NITTC/asamit-backend-test/db"
	"github.com/Asamit-NITTC/asamit-backend-test/models"
	"github.com/Asamit-NITTC/asamit-backend-test/router"
)

func main() {
	r := router.NewRouter()
	db.InitalizeDB()
	models.MigrateDB()
	models.InsertDummyData()
	r.Run(":8080")
}
