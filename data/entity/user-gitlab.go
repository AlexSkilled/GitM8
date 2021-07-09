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
	UserId    int64
	Email     string
	Token     string
	Domain    string
}

func (u *GitlabUser) ToModel() model.GitlabUser {
	return model.GitlabUser{
		UserId:   u.UserId,
		Username: u.Email,
		Token:    u.Token,
		Domain:   u.Domain,
	}
}

func (u *GitlabUser) FromModel(id int64, user model.GitlabUser) {
	u.UserId = id
	u.Email = user.Username
	u.Token = user.Token
	u.Domain = user.Domain
}
