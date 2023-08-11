package main

import (
	"github.com/Asamit-NITTC/asamit-backend-test/router"
)

func main() {
	r := router.NewRouter()
	r.Run(":8080")
}
