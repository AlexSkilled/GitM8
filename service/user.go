package service

import (
	"gitlab-tg-bot/internal/interfaces"
	"gitlab-tg-bot/internal/model"
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
	return u.UserProvider.AddGitlab(tgId, gitlab)
}
