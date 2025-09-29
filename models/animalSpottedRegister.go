package models

import (
	"gorm.io/gorm"
	"time"
)

type AnimalSpottedRegister struct {
	gorm.Model
	MissingAnimalID uint `gorm:"not null"`
	AnimalPicture   []byte
	LocalFound      string
	Time            time.Time
}
