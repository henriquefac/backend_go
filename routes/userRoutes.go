package routes

import (
	"github.com/henriquefac/backend_go/database"
	"github.com/henriquefac/backend_go/models/data_models"
	"github.com/henriquefac/backend_go/repositories"
	"github.com/henriquefac/backend_go/services"

	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

// função para receber contexto e criar usuário

func SetupUserRouter(router *gin.Engine) {
	userGroup := router.Group("/users")
	{
		userGroup.POST("/signup", signupUser)
		userGroup.GET("/login", loginUser)
	}
}

func signupUser(c *gin.Context) {
	var request data_models.CreateUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados de entrada inválidos", "details": err.Error()})
		return
	}

	userRepo := repositories.NewUserRepository(database.DB)

	err := userRepo.CreateUserFromCreateRequest(&request)

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

// Handller para busca do usuário

func loginUser(c *gin.Context) {
	var request data_models.LoginRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Dados de entrada inválidos",
			"Datails": err.Error(),
		})
		return
	}

	userRepo := repositories.NewUserRepository(database.DB)
	userServices := services.NewUserService(userRepo)

	response, err := userServices.Login(request.Email, request.Password)

	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) || errors.Is(err, services.ErrInvalidPassword) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Email ou senha incorretos",
			})
			return
		}
		c.JSON(http.StatusInternalServerError,
			gin.H{
				"error":   "Erro ao realizar login",
				"details": err.Error(),
			})
		return
	}

	c.JSON(http.StatusOK, response)

}
