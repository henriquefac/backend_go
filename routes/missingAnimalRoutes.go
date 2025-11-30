package routes

import (
	"github.com/henriquefac/backend_go/database"
	"github.com/henriquefac/backend_go/models/data_models"
	"github.com/henriquefac/backend_go/repositories"
	"github.com/henriquefac/backend_go/services"

	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func SetupMissingAnimalRouter(router *gin.Engine) {
	missingGroup := router.Group("/missingAnimal")
	{
		missingGroup.POST("/create", Create)
		missingGroup.PUT("/update", Update)
		missingGroup.GET("listAll", ListAll)
		missingGroup.GET("/listByUser/:userId", ListByUserID)

	}
}

func Create(c *gin.Context) {
	var request data_models.MissingAnimalCreateRequest

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

	response, err := missingService.Update(&request)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Dados d entrada inválidos",
			"details": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, response)
}

// apagar registro

// receber lista de registros de missinanimal baseado no id do usuário

func ListByUserID(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)

	if err != nil || userID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de usuário inválido na URL."})
		return
	}

	missingRepo := repositories.NewMissingAnimalRepository(database.DB)

	response, err := missingRepo.ListUserMissingAnimals(uint(userID))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao listar registros",
			"details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)

}

func ListAll(c *gin.Context) {
	missingRepo := repositories.NewMissingAnimalRepository(database.DB)

	response, err := missingRepo.ListAllMissingAnimals()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao listar registros",
			"details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
