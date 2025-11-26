package db_models

import (
	"gorm.io/gorm"
	"time"
)

type AnimalReturnedRegister struct {
	gorm.Model
	MissingAnimalID uint `gorm:"not null;unique"`
	RescuerID       uint `gorm:"not null"`
	ReturnDate      time.Time
	Rescuer         User `gorm:"foreignKey:RescuerID"`
}
