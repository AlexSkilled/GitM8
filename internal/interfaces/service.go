package interfaces

import (
	"gitlab-tg-bot/service/model"
)

type ServiceStorage interface {
	User() UserService
	Subscription() SubscriptionService
	MessageHandler() MessageService
}

type UserService interface {
	GetWithGitlabUsersById(id int64) (model.User, error)
	Register(user model.User) (model.User, error)
	AddGitAccount(tgId int64, gitlab model.GitUser) error
}

type TelegramWorker interface {
	Start()
	TelegramMessageSender
}

type TelegramMessageSender interface {
	SendMessages(message []model.OutputMessage)
}

type GitApiService interface {
	GetRepositories(user model.GitUser) ([]model.Repository, error)
	AddWebHook(user model.GitUser, hookInfo model.Hook) error
	GetUser(git model.GitUser, userId string) (model.GitUserInfo, error)
}

type SubscriptionService interface {
	Subscribe(user model.GitUser, chatId int64, hookInfo model.Hook) (ticketId int32, err error)
	GetRepositories(user model.GitUser) ([]model.Repository, error)
}

type Configuration interface {
	GetBool(string) bool
	GetInt(string) int
	GetInt32(string) int32
	GetInt64(string) int64
	GetString(string) string
}

type MessageService interface {
	ProcessMessage(event model.GitEvent) (msg []model.OutputMessage, err error)
}
