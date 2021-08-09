package service

import (
	"gitlab-tg-bot/internal/interfaces"
)

type Storage struct {
	interfaces.UserService
	interfaces.SubscriptionService
	interfaces.MessageService
}

var _ interfaces.ServiceStorage = (*Storage)(nil)

func NewStorage(providerStorage interfaces.ProviderStorage, conf interfaces.Configuration) *Storage {
	// Убираю доступ к апи как к сервису, напрямую.
	gitlabApiService := NewGitlabApiService(conf)

	return &Storage{
		UserService:         NewUserService(providerStorage),
		SubscriptionService: NewSubscription(conf, providerStorage, gitlabApiService),
		MessageService:      NewMessageService(gitlabApiService, providerStorage),
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
