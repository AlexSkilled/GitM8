package model

import "github.com/plouc/go-gitlab-client/gitlab"

type GitUser struct {
	Id int64
	//TgUserId int64
	Username string
	Token    string
	Domain   string
	Email    string
	GitSource
}

func (u *GitUser) GetGitlabClient() *gitlab.Gitlab {
	return gitlab.NewGitlab(u.Domain, "api/v4", u.Token)
}
