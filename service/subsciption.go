package service

import (
	"fmt"
	"gitlab-tg-bot/internal/interfaces"
	"gitlab-tg-bot/service/model"

	"github.com/sirupsen/logrus"
)

type SubscriptionService struct {
	GitlabApi interfaces.GitApiService

	ticket interfaces.TicketProvider
}

var _ interfaces.SubscriptionService = (*SubscriptionService)(nil)

func NewSubscription(conf interfaces.Configuration, provider interfaces.ProviderStorage,
	gitlabApi interfaces.GitApiService) *SubscriptionService {
	return &SubscriptionService{
		ticket: provider.Ticket(),

		GitlabApi: gitlabApi,
	}
}

func (s *SubscriptionService) Subscribe(gitlab model.GitUser, chatId int64, hookInfo model.Hook) (int32, error) {
	err := s.GitlabApi.AddWebHook(gitlab, hookInfo)
	if err != nil {
		logrus.Error(err)
		return 0, err
	}

	ticket := model.Ticket{
		MaintainerGitlabId: gitlab.Id,
		ChatIds:            []int64{chatId},
		RepositoryId:       hookInfo.RepoId,
		HookTypes: map[model.GitHookType]interface{}{
			model.HookTypePush:               hookInfo.PushEvents,
			model.HookTypeIssues:             hookInfo.IssuesEvents,
			model.HookTypeConfidentialIssues: hookInfo.ConfidentialIssuesEvents,
			model.HookTypeMergeRequests:      hookInfo.MergeRequestsEvents,
			model.HookTypeTagPush:            hookInfo.TagPushEvents,
			model.HookTypeNote:               hookInfo.NoteEvents,
			model.HookTypeJob:                hookInfo.JobEvents,
			model.HookTypePipeline:           hookInfo.PipelineEvents,
			model.HookTypeWikiPage:           hookInfo.WikiPageEvents,
		},
	}
	ticketId, err := s.ticket.AddTicket(ticket)
	if err != nil {
		return 0, err
	}

	return ticketId, nil
}

func (s *SubscriptionService) GetRepositories(user model.GitUser) ([]model.Repository, error) {
	return s.GitlabApi.GetRepositories(user)
}

func (s *SubscriptionService) ProcessMessage(event model.GitEvent) ([]model.OutputMessage, error) {
	tickets, err := s.ticket.GetTicketsToSend(event.ProjectId, event.HookType)
	fmt.Println(tickets)
	if err != nil {
		return nil, err
	}

	var messageText string

	switch event.HookType {
	case model.HookTypeMergeRequests:
		messageText = ""
	}

	messages := make([]model.OutputMessage, 0)
	messages[0].Msg = messageText

	return messages, nil
}
