package service

import (
	"gitlab-tg-bot/internal/interfaces"
	"gitlab-tg-bot/service/model"
	"gitlab-tg-bot/service/processor"
)

type MessageService struct {
	ticket   interfaces.TicketProvider
	patterns interfaces.MessagePatternProvider

	interfaces.GitApiService
}

var _ interfaces.MessageService = (*MessageService)(nil)

func NewMessageService(gitlabApi interfaces.GitApiService, data interfaces.ProviderStorage) *MessageService {
	return &MessageService{
		ticket:        data.Ticket(),
		patterns:      data.MessagePattern(),
		GitApiService: gitlabApi,
	}
}

func (s *MessageService) ProcessMessage(event model.GitEvent) ([]model.OutputMessage, error) {
	tickets, err := s.ticket.GetTicketsToSend(event.ProjectId, event.HookType)
	if err != nil {
		return nil, err
	}

	messageText, patterns, err := s.patterns.GetMessage("ru_RU", event.HookType, event.SubType)
	if err != nil {
		return nil, err
	}

	if len(event.AuthorId) != 0 {
		gitUser, err := s.ticket.GetGitUserByTicketId(tickets[0].TicketId)
		if err != nil {
			return nil, err
		}

		user, err := s.GitApiService.GetUser(gitUser, event.AuthorId)
		if err != nil {
			return nil, err
		}
		event.AuthorName = user.Name
	}

	switch event.HookType {
	case model.HookTypeMergeRequests:
		messageText, err = processor.ProcessMergeRequest(event, messageText, patterns)
	}

	messages := make([]model.OutputMessage, len(tickets))

	for i, item := range tickets {
		messages[i].Msg = messageText
		messages[i].ChatId = item.ChatId
	}

	return messages, nil
}
