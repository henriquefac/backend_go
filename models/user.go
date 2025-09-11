package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Name           string    `json:"name" gorm:"size:100;index"`
	Email          string    `json:"email" gorm:"size:150;uniqueIndex"`
	Phone          string    `json:"phone" gorm:"size:30"`
	ProfilePicture []byte    `json:"profile_picture"`
	BirthDate      time.Time `json:"birth_date"`
	Points         int       `json:"points"`
	Level          int       `json:"level"`
}
