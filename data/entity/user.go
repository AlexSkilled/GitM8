package entity

import (
	"gitlab-tg-bot/internal/model"
)

type User struct {
	tableName  struct{} `pg:"tg_user"`
	Id         int64
	Name       string
	TgUsername string
}

func (u *User) FromModel(user model.User) {
	u.Id = user.Id
	u.Name = user.Name
	u.TgUsername = user.TgUsername
}

func (u *User) ToModel() model.User {
	return model.User{
		Id:         u.Id,
		Name:       u.Name,
		TgUsername: u.TgUsername,
	}
}
