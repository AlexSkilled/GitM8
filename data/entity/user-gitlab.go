package entity

import "gitlab-tg-bot/internal/model"

type UserGitlab struct {
	Id int32

	UserId         int32
	UserName       string
	UserTgUsername string
	UserTgId       int64

	Email   string
	Token   string
	BaseUrl string
}

func (u *UserGitlab) ToModel() model.UserGitlab {
	return model.UserGitlab{
		Id: u.Id,
		User: model.User{
			Id:         u.UserId,
			Name:       u.UserName,
			TgUsername: u.UserTgUsername,
			TgId:       u.UserTgId,
		},
		Email:  u.Email,
		Token:  u.Token,
		Domain: u.BaseUrl,
	}
}
