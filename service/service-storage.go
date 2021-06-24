package service

import "gitlab-tg-bot/data"

type Holder struct {
	User  UserService
	Pipes PipeService
}

func NewStorage(holder data.ProviderStorage) Holder {
	return Holder{
		User: NewUserService(holder),
	}
}
