package model

import (
	"time"
)

type BaseModel struct {
	Id         uint64 `gorm:"AUTO_INCREMENT;PRIMARY_KEY"`
	AddTime    time.Time
	UpdateTime time.Time
}

type User struct {
	BaseModel
	Avarta   string
	Email    string
	Password string
	NickName string
	Flag     int
}

func (User) TableName() string {
	return "i_user"
}

type Chat struct {
	BaseModel
	Content  string
	UserId   uint64
	Source   string
	Status   int
	ChatId   uint64
	TaskCode string
}

func (Chat) TableName() string {
	return "i_chat"
}
