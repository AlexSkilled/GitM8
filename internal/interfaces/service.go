package interfaces

import (
	"gitlab-tg-bot/service/model"
)

type ServiceStorage interface {
	User() UserService
	Subscription() SubscriptionService
	Webhook() WebhookService
}

type UserService interface {
	GetWithGitlabUsersById(id int64) (model.User, error)
	Register(user model.User) (model.User, error)
	AddGitlabAccount(tgId int64, gitlab model.GitlabUser) error
}

type TelegramWorker interface {
	Start()
	TelegramMessageSender
}

type TelegramMessageSender interface {
	SendMessage(chatId []int32, msg string)
}

type GitApiService interface {
	GetRepositories(user model.GitlabUser) ([]model.Repository, error)
	AddWebHook(user model.GitlabUser, hookInfo model.Hook) error
}

type SubscriptionService interface {
	Subscribe(user model.GitlabUser, chatId int64, hookInfo model.Hook) (ticketId int32, err error)
	GetRepositories(user model.GitlabUser) ([]model.Repository, error)
}

type Configuration interface {
	GetBool(string) bool
	GetInt(string) int
	GetInt32(string) int32
	GetInt64(string) int64
	GetString(string) string
}

type WebhookService interface {
	ProcessMessage(event model.GitEvent) (msg string, chatIds []int32, err error)
}
