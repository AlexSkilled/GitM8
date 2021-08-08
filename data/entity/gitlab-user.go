package entity

import (
	"gitlab-tg-bot/service/model"
)

type GitlabUsers []GitUser

func (u *GitlabUsers) ToModel() []model.GitUser {
	out := make([]model.GitUser, len(*u))
	for i, item := range *u {
		out[i] = item.ToModel()
	}
	return out
}

type GitUser struct {
	tableName struct{} `pg:"git_user"`
	Id        int64
	Name      string
	UserId    int64
	Email     string
	Token     string
	Domain    string
	GitSource model.GitSource
}

func (u *GitUser) ToModel() model.GitUser {
	return model.GitUser{
		Id:       u.Id,
		Email:    u.Email,
		Username: u.Email,
		Token:    u.Token,
		Domain:   u.Domain,
	}
}

func (u *GitUser) FromModel(id int64, user model.GitUser) {
	u.UserId = id
	u.Email = user.Email
	u.Token = user.Token
	u.Domain = user.Domain
	u.Name = user.Username
}
