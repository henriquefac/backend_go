package repositories

import (
	"github.com/henriquefac/backend_go/database"
	"github.com/henriquefac/backend_go/models/data_models"
	"github.com/henriquefac/backend_go/models/db_models"
	"github.com/henriquefac/backend_go/utils"
	"gorm.io/gorm"
	"time"
)

// Fnção para passar de uma requisição para criar um usuário para
// Um objeto db_model.user para salvar no banco de dados
func CreateFromUserRequest(user *data_models.CreateUserRequest) error {
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	var db_user db_models.User

	// deve receber somente structs que fo

	db_user.Name = user.Name
	db_user.Email = user.Email
	db_user.Phone = user.Phone
	db_user.BirthDate = user.BirthDate
	db_user.Password = hashedPassword
	// Default
	db_user.RegisterDate = time.Now()
	db_user.Points = 0
	db_user.Level = 1

	if err := database.DB.Create(&db_user).Error; err != nil {
		return err
	}

	return nil
}

// Buscar usuário pelo email
// vai ser usado no endpoint de login
func GetUserByEmail(email *string) (*data_models.PublicUserResponse, error) {
	var db_user db_models.User

	result := database.DB.Where("email = ?", *email).First(&db_user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}

	publicUser := &data_models.PublicUserResponse{
		ID:           db_user.ID,
		Name:         db_user.Name,
		Email:        db_user.Email,
		Phone:        db_user.Phone,
		BirthDate:    db_user.BirthDate,
		RegisterDate: db_user.RegisterDate,
		Points:       db_user.Points,
		Level:        db_user.Level,
	}

	return publicUser, nil
}
