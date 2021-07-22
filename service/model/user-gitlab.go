package model

import "github.com/plouc/go-gitlab-client/gitlab"

type GitlabUser struct {
	Id int64
	//TgUserId int64
	Username string
	Token    string
	Domain   string
	Email    string
}

func (u *GitlabUser) GetGitlabClient() *gitlab.Gitlab {
	return gitlab.NewGitlab(u.Domain, "api/v4", u.Token)
	//return gitlab.NewGitlab("https://gitlab.ru/", ApiPath,"DFX3ppBJb7qdBsjz3DsH")
}
