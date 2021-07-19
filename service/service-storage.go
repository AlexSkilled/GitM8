package service

import (
	"gitlab-tg-bot/data"
	"gitlab-tg-bot/internal/interfaces"
)

type Storage struct {
	interfaces.UserService
	interfaces.GitlabApiService
}

func NewStorage(providerStorage data.ProviderStorage) Storage {
	return Storage{
		UserService:      NewUserService(&providerStorage),
		GitlabApiService: NewGitlabApiService(),
	}
}

func (s *Storage) User() interfaces.UserService {
	return s.UserService
}

func (s *Storage) GitlabApi() interfaces.GitlabApiService {
	return s.GitlabApiService
}
