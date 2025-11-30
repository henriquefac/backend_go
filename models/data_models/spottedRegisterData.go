package data_models

import (
	"time"
)

type SpottedRegisterCreateRequest struct {
	MissingAnimalID uint    `json:"missingAnimalId" binding:"required"`
	UserID          uint    `json:"userId" binding:"required"`
	AnimalPicture   *[]byte `json:"animalPicture"`

	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`

	SpottedTime time.Time `json:"spottedTime"`
	Description string    `json:"description"`
}

type SpottedRegisterResponse struct {
	ID              uint `json:"id"`
	UserID          uint `json:"userId"`
	MissingAnimalID uint `json:"missingAnimalId"`

	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`

	SpottedTime time.Time `json:"spottedTime"`
	Description string    `json:"description"`
}
