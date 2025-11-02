package repositories

import (
	"errors"
	"github.com/henriquefac/backend_go/models/data_models"
	"github.com/henriquefac/backend_go/models/db_models"
	"github.com/henriquefac/backend_go/utils"
	"gorm.io/gorm"
	"strings"
	"time"
)

// reporisório que deve criar no banco de dados o usuário
// a partir de data_models.CreateUserRequest

var (
	ErrUserAlreadyExists = errors.New("Usuário já registrado com esse email")
	ErrUserNotFound      = errors.New("Usuário não encontrado")
	ErrInvalidPassword   = errors.New("Senha inválida")
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUserFromCreateRequest(createRequest *data_models.CreateUserRequest) error {
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

	result := r.db.Create(&userDB)

	if result.Error != nil {
		// Trata erro de duplicidade (MySQL/MariaDB)
		if strings.Contains(result.Error.Error(), "duplicate entry value") {
			return ErrUserAlreadyExists
		}
		return result.Error
	}

	return nil
}

func (r *UserRepository) GetUserByEmail(email string,
	publicResponse *data_models.PublicUserResponse) error {
	// modelos do database que representa os Usuários
	var userDB db_models.User

	// buscar pelo email
	result := r.db.Where("email = ?", email).First(&userDB)

	// verificar error (se for erro de Record Not Found, devolver instancia de erro personalizado)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		// Se for outro erro, retorna erro
		return result.Error
	}

	// Completar "publicResponse" com as informações necessárias

	publicResponse.ID = userDB.ID
	publicResponse.Name = userDB.Name
	publicResponse.Email = userDB.Email
	publicResponse.RegisterDate = userDB.RegisterDate
	publicResponse.Phone = userDB.Phone
	publicResponse.ProfilePicture = userDB.ProfilePicture
	publicResponse.BirthDate = userDB.BirthDate
	publicResponse.Points = userDB.Points
	publicResponse.Level = userDB.Level

	publicResponse.Password = userDB.Password

	return nil

}

func (r *UserRepository) LoginByEmailAndPassword(email, password string,
	publicResponse *data_models.PublicUserResponse) error {

	var userDB db_models.User

	result := r.db.Where("email = ?", email).First(&userDB)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return result.Error
	}

	// Verificar se a senha confere
	if !utils.CheckPassword(&userDB.Password, &password) {
		return ErrInvalidPassword
	}

	publicResponse.ID = userDB.ID
	publicResponse.Name = userDB.Name
	publicResponse.Email = userDB.Email
	publicResponse.RegisterDate = userDB.RegisterDate
	publicResponse.Phone = userDB.Phone
	publicResponse.ProfilePicture = userDB.ProfilePicture
	publicResponse.BirthDate = userDB.BirthDate
	publicResponse.Points = userDB.Points
	publicResponse.Level = userDB.Level

	return nil

}
