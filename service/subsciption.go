package service

import (
	"gitlab-tg-bot/internal/interfaces"
	"gitlab-tg-bot/service/model"

	"github.com/sirupsen/logrus"
)

const (
	HookTypePushEvents               = "PushEvents"
	HookTypeIssuesEvents             = "IssuesEvents"
	HookTypeConfidentialIssuesEvents = "ConfidentialIssuesEvents"
	HookTypeMergeRequestsEvents      = "MergeRequestsEvents"
	HookTypeTagPushEvents            = "TagPushEvents"
	HookTypeNoteEvents               = "NoteEvents"
	HookTypeJobEvents                = "JobEvents"
	HookTypePipelineEvents           = "PipelineEvents"
	HookTypeWikiPageEvents           = "WikiPageEvents"
)

type SubscriptionService struct {
	TicketProvider interfaces.TicketProvider

	GitlabApi interfaces.GitApiService
}

var _ interfaces.SubscriptionService = (*SubscriptionService)(nil)

func NewSubscription(conf interfaces.Configuration, provider interfaces.ProviderStorage,
	gitlabApi interfaces.GitApiService) *SubscriptionService {
	return &SubscriptionService{
		TicketProvider: provider.Ticket(),

		GitlabApi: gitlabApi,
	}
}

func (s *SubscriptionService) Subscribe(gitlab model.GitlabUser, chatId int64, hookInfo model.Hook) (int32, error) {
	err := s.GitlabApi.AddWebHook(gitlab, hookInfo)
	if err != nil {
		logrus.Error(err)
		return 0, err
	}

	ticket := model.Ticket{
		MaintainerGitlabId: gitlab.Id,
		ChatIds:            []int64{chatId},
		RepositoryId:       hookInfo.RepoId,
		HookTypes: map[string]interface{}{
			HookTypePushEvents:               hookInfo.PushEvents,
			HookTypeIssuesEvents:             hookInfo.IssuesEvents,
			HookTypeConfidentialIssuesEvents: hookInfo.ConfidentialIssuesEvents,
			HookTypeMergeRequestsEvents:      hookInfo.MergeRequestsEvents,
			HookTypeTagPushEvents:            hookInfo.TagPushEvents,
			HookTypeNoteEvents:               hookInfo.NoteEvents,
			HookTypeJobEvents:                hookInfo.JobEvents,
			HookTypePipelineEvents:           hookInfo.PipelineEvents,
			HookTypeWikiPageEvents:           hookInfo.WikiPageEvents,
		},
	}
	ticketId, err := s.TicketProvider.AddTicket(ticket)
	if err != nil {
		return 0, err
	}

	return ticketId, nil
}

func (s *SubscriptionService) GetRepositories(user model.GitlabUser) ([]model.Repository, error) {
	return s.GitlabApi.GetRepositories(user)
}
