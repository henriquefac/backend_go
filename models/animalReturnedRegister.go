package models

import (
	"gorm.io/gorm"
	"time"
)

type AnimalReturnedRegister struct {
	gorm.Model
	MissingAnimalID uint `gorm:"not null"`
	RescuerID       uint `gorm:"not null"`
	ReturnDate      time.Time
}
