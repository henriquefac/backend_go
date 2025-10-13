package db_models

import (
	"gorm.io/gorm"
)

type Friendship struct {
	gorm.Model
	UserID   uint `json:"user_id"`
	FriendID uint `json:"friend_id"`

	User   User `gorm:"foreignKey:UserID" json:"-"`
	Friend User `gorm:"foreignKey:FriendID" json:"-"`
}
