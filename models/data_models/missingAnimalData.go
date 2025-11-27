package data_models

import "time"

// criar struct para receber requisicao de criacao de missing animal

type LastSeen struct {
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Description string  `json:"description"`
}

type MissingAnimalRequest struct {
	UserID        uint     `json:"userid" binding:"required"`
	Name          string   `json:"name" binding:"required"`
	AnimalPicture *[]byte  `json:"animalPicture"`
	Description   string   `json:"description" binding:"required"`
	DangerLevel   int      `json:"dangerLevel" binding:"required"`
	LastSeen      LastSeen `json:"initialLoc" binding:"required"`
}

// resposta para o request

type LastSeenResponse struct {
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Description string  `json:"description"`
}

// Struct para a Resposta Completa do Novo Registro
type MissingAnimalResponse struct {
	ID          uint   `json:"id"`
	UserID      uint   `json:"userId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      int    `json:"status"`
	DangerLevel int    `json:"dangerLevel"`
	// Aqui você pode adicionar a Data de Criação (do gorm.Model)
	CreatedAt time.Time `json:"createdAt"`

	// Localização inicial (Última vista)
	LastSeen LastSeenResponse `json:"lastSeen"`
}

// modelo de request para alterar um registro de missing animal
type MissingAnimalUpdateRequest struct {
	ID            uint    `json:"id" binding:"required"`
	UserID        uint    `json:"userId" binding:"required"`
	Name          *string `json:"name"`
	AnimalPicture *[]byte `json:"animalPicture"`
	Description   *string `json:"description"`
	Status        *int    `json:"status"`
	DangerLevel   *int    `json:"dangerLevel"`
}
