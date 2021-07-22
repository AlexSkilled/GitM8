package service

import (
	"gitlab-tg-bot/internal"
	"gitlab-tg-bot/internal/interfaces"
)

type Storage struct {
	interfaces.UserService
	interfaces.GitlabApiService
	interfaces.SubscriptionService
}

func NewStorage(providerStorage interfaces.ProviderStorage, conf internal.Configuration) *Storage {
	// Убираю доступ к апи как к сервису, напрямую.
	gitlabApiService := NewGitlabApiService(conf)

	return &Storage{
		UserService:         NewUserService(providerStorage),
		GitlabApiService:    gitlabApiService,
		SubscriptionService: NewSubscription(conf, providerStorage, gitlabApiService),
	}
}

func (s *Storage) User() interfaces.UserService {
	return s.UserService
}

func (s *Storage) Subscription() interfaces.SubscriptionService {
	return s.SubscriptionService
}
