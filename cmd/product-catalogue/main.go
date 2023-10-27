package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/vincentcent6/product-catalogue/pkg/config"
	dbClient "github.com/vincentcent6/product-catalogue/pkg/database"
	httpRoutes "github.com/vincentcent6/product-catalogue/pkg/http"
)

func main() {
	err := config.Init()
	if err != nil {
		log.Fatalln("failed to init config", err.Error())
	}
	err = dbClient.InitConnection()
	if err != nil {
		log.Fatalln("failed to init database", err.Error())
	}
	// init http module
	httpRoutes.Init()

	// assign all routes
	r := httpRoutes.NewRoutes()

	// start application
	fmt.Println("Application started at port :3000")
	http.ListenAndServe(":3000", r)
}
