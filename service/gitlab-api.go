package service

import (
	"gitlab-tg-bot/internal/interfaces"
	"gitlab-tg-bot/internal/model"

	gapi "github.com/xanzy/go-gitlab"
)

type GitlabApiService struct {
}

var _ interfaces.GitlabApiService = (*GitlabApiService)(nil)

func NewGitlabApiService() *GitlabApiService {
	return &GitlabApiService{}
}

func (g *GitlabApiService) GetRepositories(gitlabUser model.GitlabUser) ([]model.Repository, error) {
	client, err := gapi.NewClient(gitlabUser.Token)
	if err != nil {
		return nil, err
	}
	list, resp, err := client.Projects.ListProjects(&gapi.ListProjectsOptions{})
	if err != nil {
		return nil, err
	}
	list[0].Name = ""
	resp.Header.Get("")
	return nil, err
}
