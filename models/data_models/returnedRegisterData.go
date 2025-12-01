package data_models

import (
	"time"
)

type ReturnedRegisterCreateRequest struct {
	MissingAnimalID uint      `json:"missingAnimalId" binding:"required"`
	RescuerID       uint      `json:"rescuerID" binding:"required"`
	ReturnDate      time.Time `json:"returnDate" binding:""`
}

type ReturnedRegisterResponse struct {
	// Campos do GORM.Model:
	ID              uint      `json:"id"`
	CreatedAt       time.Time `json:"createdAt"` // Quando o registro foi salvo
	MissingAnimalID uint      `json:"missingAnimalId" binding:"required"`
	RescuerID       uint      `json:"rescuerID" binding:"required"`
	ReturnDate      time.Time `json:"returnDate" binding:""`
}
