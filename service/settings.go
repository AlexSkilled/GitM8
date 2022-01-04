package service

import (
	"gitlab-tg-bot/internal/interfaces"
	"gitlab-tg-bot/service/model"
)

type SettingsService struct {
	user interfaces.UserProvider
}

var _ interfaces.SettingsService = (*SettingsService)(nil)

func NewSettingsService(data interfaces.ProviderStorage) *SettingsService {
	return &SettingsService{
		user: data.User(),
	}
}

func (s *SettingsService) ChangeLanguage(userId int64, language string) error {
	err := s.user.Update(model.User{
		Id:     userId,
		Locale: language,
	})
	return err
}
