package service

import (
	"errors"
	"fmt"
	"strings"

	"gitlab-tg-bot/internal/interfaces"
	"gitlab-tg-bot/service/model"
)

type GitManager struct {
	sourceToService map[model.GitSource]interfaces.GitApiService
}

func NewGitManger(cont interfaces.Configuration) *GitManager {
	return &GitManager{
		sourceToService: map[model.GitSource]interfaces.GitApiService{
			model.Gitlab: NewGitlabApiService(cont),
		},
	}
}

var _ interfaces.GitApiService = (*GitManager)(nil)

func (g *GitManager) GetRepositories(user model.GitUser) ([]model.Repository, error) {
	if api, ok := g.sourceToService[user.GitSource]; ok {
		return api.GetRepositories(user)
	} else {
		return nil, errors.New(fmt.Sprintf("Для git %s не нашлось обработчика", user.GitSource))
	}
}

func (g *GitManager) AddWebHook(user model.GitUser, hookInfo model.Hook) error {
	if api, ok := g.sourceToService[user.GitSource]; ok {
		return api.AddWebHook(user, hookInfo)
	} else {
		return errors.New(fmt.Sprintf("Для git %s не нашлось обработчика", user.GitSource))
	}
}

func (g *GitManager) GetUser(git model.GitUser, userId string) (model.GitUserInfo, error) {
	if api, ok := g.sourceToService[git.GitSource]; ok {
		return api.GetUser(git, userId)
	} else {
		return model.GitUserInfo{}, errors.New(fmt.Sprintf("Для git %s не нашлось обработчика", git.GitSource))
	}
}

func (g *GitManager) GetGitType(s string) model.GitSource {
	switch {
	case strings.Contains(s, string(model.Gitlab)):
		return model.Gitlab
	}
	return model.NotImplemented
}
