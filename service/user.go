package service

import (
	"gitlab-tg-bot/internal/interfaces"
	"gitlab-tg-bot/internal/model"
)

type UserService struct {
	UserProvider interfaces.UserProvider
}

var _ interfaces.UserService = (*UserService)(nil)

func NewUserService(provider interfaces.ProviderStorage) interfaces.UserService {
	return &UserService{
		UserProvider: provider.User(),
	}
}

func (u *UserService) GetByTelegramId(tgId int64) (model.User, error) {
	return u.UserProvider.Get(tgId)
}

func (u *UserService) Register(user model.User) error {
	err := u.UserProvider.Create(user)
	return err
}
