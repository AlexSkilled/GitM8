package interfaces

import (
	"gitlab-tg-bot/service/model"
)

type ServiceStorage interface {
	User() UserService
	Subscription() SubscriptionService
	MessageHandler() MessageService
	Settings() SettingsService
	GitApi() GitApiService
}

type UserService interface {
	Register(user model.User) (model.User, error)
	AddGitAccount(tgId int64, gitlab model.GitUser) error
	GetWithGitlabUsersById(id int64) (model.User, error)
}

type SettingsService interface {
	ChangeLanguage(userId int64, language string) error
}

type TelegramWorker interface {
	Start()
	TelegramMessageSender
}

type TelegramMessageSender interface {
	SendMessages(message []model.OutputMessage)
}

type GitApiService interface {
	GetGitType(string) model.GitSource
	GetRepositories(user model.GitUser) ([]model.Repository, error)
	AddWebHook(user model.GitUser, hookInfo model.Hook) error
	GetUser(git model.GitUser, userId string) (model.GitUserInfo, error)
	GetWebhookUrl(domain string, tgUserId int64) (string, error)
}

type SubscriptionService interface {
	Subscribe(user model.GitUser, chatId int64, hookInfo model.Hook) (ticketId int32, err error)
	GetRepositories(user model.GitUser) ([]model.Repository, error)
	GetUserTickets(userId int64) (tickets []model.Ticket, err error)
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

type TokenExpiration interface {
	Start()
	AddGitIds(gits []model.GitUser)
	Stop()
}
