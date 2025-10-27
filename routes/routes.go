package routes

import (
	"github.com/henriquefac/backend_go/models/data_models"
	"github.com/henriquefac/backend_go/repositories"

	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

// função para receber contexto e criar usuário

func SetupRouter(router *gin.Engine) {
	userGroup := router.Group("/users")
	{
		userGroup.POST("/register", createUser)
	}
}

func createUser(c *gin.Context) {
	var request data_models.CreateUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados de entrada inválidos", "details": err.Error()})
		return
	}

	err := repositories.CreateUserFromCreateRequest(&request)

	if err != nil {
		if errors.Is(err, repositories.ErrUserAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar usuário", "details": err.Error()}) // 500 Internal Server Error
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Usuário criado com sucesso!"})
}
