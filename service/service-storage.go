package service

import (
	"gitlab-tg-bot/data"
	"gitlab-tg-bot/internal/interfaces"
)

type Storage struct {
	interfaces.UserService
}

func NewStorage(providerStorage data.ProviderStorage) Storage {
	return Storage{
		UserService: NewUserService(&providerStorage),
	}
}

func (s *Storage) User() interfaces.UserService {
	return s.UserService
}
