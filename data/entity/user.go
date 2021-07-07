package entity

import (
	"gitlab-tg-bot/internal/model"
)

type User struct {
	tableName  struct{} `pg:"user"`
	Id         int32
	Name       string
	TgUsername string
	TgId       int64
}

func (u *User) FromModel(user model.User) {
	u.Id = user.Id
	u.Name = user.Name
	u.TgUsername = user.TgUsername
	u.TgId = user.TgId
}

func (u *User) ToModel() model.User {
	return model.User{
		Id:         u.Id,
		Name:       u.Name,
		TgUsername: u.TgUsername,
		TgId:       u.TgId,
	}
}
