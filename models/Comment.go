package models

import (
	"fmt"
	"time"

	u "github.com/abdullahi/go-api/utils"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

type Comment struct {
	ID         string `gorm:"primary_key;type:varchar(255);"`
	Text       string `json:"text"`
	Time       string
	VotesCount int       `gorm:"column:votesCount"`
	CreatedAt  time.Time `gorm:"column:createdAt"`
	UpdatedAt  time.Time `gorm:"column:updatedAt"`
	User       User      `gorm:"foreignkey:UserID"`
	UserID     string    `gorm:"column:userId"`
	PostID     string    `gorm:"column:postId"`
}

func (Comment) TableName() string {
	return "comment"
}

func (comment *Comment) BeforeCreate(scope *gorm.Scope) error {
	u1 := uuid.Must(uuid.NewV4())
	scope.SetColumn("ID", u1.String())
	return nil
}

func GetComments(id string) []*Comment {
	comments := make([]*Comment, 0)

	err := GetDB().Table("comment").Preload("User").Where("postId = ?", id).Order("createdAt DESC").Find(&comments).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return comments
}

func (comment *Comment) Create() map[string]interface{} {
	post := GetPost(comment.PostID)

	GetDB().Model(&post).Update("commentsCount", post.CommentsCount+1)
	GetDB().Create(comment)

	response := u.Message(true, "Post has been created")

	response["comment"] = comment

	return response
}
