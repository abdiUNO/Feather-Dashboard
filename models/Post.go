package models

import (
	"fmt"
	"time"

	u "github.com/abdullahi/go-api/utils"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

type Post struct {
	ID            string `gorm:"primary_key;type:varchar(255);"`
	Text          string `json:"text"`
	Category      string `json:"category"`
	VotesCount    int    `gorm:"column:votesCount"`
	CommentsCount int    `gorm:"column:commentsCount"`
	Time          string `sql:"DEFAULT:'0'"`
	Color         string
	Image         string    `sql:"DEFAULT:'null'"`
	CreatedAt     time.Time `gorm:"column:createdAt"`
	UpdatedAt     time.Time `gorm:"column:updatedAt"`
	User          User      `gorm:"foreignkey:UserID"`
	UserID        string    `gorm:"column:userId"`
}

func (Post) TableName() string {
	return "post"
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func (post *Post) BeforeCreate(scope *gorm.Scope) error {
	u1 := uuid.Must(uuid.NewV4())
	scope.SetColumn("ID", u1.String())
	return nil
}

func (post *Post) Validate() (map[string]interface{}, bool) {
	if post.Text == "" {
		return u.Message(false, "Post should contain text"), false
	}

	return u.Message(true, "success"), true
}

func (post *Post) Create() map[string]interface{} {
	if resp, ok := post.Validate(); !ok {
		return resp
	}

	GetDB().Create(post)

	response := u.Message(true, "Post has been created")

	response["post"] = post

	return response
}

func GetPost(id string) *Post {

	post := &Post{}
	err := GetDB().Table("post").Preload("User").Where("id = ?", id).First(post).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return post
}

func GetPosts() []*Post {

	posts := make([]*Post, 0)

	err := GetDB().Table("post").Preload("User").Order("createdAt DESC").Find(&posts).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return posts
}
