package interfaces

import "gitlab-tg-bot/internal/model"

type ServiceStorage interface {
	User() UserService
}

type UserService interface {
	Register(user model.User) error
}

type TelegramWorker interface {
	Start()
	SendMessage(chatId []int32, msg string)
}
