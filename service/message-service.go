package service

import (
	"gitlab-tg-bot/internal/interfaces"
	"gitlab-tg-bot/service/model"
	"gitlab-tg-bot/service/processor"

	"github.com/sirupsen/logrus"
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

	event.AuthorName, err = s.findAuthor(event.AuthorId, tickets)
	if err != nil {
		return nil, err
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

func (s *MessageService) findAuthor(authorId string, tickets []model.TicketChatId) (authorName string, err error) {
	if len(authorId) == 0 {
		return "", nil
	}
	failedToAccess := make([]model.GitUser, 0)

	gitUsers := make([]model.GitUser, 0)

	for _, item := range tickets {
		gits, err := s.ticket.GetGitUsersByTicketId(item.TicketId)
		if err != nil {
			logrus.Errorf("Ошибка при попытке получить пользователей с помощью тикета(#%d):%v", item.TicketId, err)
			continue
		}
		gitUsers = append(gitUsers, gits...)
	}

	for _, gitUser := range gitUsers {
		user, err := s.GitApiService.GetUser(gitUser, authorId)
		if err != nil {
			failedToAccess = append(failedToAccess, gitUser)
			logrus.Error(err)
			continue
		}
		authorName = user.Name
		break
	}

	return authorName, nil
}
