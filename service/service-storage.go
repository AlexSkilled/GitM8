package service

import (
	"gitlab-tg-bot/internal/interfaces"
)

type Storage struct {
	interfaces.UserService
	interfaces.SubscriptionService
	interfaces.MessageService
	interfaces.SettingsService
	interfaces.GitApiService
}

var _ interfaces.ServiceStorage = (*Storage)(nil)

func NewStorage(providerStorage interfaces.ProviderStorage, conf interfaces.Configuration) *Storage {
	// Убираю доступ к апи как к сервису, напрямую.
	gitApiService := NewGitManger(conf, providerStorage)

	return &Storage{
		UserService:         NewUserService(providerStorage),
		SubscriptionService: NewSubscription(providerStorage, gitApiService),
		MessageService:      NewMessageService(gitApiService, providerStorage),
		SettingsService:     NewSettingsService(providerStorage),
		GitApiService:       gitApiService,
	}
}

func (s *Storage) User() interfaces.UserService {
	return s.UserService
}

func (s *Storage) Subscription() interfaces.SubscriptionService {
	return s.SubscriptionService
}

func (s *Storage) MessageHandler() interfaces.MessageService {
	return s.MessageService
}

func (s *Storage) Settings() interfaces.SettingsService {
	return s.SettingsService
}

func (s *Storage) GitApi() interfaces.GitApiService {
	return s.GitApiService
}
