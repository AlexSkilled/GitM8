package service

import (
	config "gitlab-tg-bot/conf"
	"gitlab-tg-bot/internal"
	"gitlab-tg-bot/internal/interfaces"
	"gitlab-tg-bot/service/model"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"

	gapi "github.com/plouc/go-gitlab-client/gitlab"
)

const StandardApiLevel = "api/v4"

type GitlabApiService struct {
	webHookUrl string
}

var _ interfaces.GitApiService = (*GitlabApiService)(nil)

func NewGitlabApiService(conf interfaces.Configuration) *GitlabApiService {
	return &GitlabApiService{
		webHookUrl: conf.GetString(config.WebHookUrl) + model.Gitlab.GetUri(),
	}
}

func (g *GitlabApiService) GetRepositories(gitlabUser model.GitUser) ([]model.Repository, error) {
	client := gapi.NewGitlab(gitlabUser.Domain, StandardApiLevel, gitlabUser.Token)

	list, _, err := client.Projects(&gapi.ProjectsOptions{Membership: true})
	if err != nil {
		return nil, err
	}

	return g.toModelProjects(list), err
}

func (g *GitlabApiService) AddWebHook(gitlabUser model.GitUser, hookInfo model.Hook) error {
	client := gapi.NewGitlab(gitlabUser.Domain, StandardApiLevel, gitlabUser.Token)

	addr, _ := model.GetMyInterfaceAddr()
	logrus.Info(addr)
	hookPayload := gapi.HookAddPayload{
		Url:                      g.webHookUrl,
		PushEvents:               hookInfo.PushEvents,
		IssuesEvents:             hookInfo.IssuesEvents,
		ConfidentialIssuesEvents: hookInfo.ConfidentialIssuesEvents,
		MergeRequestsEvents:      hookInfo.MergeRequestsEvents,
		TagPushEvents:            hookInfo.TagPushEvents,
		NoteEvents:               hookInfo.NoteEvents,
		JobEvents:                hookInfo.JobEvents,
		PipelineEvents:           hookInfo.PipelineEvents,
		WikiPageEvents:           hookInfo.WikiPageEvents,
		Token:                    gitlabUser.Token,
	}

	_, _, err := client.AddProjectHook(hookInfo.RepoId, &hookPayload)
	if err != nil {
		return err
	}
	return nil
}

func (g *GitlabApiService) GetUser(git model.GitUser, userId string) (model.GitUserInfo, error) {
	client := gapi.NewGitlab(git.Domain, StandardApiLevel, git.Token)
	user, _, err := client.User(userId)
	if err != nil {
		if strings.Contains(err.Error(), "<401>") {
			return model.GitUserInfo{}, internal.ErrorAccessDenied
		}
		return model.GitUserInfo{}, err
	}

	return model.GitUserInfo{
		Id:   int32(user.Id),
		Name: user.Name,
		Url:  user.WebUrl,
	}, nil
}

func (g *GitlabApiService) toModelProjects(in *gapi.ProjectCollection) []model.Repository {
	out := make([]model.Repository, len(in.Items))
	for i, item := range in.Items {
		out[i] = model.Repository{
			Id:   strconv.Itoa(item.Id),
			Name: item.Name,
		}
	}
	return out
}
