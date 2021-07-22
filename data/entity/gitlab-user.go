package entity

import "gitlab-tg-bot/internal/model"

type GitlabUsers []GitlabUser

func (u *GitlabUsers) ToModel() []model.GitlabUser {
	out := make([]model.GitlabUser, len(*u))
	for i, item := range *u {
		out[i] = item.ToModel()
	}
	return out
}

type GitlabUser struct {
	tableName struct{} `pg:"gitlab_user"`
	Id        int64
	Name      string
	UserId    int64
	Email     string
	Token     string
	Domain    string
}

func (u *GitlabUser) ToModel() model.GitlabUser {
	return model.GitlabUser{
		Id:       u.Id,
		Email:    u.Email,
		Username: u.Email,
		Token:    u.Token,
		Domain:   u.Domain,
	}
}

func (u *GitlabUser) FromModel(id int64, user model.GitlabUser) {
	u.UserId = id
	u.Email = user.Email
	u.Token = user.Token
	u.Domain = user.Domain
	u.Name = user.Username
}
