package model

import "time"

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
