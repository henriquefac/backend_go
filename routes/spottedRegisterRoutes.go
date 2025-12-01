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

func SetupSpottedAnimalRegisterRouter(router *gin.Engine) {
	spottedGroup := router.Group("/spottedRegister")
	{
		spottedGroup.POST("/create", CreateSpotted)
		spottedGroup.GET("/listByAnimal/:animalId", ListSpottedByAnimal)
		spottedGroup.GET("/listByUser/:userId", ListSpottedByUserID)
	}
}

func CreateSpotted(c *gin.Context) {
	var request data_models.SpottedRegisterCreateRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados de entrada inválidos",
			"details": err.Error()})
		return
	}

	spottedRepo := repositories.NewSpottedRegisterRepository(database.DB)
	spottedServie := services.NewSpottedRegisterService(spottedRepo)

	response, err := spottedServie.Create(&request)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar um registro de animal avistado",
			"details": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, response)
}

func ListSpottedByAnimal(c *gin.Context) {
	animalIDStr := c.Param("animalId")
	animalID, err := strconv.ParseUint(animalIDStr, 10, 32)

	if err != nil || animalID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de registro de animla perdido inváldo na URL"})
		return
	}

	spottedRepo := repositories.NewSpottedRegisterRepository(database.DB)

	response, err := spottedRepo.ListSpottedRegistersByAnimalID(uint(animalID))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao listar registros",
			"details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func ListSpottedByUserID(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)

	if err != nil || userID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de usuário inválido na URL."})
		return
	}

	spottedRepo := repositories.NewSpottedRegisterRepository(database.DB)

	response, err := spottedRepo.ListSpottedRegisterByUserID(uint(userID))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao listar registros",
			"details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
