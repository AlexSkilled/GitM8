package service

import (
	"gitlab-tg-bot/internal/interfaces"
	"gitlab-tg-bot/internal/model"
	"strings"

	gapi "github.com/plouc/go-gitlab-client/gitlab"
)

type UserService struct {
	UserProvider interfaces.UserProvider
}

var _ interfaces.UserService = (*UserService)(nil)

func NewUserService(provider interfaces.ProviderStorage) *UserService {
	return &UserService{
		UserProvider: provider.User(),
	}
}

func (u *UserService) GetWithGitlabUsersByTgId(tgId int64) (model.User, error) {
	return u.UserProvider.GetWithGitlabUsers(tgId)
}

func (u *UserService) Register(user model.User) (model.User, error) {
	return u.UserProvider.Create(user)
}

func (u *UserService) AddGitlabAccount(tgId int64, gitlab model.GitlabUser) error {
	if !strings.HasPrefix(gitlab.Domain, "http") {
		gitlab.Domain = "https://" + gitlab.Domain
	}

	if !strings.HasSuffix(gitlab.Domain, "/") {
		gitlab.Domain += "/"
	}

	client := gapi.NewGitlab(gitlab.Domain, StandardApiLevel, gitlab.Token)
	user, _, err := client.CurrentUser()
	if err != nil {
		return err
	}
	gitlab.Email = user.Email

	return u.UserProvider.AddGitlab(tgId, gitlab)
}
