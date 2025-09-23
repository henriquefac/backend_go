package models

import (
	"gorm.io/gorm"
	"time"
)

type Achievement struct {
	gorm.Model
	UserID       uint   `json:"id_user"`
	Name         string `gorm:"size:150" json:"name"`
	Description  string `gorm:"size:600" json:"description"`
	ObtainedDate time.Time
	Icon         []byte
	Points       int
	Rarity       int
}
