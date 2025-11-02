package data_models

import "time"

// Aqui serão as structs que devem representar os
// dados em trasinção em tempo de execução do código

// Estrutura que recebe diretamente dados da requisição
// com biding
// Tem o propósito de criar um novo usuário
type CreateUserRequest struct {
	Name      string    `json:"name" binding:"required"`
	Email     string    `json:"email" binding:"required"`
	Phone     string    `json:"phone" binding:"required"`
	BirthDate time.Time `json:"birthDate" binding:"required"`
	Password  string    `json:"password" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Estrutura que deve representar a resposta para
// A criação do usuário

type CreateUserResponse struct {
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	BirthDate time.Time `json:"birthDate"`
}

// Estrutura para representar dados em transito
// Serve para representar tudo sobre o usuário que seria público

// por exemplo, atende a uma requisição "get_user"

type PublicUserResponse struct {
	ID             uint      `json:"id"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	RegisterDate   time.Time `json:"registerData"`
	Phone          string    `json:"phone"`
	ProfilePicture []byte    `json:"profilePicture,omitempty"`
	BirthDate      time.Time `json:"birthDate"`
	Points         int       `json:"points"`
	Level          int       `json:"level"`
	Password       string    `json:"-"`
}
