package service

import (
	"gitlab-tg-bot/internal/interfaces"
	"gitlab-tg-bot/service/model"
)

type WebhookService struct {
	ticket interfaces.TicketProvider
}

var _ interfaces.WebhookService = (*WebhookService)(nil)

func NewWebhook() interfaces.WebhookService {
	return &WebhookService{}
}

func (w *WebhookService) ProcessMessage(event model.GitEvent) (messages []model.OutputMessage, err error) {

	switch event.HookType {
	//case model.GitEventPush:
	//case model.GitEventPushTag:
	//case model.GitEventPushIssue:
	//case model.GitEventPushNote:
	case model.GitEventMergeRequest:

		break
		//case model.GitEventWiki:
		//case model.GitEventPipeline:
		//case model.GitEventJob:
		//case model.GitEventDeployment:
		//case model.GitEventMember:
		//case model.GitEventSubgroup:
		//case model.GitEventFeatureFlag:
		//case model.GitEventRelease:
	}
	messages = make([]model.OutputMessage, 0)
	return messages, nil
}
