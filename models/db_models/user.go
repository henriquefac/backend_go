package db_models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Name           string `gorm:"size:100;index"`
	Password       string
	RegisterDate   time.Time
	Email          string `gorm:"size:150;uniqueIndex"`
	Phone          string `gorm:"size:30"`
	ProfilePicture []byte
	BirthDate      time.Time
	Points         int
	Level          int

	MissingAnimals   []MissingAnimal
	SpottedRegisters []AnimalSpottedRegister
	AnimalsReturned  []AnimalReturnedRegister
}
