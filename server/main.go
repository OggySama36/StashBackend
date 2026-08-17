package main

import (
	"Stash/config"
	"Stash/internal/router"
	"fmt"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	config.Load()
	config.Init()
	r := router.SetupRoutes()
	port := fmt.Sprintf(":%s", config.App.Port)
	r.Run(port)
}
