package service

import "gitlab-tg-bot/data"

type Storage struct {
	User UserService
}

func NewStorage(holder data.ProviderStorage) Storage {
	return Storage{
		User: NewUserService(holder),
	}
}
