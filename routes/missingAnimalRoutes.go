package routes

import (
	"github.com/henriquefac/backend_go/database"
	"github.com/henriquefac/backend_go/models/data_models"
	"github.com/henriquefac/backend_go/repositories"
	"github.com/henriquefac/backend_go/services"

	"github.com/gin-gonic/gin"
	"net/http"
)

func SetupMissingAnimalRouter(router *gin.Engine) {
	missingGroup := router.Group("/missingAnimal")
	{
		missingGroup.POST("/create", Create)

	}
}

func Create(c *gin.Context) {
	var request data_models.MissingAnimalRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados de entrada inválidos",
			"details": err.Error()})
		return
	}

	missingRepo := repositories.NewMissingAnimalRepository(database.DB)
	missingService := services.NewMissingAnimalService(missingRepo)

	response, err := missingService.Create(&request)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar um registro de animal perdido",
			"details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response)
}

func Update(c *gin.Context) {
	var request data_models.MissingAnimalUpdateRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados de entrada inválidos",
			"details": err.Error()})
		return
	}

	missingRepo := repositories.NewMissingAnimalRepository(database.DB)
	missingService := services.NewMissingAnimalService(missingRepo)

	response, err := missingService.Update()
}
