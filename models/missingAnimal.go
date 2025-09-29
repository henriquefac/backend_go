package models

import "gorm.io/gorm"

type MissingAnimal struct {
	gorm.Model
	UserID        uint   `gorm:"not null"`
	Name          string `gorm:"size:100"`
	AnimalPicture []byte
	Description   string
	LastSeen      string
	Status        int
	DangerLevel   int

	SpottedRegister  []AnimalSpottedRegister
	ReturnedRegister []AnimalReturnedRegister
}
