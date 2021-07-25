package service

import (
	"gitlab-tg-bot/internal/interfaces"
	"gitlab-tg-bot/service/model"
)

type WebhookService struct {
}

var _ interfaces.WebhookService = (*WebhookService)(nil)

func NewWebhook() interfaces.WebhookService {
	return &WebhookService{}
}

func (w WebhookService) ProcessMessage(event model.GitEvent) (msg string, chatIds []int32, err error) {

	return "", nil, err
}
