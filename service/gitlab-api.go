package service

import (
	"gitlab-tg-bot/internal/interfaces"
	"gitlab-tg-bot/internal/model"

	"github.com/sirupsen/logrus"

	gapi "github.com/plouc/go-gitlab-client/gitlab"
)

const StandardApiLevel = "api/v4"

type GitlabApiService struct {
}

var _ interfaces.GitlabApiService = (*GitlabApiService)(nil)

func NewGitlabApiService() *GitlabApiService {
	return &GitlabApiService{}
}

func (g *GitlabApiService) GetRepositories(gitlabUser model.GitlabUser) ([]model.Repository, error) {
	client := gapi.NewGitlab(gitlabUser.Domain, StandardApiLevel, gitlabUser.Token)

	list, resp, err := client.Projects(&gapi.ProjectsOptions{Membership: true})
	if err != nil {
		return nil, err
	}

	logrus.Info(list)
	logrus.Info(resp)
	return g.toModelProjects(list), err
}

func (g *GitlabApiService) toModelProjects(in *gapi.ProjectCollection) []model.Repository {
	out := make([]model.Repository, len(in.Items))
	for i, item := range in.Items {
		out[i] = model.Repository{
			Id:   int32(item.Id),
			Name: item.Name,
		}
	}
	return out
}
