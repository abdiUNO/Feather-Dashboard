package models

import (
	"time"
)

type Vote struct {
	ID        string `gorm:"primary_key;type:varchar(255);"`
	Dir       int
	PostId    uint `json:"user_id"`
	UserId    string
	User      User
	CreatedAt time.Time `gorm:"column:createdAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt"`
}

func (Vote) TableName() string {
	return "vote"
}
