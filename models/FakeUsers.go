package models

import (
	"fmt"
	"time"

	u "github.com/abdullahi/go-api/utils"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

//a struct to rep user
type FakeUser struct {
	ID        string    `gorm:"primary_key;type:varchar(255);"`
	Username  string    `json:"username";gorm:"column:username"`
	User      User      `gorm:"foreignkey:UserId"`
	UserId    string    `gorm:"column:userId"`
	CreatedAt time.Time `gorm:"column:createdAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt"`
}

func (FakeUser) TableName() string {
	return "fake_user"
}

func (user *FakeUser) BeforeCreate(scope *gorm.Scope) error {
	u1 := uuid.Must(uuid.NewV4())
	scope.SetColumn("ID", u1.String())
	return nil
}

func (fakeUser *FakeUser) Create() map[string]interface{} {
	GetDB().Create(fakeUser)

	response := u.Message(true, "User has been created")
	return response
}

func GetFakerUsers() []*FakeUser {
	fakeUsers := make([]*FakeUser, 0)

	err := GetDB().Table("fake_user").Preload("User").Find(&fakeUsers).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return fakeUsers
}
