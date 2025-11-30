package db_models

import "gorm.io/gorm"

type MissingAnimal struct {
	gorm.Model
	UserID        uint   `gorm:"not null"`
	Name          string `gorm:"size:100"`
	AnimalPicture []byte
	Description   string
	Status        int
	DangerLevel   int

	User User

	SpottedRegister []AnimalSpottedRegister `gorm:"foreignKey:MissingAnimalID"`

	ReturnedRegister AnimalReturnedRegister `gorm:"foreignKey:MissingAnimalID"`
}
