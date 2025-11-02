package services

import (
	"errors"
	"github.com/henriquefac/backend_go/models/data_models"
	"github.com/henriquefac/backend_go/repositories"
	"github.com/henriquefac/backend_go/utils"
)

var ErrInvalidPassword = errors.New("Senha inválida")

type UserService struct {
	userRepo *repositories.UserRepository
}

func NewUserService(ur *repositories.UserRepository) *UserService {
	return &UserService{userRepo: ur}
}

func (us *UserService) Login(email,
	password string) (*data_models.PublicUserResponse, error) {
	var publicResponse data_models.PublicUserResponse

	err := us.userRepo.GetUserByEmail(email, &publicResponse)
	if err != nil {
		return nil, err // pode ser ErrUserNotFound
	}

	if !utils.CheckPassword(&publicResponse.Password, &password) {
		return nil, ErrInvalidPassword
	}

	// Remove a senha antes de retornar
	publicResponse.Password = ""
	return &publicResponse, nil
}
