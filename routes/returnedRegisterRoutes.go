package routes

import (
	"github.com/henriquefac/backend_go/database"
	"github.com/henriquefac/backend_go/models/data_models"
	"github.com/henriquefac/backend_go/repositories"
	"github.com/henriquefac/backend_go/services"

	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func SetupReturnedRegisterRouter(router *gin.Engine) {
	returnedGroup := router.Group("/returnedRegister")
	{
		returnedGroup.POST("/create", CreateReturned)

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
	missingRepo := repositories.NewMissingAnimalRepository(database.DB) // Necessário para o Service/Transação

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
	// Status 201 Created é apropriado para a criação de um novo recurso.
	c.JSON(http.StatusCreated, response)
}
