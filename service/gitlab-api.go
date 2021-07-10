package service

import (
	"gitlab-tg-bot/internal/interfaces"
	"gitlab-tg-bot/internal/model"

	"github.com/sirupsen/logrus"

	gapi "github.com/plouc/go-gitlab-client/gitlab"
)

const standardApiLevel = "api/v4"

type GitlabApiService struct {
}

var _ interfaces.GitlabApiService = (*GitlabApiService)(nil)

func NewGitlabApiService() *GitlabApiService {
	return &GitlabApiService{}
}

func (g *GitlabApiService) GetRepositories(gitlabUser model.GitlabUser) ([]model.Repository, error) {
	//client := gapi.NewGitlab("https://gitlab.ru/", "api/v4", "tBm_wxrKuwQxEPAzRQN4")
	client := gapi.NewGitlab(gitlabUser.Domain, standardApiLevel, gitlabUser.Token)

	list, resp, err := client.Projects(&gapi.ProjectsOptions{Membership: true})
	if err != nil {
		return nil, err
	}

	logrus.Info(list)
	logrus.Info(resp)
	return nil, err
}
