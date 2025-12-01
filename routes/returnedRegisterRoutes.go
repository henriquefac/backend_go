package routes

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/henriquefac/backend_go/database"
	"github.com/henriquefac/backend_go/models/data_models"
	"github.com/henriquefac/backend_go/repositories"
	"github.com/henriquefac/backend_go/services"
)

func SetupReturnedRegisterRouter(router *gin.Engine) {
	returnedGroup := router.Group("/returnedRegister")
	{
		returnedGroup.POST("/create", CreateReturned)
		returnedGroup.GET("/animal/:animalId", GetReturnedByAnimalId)
		returnedGroup.GET("/user/:userId", LIstReturnedRegisterByUser)
	}
}

func CreateReturned(c *gin.Context) {
	var request data_models.ReturnedRegisterCreateRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados de entrada inválidos",
			"details": err.Error()})
		return
	}

	// 1. Inicializa Repositórios (Dependências)
	returnedRepo := repositories.NewReturnedRegisterRepository(database.DB)
	missingRepo := repositories.NewMissingAnimalRepository(database.DB)

	// 2. Inicializa o Serviço
	returnedService := services.NewReturnedRegisterService(returnedRepo, missingRepo)

	// 3. Executa a Lógica de Negócio (Criação e Atualização de Status)
	response, err := returnedService.Create(&request)

	if err != nil {
		// Trata erros de integridade referencial e duplicação

		if errors.Is(err, repositories.ErrMissingAnimalNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Animal não encontrado ou já apagado."})
			return
		}
		if errors.Is(err, repositories.ErrReturnedRegisterExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "Este animal já possui um registro de retorno."})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao registrar retorno do animal",
			"details": err.Error()})
		return
	}

	// 4. Resposta de Sucesso
	c.JSON(http.StatusCreated, response)
}

func GetReturnedByAnimalId(c *gin.Context) {
	animalIDStr := c.Param("animalId")
	animalID, err := strconv.ParseUint(animalIDStr, 10, 32)

	if err != nil || animalID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de registro de animal inválido na URL"})
		return
	}

	// 1. Inicializa Repositórios (Dependências)
	returnedRepo := repositories.NewReturnedRegisterRepository(database.DB)
	missingRepo := repositories.NewMissingAnimalRepository(database.DB)

	// 2. Inicializa o Serviço
	returnedService := services.NewReturnedRegisterService(returnedRepo, missingRepo)

	// 3. Executa a busca
	// O service retorna nil se o registro não for encontrado.
	response, err := returnedService.GetReturnedRegisterByAnimalID(uint(animalID))

	// 4. Tratamento de Erro do Serviço
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro interno ao buscar registro de retorno",
			"details": err.Error()})
		return
	}

	// 5. Tratamento de Not Found (O serviço retorna 'nil' se o registro não existir)
	if response == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Registro de retorno não encontrado para este animal."})
		return
	}

	// 6. Resposta de Sucesso
	c.JSON(http.StatusOK, response)
}

func LIstReturnedRegisterByUser(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)

	if err != nil || userID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de usuário inválido na URL."})
		return
	}

	registerRepo := repositories.NewReturnedRegisterRepository(database.DB)

	response, err := registerRepo.ListReturnedRegisterFromUser(uint(userID))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao listar registros",
			"details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
