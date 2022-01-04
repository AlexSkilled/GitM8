package service

import (
	"gitlab-tg-bot/internal/interfaces"
	"gitlab-tg-bot/service/model"

	"github.com/sirupsen/logrus"
)

type SubscriptionService struct {
	GitlabApi interfaces.GitApiService

	ticket   interfaces.TicketProvider
	patterns interfaces.MessagePatternProvider
}

var _ interfaces.SubscriptionService = (*SubscriptionService)(nil)

func NewSubscription(provider interfaces.ProviderStorage,
	gitlabApi interfaces.GitApiService) *SubscriptionService {
	return &SubscriptionService{
		ticket:    provider.Ticket(),
		patterns:  provider.MessagePattern(),
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
