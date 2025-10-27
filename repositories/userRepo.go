package repositories

import (
	"errors"
	"github.com/henriquefac/backend_go/database"
	"github.com/henriquefac/backend_go/models/data_models"
	"github.com/henriquefac/backend_go/models/db_models"
	"github.com/henriquefac/backend_go/utils"
	"gorm.io/gorm"
	"strings"
	"time"
)

// reporisório que deve criar no banco de dados o usuário
// a partir de data_models.CreateUserRequest

var ErrUserAlreadyExists = errors.New("Usuário já registrado com esse email")
var ErrUserNotFound = errors.New("Usuário não encontrado")

func CreateUserFromCreateRequest(createRequest *data_models.CreateUserRequest) error {
	hashedPassword, err := utils.HashPassword(&createRequest.Password)

	if err != nil {
		return err
	}

	userDB := db_models.User{
		Name:         createRequest.Name,
		Password:     hashedPassword,
		RegisterDate: time.Now(),
		Email:        createRequest.Email,
		Phone:        createRequest.Phone,
		BirthDate:    createRequest.BirthDate,
		Points:       0,
		Level:        1,
	}

	result := database.DB.Create(&userDB)

	if result.Error != nil {
		// Trata erro de duplicidade (MySQL/MariaDB)
		if strings.Contains(result.Error.Error(), "duplicate entry value") {
			return ErrUserAlreadyExists
		}
		return result.Error
	}

	return nil
}

// procurar usuário pelo email
// Retornar data_models.PublicUserResponse
func GetUserByEmai(email string) (*data_models.PublicUserResponse, error) {
	var userDB db_models.User

	result := database.DB.Where("email = ?", email).First(&userDB)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}

		return nil, result.Error
	}

	response := &data_models.PublicUserResponse{
		Name:           userDB.Name,
		Email:          userDB.Email,
		RegisterDate:   userDB.RegisterDate,
		Phone:          userDB.Phone,
		ProfilePicture: userDB.ProfilePicture,
		BirthDate:      userDB.BirthDate,
		Points:         userDB.Points,
		Level:          userDB.Level,
	}

	return response, nil
}
