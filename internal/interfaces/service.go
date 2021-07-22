package interfaces

import "gitlab-tg-bot/internal/model"

type ServiceStorage interface {
	User() UserService
	Subscription() SubscriptionService
}

type UserService interface {
	GetWithGitlabUsersById(id int64) (model.User, error)
	Register(user model.User) (model.User, error)
	AddGitlabAccount(tgId int64, gitlab model.GitlabUser) error
}

type TelegramWorker interface {
	Start()
	SendMessage(chatId []int32, msg string)
}

type GitlabApiService interface {
	GetRepositories(user model.GitlabUser) ([]model.Repository, error)
	AddWebHook(user model.GitlabUser, hookInfo model.Hook) error
}

type SubscriptionService interface {
	Subscribe(user model.GitlabUser, tgUserId int64, hookInfo model.Hook) (ticketId int32, err error)
	GetRepositories(user model.GitlabUser) ([]model.Repository, error)
}
