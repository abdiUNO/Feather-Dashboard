package models

import "time"

type MicroLink struct {
	ID          string `gorm:"primary_key;type:varchar(255);"`
	Title       string
	Image       string
	Post        Post
	Description string
	Logo        string
	Url         string
	CreatedAt   time.Time `gorm:"column:createdAt"`
	UpdatedAt   time.Time `gorm:"column:updatedAt"`
}

func (MicroLink) TableName() string {
	return "micro_link"
}
