package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/quocphan74/gingo.git/database"
	"github.com/quocphan74/gingo.git/routes"
)

func main() {
	err := godotenv.Load()
	database.ConnectDB()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	routes := routes.SetupRoutes()

	port := os.Getenv("PORT")
	routes.Run(":" + port)

}
