package db_models

import (
	"gorm.io/gorm"
	"time"
)

type AnimalSpottedRegister struct {
	gorm.Model
	MissingAnimalID uint `gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	UserID          uint `gorm:"not null"`

	AnimalPicture []byte

	Latitude  float64 `gorm:"not null"`
	Longitude float64 `gorm:"not null"`

	SpottedTime time.Time
	Description string

	MissingAnimal MissingAnimal
	User          User
}
