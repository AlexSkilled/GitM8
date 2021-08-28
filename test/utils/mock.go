package utils

import (
	"gitlab-tg-bot/internal/interfaces"
	"gitlab-tg-bot/service/model"
)

type TgMessengerMock struct {
	messages []model.OutputMessage
}

var _ interfaces.TelegramMessageSender = (*TgMessengerMock)(nil)

func (t *TgMessengerMock) SendMessages(message []model.OutputMessage) {
	message = append(message, message...)
}
