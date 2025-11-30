package main

import (
	"github.com/gin-gonic/gin"
	"github.com/henriquefac/backend_go/database"
	"github.com/henriquefac/backend_go/routes"
	"log"
)

func main() {
	database.Connect()

	router := gin.Default()

	routes.SetupUserRouter(router)
	routes.SetupMissingAnimalRouter(router)

	log.Println("Sevidor iniciado na porta 8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}
