package interfaces

import "gitlab-tg-bot/internal/model"

type ServiceStorage interface {
	User() UserService
}

type UserService interface {
	GetByTelegramId(tgId int64) (model.User, error)
	Register(user model.User) error
	AddGitlabAccount(tgId int64, gitlab model.GitlabUser) error
}

type TelegramWorker interface {
	Start()
	SendMessage(chatId []int32, msg string)
}
